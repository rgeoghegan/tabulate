package main

import (
    "reflect"
    "fmt"
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

func main() {
    a := 1
    aType := reflect.TypeOf(a)
    aValue := reflect.ValueOf(a)
    fmt.Printf("a kind? %v\n", aType.Kind())

    caster, err := guessCaster(aType)
    if err != nil {panic(err)}
    fmt.Printf("a: %v\n", caster(aValue))
}
