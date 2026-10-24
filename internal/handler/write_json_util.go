package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

// isValidationErr checks if the error originated from a validation rule.
func isValidationErr(err error) bool {
	return strings.HasPrefix(err.Error(), "validation:")
}

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) error {
	return writeJson(w, status, apiError{Code: status, Message: err.Error()})
}

// APIError is the standard error response format for the API.
type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
