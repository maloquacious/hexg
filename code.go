// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

//import (
//	"fmt"
//	"math"
//)
//
//// LayoutPointyTypeHorizontalOddRight_t represents the orientation of a hexagonal grid.
//// Layouts default to point-top, horizontal, odd rows pushed right.
//type LayoutPointyTypeHorizontalOddRight_t struct {
//	orientation orientation
//
//	// size is used for scaling (eg, matching pixel sprite sizes)
//	// set it to Point(size, size) if you need uniform scaling.
//	// todo: this makes no sense. i need to read the page again.
//	size Point
//
//	// origin is the center of the q=0,r=0 hexagon.
//	// set it to Point(0, 0) if you do not need to translate the transformation.
//	origin Point
//
//	// offset is used for OffsetCoord neighbors
//	offset layoutOffset
//}
//
//// HashTable is a map of Hex indexed by the hash of the Hex.
//// I think that it's here for examples of how to use hex.Hash.
//type HashTable map[uint64]Hex
//
//// ConciseString returns the coordinates with signs.
//// It returns the coordinates formatted as (+q+r+s).
//func (h Hex) ConciseString() string {
//	return fmt.Sprintf("%+d%+d%+d", h.q, h.r, h.s)
//}
//
//// String implements the Stringer interface.
//// It returns the coordinates formatted as (q,r,s).
//func (h Hex) String() string {
//	return fmt.Sprintf("%d,%d,%d", h.q, h.r, h.s)
//}
//
//// NewLayoutFlat returns a layout with flat-top hexes, vertical layout, odd columns pushed down.
//// Size and origin are used when calculating screen pixels.
//func NewLayoutFlat(size, origin Point) LayoutPointyTypeHorizontalOddRight_t {
//	return LayoutPointyTypeHorizontalOddRight_t{
//		orientation: layout_flat,
//		size:        size,
//		origin:      origin,
//		offset:      odd_q,
//	}
//}
//
//// NewLayoutPointy returns a layout with pointy-top hexes, horizontal layout, odd rows pushed right.
//// Size and origin are used when calculating screen pixels.
//func NewLayoutPointy(size, origin Point, shoveOddRowsRight bool) LayoutPointyTypeHorizontalOddRight_t {
//	return LayoutPointyTypeHorizontalOddRight_t{
//		orientation: layout_pointy,
//		size:        size,
//		origin:      origin,
//		offset:      odd_r,
//	}
//}
//
//// NewLayoutEvenQ returns a layout with flat-top hexes, vertical layout, shoves even columns down.
//// Size and origin are used when calculating screen pixels.
//func NewLayoutEvenQ(size, origin Point) LayoutPointyTypeHorizontalOddRight_t {
//	return LayoutPointyTypeHorizontalOddRight_t{
//		orientation: layout_flat,
//		size:        size,
//		origin:      origin,
//		offset:      even_q,
//	}
//}
//
//// NewLayoutEvenR returns a layout with pointy-top hexes, horizontal layout, shoves even rows right.
//// Size and origin are used when calculating screen pixels.
//func NewLayoutEvenR(size, origin Point) LayoutPointyTypeHorizontalOddRight_t {
//	return LayoutPointyTypeHorizontalOddRight_t{
//		orientation: layout_pointy,
//		size:        size,
//		origin:      origin,
//		offset:      even_r,
//	}
//}
//
//// NewLayoutOddQ returns a layout with flat-top hexes, vertical layout, shoves odd columns down.
//// Size and origin are used when calculating screen pixels.
//func NewLayoutOddQ(size, origin Point) LayoutPointyTypeHorizontalOddRight_t {
//	return LayoutPointyTypeHorizontalOddRight_t{
//		orientation: layout_flat,
//		size:        size,
//		origin:      origin,
//		offset:      odd_q,
//	}
//}
//
//// NewLayoutOddR returns a layout with pointy-top hexes, horizontal layout, shoves odd rows right.
//// Size and origin are used when calculating screen pixels.
//func NewLayoutOddR(size, origin Point) LayoutPointyTypeHorizontalOddRight_t {
//	return LayoutPointyTypeHorizontalOddRight_t{
//		orientation: layout_pointy,
//		size:        size,
//		origin:      origin,
//		offset:      odd_r,
//	}
//}
//
//// IsFlatTop returns true if the layout was created with flat-top hexes.
//func (layout LayoutPointyTypeHorizontalOddRight_t) IsFlatTop() bool {
//	return layout.offset == odd_q || layout.offset == even_q
//}
//
//// IsPointyTop returns true if the layout was created with point-top hexes.
//func (layout LayoutPointyTypeHorizontalOddRight_t) IsPointyTop() bool {
//	return layout.offset == odd_r || layout.offset == even_r
//}
//
//// IsEvenQ returns true if the layout supports even-q offset coordinates (vertical layout, shoves even columns down).
//func (layout LayoutPointyTypeHorizontalOddRight_t) IsEvenQ() bool {
//	return layout.offset == even_q
//}
//
//// IsEvenR returns true if the layout supports even-r offset coordinates (horizontal layout, shoves even rows right).
//func (layout LayoutPointyTypeHorizontalOddRight_t) IsEvenR() bool {
//	return layout.offset == even_r
//}
//
//// IsOddQ returns true if the layout supports odd-q offset coordinates (vertical layout, shoves odd columns down).
//func (layout LayoutPointyTypeHorizontalOddRight_t) IsOddQ() bool {
//	return layout.offset == odd_q
//}
//
//// IsOddR returns true if the layout supports odd-r offset coordinates (horizontal layout, shoves odd rows right).
//func (layout LayoutPointyTypeHorizontalOddRight_t) IsOddR() bool {
//	return layout.offset == odd_r
//}
//
//// NewPoint returns a new Point with specified screen coordinates
//func NewPoint(x, y float64) Point {
//	return Point{X: x, Y: y}
//}
//
//func (p Point) String() string {
//	return fmt.Sprintf("%g,%g", p.X, p.Y)
//}
//
//// ToPixel returns the origin of the hex on the grid as a Point.
//func (h Hex) ToPixel(layout LayoutPointyTypeHorizontalOddRight_t) Point {
//	return ToPixel(layout, h)
//}
//
//// ToPixel returns the origin of the hex on the grid as a Point.
//func ToPixel(layout LayoutPointyTypeHorizontalOddRight_t, h Hex) Point {
//	M := layout.orientation
//	return Point{
//		X: layout.origin.X + (M.f0*float64(h.q)+M.f1*float64(h.r))*layout.size.X,
//		Y: layout.origin.Y + (M.f2*float64(h.q)+M.f3*float64(h.r))*layout.size.Y,
//	}
//}
//
//// ToFractionalHex returns the fractional hex that encloses the pixel.
//// In theory, the origin of that fractional hex will be the pixel.
//func (p Point) ToFractionalHex(layout LayoutPointyTypeHorizontalOddRight_t) FractionalHex {
//	return PixelToFractionalHex(layout, p)
//}
//
//// PixelToFractionalHex returns the fractional hex that encloses the pixel.
//// In theory, the origin of that fractional hex will be the pixel.
//func (layout LayoutPointyTypeHorizontalOddRight_t) PixelToFractionalHex(p Point) FractionalHex {
//	return PixelToFractionalHex(layout, p)
//}
//
//func (h Hex) CornerOffset(layout LayoutPointyTypeHorizontalOddRight_t, corner int) Point {
//	return HexCornerOffset(layout, corner)
//}
//
//func (layout LayoutPointyTypeHorizontalOddRight_t) HexCornerOffset(corner int) Point {
//	return HexCornerOffset(layout, corner)
//}
//
//func (h Hex) PolygonCorners(layout LayoutPointyTypeHorizontalOddRight_t) [6]Point {
//	return PolygonCorners(layout, h)
//}
//
//func (layout LayoutPointyTypeHorizontalOddRight_t) PolygonCorners() [6]Point {
//	return PolygonCorners(layout, Hex{})
//}
//
//// 3.0 Fractional Hex
//
//// NewFractionalHex returns a FractionalHex initialized with Cube coordinates.
//// Panics on invalid input.
//func NewFractionalHex(q, r, s float64) FractionalHex {
//	if q+r+s != 0 {
//		panic("assert (q + r + s == 0)")
//	}
//	return FractionalHex{q: q, r: r, s: s}
//}
//
//// NewFractionalHexFromAxialCoords returns a FractionalHex initialized with Axial coordinates.
//func NewFractionalHexFromAxialCoords(q, r float64) FractionalHex {
//	return FractionalHex{q: q, r: r, s: -q - r}
//}
//

//// 4.3 Optimized storage
//
//// todo: translate the template RectangularPointTopMap
//
//// 5.0 Rotation
//
//func (h Hex) RotateLeft() Hex {
//	return Hex{q: -h.s, r: -h.q, s: -h.r}
//}
//
//func (h Hex) RotateRight() Hex {
//	return Hex{q: -h.r, r: -h.s, s: -h.q}
//}
//
//// 6.0 Offset coordinates
//
//// From the source:
//// For offset coordinates I need to know if a row/col is even or odd.
//// I use `a&1` (bitwise and) instead of `a%2` return 0 or +1. Why?
////
//// * On systems using two’s complement representation, which is just
////   about every system out there, `a&1` returns 0 for even a and 1 for
////   odd `a`. This is what I want. It’s not strictly portable, but should
////   work everywhere in practice.
//// * The % remainder operator has multiple variants: floored, Euclidean,
////   truncated, rounded, and ceiling.
////   * With floored or Euclidean, (-1) % 2 is +1
////   * With truncated, (-1) % 2 is -1. This will cause the algorithms
////     on this page to break for negative coordinates.
//
////func OffsetCoordFromColRow(col, row int) OffsetCoord {
////	return OffsetCoord{col: col, row: row}
////}
//
//// HexFromOffsetCoord returns a new Hex from the OffsetCoord.
//func (layout LayoutPointyTypeHorizontalOddRight_t) HexFromOffsetCoord(oc OffsetCoord) Hex {
//	col, row := oc.Col, oc.Row
//	switch layout.offset {
//	case even_q: // flat-top, vertical layout, shoves even columns down
//		q, r := col, row-(col+EVEN*(col&1))/2
//		return Hex{q: q, r: r, s: -q - r}
//	case odd_q: // flat-top, vertical layout, shoves odd columns down
//		q, r := col, row-(col+ODD*(col&1))/2
//		return Hex{q: q, r: r, s: -q - r}
//	case even_r: // pointy-top, horizontal layout, shoves even rows right
//		q, r := col-(row+EVEN*(row&1))/2, row
//		return Hex{q: q, r: r, s: -q - r}
//	case odd_r: // pointy-top, horizontal layout, shoves odd rows right
//		q, r := col-(row+ODD*(row&1))/2, row
//		return Hex{q: q, r: r, s: -q - r}
//	}
//	panic("assert(!reached)")
//}
//
//// HexFromOffsetColRow returns a new Hex using offset column and row coordinates.
//func (layout LayoutPointyTypeHorizontalOddRight_t) HexFromOffsetColRow(col, row int) Hex {
//	switch layout.offset {
//	case even_q: // flat-top, vertical layout, shoves even columns down
//		q, r := col, row-(col+EVEN*(col&1))/2
//		return Hex{q: q, r: r, s: -q - r}
//	case odd_q: // flat-top, vertical layout, shoves odd columns down
//		q, r := col, row-(col+EVEN*(col&1))/2
//		return Hex{q: q, r: r, s: -q - r}
//	case even_r: // pointy-top, horizontal layout, shoves even rows right
//		q, r := col-(row+EVEN*(row&1))/2, row
//		return Hex{q: q, r: r, s: -q - r}
//	case odd_r: // pointy-top, horizontal layout, shoves odd rows right
//		q, r := col-(row+ODD*(row&1))/2, row
//		return Hex{q: q, r: r, s: -q - r}
//	}
//	panic("assert(!reached)")
//}
//
//// there are four types of OffsetCoord
//
//type LayoutOffset_e int
//
//const (
//	OddR LayoutOffset_e = iota
//	EvenR
//	OddQ
//	EvenQ
//)
//
//func (e LayoutOffset_e) String() string {
//	switch e {
//	case OddR:
//		return "odd-r"
//	case EvenR:
//		return "even-r"
//	case OddQ:
//		return "odd-q"
//	case EvenQ:
//		return "even-q"
//	default:
//		panic(fmt.Sprintf("assert(e != %d)", e))
//	}
//}
//
//type layoutOffset int
//
//const (
//	odd_r layoutOffset = iota
//	even_r
//	odd_q
//	even_q
//)
//
//const (
//	EVEN = +1
//	ODD  = -1
//)
//
////// panics on invalid input
////func (h Hex) qoffset_from_cube(offset int) OffsetCoord {
////	if !(offset == EVEN || offset == ODD) {
////		panic("assert(offset == EVEN || offset == ODD)")
////	}
////	col := h.q
////	row := h.r + int((h.q+offset*(h.q&1))/2)
////	return OffsetCoord{col: col, row: row}
////}
////
////func (h Hex) qoffset_from_cube_even() OffsetCoord {
////	col := h.q
////	row := h.r + int((h.q+EVEN*(h.q&1))/2)
////	return OffsetCoord{col: col, row: row}
////}
////
////func (h Hex) qoffset_from_cube_odd() OffsetCoord {
////	col := h.q
////	row := h.r + int((h.q+ODD*(h.q&1))/2)
////	return OffsetCoord{col: col, row: row}
////}
////
////// panics on invalid input
////func (oc OffsetCoord) qoffset_to_cube(offset int) Hex {
////	if !(offset == EVEN || offset == ODD) {
////		panic("assert(offset == EVEN || offset == ODD)")
////	}
////	q := oc.col
////	r := oc.row - int((oc.col+offset*(oc.col&1))/2)
////	s := -q - r
////	return Hex{q: q, r: r, s: s}
////}
////
////func (oc OffsetCoord) qoffset_to_cube_even() Hex {
////	q := oc.col
////	r := oc.row - int((oc.col+EVEN*(oc.col&1))/2)
////	s := -q - r
////	return Hex{q: q, r: r, s: s}
////}
////
////func (oc OffsetCoord) qoffset_to_cube_odd() Hex {
////	q := oc.col
////	r := oc.row - int((oc.col+ODD*(oc.col&1))/2)
////	s := -q - r
////	return Hex{q: q, r: r, s: s}
////}
////
////func (oc OffsetCoord) ToCubeOdd() Hex {
////	q := oc.col
////	r := oc.row - int((oc.col+ODD*(oc.col&1))/2)
////	s := -q - r
////	return Hex{q: q, r: r, s: s}
////}
//
//// HexToOffsetCoord returns the offset coordinates of the hex.
//// Uses the offset from the layout to shift rows and columns correctly.
//func (layout LayoutPointyTypeHorizontalOddRight_t) HexToOffsetCoord(h Hex) OffsetCoord {
//	switch layout.offset {
//	case even_q:
//		col, row := h.q, h.r+(h.q+EVEN*(h.q&1))/2
//		return OffsetCoord{Col: col, Row: row}
//	case odd_q:
//		col, row := h.q, h.r+(h.q+ODD*(h.q&1))/2
//		return OffsetCoord{Col: col, Row: row}
//	case even_r:
//		col, row := h.q+(h.r+EVEN*(h.r&1))/2, h.r
//		return OffsetCoord{Col: col, Row: row}
//	case odd_r:
//		col, row := h.q+(h.r+ODD*(h.r&1))/2, h.r
//		return OffsetCoord{Col: col, Row: row}
//	}
//	panic("assert(!reached)")
//}
//
////// panics on invalid input
////func (h Hex) roffset_from_cube(offset int) OffsetCoord {
////	if !(offset == EVEN || offset == ODD) {
////		panic("assert(offset == EVEN || offset == ODD)")
////	}
////	col := h.q + int((h.r+offset*(h.r&1))/2)
////	row := h.r
////	return OffsetCoord{col: col, row: row}
////}
////
////func (h Hex) roffset_from_cube_even() OffsetCoord {
////	col := h.q + int((h.r+EVEN*(h.r&1))/2)
////	row := h.r
////	return OffsetCoord{col: col, row: row}
////}
////
////func (h Hex) roffset_from_cube_odd() OffsetCoord {
////	col := h.q + int((h.r+ODD*(h.r&1))/2)
////	row := h.r
////	return OffsetCoord{col: col, row: row}
////}
////
////// panics on invalid input
////func (oc OffsetCoord) roffset_to_cube(offset int) Hex {
////	if !(offset == EVEN || offset == ODD) {
////		panic("assert(offset == EVEN || offset == ODD)")
////	}
////	q := oc.col - int((oc.row+offset*(oc.row&1))/2)
////	r := oc.row
////	s := -q - r
////	return Hex{q: q, r: r, s: s}
////}
////
////func (oc OffsetCoord) roffset_to_cube_even() Hex {
////	q := oc.col - int((oc.row+EVEN*(oc.row&1))/2)
////	r := oc.row
////	s := -q - r
////	return Hex{q: q, r: r, s: s}
////}
////
////func (oc OffsetCoord) roffset_to_cube_odd() Hex {
////	q := oc.col - int((oc.row+ODD*(oc.row&1))/2)
////	r := oc.row
////	s := -q - r
////	return Hex{q: q, r: r, s: s}
////}
//
//// 7.0 Notes
//
//// From the source:
//// * In languages that don’t support `a>>1`, you can use `floor(a/2)` instead.
//// * Most of the functions are small and should be inlined in languages that support it.
//
//// 7.1 Cube vs Axial
//
//// 7.2 C++
//
//// 7.3 Python, Javascript
//
//// 8.0 Source Code
//
//// 8.1 Code from this page
//
//// 8.2 Other libraries
//// Go - github.com/pmcxs/hexgrid
//// Go - github.com/hautenessa/hexagolang
//
//// 9.0 Other Systems
//
//// 9.1 Compass Points
//
//// The compass points are approximations, but close to the rose.
//// ✅ = supported, 🚫 = not supported.
////
//// +------+-------+---------+---------+
//// | Name | Angle |  Flat?  | Pointy? |
//// +------+-------+---------+---------+
//// | N    |   0°  |   ✅    |   🚫    |
//// | NNE  |  30°  |   🚫    |   ✅    |
//// | ENE  |  60°  |   ✅    |   🚫    |
//// | E    |  90°  |   🚫    |   ✅    |
//// | ESE  | 120°  |   ✅    |   🚫    |
//// | SSE  | 150°  |   🚫    |   ✅    |
//// | S    | 180°  |   ✅    |   🚫    |
//// | SSW  | 210°  |   🚫    |   ✅    |
//// | WSW  | 240°  |   ✅    |   🚫    |
//// | W    | 270°  |   🚫    |   ✅    |
//// | WNW  | 300°  |   ✅    |   🚫    |
//// | NNW  | 330°  |   🚫    |   ✅    |
//// +------+-------+---------+---------+
////
//// Using an unsupported value for your layout may cause unexpected results.
//const (
//	N   int = 2
//	NNE int = 1
//	ENE int = 1
//	E   int = 0
//	ESE int = 0
//	SSE int = 5
//	S   int = 5
//	SSW int = 4
//	WSW int = 4
//	W   int = 3
//	WNW int = 3
//	NNW int = 2
//)
//
//var (
//	bearingToDirection = map[string]int{
//		"N":   2,
//		"NNE": 1,
//		"ENE": 1,
//		"E":   0,
//		"ESE": 0,
//		"SSE": 5,
//		"S":   5,
//		"SSW": 4,
//		"WSW": 4,
//		"W":   3,
//		"WNW": 3,
//		"NNW": 2,
//	}
//)
//
//// BearingToDirection returns the direction for a compass point.
//// Panics on invalid input.
//func BearingToDirection(p string) int {
//	dir, ok := bearingToDirection[p]
//	if !ok {
//		panic(fmt.Sprintf("assert(p != %q)", p))
//	}
//	return dir
//}
//
//var (
//	// horizontalDirectionToBearing maps a direction to the compass point for a horizontal layout
//	horizontalDirectionToBearing = []string{
//		"E", "NNE", "NNW", "W", "SSW", "SSE",
//	}
//	// verticalDirectionToBearing maps a direction to the compass point for a vertical layout
//	verticalDirectionToBearing = []string{
//		"ESE", "ENE", "N", "WNW", "WSW", "S",
//	}
//)
//
//// 9.2 TribeNet Coordinates
//
//// TribeNet coordinates are in the form "AB 0102":
//// - "A" (grid row) and "B" (grid column) identify a sub-map.
//// - "0102" is the in-submap position: column 01 (1-based) and row 02 (1-based).
//// Each sub-map is 30 columns wide and 21 rows tall,
//// with "0101" as the upper-left and "3021" as the lower-right corner.
////
//// The global map origin is (1,1) at the upper-left.
//// Even-numbered columns are vertically offset (odd-q layout), so (2,1) is southeast of (1,1).
////
//// TribeNet coordinates are converted to OffsetCoord using "odd-q" layout,
//// with the origin translated by (-1, -1) so "AA 0101" becomes (0,0).
//
//// NewLayoutTribeNet returns a layout with flat-top hexes for TribeNet.
////
//// All TribeNet maps are "odd-q,", origin is (1,1) and will
//// be translated to (0,0), sub-maps are 21 rows x 30 columns.
////
//// You may need to translate the origin from (0,0) to (1,1) when displaying TribeNet coordinates.
//func NewLayoutTribeNet() LayoutPointyTypeHorizontalOddRight_t {
//	return NewLayoutEvenQ(Point{1, 1}, Point{0, 0})
//}
//
//// HexFromCubeCoords returns a Hex initialized from cube coordinates.
//// Panics if q + r + s != 0.
//func HexFromCubeCoords(q, r, s int) Hex { // Cube constructor
//	if q+r+s != 0 {
//		panic("assert(q + r + s == 0)")
//	}
//	return Hex{q: q, r: r, s: -q - r}
//}

//// todo: an alternate implementation uses origin(0,0) and size(1,0) with chained transformations
//
//// hex→pixel: hex→cartesian, then scale the cartesian coordinate by multiplying by the desired scale, and then translate it to the desired origin.
//// pixel→hex: undo the translate by subtracting the origin, then undo the scale by dividing by the scale, then run cartesian→hex.
//
//// Offset implements offset coordinates for hexes. This can
//// be used to display player-friendly coordinates.
//type Offset struct {
//	Col, Row int
//}
//
//// Map_i implements an interface for additional storage like terrain,
//// objects, units, etc.
//type Map_i interface{}
//
//// Screen_i defines the interface for converting hex coordinates into
//// screen space (pixels).
////
//// Orientation is important for offset coordinates and every layout
//// that implements this interface is expected to implement that per
//// the Red Blob Games guide.
//type Screen_i interface {
//	// IsHorizontal returns true if the layout has horizontal rows.
//	// Horizontal layouts have pointy-top hexes, staggered columns, and horizontal rows.
//	IsHorizontal() bool
//
//	// IsVertical returns true if the layout has vertical columns.
//	// Vertical layouts have flat-top hexes, vertical columns, and staggered rows.
//	IsVertical() bool
//
//	// OffsetType returns the type of offset used for columns and rows.
//	OffsetType() LayoutOffset_e
//
//	// DirectionToBearing returns the bearing of a direction in the layout
//	DirectionToBearing(direction int) string
//
//	// HexagonalGrid returns a grid centered about a hex.
//	HexagonalGrid(center Hex, radius int) HashTable
//
//	// HexCorner returns the screen coordinates of the hex corner.
//	// We should define what "corner" means in this context.
//	HexCorner(h Hex, corner int) Point
//
//	// HexCorners returns the screen coordinates for every corner of the hex.
//	HexCorners(h Hex) [6]Point
//
//	// HexToOffsetCoord returns the offset coordinates of the hex.
//	// Uses the offset from the layout to shift rows and columns correctly.
//	HexToOffsetCoord(h Hex) OffsetCoord
//
//	// HexToPixel returns the origin of the hex on the screen as a pixel.
//	HexToPixel(h Hex) Point
//
//	// OffsetColRowToHex returns a new Hex using offset column and row coordinates.
//	OffsetColRowToHex(col, row int) Hex
//
//	// OffsetCoordToHex returns a new Hex from the OffsetCoord.
//	OffsetCoordToHex(oc OffsetCoord) Hex
//
//	// ParallelogramGrid returns a grid originating at (0,0,0).
//	// I don't understand the comment in the source about there
//	// being three coordinates and the caller has to choose two.
//	// does that mean the grid has three orientations?
//	ParallelogramGrid(q1, r1, q2, r2 int) HashTable
//
//	// PixelToHexRounded turns a fractional hex into a regular hex coordinate:
//	PixelToHexRounded(p Point) Hex
//
//	// PixelToFractionalHex returns the fractional hex that encloses the pixel.
//	// In theory, the origin of that fractional hex will be the pixel.
//	PixelToFractionalHex(p Point) FractionalHex
//
//	// PolygonCornerOffset returns the offset from the center of a hex to a corner.
//	// We should define what the parameter "corner" means. Which corner?
//	PolygonCornerOffset(corner int) Point
//
//	// PolygonCornerOffsets returns the offset for every corner of a hex.
//	PolygonCornerOffsets() [6]Point
//
//	// RectangularGrid returns a grid centered about a hex.
//	RectangularGrid(center Hex, left, right, top, bottom int) HashTable
//
//	// TriagonalGrid returns a grid originating at (0,0,0).
//	// there's a comment in the source about flipping the y-axis to
//	// change the direction of the triangle, but I don't understand
//	// how to implement that.
//	TriagonalGrid(side_length int) HashTable
//}
