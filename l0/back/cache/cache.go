package cache

type Cache struct {
	storage map[string]string
}

func InitCache() *Cache {
	m := map[string]string{}
	return &Cache{m}
}

func (c *Cache) Get(key string) (string, bool) {
	val, ok := c.storage[key]
	return val, ok
}

func (c *Cache) Has(key string) bool {
	_, existed := c.storage[key]
	return existed
}

func (c *Cache) Set(key string, val string) {
	c.storage[key] = val
}
