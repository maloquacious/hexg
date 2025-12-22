// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import "math"

// NewFractionalHex returns a FractionalHex initialized from cube coordinates.
func NewFractionalHex(q, r, s float64) FractionalHex {
	return FractionalHex{q: q, r: r, s: s}
}

// FractionalHex implements a hex with floating point coordinates.
// Used for linear interpolation and rounding.
type FractionalHex struct {
	q, r, s float64
}

// Round turns a fractional hex coordinate into the nearest integer
// hex coordinate.
func (h FractionalHex) Round() Hex {
	q := int(math.Round(h.q))
	r := int(math.Round(h.r))
	s := int(math.Round(h.s))
	qDiff := math.Abs(float64(q) - h.q)
	rDiff := math.Abs(float64(r) - h.r)
	sDiff := math.Abs(float64(s) - h.s)
	if qDiff > rDiff && qDiff > sDiff {
		q = -r - s
	} else if rDiff > sDiff {
		r = -q - s
	} else {
		s = -q - r
	}
	return Hex{q: q, r: r, s: s}
}

// Lerp returns the linear interpolation of points on the line between two hexes
func (h FractionalHex) Lerp(b FractionalHex, t float64) FractionalHex {
	return FractionalHex{
		q: lerp((h.q), (b.q), t),
		r: lerp((h.r), (b.r), t),
		s: lerp((h.s), (b.s), t),
	}
}

// Scale returns a FractionalHex scaled by a float64 factor.
func (h FractionalHex) Scale(factor float64) FractionalHex {
	return FractionalHex{q: h.q * factor, r: h.r * factor, s: h.s * factor}
}

// ReflectQ reflects the fractional hex across the Q axis.
func (h FractionalHex) ReflectQ() FractionalHex {
	return FractionalHex{q: h.q, r: h.s, s: h.r}
}

// ReflectR reflects the fractional hex across the R axis.
func (h FractionalHex) ReflectR() FractionalHex {
	return FractionalHex{q: h.s, r: h.r, s: h.q}
}

// ReflectS reflects the fractional hex across the S axis.
func (h FractionalHex) ReflectS() FractionalHex {
	return FractionalHex{q: h.r, r: h.q, s: h.s}
}
