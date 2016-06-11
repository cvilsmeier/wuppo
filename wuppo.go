// Package wuppo provides a web framework for Go, ridiculously simple.
package wuppo

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"
)

// Req provides information about a HTTP request and stores response data.
type Req struct {
	w        http.ResponseWriter
	r        *http.Request
	store    SessionStore
	sid      string
	Method   string
	Path     string
	Model    map[string]interface{}
	Html     string
	Template string
	Redirect string
	Status   int
}

func newReq(w http.ResponseWriter, r *http.Request, store SessionStore) *Req {
	sid := ""
	if c, err := r.Cookie("WUPPO_SESSION_ID"); err == nil {
		sid = c.Value
		store.TouchSession(sid)
	}
	req := Req{
		w:      w,
		r:      r,
		store:  store,
		sid:    sid,
		Method: r.Method,
		Path:   r.URL.Path,
		Model:  make(map[string]interface{}),
	}
	return &req
}

// Param returns the value of a request parameter.
func (req *Req) Param(name string) string {
	r := req.r
	if r == nil {
		panic("r is nil")
	}
	r.ParseForm()
	return r.FormValue(name)
}

// PutSessionValue sets a keyed session value.
func (req *Req) PutSessionValue(key string, value string) {
	newSid := req.store.PutValue(req.sid, key, value)
	if newSid != req.sid {
		req.sid = newSid
		cookie := http.Cookie{
			Name:     "WUPPO_SESSION_ID",
			Value:    newSid,
			MaxAge:   0,
			HttpOnly: true,
		}
		http.SetCookie(req.w, &cookie)
	}
}

// GetSessionValue returns a keyed session value.
// If not found, it returns the empty string.
func (req *Req) GetSessionValue(key string) string {
	return req.store.GetValue(req.sid, key)
}

// KillSession kills the session associated with this request
func (req *Req) KillSession() {
	req.store.KillSession(req.sid)
}

//
//
// --------------------------------------------------
//
//

// Handler is a net/http/Handler implementation that serves as entry
// point to wuppo.
type Handler struct {
	serve ServeFunc
	store SessionStore
}

// ServeFunc is a callback method that responds to an incoming HTTP request.
type ServeFunc func(req *Req)

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
// It prepares a Req and sends it to the user-defnied ServeFunc.
func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
	t1 := time.Now()
	handler.store.ExpireSessions()
	req := newReq(w, r, handler.store)
	handler.serve(req)
	if req.Html != "" {
		io.WriteString(w, req.Html)
	} else if req.Template != "" {
		t, err := template.ParseGlob("*.html")
		if err != nil {
			panic(err)
		}
		err = t.ExecuteTemplate(w, req.Template, req.Model)
		if err != nil {
			panic(err)
		}
	} else if req.Redirect != "" {
		http.Redirect(w, r, req.Redirect, http.StatusFound)
	} else if req.Status != 0 {
		msg := http.StatusText(req.Status)
		http.Error(w, msg, req.Status)
	} else {
		io.WriteString(w, "no result")
	}
	d := time.Since(t1)
	fmt.Printf("%s %s %s - %f s\n", r.RemoteAddr, r.Method, r.URL.Path, float64(d)/1e9)
}

// eof
