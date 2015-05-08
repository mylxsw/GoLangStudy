package strategy

import "fmt"

type ConcreteStrategyB struct {
}

func (strategy *ConcreteStrategyB) Algorithm() {
	fmt.Println("相关策略B")
}
