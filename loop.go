package main

//func optimizeZeroIndex(behaviour map[int]int) []Routine {
//	keys := make([]int, 0)
//	for k := range behaviour {
//		keys = append(keys, k)
//	}
//	sort.Ints(keys)
//	newOps := make([]Routine, 0)
//
//	newOps = append(newOps, Zero)
//	for _, k := range keys {
//		if k == 0 { // Already handled.
//			continue
//		}
//		v := behaviour[k]
//		if abs(v) > math.MaxInt32 {
//			panic("THE OFFSET IS TOO LARGE! I CANNOT HANDLE THIS")
//		}
//
//		newOps = append(newOps, OpWithArgOffset{
//			op:     DataIncArgOffset,
//			arg:    v,
//			offset: k,
//		})
//	}
//
//	return newOps
//}

func abs(i int) int {
	if i >= 0 {
		return i
	} else {
		return -i
	}
}
