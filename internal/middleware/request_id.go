package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const HeaderRequestID = "X-Request-ID"

type ctxKey struct{}

func RequestID(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(HeaderRequestID)
		if !valid(id) {
			id = new()
		}

		w.Header().Set(HeaderRequestID, id)
		ctx := context.WithValue(r.Context(), ctxKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func new() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey{}).(string)
	return id
}

// valid only accepts short ids made of safe characters, so a client
// can't inject newlines or huge values into our logs and headers.
func valid(id string) bool {
	if id == "" || len(id) > 64 {
		return false
	}
	for _, c := range id {
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9',
			c == '-', c == '_', c == '.':
		default:
			return false
		}
	}
	return true
}
