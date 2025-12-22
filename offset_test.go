// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"testing"

	"github.com/maloquacious/hexg"
)

func TestNewOffsetCoord(t *testing.T) {
	oc := hexg.NewOffsetCoord(3, 5)
	if oc.Col != 3 || oc.Row != 5 {
		t.Errorf("NewOffsetCoord(3, 5) = (%d, %d), want (3, 5)", oc.Col, oc.Row)
	}
}

func TestOffsetCoord_String(t *testing.T) {
	for _, tc := range []struct {
		col, row int
		want     string
	}{
		{0, 0, "(0,0)"},
		{3, 5, "(3,5)"},
		{-2, 4, "(-2,4)"},
	} {
		oc := hexg.NewOffsetCoord(tc.col, tc.row)
		if got := oc.String(); got != tc.want {
			t.Errorf("OffsetCoord{%d,%d}.String() = %q, want %q", tc.col, tc.row, got, tc.want)
		}
	}
}

func TestOffsetCoord_ConciseString(t *testing.T) {
	for _, tc := range []struct {
		col, row int
		want     string
	}{
		{0, 0, "+0+0"},
		{3, 5, "+3+5"},
		{-2, 4, "-2+4"},
	} {
		oc := hexg.NewOffsetCoord(tc.col, tc.row)
		if got := oc.ConciseString(); got != tc.want {
			t.Errorf("OffsetCoord{%d,%d}.ConciseString() = %q, want %q", tc.col, tc.row, got, tc.want)
		}
	}
}

func TestCubeToQOffset_OddQ(t *testing.T) {
	for _, tc := range []struct {
		q, r    int
		wantCol int
		wantRow int
	}{
		{0, 0, 0, 0},
		{1, 0, 1, 0},
		{1, 1, 1, 1},
		{2, 0, 2, 1},
		{-1, 0, -1, -1},
	} {
		hex := hexg.NewHex(tc.q, tc.r)
		oc := hex.CubeToQOffset(false)
		if oc.Col != tc.wantCol || oc.Row != tc.wantRow {
			t.Errorf("CubeToQOffset(odd) for (%d,%d): got (%d,%d), want (%d,%d)",
				tc.q, tc.r, oc.Col, oc.Row, tc.wantCol, tc.wantRow)
		}
	}
}

func TestCubeToQOffset_EvenQ(t *testing.T) {
	for _, tc := range []struct {
		q, r    int
		wantCol int
		wantRow int
	}{
		{0, 0, 0, 0},
		{1, 0, 1, 1},
		{2, 0, 2, 1},
	} {
		hex := hexg.NewHex(tc.q, tc.r)
		oc := hex.CubeToQOffset(true)
		if oc.Col != tc.wantCol || oc.Row != tc.wantRow {
			t.Errorf("CubeToQOffset(even) for (%d,%d): got (%d,%d), want (%d,%d)",
				tc.q, tc.r, oc.Col, oc.Row, tc.wantCol, tc.wantRow)
		}
	}
}

func TestQOffsetToCube_RoundTrip(t *testing.T) {
	for q := -3; q <= 3; q++ {
		for r := -3; r <= 3; r++ {
			hex := hexg.NewHex(q, r)

			ocOdd := hex.CubeToQOffset(false)
			backOdd := ocOdd.QOffsetToCube(false)
			if !hex.Equals(backOdd) {
				t.Errorf("odd-q round trip failed: %v -> %v -> %v", hex, ocOdd, backOdd)
			}

			ocEven := hex.CubeToQOffset(true)
			backEven := ocEven.QOffsetToCube(true)
			if !hex.Equals(backEven) {
				t.Errorf("even-q round trip failed: %v -> %v -> %v", hex, ocEven, backEven)
			}
		}
	}
}

func TestCubeToROffset_OddR(t *testing.T) {
	for _, tc := range []struct {
		q, r    int
		wantCol int
		wantRow int
	}{
		{0, 0, 0, 0},
		{0, 1, 0, 1},
		{1, 1, 1, 1},
		{0, 2, 1, 2},
	} {
		hex := hexg.NewHex(tc.q, tc.r)
		oc := hex.CubeToROffset(false)
		if oc.Col != tc.wantCol || oc.Row != tc.wantRow {
			t.Errorf("CubeToROffset(odd) for (%d,%d): got (%d,%d), want (%d,%d)",
				tc.q, tc.r, oc.Col, oc.Row, tc.wantCol, tc.wantRow)
		}
	}
}

func TestCubeToROffset_EvenR(t *testing.T) {
	for _, tc := range []struct {
		q, r    int
		wantCol int
		wantRow int
	}{
		{0, 0, 0, 0},
		{0, 1, 1, 1},
		{0, 2, 1, 2},
	} {
		hex := hexg.NewHex(tc.q, tc.r)
		oc := hex.CubeToROffset(true)
		if oc.Col != tc.wantCol || oc.Row != tc.wantRow {
			t.Errorf("CubeToROffset(even) for (%d,%d): got (%d,%d), want (%d,%d)",
				tc.q, tc.r, oc.Col, oc.Row, tc.wantCol, tc.wantRow)
		}
	}
}

func TestROffsetToCube_RoundTrip(t *testing.T) {
	for q := -3; q <= 3; q++ {
		for r := -3; r <= 3; r++ {
			hex := hexg.NewHex(q, r)

			ocOdd := hex.CubeToROffset(false)
			backOdd := ocOdd.ROffsetToCube(false)
			if !hex.Equals(backOdd) {
				t.Errorf("odd-r round trip failed: %v -> %v -> %v", hex, ocOdd, backOdd)
			}

			ocEven := hex.CubeToROffset(true)
			backEven := ocEven.ROffsetToCube(true)
			if !hex.Equals(backEven) {
				t.Errorf("even-r round trip failed: %v -> %v -> %v", hex, ocEven, backEven)
			}
		}
	}
}
