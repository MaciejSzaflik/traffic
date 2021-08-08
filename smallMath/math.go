package smallMath

import "math"

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Min(nums ...int) (int, int) {
	min := math.MaxInt64
	index := 0
	for i, num := range nums {
		if num < min {
			min = num
			index = i
		}
	}

	return min, index
}

func Map(value, min1, max1, min2, max2 float32) float32 {
	return min2 + (value-min1)*(max2-min2)/(max1-min1)
}

func MinFloat(nums ...float32) (float32, int) {
	min := float32(math.MaxFloat32)
	index := 0
	for i, num := range nums {
		if num < min {
			min = num
			index = i
		}
	}

	return min, index
}

func FloorToInt(a float64) int {
	return int(math.Floor(a))
}

func FloorToInt32(a float32) int {
	return int(math.Floor(float64(a)))
}

func CantorPair(x, y int) int {
	return ((x+y)*(x+y+1))/2 + y
}

func CantorDepair(pair int) (int, int) {
	t := FloorToInt((math.Sqrt(float64(8*pair+1)) - 1) / 2)
	return t*(t+3)/2 - pair, pair - t*(t+1)/2
}
