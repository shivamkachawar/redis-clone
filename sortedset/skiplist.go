package sortedset

const (
	MaxLevel    = 32
	Probability = 0.25
)

type SkipListNode struct {
	Member  string
	Score   float64
	Forward []*SkipListNode
}
type SkipList struct {
	Head   *SkipListNode
	Tail   *SkipListNode
	Length uint64
	Level  int
}
