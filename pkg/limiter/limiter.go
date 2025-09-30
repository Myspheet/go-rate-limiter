package limiter

type Limiter interface {
	Allow(key string) bool // check if request is allowed
}
