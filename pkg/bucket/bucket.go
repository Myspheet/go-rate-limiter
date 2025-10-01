package bucket

import (
	"time"
)

type TokenBucketType struct {
	Capacity   int
	RefillRate float64 // tokens per second
	Tokens     int
	LastRefill time.Time
}

type FixedWindowBucketType struct {
	CurrentWindow int64
	WindowTokens  int
	Capacity      int
}

type Bucket[T any] interface {
	Get(key string) *T
	Set(key string, bucket *T) error
	Delete(key string) error
	Clear()
}
