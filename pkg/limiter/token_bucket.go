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
