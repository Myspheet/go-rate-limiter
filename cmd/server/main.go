package main

import (
	"fmt"
	"net/http"

	"github.com/Myspheet/go-rate-limiter/internal/middleware"
	"github.com/Myspheet/go-rate-limiter/pkg/bucket"
	"github.com/Myspheet/go-rate-limiter/pkg/limiter"
)

func main() {
	b := bucket.NewInMemoryBucket()
	limiter := limiter.NewTokenBucketLimiter(b, limiter.BucketConfig{
		Capacity:   5,
		RefillRate: 0.5,
		Tokens:     5,
	})

	rl := middleware.NewRateLimiter(limiter)

	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})

	http.Handle("/", rl.Middleware(helloHandler))
	http.ListenAndServe(":8080", nil)
}
