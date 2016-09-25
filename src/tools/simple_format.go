package main

import (
	"fmt"
	"tabulate"
)

type Blah struct {
	name   string
	amount int
}

func main() {
	var blahs []*Blah

	blahs = append(blahs, &Blah{"Apple", 15})
	blahs = append(blahs, &Blah{"Orange", 1})

	fmt.Println("No Format:")
	layout := tabulate.NoFormatLayout()
	table, err := tabulate.Tabulate(blahs, layout)
	if err != nil {
		panic(err)
	}
	fmt.Println(table)

	fmt.Println("No Format, No Header:")
	layout.HideHeaders = true
	table, err = tabulate.Tabulate(blahs, layout)
	if err != nil {
		panic(err)
	}
	fmt.Println(table)

	fmt.Println("Plain Format:")
	layout = tabulate.PlainLayout()
	table, err = tabulate.Tabulate(blahs, layout)
	if err != nil {
		panic(err)
	}
	fmt.Println(table)

	fmt.Println("Plain Format, No Header:")
	layout.HideHeaders = true
	table, err = tabulate.Tabulate(blahs, layout)
	if err != nil {
		panic(err)
	}
	fmt.Println(table)

	fmt.Println("Simple Format:")
	table, err = tabulate.Tabulate(blahs, tabulate.SimpleLayout())
	if err != nil {
		panic(err)
	}
	fmt.Println(table)

	fmt.Println("Grid Format:")
	table, err = tabulate.Tabulate(blahs, tabulate.GridLayout())
	if err != nil {
		panic(err)
	}
	fmt.Println(table)

	fmt.Println("Fancy Grid Format:")
	table, err = tabulate.Tabulate(blahs, tabulate.FancyGridLayout())
	if err != nil {
		panic(err)
	}
	fmt.Println(table)

	fmt.Printf("Done!\n")
}
