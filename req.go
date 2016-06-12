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

	// SetSessionValue puts a named value in the session associated
	// with this request. If the request has no valid session, it creates
	// one. If the keyed value already exists, it is replaced.
	SetSessionValue(name string, value string)

	// SessionValue returns the named session value, or the empty string
	// if the key was not found or this reqeust has no valid session.
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

// reqImpl is the default implementation of Req. It's based on a
// http.Request and a http.ResponseWriter
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

// ReqStub implements Req but can be created and manipulated programmatically.
// Use it for unit-testing your controller logic.
type ReqStub struct {
	MethodString string
	PathString   string
	FormValueMap map[string]string
	ModelMap     map[string]interface{}
	SessionMap   map[string]string
	Html         string
	Template     string
	Redirect     string
	Status       int
}

func NewReqStub(method string, path string) *ReqStub {
	req := ReqStub{
		MethodString: method,
		PathString:   path,
		FormValueMap: make(map[string]string),
		ModelMap:     make(map[string]interface{}),
		SessionMap:   make(map[string]string),
	}
	return &req
}

func (req *ReqStub) Method() string {
	return req.MethodString
}

func (req *ReqStub) IsGet() bool {
	return req.MethodString == "GET"
}

func (req *ReqStub) IsPost() bool {
	return req.MethodString == "POST"
}

func (req *ReqStub) Path() string {
	return req.PathString
}

func (req *ReqStub) FormValue(name string) string {
	return req.FormValueMap[name]
}

func (req *ReqStub) SetModelValue(name string, value interface{}) {
	req.ModelMap[name] = value
}

func (req *ReqStub) ModelValue(name string) interface{} {
	return req.ModelMap[name]
}

func (req *ReqStub) Model() map[string]interface{} {
	return req.ModelMap
}

func (req *ReqStub) SetSessionValue(name string, value string) {
	req.SessionMap[name] = value
}

func (req *ReqStub) SessionValue(name string) string {
	return req.SessionMap[name]
}

func (req *ReqStub) KillSession() {
	req.SessionMap = nil
}

func (req *ReqStub) SetHtml(html string) {
	req.Html = html
}

func (req *ReqStub) SetTemplate(template string) {
	req.Template = template
}

func (req *ReqStub) SetRedirect(url string) {
	req.Redirect = url
}

func (req *ReqStub) SetStatus(code int) {
	req.Status = code
}
