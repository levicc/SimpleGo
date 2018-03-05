package controllers

import (
	"SimpleGo/simplego/Router"
	"fmt"
)

type LYLoginController struct {
	Router.Controller
}

func (c *LYLoginController) Get() {
	cookies := c.Ctx.R.Cookies()
	for _, cookie := range cookies {
		fmt.Println(cookie.Name, cookie.Value)
	}

	c.IsNeedRender = true
	c.TplName = "login.html"
}

func (c *LYLoginController) Post() {
	c.Data = map[interface{}]interface{}{"name": c.Input().Get("name")}
	c.IsNeedRender = true
	c.TplName = "login.html"
}
