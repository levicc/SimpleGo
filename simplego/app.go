package simplego

import (
	"SimpleGo/simplego/Router"
	"net/http"
)

var (
	SimpleApp *App
)

func init() {
	SimpleApp = NewApp()
}

type App struct {
	Handlers *Router.MyMux
	Server   *http.Server
}

func NewApp() *App {
	app := &App{Handlers: new(Router.MyMux), Server: &http.Server{}}
	return app
}

func (app *App) Run() {
	http.ListenAndServe("127.0.0.1:9091", SimpleApp.Handlers)
}

func Run() {
	SimpleApp.Run()
}

func SetStaticPath(url string, path string) {
	Router.StaticPathMap[url] = path
}

func Add(pattern string, controller Router.ControllerInterface) {
	SimpleApp.Handlers.Router(pattern, controller)
}

func Get(pattern string, f Router.MuxFunc) {
	SimpleApp.Handlers.Get(pattern, f)
}

func Post(pattern string, f Router.MuxFunc) {
	SimpleApp.Handlers.Post(pattern, f)
}
