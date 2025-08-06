package utils

import (
	"encoding/json"
	"net/http"
)

// Respond sends a JSON response with the specified HTTP status code and data.
func Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		byteData, ok := data.([]byte)
		if !ok {
			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			err := enc.Encode(data)
			if err != nil {
				return
			}

			w.Header().Set("Content-Type", "application/json")
			return
		}
		_, err := w.Write(byteData)
		if err != nil {
			return
		}
	}
}
