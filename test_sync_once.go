package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once
	onceBody := func() {
		fmt.Println("Only once")
	}

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(i int) {
			once.Do(onceBody)
			fmt.Printf("执行 %d.\n", i)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
