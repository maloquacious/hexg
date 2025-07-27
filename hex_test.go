// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"testing"

	"github.com/maloquacious/hexg"
)

func TestHex_Distance(t *testing.T) {
	for _, tc := range []struct {
		id       int
		a, b     hexg.Hex
		distance int
	}{
		{id: 1, a: hexg.NewHex(0, 0, 0), b: hexg.NewHex(0, 0, 0), distance: 0},
		{id: 2, a: hexg.NewHex(0, 0, 0), b: hexg.NewHex(0, 1, -1), distance: 1},
		{id: 3, a: hexg.NewHex(0, 0, 0), b: hexg.NewHex(0, 2, -2), distance: 2},
		{id: 4, a: hexg.NewHex(0, 0, 0), b: hexg.NewHex(0, 3, -3), distance: 3},
		{id: 5, a: hexg.NewHex(0, 0, 0), b: hexg.NewHex(0, 4, -4), distance: 4},
		{id: 6, a: hexg.NewHex(0, 0, 0), b: hexg.NewHex(1, 3, -4), distance: 4},
		{id: 7, a: hexg.NewHex(0, 0, 0), b: hexg.NewHex(2, 2, -4), distance: 4},
		{id: 8, a: hexg.NewHex(0, 0, 0), b: hexg.NewHex(3, 1, -4), distance: 4},
		{id: 9, a: hexg.NewHex(-3, -1, 4), b: hexg.NewHex(4, -1, -3), distance: 7},
		{id: 10, a: hexg.NewHex(-1, -3, 4), b: hexg.NewHex(1, 3, -4), distance: 8},
	} {
		if distance := tc.a.Distance(tc.b); distance != tc.distance {
			t.Errorf("distance: from %q: to %q: got %d, want %d\n", tc.a.ConciseString(), tc.b.ConciseString(), distance, tc.distance)
		}
		if distance := tc.b.Distance(tc.a); distance != tc.distance {
			t.Errorf("distance: from %q: to %q: got %d, want %d\n", tc.b.ConciseString(), tc.a.ConciseString(), distance, tc.distance)
		}
	}
}

func TestHex_Neighbor(t *testing.T) {
	from := hexg.NewHex(0, 0, 0)
	for _, move := range []struct {
		id        int
		direction int
		expect    string
	}{
		// move one hex and then back
		{id: 1, direction: hexg.E, expect: "+1+0-1"}, {id: 2, direction: hexg.W, expect: "+0+0+0"},
		{id: 3, direction: hexg.NNE, expect: "+1-1+0"}, {id: 4, direction: hexg.SSW, expect: "+0+0+0"},
		{id: 5, direction: hexg.NNW, expect: "+0-1+1"}, {id: 6, direction: hexg.SSE, expect: "+0+0+0"},
		{id: 7, direction: hexg.W, expect: "-1+0+1"}, {id: 8, direction: hexg.E, expect: "+0+0+0"},
		{id: 9, direction: hexg.SSW, expect: "-1+1+0"}, {id: 10, direction: hexg.NNE, expect: "+0+0+0"},
		{id: 11, direction: hexg.SSE, expect: "+0+1-1"}, {id: 12, direction: hexg.NNW, expect: "+0+0+0"},
		// circle around
		{id: 13, direction: hexg.E, expect: "+1+0-1"},
		{id: 14, direction: hexg.NNE, expect: "+2-1-1"},
		{id: 15, direction: hexg.NNW, expect: "+2-2+0"},
		{id: 16, direction: hexg.W, expect: "+1-2+1"},
		{id: 17, direction: hexg.SSW, expect: "+0-1+1"},
		{id: 18, direction: hexg.SSE, expect: "+0+0+0"},
		// move three hexes in each direction
		{id: 19, direction: hexg.SSW, expect: "-1+1+0"}, {id: 20, direction: hexg.SSW, expect: "-2+2+0"}, {id: 21, direction: hexg.SSW, expect: "-3+3+0"},
		{id: 22, direction: hexg.NNW, expect: "-3+2+1"}, {id: 23, direction: hexg.NNW, expect: "-3+1+2"}, {id: 24, direction: hexg.NNW, expect: "-3+0+3"},
		{id: 25, direction: hexg.NNE, expect: "-2-1+3"}, {id: 26, direction: hexg.NNE, expect: "-1-2+3"}, {id: 27, direction: hexg.NNE, expect: "+0-3+3"},
		{id: 28, direction: hexg.E, expect: "+1-3+2"}, {id: 29, direction: hexg.E, expect: "+2-3+1"}, {id: 30, direction: hexg.E, expect: "+3-3+0"},
		{id: 31, direction: hexg.SSE, expect: "+3-2-1"}, {id: 32, direction: hexg.SSE, expect: "+3-1-2"}, {id: 33, direction: hexg.SSE, expect: "+3+0-3"},
		{id: 34, direction: hexg.SSW, expect: "+2+1-3"}, {id: 35, direction: hexg.SSW, expect: "+1+2-3"}, {id: 36, direction: hexg.SSW, expect: "+0+3-3"},
		{id: 37, direction: hexg.W, expect: "-1+3-2"}, {id: 37, direction: hexg.W, expect: "-2+3-1"}, {id: 38, direction: hexg.W, expect: "-3+3+0"},
	} {
		to := from.Neighbor(move.direction)
		if to.ConciseString() != move.expect {
			t.Fatalf("%d: from %q: to %d: got %q, want %q\n", move.id, from.ConciseString(), move.direction, to.ConciseString(), move.expect)
		}
		from = to
	}
}
