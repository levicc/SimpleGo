package Router

import (
	"net/http"
	"sync"
)

type muxFunc func(http.ResponseWriter, *http.Request)

type MyMux struct {
	mu sync.RWMutex
	m  map[string]*controllerInfo
}

type controllerInfo struct {
	controller ControllerInterface
	f          muxFunc
	pattern    string
}

func (mux *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conInfo := mux.configControllerInfo(r)
	if conInfo.controller != nil {
		controller := conInfo.controller
		controller.Init(w, r)

		if r.Method == "GET" {
			controller.Get()
		} else if r.Method == "POST" {
			controller.Post()
		} else {
			http.Error(w, "Method Not Match", 405)
		}
	} else if conInfo.f != nil {
		conInfo.f(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (mux *MyMux) configControllerInfo(r *http.Request) *controllerInfo {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	v, ok := mux.m[r.URL.Path]

	if ok {
		return v
	}

	return &controllerInfo{}
}

func (mux *MyMux) AddController(pattern string, controller ControllerInterface) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	if mux.m == nil {
		mux.m = make(map[string]*controllerInfo)
	}

	entry := &controllerInfo{controller: controller, pattern: pattern}
	mux.m[pattern] = entry
}

func (mux *MyMux) Get(pattern string, f muxFunc) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	if mux.m == nil {
		mux.m = make(map[string]*controllerInfo)
	}

	entry := &controllerInfo{f: f, pattern: pattern}
	mux.m[pattern] = entry
}
