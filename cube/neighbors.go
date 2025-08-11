// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cube

/////////////////////////////////////////////////////////////////////////////
// neighbors
// https://www.redblobgames.com/grids/hexagons/#neighbors

// there is no "up" or "down" until we have a layout, so this is abstract.
// We assign 0 and then go counter-clockwise for the remaining.
var cube_direction_vectors = [6]Cube{
	/* direction 0 */ {q: +1, r: 0, s: -1},
	/* direction 1 */ {q: +1, r: -1, s: 0},
	/* direction 2 */ {q: 0, r: -1, s: +1},
	/* direction 3 */ {q: -1, r: 0, s: +1},
	/* direction 4 */ {q: -1, r: +1, s: 0},
	/* direction 5 */ {q: 0, r: +1, s: -1},
}

// function cube_direction(direction):
//   return cube_direction_vectors[direction]

func Direction(direction int) Cube {
	return cube_direction_vectors[(6+(direction%6))%6]
}

// function cube_add(hex, vec):
//   return Cube(hex.q + vec.q, hex.r + vec.r, hex.s + vec.s)

func (hex Cube) Add(vec Cube) Cube {
	return Cube{q: hex.q + vec.q, r: hex.r + vec.r, s: hex.s + vec.s}
}

func (hex FloatCube) Add(vec FloatCube) FloatCube {
	return FloatCube{q: hex.q + vec.q, r: hex.r + vec.r, s: hex.s + vec.s}
}

// function cube_neighbor(cube, direction):
//   return cube_add(cube, cube_direction(direction))

func (cube Cube) Neighbor(direction int) Cube {
	return cube.Add(Direction(direction))
}

// diagonals
// * https://www.redblobgames.com/grids/hexagons/#neighbors-diagonal

var cube_diagonal_vectors = [6]Cube{
	/* direction 0 */ {q: +2, r: -1, s: -1},
	/* direction 0 */ {q: +1, r: -2, s: +1},
	/* direction 0 */ {q: -1, r: -1, s: +2},
	/* direction 0 */ {q: -2, r: +1, s: +1},
	/* direction 0 */ {q: -1, r: +2, s: -1},
	/* direction 0 */ {q: +1, r: +1, s: -2},
}

func DiagonalDirection(direction int) Cube {
	return cube_diagonal_vectors[(6+(direction%6))%6]
}

//function cube_diagonal_neighbor(cube, direction):
//    return cube_add(cube, cube_diagonal_vectors[direction])

func (cube Cube) DiagonalNeighbor(direction int) Cube {
	return cube.Add(DiagonalDirection(direction))
}
