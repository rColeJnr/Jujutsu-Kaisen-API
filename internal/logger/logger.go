package logger

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/rcolejnr/jjk-api/internal/middleware"
)

const (
	maxBodyBytes = 1 << 20 // requests larger than 1 MiB are rejected
	maxLoggedLen = 512     // only the first 512 bytes of a body are logged
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
			if r.Body != nil && !strings.HasPrefix(r.URL.Path, "/api/auth/") {
				bodyBytes, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodyBytes))
				if err != nil {
					wrapped.Header().Set("Content-Type", "application/json")
					wrapped.WriteHeader(http.StatusRequestEntityTooLarge)
					_, _ = wrapped.Write([]byte(`{"code":413,"message":"request body too large"}` + "\n"))
					logger.Info(http.StatusText(http.StatusRequestEntityTooLarge),
						"X-Request-ID", middleware.FromContext(r.Context()),
						"method", r.Method,
						"path", r.URL.Path,
						"status", http.StatusRequestEntityTooLarge,
						"remote_address", r.RemoteAddr,
					)
					return
				}
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				bodyStr = string(bodyBytes)
				if len(bodyStr) > maxLoggedLen {
					bodyStr = bodyStr[:maxLoggedLen] + "...(truncated)"
				}
			} else if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
			}

			next.ServeHTTP(wrapped, r)

			logger.Info(
				http.StatusText(wrapped.statusCode),
				"X-Request-ID", middleware.FromContext(r.Context()),
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
