package smallMath

import "math"

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Min(nums ...int) int {
	min := math.MaxInt64
	index := 0
	for i, num := range nums {
		if num < min {
			min = num
			index = i
		}
	}

	return index
}
