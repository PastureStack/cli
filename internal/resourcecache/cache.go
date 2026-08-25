package resourcecache

import (
	"sync"
	"time"
)

type entry struct {
	value     any
	expiresAt time.Time
}

// Cache is the small expiring resource cache needed by the event monitor.
// Expired entries are removed on access, so no background goroutine is needed.
type Cache struct {
	mu      sync.RWMutex
	ttl     time.Duration
	entries map[string]entry
}

func New(ttl time.Duration) *Cache {
	return &Cache{ttl: ttl, entries: make(map[string]entry)}
}

func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	c.entries[key] = entry{value: value, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	item, ok := c.entries[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().Before(item.expiresAt) {
		return item.value, true
	}
	c.mu.Lock()
	if current, exists := c.entries[key]; exists && !time.Now().Before(current.expiresAt) {
		delete(c.entries, key)
	}
	c.mu.Unlock()
	return nil, false
}
