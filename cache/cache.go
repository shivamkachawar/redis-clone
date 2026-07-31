package cache

import (
	"go-redis/protocol"
	"time"
)

type Entry struct {
	Value      protocol.RedisValue
	Expiry     time.Time
	LastAccess time.Time
}

type Cache struct {
	data    map[string]Entry
	maxKeys int
}

func NewCache(maxKeys int) *Cache {
	return &Cache{
		data:    make(map[string]Entry),
		maxKeys: maxKeys,
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
	entry.LastAccess = time.Now()
	c.data[key] = entry

	return &entry, true
}

// Client APIs

func (c *Cache) Entries() map[string]Entry {

	entries := make(map[string]Entry)

	for key := range c.data {
		entry, exists := c.getEntry(key)
		if exists {
			entries[key] = *entry
		}
	}

	return entries
}
