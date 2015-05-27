package main

import (
	"fmt"
	"net/http"
)

type MyMux struct{}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		sayHelloName(w, req)
		return
	}
	http.NotFound(w, req)
	return
}

func sayHelloName(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, my route!")
}

func main() {
	mux := &MyMux{}
	http.ListenAndServe(":8080", mux)
}
