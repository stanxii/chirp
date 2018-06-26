package api

import (
	"encoding/json"
	"net/http"

	"chirp.com/errors"
)

// type ResponseData struct {
// 	Data        interface{}  `json:"data"`
// 	CurrentUser *models.User `json:"user"`
// }

func Render(w http.ResponseWriter, data interface{}) {
	RenderJSON(w, data, http.StatusOK)
}

func RenderAPIError(w http.ResponseWriter, err *errors.APIError) {
	RenderJSON(w, err, err.Status)
}

func RenderJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
