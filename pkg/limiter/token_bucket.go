package limiter

import (
	"time"

	"github.com/Myspheet/go-rate-limiter/pkg/bucket"
)

type BucketConfig struct {
	Capacity   int
	RefillRate float64
	Tokens     int
}
type TokenBucketLimiter struct {
	bucket     bucket.Bucket[bucket.TokenBucketType]
	capacity   int
	refillRate float64
	tokens     int
}

func init() {
	RegisterLimiter("token_bucket", func(cfg map[string]any) Limiter {
		return NewTokenBucketLimiter(bucket.NewInMemoryBucket[bucket.TokenBucketType](), BucketConfig{
			Capacity:   cfg["capacity"].(int),
			RefillRate: cfg["refill_rate"].(float64),
			Tokens:     cfg["tokens"].(int),
		})
	})
}

// NewTokenBucketLimiter creates a new TokenBucketLimiter with the given tokenBucket and bucketConfig.
// If the bucketConfig's Capacity, RefillRate, or Tokens is 0, it will be set to the default values of 5, 1, and 5 respectively.
func NewTokenBucketLimiter(tokenBucket bucket.Bucket[bucket.TokenBucketType], bucketConfig BucketConfig) *TokenBucketLimiter {
	if bucketConfig.Capacity == 0 {
		bucketConfig.Capacity = 5
	}

	if bucketConfig.RefillRate == 0 {
		bucketConfig.RefillRate = 1
	}

	if bucketConfig.Tokens == 0 {
		bucketConfig.Tokens = 5
	}
	return &TokenBucketLimiter{
		bucket:     tokenBucket,
		capacity:   bucketConfig.Capacity,
		refillRate: bucketConfig.RefillRate,
		tokens:     bucketConfig.Tokens,
	}
}

// Allow returns true if the key is allowed to be accessed, false otherwise.
// If the key doesn't exist, it will be created and allowed.
// If the key exists, it will check the elapsed time since last refill and add tokens accordingly.
// It will then deduct a token for this request and return true if the key is allowed, false otherwise.
func (tb *TokenBucketLimiter) Allow(key string) bool {
	now := time.Now()
	// get bucket from bucket store
	tokenBucket := tb.bucket.Get(key)

	// the key doesn't exist so create it and allow it
	if tokenBucket == nil {
		tokenBucket = &bucket.TokenBucketType{
			Capacity:   tb.capacity,
			RefillRate: tb.refillRate,
			Tokens:     tb.tokens,
			LastRefill: now,
		}

		tb.bucket.Set(key, tokenBucket)
	}

	// if key exists check the elapsed time since last refill
	elapsed := now.Sub(tokenBucket.LastRefill).Seconds()

	// calculate how many tokens to add
	addedTokens := int(elapsed * float64(tokenBucket.RefillRate))

	if addedTokens > 0 {
		tokenBucket.Tokens = min(tokenBucket.Capacity, tokenBucket.Tokens+addedTokens)
		tokenBucket.LastRefill = now
	}

	// deduct a token for this request
	if tokenBucket.Tokens > 0 {
		tokenBucket.Tokens--

		return true
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
