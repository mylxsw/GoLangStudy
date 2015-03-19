package main

import (
	"fmt"
	"github.com/hprose/hprose-go/hprose"
)

type clientStub struct {
	Hello      func(string) string
	AsyncHello func(string) <-chan string `name:"hello"`
}

func main() {
	client := hprose.NewClient("http://127.0.0.1:8080/")
	var ro *clientStub
	client.UseService(&ro)

	fmt.Println(ro.Hello("Synchronmous Invoking"))
	fmt.Println(<-ro.AsyncHello("Asynchronmous Invoking"))
}
