package main

import (
	"fmt"
	"time"
	// "runtime"
	// "runtime/debug"
)

func main() {
	// go fmt.Println("Go! Goroutine!")
	// time.Sleep(time.Millisecond)
	// runtime.Gosched()

	// names := []string{"Eric", "Harry", "Robert", "Jim", "Mark"}
	// for _, name := range names {
	// 	go func(name string) {
	// 		fmt.Printf("Hello, %s \n", name)
	// 	}(name)
	// }
	//
	// fmt.Printf("Number of Goroutine: %d\n", runtime.NumGoroutine())
	// fmt.Printf("Number of CPU: %d\n", runtime.NumCPU())
	//
	// debug.FreeOSMemory()
	// runtime.GC()
	// runtime.Gosched()

	ch := make(chan int, 5)
	sign := make(chan byte, 2)

	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
			time.Sleep(1 * time.Second)
		}

		close(ch)
		fmt.Println("The channel is closed. ")
		sign <- 0
	}()

	go func() {
		for {
			e, ok := <-ch
			fmt.Printf("%d (%v)\n", e, ok)

			if !ok {
				break
			}

			time.Sleep(2 * time.Second)
		}

		fmt.Println("Done.")
		sign <- 1
	}()

	<-sign
	<-sign
}
