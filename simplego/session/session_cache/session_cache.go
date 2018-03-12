package session_cache

import (
	"SimpleGo/simplego/session"
	"sync"
	"time"

	"github.com/muesli/cache2go"
)

var cacheTable = cache2go.Cache("SessionCache")
var cacheProvide = &SessionProvide{}

func init() {
	session.Register("cache", cacheProvide)
}

type SessionStroe struct {
	sid   string
	value map[interface{}]interface{}
	lock  sync.Mutex
}

func (store *SessionStroe) Set(key, value interface{}) {
	store.lock.Lock()
	defer store.lock.Unlock()

	store.value[key] = value
}

func (store *SessionStroe) Get(key interface{}) interface{} {
	store.lock.Lock()
	defer store.lock.Unlock()

	if v, ok := store.value[key]; ok {
		return v
	}
	return nil
}

func (store *SessionStroe) Delete(key interface{}) {
	store.lock.Lock()
	defer store.lock.Unlock()

	delete(store.value, key)
}

func (store *SessionStroe) SessionId() string {
	return store.sid
}

type SessionProvide struct {
	lock        sync.Mutex
	maxLifeTime int64
}

func (provide *SessionProvide) SessionInit(sid string, maxLifeTime int64) (session.Session, error) {
	provide.lock.Lock()
	defer provide.lock.Unlock()

	value := make(map[interface{}]interface{}, 0)
	session := &SessionStroe{sid: sid, value: value}

	provide.maxLifeTime = maxLifeTime
	cacheTable.Add(sid, time.Duration(provide.maxLifeTime), session)

	return session, nil
}

func (provide *SessionProvide) SessionRead(sid string) (session.Session, error) {
	provide.lock.Lock()
	defer provide.lock.Unlock()

	if cacheTable.Exists(sid) {
		value, err := cacheTable.Value(sid)
		return value.Data().(*SessionStroe), err
	} else {
		return provide.SessionInit(sid, provide.maxLifeTime)
	}

	return nil, nil
}

func (provide *SessionProvide) SessionDelete(sid string) error {
	provide.lock.Lock()
	defer provide.lock.Unlock()

	if cacheTable.Exists(sid) {
		_, err := cacheTable.Delete(sid)
		return err
	}

	return nil
}

func (provide *SessionProvide) SessionGC(maxLifeTime int64) {

}
