// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cube

import (
	"testing"
)

func TestCube_Neighbor(t *testing.T) {
	from := Cube{q: 0, r: 0, s: 0}
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
		{id: 37, direction: 3, expect: "-1+3-2"}, {id: 37, direction: 3, expect: "-2+3-1"}, {id: 38, direction: 3, expect: "-3+3+0"},
	} {
		to := from.Neighbor(move.direction)
		if to.ConciseString() != move.expect {
			t.Fatalf("%d: from %q: to %d: got %q, want %q\n", move.id, from.ConciseString(), move.direction, to.ConciseString(), move.expect)
		}
		from = to
	}
}
