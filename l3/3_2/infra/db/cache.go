package db

import "sync"

type MemoryCache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]string),
	}
}

func (c *MemoryCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *MemoryCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}
