package decorator

import "fmt"

type ConcreteComponent struct {
}

func (comp *ConcreteComponent) Operation() {

	fmt.Println("相关组件执行操作...")

}
