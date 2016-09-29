package tabulate

import (
	"bytes"
)

// TableFormatterInterface determines how a layout will format the table.
// Create your own implementation if you need a custom format.
type TableFormatterInterface interface {
	// Passed in a list of column widths (including the header if shown)
	// before drawing the table. Save the widths if you need (for example)
	// to show a bar across a row.
	RegisterWidths([]int)

	// Spacer returns the string to join between the columns (but not
	// before the first column or after the last one)
	Spacer() string
	// Line prefix is shown before the first column.
	LinePrefix() string
	// Line prefix is shown after the last column. Should not contain a return
	// line.
	LinePostfix() string

	// This string appears at the top of the table. Should not contain a
	// return line.
	AboveTable() string
	// This string appears in the table, right after the header. Should not
	// contain a return line.
	BelowHeader() string
	// This string appears in the table, between every "normal" row. Should
	// not contain a return line.
	BetweenRow(index int) string
	// This string appears at the bottom of the table. Should not contain a
	// return line.
	BelowTable() string
}

type spacerFormatting string

func (s spacerFormatting) Spacer() string {
	return string(s)
}

func (s spacerFormatting) RegisterWidths([]int)        {}
func (s spacerFormatting) LinePrefix() string          { return "" }
func (s spacerFormatting) LinePostfix() string         { return "" }
func (s spacerFormatting) AboveTable() string          { return "" }
func (s spacerFormatting) BelowHeader() string         { return "" }
func (s spacerFormatting) BetweenRow(index int) string { return "" }
func (s spacerFormatting) BelowTable() string          { return "" }

type barFormat struct {
	leftCorner  string
	bar         rune
	spacer      string
	rightCorner string
}

func (b *barFormat) draw(colSizes []int) string {
	var bar bytes.Buffer
	var err error

	bar.WriteString(b.leftCorner)

	for i, col := range colSizes {
		if i > 0 {
			_, err = bar.WriteString(b.spacer)
			if err != nil {
				panic(err)
			}
		}
		for j := 0; j < col; j++ {
			_, err = bar.WriteRune(b.bar)
			if err != nil {
				panic(err)
			}
		}
	}

	bar.WriteString(b.rightCorner)

	return bar.String()
}

type headerFormatting struct {
	spacerFormatting
	barSymbol rune
	colSizes  []int
}

func (h *headerFormatting) RegisterWidths(colSizes []int) {
	h.colSizes = colSizes
}

func (h *headerFormatting) BelowHeader() string {
	format := &barFormat{"", h.barSymbol, h.Spacer(), ""}
	return format.draw(h.colSizes)
}

type gridFormatting struct {
	spacerStr string
	leftEdge  string
	rightEdge string

	top    *barFormat
	header *barFormat
	body   *barFormat
	bottom *barFormat

	colSizes []int
}

func (g *gridFormatting) RegisterWidths(colSizes []int) {
	g.colSizes = colSizes
}

func (g *gridFormatting) Spacer() string      { return g.spacerStr }
func (g *gridFormatting) LinePrefix() string  { return g.leftEdge }
func (g *gridFormatting) LinePostfix() string { return g.rightEdge }

func (g *gridFormatting) AboveTable() string {
	return g.top.draw(g.colSizes)
}
func (g *gridFormatting) BelowHeader() string {
	return g.header.draw(g.colSizes)
}
func (g *gridFormatting) BetweenRow(index int) string {
	return g.body.draw(g.colSizes)
}
func (g *gridFormatting) BelowTable() string {
	return g.bottom.draw(g.colSizes)
}

func newGridFormat(left, spacer, right string, top, header, body, bottom *barFormat) *gridFormatting {
	bars := []*barFormat{top, header, body, bottom}

	for _, bar := range bars {
		bar.leftCorner = leftAlign(
			bar.leftCorner, bar.bar, utf8Len(left),
		)
		bar.rightCorner = rightAlign(
			bar.rightCorner, bar.bar, utf8Len(right),
		)
		bar.spacer = center(bar.spacer, bar.bar, utf8Len(spacer))
	}

	return &gridFormatting{
		spacer, left, right,
		top, header, body, bottom,
		nil,
	}
}

func makebarFormat(left rune, bar rune, realSpacer string, middle rune,
	right rune) *barFormat {
	var spacer bytes.Buffer
	spacerLength := len(realSpacer)

	if spacerLength%2 == 0 {
		spacerLength -= 1
	}
	half := spacerLength / 2

	for i := 0; i < half; i++ {
		spacer.WriteRune(bar)
	}
	spacer.WriteRune(middle)
	for i := 0; i < half; i++ {
		spacer.WriteRune(bar)
	}
	if len(realSpacer)%2 == 0 {
		spacer.WriteRune(bar)
	}

	return &barFormat{string(left), bar, spacer.String(), string(right)}
}

var noFormat spacerFormatting = ""
var plainFormat spacerFormatting = " "

var simpleFormat *headerFormatting = &headerFormatting{" ", '-', nil}
var gridFormat *gridFormatting = newGridFormat(
	"| ", " | ", " |",

	&barFormat{"+", '-', "+", "+"},
	&barFormat{"+", '=', "+", "+"},
	&barFormat{"+", '-', "+", "+"},
	&barFormat{"+", '-', "+", "+"},
)
var fancyGridFormat *gridFormatting = newGridFormat(
	"\u2502 ", " \u2502 ", " \u2502",

	&barFormat{"\u2552", '\u2550', "\u2564", "\u2555"},
	&barFormat{"\u255e", '\u2550', "\u256a", "\u2561"},
	&barFormat{"\u251c", '\u2500', "\u253c", "\u2524"},
	&barFormat{"\u2558", '\u2550', "\u2567", "\u255b"},
)
