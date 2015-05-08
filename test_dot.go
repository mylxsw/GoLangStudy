package main

import (
	"fmt"
)

func test_func_2(m ...string) {
	fmt.Println(m)
}

func test_func(m ...string) {
	test_func_2(m...)
}

func main() {

	test_func("user", "password", "code")
}
