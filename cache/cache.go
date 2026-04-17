package cache

import (
	"sync"
	"time"
)

// CacheItem хранит значение и момент истечения срока жизни записи.
type CacheItem[T any] struct {
	Value T
	Ttl   int64
}

// Cache представляет простой in-memory cache с TTL для каждого значения.
type Cache[T any] struct {
	mu    sync.RWMutex
	items map[string]CacheItem[T]
	ttl   time.Duration
}

// NewCache создаёт новый cache с общим TTL для всех записей.
func NewCache[T any](ttl time.Duration) *Cache[T] {
	return &Cache[T]{
		items: make(map[string]CacheItem[T]),
		ttl:   ttl,
	}
}

// Get возвращает значение по ключу, если запись существует и ещё не истекла.
func (c *Cache[T]) Get(key string) (T, bool) {
	var v T

	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok || time.Now().UnixNano() > item.Ttl {
		return v, false
	}

	return item.Value, true
}

// Evict удаляет запись по ключу.
func (c *Cache[T]) Evict(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// Set сохраняет значение в cache с TTL, заданным при создании cache.
func (c *Cache[T]) Set(key string, value T) {
	c.mu.Lock()
	c.items[key] = CacheItem[T]{
		Value: value,
		Ttl:   time.Now().Add(c.ttl).UnixNano(),
	}
	c.mu.Unlock()
}

// Cleanup удаляет из cache все записи, срок жизни которых истёк.
func (c *Cache[T]) Cleanup() {
	now := time.Now().UnixNano()

	c.mu.Lock()
	for k, v := range c.items {
		if now > v.Ttl {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}
