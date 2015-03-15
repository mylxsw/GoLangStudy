package main

import "fmt"

type Person interface {
	speak(words string)
	eat(food string)
}

type Programmer struct {
	username string
	age      int
}

func (person *Programmer) speak(words string) {
	fmt.Println("Speak: ", words, "\n")
}

func (person *Programmer) eat(food string) {
	fmt.Println("Eat: ", food, ", Very delicious!\n")
}

func main() {
	var person Person
	programmer := &Programmer{username: "Tom", age: 20}

	person = programmer
	person.speak("Hello, world")

}
