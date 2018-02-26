package Router

import (
	"net/http"
)

type Controller struct {
	R *http.Request
	W http.ResponseWriter
}

type ControllerInterface interface {
	Init(w http.ResponseWriter, r *http.Request)
	Get()
	Post()
}

func (c *Controller) Init(w http.ResponseWriter, r *http.Request) {
	c.W = w
	c.R = r
}

func (c *Controller) Get() {
	http.Error(c.W, "Method Not Allowed", 405)
}

func (c *Controller) Post() {
	http.Error(c.W, "Method Not Allowed", 405)
}
