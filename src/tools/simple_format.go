package main

import (
    "fmt"
    "tabulate"
)

type Blah struct {
    name string
    amount int
}

func main() {
    var blahs []*Blah

    blahs = append(blahs, &Blah{"apple", 15})
    blahs = append(blahs, &Blah{"pommegranit", 1})

    table, err := tabulate.Tabulate(blahs)
    if err != nil {panic(err)}

    fmt.Println(table)
    fmt.Printf("Done!\n")
}
