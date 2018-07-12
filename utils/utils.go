package utils

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"
)

func RenderJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

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
