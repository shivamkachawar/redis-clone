package main

import (
	"sync"
	"time"
)

type Entry struct {
	Value  RedisValue
	Expiry time.Time
}

type Cache struct {
	data  map[string]Entry
	mutex sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]Entry),
	}
}

// --------------------
// Private Helper
// --------------------

func (c *Cache) getEntry(key string) (*Entry, bool) {
	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if !entry.Expiry.IsZero() && time.Now().After(entry.Expiry) {
		delete(c.data, key)
		return nil, false
	}

	return &entry, true
}

// Client APIs

func (c *Cache) Entries() map[string]Entry {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entries := make(map[string]Entry)

	for key := range c.data {
		entry, exists := c.getEntry(key)
		if exists {
			entries[key] = *entry
		}
	}

	return entries
}
