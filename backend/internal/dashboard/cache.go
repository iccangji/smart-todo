package dashboard

import (
	"sync"
	"time"
)

type Cache interface {
	Get(key string) ([]string, bool)
	Set(key string, value []string, ttl time.Duration)
	Delete(key string)
}

type CacheItem struct {
	Value     []string
	ExpiresAt time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	store map[string]CacheItem
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		store: make(map[string]CacheItem),
	}
}

func (c *MemoryCache) Get(key string) ([]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.store[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Value, true
}

func (c *MemoryCache) Set(key string, value []string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
}
