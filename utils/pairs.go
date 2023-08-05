package utils

import (
	"math"
)

type IntPair struct {
	X, Y int
}
type FloatPair struct {
	X, Y float64
}

func (pair FloatPair) Normalized() FloatPair {
	len := math.Sqrt((pair.X * pair.X) + (pair.Y * pair.Y))
	if len == 0.0 {
		return FloatPair{}
	}
	return FloatPair {
		X: pair.X / len,
		Y: pair.Y / len,
	}
}
func (pair FloatPair) Negative() FloatPair {
	return FloatPair {
		X: -pair.X,
		Y: -pair.Y,
	}
}
func (pair FloatPair) Minus(other FloatPair) FloatPair {
	return FloatPair {
		X: pair.X - other.X,
		Y: pair.Y - other.Y,
	}
}
func (pair FloatPair) Plus(other FloatPair) FloatPair {
	return FloatPair {
		X: pair.X + other.X,
		Y: pair.Y + other.Y,
	}
}
func (pair FloatPair) Multiply(other float64) FloatPair {
	return FloatPair {
		X: pair.X * other,
		Y: pair.Y *other,
	}
}
func (pair IntPair) Minus(other IntPair) IntPair {
	return IntPair {
		X: pair.X - other.X,
		Y: pair.Y - other.Y,
	}
}
func (pair IntPair) Plus(other IntPair) IntPair {
	return IntPair {
		X: pair.X + other.X,
		Y: pair.Y + other.Y,
	}
}

func (pair FloatPair) DistanceTo(other FloatPair) float64 {
	p := pair.Minus(other)
	return math.Sqrt((p.X * p.X) + (p.Y * p.Y))
}

func (pair FloatPair) Length() float64 {
	return math.Sqrt((pair.X * pair.X) + (pair.Y * pair.Y))
}

func (pair FloatPair) ToInt() IntPair {
	return IntPair {
		X: int(pair.X),
		Y: int(pair.Y),
	}
}

func (pair IntPair) ToFloat() FloatPair {
	return FloatPair {
		X: float64(pair.X),
		Y: float64(pair.Y),
	}
}
