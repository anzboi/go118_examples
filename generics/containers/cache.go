package containers

import "sync"

type Cache[T any] interface {
	Put(T)
	Get() T
}

type cache[T any] struct {
	sync.RWMutex
	value T
}

func NewCache[T any](value T) Cache[T] {
	return &cache[T]{
		value: value,
	}
}

func (c *cache[T]) Put(value T) {
	c.Lock()
	defer c.Unlock()
	c.value = value
}

func (c *cache[T]) Get() T {
	c.RLock()
	defer c.RUnlock()
	return c.value
}
