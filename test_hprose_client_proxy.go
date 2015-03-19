package main

import (
	"github.com/hprose/hprose-go/hprose"
	"net/http"
)

type proxyStub struct {
	Hello func(string) (string, error)
	Swap  func(int, int) (int, int)
	Sum   func(...int) int
}

func main() {
	client := hprose.NewClient("http://127.0.0.1:8080/")
	var ro *proxyStub
	client.UseService(&ro)
	service := hprose.NewHttpService()
	service.AddMethods(ro)
	http.ListenAndServe(":8181", service)
}
