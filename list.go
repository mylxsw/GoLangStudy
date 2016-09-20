package main

import (
    "fmt"
    "container/list"
)

func main () {
    items := list.New()
    items.PushFront("Hello")
    items.PushFront("World")

    for item := items.Front(); item != nil; item = item.Next() {
        fmt.Println(item.Value)
    }
}