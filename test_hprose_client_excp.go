package main

import (
	"fmt"
	"github.com/hprose/hprose-go/hprose"
)

type clientStub struct {
	Sum func(...int) int
}

func main() {
	client := hprose.NewClient("http://127.0.0.1:8080/")
	var ro *clientStub
	client.UseService(&ro)

	fmt.Println(ro.Sum(1, 2, 3))
}
