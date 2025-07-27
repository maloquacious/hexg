// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import (
	"math"
)

func NewVerticalOddQLayout(size, origin Point) VerticalOddQLayout {
	return VerticalOddQLayout{
		size:   size,
		origin: origin,
	}
}

// VerticalOddQLayout returns a layout with vertical layout (flat-top hexes)
// that shoves odd columns down.
type VerticalOddQLayout struct {
	// size and origin are used when calculating screen pixels.
	size   Point
	origin Point
}

func (l VerticalOddQLayout) DirectionToBearing(direction int) string {
	// we must coerce direction to 0 ... 5
	return verticalDirectionToBearing[(6+(direction%6))%6]
}

func (l VerticalOddQLayout) HexagonalGrid(center Hex, radius int) GridStore {
	gs := GridStore{}
	N := radius
	for q := -N; q <= N; q++ {
		r1 := max(-N, -q-N)
		r2 := min(N, -q+N)
		for r := r1; r <= r2; r++ {
			hex := center.Add(NewHexFromAxialCoords(q, r))
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

func (l VerticalOddQLayout) HexCorner(h Hex, corner int) Point {
	center := l.HexToPixel(h)
	offset := l.PolygonCornerOffset(corner)
	return Point{X: center.X + offset.X, Y: center.Y + offset.Y}
}

func (l VerticalOddQLayout) HexCorners(h Hex) [6]Point {
	center := l.HexToPixel(h)
	corners := l.PolygonCornerOffsets()
	for i := 0; i < 6; i++ {
		corners[i].X, corners[i].Y = center.X+corners[i].X, center.Y+corners[i].Y
	}
	return corners
}

// HexToOffsetCoord returns the offset coordinates of the hex.
// Uses the offset from the layout to shift rows and columns correctly.
func (l VerticalOddQLayout) HexToOffsetCoord(h Hex) OffsetCoord {
	col, row := h.q, h.r+(h.q+ODD*(h.q&1))/2
	return OffsetCoord{Col: col, Row: row}
}

func (l VerticalOddQLayout) HexToPixel(h Hex) Point {
	M := verticalOrientation
	return Point{
		X: l.origin.X + (M.f0*float64(h.q)+M.f1*float64(h.r))*l.size.X,
		Y: l.origin.Y + (M.f2*float64(h.q)+M.f3*float64(h.r))*l.size.Y,
	}
}

func (l VerticalOddQLayout) IsHorizontal() bool {
	return false
}

func (l VerticalOddQLayout) IsVertical() bool {
	return true
}

func (l VerticalOddQLayout) OffsetColRowToHex(col, row int) Hex {
	q, r := col, row-(col+ODD*(col&1))/2
	return Hex{q: q, r: r, s: -q - r}
}

func (l VerticalOddQLayout) OffsetCoordToHex(oc OffsetCoord) Hex {
	return l.OffsetColRowToHex(oc.Col, oc.Row)
}

func (l VerticalOddQLayout) OffsetType() LayoutOffset_e {
	return OddQ
}

func (l VerticalOddQLayout) ParallelogramGrid(q1, r1, q2, r2 int) GridStore {
	gs := GridStore{}
	for q := q1; q <= q2; q++ {
		for r := r1; r <= r2; r++ {
			hex := NewHexFromAxialCoords(q, r)
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

func (l VerticalOddQLayout) PixelToFractionalHex(p Point) FractionalHex {
	M := verticalOrientation
	pt := Point{X: (p.X - l.origin.X) / l.size.X, Y: (p.Y - l.origin.Y) / l.size.Y}
	q := M.b0*pt.X + M.b1*pt.Y
	r := M.b2*pt.X + M.b3*pt.Y
	return FractionalHex{q: q, r: r, s: -q - r}
}

func (l VerticalOddQLayout) PixelToHexRounded(p Point) Hex {
	return l.PixelToFractionalHex(p).Round()
}

func (l VerticalOddQLayout) PolygonCornerOffset(corner int) Point {
	M := verticalOrientation
	size := l.size
	// todo: maybe explain why adding corner to start_angle is correct
	angle := 2.0 * math.Pi * (M.start_angle + float64(corner)) / 6
	return Point{X: size.X * math.Cos(angle), Y: size.Y * math.Sin(angle)}
}

func (l VerticalOddQLayout) PolygonCornerOffsets() [6]Point {
	var corners [6]Point
	center := l.HexToPixel(Hex{})
	for i := 0; i < 6; i++ {
		offset := l.PolygonCornerOffset(i)
		corners[i] = Point{X: center.X + offset.X, Y: center.Y + offset.Y}
	}
	return corners
}

func (l VerticalOddQLayout) RectangularGrid(center Hex, left, right, top, bottom int) GridStore {
	gs := GridStore{}
	for q := left; q <= right; q++ {
		q_offset := q >> 1 // or math.Floor(float64(q) / 2.0)
		for r := top - q_offset; r <= bottom-q_offset; r++ {
			hex := center.Add(NewHexFromAxialCoords(q, r))
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

func (l VerticalOddQLayout) TriagonalGrid(side_length int) GridStore {
	gs := GridStore{}
	map_size := side_length
	for q := 0; q <= map_size; q++ {
		for r := map_size - q; r <= map_size; r++ {
			hex := NewHexFromAxialCoords(q, r)
			gs[hex.Hash()] = hex
		}
	}
	return gs
}
