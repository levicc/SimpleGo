package session_cache

import (
	"SimpleGo/simplego/session"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionCache(t *testing.T) {
	globalSession, err := session.NewManager("cache", "sessionId", 3600)

	if err != nil || globalSession == nil {
		t.Error("cacheSession未添加")
	}

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	session := globalSession.StartSession(w, r)

	if session == nil {
		t.Error("未能新建session")
	}

	session.Set("username", "luyang")

	if session.Get("username") != "luyang" {
		t.Error("get Error")
	}
}
