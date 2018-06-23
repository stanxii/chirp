package api

import (
	"encoding/json"
	"net/http"

	"chirp.com/errors"
)

func RenderJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

func RenderAPIError(w http.ResponseWriter, err *errors.APIError) {
	RenderJSON(w, err, err.Status)
}

// func RenderValidationError(w http.ResponseWriter, vErr *errors.ValidationError) {
// 	RenderJSON(w, vErr, http.StatusUnprocessableEntity)
// }
