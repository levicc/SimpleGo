package Router

import (
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type MuxFunc func(http.ResponseWriter, *http.Request)

type MyMux struct {
	mu     sync.RWMutex
	routes []*controllerInfo
}

type controllerInfo struct {
	params     map[int]string
	regex      *regexp.Regexp
	controller ControllerInterface
	f          MuxFunc
}

func (mux *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conInfo, params := mux.configControllerInfo(r)

	if conInfo.controller != nil {
		controller := conInfo.controller
		controller.Init(w, r, params)

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

func (mux *MyMux) configControllerInfo(r *http.Request) (*controllerInfo, map[string]string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	requestPath := r.URL.Path
	for _, conInfor := range mux.routes {
		if !conInfor.regex.MatchString(requestPath) {
			continue
		}

		matches := conInfor.regex.FindStringSubmatch(requestPath)
		if len(matches[0]) != len(requestPath) {
			continue
		}

		params := make(map[string]string)
		if len(conInfor.params) == len(matches)-1 {
			for i, match := range matches[1:] {
				params[conInfor.params[i]] = match
			}
		}

		return conInfor, params
	}

	return &controllerInfo{}, nil
}

func (mux *MyMux) Router(pattern string, controller ControllerInterface) {
	mux.addController(pattern, controller, nil)
}

func (mux *MyMux) Get(pattern string, f MuxFunc) {
	mux.addController(pattern, nil, f)
}

func (mux *MyMux) addController(pattern string, controller ControllerInterface, f MuxFunc) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	parts := strings.Split(pattern, "/")

	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			if index := strings.Index(part, "("); index != -1 {
				params[j] = part[:index]
				parts[i] = part[index:]
				j++
			} else {
				panic("strings error")
				return
			}
		}
	}

	pattern = strings.Join(parts, "/")
	regex, regexErr := regexp.Compile(pattern)
	if regexErr != nil {
		panic(regexErr)
		return
	}

	conInfor := &controllerInfo{params: params, regex: regex, controller: controller, f: f}
	mux.routes = append(mux.routes, conInfor)
}
