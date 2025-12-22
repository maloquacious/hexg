// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import (
	"fmt"
	"math"
)

type Orientation int

const (
	// LayoutPointy is a pointy top orientation:
	// * staggered columns
	// * horizontal rows
	LayoutPointy Orientation = iota

	// LayoutFlat is a flat top orientation:
	// * vertical columns
	// * staggered rows
	LayoutFlat
)

// LayoutOffset describes the layout (horizontal row or vertical columns)
// and which way odd rows and columns are pushed or pulled.
type LayoutOffset int

const (
	// OddR is a pointy-top layout with horizontal rows that shoves odd rows right
	OddR LayoutOffset = iota
	// EvenR is a pointy-top layout with horizontal rows that shoves even rows right
	EvenR
	// OddQ is a flat-top layout with vertical columns that shoves odd columns down
	OddQ
	// EvenQ is a flat-top layout with vertical columns that shoves even columns down
	EvenQ
)

type orientation struct {
	f0, f1, f2, f3 float64
	b0, b1, b2, b3 float64
	start_angle    float64 // in multiples of 60°
}

// Layout implements the conversion between hex (q,r,s) and screen (x,y)
// coordinates. The conversion uses a matrix as well as the inverse of
// the matrix for each orientation.
type Layout struct {
	orientation orientation
	offset      LayoutOffset

	// origin is center of the (0,0,0) hex. It is used in translate
	// transformation to move hexes on the screen. If you don’t need
	// this, set it to Point(0, 0).
	origin Point

	// size is used for non-uniform scaling, especially for matching
	// pixel sprite sizes. It’s a scale transform. If you need uniform
	// scaling, set it to Point(width, height), where width == height.
	size Point
}

// NewLayout returns a layout with either flat-top or pointy-top hexes.
//
// Flat-top layouts have vertical columns and staggered rows.
// Drawing corners start at  0°.
//
// Pointy-top layouts have staggered columns and horizontal rows.
// Drawing corners start at 30°.
func NewLayout(offset LayoutOffset, size, origin Point) Layout {
	switch offset {
	case OddR: // OddR is a pointy-top layout with horizontal rows that shoves odd rows right
		return Layout{
			offset: OddR,
			orientation: orientation{
				f0: math.Sqrt(3.0), f1: math.Sqrt(3.0) / 2.0, f2: 0.0, f3: 3.0 / 2.0,
				b0: math.Sqrt(3.0) / 3.0, b1: -1.0 / 3.0, b2: 0.0, b3: 2.0 / 3.0,
				start_angle: 0.5,
			},
			origin: origin,
			size:   size,
		}
	case EvenR: // EvenR is a pointy-top layout with horizontal rows that shoves even rows right
		return Layout{
			offset: EvenR,
			orientation: orientation{
				f0: math.Sqrt(3.0), f1: math.Sqrt(3.0) / 2.0, f2: 0.0, f3: 3.0 / 2.0,
				b0: math.Sqrt(3.0) / 3.0, b1: -1.0 / 3.0, b2: 0.0, b3: 2.0 / 3.0,
				start_angle: 0.5,
			},
			origin: origin,
			size:   size,
		}
	case OddQ: // OddQ is a flat-top layout with vertical columns that shoves odd columns down
		return Layout{
			offset: OddQ,
			orientation: orientation{
				f0: 3.0 / 2.0, f1: 0.0, f2: math.Sqrt(3.0) / 2.0, f3: math.Sqrt(3.0),
				b0: 2.0 / 3.0, b1: 0.0, b2: -1.0 / 3.0, b3: math.Sqrt(3.0) / 3.0,
				start_angle: 0.0,
			},
			origin: origin,
			size:   size,
		}
	case EvenQ: // EvenQ is a flat-top layout with vertical columns that shoves even columns down
		return Layout{
			offset: EvenQ,
			orientation: orientation{
				f0: 3.0 / 2.0, f1: 0.0, f2: math.Sqrt(3.0) / 2.0, f3: math.Sqrt(3.0),
				b0: 2.0 / 3.0, b1: 0.0, b2: -1.0 / 3.0, b3: math.Sqrt(3.0) / 3.0,
				start_angle: 0.0,
			},
			origin: origin,
			size:   size,
		}
	}
	panic("assert(offset in (OddR, EvenR, OddQ, EvenQ)")
}

func (layout Layout) IsEven() bool {
	return layout.offset == EvenQ || layout.offset == EvenR
}

// IsEvenQ returns true if the layout supports even-q offset coordinates (vertical layout, shoves even columns down).
func (layout Layout) IsEvenQ() bool {
	return layout.offset == EvenQ
}

// IsEvenR returns true if the layout supports even-r offset coordinates (horizontal layout, shoves even rows right).
func (layout Layout) IsEvenR() bool {
	return layout.offset == EvenR
}

func (layout Layout) IsFlat() bool {
	return layout.offset == OddQ || layout.offset == EvenQ
}

func (layout Layout) IsHorizontal() bool {
	return layout.offset == OddR || layout.offset == EvenR
}

func (layout Layout) IsOdd() bool {
	return layout.offset == OddQ || layout.offset == OddR
}

// IsOddQ returns true if the layout supports odd-q offset coordinates (vertical layout, shoves odd columns down).
func (layout Layout) IsOddQ() bool {
	return layout.offset == OddQ
}

// IsOddR returns true if the layout supports odd-r offset coordinates (horizontal layout, shoves odd rows right).
func (layout Layout) IsOddR() bool {
	return layout.offset == OddR
}

func (layout Layout) IsPointy() bool {
	return layout.offset == OddR || layout.offset == EvenR
}

func (layout Layout) IsVertical() bool {
	return layout.offset == OddQ || layout.offset == EvenQ
}

// BoundingBox returns the bounding box of a list of hexes.
// If the list is empty, returns (0,0,0), (0,0,0)
func (l *Layout) BoundingBox(hexes ...Hex) (upperLeft, lowerRight Hex) {
	if len(hexes) == 0 {
		return Hex{}, Hex{}
	}

	var upperLeftCoords, lowerRightCoords OffsetCoord

	for n, h := range hexes {
		oc := l.CubeToOffset(h)
		if n == 0 {
			upperLeftCoords, lowerRightCoords = oc, oc
			continue
		}
		if oc.Col < upperLeftCoords.Col {
			upperLeftCoords.Col = oc.Col
		}
		if oc.Row < upperLeftCoords.Row {
			upperLeftCoords.Row = oc.Row
		}
		if oc.Col > lowerRightCoords.Col {
			lowerRightCoords.Col = oc.Col
		}
		if oc.Row > lowerRightCoords.Row {
			lowerRightCoords.Row = oc.Row
		}
	}

	return l.OffsetToCube(upperLeftCoords), l.OffsetToCube(lowerRightCoords)
}

func (l *Layout) CubeToOffset(h Hex) OffsetCoord {
	switch l.offset {
	case OddR:
		return h.CubeToROffset(false)
	case EvenR:
		return h.CubeToROffset(true)
	case OddQ:
		return h.CubeToQOffset(false)
	case EvenQ:
		return h.CubeToQOffset(true)
	}
	panic(fmt.Sprintf("assert(offset != %d)", l.offset))
}

func (l *Layout) OffsetToCube(oc OffsetCoord) Hex {
	switch l.offset {
	case OddR:
		return oc.ROffsetToCube(false)
	case EvenR:
		return oc.ROffsetToCube(true)
	case OddQ:
		return oc.QOffsetToCube(false)
	case EvenQ:
		return oc.QOffsetToCube(true)
	}
	panic(fmt.Sprintf("assert(offset != %d)", l.offset))
}

func (l *Layout) RotateLeft(h Hex) Hex {
	switch l.offset {
	case OddR:
		return h.RotateLeft()
	case EvenR:
		return h.RotateLeft()
	case OddQ:
		return h.RotateRight()
	case EvenQ:
		return h.RotateRight()
	}
	panic(fmt.Sprintf("assert(offset != %d)", l.offset))
}

func (l *Layout) RotateRight(h Hex) Hex {
	switch l.offset {
	case OddR:
		return h.RotateRight()
	case EvenR:
		return h.RotateRight()
	case OddQ:
		return h.RotateLeft()
	case EvenQ:
		return h.RotateLeft()
	}
	panic(fmt.Sprintf("assert(offset != %d)", l.offset))
}

// HexToPixel returns the origin of the hex on the grid as a Point.
func (layout Layout) HexToPixel(h Hex) Point {
	M := layout.orientation
	x := (M.f0*float64(h.q) + M.f1*float64(h.r)) * layout.size.X
	y := (M.f2*float64(h.q) + M.f3*float64(h.r)) * layout.size.Y
	return Point{
		X: x + layout.origin.X,
		Y: y + layout.origin.Y,
	}
}

// PixelToFractionalHex returns the fractional hex that encloses the pixel.
// In theory, the origin of that fractional hex will be the pixel.
func (layout Layout) PixelToFractionalHex(p Point) FractionalHex {
	M := layout.orientation
	pt := Point{
		X: (p.X - layout.origin.X) / layout.size.X,
		Y: (p.Y - layout.origin.Y) / layout.size.Y,
	}
	q := M.b0*pt.X + M.b1*pt.Y
	r := M.b2*pt.X + M.b3*pt.Y
	return FractionalHex{q: q, r: r, s: -q - r}
}

// To draw a hex, get the center of the hex, then find the corners.

// HexCornerOffset returns the screen location (pixel) of a corner of a hex on the grid.
func (layout Layout) HexCornerOffset(corner int) Point {
	size := layout.size
	angle := 2.0 * math.Pi *
		(layout.orientation.start_angle + float64(corner)) / 6
	return Point{
		X: size.X * math.Cos(angle),
		Y: size.Y * math.Sin(angle),
	}
}

// PolygonCorners returns the location of the six corners of the hex on the grid.
func (layout Layout) PolygonCorners(h Hex) [6]Point {
	var corners [6]Point
	center := layout.HexToPixel(h)

	for i := 0; i < 6; i++ {
		offset := layout.HexCornerOffset(i)
		corners[i] = Point{
			X: center.X + offset.X,
			Y: center.Y + offset.Y,
		}
	}
	return corners
}

// PixelToHexRounded turns a fractional hex into a regular hex coordinate:
func (layout Layout) PixelToHexRounded(p Point) Hex {
	return layout.PixelToFractionalHex(p).Round()
}

func (layout Layout) ParallelogramQR(q1, r1 int, q2, r2 int) HashTable {
	gs := make(map[uint64]Hex)
	for q := q1; q <= q2; q++ {
		for r := r1; r <= r2; r++ {
			hex := Hex{q: q, r: r, s: -q - r}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

func (layout Layout) ParallelogramQS(q1, s1 int, q2, s2 int) HashTable {
	gs := make(map[uint64]Hex)
	for q := q1; q <= q2; q++ {
		for s := s1; s <= s2; s++ {
			hex := Hex{q: q, r: -q - s, s: s}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

func (layout Layout) ParallelogramRS(r1, s1 int, r2, s2 int) HashTable {
	gs := make(map[uint64]Hex)
	for r := r1; r <= r2; r++ {
		for s := s1; s <= s2; s++ {
			hex := Hex{q: -r - s, r: r, s: s}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

// TriangleUpDown returns a grid originating at (0,0,0).
// `map_size` is the length of a side.
func (layout Layout) TriangleUpDown(map_size int) HashTable {
	gs := HashTable{}
	for q := 0; q <= map_size; q++ {
		for r := 0; r <= map_size-q; r++ {
			hex := Hex{q: q, r: r, s: -q - r}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

// TriangleLeftRight returns a grid originating at (0,0,0).
// `map_size` is the length of a side.
func (layout Layout) TriangleLeftRight(map_size int) HashTable {
	gs := HashTable{}
	for q := 0; q <= map_size; q++ {
		for r := map_size - q; r <= map_size; r++ {
			hex := Hex{q: q, r: r, s: -q - r}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

// Hexagon returns a grid centered about (0,0,0).
// does not depend on the orientation of the grid.
func (layout Layout) Hexagon(radius int) HashTable {
	gs := HashTable{}
	N := radius
	for q := -N; q <= N; q++ {
		r1 := max(-N, -q-N)
		r2 := min(N, -q+N)
		for r := r1; r <= r2; r++ {
			hex := Hex{q: q, r: r, s: -q - r}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

// Rectangle returns a grid centered about (0,0,0).
// the internal logic depends on the orientation of the grid.
func (layout Layout) Rectangle(left, right, top, bottom int) HashTable {
	gs := HashTable{}
	if layout.IsPointy() {
		for r := top; r <= bottom; r++ {
			r_offset := r >> 1 // or math.Floor(float64(r) / 2.0)
			for q := left - r_offset; q <= right-r_offset; q++ {
				hex := Hex{q: q, r: r, s: -q - r}
				gs[hex.Hash()] = hex
			}
		}
		return gs
	}
	for q := left; q <= right; q++ {
		q_offset := q >> 1 // or math.Floor(float64(q) / 2.0)
		for r := top - q_offset; r <= bottom-q_offset; r++ {
			hex := Hex{q: q, r: r, s: -q - r}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}
