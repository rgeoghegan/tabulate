package tabulate

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	noHeaders int = iota
	fromStruct
)

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

func fetchStructColumn(rowType reflect.Type, table reflect.Value, colDepth int,
	customHeaders []string, index int) ([]string, error) {
	var output []string
	var headerName string

	header := rowType.Field(index)
	if customHeaders == nil {
		headerName = header.Name
	} else {
		headerName = customHeaders[index]
	}
	output = append(output, headerName)
	caster, err := guessCaster(header.Type)

	if err != nil {
		return nil, err
	}

	for i := 0; i < colDepth; i++ {
		row := table.Index(i).Elem()
		output = append(output, caster(row.Field(index)))
	}
	if header.Type.Kind() == reflect.Float32 ||
		header.Type.Kind() == reflect.Float64 {
		alignFloats(output[1:len(output)])
	}
	return output, nil
}

func fetchMatrixColumn(rowType reflect.Type, table reflect.Value, colDepth int,
	customHeaders []string, index int) ([]string, error) {

	var output []string

	output = append(output, customHeaders[index])
	cellType := rowType.Elem()
	caster, err := guessCaster(cellType)

	if err != nil {
		return nil, err
	}

	for i := 0; i < colDepth; i++ {
		cell := table.Index(i).Index(index)
		output = append(output, caster(cell))
	}
	if cellType.Kind() == reflect.Float32 || cellType.Kind() == reflect.Float64 {
		alignFloats(output[1:len(output)])
	}
	return output, nil
}

type table [][]string

func (t table) columnWidths(countHeaders bool) []int {
	var colWidths []int

	for _, col := range t {
		colLength := 0
		for rowI, cell := range col {
			if (countHeaders) || (rowI > 0) {
				if len(cell) > colLength {
					colLength = len(cell)
				}
			}
		}
		colWidths = append(colWidths, colLength)
	}
	return colWidths
}

func (t table) align(widths []int, showHeaders bool) {
	for colI, col := range t {
		start := 0
		if !showHeaders {
			// Skip first row because it contains the headers
			start = 1
		}
		for i := start; i < len(col); i++ {
			col[i] = fmt.Sprintf("%[1]*[2]s", widths[colI], col[i])
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

	joinCols := func(rowI int) string {
		var parts []string
		for _, col := range t {
			parts = append(parts, col[rowI])
		}
		return format.LinePrefix() +
			strings.Join(parts, format.Spacer()) +
			format.LinePostfix()
	}

	output = appendRow(output, format.AboveTable())
	if showHeaders {
		output = append(output, joinCols(0))
		output = appendRow(output, format.BelowHeader())
	}

	for rowI := 1; rowI < len(t[0]); rowI++ {
		output = appendRow(output, joinCols(rowI))

		if rowI < len(t[0])-1 {
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
		if layout.Headers == nil {
			return nil, fmt.Errorf("Must provide headers in layout with slice of slices.")
		}
		isStruct = false
		colCount = len(layout.Headers)
	default:
		return nil, fmt.Errorf("Inputted data must be a slice of slices or a slice of structs.")
	}

	for col := 0; col < colCount; col++ {
		var rows []string
		var err error

		if isStruct {
			rows, err = fetchStructColumn(
				rowType, tableV, tableLength, layout.Headers, col,
			)
		} else {
			rows, err = fetchMatrixColumn(
				rowType, tableV, tableLength, layout.Headers, col,
			)
		}

		if err != nil {
			return nil, fmt.Errorf("Error with col %d: %s", col, err)
		}
		columns = append(columns, rows)
	}

	return columns, nil
}

// Implemented Layouts w/ Formats
func NoFormatLayout() *Layout  { return &Layout{Format: NoFormat} }
func PlainLayout() *Layout     { return &Layout{Format: PlainFormat} }
func SimpleLayout() *Layout    { return &Layout{Format: SimpleFormat} }
func GridLayout() *Layout      { return &Layout{Format: GridFormat} }
func FancyGridLayout() *Layout { return &Layout{Format: FancyGridFormat} }

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
