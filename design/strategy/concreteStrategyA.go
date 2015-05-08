package strategy

import "fmt"

type ConcreteStrategyA struct {
}

func (strategy *ConcreteStrategyA) Algorithm() {
	fmt.Println("相关策略A")
}
