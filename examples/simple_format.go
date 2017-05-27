package main

import (
	"fmt"
	"github.com/rgeoghegan/tabulate"
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
	layout := &tabulate.Layout{Format: tabulate.NoFormat}
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
	layout = &tabulate.Layout{Format: tabulate.PlainFormat}
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
	layout = &tabulate.Layout{Format: tabulate.SimpleFormat}
	table, err = tabulate.Tabulate(blahs, layout)
	if err != nil {
		panic(err)
	}
	fmt.Println(table)

	fmt.Println("Grid Format:")
	layout = &tabulate.Layout{Format: tabulate.GridFormat}
	table, err = tabulate.Tabulate(blahs, layout)
	if err != nil {
		panic(err)
	}
	fmt.Println(table)

	fmt.Println("Fancy Grid Format:")
	layout = &tabulate.Layout{Format: tabulate.FancyGridFormat}
	table, err = tabulate.Tabulate(blahs, layout)
	if err != nil {
		panic(err)
	}
	fmt.Println(table)

	fmt.Printf("Done!\n")
}
