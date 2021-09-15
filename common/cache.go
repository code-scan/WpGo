package common

type Cache struct {
	cache map[string]bool
}

func (c *Cache) Set(key string, val bool) {
	c.cache[key] = val
}
func (c *Cache) Get(key string) bool {
	return c.cache[key]
}
