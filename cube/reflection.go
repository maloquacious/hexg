// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cube

/////////////////////////////////////////////////////////////////////////////
// reflection
// * https://www.redblobgames.com/grids/hexagons/#reflection

// function reflectQ(h) { return Cube(h.q, h.s, h.r); }

func (h Cube) ReflectQ() Cube {
	return Cube{q: h.q, r: h.s, s: h.r}
}

func (h FloatCube) ReflectQ() FloatCube {
	return FloatCube{q: h.q, r: h.s, s: h.r}
}

// function reflectR(h) { return Cube(h.s, h.r, h.q); }

func (h Cube) ReflectR() Cube {
	return Cube{q: h.s, r: h.r, s: h.q}
}

func (h FloatCube) ReflectR() FloatCube {
	return FloatCube{q: h.s, r: h.r, s: h.q}
}

// function reflectS(h) { return Cube(h.r, h.q, h.s); }

func (h Cube) ReflectS() Cube {
	return Cube{q: h.r, r: h.q, s: h.s}
}

func (h FloatCube) ReflectS() FloatCube {
	return FloatCube{q: h.r, r: h.q, s: h.s}
}
