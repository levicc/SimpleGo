package main

import (
	"net/http"
	"sync"
)

type muxFunc func(http.ResponseWriter, *http.Request)

type myMux struct {
	mu sync.RWMutex
	m  map[string]muxEntry
}

type muxEntry struct {
	h       muxFunc
	pattern string
}

func (mux *myMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f := mux.configMuxFunc(r)
	f(w, r)
}

func (mux *myMux) configMuxFunc(r *http.Request) muxFunc {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	v, ok := mux.m[r.URL.Path]

	if ok {
		return v.h
	}

	return nil
}

func (mux *myMux) addMuxFunc(pattern string, handle muxFunc) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}

	entry := muxEntry{h: handle, pattern: pattern}
	mux.m[pattern] = entry
}
