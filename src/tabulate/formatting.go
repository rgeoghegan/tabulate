package tabulate

import "bytes"

type TableFormatterInterface interface {
    RegisterWidths([]int)
    Seperator() string
    LinePrefix() string
    LinePostfix() string

    AboveTable() (string, bool)
    BelowHeader() (string, bool)
    BetweenRow(index int) (string, bool)
    BelowTable() (string, bool)
}

type NoFormatting string

func (n NoFormatting) Seperator() string {
    return string(n)
}

func (n NoFormatting) RegisterWidths([]int) {}
func (n NoFormatting) LinePrefix() string {return ""}
func (n NoFormatting) LinePostfix() string {return ""}
func (n NoFormatting) AboveTable() (string, bool) {return "", false}
func (n NoFormatting) BelowHeader() (string, bool) {return "", false}
func (n NoFormatting) BetweenRow(index int) (string, bool) {return "", false}
func (n NoFormatting) BelowTable() (string, bool) {return "", false}


type HeaderFormatting struct {
    NoFormatting
    BarSymbol rune
    colSizes []int
}

func (h *HeaderFormatting) RegisterWidths(colSizes []int) {
    h.colSizes = colSizes
}

func (h *HeaderFormatting) BelowHeader() (string, bool) {
    var bar bytes.Buffer
    var err error

    for i, col := range h.colSizes {
        if i > 0 {
            _, err = bar.WriteString(h.Seperator())
            if (err != nil) {panic(err)}
        }
        for j := 0; j < col; j++ {
            _, err = bar.WriteRune(h.BarSymbol)
            if (err != nil) {panic(err)}
        }
    }

    return bar.String(), true
}

// Implemented Formats

var NoFormat NoFormatting = ""
var SimpleFormat NoFormatting = " "
var HeaderFormat *HeaderFormatting = &HeaderFormatting{" ", '-', nil}
