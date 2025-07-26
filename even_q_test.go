// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"testing"

	"github.com/maloquacious/hexg"
)

func TestEvenQLayout(t *testing.T) {
	l := hexg.NewVerticalEvenQLayout(hexg.NewPoint(1, 1), hexg.NewPoint(0, 0))

	if l.IsHorizontal() {
		t.Fatalf("even-q: isHorizontal: got %v, want %v\n", !l.IsHorizontal(), false)
	} else if !l.IsVertical() {
		t.Fatalf("even-q: isVertical: got %v, want %v\n", l.IsVertical(), true)
	} else if l.OffsetType() != hexg.EvenQ {
		t.Fatalf("even-q: offsetType: got %q, want %q\n", l.OffsetType(), hexg.EvenQ)
	}
}

func TestEvenQCompass(t *testing.T) {
	l := hexg.NewVerticalEvenQLayout(hexg.NewPoint(1, 1), hexg.NewPoint(0, 0))

	for _, tc := range []struct {
		id        int
		bearing   string
		direction int
	}{
		{id: 1, bearing: "N", direction: hexg.N},
		{id: 2, bearing: "ENE", direction: hexg.ENE},
		{id: 3, bearing: "ESE", direction: hexg.ESE},
		{id: 4, bearing: "S", direction: hexg.S},
		{id: 5, bearing: "WSW", direction: hexg.WSW},
		{id: 6, bearing: "WNW", direction: hexg.WNW},
	} {
		bearing := l.DirectionToBearing(tc.direction)
		if bearing != tc.bearing {
			t.Errorf("%d: even-q: direction %d: bearing got %q, want %q\n", tc.id, tc.direction, bearing, tc.bearing)
		}
	}
}

func TestEvenQNeighbor(t *testing.T) {
	l := hexg.NewVerticalEvenQLayout(hexg.NewPoint(1, 1), hexg.NewPoint(0, 0))

	// even-q, even column
	for _, tc := range []struct {
		id       int
		q, r, s  int
		bearing  string
		col, row int
	}{
		{id: 1, q: 0, r: 0, s: 0, bearing: "N", col: 0, row: -1},
		{id: 2, q: 0, r: 0, s: 0, bearing: "ENE", col: 1, row: 0},
		{id: 3, q: 0, r: 0, s: 0, bearing: "ESE", col: 1, row: 1},
		{id: 4, q: 0, r: 0, s: 0, bearing: "S", col: 0, row: 1},
		{id: 5, q: 0, r: 0, s: 0, bearing: "WSW", col: -1, row: 1},
		{id: 6, q: 0, r: 0, s: 0, bearing: "WNW", col: -1, row: 0},
	} {
		from := hexg.NewHex(tc.q, tc.r, tc.s)
		direction := hexg.BearingToDirection(tc.bearing)
		neighbor := from.Neighbor(direction)
		to := l.HexToOffsetCoord(neighbor)
		expect := hexg.OffsetCoordFromColRow(tc.col, tc.row)
		if to.String() != expect.String() {
			t.Errorf("%d: even-col: from %q: %-3s: %q: got %q, want %q\n", tc.id, from.ConciseString(), tc.bearing, neighbor.ConciseString(), to.String(), expect.String())
		}
	}

	// even-q, odd column
	for _, tc := range []struct {
		id       int
		q, r, s  int
		bearing  string
		col, row int
	}{
		{id: 1, q: 1, r: 0, s: -1, bearing: "N", col: 1, row: 0},
		{id: 2, q: 1, r: 0, s: -1, bearing: "ENE", col: 2, row: 0},
		{id: 3, q: 1, r: 0, s: -1, bearing: "ESE", col: 2, row: 1},
		{id: 4, q: 1, r: 0, s: -1, bearing: "S", col: 1, row: 2},
		{id: 5, q: 1, r: 0, s: -1, bearing: "WSW", col: 0, row: 1},
		{id: 6, q: 1, r: 0, s: -1, bearing: "WNW", col: 0, row: 0},
	} {
		from := hexg.NewHex(tc.q, tc.r, tc.s)
		direction := hexg.BearingToDirection(tc.bearing)
		neighbor := from.Neighbor(direction)
		to := l.HexToOffsetCoord(neighbor)
		expect := hexg.OffsetCoordFromColRow(tc.col, tc.row)
		if to.String() != expect.String() {
			t.Errorf("%d: odd-col:  from %q: %-3s: %q: got %q, want %q\n", tc.id, from.ConciseString(), tc.bearing, neighbor.ConciseString(), to.String(), expect.String())
		}
	}
}

func TestEvenQOffsetToHex(t *testing.T) {
	l := hexg.NewVerticalEvenQLayout(hexg.NewPoint(1, 1), hexg.NewPoint(0, 0))

	for _, tc := range []struct {
		id       int
		col, row int
		q, r, s  int
	}{
		{id: 1, col: 0, row: 0, q: 0, r: 0, s: 0},
		// neighboring ring of hexes
		{id: 2, col: 0, row: -1, q: 0, r: -1, s: 1},
		{id: 3, col: 1, row: 0, q: 1, r: -1, s: 0},
		{id: 4, col: 1, row: 1, q: 1, r: 0, s: -1},
		{id: 5, col: 0, row: 1, q: 0, r: 1, s: -1},
		{id: 6, col: -1, row: 1, q: -1, r: 1, s: 0},
		{id: 7, col: -1, row: 0, q: -1, r: 0, s: 1},
		// 2, 0 and down two
		{id: 8, col: 2, row: 0, q: 2, r: -1, s: -1},
		{id: 9, col: 2, row: 1, q: 2, r: 0, s: -2},
		{id: 10, col: 2, row: 2, q: 2, r: 1, s: -3},
		// -2, 0 and up two
		{id: 11, col: -2, row: 0, q: -2, r: 1, s: 1},
		{id: 12, col: -2, row: -1, q: -2, r: 0, s: 2},
		{id: 13, col: -2, row: -2, q: -2, r: -1, s: 3},
		// 3, 0 and down three
		{id: 14, col: 3, row: 0, q: 3, r: -2, s: -1},
		{id: 15, col: 3, row: 1, q: 3, r: -1, s: -2},
		{id: 16, col: 3, row: 2, q: 3, r: 0, s: -3},
		{id: 17, col: 3, row: 3, q: 3, r: 1, s: -4},
		// -3, 0 and up three
		{id: 18, col: -3, row: 0, q: -3, r: 1, s: 2},
		{id: 19, col: -3, row: -1, q: -3, r: 0, s: 3},
		{id: 20, col: -3, row: -2, q: -3, r: -1, s: 4},
		{id: 21, col: -3, row: -3, q: -3, r: -2, s: 5},
	} {
		oc := hexg.NewOffsetCoord(tc.col, tc.row)
		got := l.OffsetColRowToHex(tc.col, tc.row)
		expect := hexg.NewHex(tc.q, tc.r, tc.s)
		if got.String() != expect.String() {
			t.Errorf("%d: oc %q: got %q, want %q\n", tc.id, oc.ConciseString(), got.ConciseString(), expect.ConciseString())
		}
	}
}
