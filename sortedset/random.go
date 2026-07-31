package sortedset

import (
	"math/rand/v2"
)

func randomLevel() int {
	level := 0

	for rand.Float64() < Probability {
		if level == MaxLevel-1 {
			break
		}
		level++
	}
	return level

}
