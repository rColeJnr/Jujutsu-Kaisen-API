package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
)

const maxBodyBytes = 1 << 20 // 1 MiB

// decodeJSON reads the request body into dst, strictly:
//   - Content-Type must be application/json (415 otherwise)
//   - the body is capped at maxBodyBytes (413 otherwise)
//   - unknown fields are rejected
//   - nothing may follow the JSON value
//
// On failure it writes the error response itself and returns false,
func decodeJSON(w htp.ResponseWriter, r *http.Request, dst any) bool {
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mt != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, fmt.Errorf("Content-Type must be application/json"))
		return false
	}

	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var (
			syntaxErr *json.SyntaxError
			typeErr   *json.UnmarshalTypeError
			tooBig    *http.MaxBytesError
		)
		switch {
		case errors.As(err, &tooBig):
			writeError(w, http.StatusRequestEntityTooLarge, fmt.Errorf("request body must not be larger than %d bytes", maxBodyBytes))
		case errors.Is(err, io.EOF):
			writeError(w, http.StatusBadRequest, fmt.Errorf("request body must not be empty"))
		case errors.As(err, &syntaxErr), errors.Is(err, io.ErrUnexpectedEOF):
			writeError(w, http.StatusBadRequest, fmt.Errorf("request body contains malformed JSON"))
		case errors.As(err, &typeErr):
			writeError(w, http.StatusBadRequest, fmt.Errorf("field %q has the wrong type", typeErr.Field))
		default:
			// e.g. `json: unknown field "nmae"`
			writeError(w, http.StatusBadRequest, err)
		}
		return false
	}

	// A second Decode must hit EOF; anything else means extra data after the value.
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, fmt.Errorf("request body must contain a single JSON value"))
		return false
	}

	return true
}

// pathID parses the {id} segment of the URL and rejects anything below 1.
func pathID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("invalid id: must be a positive integer")
	}
	return id, nil
}
