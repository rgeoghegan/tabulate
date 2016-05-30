package tabulate

import (
    "bytes"
)

type TableFormatterInterface interface {
    RegisterWidths([]int)
    Spacer() string
    LinePrefix() string
    LinePostfix() string

    AboveTable() string
    BelowHeader() string
    BetweenRow(index int) string
    BelowTable() string
}

type SpacerFormatting string

func (s SpacerFormatting) Spacer() string {
    return string(s)
}

func (s SpacerFormatting) RegisterWidths([]int) {}
func (s SpacerFormatting) LinePrefix() string {return ""}
func (s SpacerFormatting) LinePostfix() string {return ""}
func (s SpacerFormatting) AboveTable() string {return ""}
func (s SpacerFormatting) BelowHeader() string {return ""}
func (s SpacerFormatting) BetweenRow(index int) string {return ""}
func (s SpacerFormatting) BelowTable() string {return ""}

type BarFormat struct {
    LeftCorner string
    Bar rune
    Spacer string
    RightCorner string
}

func (b *BarFormat) Draw(colSizes []int) string {
    var bar bytes.Buffer
    var err error

    bar.WriteString(b.LeftCorner)

    for i, col := range colSizes {
        if i > 0 {
            _, err = bar.WriteString(b.Spacer)
            if (err != nil) {panic(err)}
        }
        for j := 0; j < col; j++ {
            _, err = bar.WriteRune(b.Bar)
            if (err != nil) {panic(err)}
        }
    }

    bar.WriteString(b.RightCorner)

    return bar.String()
}


type HeaderFormatting struct {
    SpacerFormatting
    BarSymbol rune
    colSizes []int
}


func (h *HeaderFormatting) RegisterWidths(colSizes []int) {
    h.colSizes = colSizes
}


func (h *HeaderFormatting) BelowHeader() string {
    format := &BarFormat{"", h.BarSymbol, h.Spacer(), ""}
    return format.Draw(h.colSizes)
}


type GridFormatting struct {
    SpacerStr string
    LeftEdge string
    RightEdge string

    Top *BarFormat
    Header *BarFormat
    Body *BarFormat
    Bottom *BarFormat

    ColSizes []int
}


func (g *GridFormatting) RegisterWidths(colSizes []int) {
    g.ColSizes = colSizes
}

func (g *GridFormatting) Spacer() string {return g.SpacerStr}
func (g *GridFormatting) LinePrefix() string {return g.LeftEdge}
func (g *GridFormatting) LinePostfix() string {return g.RightEdge}

func (g *GridFormatting) AboveTable() string {
    return g.Top.Draw(g.ColSizes)
}
func (g *GridFormatting) BelowHeader() string {
    return g.Header.Draw(g.ColSizes)
}
func (g *GridFormatting) BetweenRow(index int) string {
    return g.Body.Draw(g.ColSizes)
}
func (g *GridFormatting) BelowTable() string {
    return g.Bottom.Draw(g.ColSizes)
}


func NewGridFormat(left string, spacer string, right string, top, header, body, bottom *BarFormat) *GridFormatting {
    bars := []*BarFormat{top, header, body, bottom}

    for _, bar := range bars {
        bar.LeftCorner = leftAlign(
            bar.LeftCorner, bar.Bar, utf8Len(left),
        )
        bar.RightCorner = rightAlign(
            bar.RightCorner, bar.Bar, utf8Len(right),
        )
        bar.Spacer = center(bar.Spacer, bar.Bar, utf8Len(spacer))
    }

    return &GridFormatting{
        spacer, left, right,
        top, header, body, bottom,
        nil,
    }
}


func makeBarFormat(left rune, bar rune, realSpacer string, middle rune, right rune) *BarFormat {
    var spacer bytes.Buffer
    spacerLength := len(realSpacer)

    if spacerLength % 2 == 0 {
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
    if len(realSpacer) % 2 == 0 {
        spacer.WriteRune(bar)
    }

    return &BarFormat{string(left), bar, spacer.String(), string(right)}
}


// Implemented Formats

var NoFormat SpacerFormatting = ""
var PlainFormat SpacerFormatting = " "
var SimpleFormat *HeaderFormatting = &HeaderFormatting{" ", '-', nil}
var GridFormat *GridFormatting = NewGridFormat(
    "| ", " | ", " |",

    &BarFormat{"+", '-', "+", "+"},
    &BarFormat{"+", '=', "+", "+"},
    &BarFormat{"+", '-', "+", "+"},
    &BarFormat{"+", '-', "+", "+"},
)
var FancyGridFormat *GridFormatting = NewGridFormat(
    "\u2502 ", " \u2502 ", " \u2502",

    &BarFormat{"\u2552", '\u2550', "\u2564", "\u2555"},
    &BarFormat{"\u255e", '\u2550', "\u256a", "\u2561"},
    &BarFormat{"\u251c", '\u2500', "\u253c", "\u2524"},
    &BarFormat{"\u2558", '\u2550', "\u2567", "\u255b"},
)
