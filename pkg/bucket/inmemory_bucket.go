package bucket

type InMemoryBucket[T AllowedTypes] struct {
	buckets map[string]*T
}

func NewInMemoryBucket[T AllowedTypes]() *InMemoryBucket[T] {
	return &InMemoryBucket[T]{
		buckets: make(map[string]*T),
	}
}

func (b *InMemoryBucket[T]) Get(key string) *T {
	return b.buckets[key]
}

func (b *InMemoryBucket[T]) Set(key string, bucket *T) error {
	b.buckets[key] = bucket
	return nil
}

func (b *InMemoryBucket[T]) Delete(key string) error {
	delete(b.buckets, key)
	return nil
}

func (b *InMemoryBucket[T]) Clear() {
	b.buckets = make(map[string]*T)
}
