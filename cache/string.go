package cache

import (
	"fmt"
	"go-redis/protocol"
	"strconv"
	"time"
)

func (c *Cache) Set(key string, value string, ttl int) {

	entry := Entry{
		Value: &protocol.StringValue{
			Value: value,
		},
	}

	if ttl > 0 {
		entry.Expiry = time.Now().Add(time.Duration(ttl) * time.Second)
	}

	c.data[key] = entry
}

func (c *Cache) Get(key string) (string, bool, error) {

	entry, exists := c.getEntry(key)
	if !exists {
		return "", false, nil
	}

	str, ok := entry.Value.(*protocol.StringValue)
	if !ok {
		return "", true, protocol.ErrWrongType
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
func (c *Cache) MGet(keys []string) ([]protocol.Response, error) {
	responses := make([]protocol.Response, len(keys))

	for i, key := range keys {
		value, exists, err := c.Get(key)
		if err != nil {
			return nil, err
		}
		if exists {
			responses[i] = protocol.NewBulkString(value)
		} else {
			responses[i] = protocol.NewNullBulkString()
		}
	}

	return responses, nil
}
func (c *Cache) MSet(keys []string, values []string) {

	for i := 0; i < len(keys); i++ {
		c.data[keys[i]] = Entry{
			Value: &protocol.StringValue{
				Value: values[i],
			},
		}
	}
}
func (c *Cache) Append(key string, value string) int {

	entry, exists := c.getEntry(key)

	if !exists {
		c.data[key] = Entry{
			Value: &protocol.StringValue{
				Value: value,
			},
		}
		return len(value)
	}
	str, err := c.GetString(entry)
	if err != nil {
		return 0
	}

	str.Value += value

	c.data[key] = *entry

	return len(str.Value)
}
func (c *Cache) Strlen(key string) int {

	entry, exists := c.getEntry(key)
	if !exists {
		return 0
	}
	str, err := c.GetString(entry)
	if err != nil {
		return 0
	}

	return len(str.Value)
}
func (c *Cache) GetSet(key string, value string) (string, bool) {

	entry, exists := c.getEntry(key)

	if !exists {
		c.data[key] = Entry{
			Value: &protocol.StringValue{
				Value: value,
			},
		}
		return "", false
	}
	str, err := c.GetString(entry)
	if err != nil {
		return "", false
	}

	oldValue := str.Value
	c.data[key] = Entry{
		Value: &protocol.StringValue{
			Value: value,
		},
	}
	return oldValue, true
}

func (c *Cache) changeBy(key string, delta int) (int, error) {

	entry, exists := c.getEntry(key)

	// Key doesn't exist
	if !exists {
		c.data[key] = Entry{
			Value: &protocol.StringValue{
				Value: strconv.Itoa(delta),
			},
		}
		return delta, nil
	}
	str, err := c.GetString(entry)
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
func (c *Cache) GetString(entry *Entry) (*protocol.StringValue, error) {
	str, ok := entry.Value.(*protocol.StringValue)
	if !ok {
		return nil, protocol.ErrWrongType
	}
	return str, nil
}
