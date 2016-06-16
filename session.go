package wuppo

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// A SessionStore is used to manage HTTP sessions.
type SessionStore interface {
	// ExpireSessions expires old sessions. A session is old if it was
	// not accessed within the last 30 minutes.
	ExpireSessions()

	// TouchSession sets the atime (last access time) of a session to the
	// current time, much like the unix 'touch' command does with files.
	TouchSession(sid string)

	// KillSession removes a session.
	KillSession(sid string)

	// PutValue puts a value into a session and returns the session id.
	// If the session with the incoming session id was not found,
	// PutValue creates a new session and returns the new session id.
	PutValue(sid string, key string, value interface{}) string

	// GetValue returns a session value. If the session or the key does not
	// exist, it returns nil.
	GetValue(sid string, key string) interface{}

	// GetSessionInfos returns map of maps containg all sessions with
	// their key/value pairs.
	GetSessionInfos() map[string]map[string]interface{}
}

// MemStore is a SessionStore that stores HTTP session data in memory.
// If the process ends, all session data will be lost.
type MemStore struct {
	mx       sync.Mutex
	sessions map[string]*session
}

// NewMemStore creates a new MemStore.
func NewMemStore() *MemStore {
	st := &MemStore{
		sessions: make(map[string]*session),
	}
	return st
}

// ExpireSessions expires old sessions. A session is old if it was
// not accessed within the last 30 minutes.
func (st *MemStore) ExpireSessions() {
	st.mx.Lock()
	defer st.mx.Unlock()
	// expire old sessions
	for sid := range st.sessions {
		s := st.sessions[sid]
		if time.Since(s.atime).Minutes() > 30 {
			// fmt.Printf("session %s expired\n", sid)
			delete(st.sessions, sid)
		}
	}
}

// TouchSession sets the atime (last access time) of a session to the
// current time, much like the unix 'touch' command does with files.
func (st *MemStore) TouchSession(sid string) {
	st.mx.Lock()
	defer st.mx.Unlock()
	s := st.sessions[sid]
	if s != nil {
		s.atime = time.Now()
	}
}

// KillSession removes a session.
func (st *MemStore) KillSession(sid string) {
	st.mx.Lock()
	defer st.mx.Unlock()
	delete(st.sessions, sid)
}

// PutValue puts a value into a session and returns the session id.
// If the session with the incoming session id was not found,
// PutValue creates a new session and returns the new session id.
func (st *MemStore) PutValue(sid string, key string, value interface{}) string {
	st.mx.Lock()
	defer st.mx.Unlock()
	s := st.sessions[sid]
	if s == nil {
		// fmt.Printf("session %s not found\n", sid)
		buf := make([]byte, 16)
		if _, err := rand.Read(buf); err != nil {
			panic(err)
		}
		s = &session{
			sid:    hex.EncodeToString(buf),
			atime:  time.Now(),
			values: make(map[string]interface{}),
		}
		st.sessions[s.sid] = s
		// fmt.Printf("created new session %s\n", s.sid)
	}
	s.values[key] = value
	return s.sid
}

// GetValue returns a session value. If the session or the key does not
// exist, it returns nil.
func (st *MemStore) GetValue(sid string, key string) interface{} {
	st.mx.Lock()
	defer st.mx.Unlock()
	s := st.sessions[sid]
	if s == nil {
		return nil
	}
	return s.values[key]
}

// GetSessionInfos returns map of maps containg all sessions with
// their key/value pairs.
func (st *MemStore) GetSessionInfos() map[string]map[string]interface{} {
	st.mx.Lock()
	defer st.mx.Unlock()
	infos := make(map[string]map[string]interface{})
	for sid := range st.sessions {
		s := st.sessions[sid]
		infos[sid] = make(map[string]interface{})
		infos[sid]["_sid"] = sid
		infos[sid]["_atime"] = s.atime.String()
		for key := range s.values {
			infos[sid][key] = s.values[key]
		}
	}
	return infos
}

type session struct {
	sid    string
	atime  time.Time
	values map[string]interface{}
}

// eof
