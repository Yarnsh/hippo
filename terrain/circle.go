package terrain

import (
	"github.com/Yarnsh/hippo/utils"
	"math"
)

type Circle struct {
	pos utils.FloatPair
	radius float64

	radius_squared float64
}

func (circ *Circle) Init(x float64, y float64, r float64) {
	circ.pos = utils.FloatPair {
		X: x,
		Y: y,
	}
	circ.radius = r

	circ.radius_squared = r * r
}

func NewCircle(x, y, r float64) Circle {
	circ := Circle{}
	circ.Init(x, y, r)
	return circ
}

func (circ *Circle) SetRadius(r float64) {
	circ.radius = r
	circ.radius_squared = r * r
}

func (circ *Circle) SetPosition(x float64, y float64) {
	circ.pos = utils.FloatPair {
		X: x,
		Y: y,
	}
}

func (circ *Circle) Translate(x float64, y float64) {
	circ.pos = utils.FloatPair {
		X: circ.pos.X + x,
		Y: circ.pos.Y + y,
	}
}

func (circ Circle) SeparationForAxisRect(other AxisRect) (utils.FloatPair, float64) {
	closest_x := utils.ClampFloat64(circ.pos.X, float64(other.x), float64(other.x2))
	closest_y := utils.ClampFloat64(circ.pos.Y, float64(other.y), float64(other.y2))

	squared_dist_to_closest := ((circ.pos.X - closest_x) * (circ.pos.X - closest_x)) + ((circ.pos.Y - closest_y) * (circ.pos.Y - closest_y))
	if (squared_dist_to_closest > circ.radius_squared) {
		return utils.FloatPair{}, 0
	}

	dist := math.Sqrt(squared_dist_to_closest)
	closest := utils.FloatPair{X: closest_x, Y: closest_y}
	return closest.Minus(circ.pos).Normalized().Multiply(-(circ.radius - dist)), math.Abs(circ.radius - dist)
}

func (circ Circle) BBIntersectsAxisRect(other AxisRect) bool {
	// Quick check to see if the bounding box of the circle intersects another rect, for faster pre-checks in quad trees
	return !((circ.pos.X + circ.radius < float64(other.x)) || (circ.pos.X - circ.radius > float64(other.x2)) || (circ.pos.Y + circ.radius < float64(other.y)) || (circ.pos.Y - circ.radius > float64(other.y2)))
}