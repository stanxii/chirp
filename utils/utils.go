package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

func RenderJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

func NormalizeText(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return text
}
