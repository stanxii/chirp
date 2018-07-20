package utils

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"

	"chirp.com/errors"
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

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Render(w http.ResponseWriter, data interface{}) {
	renderHTTP(w, data, http.StatusOK)
}

func RenderAPIError(w http.ResponseWriter, err *errors.APIError) {
	renderHTTP(w, err, err.Status)
}

func renderHTTP(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
