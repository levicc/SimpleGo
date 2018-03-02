package controllers

import (
	"SimpleGo/simplego/Router"
)

type LYLoginController struct {
	Router.Controller
}

func (c *LYLoginController) Get() {
	c.IsNeedRender = true
	c.TplName = "login.html"
}

func (c *LYLoginController) Post() {
	c.Data = map[interface{}]interface{}{"name": c.Input().Get("name")}
	c.IsNeedRender = true
	c.TplName = "login.html"
}
