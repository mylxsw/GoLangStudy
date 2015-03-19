package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var usage = `Usage: goab [options...] <url>

Your CPU numbers: %d


`

var c = flag.Int("c", 50, "并发数目")
var n = flag.Int("n", 200, "请求数目")

flag.Usage = func() {
  usage()
}

func usageAndExit(message string) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}

	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exist(1)
}
