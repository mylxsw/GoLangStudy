package main

import (
	ada "aicode.cc/design/adapter"
	"fmt"
)

func main() {
	fmt.Println("适配器模式")

	ada.ClientFunc(ada.NewAdapter())
}
