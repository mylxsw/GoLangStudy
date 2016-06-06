package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Hello, world! The time is ", time.Now())
	fmt.Println("My favorite number is ", rand.Intn(10))

	array := [...]string{
		"Hello", "World", "Ni", "Wo", "Ta",
	}

	fmt.Println(array)

}
