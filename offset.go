// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import "fmt"

type OffsetCoord struct {
	Col, Row int
}

func NewOffsetCoord(col, row int) OffsetCoord {
	return OffsetCoord{Col: col, Row: row}
}

func (oc OffsetCoord) String() string {
	return fmt.Sprintf("(%d,%d)", oc.Col, oc.Row)
}

func (oc OffsetCoord) ConciseString() string {
	return fmt.Sprintf("%+d%+d", oc.Col, oc.Row)
}

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

func (h OffsetCoord) QOffsetToCube(even bool) Hex {
	q := h.Col
	var r int
	if even {
		r = h.Row - (h.Col+1*(h.Col&1))/2
	} else {
		r = h.Row - (h.Col-1*(h.Col&1))/2
	}
	return Hex{
		q: q,
		r: r,
		s: -q - r,
	}
}

func (h Hex) CubeToROffset(even bool) OffsetCoord {
	var col int
	if even {
		col = h.q + (h.r+1*(h.r&1))/2
	} else {
		col = h.q + (h.r+1*(h.r&1))/2
	}
	return OffsetCoord{
		Col: col,
		Row: h.r,
	}
}

func (h OffsetCoord) ROffsetToCube(even bool) Hex {
	var q int
	if even {
		q = h.Col - (h.Row+1*(h.Row&1))/2
	} else {
		q = h.Col - (h.Row-1*(h.Row&1))/2
	}
	r := h.Row
	return Hex{
		q: q,
		r: r,
		s: -q - r,
	}
}
