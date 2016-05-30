package tabulate

import (
    "reflect"
    "fmt"
    "strings"
)


func getRowType(table interface{}) (reflect.Type, error) {
    tableType := reflect.TypeOf(table)
    if reflect.Slice == tableType.Kind() {
        
        rowPtrType := tableType.Elem()
        if reflect.Ptr == rowPtrType.Kind() {

            rowType := rowPtrType.Elem()
            if reflect.Struct == rowType.Kind() {
                return rowType, nil
            }
        }
    }
    return nil, fmt.Errorf("Must pass in slice of struct pointers.")
}


func fetchColumn(rowType reflect.Type, table reflect.Value, len int, index int) ([]string, error) {
    var output []string
    
    header := rowType.Field(index)
    output = append(output, header.Name)
    caster, err := guessCaster(header.Type)

    if err != nil {return nil, err}

    for i := 0; i < len; i++ {
        row := table.Index(i).Elem()
        output = append(output, caster(row.Field(index)))
    }
    return output, nil
}


type table [][]string

func (t table) columnWidths() []int {
    var colWidths []int
    for _, col := range t {
        colLength := 0
        for _, cell := range col {
            if len(cell) > colLength {
                colLength = len(cell)
            }
        }
        colWidths = append(colWidths, colLength)
    }
    return colWidths
}


func (t table) align(widths []int) {
    for colI, col := range t {
        for i := 0; i < len(col); i++ {
            col[i] = fmt.Sprintf("%[1]*[2]s", widths[colI], col[i])
        }
    }
}

func (t table) draw(format TableFormatterInterface) string {
    var output []string

    columnWidths := t.columnWidths()
    t.align(columnWidths)
    format.RegisterWidths(columnWidths)

    var row string

    for rowI := range t[0] {
        var line []string
        for _, col := range t {
            line = append(line, col[rowI])
        }

        if rowI == 0 {
            row = format.AboveTable()
            if len(row) > 0 {
                output = append(output, row)
            }
        }

        output = append(
            output,
            format.LinePrefix() +
            strings.Join(line, format.Spacer()) +
            format.LinePostfix(),
        )

        switch {
        case rowI == 0:
            row = format.BelowHeader()
        case rowI == len(t[0]) - 1:
            row = ""
        case true:
            row = format.BetweenRow(rowI)
        }
    
        if len(row) > 0 {
            output = append(output, row)
        }
    }

    row = format.BelowTable()
    if len(row) > 0 {
        output = append(output, row)
    }

    return strings.Join(output, "\n")
}


func Tabulate(data interface{}, format TableFormatterInterface) (string, error) {
    rowType, err := getRowType(data)
    if err != nil {
        return "", err
    }

    tableV := reflect.ValueOf(data)
    tableLength := tableV.Len()
    var columns table

    for col := 0; col < rowType.NumField(); col++ {
        rows, err := fetchColumn(rowType, tableV, tableLength, col)
        if err != nil {return "", fmt.Errorf("Error with col %d: %s", col, err)}

        columns = append(columns, rows)
    }

    return columns.draw(format), nil
}