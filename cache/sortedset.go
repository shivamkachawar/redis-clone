package cache

import (
	"go-redis/protocol"
	"go-redis/sortedset"
)

func (c *Cache) getOrCreateSortedSet(key string) (*sortedset.SortedSet, error) {
	entry, exists := c.getEntry(key)

	if !exists {
		ss := sortedset.NewSortedSet()

		c.data[key] = Entry{
			Value: &protocol.SortedSetValue{
				Value: ss,
			},
		}

		return ss, nil
	}

	value, ok := entry.Value.(*protocol.SortedSetValue)
	if !ok {
		return nil, protocol.ErrWrongType
	}

	return value.Value, nil
}

func (c *Cache) ZAdd(key string, score float64, member string) error {

	ss, err := c.getOrCreateSortedSet(key)
	if err != nil {
		return err
	}

	ss.Add(score, member)

	return nil
}

func (c *Cache) ZRem(key, member string) (bool, error) {

	ss, err := c.getOrCreateSortedSet(key)
	if err != nil {
		return false, err
	}

	return ss.Delete(member), nil
}

func (c *Cache) ZScore(key, member string) (float64, bool, error) {

	entry, exists := c.getEntry(key)
	if !exists {
		return 0, false, nil
	}

	value, ok := entry.Value.(*protocol.SortedSetValue)
	if !ok {
		return 0, false, protocol.ErrWrongType
	}

	score, found := value.Value.Score(member)

	return score, found, nil
}
func (c *Cache) ZCard(key string) (uint64, error) {

	entry, exists := c.getEntry(key)
	if !exists {
		return 0, nil
	}

	value, ok := entry.Value.(*protocol.SortedSetValue)
	if !ok {
		return 0, protocol.ErrWrongType
	}

	return value.Value.Card(), nil
}
func (c *Cache) ZRange(key string, start, stop int) ([]string, error) {
	entry, exists := c.getEntry(key)
	if !exists {
		return []string{}, nil
	}

	value, ok := entry.Value.(*protocol.SortedSetValue)
	if !ok {
		return nil, protocol.ErrWrongType
	}

	return value.Value.Range(start, stop), nil
}
