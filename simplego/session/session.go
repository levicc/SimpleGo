package session

import (
	"container/list"
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
type SessionStroe struct {
	sid         string
	value       map[interface{}]interface{}
	timeAccessd time.Time
}

func (store *SessionStroe) Set(key, value interface{}) {
	store.value[key] = value
}

func (store *SessionStroe) Get(key interface{}) interface{} {
	if v, ok := store.value[key]; ok {
		return v
	}
	return nil
}

func (store *SessionStroe) Delete(key interface{}) {
	delete(store.value, key)
}

func (store *SessionStroe) SessionId() string {
	return store.sid
}

type SessionProvide struct {
	lock       sync.Mutex
	sessionMap map[string]*list.Element
	lruList    *list.List
}

func (provide *SessionProvide) SessionInit(sid string) (Session, error) {
	provide.lock.Lock()
	defer provide.lock.Unlock()

	value := make(map[interface{}]interface{}, 0)
	session := &SessionStroe{sid: sid, value: value, timeAccessd: time.Now()}
	elem := provide.lruList.PushFront(session)
	provide.sessionMap[sid] = elem

	return session, nil
}

func (provide *SessionProvide) SessionRead(sid string) (Session, error) {
	provide.lock.Lock()
	defer provide.lock.Unlock()

	if elem, ok := provide.sessionMap[sid]; ok {
		return elem.Value.(*SessionStroe), nil
	} else {
		return provide.SessionInit(sid)
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
		session := elem.Value.(*SessionStroe)
		if time.Now().Unix()-session.timeAccessd.Unix() <= maxLifeTime {
			break
		}
		provide.lruList.Remove(elem)
		delete(provide.sessionMap, session.sid)
	}
}
