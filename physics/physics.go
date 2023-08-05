package physics

import (
	"github.com/Yarnsh/hippo/shapes"
)

type Body struct {
	x, y, r float64
	shapes []*shapes.Shape
}
