package limiter

import (
	"testing"
	"time"

	"github.com/Myspheet/go-rate-limiter/pkg/bucket"
)

type mockFixedWindowBucket struct {
	store map[string]*bucket.FixedWindowBucketType
}

func (m *mockFixedWindowBucket) Get(key string) *bucket.FixedWindowBucketType {
	return m.store[key]
}

func (m *mockFixedWindowBucket) Set(key string, bucket *bucket.FixedWindowBucketType) error {
	m.store[key] = bucket
	return nil
}

func (m *mockFixedWindowBucket) Delete(key string) error {
	delete(m.store, key)
	return nil
}

func (m *mockFixedWindowBucket) Clear() {
	m.store = make(map[string]*bucket.FixedWindowBucketType)
}

func TestNewFixedWindowLimiter_Defaults(t *testing.T) {
	mockBucket := &mockFixedWindowBucket{store: make(map[string]*bucket.FixedWindowBucketType)}
	limiter := NewFixedWindowLimiter(mockBucket, FixedWindowConfig{})

	if limiter.WindowSize != 1 {
		t.Errorf("expected WindowSize to be 1, got %d", limiter.WindowSize)
	}

	if limiter.WindowTokens != 5 {
		t.Errorf("expected WindowTokens to be 5, got %d", limiter.WindowTokens)
	}

	if limiter.WindowDuration != time.Second {
		t.Errorf("expected WindowDuration to be 1s, got %s", limiter.WindowDuration)
	}
}

func TestFixedWindowLimiter_Allow_NewKey(t *testing.T) {
	mockBucket := &mockFixedWindowBucket{store: make(map[string]*bucket.FixedWindowBucketType)}
	limiter := NewFixedWindowLimiter(mockBucket, FixedWindowConfig{})

	key := "newkey"
	allowed := limiter.Allow(key)
	if !allowed {
		t.Errorf("expected Allow to return true for new key")
	}

	tb := mockBucket.Get(key)
	if tb == nil {
		t.Errorf("expected token bucket to be created")
	}

	if tb != nil && tb.WindowTokens != 4 {
		t.Errorf("expected tokens=4 for new key, got %d", tb.WindowTokens)
	}

	if tb.CurrentWindow == 0 {
		t.Errorf("expected CurrentWindow to be set for new key, got %d", tb.CurrentWindow)
	}
}

func TestFixedWindowLimiter_Allow_Burst(t *testing.T) {
	mockBucket := &mockFixedWindowBucket{store: make(map[string]*bucket.FixedWindowBucketType)}
	limiter := NewFixedWindowLimiter(mockBucket, FixedWindowConfig{})

	key := "burstkey"

	// First 5 calls should be allowed (initial allow without decrement, then 4 decrements)
	for i := 1; i <= 5; i++ {
		allowed := limiter.Allow(key)
		if !allowed {
			t.Errorf("expected Allow to return true on call %d", i)
		}
	}

	tb := mockBucket.Get(key)
	if tb.WindowTokens != 0 {
		t.Errorf("expected tokens=0 after burst, got %d", tb.WindowTokens)
	}

	// Next call should be denied
	allowed := limiter.Allow(key)
	if allowed {
		t.Errorf("expected Allow to return false after burst")
	}
	if tb.WindowTokens != 0 {
		t.Errorf("expected tokens still=0, got %d", tb.WindowTokens)
	}
}
