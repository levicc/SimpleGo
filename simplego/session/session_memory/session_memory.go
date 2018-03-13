package session_memory

import (
	"SimpleGo/simplego/session"
	"container/list"
	"sync"
	"time"
)

var memoryProvider = &SessionProvide{lruList: list.New(), sessionMap: make(map[string]*list.Element, 0)}

func init() {
	session.Register("memory", memoryProvider)
}

type SessionStroe struct {
	sid         string
	value       map[interface{}]interface{}
	timeAccessd time.Time
	lock        sync.Mutex
}

func (store *SessionStroe) Set(key, value interface{}) {
	store.lock.Lock()
	defer store.lock.Unlock()

	memoryProvider.SessionUpdate(store.sid)
	store.value[key] = value
}

func (store *SessionStroe) Get(key interface{}) interface{} {
	store.lock.Lock()
	defer store.lock.Unlock()

	memoryProvider.SessionUpdate(store.sid)
	if v, ok := store.value[key]; ok {
		return v
	}
	return nil
}

func (store *SessionStroe) Delete(key interface{}) {
	store.lock.Lock()
	defer store.lock.Unlock()

	memoryProvider.SessionUpdate(store.sid)
	delete(store.value, key)
}

func (store *SessionStroe) SessionId() string {
	return store.sid
}

type SessionProvide struct {
	lock        sync.RWMutex
	sessionMap  map[string]*list.Element
	lruList     *list.List
	maxLifeTime int64
}

func (provide *SessionProvide) SessionInit(sid string, maxLifeTime int64) (session.Session, error) {
	provide.lock.Lock()
	defer provide.lock.Unlock()

	value := make(map[interface{}]interface{}, 0)
	session := &SessionStroe{sid: sid, value: value, timeAccessd: time.Now()}
	elem := provide.lruList.PushFront(session)
	provide.sessionMap[sid] = elem
	provide.maxLifeTime = maxLifeTime

	return session, nil
}

func (provide *SessionProvide) SessionRead(sid string) (session.Session, error) {
	// provide.lock.Lock()
	// defer provide.lock.Unlock()

	if elem, ok := provide.sessionMap[sid]; ok {
		return elem.Value.(*SessionStroe), nil
	} else {
		return provide.SessionInit(sid, provide.maxLifeTime)
	}

	return nil, nil
}

func (provide *SessionProvide) SessionDelete(sid string) error {
	provide.lock.Lock()
	defer provide.lock.Unlock()

	if elem, ok := provide.sessionMap[sid]; ok {
		delete(provide.sessionMap, sid)
		provide.lruList.Remove(elem)
	}
	return nil
}

func (provide *SessionProvide) SessionUpdate(sid string) {
	provide.lock.Lock()
	defer provide.lock.Unlock()

	if elem, ok := provide.sessionMap[sid]; ok {
		elem.Value.(*SessionStroe).timeAccessd = time.Now()
		provide.lruList.MoveToFront(elem)
	}
}

func (provide *SessionProvide) SessionGC(maxLifeTime int64) {
	provide.lock.Lock()
	defer provide.lock.Unlock()

	for {
		elem := provide.lruList.Back()

		if elem == nil {
			break
		}

		session := elem.Value.(*SessionStroe)
		if time.Now().Unix()-session.timeAccessd.Unix() <= maxLifeTime {
			break
		}
		provide.lruList.Remove(elem)
		delete(provide.sessionMap, session.sid)
	}
}
