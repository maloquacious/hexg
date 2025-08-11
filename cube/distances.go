// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cube

/////////////////////////////////////////////////////////////////////////////
// distances
// * https://www.redblobgames.com/grids/hexagons/#distances

// function cube_subtract(a, b):
//   return Cube(a.q - b.q, a.r - b.r, a.s - b.s)

func (a Cube) Subtract(b Cube) Cube {
	return Cube{a.q - b.q, a.r - b.r, a.s - b.s}
}

// function cube_distance(a, b):
//   var vec = cube_subtract(a, b)
//   return max(abs(vec.q), abs(vec.r), abs(vec.s))
//// or: max(abs(a.q - b.q), abs(a.r - b.r), abs(a.s - b.s))

func (a Cube) Distance(b Cube) int {
	return max(abs(a.q-b.q), abs(a.r-b.r), abs(a.s-b.s))
}
