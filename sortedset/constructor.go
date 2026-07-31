package sortedset

import "math"

func NewSkipList() *SkipList {
	head := &SkipListNode{
		Member:  "",
		Score:   math.Inf(-1),
		Forward: make([]*SkipListNode, MaxLevel),
	}

	return &SkipList{
		Head:   head,
		Tail:   nil,
		Length: 0,
		Level:  0,
	}
}
