package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/rcolejnr/jjk-api/internal/models"
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

// parsePagination reads page and limit from request.
// missing values fall back to default values.
func parsePagination(r *http.Request) (models.Pagination, error) {
	p := models.Pagination{Page: 1, Limit: models.DefaultLimit}
	q := r.URL.Query()

	if v := q.Get("page"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return p, fmt.Errorf("invalid page: must be a positive integer")
		}
		p.Page = n
	}

	if v := q.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return p, fmt.Errorf("invalid limit: must be a positive integer")
		}
		p.Limit = n
	}

	p.Normalize()
	return p, nil
}
