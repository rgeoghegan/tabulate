package tabulate

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

func utf8Len(str string) int { return utf8.RuneCountInString(str) }

func padToken(token string, padding rune, left int, right int) string {
	var output bytes.Buffer

	for i := 0; i < left; i++ {
		output.WriteRune(padding)
	}
	output.WriteString(token)
	for i := 0; i < right; i++ {
		output.WriteRune(padding)
	}
	return output.String()
}

func center(token string, padding rune, size int) string {
	padLength := size - utf8Len(token)
	half := padLength / 2

	return padToken(token, padding, half, half+(padLength%2))
}

func leftAlign(token string, padding rune, size int) string {
	return padToken(token, padding, 0, size-utf8Len(token))
}

func rightAlign(token string, padding rune, size int) string {
	return padToken(token, padding, size-utf8Len(token), 0)
}

func writePadding(combined *bytes.Buffer, length int, padding string) {
	for i := 0; i < length; i++ {
		combined.WriteString(padding)
	}
}

func CombineHorizontal(left string, right string, padding string) string {
	var combined bytes.Buffer
	leftSplit := strings.Split(left, "\n")
	rightSplit := strings.Split(right, "\n")
	max := len(leftSplit)
	if len(rightSplit) > max {
		max = len(rightSplit)
	}
	for i := 0; i < max; i++ {
		if i < len(leftSplit) && utf8.RuneCountInString(leftSplit[i]) == utf8.RuneCountInString(leftSplit[0]) {
			combined.WriteString(leftSplit[i])
		} else {
			writePadding(&combined, utf8.RuneCountInString(leftSplit[0]), " ")
		}
		if i < len(rightSplit) {
			combined.WriteString(padding)
			combined.WriteString(rightSplit[i])
		}
		combined.WriteString("\n")
	}
	return combined.String()
}

func CombineVertical(top string, bottom string) string {
	var combined bytes.Buffer
	topSplit := strings.Split(top, "\n")
	bottomSplit := strings.Split(bottom, "\n")
	length := utf8.RuneCountInString(topSplit[0])
	if length < utf8.RuneCountInString(bottomSplit[0]) {
		length = utf8.RuneCountInString(bottomSplit[0])
	}
	for i := 0; i < len(topSplit); i++ {
		combined.WriteString(topSplit[i])
		writePadding(&combined, length-utf8.RuneCountInString(topSplit[i]), " ")
		combined.WriteString("\n")
	}
	for i := 0; i < len(bottomSplit); i++ {
		combined.WriteString(bottomSplit[i])
		writePadding(&combined, length-utf8.RuneCountInString(bottomSplit[i]), " ")
		combined.WriteString("\n")
	}
	return combined.String()
}
