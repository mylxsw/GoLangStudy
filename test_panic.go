package main

import "fmt"
import "os"

func main() {
	var _recover = true
	if len(os.Args) > 1 {
		_recover = os.Args[1] == "true"
	}
	f(_recover)
	fmt.Println("Returned normally from f.")
}

func f(_recover bool) {
	defer func() {
		if _recover {

			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}
	}()

	fmt.Println("Calling g.")
	g(0)
	fmt.Println("Returned normally from g.")
}

func g(i int) {
	if i > 3 {
		fmt.Println("Panicking!")
		panic(fmt.Sprintf("%v", i))
	}

	defer fmt.Println("Defer in g", i)
	fmt.Println("Printing in g", i)
	g(i + 1)
}
