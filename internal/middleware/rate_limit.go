package middleware

import (
	"net/http"

	"github.com/Myspheet/go-rate-limiter/pkg/limiter"
)

type RateLimiter struct {
	rlimiter limiter.Limiter
}

func NewRateLimiter(rlimiter limiter.Limiter) *RateLimiter {
	return &RateLimiter{
		rlimiter: rlimiter,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("Rate limiter middleware %s", r.RemoteAddr)
		if !rl.rlimiter.Allow(r.RemoteAddr) {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
