// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"testing"

	"github.com/maloquacious/hexg"
)

func TestEvenQ_Layout(t *testing.T) {
	l := hexg.NewVerticalEvenQLayout(hexg.NewPoint(1, 1), hexg.NewPoint(0, 0))

	if l.IsHorizontal() {
		t.Fatalf("even-q: isHorizontal: got %v, want %v\n", !l.IsHorizontal(), false)
	} else if !l.IsVertical() {
		t.Fatalf("even-q: isVertical: got %v, want %v\n", l.IsVertical(), true)
	} else if l.OffsetType() != hexg.EvenQ {
		t.Fatalf("even-q: offsetType: got %q, want %q\n", l.OffsetType(), hexg.EvenQ)
	}
}

func TestEvenQ_Compass(t *testing.T) {
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

func TestEvenQ_Neighbor(t *testing.T) {
	l := hexg.NewVerticalEvenQLayout(hexg.NewPoint(1, 1), hexg.NewPoint(0, 0))

	// even-q, even column
	for _, tc := range []struct {
		id      int
		from    hexg.Hex
		bearing string
		expect  string
	}{
		{id: 1, from: hexg.NewHex(0, 0, 0), bearing: "N", expect: "+0-1"},
		{id: 2, from: hexg.NewHex(0, 0, 0), bearing: "ENE", expect: "+1+0"},
		{id: 3, from: hexg.NewHex(0, 0, 0), bearing: "ESE", expect: "+1+1"},
		{id: 4, from: hexg.NewHex(0, 0, 0), bearing: "S", expect: "+0+1"},
		{id: 5, from: hexg.NewHex(0, 0, 0), bearing: "WSW", expect: "-1+1"},
		{id: 6, from: hexg.NewHex(0, 0, 0), bearing: "WNW", expect: "-1+0"},
	} {
		direction := hexg.BearingToDirection(tc.bearing)
		neighbor := tc.from.Neighbor(direction)
		to := l.HexToOffsetCoord(neighbor)
		got := to.ConciseString()
		if got != tc.expect {
			t.Errorf("even-col: %d: from %q: %-3s: %q: got %q, want %q\n", tc.id, tc.from.ConciseString(), tc.bearing, neighbor.ConciseString(), got, tc.expect)
		}
	}

	// even-q, odd column
	for _, tc := range []struct {
		id      int
		from    hexg.Hex
		bearing string
		expect  string
	}{
		{id: 1, from: hexg.NewHex(1, 0, -1), bearing: "N", expect: "+1+0"},
		{id: 2, from: hexg.NewHex(1, 0, -1), bearing: "ENE", expect: "+2+0"},
		{id: 3, from: hexg.NewHex(1, 0, -1), bearing: "ESE", expect: "+2+1"},
		{id: 4, from: hexg.NewHex(1, 0, -1), bearing: "S", expect: "+1+2"},
		{id: 5, from: hexg.NewHex(1, 0, -1), bearing: "WSW", expect: "+0+1"},
		{id: 6, from: hexg.NewHex(1, 0, -1), bearing: "WNW", expect: "+0+0"},
	} {
		direction := hexg.BearingToDirection(tc.bearing)
		neighbor := tc.from.Neighbor(direction)
		to := l.HexToOffsetCoord(neighbor)
		got := to.ConciseString()
		if got != tc.expect {
			t.Errorf("odd-col : %d: from %q: %-3s: %q: got %q, want %q\n", tc.id, tc.from.ConciseString(), tc.bearing, neighbor.ConciseString(), got, tc.expect)
		}
	}
}

func TestEvenQ_OffsetToHex(t *testing.T) {
	l := hexg.NewVerticalEvenQLayout(hexg.NewPoint(1, 1), hexg.NewPoint(0, 0))

	for _, tc := range []struct {
		id       int
		col, row int
		expect   string
	}{
		{id: 1, col: 0, row: 0, expect: "+0+0+0"},
		// neighboring ring of hexes
		{id: 2, col: 0, row: -1, expect: "+0-1+1"},
		{id: 3, col: 1, row: 0, expect: "+1-1+0"},
		{id: 4, col: 1, row: 1, expect: "+1+0-1"},
		{id: 5, col: 0, row: 1, expect: "+0+1-1"},
		{id: 6, col: -1, row: 1, expect: "-1+1+0"},
		{id: 7, col: -1, row: 0, expect: "-1+0+1"},
		// 2, 0 and down two
		{id: 8, col: 2, row: 0, expect: "+2-1-1"},
		{id: 9, col: 2, row: 1, expect: "+2+0-2"},
		{id: 10, col: 2, row: 2, expect: "+2+1-3"},
		// -2, 0 and up two
		{id: 11, col: -2, row: 0, expect: "-2+1+1"},
		{id: 12, col: -2, row: -1, expect: "-2+0+2"},
		{id: 13, col: -2, row: -2, expect: "-2-1+3"},
		// 3, 0 and down three
		{id: 14, col: 3, row: 0, expect: "+3-2-1"},
		{id: 15, col: 3, row: 1, expect: "+3-1-2"},
		{id: 16, col: 3, row: 2, expect: "+3+0-3"},
		{id: 17, col: 3, row: 3, expect: "+3+1-4"},
		// -3, 0 and up three
		{id: 18, col: -3, row: 0, expect: "-3+1+2"},
		{id: 19, col: -3, row: -1, expect: "-3+0+3"},
		{id: 20, col: -3, row: -2, expect: "-3-1+4"},
		{id: 21, col: -3, row: -3, expect: "-3-2+5"},
	} {
		hex := l.OffsetColRowToHex(tc.col, tc.row)
		got := hex.ConciseString()
		if got != tc.expect {
			t.Errorf("%d: col %3d, row %3d: got %q, want %q\n", tc.id, tc.col, tc.row, got, tc.expect)
		}
	}
}
