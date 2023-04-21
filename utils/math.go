package utils

func ClampFloat64(val float64, min float64, max float64) float64 {
	if val <= min {
		return min
	} else if val >= max {
		return max
	}
	return val
}
