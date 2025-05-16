package utils

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Success bool `json:"success"`
	Data interface{}
}
type ErrorResponse struct {
	Error string `json:"failed"`
	Data string
}

// Convert a boolean to SQLite integer (1 for true, 0 for false)
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Returns nil for nil values or Unix timestamp for populated values
func NilOrNullInt(t *int64) interface{} {
	if t == nil {
		return nil
	}
	return *t
}

func SendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		
		Error:   message,

	})
}

// sendSuccess sends a JSON success response
func SendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data: data,
	})
}
