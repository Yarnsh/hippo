package shapes

func CrossProduct2D(x1 int, y1 int, x2 int, y2 int) float64 {
	return float64((x1*y2) - (y1*x2))
}

type Line struct {
	x int
	y int
	dirx int
	diry int

	bounding_box AxisRect
}

// Getters
func (ray Line) X() int {
	return ray.x
}
func (ray Line) Y() int {
	return ray.y
}
func (ray Line) DirX() int {
	return ray.dirx
}
func (ray Line) DirY() int {
	return ray.diry
}
func (ray Line) BoundingBox() AxisRect {
	return ray.bounding_box
}

func (ray *Line) Init(x int, y int, dirx int, diry int) {
	ray.x = x
	ray.y = y
	ray.dirx = dirx
	ray.diry = diry

	ray.bounding_box = NewAxisRect(x, y, dirx, diry)
}

func NewLine(x int, y int, x2 int, y2 int) Line {
	ray := Line{}
	ray.Init(x, y, x2 - x, y2 - y)
	return ray
}

func (ray Line) GetLineIntersectionTimeRaw(x int, y int, dirx int, diry int) float64 {
	// 2.0 is the value we use for a miss, I forget why
	denom := CrossProduct2D(ray.dirx, ray.diry, dirx, diry)
	if denom == 0.0 {
		return 2.0
	}
	to_other_x := x - ray.x
	to_other_y := y - ray.y
	t := CrossProduct2D(to_other_x, to_other_y, dirx, diry) / denom
	u := CrossProduct2D(to_other_x, to_other_y, ray.dirx, ray.diry) / denom
	if u > 1.0 || u < 0.0 || t > 1.0 || t < 0.0 {
		return 2.0
	}
	return t
}

func (ray Line) GetAxisRectIntersectionTime(other AxisRect) float64 {
	// We treat rects as filled in, so if our starting point is inside it we treat that as the hit
	if other.ContainsPoint(float64(ray.x), float64(ray.y)) {
		return 0.0
	}

	t := 2.0
	
	if ray.dirx != 0.0 {
		if ray.x < int(other.x) {
			u := ray.GetLineIntersectionTimeRaw(int(other.x), int(other.y), 0.0, int(other.h))
			if u >= 0.0 && u < t {
				t = u
			}
		} else if ray.x > int(other.x2) {
			u := ray.GetLineIntersectionTimeRaw(int(other.x2), int(other.y), 0.0, int(other.h))
			if u >= 0.0 && u < t {
				t = u
			}
		}
	}
	
	if ray.diry != 0.0 {
		if ray.y < int(other.y) {
			u := ray.GetLineIntersectionTimeRaw(int(other.x), int(other.y), int(other.w), 0.0)
			if u >= 0.0 && u < t {
				t = u
			}
		} else if ray.y > int(other.y2) {
			u := ray.GetLineIntersectionTimeRaw(int(other.x), int(other.y2), int(other.w), 0.0)
			if u >= 0.0 && u < t {
				t = u
			}
		}
	}

	if t > 1.0 {
		return 2.0
	}
	return t
}

func (ray Line) IntersectsAxisRect(other AxisRect) bool {
	return ray.GetAxisRectIntersectionTime(other) < 2.0
}