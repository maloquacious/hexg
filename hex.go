// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package hexg implements hex grids in the style of the Red Blob Games blog
// https://www.redblobgames.com/grids/hexagons/ and
// https://www.redblobgames.com/grids/hexagons/implementation.html.
package hexg

// NewHex returns a Hex initialized from axial coordinates.
func NewHex(q, r int) Hex { // Axial constructor
	return Hex{q: q, r: r, s: -q - r}
}

// Hex stores the `q`, `r`, and `s` coordinates.
type Hex struct {
	q, r, s int
}

// Q returns the q coordinate.
func (h Hex) Q() int { return h.q }

// R returns the r coordinate.
func (h Hex) R() int { return h.r }

// S returns the s coordinate.
func (h Hex) S() int { return h.s }

// QRS returns all three coordinates.
func (h Hex) QRS() (int, int, int) { return h.q, h.r, h.s }

// Equals returns true if the two hexes have the same coordinates.
func (h Hex) Equals(b Hex) bool {
	return h.q == b.q && h.r == b.r && h.s == b.s
}

// NotEquals returns true if the two hexes have different coordinates.
func (h Hex) NotEquals(b Hex) bool {
	return !h.Equals(b)
}

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

// Length is the distance from the origin to the hex.
func (h Hex) Length() int {
	return (abs(h.q) + abs(h.r) + abs(h.s)) / 2
}

// Distance between two hexes is the length of the line between them.
func (h Hex) Distance(b Hex) int {
	return h.Subtract(b).Length()
}

var directionVectors = [6]Hex{
	{1, 0, -1}, {1, -1, 0}, {0, -1, 1},
	{-1, 0, 1}, {-1, 1, 0}, {0, 1, -1},
}

// DirectionVector returns the unit vector for a hex direction.
// Direction must be in the range 0..5. Panics on invalid input.
func DirectionVector(direction int) Hex {
	return directionVectors[direction]
}

// Neighbor returns the hex that is one step away in the given direction.
// Direction is coerced to the range 0..5.
func (h Hex) Neighbor(direction int) Hex {
	direction = (6 + (direction % 6)) % 6
	return h.Add(DirectionVector(direction))
}

// Lerp returns the linear interpolation of points on the line between two hexes
func (h Hex) Lerp(b Hex, t float64) FractionalHex {
	return FractionalHex{
		q: lerp(float64(h.q), float64(b.q), t),
		r: lerp(float64(h.r), float64(b.r), t),
		s: lerp(float64(h.s), float64(b.s), t),
	}
}

// LineDraw returns the set of hexes on the line between two hexes.
func (h Hex) LineDraw(b Hex) []Hex {
	return h.lineDraw(b, false)
}

// LineDrawNudged returns the set of hexes on the line between two hexes,
// with nudging to push points on an edge in a consistent direction.
func (h Hex) LineDrawNudged(b Hex) []Hex {
	return h.lineDraw(b, true)
}

func (h Hex) lineDraw(b Hex, withNudge bool) []Hex {
	N := h.Distance(b)
	var results []Hex
	var step float64
	if N == 0 {
		step = 1.0
	} else {
		step = 1.0 / float64(N)
	}
	if withNudge {
		aNudge := FractionalHex{q: float64(h.q) + 1e-6, r: float64(h.r) + 1e-6, s: float64(h.s) - 2e-6}
		bNudge := FractionalHex{q: float64(b.q) + 1e-6, r: float64(b.r) + 1e-6, s: float64(b.s) - 2e-6}
		for i := 0; i <= N; i++ {
			results = append(results, aNudge.Lerp(bNudge, step*float64(i)).Round())
		}
		return results
	}
	for i := 0; i <= N; i++ {
		results = append(results, h.Lerp(b, step*float64(i)).Round())
	}
	return results
}

// HexSet is a set of hexes.
type HexSet map[Hex]struct{}

// RotateLeft shifts the q, r, and s coordinates.
// The effect depends on the layout.
func (h Hex) RotateLeft() Hex {
	return Hex{q: -h.s, r: -h.q, s: -h.r}
}

// RotateRight shifts the q, r, and s coordinates.
// The effect depends on the layout.
func (h Hex) RotateRight() Hex {
	return Hex{q: -h.r, r: -h.s, s: -h.q}
}
