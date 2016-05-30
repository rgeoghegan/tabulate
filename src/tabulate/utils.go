package tabulate

import (
    "unicode/utf8"
    "bytes"
)


func utf8Len(str string) int {return utf8.RuneCountInString(str)}


func padToken(token string, padding rune, left int, right int) string {
    var output bytes.Buffer

    for i := 0; i < left; i++ {output.WriteRune(padding)}
    output.WriteString(token)
    for i := 0; i < right; i++ {output.WriteRune(padding)}
    return output.String()
}


func center(token string, padding rune, size int) string {
    padLength := size - utf8Len(token)
    half := padLength / 2

    return padToken(token, padding, half, half + (padLength % 2))
}


func leftAlign(token string, padding rune, size int) string {
    return padToken(token, padding, 0, size - utf8Len(token))
}

func rightAlign(token string, padding rune, size int) string {
    return padToken(token, padding, size - utf8Len(token), 0)
}