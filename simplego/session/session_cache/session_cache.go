package session_cache

import (
	"sync"
)

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
