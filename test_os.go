package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Getwd())

    if _, err := os.Stat("./test_array.go"); err != nil {
        fmt.Println(err)
    }
}
