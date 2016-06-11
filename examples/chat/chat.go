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

//
// --------------------------------------------------
//

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
// --------------------------------------------------
//

var theDb *Db

func serveIndex(req *wuppo.Req) {
	req.Model["reason"] = req.Param("reason")
	if req.Method == "GET" {
		req.Template = "index.html"
		return
	}
	name := req.Param("name")
	name = strings.TrimSpace(name)
	if name == "" {
		req.Model["errors"] = []string{"Name must not be empty"}
		req.Template = "index.html"
		return
	}
	req.PutSessionValue("name", name)
	req.Redirect = "/chat"
}

func serveLogout(req *wuppo.Req) {
	req.KillSession()
	req.Redirect = "/"
}

func serveChat(req *wuppo.Req) {
	name := req.GetSessionValue("name")
	if name == "" {
		req.Redirect = "/?reason=notLoggedIn"
		return
	}
	req.Model["name"] = name
	if req.Method == "POST" {
		message := req.Param("message")
		message = strings.TrimSpace(message)
		if message != "" {
			db.AddMessage(name + ": " + message)
		}
	}
	req.Model["messages"] = db.GetMessages()
	req.Template = "chat.html"
}

func serveSessions(req *wuppo.Req) {
	infos := sessionStore.GetSessionInfos()
	req.Model["infos"] = infos
	req.Template = "sessions.html"
}

func serve(req *wuppo.Req) {
	switch req.Path {
	case "/":
		serveIndex(req)
	case "/logout":
		serveLogout(req)
	case "/chat":
		serveChat(req)
	case "/sessions":
		serveSessions(req)
	default:
		req.Status = 404
	}
}

var sessionStore wuppo.SessionStore
var handler wuppo.Handler
var db *Db

func main() {
	// init a database that stores our messages
	db = NewDb()
	// init wuppo, store session data in memory
	sessionStore = wuppo.NewMemStore()
	handler = wuppo.NewHandler(serve, sessionStore)
	http.Handle("/", handler)
	// serve static files (favicon)
	http.Handle("/favicon.ico", http.FileServer(http.Dir(".")))
	// for development only: exit if a go file changes
	go watchFiles()
	// start the server on port 8080
	fmt.Printf("chat server is up\n")
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
