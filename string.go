package main

import (
	"fmt"
	"strconv"
	"time"
)

func (c *Cache) Set(key string, value string, ttl int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry := Entry{
		Value: &StringValue{
			Value: value,
		},
	}

	if ttl > 0 {
		entry.Expiry = time.Now().Add(time.Duration(ttl) * time.Second)
	}

	c.data[key] = entry
}

func (c *Cache) Get(key string) (string, bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)
	if !exists {
		return "", false, nil
	}

	str, ok := entry.Value.(*StringValue)
	if !ok {
		return "", true, ErrWrongType
	}

	return str.Value, true, nil
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
func (c *Cache) MGet(keys []string) ([]Response, error) {
	responses := make([]Response, len(keys))

	for i, key := range keys {
		value, exists, err := c.Get(key)
		if err != nil {
			return nil, err
		}
		if exists {
			responses[i] = NewBulkString(value)
		} else {
			responses[i] = NewNullBulkString()
		}
	}

	return responses, nil
}
func (c *Cache) MSet(keys []string, values []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := 0; i < len(keys); i++ {
		c.data[keys[i]] = Entry{
			Value: &StringValue{
				Value: values[i],
			},
		}
	}
}
func (c *Cache) Append(key string, value string) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	if !exists {
		c.data[key] = Entry{
			Value: &StringValue{
				Value: value,
			},
		}
		return len(value)
	}
	str, err := getString(entry)
	if err != nil {
		return 0
	}

	str.Value += value

	c.data[key] = *entry

	return len(str.Value)
}
func (c *Cache) Strlen(key string) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)
	if !exists {
		return 0
	}
	str, err := getString(entry)
	if err != nil {
		return 0
	}

	return len(str.Value)
}
func (c *Cache) GetSet(key string, value string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	if !exists {
		c.data[key] = Entry{
			Value: &StringValue{
				Value: value,
			},
		}
		return "", false
	}
	str, err := getString(entry)
	if err != nil {
		return "", false
	}

	oldValue := str.Value
	c.data[key] = Entry{
		Value: &StringValue{
			Value: value,
		},
	}
	return oldValue, true
}

func (c *Cache) changeBy(key string, delta int) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	// Key doesn't exist
	if !exists {
		c.data[key] = Entry{
			Value: &StringValue{
				Value: strconv.Itoa(delta),
			},
		}
		return delta, nil
	}
	str, err := getString(entry)
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(str.Value)

	if err != nil {
		return 0, fmt.Errorf("value is not an integer")
	}

	value += delta

	str.Value = strconv.Itoa(value)
	c.data[key] = *entry

	return value, nil
}
func getString(entry *Entry) (*StringValue, error) {
	str, ok := entry.Value.(*StringValue)
	if !ok {
		return nil, ErrWrongType
	}
	return str, nil
}
