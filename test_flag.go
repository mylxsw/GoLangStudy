package main

import (
	"flag"
	"fmt"
	"os"
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var ip = flag.Int("flagname", 1234, "help message")
	fmt.Println(ip)
	Usage()
}
