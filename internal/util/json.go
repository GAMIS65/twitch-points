package util

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func SendJSON(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")

	// Check if data is a nil slice and convert it to empty slice
	v := reflect.ValueOf(data)
	if data != nil && v.Kind() == reflect.Slice && v.IsNil() {
		// Create a new non-nil slice of the same type
		emptySlice := reflect.MakeSlice(v.Type(), 0, 0).Interface()
		data = emptySlice
	}

	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}
