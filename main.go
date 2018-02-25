package main

import (
	"net/http"
)

func helloworld(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/", helloworld)
	http.ListenAndServe("127.0.0.1:9090", nil)

}
