package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	t := time.Now().UTC()
	rand.Seed(t.UnixNano())

	randFunc := func(status chan bool) {
		fmt.Println(rand.Intn(100))
		status <- true
	}

	var status chan bool = make(chan bool)
	for i := 0; i < 10; i++ {
		go randFunc(status)
	}

	for i := 0; i < 10; i++ {
		<-status
	}
}
