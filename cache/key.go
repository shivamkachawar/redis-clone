package cache

func (c *Cache) Delete(keys []string) int {

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

	count := 0
	for _, key := range keys {
		_, exists := c.getEntry(key)
		if exists {
			count++
		}
	}
	return count
}
func (c *Cache) Keys() []string {

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

	c.data = make(map[string]Entry)
}
