package main

import (
	"aicode.cc/design/observer"
	"fmt"
)

func main() {
	fmt.Println("Observer")

	var observer_1 observer.Observer
	var observer_2 observer.Observer

	observer_1 = &observer.ConcreteObserver{}
	observer_2 = &observer.ConcreteObserver{}

	subject := observer.NewConcreteSubject()
	subject.Attach(&observer_1)
	subject.Attach(&observer_2)

	subject.SomeOperation()
}
