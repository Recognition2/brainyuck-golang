package main

import (
	"sort"
)

func optimizeZeroIndex(behaviour map[int]int) []Op {
	keys := make([]int, 0)
	for k := range behaviour {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	newOps := make([]Op, 0)

	newOps = append(newOps, Zero)
	for _, k := range keys {
		if k == 0 { // Already handled.
			continue
		}
		v := behaviour[k]
		if abs(v) > 255 {
			panic("THE OFFSET IS TOO LARGE! I CANNOT HANDLE THIS")
		}
		newOps = append(newOps, Offset)
		newOps = append(newOps, Op(int8(k)))
		if v == -1 {
			newOps = append(newOps, DataDecOffset)
		} else if v == 1 {
			newOps = append(newOps, DataIncOffset)
		} else if v < 0 {
			newOps = append(newOps, DataDecArgOffset)
			newOps = append(newOps, Op(-v))
		} else if v > 0 {
			newOps = append(newOps, DataIncArgOffset)
			newOps = append(newOps, Op(v))
		}
	}

	return newOps
}
func abs(i int) int {
	if i >= 0 {
		return i
	} else {
		return -i
	}
}
