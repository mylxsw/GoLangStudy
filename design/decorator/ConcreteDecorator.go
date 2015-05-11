package decorator

import "fmt"

type ConcreteDecorator struct {
	Decorator
}

func NewConcreteDecorator(component Component) (dec *ConcreteDecorator) {
	dec = &ConcreteDecorator{}
	dec.setComponent(component)

	return
}

func (dec *ConcreteDecorator) Operation() {
	dec.securityCheck()
	dec.component.Operation()
}

func (dec *ConcreteDecorator) securityCheck() {
	fmt.Println("执行安全检查")
}
