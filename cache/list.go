package cache

import (
	"go-redis/protocol"
)

func (c *Cache) getOrCreateList(key string) (*protocol.ListValue, error) {
	entry, exists := c.getEntry(key)

	if !exists {
		list := &protocol.ListValue{
			Values: []string{},
		}

		c.data[key] = Entry{
			Value: list,
		}

		return list, nil
	}

	list, ok := entry.Value.(*protocol.ListValue)
	if !ok {
		return nil, protocol.ErrWrongType
	}

	return list, nil
}
func (c *Cache) LPush(key string, values []string) (int, error) {

	list, err := c.getOrCreateList(key)
	if err != nil {
		return 0, err
	}

	for _, value := range values {
		list.Values = append([]string{value}, list.Values...)
	}

	return len(list.Values), nil
}
func (c *Cache) LRange(key string, start int, stop int) ([]string, error) {

	entry, exists := c.getEntry(key)

	if !exists {
		return []string{}, nil
	}
	list, ok := entry.Value.(*protocol.ListValue)
	if !ok {
		return nil, protocol.ErrWrongType
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

	entry, exists := c.getEntry(key)

	if !exists {
		return 0, nil
	}

	list, ok := entry.Value.(*protocol.ListValue)
	if !ok {
		return 0, protocol.ErrWrongType
	}

	return len(list.Values), nil
}
func (c *Cache) getList(key string) (*protocol.ListValue, bool, error) {
	entry, exists := c.getEntry(key)

	if !exists {
		return nil, false, nil
	}

	list, ok := entry.Value.(*protocol.ListValue)
	if !ok {
		return nil, true, protocol.ErrWrongType
	}

	if len(list.Values) == 0 {
		delete(c.data, key)
		return nil, false, nil
	}

	return list, true, nil
}
func (c *Cache) LPop(key string) (string, bool, error) {

	list, exists, err := c.getList(key)
	if err != nil || !exists {
		return "", exists, err
	}

	value := list.Values[0]
	list.Values = list.Values[1:]

	if len(list.Values) == 0 {
		delete(c.data, key)
	}

	return value, true, nil
}
func (c *Cache) RPush(key string, values []string) (int, error) {

	list, err := c.getOrCreateList(key)
	if err != nil {
		return 0, err
	}

	list.Values = append(list.Values, values...)

	return len(list.Values), nil
}
func (c *Cache) RPop(key string) (string, bool, error) {

	list, exists, err := c.getList(key)
	if err != nil || !exists {
		return "", exists, err
	}

	last := len(list.Values) - 1
	value := list.Values[last]

	list.Values = list.Values[:last]

	if len(list.Values) == 0 {
		delete(c.data, key)
	}

	return value, true, nil
}
