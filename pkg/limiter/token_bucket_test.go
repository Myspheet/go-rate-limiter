package limiter

// test to write
// Check that a key is allowed if it doesn't exist
// check that a key is allowed if it exists but not up to the limit

import (
	"testing"
	"time"

	"github.com/Myspheet/go-rate-limiter/pkg/bucket"
)

type mockBucket struct {
	store map[string]*bucket.TokenBucketType
}

func (m *mockBucket) Get(key string) *bucket.TokenBucketType {
	return m.store[key]
}

func (m *mockBucket) Set(key string, tb *bucket.TokenBucketType) error {
	m.store[key] = tb
	return nil
}

func TestTokenBucketLimiter_Allow_NewKey(t *testing.T) {
	mockB := &mockBucket{store: make(map[string]*bucket.TokenBucketType)}

	limiter := NewTokenBucketLimiter(mockB, BucketConfig{
		Capacity:   5,
		RefillRate: 1,
		Tokens:     5,
	})

	key := "newkey"
	allowed := limiter.Allow(key)
	if !allowed {
		t.Errorf("expected Allow to return true for new key")
	}

	tb := mockB.Get(key)
	if tb == nil {
		t.Errorf("expected token bucket to be created")
	}
}

func TestTokenBucketLimiter_Allow_Burst(t *testing.T) {
	mockB := &mockBucket{store: make(map[string]*bucket.TokenBucketType)}
	limiter := NewTokenBucketLimiter(mockB, BucketConfig{
		Capacity:   5,
		RefillRate: 1,
		Tokens:     5,
	})

	key := "burstkey"

	// First 5 calls should be allowed (initial allow without decrement, then 5 decrements)
	for i := 1; i <= 5; i++ {
		allowed := limiter.Allow(key)
		if !allowed {
			t.Errorf("expected Allow to return true on call %d", i)
		}
	}

	tb := mockB.Get(key)
	if tb.Tokens != 0 {
		t.Errorf("expected tokens=0 after burst, got %d", tb.Tokens)
	}

	// Next call should be denied
	allowed := limiter.Allow(key)
	if allowed {
		t.Errorf("expected Allow to return false after burst")
	}
	if tb.Tokens != 0 {
		t.Errorf("expected tokens still=0, got %d", tb.Tokens)
	}
}

func TestTokenBucketLimiter_Allow_Refill(t *testing.T) {
	mockB := &mockBucket{store: make(map[string]*bucket.TokenBucketType)}
	limiter := NewTokenBucketLimiter(mockB, BucketConfig{
		Capacity:   5,
		RefillRate: 1,
		Tokens:     5,
	})

	key := "refillkey"

	// Exhaust the bucket
	for i := 1; i <= 5; i++ {
		limiter.Allow(key)
	}

	tb := mockB.Get(key)
	if tb.Tokens != 0 {
		t.Errorf("expected tokens=0 after exhaust, got %d", tb.Tokens)
	}

	// Simulate 2 seconds elapsed by adjusting LastRefill
	now := time.Now()
	tb.LastRefill = now.Add(-2 * time.Second)

	// Call Allow: should add 2 tokens, cap at 5, then decrement to 1
	allowed := limiter.Allow(key)
	if !allowed {
		t.Errorf("expected Allow to return true after refill")
	}
	if tb.Tokens != 1 {
		t.Errorf("expected tokens=1 after refill and decrement, got %d", tb.Tokens)
	}
	if tb.LastRefill.Sub(now) > time.Millisecond {
		t.Errorf("expected LastRefill close to now, got %v", tb.LastRefill)
	}
}

func TestTokenBucketLimiter_Allow_RefillExceedsCapacity(t *testing.T) {
	mockB := &mockBucket{store: make(map[string]*bucket.TokenBucketType)}
	limiter := NewTokenBucketLimiter(mockB, BucketConfig{
		Capacity:   5,
		RefillRate: 1,
		Tokens:     5,
	})

	key := "exceedkey"

	// Exhaust the bucket
	for i := 1; i <= 6; i++ {
		limiter.Allow(key)
	}

	tb := mockB.Get(key)
	if tb.Tokens != 0 {
		t.Errorf("expected tokens=0 after exhaust, got %d", tb.Tokens)
	}

	// Simulate 10 seconds elapsed
	now := time.Now()
	tb.LastRefill = now.Add(-10 * time.Second)

	// Call Allow: should add 10 tokens, cap at 5, then decrement to 4
	allowed := limiter.Allow(key)
	if !allowed {
		t.Errorf("expected Allow to return true after large refill")
	}
	if tb.Tokens != 4 {
		t.Errorf("expected tokens=4 after capped refill and decrement, got %d", tb.Tokens)
	}
	if tb.LastRefill.Sub(now) > time.Millisecond {
		t.Errorf("expected LastRefill close to now, got %v", tb.LastRefill)
	}
}

func TestTokenBucketLimiter_Allow_NoRefillIfZeroAdded(t *testing.T) {
	mockB := &mockBucket{store: make(map[string]*bucket.TokenBucketType)}
	limiter := NewTokenBucketLimiter(mockB, BucketConfig{
		Capacity:   5,
		RefillRate: 1,
		Tokens:     5,
	})

	key := "norefillkey"

	// Exhaust the bucket
	for i := 1; i <= 6; i++ {
		limiter.Allow(key)
	}

	tb := mockB.Get(key)
	if tb.Tokens != 0 {
		t.Errorf("expected tokens=0 after exhaust, got %d", tb.Tokens)
	}

	// Simulate 0.4 seconds elapsed (addedTokens = int(0.4 * 1) = 0)
	now := time.Now()
	tb.LastRefill = now.Add(-400 * time.Millisecond)
	lastRefill := tb.LastRefill

	// Call Allow: no add, tokens=0, false, LastRefill not updated
	allowed := limiter.Allow(key)
	if allowed {
		t.Errorf("expected Allow to return false with no refill")
	}
	if tb.Tokens != 0 {
		t.Errorf("expected tokens=0, got %d", tb.Tokens)
	}
	if !tb.LastRefill.Equal(lastRefill) {
		t.Errorf("expected LastRefill not updated when addedTokens=0")
	}
}

func TestTokenBucketLimiter_Allow_FractionalRefill(t *testing.T) {
	mockB := &mockBucket{store: make(map[string]*bucket.TokenBucketType)}
	limiter := NewTokenBucketLimiter(mockB, BucketConfig{
		Capacity:   5,
		RefillRate: 1,
		Tokens:     5,
	})

	key := "fractionalkey"

	// Exhaust the bucket
	for i := 1; i <= 6; i++ {
		limiter.Allow(key)
	}

	tb := mockB.Get(key)
	if tb.Tokens != 0 {
		t.Errorf("expected tokens=0 after exhaust, got %d", tb.Tokens)
	}

	// Simulate 1.7 seconds elapsed (addedTokens = int(1.7 * 1) = 1)
	now := time.Now()
	tb.LastRefill = now.Add(-1700 * time.Millisecond)

	// Call Allow: add 1, tokens=1, decrement to 0, true
	allowed := limiter.Allow(key)
	if !allowed {
		t.Errorf("expected Allow to return true after fractional refill")
	}
	if tb.Tokens != 0 {
		t.Errorf("expected tokens=0 after refill and decrement, got %d", tb.Tokens)
	}
	if tb.LastRefill.Sub(now) > time.Millisecond {
		t.Errorf("expected LastRefill close to now, got %v", tb.LastRefill)
	}
}

// type mockBucket struct {
// 	store map[string]*bucket.TokenBucketType
// }

// func (m *mockBucket) Get(key string) *bucket.TokenBucketType {
// 	return m.store[key]
// }

// func (m *mockBucket) Set(key string, bucket *bucket.TokenBucketType) error {
// 	m.store[key] = bucket
// 	return nil
// }

// func TestTokenBucketLimiter_Allow_NewKey(t *testing.T) {
// 	mockDB := &mockBucket{store: make(map[string]*bucket.TokenBucketType)}
// 	limiter := NewTokenBucketLimiter(mockDB, BucketConfig{
// 		Capacity:   5,
// 		RefillRate: 1,
// 		Tokens:     5,
// 	})

// 	allowed := limiter.Allow("newkey")
// 	if !allowed {
// 		t.Errorf("expected Allow to return true for new key")
// 	}

// 	tb := mockDB.Get("newkey")
// 	if tb == nil {
// 		t.Errorf("expected TokenBucketType to be created for new key")
// 	}
// }
