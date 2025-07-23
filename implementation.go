// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package hexg implements hex grids in the style of the Red Blob Games blog
// https://www.redblobgames.com/grids/hexagons/implementation.html
package hexg

import "math"

// 1.0 Hex coordinates

// Hex implements a hex that uses Cube coordinates.
type Hex struct {
	q, r, s int
}

// new_hex returns a Hex initialized with Cube coordinates.
func new_hex(q, r, s int) Hex {
	if q+r+s != 0 {
		panic("assert (q + r + s == 0)")
	}
	return Hex{q: q, r: r, s: s}
}

// new_hex_from_axial_coords returns a Hex initialized with Axial coordinates.
func new_hex_from_axial_coords(q, r int) Hex {
	return new_hex(q, r, -q-r)
}

// 1.1 Equality

// hex_equals returns true if the two hexes have the same coordinates.
func hex_equals(a, b Hex) bool {
	return a.q == b.q && a.r == b.r && a.s == b.s
}

// hex_not_equals returns true if the two hexes have different coordinates.
func hex_not_equals(a, b Hex) bool {
	return !hex_equals(a, b)
}

// 1.2 Coordinate arithmetic

// hex_add returns the sum of two hexes.
func hex_add(a, b Hex) Hex {
	return Hex{q: a.q + b.q, r: a.r + b.r, s: a.s + b.s}
}

// hex_subtract returns the difference of two hexes.
func hex_subtract(a, b Hex) Hex {
	return Hex{q: a.q - b.q, r: a.r - b.r, s: a.s - b.s}
}

// hex_multiply returns a Hex scaled by an integer.
func hex_multiply(a Hex, k int) Hex {
	return Hex{q: a.q * k, r: a.r * k, s: a.s * k}
}

// 1.3 Distance

// hex_length is the line from the origin to a hex
func hex_length(hex Hex) int {
	return int((abs(hex.q) + abs(hex.r) + abs(hex.s)) / 2)
}

// hex_distance between two hexes is the length of the line between them.
func hex_distance(a, b Hex) int {
	return hex_length(hex_subtract(a, b))
}

// 1.3.1 Neighbors

// hex_directions has the offset to the neighboring hex indexed by direction 0..5
var hex_directions = [6]Hex{
	{1, 0, -1}, {1, -1, 0}, {0, -1, 1},
	{-1, 0, 1}, {-1, 1, 0}, {0, 1, -1},
}

func hex_direction(direction int) Hex {
	return hex_directions[(6+(direction%6))%6]
}

// hex_neighbor returns the hex that is one step away in the given direction.
// Direction is coerced to the range 0..5.
func hex_neighbor(hex Hex, direction int) Hex {
	return hex_add(hex, hex_directions[(6+(direction%6))%6])
}

// 2.0 Layout

// orientation is a helper for Layouts
type orientation struct {
	f0, f1, f2, f3 float64
	b0, b1, b2, b3 float64
	start_angle    float64 // in multiples of 60°
}

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

// new_layout creates a new Layout with the orientation, size, and origin
func new_layout(o orientation, size, origin Point, offset int) Layout {
	return Layout{
		orientation: o,
		size:        size,
		origin:      origin,
		offset:      offset,
	}
}

// Point represents a coordinate on the hexagonal grid
type Point struct {
	x float64
	y float64
}

// new_point creates a new Point with specified coordinates
func new_point(x, y float64) Point {
	return Point{x: x, y: y}
}

// 2.1 Hex to screen

// hex_to_pixel returns the origin of the hex on the grid as a Point.
func hex_to_pixel(layout Layout, h Hex) Point {
	M := layout.orientation
	return Point{
		x: layout.origin.x + (M.f0*float64(h.q)+M.f1*float64(h.r))*layout.size.x,
		y: layout.origin.y + (M.f2*float64(h.q)+M.f3*float64(h.r))*layout.size.y,
	}
}

// 2.2 Screen to hex

// pixel_to_hex_fractional returns the fractional hex that encloses the pixel.
// In theory, the origin of that fractional hex will be the pixel.
func pixel_to_hex_fractional(layout Layout, p Point) FractionalHex {
	M := layout.orientation
	pt := Point{x: (p.x - layout.origin.x) / layout.size.x, y: (p.y - layout.origin.y) / layout.size.y}
	q := M.b0*pt.x + M.b1*pt.y
	r := M.b2*pt.x + M.b3*pt.y
	return FractionalHex{q: q, r: r, s: -q - r}
}

// 2.3 Drawing a hex

// hex_corner_offset returns the screen location (pixel) of a corner of a hex on the grid.
func hex_corner_offset(layout Layout, corner int) Point {
	size := layout.size
	// todo: is adding corner to start_angle correct?
	angle := 2.0 * math.Pi * (layout.orientation.start_angle + float64(corner)) / 6
	return Point{x: size.x * math.Cos(angle), y: size.y * math.Sin(angle)}
}

// polygon_corners returns the location of the six corners of the hex on the grid.
func polygon_corners(layout Layout, h Hex) [6]Point {
	var corners [6]Point
	center := hex_to_pixel(layout, h)
	for i := 0; i < 6; i++ {
		offset := hex_corner_offset(layout, i)
		corners[i] = Point{x: center.x + offset.x, y: center.y + offset.y}
	}
	return corners
}

// 2.4 Layout examples

// todo: an alternate implementation uses origin(0,0) and size(1,0) with chained transformations

// hex→pixel: hex→cartesian, then scale the cartesian coordinate by multiplying by the desired scale, and then translate it to the desired origin.
// pixel→hex: undo the translate by subtracting the origin, then undo the scale by dividing by the scale, then run cartesian→hex.

// 3.0 Fractional Hex

// FractionalHex implements a hex with floating point coordinates.
// Used for linear interpolation and rounding.
type FractionalHex struct {
	q, r, s float64
}

// new_fractional_hex returns a FractionalHex initialized with Cube coordinates.
func new_fractional_hex(q, r, s float64) FractionalHex {
	if q+r+s != 0 {
		panic("assert (q + r + s == 0)")
	}
	return FractionalHex{q: q, r: r, s: s}
}

// new_fractional_hex_from_axial_coords returns a FractionalHex initialized with Axial coordinates.
func new_fractional_hex_from_axial_coords(q, r float64) FractionalHex {
	return new_fractional_hex(q, r, -q-r)
}

// 3.1 Hex rounding

// hex_round turns a fractional hex coordinate into the nearest integer hex coordinate.
func hex_round(h FractionalHex) Hex {
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

// pixel_to_hex_rounded turns a fractional hex into a regular hex coordinate:
func pixel_to_hex_rounded(layout Layout, p Point) Hex {
	return hex_round(pixel_to_hex_fractional(layout, p))
}

// 3.2 Line drawing

// lerp returns the linear interpolation of points on the line between two hexes
func lerp(a, b, t float64) float64 {
	return a*(1-t) + b*t // better for floating point precision than a + (b - a) * t
}

func hex_lerp(a, b Hex, t float64) FractionalHex {
	return FractionalHex{
		q: lerp(float64(a.q), float64(b.q), t),
		r: lerp(float64(a.r), float64(b.r), t),
		s: lerp(float64(a.s), float64(b.s), t),
	}
}

// hex_linedraw returns the set of hexes that are between two hexes.
func hex_linedraw(a, b Hex) []Hex {
	N := hex_distance(a, b)
	var results []Hex
	var step float64
	if N == 0 {
		step = 1.0
	} else {
		step = 1.0 / float64(N)
	}
	for i := 0; i <= N; i++ {
		results = append(results, hex_round(hex_lerp(a, b, step*float64(i))))
	}
	return results
}

// hex_linedraw_with_nudge returns the set of hexes that are between two hexes.
// enable nudging to push points on an edge in a consistent direction.
func hex_linedraw_with_nudge(a, b Hex) []Hex {
	N := hex_distance(a, b)
	a_nudge := FractionalHex{q: float64(a.q) + 1e-6, r: float64(a.r) + 1e-6, s: float64(a.s) - 2e-6}
	b_nudge := FractionalHex{q: float64(b.q) + 1e-6, r: float64(b.r) + 1e-6, s: float64(b.s) - 2e-6}
	var results []Hex
	var step float64
	if N == 0 {
		step = 1.0
	} else {
		step = 1.0 / float64(N)
	}
	for i := 0; i <= N; i++ {
		results = append(results, hex_round(fractional_hex_lerp(a_nudge, b_nudge, step*float64(i))))
	}
	return results
}

func fractional_hex_lerp(a, b FractionalHex, t float64) FractionalHex {
	return FractionalHex{
		q: lerp((a.q), (b.q), t),
		r: lerp((a.r), (b.r), t),
		s: lerp((a.s), (b.s), t),
	}
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

// GridStore is a map of Hex indexed by the hash of the Hex
type GridStore map[uint64]Hex

// 4.2.1 Parallelograms

// parallelogram_grid returns a grid originating at (0,0,0).
// the internal logic depends on the orientation of the grid.
// I don't understand the comment in the source about there
// being three coordinates and the caller has to choose two.
func parallelogram_grid(layout Layout, q1, r1, q2, r2 int) GridStore {
	gs := GridStore{}
	if layout.IsPointyTop() {
		for q := q1; q <= q2; q++ {
			for r := r1; r <= r2; r++ {
				hex := new_hex(q, r, -q-r)
				gs[hex.Hash()] = hex
			}
		}
		return gs
	}
	for q := q1; q <= q2; q++ {
		for r := r1; r <= r2; r++ {
			hex := new_hex(q, r, -q-r)
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

// 4.2.2 Triangles

// triagonal_grid returns a grid originating at (0,0,0).
// the internal logic depends on the orientation of the grid.
// there's a comment in the source about flipping the y-axis to
// change the direction of the triangle, but I don't understand
// how to implement that.
func triagonal_grid(layout Layout, side_length int) GridStore {
	gs := GridStore{}
	map_size := side_length
	if layout.IsPointyTop() {
		for q := 0; q <= map_size; q++ {
			for r := 0; r <= map_size-q; r++ {
				hex := new_hex(q, r, -q-r)
				gs[hex.Hash()] = hex
			}
		}
		return gs
	}
	// flat top
	for q := 0; q <= map_size; q++ {
		for r := map_size - q; r <= map_size; r++ {
			hex := new_hex(q, r, -q-r)
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

// 4.2.3 Hexagons

// hexagonal_grid returns a grid centered about (0,0,0).
// does not depend on the orientation of the grid.
func hexagonal_grid(radius int) GridStore {
	gs := GridStore{}
	N := radius
	for q := -N; q <= N; q++ {
		r1 := max(-N, -q-N)
		r2 := min(N, -q+N)
		for r := r1; r <= r2; r++ {
			hex := new_hex(q, r, -q-r)
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

// 4.2.4 Rectangles

// rectangular_grid returns a grid centered about (0,0,0).
// the internal logic depends on the orientation of the grid.
func rectangular_grid(layout Layout, left, right, top, bottom int) GridStore {
	gs := GridStore{}
	if layout.IsPointyTop() {
		for r := top; r <= bottom; r++ {
			r_offset := r >> 1 // or math.Floor(float64(r) / 2.0)
			for q := left - r_offset; q <= right-r_offset; q++ {
				hex := new_hex(q, r, -q-r)
				gs[hex.Hash()] = hex
			}
		}
		return gs
	}
	// flat top
	for q := left; q <= right; q++ {
		q_offset := q >> 1 // or math.Floor(float64(q) / 2.0)
		for r := top - q_offset; r <= bottom-q_offset; r++ {
			hex := new_hex(q, r, -q-r)
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

// 4.3 Optimized storage

// todo: translate the template RectangularPointTopMap

// 5.0 Rotation

func hex_rotate_left(a Hex) Hex {
	return Hex{q: -a.s, r: -a.q, s: -a.r}
}

func hex_rotate_right(a Hex) Hex {
	return Hex{q: -a.r, r: -a.s, s: -a.q}
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
// * The % remainder operator has multiple variants: floored, euclidean,
//   truncated, rounded, and ceiling.
//   * With floored or euclidean, (-1) % 2 is +1
//   * With truncated, (-1) % 2 is -1. This will cause the algorithms
//     on this page to break for negative coordinates.

type OffsetCoord struct {
	col, row int
}

func new_offset_coord(col, row int) OffsetCoord {
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

func qoffset_from_cube(offset int, h Hex) OffsetCoord {
	if !(offset == EVEN || offset == ODD) {
		panic("assert(offset == EVEN || offset == ODD)")
	}
	col := h.q
	row := h.r + int((h.q+offset*(h.q&1))/2)
	return OffsetCoord{col: col, row: row}
}

func qoffset_from_cube_even(h Hex) OffsetCoord {
	col := h.q
	row := h.r + int((h.q+EVEN*(h.q&1))/2)
	return OffsetCoord{col: col, row: row}
}

func qoffset_from_cube_odd(h Hex) OffsetCoord {
	col := h.q
	row := h.r + int((h.q+ODD*(h.q&1))/2)
	return OffsetCoord{col: col, row: row}
}

func qoffset_to_cube(offset int, h OffsetCoord) Hex {
	if !(offset == EVEN || offset == ODD) {
		panic("assert(offset == EVEN || offset == ODD)")
	}
	q := h.col
	r := h.row - int((h.col+offset*(h.col&1))/2)
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

func qoffset_to_cube_even(h OffsetCoord) Hex {
	q := h.col
	r := h.row - int((h.col+EVEN*(h.col&1))/2)
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

func qoffset_to_cube_odd(h OffsetCoord) Hex {
	q := h.col
	r := h.row - int((h.col+ODD*(h.col&1))/2)
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

func roffset_from_cube(offset int, h Hex) OffsetCoord {
	if !(offset == EVEN || offset == ODD) {
		panic("assert(offset == EVEN || offset == ODD)")
	}
	col := h.q + int((h.r+offset*(h.r&1))/2)
	row := h.r
	return OffsetCoord{col: col, row: row}
}

func roffset_from_cube_even(h Hex) OffsetCoord {
	col := h.q + int((h.r+EVEN*(h.r&1))/2)
	row := h.r
	return OffsetCoord{col: col, row: row}
}

func roffset_from_cube_odd(h Hex) OffsetCoord {
	col := h.q + int((h.r+ODD*(h.r&1))/2)
	row := h.r
	return OffsetCoord{col: col, row: row}
}

func roffset_to_cube(offset int, h OffsetCoord) Hex {
	if !(offset == EVEN || offset == ODD) {
		panic("assert(offset == EVEN || offset == ODD)")
	}
	q := h.col - int((h.row+offset*(h.row&1))/2)
	r := h.row
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

func roffset_to_cube_even(h OffsetCoord) Hex {
	q := h.col - int((h.row+EVEN*(h.row&1))/2)
	r := h.row
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

func roffset_to_cube_odd(h OffsetCoord) Hex {
	q := h.col - int((h.row+ODD*(h.row&1))/2)
	r := h.row
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
