package logger

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	rBody      *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode >= 400 {
		rw.rBody.Write(b)
	}
	return rw.ResponseWriter.Write(b)
}

// Logger returns a structured logging middleware using the stdlib slog package.
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				rBody:          bytes.NewBuffer(nil),
			}

			var bodyStr string
			if r.Body != nil {
				bodyBytes, err := io.ReadAll(r.Body)
				if err == nil {
					bodyStr = string(bodyBytes)
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}

			next.ServeHTTP(wrapped, r)

			logger.Info(
				http.StatusText(wrapped.statusCode),
				"method", r.Method,
				"path", r.URL.Path,
				"body", bodyStr,
				"status", wrapped.statusCode,
				"duration_ms", time.Since(start).Milliseconds(),
				"remote_address", r.RemoteAddr,
				"error_msg", wrapped.rBody.String(),
			)

		})
	}
}
