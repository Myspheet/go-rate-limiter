package limiter

import (
	"time"

	"github.com/Myspheet/go-rate-limiter/pkg/bucket"
)

type SlidingWindowLogConfig struct {
	WindowSize     int64
	Capacity       int
	WindowDuration time.Duration
}

type SlidingWindowLogLimiter struct {
	bucket         bucket.Bucket[bucket.SlidingWindowLogBucketType]
	Capacity       int
	WindowSize     int64
	WindowDuration time.Duration
}

func init() {
	RegisterLimiter("sliding_window_log", func(cfg map[string]any) Limiter {
		return NewSlidingWindowLogLimiter(bucket.NewInMemoryBucket[bucket.SlidingWindowLogBucketType](), SlidingWindowLogConfig{
			WindowSize:     cfg["window_size"].(int64),
			Capacity:       cfg["capacity"].(int),
			WindowDuration: cfg["window_duration"].(time.Duration),
		})
	})
}

func NewSlidingWindowLogLimiter(swBucket bucket.Bucket[bucket.SlidingWindowLogBucketType], config SlidingWindowLogConfig) *SlidingWindowLogLimiter {

	if config.WindowSize == 0 {
		config.WindowSize = 1
	}

	if config.WindowDuration == 0 {
		config.WindowDuration = time.Minute
	}

	if config.Capacity == 0 {
		config.Capacity = 5
	}

	return &SlidingWindowLogLimiter{
		bucket:         swBucket,
		Capacity:       config.Capacity,
		WindowSize:     config.WindowSize,
		WindowDuration: config.WindowDuration,
	}
}

func (s *SlidingWindowLogLimiter) Allow(key string) bool {
	// check if key exists
	swl := s.bucket.Get(key)
	if swl == nil {
		swl = &bucket.SlidingWindowLogBucketType{
			WindowLog: make([]time.Time, 0),
		}
		s.bucket.Set(key, swl)
	}

	now := time.Now()
	// check if it's in the current window

	newWindowLog := make([]time.Time, 0)
	for _, t := range swl.WindowLog {
		if t.Add(time.Duration(s.WindowSize) * s.WindowDuration).After(now) {
			newWindowLog = append(newWindowLog, t)
		}
	}

	if len(newWindowLog) < s.Capacity {
		newWindowLog = append(newWindowLog, time.Now())

		swl.WindowLog = newWindowLog
		s.bucket.Set(key, swl)
		return true
	}

	return false
}
