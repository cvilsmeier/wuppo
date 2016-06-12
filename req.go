// Package wuppo provides a web framework for Go, ridiculously simple.
package wuppo

import (
	"net/http"
)

// Req provides information about a HTTP request and stores response data.
type Req interface {

	// Method returns the request method: "GET", "POST", etc..
	Method() string

	// IsGet returns true if the request method is GET, false otherwise.
	IsGet() bool

	// IsPost returns true if the request method is POST, false otherwise.
	IsPost() bool

	// Path returns the URL path of the request.
	Path() string

	// FormValue returns the value of a request parameter.
	FormValue(name string) string

	// SetModelValue
	SetModelValue(key string, value interface{})

	// ModelValue returns a keyed model value.
	ModelValue(key string) interface{}

	// Model returns the model as a map.
	Model() map[string]interface{}

	// SetSessionValue puts a named value in the session associated with this request. If the request has no valid session, it creates one. If the keyed value already exists, it is replaced.
	SetSessionValue(name string, value string)

	// SessionValue returns the named session value, or the empty string if the key was not found or this reqeust has no valid session.
	SessionValue(name string) string

	// KillSession kills the session associated with this request.
	KillSession()

	// SetHtml sets a html reponse.
	SetHtml(html string)

	// SetTemplate sets a template reponse.
	SetTemplate(template string)

	// SetRedirect sets a redirect reponse.
	SetRedirect(url string)

	// SetStatus sets a status reponse.
	SetStatus(code int)
}

type reqImpl struct {
	w        http.ResponseWriter
	r        *http.Request
	store    SessionStore
	sid      string
	model    map[string]interface{}
	html     string
	template string
	redirect string
	status   int
}

func newReqImpl(w http.ResponseWriter, r *http.Request, store SessionStore) *reqImpl {
	sid := ""
	if c, err := r.Cookie("WUPPO_SESSION_ID"); err == nil {
		sid = c.Value
		store.TouchSession(sid)
	}
	req := reqImpl{
		w:     w,
		r:     r,
		store: store,
		sid:   sid,
		model: make(map[string]interface{}),
	}
	return &req
}

func (req *reqImpl) Method() string {
	return req.r.Method
}

func (req *reqImpl) IsGet() bool {
	return req.r.Method == "GET"
}

func (req *reqImpl) IsPost() bool {
	return req.r.Method == "POST"
}

func (req *reqImpl) Path() string {
	return req.r.URL.Path
}

func (req *reqImpl) FormValue(name string) string {
	return req.r.FormValue(name)
}

func (req *reqImpl) SetModelValue(name string, value interface{}) {
	req.model[name] = value
}

func (req *reqImpl) ModelValue(name string) interface{} {
	return req.model[name]
}

func (req *reqImpl) Model() map[string]interface{} {
	return req.model
}

func (req *reqImpl) SetSessionValue(name string, value string) {
	newSid := req.store.PutValue(req.sid, name, value)
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

func (req *reqImpl) SessionValue(name string) string {
	return req.store.GetValue(req.sid, name)
}

func (req *reqImpl) KillSession() {
	req.store.KillSession(req.sid)
}

func (req *reqImpl) SetHtml(html string) {
	req.html = html
}

func (req *reqImpl) SetTemplate(template string) {
	req.template = template
}

func (req *reqImpl) SetRedirect(url string) {
	req.redirect = url
}

func (req *reqImpl) SetStatus(code int) {
	req.status = code
}
