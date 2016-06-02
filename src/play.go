package main

import (
    "fmt"
)

func main() {
    a := []string{"a", "b", "c"}
    fmt.Printf("%v\n", a)

    b := a[1:len(a)]
    fmt.Printf("%v\n", b)
}