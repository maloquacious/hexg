// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cube

import (
	"testing"
)

func TestCube_Distance(t *testing.T) {
	for _, tc := range []struct {
		id       int
		a, b     Cube
		distance int
	}{
		{id: 1, a: Cube{q: 0, r: 0, s: 0}, b: Cube{q: 0, r: 0, s: 0}, distance: 0},
		{id: 2, a: Cube{q: 0, r: 0, s: 0}, b: Cube{q: 0, r: 1, s: -1}, distance: 1},
		{id: 3, a: Cube{q: 0, r: 0, s: 0}, b: Cube{q: 0, r: 2, s: -2}, distance: 2},
		{id: 4, a: Cube{q: 0, r: 0, s: 0}, b: Cube{q: 0, r: 3, s: -3}, distance: 3},
		{id: 5, a: Cube{q: 0, r: 0, s: 0}, b: Cube{q: 0, r: 4, s: -4}, distance: 4},
		{id: 6, a: Cube{q: 0, r: 0, s: 0}, b: Cube{q: 1, r: 3, s: -4}, distance: 4},
		{id: 7, a: Cube{q: 0, r: 0, s: 0}, b: Cube{q: 2, r: 2, s: -4}, distance: 4},
		{id: 8, a: Cube{q: 0, r: 0, s: 0}, b: Cube{q: 3, r: 1, s: -4}, distance: 4},
		{id: 9, a: Cube{q: -3, r: -1, s: 4}, b: Cube{q: 4, r: -1, s: -3}, distance: 7},
		{id: 10, a: Cube{q: -1, r: -3, s: 4}, b: Cube{q: 1, r: 3, s: -4}, distance: 8},
	} {
		if distance := tc.a.Distance(tc.b); distance != tc.distance {
			t.Errorf("distance: from %q: to %q: got %d, want %d\n", tc.a.ConciseString(), tc.b.ConciseString(), distance, tc.distance)
		}
		if distance := tc.b.Distance(tc.a); distance != tc.distance {
			t.Errorf("distance: from %q: to %q: got %d, want %d\n", tc.b.ConciseString(), tc.a.ConciseString(), distance, tc.distance)
		}
	}
}
