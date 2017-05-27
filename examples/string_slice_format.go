package main

import (
	"fmt"
	"github.com/rgeoghegan/tabulate"
)

func main() {
	table := [][]string{
		[]string{"apple", "blueberry"},
		[]string{"ant", "beetle"},
	}

	fmt.Println("Simple Format:")
	layout := &tabulate.Layout{Format: tabulate.SimpleFormat}
	layout.Headers = []string{"A", "B"}
	tableText, err := tabulate.Tabulate(table, layout)
	if err != nil {
		panic(err)
	}

	fmt.Println(tableText)
}
