package main

import (
	"fmt"
	"github.com/cvilsmeier/wuppo"
	"net/http"
)

func main() {
	// init a wuppo handler
	handler := wuppo.DefaultHandler(func(req *wuppo.Req) {
		req.Html = fmt.Sprintf("<html>Hello %s</html>", req.Path)
	})
	// register it with go's http
	http.Handle("/", handler)
	// start on port 8080
    fmt.Printf("now goto http://localhost:8080\n")
	err := http.ListenAndServe(":8080", nil)
    fmt.Printf("%s\n", err)
}
