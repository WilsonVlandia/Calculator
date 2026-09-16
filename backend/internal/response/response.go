// Package response builds and serializes the JSON bodies returned by
// the API. It is the only package allowed to write a response body, so
// that every endpoint shares one consistent envelope.
package response

import (
	"encoding/json"
	"net/http"
)

// successBody is the JSON shape returned by every successful
// calculation.
type successBody struct {
	Success   bool    `json:"success"`
	Operation string  `json:"operation"`
	Result    float64 `json:"result"`
}

// WriteSuccess writes a 200 OK JSON response describing a completed
// calculation.
func WriteSuccess(w http.ResponseWriter, operation string, result float64) {
	writeJSON(w, http.StatusOK, successBody{
		Success:   true,
		Operation: operation,
		Result:    result,
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
