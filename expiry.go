package main

import "time"

func (c *Cache) Expire(key string, seconds int) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)
	if !exists {
		return false
	}

	entry.Expiry = time.Now().Add(time.Duration(seconds) * time.Second)
	c.data[key] = *entry

	return true
}

func (c *Cache) Persist(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)
	if !exists {
		return false
	}

	entry.Expiry = time.Time{}
	c.data[key] = *entry

	return true
}

func (c *Cache) TTL(key string) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)
	if !exists {
		return -2
	}

	if entry.Expiry.IsZero() {
		return -1
	}

	return int(time.Until(entry.Expiry).Seconds())
}
