package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Myspheet/go-rate-limiter/internal/middleware"
	"github.com/Myspheet/go-rate-limiter/pkg/limiter"
)

func main() {
	// Create a new TokenBucketLimiter
	// b := bucket.NewInMemoryBucket[bucket.TokenBucketType]()
	// limiter := limiter.NewTokenBucketLimiter(b, limiter.BucketConfig{
	// 	Capacity:   5,
	// 	RefillRate: 0.5,
	// 	Tokens:     5,
	// })

	// Create a new Fixed Window Limiter
	// fw := bucket.NewInMemoryBucket[bucket.FixedWindowBucketType]()
	// limiter := limiter.NewFixedWindowLimiter(fw, limiter.FixedWindowConfig{
	// 	WindowDuration: time.Minute,
	// 	WindowTokens:   5,
	// })

	// swl := bucket.NewInMemoryBucket[bucket.SlidingWindowLogBucketType]()
	// limiter := limiter.NewSlidingWindowLogLimiter(swl, limiter.SlidingWindowLogConfig{
	// 	Capacity:       5,
	// 	WindowSize:     1,
	// 	WindowDuration: time.Minute,
	// })

	// new limiter
	limiter, err := limiter.NewRateLimiter("fixed_window", map[string]interface{}{
		"window_duration": time.Minute,
		"window_tokens":   5,
		"window_size":     1,
	})
	if err != nil {
		panic(err)
	}

	rl := middleware.NewRateLimiter(limiter)

	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})

	http.Handle("/", rl.Middleware(helloHandler))
	http.ListenAndServe(":8080", nil)
}
