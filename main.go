package main

import (
	"SimpleGo/Router"
	"SimpleGo/controllers"
	"fmt"
	"net/http"
)

func main() {
	mux := new(Router.MyMux)
	mux.AddController("/", &controllers.LYLoginController{})
	mux.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hellor world! Get")
	})
	http.ListenAndServe("127.0.0.1:9090", mux)
}
