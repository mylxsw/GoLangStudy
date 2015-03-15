package main

import (
	"fmt"
	"log"
	"net/http"
)

type Hello struct{}
type String string

func (h String) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, h)
}

func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")
}

func main() {
	http.Handle("/string", String("Hello, world"))
	http.Handle("/hello", Hello{})
	err := http.ListenAndServe("localhost:4000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
