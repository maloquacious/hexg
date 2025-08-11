// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package cube implements hexes with cube coordinates (q, r, s).
package cube

/////////////////////////////////////////////////////////////////////////////
// coordinate systems
// * https://www.redblobgames.com/grids/hexagons/#coordinates

// cube coordinates
// * https://www.redblobgames.com/grids/hexagons/#coordinates-cube

// Cube implements coordinates for hexes on a grid.
// See https://www.redblobgames.com/grids/hexagons/#coordinates-cube
type Cube struct {
	q, r, s int
}

// FloatCube implements floating point coordinates for hexes on a grid.
type FloatCube struct {
	q, r, s float64
}
