package Router

import (
	"net/http"
)

type Context struct {
	R      *http.Request
	W      http.ResponseWriter
	Params map[string]string
}

type Controller struct {
	Ctx *Context
}

type ControllerInterface interface {
	Init(ctx *Context)
	Get()
	Post()
}

func (c *Controller) Init(ctx *Context) {
	c.Ctx = ctx
}

func (c *Controller) Get() {
	http.Error(c.Ctx.W, "Method Not Allowed", 405)
}

func (c *Controller) Post() {
	http.Error(c.Ctx.W, "Method Not Allowed", 405)
}
