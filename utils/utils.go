package utils

import (
	"strings"
	"unicode"
)

func NormalizeText(text string) string {
	text = strings.ToLower(text)
	text = TrimAllWhiteSpace(text)
	return text
}

func TrimAllWhiteSpace(text string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, text)
}
