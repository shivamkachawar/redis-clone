package cache

import "go-redis/protocol"

func (c *Cache) getOrCreateSet(key string) (*protocol.SetValue, error) {
	entry, exists := c.getEntry(key)

	if !exists {
		set := &protocol.SetValue{
			Members: make(map[string]struct{}),
		}
		c.data[key] = Entry{
			Value: set,
		}
		return set, nil
	}
	set, ok := entry.Value.(*protocol.SetValue)
	if !ok {
		return nil, protocol.ErrWrongType
	}

	return set, nil

}

func (c *Cache) SAdd(key, member string) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	set, err := c.getOrCreateSet(key)

	if err != nil {
		return 0, nil
	}
	if _, exists := set.Members[member]; exists {
		return 0, nil
	}
	set.Members[member] = struct{}{}

	return 1, nil
}
func (c *Cache) SIsMember(key, member string) (bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	if !exists {
		return false, nil
	}

	set, ok := entry.Value.(*protocol.SetValue)
	if !ok {
		return false, protocol.ErrWrongType
	}

	_, exists = set.Members[member]
	return exists, nil
}
func (c *Cache) SCard(key string) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)

	if !exists {
		return 0, nil
	}

	set, ok := entry.Value.(*protocol.SetValue)
	if !ok {
		return 0, protocol.ErrWrongType
	}

	return len(set.Members), nil
}
func (c *Cache) SMembers(key string) ([]string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)
	if !exists {
		return []string{}, nil
	}

	set, ok := entry.Value.(*protocol.SetValue)
	if !ok {
		return nil, protocol.ErrWrongType
	}

	members := make([]string, 0, len(set.Members))

	for member := range set.Members {
		members = append(members, member)
	}

	return members, nil
}
func (c *Cache) SRem(key, member string) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.getEntry(key)
	if !exists {
		return 0, nil
	}

	set, ok := entry.Value.(*protocol.SetValue)
	if !ok {
		return 0, protocol.ErrWrongType
	}

	if _, exists := set.Members[member]; !exists {
		return 0, nil
	}

	delete(set.Members, member)

	return 1, nil
}
