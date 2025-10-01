package limiter

import "errors"

type Limiter interface {
	Allow(key string) bool // check if request is allowed
}

// used to create the rate limiter with the default config
type LimiterFactory func(cfg map[string]any) Limiter

var limiterRegistry = make(map[string]LimiterFactory)

// Creates a new rate limiter
func NewRateLimiter(name string, cfg map[string]interface{}) (Limiter, error) {
	if constructor, ok := limiterRegistry[name]; ok {
		return constructor(cfg), nil
	}

	return nil, errors.New("Invalid Limiter: " + name)
}

// Registers a rate limiter with it's config
func RegisterLimiter(name string, factory LimiterFactory) error {
	// if it exists, just return an error
	if _, ok := limiterRegistry[name]; ok {
		return errors.New("Limiter already registered: " + name)
	}

	limiterRegistry[name] = factory
	return nil
}
