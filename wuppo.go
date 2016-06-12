// Package wuppo provides a web framework for Go, ridiculously simple.
package wuppo

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"
)

// Handler is a net/http/Handler implementation that serves as entry
// point to wuppo.
type Handler struct {
	serve ServeFunc
	store SessionStore
}

// ServeFunc is a callback method that responds to an incoming request.
type ServeFunc func(req Req)

// NewHandler creates a new handler with a ServerFunc function and a
// custom sessionStore.
func NewHandler(serve ServeFunc, sessionStore SessionStore) Handler {
	h := Handler{
		serve: serve,
		store: sessionStore,
	}
	return h
}

// DefaultHandler creates a new handler with a ServerFunc function and a
// in-memory session store.
func DefaultHandler(serve ServeFunc) Handler {
	memstore := NewMemStore()
	h := Handler{
		serve: serve,
		store: memstore,
	}
	return h
}

// ServeHTTP implements the net/http/Handler interface.
// It creates a new Req and sends it to the user-defnied ServeFunc.
func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
	t1 := time.Now()
	handler.store.ExpireSessions()
	req := newReqImpl(w, r, handler.store)
	handler.serve(req)
	if req.html != "" {
		io.WriteString(w, req.html)
	} else if req.template != "" {
		t, err := template.ParseGlob("*.html")
		if err != nil {
			panic(err)
		}
		err = t.ExecuteTemplate(w, req.template, req.model)
		if err != nil {
			panic(err)
		}
	} else if req.redirect != "" {
		http.Redirect(w, r, req.redirect, http.StatusFound)
	} else if req.status != 0 {
		msg := http.StatusText(req.status)
		http.Error(w, msg, req.status)
	} else {
		io.WriteString(w, "no result")
	}
	d := time.Since(t1)
	fmt.Printf("%s %s %s - %f s\n", r.RemoteAddr, r.Method, r.URL.Path, float64(d)/1e9)
}
