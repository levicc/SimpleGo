package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Session interface {
	Set(key, value interface{})
	Get(key interface{}) interface{}
	Delete(key interface{})
	SessionId() string
}

type Provide interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDelete(sid string) error
	SessionGC(maxLifeTime int64)
}

type SessionManager struct {
	cookieName  string
	provide     Provide
	lock        sync.Mutex
	maxLifeTime int64
}

var (
	provideMap = make(map[string]Provide)
)

func NewManager(provideName, cookieName string, maxLifeTime int64) (*SessionManager, error) {
	if provide, ok := provideMap[provideName]; ok {
		return &SessionManager{cookieName: cookieName, provide: provide, maxLifeTime: maxLifeTime, lock: sync.Mutex{}}, nil
	}
	return nil, fmt.Errorf("session: unknown provideName %q", provideName)
}

func Register(provideName string, provide Provide) {
	if _, dup := provideMap[provideName]; dup {
		panic("session: Register called twice for provide " + provideName)
	}
	provideMap[provideName] = provide
}

func (manager *SessionManager) SessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *SessionManager) StartSession(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	cookie, err := r.Cookie(manager.cookieName)

	if err != nil || cookie.Value == "" {
		sid := manager.SessionId()
		session, _ = manager.provide.SessionInit(sid)
		http.SetCookie(w, &http.Cookie{
			Name:  manager.cookieName,
			Value: url.QueryEscape(sid),
		})
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provide.SessionRead(sid)
	}
	return
}

func (manager *SessionManager) DestorySession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	sid, _ := url.QueryUnescape(cookie.Value)

	if err == nil && sid != "" {
		manager.lock.Lock()
		defer manager.lock.Unlock()

		manager.provide.SessionDelete(sid)
		http.SetCookie(w, &http.Cookie{
			Name:  manager.cookieName,
			Value: "",
		})
	}
}

func (manager *SessionManager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	manager.provide.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() {
		manager.GC()
	})
}
