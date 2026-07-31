package sortedset

func (sl *SkipList) Delete(score float64, member string) bool {
	update := sl.findUpdate(score)
	target := update[0].Forward[0]

	if target == nil || target.Score != score || target.Member != member {
		return false
	}

	for i := 0; i < len(target.Forward); i++ {
		update[i].Forward[i] = target.Forward[i]
	}

	if sl.Tail == target {
		if update[0] == sl.Head {
			sl.Tail = nil
		} else {
			sl.Tail = update[0]
		}
	}
	sl.Length--

	for sl.Level > 0 && sl.Head.Forward[sl.Level] == nil {
		sl.Level--
	}

	return true
}
