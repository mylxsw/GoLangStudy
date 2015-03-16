package main

import "fmt"
import "time"

func main() {
	year, month, day := time.Now().Date()
	fmt.Printf("%d/%d/%d\n", year, month, day)
	fmt.Println(time.Now().Format("2006/01/02"))
}
