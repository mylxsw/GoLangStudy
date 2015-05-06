package main

import (
	memento "aicode.cc/design/memento"
	"fmt"
)

func main() {
	originator := &memento.Originator{}
	originator.SetState1("state #1")
	originator.SetState2("state #2")
	originator.SetState3("state #3")

	fmt.Println(originator)

	caretaker := &memento.Caretaker{}
	caretaker.SetMemento(*originator.CreateMemento())

	originator.SetState1("架构")
	originator.SetState2("平台")
	originator.SetState3("元素")

	fmt.Println(originator)

	originator.RestoreMemento(caretaker.GetMemento())

	fmt.Println(originator)

}
