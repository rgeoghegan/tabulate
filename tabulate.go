// Package tabulate formats data into tables.
//
// Example Usage
//
// The following example will show the table shown below:
//
//     package main
//
//     import (
//         "fmt"
//         "github.com/rgeoghegan/tabulate"
//     )
//
//     type Row struct {
//         name  string
//         count int
//     }
//
//     func main() {
//         table := []*Row{
//             &Row{"alpha", 1},
//             &Row{"bravo", 2},
//         }
//         asText, _ := tabulate.Tabulate(
//		      table, &tabulate.Layout{Format:tabulate.SimpleFormat},
//		   )
//         fmt.Print(asText)
//     }
//
// Which will print out the following:
//
//      name count
//     ----- -----
//     alpha     1
//     bravo     2
//
// You can also provide a slice of slice of strings:
//
//     table := [][]string{
//         []string{"alpha", "1"},
//         []string{"bravo", "2"},
//     }
//     layout := &Layout{Headers:[]string{"name", "count"}, Format:tabulate.SimpleFormat}
//     asText, err := tabulate.Tabulate(table, layout)
package tabulate

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"golang.org/x/exp/utf8string"
)

const (
	noHeaders int = iota
	fromStruct
)

// Layout specifies the general layout of the table. Provide Headers to show a custom list of headings at the top of the table. Set HideHeaders to false to not show Headers.
type Layout struct {
	Format      TableFormatterInterface
	HideHeaders bool
	Headers     []string
}

func getRowType(table interface{}) (reflect.Type, error) {
	tableType := reflect.TypeOf(table)
	if reflect.Slice == tableType.Kind() {
		rowType := tableType.Elem()
		if reflect.Ptr == rowType.Kind() {
			rowType = rowType.Elem()
		}

		switch rowType.Kind() {
		case reflect.Struct, reflect.Slice:
			return rowType, nil
		}
	}
	return nil, fmt.Errorf("Must pass in slice of struct pointers.")
}

type column struct {
	header string
	column []string
}

func fetchStructColumn(rowType reflect.Type, table reflect.Value, colDepth int,
	customHeaders []string, index int) (*column, error) {
	var col = &column{}

	header := rowType.Field(index)
	if customHeaders == nil {
		col.header = header.Name
	} else {
		col.header = customHeaders[index]
	}
	caster, err := guessCaster(header.Type)

	if err != nil {
		return nil, err
	}

	for i := 0; i < colDepth; i++ {
		row := table.Index(i).Elem()
		col.column = append(col.column, caster(row.Field(index)))
	}
	if header.Type.Kind() == reflect.Float32 ||
		header.Type.Kind() == reflect.Float64 {
		alignFloats(col.column)
	}
	return col, nil
}

func fetchMatrixColumn(rowType reflect.Type, table reflect.Value, colDepth int,
	customHeaders []string, hideHeaders bool, index int) (*column, error) {

	col := &column{}

	if !hideHeaders {
		col.header = customHeaders[index]
	}
	cellType := rowType.Elem()
	caster, err := guessCaster(cellType)

	if err != nil {
		return nil, err
	}

	for i := 0; i < colDepth; i++ {
		cell := table.Index(i).Index(index)
		col.column = append(col.column, caster(cell))
	}
	if cellType.Kind() == reflect.Float32 || cellType.Kind() == reflect.Float64 {
		alignFloats(col.column)
	}
	return col, nil
}

type table []*column

func (t table) columnWidths(countHeaders bool) []int {
	var colWidths []int

	for _, col := range t {
		colLength := 0
		if countHeaders {
			colLength = len(col.header)
		}

		for _, cell := range col.column {
			if len(cell) > colLength {
				colLength = len(cell)
			}
		}
		colWidths = append(colWidths, colLength)
	}
	return colWidths
}

func (t table) align(widths []int, showHeaders bool) {
	for colI, col := range t {
		if showHeaders {
			col.header = fmt.Sprintf("%[1]*[2]s", widths[colI], col.header)
		}
		for i := 0; i < len(col.column); i++ {
			col.column[i] = fmt.Sprintf("%[1]*[2]s", widths[colI], col.column[i])
		}
	}
}

func (t table) draw(format TableFormatterInterface, showHeaders bool) string {
	var output []string

	columnWidths := t.columnWidths(showHeaders)
	t.align(columnWidths, showHeaders)
	format.RegisterWidths(columnWidths)

	appendRow := func(rows []string, row string) []string {
		if len(row) > 0 {
			return append(rows, row)
		}
		return rows
	}

	joinTokens := func(parts []string) string {
		return format.LinePrefix() +
			strings.Join(parts, format.Spacer()) +
			format.LinePostfix()
	}

	output = appendRow(output, format.AboveTable())
	if showHeaders {
		parts := make([]string, len(t))
		for i, col := range t {
			parts[i] = col.header
		}
		output = append(output, joinTokens(parts))
		output = appendRow(output, format.BelowHeader())
	}

	for rowI := 0; rowI < len(t[0].column); rowI++ {
		parts := make([]string, len(t))
		for i, col := range t {
			parts[i] = col.column[rowI]
		}
		output = appendRow(output, joinTokens(parts))

		if rowI < len(t[0].column)-1 {
			output = appendRow(output, format.BetweenRow(rowI))
		}
	}
	output = appendRow(output, format.BelowTable())

	return strings.Join(output, "\n") + "\n"
}

func buildTable(data interface{}, layout *Layout) (table, error) {
	rowType, err := getRowType(data)
	if err != nil {
		panic(err)
	}

	tableV := reflect.ValueOf(data)
	tableLength := tableV.Len()

	var columns table
	var isStruct bool
	var colCount int

	switch rowType.Kind() {
	case reflect.Struct:
		isStruct = true
		colCount = rowType.NumField()

	case reflect.Slice:
		isStruct = false
		if layout.HideHeaders {
			// Take the length of the first row as the tables width
			colCount = tableV.Index(0).Len()
		} else {
			if layout.Headers == nil {
				return nil, fmt.Errorf(
					"Must provide headers in layout with slice of slices.",
				)
			}
			colCount = len(layout.Headers)
		}

	default:
		return nil, fmt.Errorf(
			"Inputted data must be a slice of slices or a slice of structs.",
		)
	}

	for col := 0; col < colCount; col++ {
		var rows *column
		var err error

		if isStruct {
			rows, err = fetchStructColumn(
				rowType, tableV, tableLength, layout.Headers, col,
			)
		} else {
			rows, err = fetchMatrixColumn(
				rowType, tableV, tableLength, layout.Headers,
				layout.HideHeaders, col,
			)
		}

		if err != nil {
			return nil, fmt.Errorf("Error with col %d: %s", col, err)
		}
		columns = append(columns, rows)
	}

	return columns, nil
}

// Tabulate will tabulate the provided data with the given layout. If no
// format is specified in the layout, it will use a simple format by default.
//
// Data
//
// The data parameter must either be a slice of structs, and the table will
// use the field names of the struct as column names. If provided a slice
// of slices of strings, you will need to provide a list of Headers (mostly
// so it can figure out how many columns to size for).
//
func Tabulate(data interface{}, layout *Layout) (string, error) {
	columns, err := buildTable(data, layout)
	if err != nil {
		return "", err
	}

	format := layout.Format
	if format == nil {
		format = SimpleFormat
	}

	return columns.draw(format, !layout.HideHeaders), nil
}

func writePadding(combined *bytes.Buffer, length int, padding string) {
	for i := 0; i < length; i++ {
		combined.WriteString(padding)
	}
}

// CombineHorizontal place two tables next to one another
// like:
//
// ╒═══════════╤═══════════╤═══════════╕ ╒═══════════╤═══════════╤═══════════╕
// │         A │         B │         C │ │         A │         B │         C │
// ╞═══════════╪═══════════╪═══════════╡ ╞═══════════╪═══════════╪═══════════╡
// │ A value 1 │ B value 1 │ C value 1 │ │ A value 2 │ B value 2 │ C value 2 │
// ╘═══════════╧═══════════╧═══════════╛ ╘═══════════╧═══════════╧═══════════╛
func CombineHorizontal(left string, right string, padding string) string {
	var combined bytes.Buffer
	leftSplit := strings.Split(left, "\n")
	rightSplit := strings.Split(right, "\n")
	max := len(leftSplit)
	if len(rightSplit) > max {
		max = len(rightSplit)
	}
	for i := 0; i < max; i++ {
		if i < len(leftSplit) && utf8Len(leftSplit[i]) == utf8Len(leftSplit[0]) {
			combined.WriteString(leftSplit[i])
		} else if i < len(rightSplit) && utf8Len(rightSplit[i]) == utf8Len(rightSplit[0]) {
			writePadding(&combined, utf8Len(leftSplit[0]), padding)
		}
		if i < len(rightSplit) && utf8Len(rightSplit[i]) == utf8Len(rightSplit[0]) {
			combined.WriteString(padding)
			combined.WriteString(rightSplit[i])
		} else if i < len(leftSplit) && utf8Len(leftSplit[i]) == utf8Len(leftSplit[0]) {
			combined.WriteString(padding)
			writePadding(&combined, utf8Len(rightSplit[0]), padding)
		}
		if i < max-1 {
			combined.WriteString("\n")
		}
	}
	return combined.String()
}

// CombineVertical place two tables verticaly
// like:
//
// ╒═══════════╤═══════════╤═══════════╕
// │         A │         B │         C │
// ╞═══════════╪═══════════╪═══════════╡
// │ A value 1 │ B value 1 │ C value 1 │
// ╘═══════════╧═══════════╧═══════════╛
// ╒═══════════╤═══════════╤═══════════╕
// │         A │         B │         C │
// ╞═══════════╪═══════════╪═══════════╡
// │ A value 2 │ B value 2 │ C value 2 │
// ╘═══════════╧═══════════╧═══════════╛
func CombineVertical(top string, bottom string, padding string) string {
	var combined bytes.Buffer
	topSplit := strings.Split(top, "\n")
	bottomSplit := strings.Split(bottom, "\n")
	length := utf8Len(topSplit[0])
	if length < utf8Len(bottomSplit[0]) {
		length = utf8Len(bottomSplit[0])
	}
	for i := 0; i < len(topSplit); i++ {
		combined.WriteString(topSplit[i])
		if i < len(topSplit)-1 {
			writePadding(&combined, length-utf8Len(topSplit[i]), " ")
			combined.WriteString("\n")
		}
	}
	if padding != "" {
		paddingUtf8 := utf8string.NewString(padding)
		for i := 0; i < utf8Len(padding); i++ {
			writePadding(&combined, utf8Len(topSplit[0]), string(paddingUtf8.At(i)))
			combined.WriteString("\n")
		}
	}
	for i := 0; i < len(bottomSplit); i++ {
		combined.WriteString(bottomSplit[i])
		if utf8Len(bottomSplit[i]) == len(bottomSplit[0]) {
			writePadding(&combined, length-utf8Len(bottomSplit[i]), " ")
		}
		if i < len(bottomSplit)-1 {
			combined.WriteString("\n")
		}
	}
	return combined.String()
}
