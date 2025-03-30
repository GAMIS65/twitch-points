package util

import (
	"encoding/json"
	"net/http"
)

func SendJSON(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}
