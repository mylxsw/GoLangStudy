package main

import "fmt"
import alg "aicode.cc/algorithm"

func main() {
	var items []int
	items = []int{1, 100, 33, 43, 22, 56, 99, 5, 9, 0, 9}
	fmt.Println(alg.BubbleSort(items))
	fmt.Println(alg.SelectSort(items))
	new_items := items
	alg.QuickSort(&new_items, 0, len(new_items)-1)
	fmt.Println(new_items)
}
