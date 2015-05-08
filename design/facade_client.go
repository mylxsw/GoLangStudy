package main

import (
	"aicode.cc/design/facade"
	"fmt"
)

func main() {
	fmt.Println("外观模式")

	facade := facade.New()
	facade.Operation()

}
