// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import (
	"errors"
	"slices"
)

// Cylinder is an offset rectangle that wraps east-west and has walled north
// and south edges. Canonical columns are 0..Columns()-1 and rows are 0..Rows()-1.
// Construct it with NewCylinder; the zero value is not usable.
type Cylinder struct {
	layout  Layout
	columns int
	rows    int
}

// NewCylinder returns a cylinder using layout's offset convention; pixel size
// and origin do not affect topology. Q layouts require an even column count
// of at least four so the seam preserves column parity. R layouts require an
// odd count of at least three, following the pointy-top cylinder convention.
// Rows must be positive and columns*rows must fit in an int.
// Invalid dimensions return an error. Negative Ring/Spiral radii panic instead,
// as they do on Hex.
func NewCylinder(layout Layout, columns, rows int) (Cylinder, error) {
	if layout.IsFlat() {
		if columns < 4 || columns%2 != 0 {
			return Cylinder{}, errors.New("cylinder: Q layouts require even columns >= 4")
		}
	} else if columns < 3 || columns%2 != 1 {
		return Cylinder{}, errors.New("cylinder: R layouts require odd columns >= 3")
	}
	if rows <= 0 || rows > int(^uint(0)>>1)/columns {
		return Cylinder{}, errors.New("cylinder: rows must be positive and columns*rows must fit in an int")
	}
	return Cylinder{layout: layout, columns: columns, rows: rows}, nil
}

// Columns returns the total column count.
func (c Cylinder) Columns() int { return c.columns }

// Rows returns the total row count.
func (c Cylinder) Rows() int { return c.rows }

// Normalize wraps h's offset column into the canonical range, leaving its offset
// row unchanged. It does not clip rows, so it is defined even outside the walls.
func (c Cylinder) Normalize(h Hex) Hex {
	oc := c.layout.CubeToOffset(h)
	oc.Col %= c.columns
	if oc.Col < 0 {
		oc.Col += c.columns
	}
	return c.layout.OffsetToCube(oc)
}

// Contains reports whether h's offset row is inside the walls. Noncanonical
// columns are accepted because they identify the same hex after wrapping.
func (c Cylinder) Contains(h Hex) bool {
	row := c.layout.CubeToOffset(h).Row
	return 0 <= row && row < c.rows
}

// Neighbor returns the wrapped neighbor in direction (coerced to 0..5).
// It returns the zero Hex and false if h or the neighbor is outside the walls.
func (c Cylinder) Neighbor(h Hex, direction int) (Hex, bool) {
	if !c.Contains(h) {
		return Hex{}, false
	}
	n := c.Normalize(h).Neighbor(direction)
	if !c.Contains(n) {
		return Hex{}, false
	}
	return c.Normalize(n), true
}

// Neighbors returns the on-map neighbors in direction order 0..5, omitting
// steps through a wall. It returns an empty slice if h is outside the walls.
func (c Cylinder) Neighbors(h Hex) []Hex {
	result := make([]Hex, 0, 6)
	for direction := range 6 {
		if n, ok := c.Neighbor(h, direction); ok {
			result = append(result, n)
		}
	}
	return result
}

// nearest returns a copy of canonical b closest to canonical a. Ties prefer
// the canonical copy, then the western copy, then the eastern copy.
func (c Cylinder) nearest(a, b Hex) Hex {
	oc := c.layout.CubeToOffset(b)
	best, distance := b, a.Distance(b)
	for _, shift := range [...]int{-c.columns, c.columns} {
		copy := c.layout.OffsetToCube(OffsetCoord{Col: oc.Col + shift, Row: oc.Row})
		if d := a.Distance(copy); d < distance {
			best, distance = copy, d
		}
	}
	return best
}

// Distance returns the shortest hex distance around the cylinder. Inputs are
// normalized internally. Like Normalize, it is defined outside the row bounds;
// for on-map inputs it is also the shortest distance without crossing a wall.
func (c Cylinder) Distance(a, b Hex) int {
	a, b = c.Normalize(a), c.Normalize(b)
	return a.Distance(c.nearest(a, b))
}

// Ring returns exactly the on-map hexes at wrapped distance radius from center,
// without duplicates, sorted by Hex.Compare. An off-map center returns an empty
// slice. A negative radius panics.
func (c Cylinder) Ring(center Hex, radius int) []Hex {
	if radius < 0 {
		panic("radius must be non-negative")
	}
	center = c.Normalize(center)
	// Walking along a row and then a column gives this diameter upper bound.
	if !c.Contains(center) || radius > c.rows-1+c.columns {
		return []Hex{}
	}
	if radius == 0 {
		return []Hex{center}
	}
	result := make(HexSet)
	h := center.Add(DirectionVector(4).Multiply(radius))
	for direction := range 6 {
		for range radius {
			n := c.Normalize(h)
			// Normalizing a planar ring alone is wrong beyond half a lap:
			// some hexes have a shorter route through a different copy.
			if c.Contains(n) && c.Distance(center, n) == radius {
				result[n] = struct{}{}
			}
			h = h.Neighbor(direction)
		}
	}
	return result.Sorted()
}

// Spiral returns all on-map hexes within wrapped distance radius from center,
// without duplicates, sorted by Hex.Compare (not spiral traversal order).
// An off-map center returns an empty slice. A negative radius panics.
func (c Cylinder) Spiral(center Hex, radius int) []Hex {
	if radius < 0 {
		panic("radius must be non-negative")
	}
	center = c.Normalize(center)
	if !c.Contains(center) {
		return []Hex{}
	}
	oc := c.layout.CubeToOffset(center)
	// Every step changes each offset coordinate by at most one. Enumerate
	// each candidate once, bounded by the rectangle even for huge radii.
	firstCol, count := 0, c.columns
	if radius < c.columns/2 {
		firstCol, count = oc.Col-radius, 2*radius+1
	}
	firstRow := oc.Row - min(radius, oc.Row)
	lastRow := oc.Row + min(radius, c.rows-1-oc.Row)
	result := make([]Hex, 0)
	for row := firstRow; row <= lastRow; row++ {
		for i := range count {
			h := c.Normalize(c.layout.OffsetToCube(OffsetCoord{Col: firstCol + i, Row: row}))
			if c.Distance(center, h) <= radius {
				result = append(result, h)
			}
		}
	}
	slices.SortFunc(result, Hex.Compare)
	return result
}

// LineDraw returns a shortest on-map line from a to b, including both endpoints.
// Inputs and results are canonicalized. Off-map endpoints return an empty slice.
// Equally short copies of b prefer canonical, then west, then east. At each step
// the line takes the shortest-path neighbor closest to the ideal straight line
// in cube space; ties prefer direction order 0..5. This keeps the line inside
// the walls even when rounding a planar line would put it outside the rectangle.
func (c Cylinder) LineDraw(a, b Hex) []Hex {
	a, b = c.Normalize(a), c.Normalize(b)
	if !c.Contains(a) || !c.Contains(b) {
		return []Hex{}
	}
	b = c.nearest(a, b)
	distance := a.Distance(b)
	result := []Hex{a}
	h := a
	for step := 1; step <= distance; step++ {
		ideal := a.Lerp(b, float64(step)/float64(distance))
		var next Hex
		best := -1.0
		for direction := range 6 {
			n := h.Neighbor(direction)
			if !c.Contains(n) || n.Distance(b) != distance-step {
				continue
			}
			dq, dr, ds := float64(n.q)-ideal.q, float64(n.r)-ideal.r, float64(n.s)-ideal.s
			score := dq*dq + dr*dr + ds*ds
			if best < 0 || score < best {
				next, best = n, score
			}
		}
		h = next
		result = append(result, c.Normalize(h))
	}
	return result
}

// Index returns the dense row-major index row*Columns()+col after normalizing
// h. It returns zero and false if h is outside the walls.
func (c Cylinder) Index(h Hex) (int, bool) {
	if !c.Contains(h) {
		return 0, false
	}
	oc := c.layout.CubeToOffset(c.Normalize(h))
	return oc.Row*c.columns + oc.Col, true
}

// HexAt is the inverse of Index. It returns a canonical hex, or the zero Hex
// and false if index is outside 0..Columns()*Rows()-1.
func (c Cylinder) HexAt(index int) (Hex, bool) {
	if index < 0 || index >= c.columns*c.rows {
		return Hex{}, false
	}
	return c.layout.OffsetToCube(OffsetCoord{Col: index % c.columns, Row: index / c.columns}), true
}
