package main

import (
	"fmt"
	"net/http"
	"sync"
)

type muxFunc func(http.ResponseWriter, *http.Request)

type MyMux struct {
	mu sync.RWMutex
	m  map[string]muxEntry
}

type muxEntry struct {
	h       muxFunc
	pattern string
}

func (mux *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f := mux.configMuxFunc(r)
	if f != nil {
		f(w, r)
	}
}

func (mux *MyMux) configMuxFunc(r *http.Request) muxFunc {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	fmt.Println(r.URL.Path)
	v, ok := mux.m[r.URL.Path]

	if ok {
		return v.h
	}

	return nil
}

func (mux *MyMux) AddMuxFunc(pattern string, handle muxFunc) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}

	entry := muxEntry{h: handle, pattern: pattern}
	mux.m[pattern] = entry
}
