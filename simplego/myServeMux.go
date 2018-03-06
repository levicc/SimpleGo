package simplego

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

var (
	StaticPathMap map[string]string
)

type MuxFunc func(http.ResponseWriter, *http.Request)

type MyMux struct {
	mu     sync.RWMutex
	routes []*controllerInfo
}

type controllerInfo struct {
	params         map[int]string
	regex          *regexp.Regexp
	controllerType reflect.Type
	fmap           map[string]MuxFunc
	pattern        string
}

func init() {
	StaticPathMap = make(map[string]string)
}

func (mux *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestPath := r.URL.Path

	for url, path := range StaticPathMap {
		if strings.HasPrefix(requestPath, url) {
			file := path + requestPath[len(url):]
			http.ServeFile(w, r, file)
			return
		}
	}

	conInfo, params := mux.configControllerInfo(requestPath)
	if conInfo.controllerType != nil {
		controller := reflect.New(conInfo.controllerType)
		method := controller.MethodByName("Init")

		ctx := &Context{W: w, R: r, Params: params}
		method.Call([]reflect.Value{reflect.ValueOf(ctx)})

		isMethodMatch := true
		if r.Method == "GET" {
			controller.MethodByName("Get").Call(nil)
		} else if r.Method == "POST" {
			controller.MethodByName("Post").Call(nil)
		} else {
			isMethodMatch = false
			http.Error(w, "Method Not Match", 405)
		}

		if isMethodMatch {
			controller.MethodByName("Render").Call(nil)
		}

	} else if conInfo.fmap != nil {
		if f, isExist := conInfo.fmap[r.Method]; isExist {
			f(w, r)
		}
	} else {
		http.NotFound(w, r)
	}
}

func (mux *MyMux) Router(pattern string, controller ControllerInterface) {
	mux.addControllerInfo(pattern, controller, nil)
}

func (mux *MyMux) Get(pattern string, f MuxFunc) {
	mux.addMethods(pattern, f, "GET")
}

func (mux *MyMux) Post(pattern string, f MuxFunc) {
	mux.addMethods(pattern, f, "POST")
}

func (mux *MyMux) addMethods(pattern string, f MuxFunc, methodName string) {
	if conInfor, isExist := mux.isControllerExistWithPattern(pattern); isExist {
		fmap := conInfor.fmap
		fmap[methodName] = f
	} else {
		fmap := make(map[string]MuxFunc, 1)
		fmap[methodName] = f
		mux.addControllerInfo(pattern, nil, fmap)
	}
}

func (mux *MyMux) addControllerInfo(pattern string, controller ControllerInterface, fmap map[string]MuxFunc) {
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

	var conInfor *controllerInfo
	if controller != nil {
		conInfor = &controllerInfo{params: params, regex: regex, pattern: pattern, controllerType: reflect.Indirect(reflect.ValueOf(controller)).Type(), fmap: nil}
	} else {
		conInfor = &controllerInfo{params: params, regex: regex, pattern: pattern, controllerType: nil, fmap: fmap}
	}

	mux.routes = append(mux.routes, conInfor)
}

func (mux *MyMux) configControllerInfo(requestPath string) (*controllerInfo, map[string]string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

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

func (mux *MyMux) isControllerExistWithPattern(pattern string) (*controllerInfo, bool) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	for _, conInfor := range mux.routes {
		if conInfor.pattern != pattern {
			continue
		}

		return conInfor, true
	}

	return &controllerInfo{}, false
}
