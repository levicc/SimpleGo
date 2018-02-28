package controllers

import (
	"SimpleGo/simplego/Router"
	"fmt"
)

type LYLoginController struct {
	Router.Controller
}

func (c *LYLoginController) Get() {
	fmt.Fprintf(c.W, fmt.Sprint(c.Params))
}
