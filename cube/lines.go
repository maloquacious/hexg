// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cube

/////////////////////////////////////////////////////////////////////////////
// line drawing
// * https://www.redblobgames.com/grids/hexagons/#line-drawing

// function lerp(a, b, t): # for floats
//    return a + (b - a) * t
// implemented in math.go

// function cube_lerp(a, b, t): # for hexes
//    return Cube(lerp(a.q, b.q, t),
//                lerp(a.r, b.r, t),
//                lerp(a.s, b.s, t))

func (a Cube) Lerp(b Cube, t float64) FloatCube {
	return FloatCube{
		q: lerp(a.q, b.q, t),
		r: lerp(a.r, b.r, t),
		s: lerp(a.s, b.s, t),
	}
}

func (a FloatCube) Lerp(b FloatCube, t float64) FloatCube {
	return FloatCube{
		q: lerp(a.q, b.q, t),
		r: lerp(a.r, b.r, t),
		s: lerp(a.s, b.s, t),
	}
}

// function cube_linedraw(a, b):
//    var N = cube_distance(a, b)
//    var results = []
//    for each 0 ≤ i ≤ N:
//        results.append(cube_round(cube_lerp(a, b, 1.0/N * i)))
//    return results

func (a Cube) Linedraw(b Cube) []Cube {
	N := a.Distance(b)
	if N == 0 {
		return []Cube{}
	}
	var results []Cube
	for i := 0; i <= N; i++ {
		results = append(results, a.Lerp(b, 1.0/(float64(N*i))).Round())
	}
	return results
}

func (a Cube) LinedrawWithNudge(b Cube) []Cube {
	N := a.Distance(b)
	if N == 0 {
		return []Cube{}
	}
	step := 1.0 / float64(N)
	var results []Cube
	h_nudge := FloatCube{q: float64(a.q) + 1e-6, r: float64(a.r) + 2e-6, s: float64(a.s) - 3e-6}
	b_nudge := FloatCube{q: float64(b.q) + 1e-6, r: float64(b.r) + 2e-6, s: float64(b.s) - 3e-6}
	for i := 0; i <= N; i++ {
		results = append(results, a.Lerp(b, 1.0/(float64(N*i))).Round())
		results = append(results, h_nudge.Lerp(b_nudge, step*float64(i)).Round())

	}
	return results
}
