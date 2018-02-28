package controllers

import (
	"SimpleGo/simplego/Router"
	"fmt"
)

type LYLoginController struct {
	Router.Controller
}

func (c *LYLoginController) Get() {
	fmt.Fprintf(c.Ctx.W, fmt.Sprint(c.Ctx.Params))
}
