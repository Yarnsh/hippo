package utils

import (
	"math"
)

func ClampFloat64(val float64, min float64, max float64) float64 {
	if val <= min {
		return min
	} else if val >= max {
		return max
	}
	return val
}

func TowardsByFloat64(val, target, by float64) float64 {
	if val < target {
		return math.Min(target, val + by)
	} else {
		return math.Max(target, val - by)
	}
}
