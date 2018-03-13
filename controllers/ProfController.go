package controllers

import (
	"SimpleGo/simplego"
	"fmt"
	"net/http/pprof"
)

type ProfController struct {
	simplego.Controller
}

func (this *ProfController) Get() {
	fmt.Println(this.Ctx.Params)
	switch this.Ctx.Params[":pp"] {
	default:
		pprof.Index(this.Ctx.W, this.Ctx.R)
	case "":
		pprof.Index(this.Ctx.W, this.Ctx.R)
	case "cmdline":
		pprof.Cmdline(this.Ctx.W, this.Ctx.R)
	case "profile":
		pprof.Profile(this.Ctx.W, this.Ctx.R)
	case "symbol":
		pprof.Symbol(this.Ctx.W, this.Ctx.R)
	}
	this.Ctx.W.WriteHeader(200)
}
