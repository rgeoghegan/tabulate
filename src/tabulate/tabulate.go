package tabulate

import (
    "reflect"
    "fmt"
    "strings"
)


func guessCaster(cellType reflect.Type) ((func (reflect.Value) string), error) {
    switch cellType.Kind() {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
                reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
                reflect.Uint32, reflect.Uint64:
            return func (integer reflect.Value) string {
                return fmt.Sprintf("%d", integer.Int())
            }, nil

        case reflect.Float32, reflect.Float64:
            return func (floating reflect.Value) string {
                return fmt.Sprintf("%f", floating.Float())
            }, nil

        case reflect.Bool:
            return func (boolean reflect.Value) string {
                if boolean.Bool() {
                    return "true"
                }
                return "false"
            }, nil

        case reflect.String:
            return func (str reflect.Value) string {
                return str.String()
            }, nil
    }

    _, found := cellType.MethodByName("String")
    if ! found {
        return nil, fmt.Errorf("Column must either contain an int, a float, a bool, a string or something implementing the fmt.Stringer interface.")
    }

    return func (value reflect.Value) string {
        toString := value.MethodByName("String")
        res := toString.Call(nil)
        return res[0].String()
    }, nil
}


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
    var add bool

    for rowI := range t[0] {
        var line []string
        for _, col := range t {
            line = append(line, col[rowI])
        }

        if rowI == 0 {
            row, add = format.AboveTable()
            if add {
                output = append(output, row)
            }
        }

        output = append(
            output,
            format.LinePrefix() +
            strings.Join(line, format.Seperator()) +
            format.LinePostfix(),
        )

        if rowI == 0 {
            row, add = format.BelowHeader()
            if add {
                output = append(output, row)
            }
        }
    }

    row, add = format.BelowTable()
    if (add) {
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