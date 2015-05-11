package main

import (
	"aicode.cc/design/decorator"
	"fmt"
)

func main() {
	fmt.Println("装饰器模式")

	var comp decorator.Component
	comp = &decorator.ConcreteComponent{}

	decorator := decorator.NewConcreteDecorator(comp)
	decorator.Operation()
}
