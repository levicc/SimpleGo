package controllers

import (
	"SimpleGo/simplego"
	"fmt"
)

type LYLoginController struct {
	simplego.Controller
}

func (c *LYLoginController) Get() {
	session := c.GetSession()
	userName := session.Get("username")

	if userName == nil || userName == "" {
		c.IsNeedRender = true
		c.TplName = "login.html"
	} else {
		c.IsNeedRender = false
		fmt.Fprintln(c.Ctx.W, fmt.Sprintf("hello %s", userName))
	}
}

func (c *LYLoginController) Post() {
	userName := c.Input().Get("name")

	if userName == "luyang" {
		sess := c.GetSession()
		sess.Set("username", userName)
		c.IsNeedRender = false
		fmt.Fprintln(c.Ctx.W, fmt.Sprintf("hello %s", userName))
	} else {
		c.IsNeedRender = true
		c.TplName = "login.html"
	}
}
