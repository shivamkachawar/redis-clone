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
