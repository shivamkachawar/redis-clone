package cache

import "time"

func (c *Cache) evictLRU() {

	if len(c.data) == 0 {
		return
	}

	var victimKey string
	var oldest time.Time
	first := true

	for key, entry := range c.data {

		if first {
			victimKey = key
			oldest = entry.LastAccess
			first = false
			continue
		}

		if entry.LastAccess.Before(oldest) {
			victimKey = key
			oldest = entry.LastAccess
		}
	}

	delete(c.data, victimKey)
}

func (c *Cache) checkEviction() {
	if len(c.data) >= c.maxKeys {
		c.evictLRU()
	}
}
