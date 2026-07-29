package cache

import (
	"go-redis/protocol"
	"strconv"
)

func (c *Cache) getOrCreateHash(key string) (*protocol.HashValue, error) {
	entry, exists := c.getEntry(key)

	if !exists {
		hash := &protocol.HashValue{
			Fields: make(map[string]string),
		}
		c.data[key] = Entry{
			Value: hash,
		}
		return hash, nil
	}
	hash, ok := entry.Value.(*protocol.HashValue)
	if !ok {
		return nil, protocol.ErrWrongType
	}

	return hash, nil

}

func (c *Cache) HSet(key, field, value string) (int, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	hash, err := c.getOrCreateHash(key)

	if err != nil {
		return 0, err
	}
	_, exists := hash.Fields[field]
	hash.Fields[field] = value

	if exists {
		return 0, nil
	}

	return 1, nil
}

func (c *Cache) HGet(key, field string) (string, bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	if !exists {
		return "", false, nil
	}
	hash, ok := entry.Value.(*protocol.HashValue)

	if !ok {
		return "", false, protocol.ErrWrongType
	}
	value, exists := hash.Fields[field]

	if !exists {
		return "", false, nil
	}

	return value, true, nil
}

func (c *Cache) HDel(key, field string) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	if !exists {
		return 0, nil
	}
	hash, ok := entry.Value.(*protocol.HashValue)
	if !ok {
		return 0, protocol.ErrWrongType
	}
	_, exists = hash.Fields[field]

	if !exists {
		return 0, nil
	}
	delete(hash.Fields, field)

	return 1, nil

}
func (c *Cache) HExists(key, field string) (bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	if !exists {
		return false, nil
	}
	hash, ok := entry.Value.(*protocol.HashValue)
	if !ok {
		return false, protocol.ErrWrongType
	}
	_, exists = hash.Fields[field]

	return exists, nil
}
func (c *Cache) HLen(key string) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	if !exists {
		return 0, nil
	}
	hash, ok := entry.Value.(*protocol.HashValue)
	if !ok {
		return 0, protocol.ErrWrongType
	}
	return len(hash.Fields), nil
}
func (c *Cache) HGetAll(key string) ([]string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)
	if !exists {
		return []string{}, nil
	}
	hash, ok := entry.Value.(*protocol.HashValue)
	if !ok {
		return nil, protocol.ErrWrongType
	}
	result := []string{}

	for field, value := range hash.Fields {
		result = append(result, field)
		result = append(result, value)
	}
	return result, nil
}
func (c *Cache) HKeys(key string) ([]string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)
	if !exists {
		return []string{}, nil
	}

	hash, ok := entry.Value.(*protocol.HashValue)
	if !ok {
		return nil, protocol.ErrWrongType
	}

	keys := []string{}

	for field := range hash.Fields {
		keys = append(keys, field)
	}

	return keys, nil
}
func (c *Cache) HVals(key string) ([]string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)
	if !exists {
		return []string{}, nil
	}

	hash, ok := entry.Value.(*protocol.HashValue)
	if !ok {
		return nil, protocol.ErrWrongType
	}

	values := []string{}

	for _, value := range hash.Fields {
		values = append(values, value)
	}

	return values, nil
}
func (c *Cache) HIncrBy(key, field string, increment int) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	hash, err := c.getOrCreateHash(key)
	if err != nil {
		return 0, err
	}
	value, exists := hash.Fields[field]

	if !exists {
		hash.Fields[field] = strconv.Itoa(increment)
		return increment, nil
	}
	currentValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, protocol.ErrHashValueNotInteger
	}
	newValue := currentValue + increment
	hash.Fields[field] = strconv.Itoa(newValue)
	return newValue, nil
}
func (c *Cache) HSetNX(key, field, value string) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	hash, err := c.getOrCreateHash(key)
	if err != nil {
		return 0, err
	}
	_, exists := hash.Fields[field]
	if exists {
		return 0, nil
	}
	hash.Fields[field] = value
	return 1, nil
}
