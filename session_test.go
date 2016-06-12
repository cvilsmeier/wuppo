package wuppo

import (
	"testing"
	"time"
)

func TestNewMemStore(t *testing.T) {
	store := NewMemStore()
	if len(store.sessions) != 0 {
		t.Errorf("must have no sessions")
	}
}

func TestExpireSession(t *testing.T) {
	store := NewMemStore()
	store.sessions["sid1"] = &session{
		sid:    "sid1",
		atime:  time.Now().Add(-40 * time.Minute),
		values: make(map[string]string),
	}
	store.sessions["sid2"] = &session{
		sid:    "sid2",
		atime:  time.Now().Add(-30 * time.Minute),
		values: make(map[string]string),
	}
	store.sessions["sid3"] = &session{
		sid:    "sid3",
		atime:  time.Now().Add(-20 * time.Minute),
		values: make(map[string]string),
	}
	store.ExpireSessions()
	if len(store.sessions) != 2 {
		t.Errorf("expected 2 sessions to survive")
	}
	if _, ok := store.sessions["sid2"]; !ok {
		t.Errorf("sid2 not found")
	}
	if _, ok := store.sessions["sid3"]; !ok {
		t.Errorf("sid3 not found")
	}
}

func TestTouchSession(t *testing.T) {
	store := NewMemStore()
	t1 := time.Now().Add(-10 * time.Minute)
	s := &session{
		sid:    "aaa",
		atime:  t1,
		values: make(map[string]string),
	}
	store.sessions["aaa"] = s
	store.TouchSession("aaa")
	if !s.atime.After(t1) {
		t.Errorf("touch did not reset atime")
	}
}

func TestPutValue(t *testing.T) {
	store := NewMemStore()
	// put very first value -> must create new session
	sid1 := store.PutValue("ff0a", "name", "chris")
	if len(sid1) != 32 {
		t.Errorf("expected len(sid1)==32")
	}
	// put value -> must return the session id
	sid2 := store.PutValue(sid1, "name", "chris")
	if sid2 != sid1 {
		t.Errorf("expected sid2 == sid1")
	}
	// kill session
	store.KillSession(sid2)
	// put value -> must return a new session id
	sid3 := store.PutValue(sid2, "name", "chris")
	if sid3 == sid1 || sid3 == sid2 {
		t.Errorf("expected new sid3")
	}
}

func TestGetValue(t *testing.T) {
	store := NewMemStore()
	sid := store.PutValue("", "name", "chris")
	if store.GetValue(sid, "name") != "chris" {
		t.Errorf("wanted chris")
	}
	store.PutValue(sid, "name", "bi")
	if store.GetValue(sid, "name") != "bi" {
		t.Errorf("wanted bi")
	}
	store.KillSession(sid)
	if store.GetValue(sid, "name") != "" {
		t.Errorf("wanted empty string")
	}
}
