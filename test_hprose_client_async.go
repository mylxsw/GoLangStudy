package main

import (
	"fmt"
	"github.com/hprose/hprose-go/hprose"
)

type clientStub struct {
	Sum func(...int) (<-chan int, <-chan error)
}

func main() {
	client := hprose.NewClient("http://127.0.0.1:8080/")
	var ro *clientStub
	client.UseService(&ro)

	sum, err := ro.Sum(1, 2, 3, 4)
	fmt.Println(<-sum, <-err)

	sum, err = ro.Sum(1)
	fmt.Println(<-sum, <-err)
}
