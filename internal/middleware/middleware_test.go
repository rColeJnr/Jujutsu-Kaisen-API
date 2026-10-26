package middleware

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

var ok = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

// ---------- request id ----------

func TestRequestIDGeneratedWhenMissing(t *testing.T) {
	rec := httptest.NewRecorder()
	RequestID(ok).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	id := rec.Header().Get(HeaderRequestID)
	if !regexp.MustCompile(`^[0-9a-f]{32}$`).MatchString(id) {
		t.Fatalf("expected a 32 char hex id, got %q", id)
	}
}

func TestRequestIDReusedWhenValid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(HeaderRequestID, "trace-abc_123.x")
	rec := httptest.NewRecorder()

	var seen string
	RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = FromContext(r.Context())
	})).ServeHTTP(rec, req)

	if got := rec.Header().Get(HeaderRequestID); got != "trace-abc_123.x" {
		t.Fatalf("response header = %q", got)
	}
	if seen != "trace-abc_123.x" {
		t.Fatalf("context id = %q", seen)
	}
}

func TestRequestIDReplacedWhenInvalid(t *testing.T) {
	for _, bad := range []string{"has space", "semi;colon", strings.Repeat("a", 65)} {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(HeaderRequestID, bad)
		rec := httptest.NewRecorder()
		RequestID(ok).ServeHTTP(rec, req)

		if got := rec.Header().Get(HeaderRequestID); got == bad || got == "" {
			t.Fatalf("invalid id %q should be replaced, got %q", bad, got)
		}
	}
}

// ---------- rate limiter ----------

func newTestLimiter(limit rate.Limit, burst int, header string) (*RateLimiter, *time.Time) {
	now := time.Unix(1_700_000_000, 0)
	rl := NewRateLimiter(limit, burst, header)
	rl.now = func() time.Time { return now }
	return rl, &now
}

func do(h http.Handler, remote string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = remote
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestRateLimitBlocksAfterBurstAndRecovers(t *testing.T) {
	rl, now := newTestLimiter(1, 2, "") // 1 req/s, burst of 2
	h := rl.Middleware(ok)

	for i := 0; i < 2; i++ {
		if rec := do(h, "1.2.3.4:5000", nil); rec.Code != http.StatusOK {
			t.Fatalf("request %d: got %d, want 200", i+1, rec.Code)
		}
	}

	rec := do(h, "1.2.3.4:5000", nil)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("got %d, want 429", rec.Code)
	}
	if rec.Header().Get("Retry-After") != "1" {
		t.Fatalf("Retry-After = %q, want 1", rec.Header().Get("Retry-After"))
	}

	*now = now.Add(1100 * time.Millisecond) // one token refilled
	if rec := do(h, "1.2.3.4:5000", nil); rec.Code != http.StatusOK {
		t.Fatalf("after refill got %d, want 200", rec.Code)
	}
}

func TestRateLimitIsPerClient(t *testing.T) {
	rl, _ := newTestLimiter(1, 1, "")
	h := rl.Middleware(ok)

	if do(h, "1.1.1.1:1", nil).Code != http.StatusOK {
		t.Fatal("first client's first request should pass")
	}
	if do(h, "1.1.1.1:2", nil).Code != http.StatusTooManyRequests {
		t.Fatal("same IP on another port must share the bucket")
	}
	if do(h, "2.2.2.2:1", nil).Code != http.StatusOK {
		t.Fatal("a different IP must have its own bucket")
	}
}

func TestRateLimitTrustedHeader(t *testing.T) {
	rl, _ := newTestLimiter(1, 1, "Fly-Client-IP")
	h := rl.Middleware(ok)
	proxy := "10.0.0.1:443" // every request arrives from the proxy

	if do(h, proxy, map[string]string{"Fly-Client-IP": "9.9.9.9"}).Code != http.StatusOK {
		t.Fatal("first client should pass")
	}
	if do(h, proxy, map[string]string{"Fly-Client-IP": "9.9.9.9"}).Code != http.StatusTooManyRequests {
		t.Fatal("same real client behind the proxy should be limited")
	}
	if do(h, proxy, map[string]string{"Fly-Client-IP": "8.8.8.8"}).Code != http.StatusOK {
		t.Fatal("another real client behind the proxy should pass")
	}
	// garbage in the header falls back to the TCP address instead of crashing
	if do(h, "7.7.7.7:1", map[string]string{"Fly-Client-IP": "not-an-ip"}).Code != http.StatusOK {
		t.Fatal("invalid header should fall back to RemoteAddr")
	}
}

func TestRateLimitSweepsIdleClients(t *testing.T) {
	rl, now := newTestLimiter(10, 10, "")
	rl.allow("a")
	rl.allow("b")
	*now = now.Add(idleTTL + 2*sweepEvery)
	rl.allow("c") // triggers the sweep

	if _, stillThere := rl.clients["a"]; stillThere {
		t.Fatal("idle client should have been removed")
	}
	if len(rl.clients) != 1 {
		t.Fatalf("expected only the active client, got %d", len(rl.clients))
	}
}
