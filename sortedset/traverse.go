package sortedset

func (sl *SkipList) findUpdate(score float64) []*SkipListNode {
	update := make([]*SkipListNode, MaxLevel)
	current := sl.Head

	for level := sl.Level; level >= 0; level-- {
		for current.Forward[level] != nil && current.Forward[level].Score < score {
			current = current.Forward[level]
		}
		update[level] = current
	}
	return update
}
