// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import "fmt"

// OffsetCoord represents offset coordinates (col, row) for hex grids.
type OffsetCoord struct {
	Col, Row int
}

// NewOffsetCoord returns an OffsetCoord initialized with the given column and row.
func NewOffsetCoord(col, row int) OffsetCoord {
	return OffsetCoord{Col: col, Row: row}
}

// String returns the coordinates formatted as "(col,row)".
func (oc OffsetCoord) String() string {
	return fmt.Sprintf("(%d,%d)", oc.Col, oc.Row)
}

// ConciseString returns the coordinates formatted as "+col+row".
func (oc OffsetCoord) ConciseString() string {
	return fmt.Sprintf("%+d%+d", oc.Col, oc.Row)
}

// CubeToQOffset converts cube coordinates to q-offset coordinates.
// If even is true, uses even-q offset; otherwise uses odd-q offset.
func (h Hex) CubeToQOffset(even bool) OffsetCoord {
	col := h.q
	var row int
	if even {
		row = h.r + (h.q+1*(h.q&1))/2
	} else {
		row = h.r + (h.q-1*(h.q&1))/2
	}
	return OffsetCoord{Col: col, Row: row}
}

// QOffsetToCube converts q-offset coordinates to cube coordinates.
// If even is true, uses even-q offset; otherwise uses odd-q offset.
func (oc OffsetCoord) QOffsetToCube(even bool) Hex {
	q := oc.Col
	var r int
	if even {
		r = oc.Row - (oc.Col+1*(oc.Col&1))/2
	} else {
		r = oc.Row - (oc.Col-1*(oc.Col&1))/2
	}
	return Hex{
		q: q,
		r: r,
		s: -q - r,
	}
}

// CubeToROffset converts cube coordinates to r-offset coordinates.
// If even is true, uses even-r offset; otherwise uses odd-r offset.
func (h Hex) CubeToROffset(even bool) OffsetCoord {
	var col int
	if even {
		col = h.q + (h.r+1*(h.r&1))/2
	} else {
		col = h.q + (h.r-1*(h.r&1))/2
	}
	return OffsetCoord{
		Col: col,
		Row: h.r,
	}
}

// ROffsetToCube converts r-offset coordinates to cube coordinates.
// If even is true, uses even-r offset; otherwise uses odd-r offset.
func (oc OffsetCoord) ROffsetToCube(even bool) Hex {
	var q int
	if even {
		q = oc.Col - (oc.Row+1*(oc.Row&1))/2
	} else {
		q = oc.Col - (oc.Row-1*(oc.Row&1))/2
	}
	r := oc.Row
	return Hex{
		q: q,
		r: r,
		s: -q - r,
	}
}
