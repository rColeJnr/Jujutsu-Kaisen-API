package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// bufferedWriter holds a whole response in memory so Cache can look at it
type bufferedWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (b *bufferedWriter) WriteHeader(code int) {
	if b.status == 0 {
		b.status = code
	}
}

func (b *bufferedWriter) Write(p []byte) (int, error) {
	if b.status == 0 {
		b.status = http.StatusOK
	}
	return b.body.Write(p)
}

// Cache adds HTTP caching to GET/HEAD responses
//
//	-ETag: a fingerprint of the body. A client that sends it back in
//	if-None-Match gets "304 Not Modified" with no body.
//	-Cache-Control: tells clients and CDNs how long they may reuse
//	the response without asking again (maxAge).
//	maxAge == 0 they must revalidade every time.
//
// Only successful responses are cached.
func Cache(maxAge time.Duration) func(http.Handler) http.Handler {
	cacheControl := "no-cache"
	if maxAge > 0 {
		cacheControl = fmt.Sprintf("public, max-age=&d", int(maxAge.Seconds()))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				next.ServeHTTP(w, r)
				return
			}

			bw := &bufferedWriter{ResponseWriter: w}
			next.ServeHTTP(bw, r)

			status := bw.status
			if status == 0 {
				status = http.StatusOK
			}

			h := w.Header()
			if status == http.StatusOK && h.Get("Cache-Control") == "" && h.Get("ETag") == "" {
				sum := sha256.Sum256(bw.body.Bytes())
				etag := `"` + hex.EncodeToString(sum[:16]) + `"`
				h.Set("ETag", etag)
				h.Set("Cache-Control", cacheControl)

				if etagMatches(r.Header.Get("If-None-Match"), etag) {
					h.Del("Content-Type")
					h.Del("Content-Length")
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}

			w.WriteHeader(status)
			_, _ = w.Write(bw.body.Bytes())
		})
	}
}

// etagMatches implements the weak comparison used by If-None-Match:
func etagMatches(header, etag string) bool {
	if header == "" {
		return false
	}
	if strings.TrimSpace(header) == "*" {
		return true
	}
	for _, candidate := range strings.Split(header, ",") {
		candidate = strings.TrimSpace(candidate)
		candidate = strings.TrimPrefix(candidate, "w/")
		if candidate == etag {
			return true
		}
	}
	return false
}
