package main

import (
	"fmt"
	"github.com/hprose/hprose-go/hprose"
)

type ClientStub struct {
	Sum   func(...int) int
	Hello func(name string) string
	Swap  func(int, int) (int, int)
}

func main() {
	client := hprose.NewClient("http://127.0.0.1:8181/")
	var ro *ClientStub
	client.UseService(&ro)
	fmt.Println(ro.Sum(1, 2, 3, 5, 67))
	fmt.Println(ro.Hello("hprose"))
	fmt.Println(ro.Swap(3, 4))
}
