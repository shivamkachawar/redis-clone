package sortedset

type SortedSet struct {
	skipList *SkipList
	dict     map[string]*SkipListNode
}

func NewSortedSet() *SortedSet {
	return &SortedSet{
		skipList: NewSkipList(),
		dict:     make(map[string]*SkipListNode),
	}
}

func (ss *SortedSet) Add(score float64, member string) {
	existing, exists := ss.dict[member]

	if exists {
		if existing.Score == score {
			return
		}
		ss.skipList.Delete(existing.Score, member)
	}
	node := ss.skipList.Insert(score, member)

	ss.dict[member] = node
}

func (ss *SortedSet) Delete(member string) bool {

	node, exists := ss.dict[member]

	if !exists {
		return false
	}

	if !ss.skipList.Delete(node.Score, member) {
		return false
	}

	delete(ss.dict, member)

	return true
}

func (ss *SortedSet) Score(member string) (float64, bool) {

	node, exists := ss.dict[member]

	if !exists {
		return 0, false
	}

	return node.Score, true
}

func (ss *SortedSet) Card() uint64 {
	return ss.skipList.Length
}
func (ss *SortedSet) Range(start, stop int) []string {
	result := []string{}

	length := int(ss.skipList.Length)

	if length == 0 {
		return result
	}

	// Convert negative indices
	if start < 0 {
		start += length
	}
	if stop < 0 {
		stop += length
	}

	// Clamp to valid range
	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}

	// Invalid range
	if start > stop || start >= length {
		return result
	}

	current := ss.skipList.Head.Forward[0]
	index := 0

	for current != nil && index <= stop {
		if index >= start {
			result = append(result, current.Member)
		}

		current = current.Forward[0]
		index++
	}

	return result
}
func (ss *SortedSet) Rank(member string) (int, bool) {
	node, exists := ss.dict[member]
	if !exists {
		return 0, false
	}

	current := ss.skipList.Head.Forward[0]
	rank := 0

	for current != nil {
		if current == node {
			return rank, true
		}

		current = current.Forward[0]
		rank++
	}

	return 0, false
}
