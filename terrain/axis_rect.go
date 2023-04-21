package terrain

type AxisRect struct {
	x int
	y int
	w int
	h int

	x2 int
	y2 int
}

// Getters
func (rect AxisRect) X() int {
	return rect.x
}
func (rect AxisRect) Y() int {
	return rect.y
}
func (rect AxisRect) W() int {
	return rect.w
}
func (rect AxisRect) H() int {
	return rect.h
}
func (rect AxisRect) X2() int {
	return rect.x2
}
func (rect AxisRect) Y2() int {
	return rect.y2
}
// End getters

func NewAxisRect(x int, y int, w int, h int) AxisRect {
	rect := AxisRect{}
	rect.x = x
	rect.y = y
	rect.w = w
	rect.h = h

	// Flip the rectangle around in the case of negative size
	if (w < 0) {
		rect.x += w
		rect.w = -w
	}
	if (h < 0) {
		rect.y += h
		rect.h = -h
	}

	rect.x2 = rect.x + rect.w
	rect.y2 = rect.y + rect.h

	return rect
}

func (rect *AxisRect) SetPosition(x int, y int) {
	rect.x = x
	rect.y = y

	rect.x2 = rect.x + rect.w
	rect.y2 = rect.y + rect.h
}

func (rect AxisRect) Translated(x int, y int) AxisRect {
	return NewAxisRect(rect.x + x, rect.y + y, rect.w, rect.h)
}

func (rect *AxisRect) SetSize(w int, h int) {
	rect.w = w
	rect.h = h

	// Flip the rectangle around in the case of negative size
	if (w < 0) {
		rect.x += w
		rect.w = -w
	}
	if (h < 0) {
		rect.y += h
		rect.h = -h
	}

	rect.x2 = rect.x + rect.w
	rect.y2 = rect.y + rect.h
}

/*func (rect AxisRect) Intersects(other Shape) bool { // TODO: handle pointer types too
	switch other.(type) {
		case AxisRect:
			return rect.IntersectsAxisRect(other.(AxisRect))
		case Circle:
			return rect.IntersectsCircle(other.(Circle))
		case Ray:
			return other.(Ray).Intersects(rect)
		default:
		    // Don't know how to handle this
		    return false
	}
}

func (rect AxisRect) Encloses(other Shape) bool { // TODO: handle pointer types too
	switch other.(type) {
		case AxisRect:
			return rect.EnclosesAxisRect(other.(AxisRect))
		case Circle:
			return rect.EnclosesCircle(other.(Circle))
		case Ray:
			return other.(Ray).Encloses(rect)
		default:
		    // Don't know how to handle this
		    return false
	}
}*/

func (rect AxisRect) IntersectsAxisRect(other AxisRect) bool {
	bigx := other.x2 + rect.w
	if (rect.x2 >= other.x && rect.x2 <= bigx) {
		bigy := other.y2 + rect.h
		return (rect.y2 >= other.y && rect.y2 <= bigy)
	}
	return false
}

func (rect AxisRect) EnclosesAxisRect(other AxisRect) bool {
	return (other.x >= rect.x && other.x2 <= rect.x2 && other.y >= rect.y && other.y2 <= rect.y2)
}

/*func (rect AxisRect) IntersectsCircle(other Circle) bool {
	return other.Intersects(rect)
}

func (rect AxisRect) EnclosesCircle(other Circle) bool {
	return (other.x - other.r >= float64(rect.x) && other.x + other.r <= float64(rect.x2) && other.y - other.r >= float64(rect.y) && other.y + other.r <= float64(rect.y2))
}*/

func (rect AxisRect) ContainsPoint(x float64, y float64) bool {
	return x >= float64(rect.x) && x <= float64(rect.x2) && y >= float64(rect.y) && y <= float64(rect.y2)
}

/*func (rect AxisRect) SweepAxisRect(dirx float64, diry float64, other AxisRect) (float64, float64, float64) {
	// Check if its already inside
	if rect.IntersectsAxisRect(other) {
		return 0.0, 0.0, 0.0
	}

	// returns time of collision, and collision normal (x, y)
	var dxEntry, dxExit, dyEntry, dyExit float64
	var txEntry, txExit, tyEntry, tyExit float64

	if dirx > 0.0 {
		dxEntry = other.x - rect.x2
		dxExit = other.x2 - rect.x
	} else {
		dxEntry = other.x2 - rect.x
		dxExit = other.x - rect.x2
	}
	if diry > 0.0 {
		dyEntry = other.y - rect.y2
		dyExit = other.y2 - rect.y
	} else {
		dyEntry = other.y2 - rect.y
		dyExit = other.y - rect.y2
	}

	if dirx == 0.0 {
		if !(rect.x2 < other.x || rect.x > other.x) {
			txEntry = 0.0
		} else {
			txEntry = 2.0
		}
		txExit = 2.0
	} else {
		txEntry = dxEntry / dirx
		txExit = dxExit / dirx
	}
	if diry == 0.0 {
		if !(rect.y2 < other.y || rect.y > other.y) {
			tyEntry = 0.0
		} else {
			tyEntry = 2.0
		}
		tyExit = 2.0
	} else {
		tyEntry = dyEntry / diry
		tyExit = dyExit / diry
	}

	tEntry := math.Max(txEntry, tyEntry)
	tExit := math.Min(txExit, tyExit)

	if tEntry > tExit || tEntry < 0.0 || txEntry > 1.0 || tyEntry > 1.0{
		return 2.0, 0.0, 0.0
	}

	normalX := 0.0
	normalY := 0.0
	if txEntry < tyEntry {
		if dyEntry > 0.0 {
			normalY = -1.0
		} else {
			normalY = 1.0
		}
	} else {
		if dxEntry > 0.0 {
			normalX = -1.0
		} else {
			normalX = 1.0
		}
	}

	return tEntry, normalX, normalY
}*/

