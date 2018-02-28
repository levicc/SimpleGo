package main

import (
	"SimpleGo/controllers"
	"SimpleGo/simplego"
	"fmt"
	"net/http"
)

func main() {
	simplego.Add("/", &controllers.LYLoginController{})
	simplego.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hellor world! Get")
	})
	simplego.Add("/aaa/bbb/:id([\\w]+)/:username([1-9]+)", &controllers.LYLoginController{})
	simplego.Run()
}
