package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	value    string
	expireAt int64
}

type Cache struct {
	m       *sync.RWMutex
	storage map[string]CacheItem
	ttl     time.Duration
}

type CacheInterface interface {
	Get(key string) (string, bool)
	Has(key string) bool
	Set(key string, val string)
}

func InitCache(ttl time.Duration) *Cache {
	var m sync.RWMutex
	s := make(map[string]CacheItem)
	cache := &Cache{&m, s, ttl}

	go cache.cleanup()

	return cache
}

func (c *Cache) Get(key string) (string, bool) {
	c.m.RLock()
	defer c.m.RUnlock()
	item, ok := c.storage[key]
	if !ok {
		return "", false
	}
	if time.Now().Unix() > item.expireAt {
		return "", false
	}
	return item.value, true
}

func (c *Cache) Has(key string) bool {
	c.m.RLock()
	defer c.m.RUnlock()
	item, existed := c.storage[key]
	if !existed {
		return false
	}
	if time.Now().Unix() > item.expireAt {
		return false
	}
	return true
}

func (c *Cache) Set(key string, val string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.storage[key] = CacheItem{
		value:    val,
		expireAt: time.Now().Add(c.ttl).Unix(),
	}
}

func (c *Cache) cleanup() {
	for {
		time.Sleep(c.ttl)
		c.m.Lock()
		for key, item := range c.storage {
			if time.Now().Unix() > item.expireAt {
				delete(c.storage, key)
			}
		}
		c.m.Unlock()
	}
}
