package main

import (
	"fmt"
	"path"
)

func main() {
	paths := []string{
		"a/b",
		"a//b",
		"a/c/.",
		"a/c/b/..",
		"/../a/c",
		"/../a/b/../././/c",
	}

	for _, p := range paths {
		fmt.Printf("Clean(%q) = %q\n", p, path.Clean(p))
	}

	fmt.Println("Path=", path.Join("/usr/local", "env"))
	fmt.Print("org=" + "/usr/loca/env, res=")
	fmt.Println(path.Split("/usr/local/env"))
}
