package limiter

import (
	"time"

	"github.com/Myspheet/go-rate-limiter/pkg/bucket"
)

type FixedWindowConfig struct {
	WindowDuration time.Duration // duration of the window, seconds, minutes etc
	WindowTokens   int           // number of tokens per window
	WindowSize     int           // number to multiply the duration by
}

type FixedWindowLimiter struct {
	bucket         bucket.Bucket[bucket.FixedWindowBucketType]
	WindowDuration time.Duration
	WindowSize     int
	WindowTokens   int
}

func init() {
	RegisterLimiter("fixed_window", func(cfg map[string]any) Limiter {
		return NewFixedWindowLimiter(bucket.NewInMemoryBucket[bucket.FixedWindowBucketType](), FixedWindowConfig{
			WindowDuration: cfg["window_duration"].(time.Duration),
			WindowTokens:   cfg["window_tokens"].(int),
			WindowSize:     cfg["window_size"].(int),
		})
	})
}

func NewFixedWindowLimiter(fwBucket bucket.Bucket[bucket.FixedWindowBucketType], fwConfig FixedWindowConfig) *FixedWindowLimiter {

	if fwConfig.WindowSize == 0 {
		fwConfig.WindowSize = 1
	}

	if fwConfig.WindowTokens == 0 {
		fwConfig.WindowTokens = 5
	}

	if fwConfig.WindowDuration == 0 {
		fwConfig.WindowDuration = time.Second
	}

	return &FixedWindowLimiter{
		bucket:         fwBucket,
		WindowDuration: fwConfig.WindowDuration,
		WindowTokens:   fwConfig.WindowTokens,
		WindowSize:     fwConfig.WindowSize,
	}
}

func (f *FixedWindowLimiter) Allow(key string) bool {
	currentWindow := getCurrentWindow(time.Duration(f.WindowSize) * f.WindowDuration)

	// check if the key exists
	fw := f.bucket.Get(key)
	if fw == nil {
		fw = &bucket.FixedWindowBucketType{
			CurrentWindow: currentWindow,
			WindowTokens:  f.WindowTokens,
			Capacity:      f.WindowTokens,
		}
		f.bucket.Set(key, fw)
	}

	// check if it's in the current window
	if fw.CurrentWindow == currentWindow {
		// check if there are tokens left
		if fw.WindowTokens > 0 {
			fw.WindowTokens--
			f.bucket.Set(key, fw)
			return true
		}
	} else {
		fw.CurrentWindow = currentWindow
		fw.WindowTokens = f.WindowTokens - 1
		f.bucket.Set(key, fw)
		return true
	}

	return false
}

func getCurrentWindow(windowSize time.Duration) int64 {
	now := time.Now().Unix()
	return now / int64(windowSize.Seconds())
}
