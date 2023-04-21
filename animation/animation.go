package animation

import (
	"image"
	_ "image/png"
	"math"
	"sort"
	"strconv"
	"os"
	"io/fs"
	"encoding/json"
	"path/filepath"
	"bytes"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	FileSystem = os.DirFS(".")
)

func NewEbitenImage(path string) (*ebiten.Image) {
	dat, err := fs.ReadFile(FileSystem, filepath.ToSlash(path))
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(bytes.NewReader(dat))
	if err != nil {
		panic(err)
	}
	result := ebiten.NewImageFromImage(img)
	return result
}

type AnimationDefinition struct {
    SheetPath string `json:"sheet_path"`
    FrameWidth int `json:"frame_width"`
    FrameHeight int `json:"frame_height"`
    FrameOffsetX int `json:"frame_offset_x"`
    FrameOffsetY int `json:"frame_offset_y"`

    Length float64 `json:"length"`
    Frames map[string]int  `json:"frames"`
    XOffsets map[string]float64  `json:"x_offsets"`
    YOffsets map[string]float64  `json:"y_offsets"`
    WScales map[string]float64  `json:"w_scales"`
    HScales map[string]float64  `json:"h_scales"`
    WMirrors map[string]bool  `json:"w_mirrors"`
    HMirrors map[string]bool  `json:"h_mirrors"`
    Rotations map[string]float64  `json:"rotations"`
}

type AnimationMapDefinition struct {
	Animations map[string]AnimationDefinition `json:"animations"`
}

func LoadAnimationMap(path string) (map[string]Animation, error) {
	bytes, err := fs.ReadFile(FileSystem, path)
	if err != nil {
		return nil, err
	}

	var def AnimationMapDefinition
	err = json.Unmarshal(bytes, &def)
	if err != nil {
		return nil, err
	}

	animations := make(map[string]Animation)
	for name, animdef := range def.Animations {
		animations[name] = NewAnimationFromDefinition(animdef, filepath.Dir(path))
	}

	return animations, nil
}

func getReverseSortedSliceOfKeys1(data map[string]int) ([]float64){ // TODO: surely the value type of the map doesn't matter, how does golang work help
	keys := make([]float64, len(data))

	i := 0
	for k := range data {
		key_float, err := strconv.ParseFloat(k, 64)
		if err != nil {
			panic(err)
		}
	    keys[i] = key_float
	    i++
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] > keys[j] })
	return keys
}

func getReverseSortedSliceOfKeys2(data map[string]float64) ([]float64){
	keys := make([]float64, len(data))

	i := 0
	for k := range data {
		key_float, err := strconv.ParseFloat(k, 64)
		if err != nil {
			panic(err)
		}
	    keys[i] = key_float
	    i++
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] > keys[j] })
	return keys
}

func getReverseSortedSliceOfKeys3(data map[string]bool) ([]float64){
	keys := make([]float64, len(data))

	i := 0
	for k := range data {
		key_float, err := strconv.ParseFloat(k, 64)
		if err != nil {
			panic(err)
		}
	    keys[i] = key_float
	    i++
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] > keys[j] })
	return keys
}

func floatifyKeys2(data map[string]float64) (map[float64]float64) {
	result := make(map[float64]float64)
	for k, v := range data {
		key_float, err := strconv.ParseFloat(k, 64)
		if err != nil {
			panic(err)
		}
		result[key_float] = v
	}
	return result
}

func floatifyKeys3(data map[string]bool) (map[float64]bool) {
	result := make(map[float64]bool)
	for k, v := range data {
		key_float, err := strconv.ParseFloat(k, 64)
		if err != nil {
			panic(err)
		}
		result[key_float] = v
	}
	return result
}

func NewAnimationFromDefinition(def AnimationDefinition, parent_path string) (Animation) {
	sheet := NewEbitenImage(filepath.Join(parent_path, def.SheetPath)) // TODO: cache the sheet instead of loading a new one for each animation
	width_in_frames := (sheet.Bounds().Max.X - def.FrameOffsetX) / def.FrameWidth
	
	// Some defaults in case we are missing values
	// This feels messy, maybe there is a better way
	if def.Frames == nil {
		def.Frames = make(map[string]int)
	}
	if len(def.Frames) == 0 {
		def.Frames["0.0"] = 0
	}
	if def.XOffsets == nil {
		def.XOffsets = make(map[string]float64)
	}
	if len(def.XOffsets) == 0 {
		def.XOffsets["0.0"] = 0.0
	}
	if def.YOffsets == nil {
		def.YOffsets = make(map[string]float64)
	}
	if len(def.YOffsets) == 0 {
		def.YOffsets["0.0"] = 0.0
	}
	if def.WScales == nil {
		def.WScales = make(map[string]float64)
	}
	if len(def.WScales) == 0 {
		def.WScales["0.0"] = 1.0
	}
	if def.HScales == nil {
		def.HScales = make(map[string]float64)
	}
	if len(def.HScales) == 0 {
		def.HScales["0.0"] = 1.0
	}
	if def.WMirrors == nil {
		def.WMirrors = make(map[string]bool)
	}
	if len(def.WMirrors) == 0 {
		def.WMirrors["0.0"] = false
	}
	if def.HMirrors == nil {
		def.HMirrors = make(map[string]bool)
	}
	if len(def.HMirrors) == 0 {
		def.HMirrors["0.0"] = false
	}
	if def.Rotations == nil {
		def.Rotations = make(map[string]float64)
	}
	if len(def.Rotations) == 0 {
		def.Rotations["0.0"] = 0.0
	}

	frames := make(map[float64]image.Rectangle)
	for key, frame := range def.Frames {
		key_float, err := strconv.ParseFloat(key, 64)
		if err != nil {
			panic(err)
		}
		x := frame % width_in_frames
		y := frame / width_in_frames
		frames[key_float] = image.Rect((x * def.FrameWidth) + def.FrameOffsetX, (y * def.FrameHeight) + def.FrameOffsetY, ((x + 1) * def.FrameWidth) + def.FrameOffsetX, ((y + 1) * def.FrameHeight) + def.FrameOffsetY)
	}
	frames_keys := getReverseSortedSliceOfKeys1(def.Frames)

	return Animation{
		length: def.Length,
		Sheet: sheet,
		frames: frames,
		frames_keys: frames_keys,
		x_offset: floatifyKeys2(def.XOffsets),
		x_offset_keys: getReverseSortedSliceOfKeys2(def.XOffsets),
		y_offset: floatifyKeys2(def.YOffsets),
		y_offset_keys: getReverseSortedSliceOfKeys2(def.YOffsets),
		w_scale: floatifyKeys2(def.WScales),
		w_scale_keys: getReverseSortedSliceOfKeys2(def.WScales),
		h_scale: floatifyKeys2(def.HScales),
		h_scale_keys: getReverseSortedSliceOfKeys2(def.HScales),
		w_mirror: floatifyKeys3(def.WMirrors),
		w_mirror_keys: getReverseSortedSliceOfKeys3(def.WMirrors),
		h_mirror: floatifyKeys3(def.HMirrors),
		h_mirror_keys: getReverseSortedSliceOfKeys3(def.HMirrors),
		rotation: floatifyKeys2(def.Rotations),
		rotation_keys: getReverseSortedSliceOfKeys2(def.Rotations),
	}
}

type Animation struct {
	length float64
	Sheet *ebiten.Image
	frames map[float64]image.Rectangle
	frames_keys []float64 // Keys should be sorted from highest to lowest to simplify finding frame values
	x_offset map[float64]float64
	x_offset_keys []float64 // All other keys should also be sorted highest to lowest to keep things consistent
	y_offset map[float64]float64
	y_offset_keys []float64
	w_scale map[float64]float64
	w_scale_keys []float64
	h_scale map[float64]float64
	h_scale_keys []float64
	w_mirror map[float64]bool
	w_mirror_keys []float64
	h_mirror map[float64]bool
	h_mirror_keys []float64
	rotation map[float64]float64
	rotation_keys []float64
}

func (anim Animation) GetFrameRect(time float64) (image.Rectangle) {
	time = math.Mod(time, anim.length)
	for _, key := range anim.frames_keys {
		if key <= time {
			return anim.frames[key]
		}
	}
	return anim.frames[anim.frames_keys[len(anim.frames_keys)-1]]
}

func interpolate(start, start_time, end, end_time, time float64) (float64) {
	fraction := (time - start_time) / (end_time - start_time)
	return start + ((end - start) * fraction)
}

// TODO: functions for interpolation besides linear
func getInterpolatedValueFromReversedTimeKeysAndValueMap(values map[float64]float64, keys []float64, time, max_time float64) (float64) {
	time = math.Mod(time, max_time)
	for idx, start_time := range keys {
		if start_time <= time {
			end_time := max_time
			end_value := 0.0
			if idx > 0 {
				end_time = keys[idx - 1]
				end_value = values[end_time]
			} else {
				end_value = values[keys[0]]
			}
			return interpolate(values[start_time], start_time, end_value, end_time, time)
		}
	}
	return values[keys[len(keys)-1]]
}

func getBoolAtTime(values map[float64]bool, keys []float64, time, max_time float64) (bool) {
	time = math.Mod(time, max_time)
	for _, key := range keys {
		if key <= time {
			return values[key]
		}
	}
	return values[keys[len(keys)-1]]
}

func (anim Animation) GetXOffset(time float64) (float64) {
	return getInterpolatedValueFromReversedTimeKeysAndValueMap(anim.x_offset, anim.x_offset_keys, time, anim.length)
}

func (anim Animation) GetYOffset(time float64) (float64) {
	return getInterpolatedValueFromReversedTimeKeysAndValueMap(anim.y_offset, anim.y_offset_keys, time, anim.length)
}

func (anim Animation) GetWScale(time float64) (float64) {
	return getInterpolatedValueFromReversedTimeKeysAndValueMap(anim.w_scale, anim.w_scale_keys, time, anim.length)
}

func (anim Animation) GetHScale(time float64) (float64) {
	return getInterpolatedValueFromReversedTimeKeysAndValueMap(anim.h_scale, anim.h_scale_keys, time, anim.length)
}

func (anim Animation) GetWMirror(time float64) (bool) {
	return getBoolAtTime(anim.w_mirror, anim.w_mirror_keys, time, anim.length)
}

func (anim Animation) GetHMirror(time float64) (bool) {
	return getBoolAtTime(anim.h_mirror, anim.h_mirror_keys, time, anim.length)
}

func (anim Animation) GetRotation(time float64) (float64) {
	return getInterpolatedValueFromReversedTimeKeysAndValueMap(anim.rotation, anim.rotation_keys, time, anim.length)
}

// TODO: currently animations are always drawn with the origin at the bottom middle of the sprite, give other options
// TODO: need to have draw methods actually just add to a depth sorted queue to draw with later, so we can meaningfully have a zpos
func (anim Animation) Draw(target *ebiten.Image, xpos, ypos, scale, time float64) {
	subrect := anim.GetFrameRect(time)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(-subrect.Dx()/2.0), float64(-subrect.Dy()/2.0))
	// TODO: HMirror introduces a slight verticle wobble, might just have to do with how we translate before it
	op.GeoM.Scale(maybeNegate(anim.GetWScale(time), anim.GetWMirror(time)), maybeNegate(anim.GetHScale(time), anim.GetHMirror(time)))
	op.GeoM.Translate(0.0, float64(-subrect.Dy()/2.0))
	op.GeoM.Scale(scale, scale)
	op.GeoM.Rotate(anim.GetRotation(time))
	op.GeoM.Translate(xpos, ypos)
	op.GeoM.Translate(anim.GetXOffset(time), anim.GetYOffset(time))

	target.DrawImage(anim.Sheet.SubImage(subrect).(*ebiten.Image), op)
}
