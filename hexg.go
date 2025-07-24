// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package hexg implements hex grids in the style of the Red Blob Games blog
// https://www.redblobgames.com/grids/hexagons/ and
// https://www.redblobgames.com/grids/hexagons/implementation.html.
//
// Many comments are lifted from the Red Blog Games page and are copyright
// by them.
package hexg

import (
	"fmt"
	"math"
	"strconv"
	"unicode"
)

// types

// Hex implements Cube coordinates for hexes.
type Hex struct {
	q, r, s int
}

// FractionalHex implements a hex with floating point coordinates.
// Used for linear interpolation and rounding.
type FractionalHex struct {
	q, r, s float64
}

// OffsetCoord implements offset coordinates for hexes.
type OffsetCoord struct {
	col, row int
}

// Point represents a screen coordinate
type Point struct {
	X float64
	Y float64
}

// Layout represents the orientation of a hexagonal grid
type Layout struct {
	orientation orientation

	// size is used for scaling (eg, matching pixel sprite sizes)
	// set it to Point(size, size) if you need uniform scaling.
	// todo: this makes no sense. i need to read the page again.
	size Point

	// origin is the center of the q=0,r=0 hexagon.
	// set it to Point(0, 0) if you do not need to translate the transformation.
	origin Point

	// offset is
	offset int
}

// GridStore is a map of Hex indexed by the hash of the Hex.
// I think that it's here for examples of how to use hex.Hash.
type GridStore map[uint64]Hex

// orientation is a helper for Layouts
type orientation struct {
	f0, f1, f2, f3 float64
	b0, b1, b2, b3 float64
	start_angle    float64 // in multiples of 60°
}

// 1.0 Hex coordinates

// constructors

// NewHex returns a Hex initialized with Cube coordinates.
// Panics on invalid inputs.
func NewHex(q, r, s int) Hex {
	if q+r+s != 0 {
		panic("assert (q + r + s == 0)")
	}
	return Hex{q: q, r: r, s: s}
}

// NewHexFromAxialCoords returns a Hex initialized with Axial coordinates.
// Computes s from q and s, so will never panic on inputs.
func NewHexFromAxialCoords(q, r int) Hex {
	return Hex{q: q, r: r, s: -q - r}
}

// ConciseString returns the coordinates with signs.
// It returns the coordinates formatted as (+q+r+s).
func (h Hex) ConciseString() string {
	return fmt.Sprintf("%+d%+d%+d", h.q, h.r, h.s)
}

// String implements the Stringer interface.
// It returns the coordinates formatted as (q,r,s).
func (h Hex) String() string {
	return fmt.Sprintf("%d,%d,%d", h.q, h.r, h.s)
}

// 1.1 Equality

// Equals returns true if the two hexes have the same coordinates.
func (h Hex) Equals(b Hex) bool {
	return h.q == b.q && h.r == b.r && h.s == b.s
}

// NotEquals returns true if the two hexes have different coordinates.
func (h Hex) NotEquals(b Hex) bool {
	return h.q != b.q && h.r != b.r && h.s != b.s
}

// 1.2 Coordinate arithmetic

// Add returns the sum of two hexes.
func (h Hex) Add(b Hex) Hex {
	return Hex{q: h.q + b.q, r: h.r + b.r, s: h.s + b.s}
}

// Subtract returns the difference of two hexes.
func (h Hex) Subtract(b Hex) Hex {
	return Hex{q: h.q - b.q, r: h.r - b.r, s: h.s - b.s}
}

// Multiply returns a Hex scaled by an integer.
func (h Hex) Multiply(k int) Hex {
	return Hex{q: h.q * k, r: h.r * k, s: h.s * k}
}

// 1.3 Distance

// Length is the line from the origin to a hex
func (h Hex) Length() int {
	return int((abs(h.q) + abs(h.r) + abs(h.s)) / 2)
}

// Distance between two hexes is the length of the line between them.
func (h Hex) Distance(b Hex) int {
	return h.Subtract(b).Length()
}

// 1.3.1 Neighbors

// Neighbor returns the hex that is one step away in the given direction.
// Direction is coerced to the range 0..5.
func (h Hex) Neighbor(direction int) Hex {
	return h.Add(Direction(direction))
}

// hex_directions has the offset to the neighboring hex indexed by direction 0..5
var hex_directions = [6]Hex{
	{1, 0, -1}, {1, -1, 0}, {0, -1, 1},
	{-1, 0, 1}, {-1, 1, 0}, {0, 1, -1},
}

// Direction returns the q, r, and s offsets to use based on the direction
func Direction(direction int) Hex {
	return hex_directions[(6+(direction%6))%6]
}

// 2.0 Layout

// there really are only two orientations - pointy top and flat top.
// define them here to prevent errors.

var (
	// layout_pointy returns a pointy top orientation.
	layout_pointy = orientation{
		f0: math.Sqrt(3.0), f1: math.Sqrt(3.0) / 2.0, f2: 0.0, f3: 3.0 / 2.0,
		b0: math.Sqrt(3.0) / 3.0, b1: -1.0 / 3.0, b2: 0.0, b3: 2.0 / 3.0,
		start_angle: 0.5,
	}

	// layout_flat returns a flat top orientation.
	layout_flat = orientation{
		f0: 3.0 / 2.0, f1: 0.0, f2: math.Sqrt(3.0) / 2.0, f3: math.Sqrt(3.0),
		b0: 2.0 / 3.0, b1: 0.0, b2: -1.0 / 3.0, b3: math.Sqrt(3.0) / 3.0,
		start_angle: 0.0,
	}
)

// new_layout creates a new Layout with the orientation, size, and origin.
// it is not exported since it is only used as a helper for the NewLayoutXXX constructors.
func new_layout(o orientation, size, origin Point, offset int) Layout {
	return Layout{
		orientation: o,
		size:        size,
		origin:      origin,
		offset:      offset,
	}
}

// NewLayoutFlat returns a layout with flat-top hexes.
// Size and origin are used when calculating screen pixels.
func NewLayoutFlat(size, origin Point, shoveOddColumnsDown bool) Layout {
	if shoveOddColumnsDown {
		return new_layout(layout_flat, size, origin, odd_q)
	}
	return new_layout(layout_flat, size, origin, even_q)
}

// NewLayoutPointy returns a layout with point-top hexes.
// Size and origin are used when calculating screen pixels.
func NewLayoutPointy(size, origin Point, shoveOddRowsRight bool) Layout {
	if shoveOddRowsRight {
		return new_layout(layout_pointy, size, origin, odd_r)
	}
	return new_layout(layout_pointy, size, origin, even_r)
}

// IsFlatTop returns true if the layout was created with flat-top hexes.
func (layout Layout) IsFlatTop() bool {
	return layout.offset == odd_q || layout.offset == even_q
}

// IsPointyTop returns true if the layout was created with point-top hexes.
func (layout Layout) IsPointyTop() bool {
	return layout.offset == odd_r || layout.offset == even_r
}

// NewPoint returns a new Point with specified screen coordinates
func NewPoint(x, y float64) Point {
	return Point{X: x, Y: y}
}

func (p Point) String() string {
	return fmt.Sprintf("%g,%g", p.X, p.Y)
}

// 2.1 Hex to screen

// ToPixel returns the origin of the hex on the grid as a Point.
func (h Hex) ToPixel(layout Layout) Point {
	return ToPixel(layout, h)
}

// HexToPixel returns the origin of the hex on the grid as a Point.
func (layout Layout) HexToPixel(h Hex) Point {
	return ToPixel(layout, h)
}

// ToPixel returns the origin of the hex on the grid as a Point.
func ToPixel(layout Layout, h Hex) Point {
	M := layout.orientation
	return Point{
		X: layout.origin.X + (M.f0*float64(h.q)+M.f1*float64(h.r))*layout.size.X,
		Y: layout.origin.Y + (M.f2*float64(h.q)+M.f3*float64(h.r))*layout.size.Y,
	}
}

// 2.2 Screen to hex

// ToFractionalHex returns the fractional hex that encloses the pixel.
// In theory, the origin of that fractional hex will be the pixel.
func (p Point) ToFractionalHex(layout Layout) FractionalHex {
	return PixelToFractionalHex(layout, p)
}

// PixelToFractionalHex returns the fractional hex that encloses the pixel.
// In theory, the origin of that fractional hex will be the pixel.
func (layout Layout) PixelToFractionalHex(p Point) FractionalHex {
	return PixelToFractionalHex(layout, p)
}

// PixelToFractionalHex returns the fractional hex that encloses the pixel.
// In theory, the origin of that fractional hex will be the pixel.
func PixelToFractionalHex(layout Layout, p Point) FractionalHex {
	M := layout.orientation
	pt := Point{X: (p.X - layout.origin.X) / layout.size.X, Y: (p.Y - layout.origin.Y) / layout.size.Y}
	q := M.b0*pt.X + M.b1*pt.Y
	r := M.b2*pt.X + M.b3*pt.Y
	return FractionalHex{q: q, r: r, s: -q - r}
}

// 2.3 Drawing a hex

func (h Hex) CornerOffset(layout Layout, corner int) Point {
	return HexCornerOffset(layout, corner)
}

func (layout Layout) HexCornerOffset(corner int) Point {
	return HexCornerOffset(layout, corner)
}

func (h Hex) PolygonCorners(layout Layout) [6]Point {
	return PolygonCorners(layout, h)
}

func (layout Layout) PolygonCorners() [6]Point {
	return PolygonCorners(layout, Hex{})
}

// HexCornerOffset returns the screen location (pixel) of a corner of a hex on the grid.
func HexCornerOffset(layout Layout, corner int) Point {
	size := layout.size
	// todo: is adding corner to start_angle correct?
	angle := 2.0 * math.Pi * (layout.orientation.start_angle + float64(corner)) / 6
	return Point{X: size.X * math.Cos(angle), Y: size.Y * math.Sin(angle)}
}

// PolygonCorners returns the location of the six corners of the hex on the grid.
func PolygonCorners(layout Layout, h Hex) [6]Point {
	var corners [6]Point
	center := layout.HexToPixel(h)
	for i := 0; i < 6; i++ {
		offset := layout.HexCornerOffset(i)
		corners[i] = Point{X: center.X + offset.X, Y: center.Y + offset.Y}
	}
	return corners
}

// 2.4 Layout examples

// todo: an alternate implementation uses origin(0,0) and size(1,0) with chained transformations

// hex→pixel: hex→cartesian, then scale the cartesian coordinate by multiplying by the desired scale, and then translate it to the desired origin.
// pixel→hex: undo the translate by subtracting the origin, then undo the scale by dividing by the scale, then run cartesian→hex.

// 3.0 Fractional Hex

// NewFractionalHex returns a FractionalHex initialized with Cube coordinates.
// Panics on invalid input.
func NewFractionalHex(q, r, s float64) FractionalHex {
	if q+r+s != 0 {
		panic("assert (q + r + s == 0)")
	}
	return FractionalHex{q: q, r: r, s: s}
}

// NewFractionalHexFromAxialCoords returns a FractionalHex initialized with Axial coordinates.
func NewFractionalHexFromAxialCoords(q, r float64) FractionalHex {
	return FractionalHex{q: q, r: r, s: -q - r}
}

// 3.1 Hex rounding

// Round turns a fractional hex coordinate into the nearest integer hex coordinate.
func (h FractionalHex) Round() Hex {
	q := int(math.Round(h.q))
	r := int(math.Round(h.r))
	s := int(math.Round(h.s))
	q_diff := abs(float64(q) - h.q)
	r_diff := abs(float64(r) - h.r)
	s_diff := abs(float64(s) - h.s)
	if q_diff > r_diff && q_diff > s_diff {
		q = -r - s
	} else if r_diff > s_diff {
		r = -q - s
	} else {
		s = -q - r
	}
	return Hex{q: q, r: r, s: s}
}

// ToHexRounded turns a fractional hex into a regular hex coordinate:
func (p Point) ToHexRounded(layout Layout) Hex {
	return layout.PixelToFractionalHex(p).Round()
}

// PixelToHexRounded turns a fractional hex into a regular hex coordinate:
func (layout Layout) PixelToHexRounded(p Point) Hex {
	return layout.PixelToFractionalHex(p).Round()
}

// PixelToHexRounded turns a fractional hex into a regular hex coordinate:
func PixelToHexRounded(layout Layout, p Point) Hex {
	return layout.PixelToFractionalHex(p).Round()
}

// 3.2 Line drawing

// lerp returns the linear interpolation of points on the line between two hexes
func lerp(a, b, t float64) float64 {
	return a*(1-t) + b*t // better for floating point precision than a + (b - a) * t
}

// Lerp returns the linear interpolation of points on the line between two hexes
func (h Hex) Lerp(b Hex, t float64) FractionalHex {
	return FractionalHex{
		q: lerp(float64(h.q), float64(b.q), t),
		r: lerp(float64(h.r), float64(b.r), t),
		s: lerp(float64(h.s), float64(b.s), t),
	}
}

// Lerp returns the linear interpolation of points on the line between two hexes
func (h FractionalHex) Lerp(b FractionalHex, t float64) FractionalHex {
	return FractionalHex{
		q: lerp((h.q), (b.q), t),
		r: lerp((h.r), (b.r), t),
		s: lerp((h.s), (b.s), t),
	}
}

// Linedraw returns the set of hexes that are between two hexes.
// Enable nudging to push points on an edge in a consistent direction.
func (h Hex) Linedraw(b Hex, withNudge bool) []Hex {
	N := h.Distance(b)
	var results []Hex
	var step float64
	if N == 0 {
		step = 1.0
	} else {
		step = 1.0 / float64(N)
	}
	if withNudge {
		h_nudge := FractionalHex{q: float64(h.q) + 1e-6, r: float64(h.r) + 1e-6, s: float64(h.s) - 2e-6}
		b_nudge := FractionalHex{q: float64(b.q) + 1e-6, r: float64(b.r) + 1e-6, s: float64(b.s) - 2e-6}
		for i := 0; i <= N; i++ {
			results = append(results, h_nudge.Lerp(b_nudge, step*float64(i)).Round())
		}
		return results
	}
	for i := 0; i <= N; i++ {
		results = append(results, h.Lerp(b, step*float64(i)).Round())
	}
	return results
}

// 4.0 Map storage

// 4.1 Map storage

// Key returns a hashable uint64 value derived from the axial
// coordinates (q, r) using a variation of Boost's hash_combine
// and MurmurHash3 finalization constants.
//
// The s coordinate is omitted because it is redundant in cube
// coordinates (s = -q - r).
//
// Casting to int64 before converting to uint64 preserves the signed
// bit pattern of negative values, maintaining good distribution
// across the entire hex grid, including negative coordinates.
func Key(q, r int) uint64 {
	const c1 = 0x9E3779B97F4A7C15 // golden ratio
	const c2 = 0xBF58476D1CE4E5B9
	const c3 = 0x94D049BB133111EB

	q64 := uint64(int64(q))
	r64 := uint64(int64(r))

	z := q64 ^ (r64 + c1 + (q64 << 6) + (q64 >> 2))
	z = (z ^ (z >> 30)) * c2
	z = (z ^ (z >> 27)) * c3
	return z ^ (z >> 31)
}

// Hash returns a hashable value for a Hex, derived from the
// q and r coordinates. It delegates to Key(q, r) to avoid
// duplication.
func (h Hex) Hash() uint64 {
	return Key(h.q, h.r)
}

// example using hex_hash to create a map of floats keyed by hex
// var heights map[uint64]float64
// heights[new_hex(1, -2, 3).Hash()] = 4.3

// 4.2 Map shapes

// 4.2.1 Parallelograms

// ParallelogramGrid returns a grid originating at (0,0,0).
// the internal logic depends on the orientation of the grid.
// I don't understand the comment in the source about there
// being three coordinates and the caller has to choose two.
func (layout Layout) ParallelogramGrid(q1, r1, q2, r2 int) GridStore {
	gs := GridStore{}
	if layout.IsPointyTop() {
		for q := q1; q <= q2; q++ {
			for r := r1; r <= r2; r++ {
				hex := NewHexFromAxialCoords(q, r)
				gs[hex.Hash()] = hex
			}
		}
		return gs
	}
	for q := q1; q <= q2; q++ {
		for r := r1; r <= r2; r++ {
			hex := NewHexFromAxialCoords(q, r)
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

// 4.2.2 Triangles

// TriagonalGrid returns a grid originating at (0,0,0).
// the internal logic depends on the orientation of the grid.
// there's a comment in the source about flipping the y-axis to
// change the direction of the triangle, but I don't understand
// how to implement that.
func (layout Layout) TriagonalGrid(side_length int) GridStore {
	gs := GridStore{}
	map_size := side_length
	if layout.IsPointyTop() {
		for q := 0; q <= map_size; q++ {
			for r := 0; r <= map_size-q; r++ {
				hex := NewHexFromAxialCoords(q, r)
				gs[hex.Hash()] = hex
			}
		}
		return gs
	}
	// flat top
	for q := 0; q <= map_size; q++ {
		for r := map_size - q; r <= map_size; r++ {
			hex := NewHexFromAxialCoords(q, r)
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

// 4.2.3 Hexagons

// HexagonalGrid returns a grid centered about (0,0,0).
// does not depend on the orientation of the grid.
func HexagonalGrid(radius int) GridStore {
	gs := GridStore{}
	N := radius
	for q := -N; q <= N; q++ {
		r1 := max(-N, -q-N)
		r2 := min(N, -q+N)
		for r := r1; r <= r2; r++ {
			hex := NewHexFromAxialCoords(q, r)
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

// 4.2.4 Rectangles

// RectangularGrid returns a grid centered about (0,0,0).
// the internal logic depends on the orientation of the grid.
func (layout Layout) RectangularGrid(left, right, top, bottom int) GridStore {
	gs := GridStore{}
	if layout.IsPointyTop() {
		for r := top; r <= bottom; r++ {
			r_offset := r >> 1 // or math.Floor(float64(r) / 2.0)
			for q := left - r_offset; q <= right-r_offset; q++ {
				hex := NewHexFromAxialCoords(q, r)
				gs[hex.Hash()] = hex
			}
		}
		return gs
	}
	// flat top
	for q := left; q <= right; q++ {
		q_offset := q >> 1 // or math.Floor(float64(q) / 2.0)
		for r := top - q_offset; r <= bottom-q_offset; r++ {
			hex := NewHexFromAxialCoords(q, r)
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

// 4.3 Optimized storage

// todo: translate the template RectangularPointTopMap

// 5.0 Rotation

func (h Hex) RotateLeft() Hex {
	return Hex{q: -h.s, r: -h.q, s: -h.r}
}

func (h Hex) RotateRight() Hex {
	return Hex{q: -h.r, r: -h.s, s: -h.q}
}

// 6.0 Offset coordinates

// From the source:
// For offset coordinates I need to know if a row/col is even or odd.
// I use `a&1` (bitwise and) instead of `a%2` return 0 or +1. Why?
//
// * On systems using two’s complement representation, which is just
//   about every system out there, `a&1` returns 0 for even a and 1 for
//   odd `a`. This is what I want. It’s not strictly portable, but should
//   work everywhere in practice.
// * The % remainder operator has multiple variants: floored, Euclidean,
//   truncated, rounded, and ceiling.
//   * With floored or Euclidean, (-1) % 2 is +1
//   * With truncated, (-1) % 2 is -1. This will cause the algorithms
//     on this page to break for negative coordinates.

// NewOffsetCoord returns a new OffsetCord.
func NewOffsetCoord(col, row int) OffsetCoord {
	return OffsetCoord{col: col, row: row}
}

// there are four types of OffsetCoord
const (
	odd_r int = iota
	even_r
	odd_q
	even_q

	EVEN = +1
	ODD  = -1
)

// panics on invalid input
func (h Hex) qoffset_from_cube(offset int) OffsetCoord {
	if !(offset == EVEN || offset == ODD) {
		panic("assert(offset == EVEN || offset == ODD)")
	}
	col := h.q
	row := h.r + int((h.q+offset*(h.q&1))/2)
	return OffsetCoord{col: col, row: row}
}

func (h Hex) qoffset_from_cube_even() OffsetCoord {
	col := h.q
	row := h.r + int((h.q+EVEN*(h.q&1))/2)
	return OffsetCoord{col: col, row: row}
}

func (h Hex) qoffset_from_cube_odd() OffsetCoord {
	col := h.q
	row := h.r + int((h.q+ODD*(h.q&1))/2)
	return OffsetCoord{col: col, row: row}
}

// panics on invalid input
func (oc OffsetCoord) qoffset_to_cube(offset int) Hex {
	if !(offset == EVEN || offset == ODD) {
		panic("assert(offset == EVEN || offset == ODD)")
	}
	q := oc.col
	r := oc.row - int((oc.col+offset*(oc.col&1))/2)
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

func (oc OffsetCoord) qoffset_to_cube_even() Hex {
	q := oc.col
	r := oc.row - int((oc.col+EVEN*(oc.col&1))/2)
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

func (oc OffsetCoord) qoffset_to_cube_odd() Hex {
	q := oc.col
	r := oc.row - int((oc.col+ODD*(oc.col&1))/2)
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

func (oc OffsetCoord) ToCubeOdd() Hex {
	q := oc.col
	r := oc.row - int((oc.col+ODD*(oc.col&1))/2)
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

// panics on invalid input
func (h Hex) roffset_from_cube(offset int) OffsetCoord {
	if !(offset == EVEN || offset == ODD) {
		panic("assert(offset == EVEN || offset == ODD)")
	}
	col := h.q + int((h.r+offset*(h.r&1))/2)
	row := h.r
	return OffsetCoord{col: col, row: row}
}

func (h Hex) roffset_from_cube_even() OffsetCoord {
	col := h.q + int((h.r+EVEN*(h.r&1))/2)
	row := h.r
	return OffsetCoord{col: col, row: row}
}

func (h Hex) roffset_from_cube_odd() OffsetCoord {
	col := h.q + int((h.r+ODD*(h.r&1))/2)
	row := h.r
	return OffsetCoord{col: col, row: row}
}

// panics on invalid input
func (oc OffsetCoord) roffset_to_cube(offset int) Hex {
	if !(offset == EVEN || offset == ODD) {
		panic("assert(offset == EVEN || offset == ODD)")
	}
	q := oc.col - int((oc.row+offset*(oc.row&1))/2)
	r := oc.row
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

func (oc OffsetCoord) roffset_to_cube_even() Hex {
	q := oc.col - int((oc.row+EVEN*(oc.row&1))/2)
	r := oc.row
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

func (oc OffsetCoord) roffset_to_cube_odd() Hex {
	q := oc.col - int((oc.row+ODD*(oc.row&1))/2)
	r := oc.row
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

// 7.0 Notes

// From the source:
// * In languages that don’t support `a>>1`, you can use `floor(a/2)` instead.
// * Most of the functions are small and should be inlined in languages that support it.

// 7.1 Cube vs Axial

// 7.2 C++

// 7.3 Python, Javascript

// 8.0 Source Code

// 8.1 Code from this page

// 8.2 Other libraries
// Go - github.com/pmcxs/hexgrid
// Go - github.com/hautenessa/hexagolang

// 9.0 Other Systems

// 9.1 TribeNet Coordinates

// TribeNet coordinates are in the form "AB 0102":
// - "A" (grid row) and "B" (grid column) identify a sub-map.
// - "0102" is the in-submap position: column 01 (1-based) and row 02 (1-based).
// Each sub-map is 30 columns wide and 21 rows tall,
// with "0101" as the upper-left and "3021" as the lower-right corner.
//
// The global map origin is (1,1) at the upper-left.
// Even-numbered columns are vertically offset (odd-q layout), so (2,1) is southeast of (1,1).
//
// TribeNet coordinates are converted to OffsetCoord using "odd-q" layout,
// with the origin translated by (-1, -1) so "AA 0101" becomes (0,0).

const (
	tnRowsPerGrid  = 21
	tnColsPerGrid  = 30
	tnMaxGridIndex = 26 // A ... Z -> 1 ... 26
)

// NewTribeNetCoord parses a Tribenet coordinate (eg, "AB 0102") and returns
// an OffsetCoord. All TribeNet maps are "odd-q,", origin is (1,1) and will
// be translated to (0,0), sub-maps are 21 rows x 30 columns.
//
// Invalid inputs will return an error.
//
// TribeNet coordinate "AA 0101" corresponds to OffsetCoord (  0,  0).
// TribeNet coordinate "BC 0824" corresponds to OffsetCoord ( 61, 43).
// TribeNet coordinate "JK 0609" corresponds to OffsetCoord (163,104).
// TribeNet coordinate "ZZ 3021" corresponds to OffsetCoord (779,545).

func NewTribeNetCoord(input string) (OffsetCoord, error) {
	if len(input) != 7 || input[2] != ' ' {
		return OffsetCoord{}, fmt.Errorf("invalid format: expected 'AB 0102'")
	}
	gridRowRune, gridColRune := rune(input[0]), rune(input[1])
	subColStr, subRowStr := input[3:5], input[5:7]

	// convert letters to 0-based grid index (A=0, B=1, ...)
	if !(unicode.IsUpper(gridRowRune) && unicode.IsLetter(gridRowRune)) {
		return OffsetCoord{}, fmt.Errorf("invalid grid row: must be uppercase A-Z")
	} else if !(unicode.IsUpper(gridColRune) && unicode.IsLetter(gridColRune)) {
		return OffsetCoord{}, fmt.Errorf("invalid grid column: must be uppercase A-Z")
	}
	gridRowOffset := int(gridRowRune-'A') * tnRowsPerGrid
	gridColOffset := int(gridColRune-'A') * tnColsPerGrid

	// subCol is 1-based in the input; converted to 0-based for OffsetCoord
	subCol, err := strconv.Atoi(subColStr)
	if !(err == nil && 1 <= subCol && subCol <= tnColsPerGrid) {
		return OffsetCoord{}, fmt.Errorf("invalid sub-map column: %s", subColStr)
	}
	subCol -= 1

	// subRow is 1-based in the input; converted to 0-based for OffsetCoord
	subRow, err := strconv.Atoi(subRowStr)
	if !(err == nil && 1 <= subRow && subRow <= tnRowsPerGrid) {
		return OffsetCoord{}, fmt.Errorf("invalid sub-map row: %s", subRowStr)
	}
	subRow -= 1

	// convert global row and column to an offset coordinate
	return OffsetCoord{
		col: gridColOffset + subCol,
		row: gridRowOffset + subRow,
	}, nil
}

// ToTribeNetCoord converts an OffsetCoord to a TribeNet coordinate string ("AB 0102").
//
// The conversion is based on:
// - A 0-based global offset grid with "odd-q" layout and origin (0,0) = "AA 0101".
// - Each 30×21 cell in the global grid is a sub-map labeled by row and column letters (A–Z).
// - The function converts the offset col/row to sub-map ID and in-map sub-coordinates.
//
// Returns an error if the coordinate falls outside the supported 26×26 letter grid.
func (oc OffsetCoord) ToTribeNetCoord() (string, error) {
	if oc.col < 0 || oc.row < 0 {
		return "", fmt.Errorf("invalid offset coordinates: %v", oc)
	}

	gridRow, gridCol := oc.row/tnRowsPerGrid, oc.col/tnColsPerGrid
	if gridRow >= tnMaxGridIndex || gridCol >= tnMaxGridIndex {
		return "", fmt.Errorf("coordinates out of range for A-Z grid system")
	}
	gridRowChar, gridColChar := 'A'+rune(gridRow), 'A'+rune(gridCol)
	subCol, subRow := (oc.col%tnColsPerGrid)+1, (oc.row%tnRowsPerGrid)+1

	return fmt.Sprintf("%c%c %02d%02d", gridRowChar, gridColChar, subCol, subRow), nil
}
