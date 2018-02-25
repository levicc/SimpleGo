package main

import (
	"fmt"
	"net/http"
)

func helloworld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hellor world")
}

func main() {
	mux := new(MyMux)
	mux.AddMuxFunc("/", helloworld)
	//http.HandleFunc("/", helloworld)
	http.ListenAndServe("127.0.0.1:9090", mux)

}
