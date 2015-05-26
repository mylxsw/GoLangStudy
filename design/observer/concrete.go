package observer

import (
	"container/list"
	"fmt"
)

type ConcreteObserver struct {
	Observer
}

func (observer *ConcreteObserver) Update(subject *Subject) {
	fmt.Println("接到通知...")
}

type ConcreteSubject struct {
	Subject
}

func NewConcreteSubject() *ConcreteSubject {
	subject := &ConcreteSubject{}
	subject.observers = list.New()
	return subject
}

func (subject *ConcreteSubject) SomeOperation() {
	fmt.Println("Subject执行某些操作")
	subject.Notify()
}
