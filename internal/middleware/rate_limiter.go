package middleware

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	sweepEvery = time.Minute
	idleTTL    = 10 * time.Minute
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter is a per-client token bucket. Each client (identified by IOP)
// gets `burst` request at once and then `limit` request per second.
type RateLimiter struct {
	sync.Mutex
	clients     map[string]*client
	limit       rate.Limit
	burst       int
	trustHeader string
	retryAfter  int
	lastSweep   time.Time
	now         func() time.Time
}

// NewRateLimiter creates a limiter allowing `limit` requests per second with
// burst up to `burst`. trustHeader is the header your reverse proxy sets with
// the real client IP; leave it empty when the server is exposed directly,
// so the TCP peer address is used.
func NewRateLimiter(limit rate.Limit, burst int, trustHeader string) *RateLimiter {
	retry := 1
	if limit > 0 {
		retry = int(math.Ceil(1 / float64(limit)))
	}
	if retry < 1 {
		retry = 1
	}
	return &RateLimiter{
		clients:     make(map[string]*client),
		limit:       limit,
		burst:       burst,
		trustHeader: trustHeader,
		retryAfter:  retry,
		now:         time.Now,
	}
}

// Middleware rejects request above the limit with 429 and a Retry-After header.
func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.allow(ClientIP(r, rl.trustHeader)) {
			w.Header().Set("Retry-After", strconv.Itoa(rl.retryAfter))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"code":    http.StatusTooManyRequests,
				"message": fmt.Sprintf("rate limit exceeded, retry in %ds", rl.retryAfter),
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(key string) bool {
	rl.Lock()
	defer rl.Unlock()

	now := rl.now()

	if now.Sub(rl.lastSweep) > sweepEvery {
		for k, c := range rl.clients {
			if now.Sub(c.lastSeen) > idleTTL {
				delete(rl.clients, k)
			}
		}
		rl.lastSweep = now
	}

	c, ok := rl.clients[key]
	if !ok {
		c = &client{limiter: rate.NewLimiter(rl.limit, rl.burst)}
		rl.clients[key] = c
	}
	c.lastSeen = now

	return c.limiter.AllowN(now, 1)
}

// ClientIP returns the caller's IP.
func ClientIP(r *http.Request, trustHeader string) string {
	if trustHeader != "" {
		if v := r.Header.Get(trustHeader); v != "" {
			if i := strings.LastIndex(v, ","); i >= 0 {
				v = v[i+1:]
			}
			if ip := net.ParseIP(strings.TrimSpace(v)); ip != nil {
				return ip.String()
			}
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
