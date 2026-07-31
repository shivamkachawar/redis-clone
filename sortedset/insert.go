package sortedset

func (sl *SkipList) Insert(score float64, member string) *SkipListNode {
	update := sl.findUpdate(score)
	level := randomLevel()
	node := &SkipListNode{
		Member:  member,
		Score:   score,
		Forward: make([]*SkipListNode, level+1),
	}

	if level > sl.Level {
		for i := sl.Level + 1; i <= level; i++ {
			update[i] = sl.Head
		}
		sl.Level = level
	}

	for i := 0; i <= level; i++ {
		node.Forward[i] = update[i].Forward[i]
		update[i].Forward[i] = node
	}
	sl.Length++
	if node.Forward[0] == nil {
		sl.Tail = node
	}
	return node
}
