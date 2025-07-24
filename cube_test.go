// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"testing"

	"github.com/maloquacious/hexg"
)

func TestCubeDistance(t *testing.T) {
	for _, tc := range []struct {
		id       int
		from, to hexg.Hex
		expect   int
	}{
		{1, hexg.NewHexFromAxialCoords(0, 0), hexg.NewHexFromAxialCoords(0, 0), 0},
		{2, hexg.NewHexFromAxialCoords(0, 0), hexg.NewHexFromAxialCoords(0, 1), 1},
		{3, hexg.NewHexFromAxialCoords(0, 0), hexg.NewHexFromAxialCoords(0, 2), 2},
		{4, hexg.NewHexFromAxialCoords(0, 0), hexg.NewHexFromAxialCoords(0, 3), 3},
		{5, hexg.NewHexFromAxialCoords(0, 0), hexg.NewHexFromAxialCoords(0, 4), 4},
		{6, hexg.NewHexFromAxialCoords(0, 0), hexg.NewHexFromAxialCoords(1, 3), 4},
		{7, hexg.NewHexFromAxialCoords(0, 0), hexg.NewHexFromAxialCoords(2, 2), 4},
		{8, hexg.NewHexFromAxialCoords(0, 0), hexg.NewHexFromAxialCoords(3, 1), 4},
		{9, hexg.NewHexFromAxialCoords(-3, -1), hexg.NewHexFromAxialCoords(4, -1), 7},
		{10, hexg.NewHexFromAxialCoords(-1, -3), hexg.NewHexFromAxialCoords(1, 3), 8},
	} {
		if distance := tc.from.Distance(tc.to); distance != tc.expect {
			t.Errorf("distance: from %s: to %s: got %d, want %d\n", tc.from.ConciseString(), tc.to.ConciseString(), distance, tc.expect)
		}
		if distance := tc.to.Distance(tc.from); distance != tc.expect {
			t.Errorf("distance: from %s: to %s: got %d, want %d\n", tc.to.ConciseString(), tc.from.ConciseString(), distance, tc.expect)
		}
	}
}

func TestCubeMovement(t *testing.T) {
	from := hexg.NewHex(0, 0, 0)
	if from.ConciseString() != "+0+0+0" {
		t.Fatalf("move: %3d: got %q, want %q\n", 0, from.ConciseString(), "+0+0+0")
	}
	for _, move := range []struct {
		id        int
		direction int
		expect    string
	}{
		// move one hex and then back
		{1, 0, "+1+0-1"}, {2, 3, "+0+0+0"},
		{3, 1, "+1-1+0"}, {4, 4, "+0+0+0"},
		{5, 2, "+0-1+1"}, {6, 5, "+0+0+0"},
		{7, 3, "-1+0+1"}, {8, 0, "+0+0+0"},
		{9, 4, "-1+1+0"}, {10, 1, "+0+0+0"},
		{11, 5, "+0+1-1"}, {12, 2, "+0+0+0"},
		// circle around
		{13, 0, "+1+0-1"},
		{14, 1, "+2-1-1"},
		{15, 2, "+2-2+0"},
		{16, 3, "+1-2+1"},
		{17, 4, "+0-1+1"},
		{18, 5, "+0+0+0"},
		// move three hexes in each direction
		{19, 4, "-1+1+0"}, {20, 4, "-2+2+0"}, {21, 4, "-3+3+0"},
		{22, 2, "-3+2+1"}, {23, 2, "-3+1+2"}, {24, 2, "-3+0+3"},
		{25, 1, "-2-1+3"}, {26, 1, "-1-2+3"}, {27, 1, "+0-3+3"},
		{28, 0, "+1-3+2"}, {29, 0, "+2-3+1"}, {30, 0, "+3-3+0"},
		{31, 5, "+3-2-1"}, {32, 5, "+3-1-2"}, {33, 5, "+3+0-3"},
		{34, 4, "+2+1-3"}, {35, 4, "+1+2-3"}, {36, 4, "+0+3-3"},
		{37, 3, "-1+3-2"}, {37, 3, "-2+3-1"}, {38, 3, "-3+3+0"},
	} {
		to := from.Neighbor(move.direction)
		if to.ConciseString() != move.expect {
			t.Fatalf("move: %3d: from %s: to %d: got %q, want %q\n", move.id, from.ConciseString(), move.direction, to.ConciseString(), move.expect)
		}
		from = to
	}
}

func TestCubeTribeNetMovement(t *testing.T) {
	from := hexg.NewHex(0, 0, 0)
	if from.ConciseString() != "+0+0+0" {
		t.Fatalf("move: %3d: got %q, want %q\n", 0, from.ConciseString(), "+0+0+0")
	}
	for _, move := range []struct {
		id        int
		direction int
		expect    string
	}{
		// move one hex and then back
		{1, hexg.TNSouthEast, "+1+0-1"}, {2, hexg.TNNorthWest, "+0+0+0"},
		{3, hexg.TNNorthEast, "+1-1+0"}, {4, hexg.TNSouthWest, "+0+0+0"},
		{5, hexg.TNNorth, "+0-1+1"}, {6, hexg.TNSouth, "+0+0+0"},
		{7, hexg.TNNorthWest, "-1+0+1"}, {8, hexg.TNSouthEast, "+0+0+0"},
		{9, hexg.TNSouthWest, "-1+1+0"}, {10, hexg.TNNorthEast, "+0+0+0"},
		{11, hexg.TNSouth, "+0+1-1"}, {12, hexg.TNNorth, "+0+0+0"},
	} {
		to := from.Neighbor(move.direction)
		if to.ConciseString() != move.expect {
			t.Fatalf("move: %3d: from %s: move %q: got %q, want %q\n", move.id, from.ConciseString(), hexg.TribeNetDirectionString(move.direction), to.ConciseString(), move.expect)
		}
		from = to
	}
}
