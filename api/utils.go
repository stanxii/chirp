package api

import (
	"encoding/json"
	"net/http"

	"chirp.com/errors"
)

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
