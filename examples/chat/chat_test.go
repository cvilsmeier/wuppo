package main

import (
	"fmt"
	"github.com/cvilsmeier/wuppo"
	"runtime"
	"testing"
)

func TestGetIndex(t *testing.T) {
	req := wuppo.NewReqStub("GET", "/")
	req.FormValueMap["reason"] = "loggedOut"
	serve(req)
	assert(t, req.ModelMap["reason"] == "loggedOut", "wrong reason", req.ModelMap["reason"])
	assert(t, req.Template == "index.html", "wrong template", req.Template)
}

func TestPostIndex(t *testing.T) {
	req := wuppo.NewReqStub("POST", "/")
	req.FormValueMap["name"] = "CV"
	serve(req)
	assert(t, req.Redirect == "/chat", "wrong redirect", req.Redirect)
}

func TestPostIndexWithEmptyName(t *testing.T) {
	req := wuppo.NewReqStub("POST", "/")
	req.FormValueMap["name"] = " "
	serve(req)
	errors := req.ModelMap["errors"].([]string)
	assert(t, len(errors) == 1, "wrong len errors", len(errors))
	assert(t, errors[0] == "Name must not be empty", "wrong errors[0]", errors[0])
	assert(t, req.Template == "index.html", "wrong template", req.Template)
}

// ... and so on, you get the idea

func assert(t *testing.T, condition bool, args ...interface{}) {
	if !condition {
		_, _, line, _ := runtime.Caller(1)
		loc := fmt.Sprintf("in line %d", line)
		t.Error(loc, args)
	}
}
