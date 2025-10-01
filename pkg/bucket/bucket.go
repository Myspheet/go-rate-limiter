package bucket

import (
	"time"
)

type TokenBucketType struct {
	Capacity   int       // total number of tokens
	RefillRate float64   // tokens per second
	Tokens     int       // number of tokens left
	LastRefill time.Time // last time the bucket was refilled
}

type FixedWindowBucketType struct {
	CurrentWindow int64 // current window
	WindowTokens  int   // number of tokens left in the current window
	Capacity      int   // total number of tokens in the window
}

type SlidingWindowLogBucketType struct {
	WindowLog []time.Time
	Capacity  int
}

type AllowedTypes interface {
	TokenBucketType | FixedWindowBucketType | SlidingWindowLogBucketType
}

type Bucket[T AllowedTypes] interface {
	Get(key string) *T
	Set(key string, bucket *T) error
	Delete(key string) error
	Clear()
}
