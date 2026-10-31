package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

type trackingWriter struct {
	http.ResponseWriter
	wrote bool
}

func (t *trackingWriter) WriteHeader(code int) {
	t.wrote = true
	t.ResponseWriter.WriteHeader(code)
}

func (t *trackingWriter) Write(b []byte) (int, error) {
	t.wrote = true
	return t.ResponseWriter.Write(b)
}

func (t *trackingWriter) Unwrap() http.ResponseWriter { return t.ResponseWriter }

// Recover adds into the net/http recover support, turrning a panic into a logged 500 response.
func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tw := &trackingWriter{ResponseWriter: w}

			defer func() {
				rec := recover()
				if rec == nil {
					return
				}

				if rec == http.ErrAbortHandler {
					panic(rec)
				}

				logger.Error("panic recovered",
					"X-Request-ID", FromContext(r.Context()),
					"method", r.Method,
					"path", r.URL.Path,
					"panic", rec,
					"stack", string(debug.Stack()),
				)

				if tw.wrote {
					panic(http.ErrAbortHandler)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"code":    http.StatusInternalServerError,
					"message": "internal server error",
				})
			}()

			next.ServeHTTP(tw, r)
		})
	}
}
