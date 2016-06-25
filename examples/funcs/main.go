package main

import (
	"fmt"
	"github.com/cvilsmeier/wuppo"
	"log"
	"net/http"
)

func reverse(s string) string {
    runes := []rune(s)
    count := len(runes)
    reversed := make([]rune, count)
    for i := 0; i < count ; i++ {
        reversed[i] = runes[count-1-i]
    }
    return string(reversed)
}

func serve(req wuppo.Req) {
	switch req.Path() {
	case "/":
        phrase := req.FormValue("phrase")
        log.Print("phrase ", phrase)
        req.SetModelValue("phrase", phrase)
        req.SetTemplate("index.html")
	default:
		req.SetStatus(http.StatusNotFound)
	}
}

func main() {
	// init wuppo handler
    sessionStore := wuppo.NewMemStore()
    funcmap := map[string]interface{}{
        "reverse": reverse,
    }
    h := wuppo.NewHandler(serve, sessionStore, "*.html", funcmap)
	http.Handle("/", h)
	fmt.Printf("goto http://localhost:8080\n")
	log.Panic(http.ListenAndServe(":8080", nil))
}
