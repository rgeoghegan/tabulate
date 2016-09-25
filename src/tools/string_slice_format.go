package main

import (
	"fmt"
	"tabulate"
)

func main() {
	table := [][]string{
		[]string{"apple", "blueberry"},
		[]string{"ant", "beetle"},
	}

	fmt.Println("Simple Format:")
	layout := tabulate.SimpleLayout()
	layout.Headers = []string{"A", "B"}
	tableText, err := tabulate.Tabulate(table, layout)
	if err != nil {
		panic(err)
	}

	fmt.Println(tableText)
}
