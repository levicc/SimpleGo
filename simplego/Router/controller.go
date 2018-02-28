package Router

import (
	"net/http"
)

type Controller struct {
	R      *http.Request
	W      http.ResponseWriter
	Params map[string]string
}

type ControllerInterface interface {
	Init(w http.ResponseWriter, r *http.Request, params map[string]string)
	Get()
	Post()
}

func (c *Controller) Init(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.W = w
	c.R = r
	c.Params = params
}

func (c *Controller) Get() {
	http.Error(c.W, "Method Not Allowed", 405)
}

func (c *Controller) Post() {
	http.Error(c.W, "Method Not Allowed", 405)
}
