package main

import (
	"fmt"
	"github.com/rgeoghegan/tabulate"
)

type Bloom struct {
	first string
	last  string
}

func (b *Bloom) String() string {
	return fmt.Sprintf("%s %s", b.first, b.last)
}

type Blah struct {
	Name        *Bloom
	Amount      int
	Location    string
	Done        bool
	SurfaceArea float64
}

func main() {
	var blahs []*Blah

	blahs = append(blahs, &Blah{&Bloom{"Roy", "Smith"}, 15, "Washington D.C.",
		true, 0.3453})
	blahs = append(blahs, &Blah{&Bloom{"Fred", "Flanders"}, 100, "Montreal",
		false, 1.0})
	blahs = append(blahs, &Blah{&Bloom{"Bobby", "Smith"}, -2, "San Fransisco",
		false, 124353.23333333})
	blahs = append(blahs, &Blah{&Bloom{"Jolene", "Lee"}, 234, "Guyene", true,
		11.0000000000001})

	table, err := tabulate.Tabulate(
		blahs, &tabulate.Layout{Format: tabulate.PipeFormat},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(table)
}
