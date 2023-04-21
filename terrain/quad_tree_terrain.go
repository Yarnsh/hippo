package terrain

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"

	"github.com/fzipp/astar"
	"github.com/Yarnsh/hippo/utils"
)

const (
	Sqrt2 = 1.414213
)

type void struct{}
var void_item void

type QuadTreeTerrain struct {
	space AxisRect
	leaf bool
	leaf_value int
	sub_trees [4]*QuadTreeTerrain
	pixel_x, pixel_y, pixel_width int

	dirty bool
}

// Quad tree terrain should be square, hence only width
// They should also be powers of 2 in size, should maybe fix that
func NewQuadTreeTerrain(x int, y int, w int) *QuadTreeTerrain {
	tree := QuadTreeTerrain{}
	tree.pixel_x = x
	tree.pixel_y = y
	tree.pixel_width = w
	tree.space = NewAxisRect(x, y, w, w)
	tree.leaf = true
	tree.dirty = false

	return &tree
}

func (tree QuadTreeTerrain) materialFromColor(color color.Color) int {
	r,g,b,_ := color.RGBA()
	if r > 0 || g > 0 || b > 0 {
		return 1
	}
	return 0
}

func (tree *QuadTreeTerrain) LoadImageData(image *ebiten.Image) {
	if !tree.leaf { // We are overwriting everything anyway so just collapse it down
		tree.Join()
	}

	startx := tree.pixel_x
	starty := tree.pixel_y

	color := image.At(startx, starty)
	value := tree.materialFromColor(color)
	for y := starty; y < starty + tree.pixel_width; y++ {
		for x := startx; x < startx + tree.pixel_width; x++ {
			color = image.At(x, y)
			val_new := tree.materialFromColor(color)

			if val_new != value { // We can't paint this tree all one color so we must split
				tree.Split()
				for _, sub_tree := range tree.sub_trees {
					sub_tree.LoadImageData(image)
				}
				return // Sub trees will end up painting whatever they need so we can just exit out
			}
		}
	}
	// If we got here this whole tree has a single color value so we can paint it
	tree.leaf_value = value
}

// Returns false if splitting failed (probably due to reaching max depth)
func (tree *QuadTreeTerrain) Split() bool {
	if tree.pixel_width < 2 || !tree.leaf {
		return false
	}

	half_w := tree.pixel_width / 2

	tree.sub_trees[0] = NewQuadTreeTerrain(tree.pixel_x, tree.pixel_y, half_w)
	tree.sub_trees[0].leaf_value = tree.leaf_value
	tree.sub_trees[1] = NewQuadTreeTerrain(tree.pixel_x + half_w, tree.pixel_y, half_w)
	tree.sub_trees[1].leaf_value = tree.leaf_value
	tree.sub_trees[2] = NewQuadTreeTerrain(tree.pixel_x, tree.pixel_y + half_w, half_w)
	tree.sub_trees[2].leaf_value = tree.leaf_value
	tree.sub_trees[3] = NewQuadTreeTerrain(tree.pixel_x + half_w, tree.pixel_y + half_w, half_w)
	tree.sub_trees[3].leaf_value = tree.leaf_value

	tree.leaf = false

	return true
}

func (tree *QuadTreeTerrain) Join() {
	if tree.leaf {
		return
	}
	tree.leaf = true
	tree.leaf_value = tree.sub_trees[0].leaf_value
	for _, sub_tree := range tree.sub_trees {
		sub_tree.Join()
	}
	tree.sub_trees[0] = nil
	tree.sub_trees[1] = nil
	tree.sub_trees[2] = nil
	tree.sub_trees[3] = nil
}

// Returns dirtiness
/*func (tree *QuadTreeTerrain) SetShape(shape shapes.Shape, value int) bool {
	if shape == nil {
		fmt.Println("SETTING A NIL SHAPE SOMEHOW")
		return true
	}

	if tree.leaf && tree.leaf_value == value {
		// Don't bother doing anything
		return tree.dirty
	}

	if tree.leaf {
		if shape.Encloses(tree.space) {
			tree.leaf_value = value
			tree.dirty = true
		} else if shape.Intersects(tree.space) {
			split := tree.Split()
			if !split{
				tree.leaf_value = value
				tree.dirty = true
			} else {
				for _, st := range tree.sub_trees{
					if st.SetShape(shape, value){
						tree.dirty = true
					}
				}
			}
		}
	} else {
		if shape.Encloses(tree.space) {
			tree.Join()
			tree.leaf_value = value
			tree.dirty = true
		} else if shape.Intersects(tree.space) {
			for _, st := range tree.sub_trees{
				if st.SetShape(shape, value){
					tree.dirty = true
				}
			}
		}
		
		// optimize as we climb back out of the tree
		if !tree.leaf && tree.sub_trees[0].leaf {
			join := true
			tl_val := tree.sub_trees[0].leaf_value
			for i := 1; i < 4; i++ {
				if !tree.sub_trees[i].leaf || tree.sub_trees[i].leaf_value != tl_val {
					join = false
					break // not everything in this tree has the same value
				}
			}
			if join {
				tree.Join()
				tree.leaf_value = tl_val
			}
		}
	}
	return tree.dirty
}*/

/*func (tree QuadTreeTerrain) GetRayCollision(ray shapes.Ray, ignore_values []int) (float64, int){
	// Returns time of hit and leaf_value of the rect we hit
	if !ray.BoundingBox().Intersects(tree.space) {
		return 2.0, 0
	}

	if tree.leaf {
		if utils.IntSliceContains(ignore_values, tree.leaf_value) {
			return 2.0, 0
		}
		hit := ray.GetAxisRectIntersectionTime(tree.space)
		return hit, tree.leaf_value
	} else {
		t := 2.0
		v := 0
		for _, st := range(tree.sub_trees) {
			hit, val := st.GetRayCollision(ray, ignore_values)
			if hit < t {
				t = hit
				v = val
			}
		}
		return t, v
	}
}

// TODO: we can probably check fewer rect sides if we somehow consider neighbors
func (tree QuadTreeTerrain) GetRayCollisions(ray shapes.Ray, ignore_values []int, results *shapes.RayResultList){
	// Returns time of hit and leaf_value of the rect we hit
	if !ray.BoundingBox().Intersects(tree.space) {
		return
	}

	if tree.leaf {
		if utils.IntSliceContains(ignore_values, tree.leaf_value) {
			return
		}
		t, u := ray.GetAxisRectIntersections(tree.space)
		if t <= 1.0 {
			results.Add(shapes.RayResult{t, u, tree.leaf_value})
		}
	} else {
		for _, st := range(tree.sub_trees) {
			st.GetRayCollisions(ray, ignore_values, results)
		}
	}
}

// TODO: generalize this to other shapes
func (tree QuadTreeTerrain) SweepShape(shape shapes.AxisRect, dirx float64, diry float64, ignore_values []int) (float64, float64, float64) {
	// TODO: broad phase check probably helps performance, and we dont need the hit info unless its a leaf
	t, nx, ny := shape.SweepAxisRect(dirx, diry, tree.space)

	if tree.leaf {
		if utils.IntSliceContains(ignore_values, tree.leaf_value) && t < 2.0 {
			return 2.0, 0.0, 0.0
		}
		return t, nx, ny
	} else if t < 2.0 { // We hit this space so need to search deeper
		mint := 2.0
		wnx := 0.0
		wny := 0.0
		for _, st := range(tree.sub_trees) {
			t, nx, ny = st.SweepShape(shape, dirx, diry, ignore_values)
			if t < mint {
				mint = t
				wnx = nx
				wny = ny
			}
		}
		return mint, wnx, wny
	}
	return 2.0, 0.0, 0.0
}*/

func (tree QuadTreeTerrain) DoesLineCollide(ray Line) bool {
	// TODO: for path finding we need a check like that that considers touching the side of a rectangle as not a collision
	if tree.leaf && tree.leaf_value == 0 {
		return false
	}

	if !ray.BoundingBox().IntersectsAxisRect(tree.space) {
		return false
	}

	if tree.leaf {
		if ray.IntersectsAxisRect(tree.space) {
			return true
		}
	} else {
		for _, st := range(tree.sub_trees) {
			if st.DoesLineCollide(ray) {
				return true
			}
		}
	}

	return false
}

func (tree QuadTreeTerrain) ImprovePath(path []utils.IntPair) []utils.IntPair {
	if len(path) <= 2 {
		return path
	}
	for idx := 0; idx < len(path) - 2; {
		if !tree.DoesLineCollide(NewLine(path[idx].X, path[idx].Y, path[idx+1].X, path[idx+1].Y)) {
			// remove idx+1 from the path
			path = append(path[:idx+1], path[idx+2:]...)
		} else {
			idx += 1
		}
	}
	return path
}

func (tree QuadTreeTerrain) ImprovePathBeginning(path []utils.IntPair) []utils.IntPair {
	// like ImprovePath, but we stop after our first ray hit. The idea is that a pathfinding user will be calling this as they follow the path
	if len(path) <= 2 {
		return path
	}
	for idx := 0; idx < len(path) - 2; {
		if !tree.DoesLineCollide(NewLine(path[idx].X, path[idx].Y, path[idx+1].X, path[idx+1].Y)) {
			// remove idx+1 from the path
			path = append(path[:idx+1], path[idx+2:]...)
		} else {
			return path
		}
	}
	return path
}

// Returns false if starting position is not inside this tree, returned position will not be useful in that case
// Doesn't actually return the closest corner, just the closest corner of the leaf we are in, which is close enough
func (tree QuadTreeTerrain) GetClosestCorner(x, y float64) (bool, int, int) {
	if !tree.space.ContainsPoint(x, y) {
		return false, 0, 0
	}

	if !tree.leaf {
		inside, rx, ry := tree.sub_trees[0].GetClosestCorner(x, y)
		if inside {
			return true, rx, ry
		}

		inside, rx, ry = tree.sub_trees[1].GetClosestCorner(x, y)
		if inside {
			return true, rx, ry
		}

		inside, rx, ry = tree.sub_trees[2].GetClosestCorner(x, y)
		if inside {
			return true, rx, ry
		}

		inside, rx, ry = tree.sub_trees[3].GetClosestCorner(x, y)
		if inside {
			return true, rx, ry
		}
	}

	rx := tree.space.X()
	ry := tree.space.Y()
	if x > float64(rx + (tree.space.W() / 2)) {
		rx = tree.space.X2()
	}
	if y > float64(ry + (tree.space.H() / 2)) {
		ry = tree.space.Y2()
	}

	return true, rx, ry
}

func (tree QuadTreeTerrain) Neighbours(n utils.IntPair) []utils.IntPair {
	// func to implement the astar library's graph interface
	return tree.GetAdjacentCorners(n.X, n.Y)
}

func (tree QuadTreeTerrain) GetAdjacentCorners(x, y int) []utils.IntPair {
	// This function assumes you are passing in an actual corner
	if !tree.leaf {
		flags := [4]bool{}
		result := [4]utils.IntPair{}
		min_pos_x := math.MaxInt
		min_neg_x := math.MinInt
		min_pos_y := math.MaxInt
		min_neg_y := math.MinInt
		for _, subtree := range tree.sub_trees { // TODO: we know enough about the max number of adjacent corners to reasonably be able to avoid all this slice fiddling
			if subtree.space.ContainsPoint(float64(x), float64(y)) {
				nodes := subtree.GetAdjacentCorners(x, y)
				for _, newnode := range nodes {
					if newnode.X == x {
						if newnode.Y > y && newnode.Y < min_pos_y {
							min_pos_y = newnode.Y
							result[0] = newnode
							flags[0] = true
						} else if newnode.Y < y && newnode.Y > min_neg_y {
							min_neg_y = newnode.Y
							result[1] = newnode
							flags[1] = true
						}
					} else if newnode.Y == y {
						if newnode.X > x && newnode.X < min_pos_x {
							min_pos_x = newnode.X
							result[2] = newnode
							flags[2] = true
						} else if newnode.X < x && newnode.X > min_neg_x {
							min_neg_x = newnode.X
							result[3] = newnode
							flags[3] = true
						}
					}
				}
			}
		}

		slice_result := make([]utils.IntPair, 0, 4)
		for idx, flag := range flags {
			if flag {
				slice_result = append(slice_result, result[idx])
			}
		}
		return slice_result
	}

	// TODO: don't hard code what is walkable like this probably
	// TODO: this isnt really an accurate way to do this, but it might be good enough, should explore alternatives though
	if tree.leaf_value == 1 {
		return []utils.IntPair{}
	}

	result := make([]utils.IntPair, 0, 2)
	if tree.space.X() == x {
		if tree.space.Y() == y {
			result = append(result, utils.IntPair{X: tree.space.X(), Y: tree.space.Y2()})
			result = append(result, utils.IntPair{X: tree.space.X2(), Y: tree.space.Y()})
		} else if tree.space.Y2() == y {
			result = append(result, utils.IntPair{X: tree.space.X(), Y: tree.space.Y()})
			result = append(result, utils.IntPair{X: tree.space.X2(), Y: tree.space.Y2()})
		}
	} else if tree.space.X2() == x {
		if tree.space.Y() == y {
			result = append(result, utils.IntPair{X: tree.space.X2(), Y: tree.space.Y2()})
			result = append(result, utils.IntPair{X: tree.space.X(), Y: tree.space.Y()})
		} else if tree.space.Y2() == y {
			result = append(result, utils.IntPair{X: tree.space.X2(), Y: tree.space.Y()})
			result = append(result, utils.IntPair{X: tree.space.X(), Y: tree.space.Y2()})
		}
	}

	return result
}

func (tree QuadTreeTerrain) CircleSeparation(circ Circle) (utils.FloatPair, float64) {
	if tree.leaf && tree.leaf_value == 0 {
		return utils.FloatPair{}, 0
	}

	if !circ.BBIntersectsAxisRect(tree.space) {
		return utils.FloatPair{}, 0
	}

	if tree.leaf {
		return circ.SeparationForAxisRect(tree.space)
	} else {
		maxvec, maxlen := tree.sub_trees[0].CircleSeparation(circ)
		
		for _, st := range(tree.sub_trees[1:]) {
			vec, len := st.CircleSeparation(circ)
			if len > maxlen {
				maxlen = len
				maxvec = vec
			}
		}

		return maxvec, maxlen
	}
}

func OctileDistance(s, e utils.IntPair) (float64) {
	dx := math.Abs(float64(s.X - e.X))
    dy := math.Abs(float64(s.Y - e.Y))
    return (dx + dy) + (Sqrt2 - 2) * math.Min(dx, dy)
}

func ManhattanDistance(s, e utils.IntPair) (float64) {
	dx := math.Abs(float64(s.X - e.X))
    dy := math.Abs(float64(s.Y - e.Y))
    return (dx + dy)
}

func EuclidianDistance(s, e utils.IntPair) (float64) {
	dx := math.Abs(float64(s.X - e.X))
    dy := math.Abs(float64(s.Y - e.Y))
    return math.Sqrt((dx * dx) + (dy * dy))
}

func (tree QuadTreeTerrain) FindPath(s, e utils.IntPair) []utils.IntPair {
	// check if end is inside unwalkable terrain to save a lot of time
	_, sx, sy := tree.GetClosestCorner(float64(s.X), float64(s.Y))
	_, ex, ey := tree.GetClosestCorner(float64(e.X), float64(e.Y))
	s = utils.IntPair{X: sx, Y:sy,}
	e = utils.IntPair{X: ex, Y:ey,}
	end_adjacent := tree.GetAdjacentCorners(e.X, e.Y)
	if len(end_adjacent) <= 0 {
		return []utils.IntPair{}
	}

	result := astar.FindPath[utils.IntPair](tree, s, e, ManhattanDistance, EuclidianDistance)
	result = append([]utils.IntPair{s}, result...)
	result = append(result, e)
	return result
}

func (tree QuadTreeTerrain) DebugDrawTree(target *ebiten.Image) {
	ebitenutil.DrawLine(target, float64(tree.space.X()), float64(tree.space.Y()), float64(tree.space.X2()), float64(tree.space.Y()), color.RGBA{255, 0, 0, 255})
	ebitenutil.DrawLine(target, float64(tree.space.X()), float64(tree.space.Y()), float64(tree.space.X()), float64(tree.space.Y2()), color.RGBA{255, 0, 0, 255})
	if !tree.leaf {
		tree.sub_trees[0].DebugDrawTree(target)
		tree.sub_trees[1].DebugDrawTree(target)
		tree.sub_trees[2].DebugDrawTree(target)
		tree.sub_trees[3].DebugDrawTree(target)
	}
}

func (tree QuadTreeTerrain) DebugDrawTree2(target *ebiten.Image) {
	if !tree.leaf {
		tree.sub_trees[0].DebugDrawTree2(target)
		tree.sub_trees[1].DebugDrawTree2(target)
		tree.sub_trees[2].DebugDrawTree2(target)
		tree.sub_trees[3].DebugDrawTree2(target)
	} else {
		if tree.leaf_value == 0 {
			ebitenutil.DrawRect(target, float64(tree.pixel_x), float64(tree.pixel_y), float64(tree.pixel_width), float64(tree.pixel_width), color.RGBA{200, 200, 200, 255})
		} else {
			ebitenutil.DrawRect(target, float64(tree.pixel_x), float64(tree.pixel_y), float64(tree.pixel_width), float64(tree.pixel_width), color.RGBA{50, 50, 50, 255})
		}
	}
}
