package main

import (
	"SimpleGo/controllers"
	"SimpleGo/models"
	"SimpleGo/simplego"
	_ "SimpleGo/simplego/session/session_memory"
	"fmt"
	_ "net/http/pprof"
)

func main() {
	user, _ := models.GetUserWithId(10)
	fmt.Println(user.Id, user.Username, user.Password)

	simplego.SetStaticPath("/asset", "uploads")

	simplego.Add("/", &controllers.LYLoginController{})
	simplego.Add("/aaa/bbb/:id([\\w]+)/:username([1-9]+)", &controllers.LYLoginController{})
	simplego.Add("/debug/pprof", &controllers.ProfController{})
	simplego.Add("/debug/pprof/:pp([\\w]+)", &controllers.ProfController{})
	//simplego.Run()
}
