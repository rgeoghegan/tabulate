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
    blahs = append(blahs, &Blah{"Orange", 1})

    fmt.Println("No Format:")
    table, err := tabulate.Tabulate(blahs, tabulate.NoFormat)
    if err != nil {panic(err)}
    fmt.Println(table)

    fmt.Println("Plain Format:")
    table, err = tabulate.Tabulate(blahs, tabulate.PlainFormat)
    if err != nil {panic(err)}
    fmt.Println(table)

    fmt.Println("Simple Format:")
    table, err = tabulate.Tabulate(blahs, tabulate.SimpleFormat)
    if err != nil {panic(err)}
    fmt.Println(table)

    fmt.Println("Grid Format:")
    table, err = tabulate.Tabulate(blahs, tabulate.GridFormat)
    if err != nil {panic(err)}
    fmt.Println(table)

    fmt.Println("Fancy Grid Format:")
    table, err = tabulate.Tabulate(blahs, tabulate.FancyGridFormat)
    if err != nil {panic(err)}
    fmt.Println(table)

    fmt.Printf("Done!\n")
}
