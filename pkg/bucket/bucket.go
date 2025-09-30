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

type Bucket interface {
	Get(key string) *TokenBucketType
	Set(key string, bucket *TokenBucketType) error
}
