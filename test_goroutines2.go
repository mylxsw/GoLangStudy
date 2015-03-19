package main

import (
	"log"
	"time"
)

func main() {

	for i := 0; i < 10; i++ {
		go func(counter int) {
			time.Sleep(100 * time.Millisecond)
			log.Println(counter)
		}(i)
	}

	time.Sleep(1000 * time.Millisecond)
}
