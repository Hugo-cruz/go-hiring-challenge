package api

import (
	"encoding/json"
	"net/http"
)

// OKResponse writes a successful JSON response with the provided data.
func OKResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ErrorResponse writes an error JSON response with the provided status and message.
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errorData := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(errorData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
