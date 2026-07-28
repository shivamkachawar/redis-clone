package main

import (
	"errors"
	"fmt"
	"strconv"
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

// Client APIs

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

func (c *Cache) Delete(keys []string) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	deleted := 0

	for _, key := range keys {
		_, exists := c.getEntry(key)
		if exists {
			delete(c.data, key)
			deleted++
		}
	}

	return deleted
}

func (c *Cache) Exists(keys []string) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	count := 0
	for _, key := range keys {
		_, exists := c.getEntry(key)
		if exists {
			count++
		}
	}
	return count
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
func (c *Cache) Keys() []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	keys := make([]string, 0, len(c.data))

	for key := range c.data {
		_, exists := c.getEntry(key)
		if exists {
			keys = append(keys, key)
		}
	}
	return keys
}
func (c *Cache) Size() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	count := 0

	for key := range c.data {
		_, exists := c.getEntry(key)
		if exists {
			count++
		}
	}

	return count
}
func (c *Cache) Flush() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]Entry)
}
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
func (c *Cache) LPush(key string, values []string) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, exists := c.getEntry(key)
	var list *ListValue

	if !exists {
		list = &ListValue{
			Values: []string{},
		}

		c.data[key] = Entry{
			Value: list,
		}
	} else {
		var ok bool
		list, ok = entry.Value.(*ListValue)
		if !ok {
			return 0, errors.New("WRONGTYPE")
		}
	}

	for _, value := range values {
		list.Values = append([]string{value}, list.Values...)
	}

	return len(list.Values), nil
}
func (c *Cache) LRange(key string, start int, stop int) ([]string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	if !exists {
		return []string{}, nil
	}
	list, ok := entry.Value.(*ListValue)
	if !ok {
		return nil, ErrWrongType
	}
	if len(list.Values) == 0 {
		return []string{}, nil
	}
	if start < 0 {
		start = len(list.Values) + start
	}
	if stop < 0 {
		stop = len(list.Values) + stop
	}
	if start < 0 {
		start = 0
	}
	if stop < 0 {
		stop = 0
	}
	if start >= len(list.Values) {
		return []string{}, nil
	}
	if stop >= len(list.Values) {
		stop = len(list.Values) - 1
	}
	if start > stop {
		return []string{}, nil
	}

	return list.Values[start : stop+1], nil
}
func (c *Cache) LLen(key string) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	if !exists {
		return 0, nil
	}

	list, ok := entry.Value.(*ListValue)
	if !ok {
		return 0, ErrWrongType
	}

	return len(list.Values), nil
}
func (c *Cache) LPop(key string) (string, bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	if !exists {
		return "", false, nil
	}

	list, ok := entry.Value.(*ListValue)
	if !ok {
		return "", true, ErrWrongType
	}

	if len(list.Values) == 0 {
		return "", false, nil
	}
	value := list.Values[0]
	list.Values = list.Values[1:]
	if len(list.Values) == 0 {
		delete(c.data, key)
	}

	return value, true, nil
}
