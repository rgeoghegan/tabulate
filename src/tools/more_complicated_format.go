package main

import (
    "fmt"
    "tabulate"
)

type Bloom struct {
    first string
    last string
}

func (b *Bloom) String() string {
    return fmt.Sprintf("%s %s", b.first, b.last)
}

type Blah struct {
    Name *Bloom
    Amount int
    Location string
    Done bool
}

func main() {
    var blahs []*Blah

    blahs = append(blahs, &Blah{&Bloom{"Roy", "Smith"}, 15, "Washington D.C.", true})
    blahs = append(blahs, &Blah{&Bloom{"Fred", "Flanders"}, 100, "Montreal", false})
    blahs = append(blahs, &Blah{&Bloom{"Bobby", "Smith"}, -2, "San Fransisco", false})
    blahs = append(blahs, &Blah{&Bloom{"Jolene", "Lee"}, 234, "Guyene", true})

    table, err := tabulate.Tabulate(blahs, tabulate.HeaderFormat)
    if err != nil {panic(err)}

    fmt.Println(table)
}
