package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("Cpu: %d\n", runtime.NumCPU())
}
