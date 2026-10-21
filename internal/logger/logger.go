package logger

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger returns a structured logging middleware using the stdlib slog package.
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			logger.Info(
				http.StatusText(wrapped.statusCode),
				"method", r.Method,
				"path", r.URL.Path,
				"body", r.Body,
				"status", wrapped.statusCode,
				"duration_ms", time.Since(start).Milliseconds(),
				"remote_address", r.RemoteAddr,
			)
		})
	}
}
