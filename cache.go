package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Entry struct {
	Value  string
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
func (c *Cache) changeBy(key string, delta int) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	// Key doesn't exist
	if !exists {
		c.data[key] = Entry{
			Value: strconv.Itoa(delta),
		}
		return delta, nil
	}

	value, err := strconv.Atoi(entry.Value)

	if err != nil {
		return 0, fmt.Errorf("value is not an integer")
	}

	value += delta

	entry.Value = strconv.Itoa(value)
	c.data[key] = *entry

	return value, nil
}

// Client APIs

func (c *Cache) Set(key string, value string, ttl int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry := Entry{
		Value: value,
	}

	if ttl > 0 {
		entry.Expiry = time.Now().Add(time.Duration(ttl) * time.Second)
	}

	c.data[key] = entry
}

func (c *Cache) Get(key string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)
	if !exists {
		return "", false
	}

	return entry.Value, true
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

func (c *Cache) Exists(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, exists := c.getEntry(key)
	return exists
}

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
func (c *Cache) Increment(key string) (int, error) {
	return c.changeBy(key, 1)
}
func (c *Cache) Decrement(key string) (int, error) {
	return c.changeBy(key, -1)
}
func (c *Cache) IncrementBy(key string, delta int) (int, error) {
	return c.changeBy(key, delta)
}
func (c *Cache) DecrementBy(key string, delta int) (int, error) {
	return c.changeBy(key, -delta)
}
func (c *Cache) MGet(keys []string) []string {
	values := make([]string, len(keys))

	for i, key := range keys {
		value, exists := c.Get(key)

		if exists {
			values[i] = value
		} else {
			values[i] = "(nil)"
		}
	}

	return values
}
func (c *Cache) MSet(keys []string, values []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := 0; i < len(keys); i++ {
		c.data[keys[i]] = Entry{
			Value: values[i],
		}
	}
}
func (c *Cache) Append(key string, value string) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	if !exists {
		c.data[key] = Entry{
			Value: value,
		}
		return len(value)
	}
	entry.Value += value
	c.data[key] = *entry

	return len(entry.Value)
}
func (c *Cache) Strlen(key string) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)
	if !exists {
		return 0
	}
	return len(entry.Value)
}
func (c *Cache) GetSet(key string, value string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	if !exists {
		c.data[key] = Entry{
			Value: value,
		}
		return "", false
	}
	oldValue := entry.Value
	c.data[key] = Entry{
		Value: value,
	}
	return oldValue, true
}
