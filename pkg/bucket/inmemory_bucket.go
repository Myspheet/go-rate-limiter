package bucket

type InMemoryBucket struct {
	buckets map[string]*TokenBucketType
}

func NewInMemoryBucket() *InMemoryBucket {
	return &InMemoryBucket{
		buckets: make(map[string]*TokenBucketType),
	}
}

func (b *InMemoryBucket) Get(key string) *TokenBucketType {
	return b.buckets[key]
}

func (b *InMemoryBucket) Set(key string, bucket *TokenBucketType) error {
	b.buckets[key] = bucket
	return nil
}
