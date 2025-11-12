package cache

import (
	"sync"
	"time"
)

type item struct {
	value     interface{}
	expiresAt time.Time
}

type Cache struct {
	items map[string]item
	mu    sync.RWMutex
}

func New() *Cache {
	return &Cache{
		items: make(map[string]item),
	}
}

// Set adds an item to the cache with a TTL
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Check if item has expired
	if time.Now().After(item.expiresAt) {
		delete(c.items, key)
		return nil, false
	}

	return item.value, true
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]item)
}
