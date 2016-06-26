

# wuppo

[![GoDoc](https://godoc.org/github.com/cvilsmeier/wuppo?status.svg)](https://godoc.org/github.com/cvilsmeier/wuppo)

A web framework for Go, ridiculously simple.


## Warning

Wuppo is still under development and not suited for production.


## Usage

A very basic usage is this:


```go
package main

import (
	"fmt"
	"github.com/cvilsmeier/wuppo"
	"log"
	"net/http"
)

func serve(req wuppo.Req) {
	html := fmt.Sprintf("<html>%s %s</html>", req.Method(), req.Path())
	req.SetHTML(html)
}

func main() {
	// register a default wuppo http.Handler
	// default means: store session data in memory and
	// search html templates in current directory
	http.Handle("/", wuppo.NewDefaultHandler(serve))
	// start on port 8080
	fmt.Printf("server is up, now goto http://localhost:8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

See the examples folder for more usage examples.

## Licence

Hell, it's free! Do whatevery you like.

