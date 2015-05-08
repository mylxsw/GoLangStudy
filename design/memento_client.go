package main

import (
	memento "aicode.cc/design/memento"
	"fmt"
)

// 备忘录模式

func main() {
	// 发起人对象
	originator := &memento.Originator{}
	originator.SetState1("state #1")
	originator.SetState2("state #2")

	fmt.Println(originator)

	// 管理者对象
	caretaker := &memento.Caretaker{}
	caretaker.SetMemento(originator.CreateMemento())

	originator.SetState1("架构")
	originator.SetState2("平台")

	fmt.Println(originator)

	originator.RestoreMemento(caretaker.GetMemento())

	fmt.Println(originator)

}
