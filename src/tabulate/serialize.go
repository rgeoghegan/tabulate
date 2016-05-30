package tabulate

import (
    "reflect"
    "strconv"
    "fmt"
)


func intToString(integer reflect.Value) string {
    return strconv.Itoa(int(integer.Int()))
}


func floatToString(floating reflect.Value) string {
    return strconv.FormatFloat(floating.Float(), 'f', -1, floating.Type().Bits())
}


func boolToString(boolean reflect.Value) string {
    if boolean.Bool() {
        return "true"
    }
    return "false"
}


func stringToString(str reflect.Value) string {
    return str.String()
}


func callString(value reflect.Value) string {
    toString := value.MethodByName("String")
    res := toString.Call(nil)
    return res[0].String()
}

func guessCaster(cellType reflect.Type) ((func (reflect.Value) string), error) {
    switch cellType.Kind() {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
                reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
                reflect.Uint32, reflect.Uint64:
            return intToString, nil

        case reflect.Float32, reflect.Float64:
            return floatToString, nil

        case reflect.Bool:
            return boolToString, nil

        case reflect.String:
            return stringToString, nil
    }

    _, found := cellType.MethodByName("String")
    if ! found {
        return nil, fmt.Errorf("Column must either contain an int, a float, a bool, a string or something implementing the fmt.Stringer interface.")
    }

    return callString, nil
}