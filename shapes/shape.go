package shapes

import (
	"github.com/Yarnsh/hippo/utils"
)

type Shape interface {
	Translated() Shape
	Rotated() Shape
	// collision depth, and two manifold points
	TestCollision(Shape) (float64, utils.FloatPair, utils.FloatPair)
}
