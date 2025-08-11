// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cube

import "fmt"

/////////////////////////////////////////////////////////////////////////////
// rings
// * https://www.redblobgames.com/grids/hexagons/#rings

// single rings

// function cube_scale(hex, factor):
//    return Cube(hex.q * factor, hex.r * factor, hex.s * factor)

func (hex Cube) Scale(factor int) Cube {
	return Cube{q: hex.q * factor, r: hex.r * factor, s: hex.s * factor}
}

func (hex FloatCube) Scale(factor float64) FloatCube {
	return FloatCube{q: hex.q * factor, r: hex.r * factor, s: hex.s * factor}
}

// function cube_ring(center, radius):
//    var results = []
//    # this code doesn't work for radius == 0; can you see why?
//    var hex = cube_add(center,
//                        cube_scale(cube_direction(4), radius))
//    for each 0 ≤ i < 6:
//        for each 0 ≤ j < radius:
//            results.append(hex)
//            hex = cube_neighbor(hex, i)
//    return results

func (center Cube) Ring(radius int) []Cube {
	if radius < 0 {
		panic(fmt.Sprintf("assert(radius != %d)", radius))
	} else if radius == 0 {
		return []Cube{center}
	}
	hex := center.Add(Direction(4).Scale(radius))
	var results []Cube
	for i := 0; i < 6; i++ {
		for j := 0; j < radius; j++ {
			results = append(results, hex)
		}
	}
	return results
}

// spiral rings

// function cube_spiral(center, radius):
//    var results = list(center)
//    for each 1 ≤ k ≤ radius:
//        results = list_append(results, cube_ring(center, k))
//    return results

func (center Cube) Spiral(radius int) []Cube {
	if radius < 0 {
		panic(fmt.Sprintf("assert(radius != %d)", radius))
	}
	results := []Cube{center}
	for k := 1; k <= radius; k++ {
		results = append(results, center.Ring(k)...)
	}
	return results
}

// spiral coordinates
// will not implement
