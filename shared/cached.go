package shared

import "time"

type ErrCacheExpired struct{}

func (ErrCacheExpired) Error() string {
	return "cached value expired"
}

type Cached[T any] struct {
	V         T
	ExpiresAt time.Time
}

func Cache[T any](v T, d time.Duration) Cached[T] {
	return Cached[T]{v, time.Now().Add(d)}
}

func (c Cached[T]) Get() (t T, err error) {
	if c.IsExpired() {
		return t, ErrCacheExpired{}
	}
	return c.V, nil
}

func (c Cached[T]) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

func (c *Cached[T]) Expire() {
	c.ExpiresAt = time.Now()
}
