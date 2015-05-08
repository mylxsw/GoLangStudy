package main

import (
	strategy "aicode.cc/design/strategy"
	"fmt"
)

func main() {
	fmt.Println("策略模式")

	context := &strategy.Context{}
	context.SetStrategy(&strategy.ConcreteStrategyA{})
	context.Algorithm()

	context.SetStrategy(&strategy.ConcreteStrategyB{})
	context.Algorithm()
}
