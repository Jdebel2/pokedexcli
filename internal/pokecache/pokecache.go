package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cache map[string]CacheEntry
	mu    sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	duration := interval * time.Nanosecond
	var c = Cache{
		cache: map[string]CacheEntry{},
		mu:    sync.Mutex{},
	}
	go c.reapLoop(duration)
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) (val []byte, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cache[key]
	if ok {
		return entry.val, ok
	}
	return []byte{}, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	tick := time.NewTicker(interval)
	for {
		select {
		case <-tick.C:
			t := time.Now()
			c.mu.Lock()
			for key := range c.cache {
				if t.Sub(c.cache[key].createdAt) > interval {
					delete(c.cache, key)
				}
			}
			c.mu.Unlock()
		}
	}
}
