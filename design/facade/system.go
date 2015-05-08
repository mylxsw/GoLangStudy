package facade

import "fmt"

type SystemA struct{}

func (sys *SystemA) operationA() {
	fmt.Println("SystemA operationA")
}

type SystemB struct{}

func (sys *SystemB) operationB() {
	fmt.Println("SystemB operationB")
}

type SystemC struct{}

func (sys *SystemC) operationC() {
	fmt.Println("SystemC operationC")
}
