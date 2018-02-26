package controllers

import (
	"SimpleGo/Router"
	"fmt"
)

type LYLoginController struct {
	Router.Controller
}

func (c *LYLoginController) Get() {
	fmt.Fprintf(c.W, "hellor world! Controller")
}
