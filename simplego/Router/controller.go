package Router

import (
	"html/template"
	"net/http"
	"net/url"
)

type Context struct {
	R      *http.Request
	W      http.ResponseWriter
	Params map[string]string
}

type Controller struct {
	Data         map[interface{}]interface{}
	Ctx          *Context
	TplName      string
	IsNeedRender bool
}

type ControllerInterface interface {
	Init(ctx *Context)
	Get()
	Post()
	Render()
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

func (c *Controller) Render() {
	if c.IsNeedRender {
		if len(c.TplName) > 0 {
			if t, error := template.ParseFiles("views/" + c.TplName); error == nil {
				//c.Ctx.W.Header().Set("Content-Type", "text/html")
				t.Execute(c.Ctx.W, c.Data)
			}
		}
	}
}

func (c *Controller) Input() url.Values {
	if c.Ctx.R.Form == nil {
		c.Ctx.R.ParseForm()
	}
	return c.Ctx.R.Form
}
