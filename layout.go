// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import "math"

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
	panic("invalid offset")
}

// IsEven returns true if the layout uses even offset coordinates.
func (l *Layout) IsEven() bool {
	return l.offset == EvenQ || l.offset == EvenR
}

// IsEvenQ returns true if the layout supports even-q offset coordinates (vertical layout, shoves even columns down).
func (l *Layout) IsEvenQ() bool {
	return l.offset == EvenQ
}

// IsEvenR returns true if the layout supports even-r offset coordinates (horizontal layout, shoves even rows right).
func (l *Layout) IsEvenR() bool {
	return l.offset == EvenR
}

// IsFlat returns true if the layout uses flat-top hexes.
func (l *Layout) IsFlat() bool {
	return l.offset == OddQ || l.offset == EvenQ
}

// IsHorizontal returns true if the layout has horizontal rows.
func (l *Layout) IsHorizontal() bool {
	return l.offset == OddR || l.offset == EvenR
}

// IsOdd returns true if the layout uses odd offset coordinates.
func (l *Layout) IsOdd() bool {
	return l.offset == OddQ || l.offset == OddR
}

// IsOddQ returns true if the layout supports odd-q offset coordinates (vertical layout, shoves odd columns down).
func (l *Layout) IsOddQ() bool {
	return l.offset == OddQ
}

// IsOddR returns true if the layout supports odd-r offset coordinates (horizontal layout, shoves odd rows right).
func (l *Layout) IsOddR() bool {
	return l.offset == OddR
}

// IsPointy returns true if the layout uses pointy-top hexes.
func (l *Layout) IsPointy() bool {
	return l.offset == OddR || l.offset == EvenR
}

// IsVertical returns true if the layout has vertical columns.
func (l *Layout) IsVertical() bool {
	return l.offset == OddQ || l.offset == EvenQ
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

// CubeToOffset converts cube coordinates to offset coordinates based on the layout.
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
	panic("invalid offset")
}

// OffsetToCube converts offset coordinates to cube coordinates based on the layout.
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
	panic("invalid offset")
}

// RotateLeft rotates the hex 60° counter-clockwise around the origin.
// The underlying operation depends on the layout orientation.
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
	panic("invalid offset")
}

// RotateRight rotates the hex 60° clockwise around the origin.
// The underlying operation depends on the layout orientation.
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
	panic("invalid offset")
}

// HexToPixel returns the origin of the hex on the grid as a Point.
func (l *Layout) HexToPixel(h Hex) Point {
	M := l.orientation
	x := (M.f0*float64(h.q) + M.f1*float64(h.r)) * l.size.X
	y := (M.f2*float64(h.q) + M.f3*float64(h.r)) * l.size.Y
	return Point{
		X: x + l.origin.X,
		Y: y + l.origin.Y,
	}
}

// PixelToFractionalHex returns the fractional hex that encloses the pixel.
// In theory, the origin of that fractional hex will be the pixel.
func (l *Layout) PixelToFractionalHex(p Point) FractionalHex {
	M := l.orientation
	pt := Point{
		X: (p.X - l.origin.X) / l.size.X,
		Y: (p.Y - l.origin.Y) / l.size.Y,
	}
	q := M.b0*pt.X + M.b1*pt.Y
	r := M.b2*pt.X + M.b3*pt.Y
	return FractionalHex{q: q, r: r, s: -q - r}
}

// HexCornerOffset returns the screen location (pixel) of a corner of a hex on the grid.
func (l *Layout) HexCornerOffset(corner int) Point {
	size := l.size
	angle := 2.0 * math.Pi *
		(l.orientation.start_angle + float64(corner)) / 6
	return Point{
		X: size.X * math.Cos(angle),
		Y: size.Y * math.Sin(angle),
	}
}

// PolygonCorners returns the location of the six corners of the hex on the grid.
func (l *Layout) PolygonCorners(h Hex) [6]Point {
	var corners [6]Point
	center := l.HexToPixel(h)

	for i := 0; i < 6; i++ {
		offset := l.HexCornerOffset(i)
		corners[i] = Point{
			X: center.X + offset.X,
			Y: center.Y + offset.Y,
		}
	}
	return corners
}

// PixelToHexRounded converts a screen pixel to the nearest hex coordinate.
func (l *Layout) PixelToHexRounded(p Point) Hex {
	return l.PixelToFractionalHex(p).Round()
}

// ParallelogramQR returns a parallelogram-shaped grid using q and r axes.
func (l *Layout) ParallelogramQR(q1, r1 int, q2, r2 int) HexSet {
	gs := make(HexSet)
	for q := q1; q <= q2; q++ {
		for r := r1; r <= r2; r++ {
			gs[Hex{q: q, r: r, s: -q - r}] = struct{}{}
		}
	}
	return gs
}

// ParallelogramQS returns a parallelogram-shaped grid using q and s axes.
func (l *Layout) ParallelogramQS(q1, s1 int, q2, s2 int) HexSet {
	gs := make(HexSet)
	for q := q1; q <= q2; q++ {
		for s := s1; s <= s2; s++ {
			gs[Hex{q: q, r: -q - s, s: s}] = struct{}{}
		}
	}
	return gs
}

// ParallelogramRS returns a parallelogram-shaped grid using r and s axes.
func (l *Layout) ParallelogramRS(r1, s1 int, r2, s2 int) HexSet {
	gs := make(HexSet)
	for r := r1; r <= r2; r++ {
		for s := s1; s <= s2; s++ {
			gs[Hex{q: -r - s, r: r, s: s}] = struct{}{}
		}
	}
	return gs
}

// TriangleUpDown returns a grid originating at (0,0,0).
// mapSize is the length of a side.
func (l *Layout) TriangleUpDown(mapSize int) HexSet {
	gs := make(HexSet)
	for q := 0; q <= mapSize; q++ {
		for r := 0; r <= mapSize-q; r++ {
			gs[Hex{q: q, r: r, s: -q - r}] = struct{}{}
		}
	}
	return gs
}

// TriangleLeftRight returns a grid originating at (0,0,0).
// mapSize is the length of a side.
func (l *Layout) TriangleLeftRight(mapSize int) HexSet {
	gs := make(HexSet)
	for q := 0; q <= mapSize; q++ {
		for r := mapSize - q; r <= mapSize; r++ {
			gs[Hex{q: q, r: r, s: -q - r}] = struct{}{}
		}
	}
	return gs
}

// Hexagon returns a grid centered about (0,0,0).
// Does not depend on the orientation of the grid.
func (l *Layout) Hexagon(radius int) HexSet {
	gs := make(HexSet)
	for q := -radius; q <= radius; q++ {
		r1 := max(-radius, -q-radius)
		r2 := min(radius, -q+radius)
		for r := r1; r <= r2; r++ {
			gs[Hex{q: q, r: r, s: -q - r}] = struct{}{}
		}
	}
	return gs
}

// Rectangle returns a grid centered about (0,0,0).
// The internal logic depends on the orientation of the grid.
func (l *Layout) Rectangle(left, right, top, bottom int) HexSet {
	gs := make(HexSet)
	if l.IsPointy() {
		for r := top; r <= bottom; r++ {
			rOffset := r >> 1
			for q := left - rOffset; q <= right-rOffset; q++ {
				gs[Hex{q: q, r: r, s: -q - r}] = struct{}{}
			}
		}
		return gs
	}
	for q := left; q <= right; q++ {
		qOffset := q >> 1
		for r := top - qOffset; r <= bottom-qOffset; r++ {
			gs[Hex{q: q, r: r, s: -q - r}] = struct{}{}
		}
	}
	return gs
}
