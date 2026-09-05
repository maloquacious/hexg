// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"slices"
	"testing"

	"github.com/maloquacious/hexg"
)

func TestHex_Stringers(t *testing.T) {
	for _, tc := range []struct {
		id    int
		a     hexg.Hex
		s, cs string
	}{
		{id: 1, a: hexg.NewHex(0, 0), s: "(0, 0, 0)", cs: "+0+0+0"},
		{id: 2, a: hexg.NewHex(0, 1), s: "(0, 1, -1)", cs: "+0+1-1"},
		{id: 3, a: hexg.NewHex(0, 2), s: "(0, 2, -2)", cs: "+0+2-2"},
		{id: 4, a: hexg.NewHex(0, 3), s: "(0, 3, -3)", cs: "+0+3-3"},
		{id: 5, a: hexg.NewHex(0, 4), s: "(0, 4, -4)", cs: "+0+4-4"},
		{id: 6, a: hexg.NewHex(1, 3), s: "(1, 3, -4)", cs: "+1+3-4"},
		{id: 7, a: hexg.NewHex(2, 2), s: "(2, 2, -4)", cs: "+2+2-4"},
		{id: 8, a: hexg.NewHex(-2, 1), s: "(-2, 1, 1)", cs: "-2+1+1"},
		{id: 9, a: hexg.NewHex(-3, -1), s: "(-3, -1, 4)", cs: "-3-1+4"},
		{id: 10, a: hexg.NewHex(-1, -3), s: "(-1, -3, 4)", cs: "-1-3+4"},
	} {
		if want, got := tc.s, tc.a.String(); want != got {
			t.Errorf("%d: string: want %q, got %q\n", tc.id, want, got)
		}
		if want, got := tc.cs, tc.a.ConciseString(); want != got {
			t.Errorf("%d: concise: want %q, got %q\n", tc.id, want, got)
		}
	}
}

func TestHex_Distance(t *testing.T) {
	for _, tc := range []struct {
		id       int
		a, b     hexg.Hex
		distance int
	}{
		{id: 1, a: hexg.NewHex(0, 0), b: hexg.NewHex(0, 0), distance: 0},
		{id: 2, a: hexg.NewHex(0, 0), b: hexg.NewHex(0, 1), distance: 1},
		{id: 3, a: hexg.NewHex(0, 0), b: hexg.NewHex(0, 2), distance: 2},
		{id: 4, a: hexg.NewHex(0, 0), b: hexg.NewHex(0, 3), distance: 3},
		{id: 5, a: hexg.NewHex(0, 0), b: hexg.NewHex(0, 4), distance: 4},
		{id: 6, a: hexg.NewHex(0, 0), b: hexg.NewHex(1, 3), distance: 4},
		{id: 7, a: hexg.NewHex(0, 0), b: hexg.NewHex(2, 2), distance: 4},
		{id: 8, a: hexg.NewHex(0, 0), b: hexg.NewHex(3, 1), distance: 4},
		{id: 9, a: hexg.NewHex(-3, -1), b: hexg.NewHex(4, -1), distance: 7},
		{id: 10, a: hexg.NewHex(-1, -3), b: hexg.NewHex(1, 3), distance: 8},
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
	from := hexg.NewHex(0, 0)
	for _, move := range []struct {
		id        int
		direction int
		expect    string
	}{
		// move one hex and then back
		{id: 1, direction: 0, expect: "+1+0-1"}, {id: 2, direction: 3, expect: "+0+0+0"},
		{id: 3, direction: 1, expect: "+1-1+0"}, {id: 4, direction: 4, expect: "+0+0+0"},
		{id: 5, direction: 2, expect: "+0-1+1"}, {id: 6, direction: 5, expect: "+0+0+0"},
		{id: 7, direction: 3, expect: "-1+0+1"}, {id: 8, direction: 0, expect: "+0+0+0"},
		{id: 9, direction: 4, expect: "-1+1+0"}, {id: 10, direction: 1, expect: "+0+0+0"},
		{id: 11, direction: 5, expect: "+0+1-1"}, {id: 12, direction: 2, expect: "+0+0+0"},
		// circle around
		{id: 13, direction: 0, expect: "+1+0-1"},
		{id: 14, direction: 1, expect: "+2-1-1"},
		{id: 15, direction: 2, expect: "+2-2+0"},
		{id: 16, direction: 3, expect: "+1-2+1"},
		{id: 17, direction: 4, expect: "+0-1+1"},
		{id: 18, direction: 5, expect: "+0+0+0"},
		// move three hexes in each direction
		{id: 19, direction: 4, expect: "-1+1+0"}, {id: 20, direction: 4, expect: "-2+2+0"}, {id: 21, direction: 4, expect: "-3+3+0"},
		{id: 22, direction: 2, expect: "-3+2+1"}, {id: 23, direction: 2, expect: "-3+1+2"}, {id: 24, direction: 2, expect: "-3+0+3"},
		{id: 25, direction: 1, expect: "-2-1+3"}, {id: 26, direction: 1, expect: "-1-2+3"}, {id: 27, direction: 1, expect: "+0-3+3"},
		{id: 28, direction: 0, expect: "+1-3+2"}, {id: 29, direction: 0, expect: "+2-3+1"}, {id: 30, direction: 0, expect: "+3-3+0"},
		{id: 31, direction: 5, expect: "+3-2-1"}, {id: 32, direction: 5, expect: "+3-1-2"}, {id: 33, direction: 5, expect: "+3+0-3"},
		{id: 34, direction: 4, expect: "+2+1-3"}, {id: 35, direction: 4, expect: "+1+2-3"}, {id: 36, direction: 4, expect: "+0+3-3"},
		{id: 37, direction: 3, expect: "-1+3-2"}, {id: 38, direction: 3, expect: "-2+3-1"}, {id: 39, direction: 3, expect: "-3+3+0"},
	} {
		to := from.Neighbor(move.direction)
		if to.ConciseString() != move.expect {
			t.Fatalf("%d: from %q: to %d: got %q, want %q\n", move.id, from.ConciseString(), move.direction, to.ConciseString(), move.expect)
		}
		from = to
	}
}

func TestHex_Compare(t *testing.T) {
	for _, tc := range []struct {
		id   int
		a, b hexg.Hex
		want int
	}{
		// identical hexes
		{id: 1, a: hexg.NewHex(0, 0), b: hexg.NewHex(0, 0), want: 0},
		{id: 2, a: hexg.NewHex(-3, 7), b: hexg.NewHex(-3, 7), want: 0},
		// same r, ordered by q
		{id: 3, a: hexg.NewHex(0, 0), b: hexg.NewHex(1, 0), want: -1},
		{id: 4, a: hexg.NewHex(1, 0), b: hexg.NewHex(0, 0), want: +1},
		{id: 5, a: hexg.NewHex(-5, 2), b: hexg.NewHex(-4, 2), want: -1},
		// r dominates q: a larger q still sorts first when r is smaller.
		// These cases fail under a q-then-r ordering.
		{id: 6, a: hexg.NewHex(9, 0), b: hexg.NewHex(-9, 1), want: -1},
		{id: 7, a: hexg.NewHex(-9, 1), b: hexg.NewHex(9, 0), want: +1},
		{id: 8, a: hexg.NewHex(0, -1), b: hexg.NewHex(-100, 0), want: -1},
		// negative coordinates
		{id: 9, a: hexg.NewHex(-1, -2), b: hexg.NewHex(-1, -1), want: -1},
		{id: 10, a: hexg.NewHex(3, -4), b: hexg.NewHex(2, -4), want: +1},
	} {
		if want, got := tc.want, tc.a.Compare(tc.b); want != got {
			t.Errorf("%d: compare: want %d, got %d\n", tc.id, want, got)
		}
		// Compare must be antisymmetric.
		if want, got := -tc.want, tc.b.Compare(tc.a); want != got {
			t.Errorf("%d: reversed: want %d, got %d\n", tc.id, want, got)
		}
		// Compare and Equals must agree.
		if want, got := tc.a.Equals(tc.b), tc.a.Compare(tc.b) == 0; want != got {
			t.Errorf("%d: equals: want %v, got %v\n", tc.id, want, got)
		}
	}
}

func TestHexSet_Sorted(t *testing.T) {
	layout := hexg.NewLayout(hexg.OddR, hexg.Point{X: 10, Y: 10}, hexg.Point{})

	for _, tc := range []struct {
		id   int
		gs   hexg.HexSet
		want []hexg.Hex
	}{
		{id: 1, gs: hexg.HexSet{}, want: []hexg.Hex{}},
		{id: 2, gs: nil, want: []hexg.Hex{}},
		{
			id: 3,
			gs: hexg.HexSet{
				hexg.NewHex(0, 0):  struct{}{},
				hexg.NewHex(1, -1): struct{}{},
				hexg.NewHex(-1, 0): struct{}{},
				hexg.NewHex(-1, 1): struct{}{},
				hexg.NewHex(1, 0):  struct{}{},
				hexg.NewHex(0, -1): struct{}{},
				hexg.NewHex(0, 1):  struct{}{},
			},
			// row-major: r ascending, q ascending within a row
			want: []hexg.Hex{
				hexg.NewHex(0, -1), hexg.NewHex(1, -1),
				hexg.NewHex(-1, 0), hexg.NewHex(0, 0), hexg.NewHex(1, 0),
				hexg.NewHex(-1, 1), hexg.NewHex(0, 1),
			},
		},
		{
			id: 4,
			gs: layout.Rectangle(0, 1, 0, 1),
			want: []hexg.Hex{
				hexg.NewHex(0, 0), hexg.NewHex(1, 0),
				hexg.NewHex(0, 1), hexg.NewHex(1, 1),
			},
		},
	} {
		got := tc.gs.Sorted()
		if got == nil {
			t.Errorf("%d: sorted: want non-nil slice, got nil\n", tc.id)
			continue
		}
		if want, got := len(tc.want), len(got); want != got {
			t.Errorf("%d: length: want %d, got %d\n", tc.id, want, got)
			continue
		}
		for i := range tc.want {
			if want, got := tc.want[i], got[i]; want.NotEquals(got) {
				t.Errorf("%d: index %d: want %s, got %s\n", tc.id, i, want.ConciseString(), got.ConciseString())
			}
		}
	}
}

// TestHexSet_Sorted_Deterministic pins the property the helper exists for:
// ranging a HexSet is randomized, Sorted is not. The generators are the sets
// callers actually hold, so exercise those rather than a literal.
func TestHexSet_Sorted_Deterministic(t *testing.T) {
	layout := hexg.NewLayout(hexg.OddR, hexg.Point{X: 10, Y: 10}, hexg.Point{})

	for _, tc := range []struct {
		id int
		gs hexg.HexSet
	}{
		{id: 1, gs: layout.Hexagon(3)},
		{id: 2, gs: layout.Rectangle(-2, 3, -2, 3)},
		{id: 3, gs: layout.ParallelogramQR(-2, -2, 2, 2)},
		{id: 4, gs: layout.TriangleUpDown(4)},
	} {
		first := tc.gs.Sorted()
		if want, got := len(tc.gs), len(first); want != got {
			t.Errorf("%d: length: want %d, got %d\n", tc.id, want, got)
			continue
		}
		// every element of the set appears exactly once, in Compare order
		for i := 1; i < len(first); i++ {
			if want, got := -1, first[i-1].Compare(first[i]); want != got {
				t.Errorf("%d: index %d: want %d, got %d\n", tc.id, i, want, got)
			}
		}
		for h := range tc.gs {
			if _, found := slices.BinarySearchFunc(first, h, hexg.Hex.Compare); !found {
				t.Errorf("%d: missing %s\n", tc.id, h.ConciseString())
			}
		}
		// repeated calls agree, where ranging the map would not
		for pass := 0; pass < 20; pass++ {
			if want, got := first, tc.gs.Sorted(); !slices.Equal(want, got) {
				t.Errorf("%d: pass %d: order changed\n", tc.id, pass)
			}
		}
	}
}
