package simplego

import (
	"SimpleGo/simplego/session"
	"net/http"
)

var (
	SimpleApp     *App
	GlobalSession *session.SessionManager
)

func init() {
	SimpleApp = NewApp()
}

type App struct {
	Handlers *MyMux
	Server   *http.Server
}

func NewApp() *App {
	app := &App{Handlers: new(MyMux), Server: &http.Server{}}
	return app
}

func GetGlobelSession() *session.SessionManager {
	if GlobalSession == nil {
		GlobalSession, _ = session.NewManager("memory", "sessionId", 3600)
		go GlobalSession.GC()
	}
	return GlobalSession
}

func (app *App) Run() {
	http.ListenAndServe("127.0.0.1:9091", SimpleApp.Handlers)
}

func Run() {
	SimpleApp.Run()
}

func SetStaticPath(url string, path string) {
	StaticPathMap[url] = path
}

func Add(pattern string, controller ControllerInterface) {
	SimpleApp.Handlers.Router(pattern, controller)
}

func Get(pattern string, f MuxFunc) {
	SimpleApp.Handlers.Get(pattern, f)
}

func Post(pattern string, f MuxFunc) {
	SimpleApp.Handlers.Post(pattern, f)
}
