package main

import (
	"fmt"
	"github.com/cvilsmeier/wuppo"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Db is a simple in-memory database for the wuppo chat sample
type Db struct {
	mx       sync.Mutex
	messages []string
}

func NewDb() *Db {
	return &Db{
		messages: make([]string, 0),
	}
}

func (db *Db) AddMessage(message string) {
	db.mx.Lock()
	defer db.mx.Unlock()
	db.messages = append([]string{message}, db.messages...)
}

func (db *Db) GetMessages() []string {
	db.mx.Lock()
	defer db.mx.Unlock()
	all := make([]string, len(db.messages))
	copy(all, db.messages)
	return all
}

//
// serve() methods
//

func serveIndex(req wuppo.Req) {
	reason := req.FormValue("reason")
	req.SetModelValue("reason", reason)
	if req.IsGet() {
		req.SetTemplate("index.html")
		return
	}
	name := req.FormValue("name")
	name = strings.TrimSpace(name)
	if name == "" {
		req.SetModelValue("errors", []string{"Name must not be empty"})
		req.SetTemplate("index.html")
		return
	}
	req.SetSessionValue("name", name)
	req.SetRedirect("/chat")
}

func serveLogout(req wuppo.Req) {
	req.KillSession()
	req.SetRedirect("/")
}

func serveChat(req wuppo.Req) {
	name := req.SessionValue("name")
	if name == "" {
		req.SetRedirect("/?reason=notLoggedIn")
		return
	}
	req.SetModelValue("name", name)
	if req.IsPost() {
		message := req.FormValue("message")
		message = strings.TrimSpace(message)
		if message != "" {
			theDb.AddMessage(name + ": " + message)
		}
	}
	req.SetModelValue("messages", theDb.GetMessages())
	req.SetTemplate("chat.html")
}

func serveSessions(req wuppo.Req) {
	infos := theSessionStore.GetSessionInfos()
	req.SetModelValue("infos", infos)
	req.SetTemplate("sessions.html")
}

func serve(req wuppo.Req) {
	switch req.Path() {
	case "/":
		serveIndex(req)
	case "/logout":
		serveLogout(req)
	case "/chat":
		serveChat(req)
	case "/sessions":
		serveSessions(req)
	default:
		req.SetStatus(http.StatusNotFound)
	}
}

//
// globals
//
var theSessionStore = wuppo.NewMemStore()
var theDb = NewDb()

func main() {
	// init wuppo handler
	http.Handle("/", wuppo.NewHandler(serve, theSessionStore, "*.html"))
	// serve static files (favicon)
	http.Handle("/favicon.ico", http.FileServer(http.Dir(".")))
	// for development only: exit if a go file changes
	go watchFiles()
	// start the server on port 8080
	fmt.Printf("chat server is up, goto http://localhost:8080\n")
	log.Panic(http.ListenAndServe(":8080", nil))
}

//
// file watchdog
//

func watchFiles() {
	// stop if a *.go file changes
	lastCheck := ""
	for {
		check := ""
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".go") {
				check = check + info.ModTime().String() + "_"
			}
			return nil
		})
		if lastCheck != "" && lastCheck != check {
			fmt.Printf("some watched file changed, game over\n")
			os.Exit(0)
		}
		lastCheck = check
		time.Sleep(500 * time.Millisecond)
	}
}

// eof
