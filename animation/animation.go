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
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	FileSystem = os.DirFS(".")

	SHEET_CACHE = make(map[string]*ebiten.Image) // We never unload cached data yet, but if need be we can just clear this variable
	ANIMATION_MAP_CACHE = make(map[string](map[string]Animation)) // We cache the full maps too to save on json parsing and file reading
	META_ANIMATION_MAP_CACHE = make(map[string](map[string]MetaAnimationList))
)

func NewEbitenImage(path string) (*ebiten.Image) {
	cached_image, cached := SHEET_CACHE[path]
	if cached {
		return cached_image
	}

	dat, err := fs.ReadFile(FileSystem, filepath.ToSlash(path))
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(bytes.NewReader(dat))
	if err != nil {
		panic(err)
	}
	result := ebiten.NewImageFromImage(img)
	SHEET_CACHE[path] = result
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

type MetaAnimationDefinition struct {
    Length float64 `json:"length"`
    AnimNames map[string]string  `json:"anim_names"`
    XOffsets map[string]float64  `json:"x_offsets"`
    YOffsets map[string]float64  `json:"y_offsets"`
    Scales map[string]float64  `json:"scales"`
    Times map[string]float64  `json:"times"` // currently just functions as a time offset
}

type MetaAnimationMapDefinition struct {
	AnimationsPath string `json:"animations_path"`
	MetaAnimations map[string][]MetaAnimationDefinition `json:"meta_animations"`
}

type Font struct {
	Animations []Animation
	CharacterWidth int
	CharacterHeight int
}

func NewFont(image_path string, char_w, char_h, char_per_line, height int) (Font, error) {
	result := Font{}
	result.CharacterWidth = char_w
	result.CharacterHeight = char_h

	// A quick way to create an animation map for a static monospaced font
	// loads animations into the array in english reading order, having the sprite sheet be in ASCII character order would be smart
	var animations []Animation

	for y := 0; y < height; y++ {
		for x := 0; x < char_per_line; x++ {
			frames := make(map[string]int)
			frames["0"] = 0
			def := AnimationDefinition{
				SheetPath: image_path,
				FrameWidth: char_w,
				FrameHeight: char_h,
				FrameOffsetX: x * char_w,
				FrameOffsetY: y * char_h,
				Length: 1.0,
				Frames: frames,
			}
			animations = append(animations, NewAnimationFromDefinition(def, ""))
		}
	}

	result.Animations = animations

	return result, nil
}

func (font Font) DrawText(target *ebiten.Image, xpos, ypos, scale, time float64, runes []rune) {
	// TODO: handle newlines and such
	for column, char := range runes {
		// TODO: For some reason we are a row off, hence "- 32", should probably fix that before any numbered version, since it will break things
		font.Animations[char - 32].Draw(
			target,
			xpos + float64(column * font.CharacterWidth) * scale,
			ypos,
			scale,
			time)
	}
}

func LoadAnimationMap(path string) (map[string]Animation, error) {
	cached_map, cached := ANIMATION_MAP_CACHE[path]
	if cached {
		return cached_map, nil
	}

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

	ANIMATION_MAP_CACHE[path] = animations
	return animations, nil
}

func LoadMetaAnimationMap(path string) (map[string]MetaAnimationList, error) {
	cached_map, cached := META_ANIMATION_MAP_CACHE[path]
	if cached {
		return cached_map, nil
	}

	bytes, err := fs.ReadFile(FileSystem, path)
	if err != nil {
		return nil, err
	}

	var def MetaAnimationMapDefinition
	err = json.Unmarshal(bytes, &def)
	if err != nil {
		return nil, err
	}

	anims, aerr := LoadAnimationMap(def.AnimationsPath) // TODO: should make this load relative to the metaanim file, right now its relative where we are running
	if aerr != nil {
		return nil, aerr
	}

	animations := make(map[string]MetaAnimationList)
	for name, animdef := range def.MetaAnimations {
		animations[name] = NewMetaAnimationFromDefinition(animdef, filepath.Dir(path), anims)
	}

	META_ANIMATION_MAP_CACHE[path] = animations
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

func getReverseSortedSliceOfKeys4(data map[string]string) ([]float64){
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

func floatifyKeys4(data map[string]string) (map[float64]string) {
	result := make(map[float64]string)
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

type MetaAnimationList struct {
	meta_anims []MetaAnimation
}

func (m MetaAnimationList) Draw(target *ebiten.Image, xpos, ypos, scale, time float64) {
	for _, anim := range m.meta_anims {
		anim.Draw(target, xpos, ypos, scale, time)
	}
}

func (m MetaAnimationList) GetLength() float64 {
	if len(m.meta_anims) > 0 {
		return m.meta_anims[0].GetLength()
	}
	return 0.0
}

func NewMetaAnimationFromDefinition(defs []MetaAnimationDefinition, parent_path string, anims map[string]Animation) (MetaAnimationList) {
	result := MetaAnimationList{}
	for _, def := range defs {
		// Some defaults in case we are missing values
		// This feels messy, maybe there is a better way
		if def.AnimNames == nil {
			def.AnimNames = make(map[string]string)
		}
		if len(def.AnimNames) == 0 {
			def.AnimNames["0.0"] = "default"
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
		if def.Scales == nil {
			def.Scales = make(map[string]float64)
		}
		if len(def.Scales) == 0 {
			def.Scales["0.0"] = 1.0
		}
		if def.Times == nil {
			def.Times = make(map[string]float64)
		}
		if len(def.Times) == 0 {
			def.Times["0.0"] = 0.0
		}

		result.meta_anims = append(result.meta_anims, MetaAnimation{
			animations: anims,
			length: def.Length,
			anim_name: floatifyKeys4(def.AnimNames),
			anim_name_keys: getReverseSortedSliceOfKeys4(def.AnimNames),
			x_offset: floatifyKeys2(def.XOffsets),
			x_offset_keys: getReverseSortedSliceOfKeys2(def.XOffsets),
			y_offset: floatifyKeys2(def.YOffsets),
			y_offset_keys: getReverseSortedSliceOfKeys2(def.YOffsets),
			scale: floatifyKeys2(def.Scales),
			scale_keys: getReverseSortedSliceOfKeys2(def.Scales),
			time: floatifyKeys2(def.Times),
			time_keys: getReverseSortedSliceOfKeys2(def.Times),
		})
	}

	return result
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

func getStringAtTime(values map[float64]string, keys []float64, time, max_time float64) (string) {
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

func maybeNegate(value float64, negate bool) (float64) {
	if negate {
		return -value
	}
	return value
}

// TODO: currently animations are always drawn with the origin at the bottom middle of the sprite, give other options
// TODO: need to have draw methods actually just add to a depth sorted queue to draw with later, so we can meaningfully have a zpos
func (anim Animation) Draw(target *ebiten.Image, xpos, ypos, scale, time float64) {
	if len(anim.frames_keys) == 0 {
		fmt.Println("Attempting to play animation with no frames!")
		return
	}
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

func (anim Animation) GetLength() float64 {
	return anim.length
}

type MetaAnimation struct {
	animations map[string]Animation
	length float64
	anim_name map[float64]string
	anim_name_keys []float64
	x_offset map[float64]float64
	x_offset_keys []float64 // All other keys should also be sorted highest to lowest to keep things consistent
	y_offset map[float64]float64
	y_offset_keys []float64
	scale map[float64]float64
	scale_keys []float64
	time map[float64]float64
	time_keys []float64
}

func (anim MetaAnimation) GetXOffset(time float64) (float64) {
	return getInterpolatedValueFromReversedTimeKeysAndValueMap(anim.x_offset, anim.x_offset_keys, time, anim.length)
}

func (anim MetaAnimation) GetYOffset(time float64) (float64) {
	return getInterpolatedValueFromReversedTimeKeysAndValueMap(anim.y_offset, anim.y_offset_keys, time, anim.length)
}

func (anim MetaAnimation) GetScale(time float64) (float64) {
	return getInterpolatedValueFromReversedTimeKeysAndValueMap(anim.scale, anim.scale_keys, time, anim.length)
}

func (anim MetaAnimation) GetAnimName(time float64) (string) {
	return getStringAtTime(anim.anim_name, anim.anim_name_keys, time, anim.length)
}

func (anim MetaAnimation) GetTime(time float64) (float64) {
	return getInterpolatedValueFromReversedTimeKeysAndValueMap(anim.time, anim.time_keys, time, anim.length)
}

func (anim MetaAnimation) GetLength() float64 {
	return anim.length
}

func (anim MetaAnimation) Draw(target *ebiten.Image, xpos, ypos, scale, time float64) {
	anim.animations[anim.GetAnimName(time)].Draw(
		target,
		anim.GetXOffset(time) + xpos,
		anim.GetYOffset(time) + ypos,
		anim.GetScale(time) * scale,
		anim.GetTime(time) + time)
}


type PlayableAnimation interface {
	Draw(*ebiten.Image, float64, float64, float64, float64)
	GetLength() float64
}

type AnimationPlayer struct {
	anim PlayableAnimation
	start_time float64
	xpos, ypos, scale float64
}

// TODO: would be nice to have a handler that stores these animations and manages deleting them and such
// We would need to handle camera positioning and such in here instead of letting the client handle it
// Would be a good place to handle depth as well so client doesn't need to sort things themselves
// This is basically making this a higher level engine, as long as its optional probably a good thing

func NewAnimationPlayer(anim PlayableAnimation, xpos, ypos, scale, time float64) AnimationPlayer {
	return AnimationPlayer {
		anim: anim,
		start_time: time,
		xpos: xpos,
		ypos: ypos,
		scale: scale,
	}
}

func (p AnimationPlayer) Draw(target *ebiten.Image, time float64) bool {
	// Time adjusted to the start_time and forced to not loop the animation
	// Returns if the animation is still playing (time hasnt passed the end)
	t := time - p.start_time
	if t > p.anim.GetLength() {
		p.anim.Draw(target, p.xpos, p.ypos, p.scale, p.anim.GetLength())
		return false
	} 

	p.anim.Draw(target, p.xpos, p.ypos, p.scale, t)
	return true
}
