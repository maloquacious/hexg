// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cube

/////////////////////////////////////////////////////////////////////////////
// rounding to nearest hex
// * https://www.redblobgames.com/grids/hexagons/#rounding

// function cube_round(frac):
//    var q = round(frac.q)
//    var r = round(frac.r)
//    var s = round(frac.s)
//
//    var q_diff = abs(q - frac.q)
//    var r_diff = abs(r - frac.r)
//    var s_diff = abs(s - frac.s)
//
//    if q_diff > r_diff and q_diff > s_diff:
//        q = -r-s
//    else if r_diff > s_diff:
//        r = -q-s
//    else:
//        s = -q-r
//
//    return Cube(q, r, s)

func (c Cube) Round() Cube {
	return c
}

func (frac FloatCube) Round() Cube {
	q, r, s := round(frac.q), round(frac.r), round(frac.s)

	q_diff, r_diff, s_diff := abs(float64(q)-frac.q), abs(float64(r)-frac.r), abs(float64(s)-frac.s)

	if q_diff > r_diff && q_diff > s_diff {
		q = -r - s
	} else if r_diff > s_diff {
		r = -q - s
	} else {
		s = -q - r
	}
	return Cube{q: q, r: r, s: s}
}
