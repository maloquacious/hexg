// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/maloquacious/hexg"
)

func TestBearing(t *testing.T) {
	for _, tc := range []struct {
		id       int
		bearing  hexg.Bearing
		name     string
		opposite hexg.Bearing
	}{
		{id: 1, bearing: hexg.North, name: "N", opposite: hexg.South},
		{id: 2, bearing: hexg.NorthEast, name: "NE", opposite: hexg.SouthWest},
		{id: 3, bearing: hexg.East, name: "E", opposite: hexg.West},
		{id: 4, bearing: hexg.SouthEast, name: "SE", opposite: hexg.NorthWest},
		{id: 5, bearing: hexg.South, name: "S", opposite: hexg.North},
		{id: 6, bearing: hexg.SouthWest, name: "SW", opposite: hexg.NorthEast},
		{id: 7, bearing: hexg.West, name: "W", opposite: hexg.East},
		{id: 8, bearing: hexg.NorthWest, name: "NW", opposite: hexg.SouthEast},
	} {
		if want, got := tc.name, tc.bearing.String(); want != got {
			t.Errorf("%d: string: want %q, got %q", tc.id, want, got)
		}
		if want, got := tc.opposite, tc.bearing.Opposite(); want != got {
			t.Errorf("%d: opposite: want %s, got %s", tc.id, want, got)
		}
		for _, input := range []string{tc.name, strings.ToLower(tc.name), " \t" + tc.name + "\n"} {
			got, err := hexg.ParseBearing(input)
			if got != tc.bearing || err != nil {
				t.Errorf("%d: parse %q: want %s, got %s, %v", tc.id, input, tc.bearing, got, err)
			}
		}
	}
	for _, input := range []string{"", " ", "north", "N E", "NN", "NE SW", "1"} {
		if _, err := hexg.ParseBearing(input); err == nil {
			t.Errorf("parse %q: want an error", input)
		}
	}
}

func TestBearing_InvalidValue(t *testing.T) {
	for _, tc := range []struct {
		id      int
		bearing hexg.Bearing
		name    string
	}{
		{id: 1, bearing: -1, name: "Bearing(-1)"},
		{id: 2, bearing: 8, name: "Bearing(8)"},
	} {
		if want, got := tc.name, tc.bearing.String(); want != got {
			t.Errorf("%d: string: want %q, got %q", tc.id, want, got)
		}
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%d: opposite: want a panic", tc.id)
				}
			}()
			tc.bearing.Opposite()
		}()
	}
}

func TestLayout_Bearings(t *testing.T) {
	for _, tc := range []struct {
		id     int
		offset hexg.LayoutOffset
		flat   bool
	}{
		{id: 1, offset: hexg.OddQ, flat: true},
		{id: 2, offset: hexg.EvenQ, flat: true},
		{id: 3, offset: hexg.OddR},
		{id: 4, offset: hexg.EvenR},
	} {
		// Nonuniform size and a nonzero origin must not change compass names.
		layout := hexg.NewLayout(tc.offset, hexg.Point{X: 17, Y: 31}, hexg.Point{X: -50, Y: 100})
		bearings := [6]hexg.Bearing{hexg.East, hexg.NorthEast, hexg.NorthWest, hexg.West, hexg.SouthWest, hexg.SouthEast}
		clockwise := [6]hexg.Hex{
			hexg.NewHex(1, -1), hexg.NewHex(1, 0), hexg.NewHex(0, 1),
			hexg.NewHex(-1, 1), hexg.NewHex(-1, 0), hexg.NewHex(0, -1),
		}
		if tc.flat {
			bearings = [6]hexg.Bearing{hexg.SouthEast, hexg.NorthEast, hexg.North, hexg.NorthWest, hexg.SouthWest, hexg.South}
			clockwise = [6]hexg.Hex{
				hexg.NewHex(0, -1), hexg.NewHex(1, -1), hexg.NewHex(1, 0),
				hexg.NewHex(0, 1), hexg.NewHex(-1, 1), hexg.NewHex(-1, 0),
			}
		}
		if got := layout.CompassVectors(); got != clockwise {
			t.Errorf("%d: clockwise vectors: want %v, got %v", tc.id, clockwise, got)
		}
		for direction, want := range bearings {
			got, err := layout.DirectionToBearing(direction)
			if got != want || err != nil {
				t.Errorf("%d: direction %d: want %s, got %s, %v", tc.id, direction, want, got, err)
			}
			// Translate away from the origin to catch treating b as a vector
			// instead of computing the displacement from a.
			a := hexg.NewHex(-4, 7)
			b := a.Neighbor(direction)
			if got, err := layout.BearingBetween(a, b); got != want || err != nil {
				t.Errorf("%d: bearing between %v and %v: want %s, got %s, %v", tc.id, a, b, want, got, err)
			}
			if got, err := layout.BearingBetween(b, a); got != want.Opposite() || err != nil {
				t.Errorf("%d: reverse bearing: want %s, got %s, %v", tc.id, want.Opposite(), got, err)
			}
		}
		for bearing := hexg.North; bearing <= hexg.NorthWest; bearing++ {
			want := slices.Index(bearings[:], bearing)
			got, err := layout.BearingToDirection(bearing)
			if want < 0 {
				if err == nil {
					t.Errorf("%d: unsupported %s: want an error", tc.id, bearing)
				}
			} else if got != want || err != nil {
				t.Errorf("%d: bearing %s: want %d, got %d, %v", tc.id, bearing, want, got, err)
			}
		}
		for _, direction := range []int{-1, 6, 12} {
			if _, err := layout.DirectionToBearing(direction); err == nil {
				t.Errorf("%d: direction %d: want an error", tc.id, direction)
			}
		}
		for _, bearing := range []hexg.Bearing{-1, 8, 16} {
			if _, err := layout.BearingToDirection(bearing); err == nil {
				t.Errorf("%d: bearing %s: want an error", tc.id, bearing)
			}
		}
		for _, b := range []hexg.Hex{hexg.NewHex(-4, 7), hexg.NewHex(-2, 7), hexg.NewHex(0, -3)} {
			if _, err := layout.BearingBetween(hexg.NewHex(-4, 7), b); err == nil {
				t.Errorf("%d: nonadjacent endpoint %v: want an error", tc.id, b)
			}
		}
	}
}

func ExampleBearing() {
	layout := hexg.NewLayout(hexg.OddQ, hexg.Point{}, hexg.Point{})
	bearing, err := hexg.ParseBearing("ne")
	if err != nil {
		panic(err)
	}
	direction, err := layout.BearingToDirection(bearing)
	if err != nil {
		panic(err)
	}
	from := hexg.NewHex(2, 3)
	to := from.Neighbor(direction)
	fmt.Println(bearing, direction, to)
	fmt.Println(bearing.Opposite())
	reported, err := layout.BearingBetween(from, to)
	if err != nil {
		panic(err)
	}
	fmt.Println(reported)
	_, err = layout.BearingToDirection(hexg.East)
	fmt.Println(err != nil) // flat-top hexes have no due-east neighbor
	// Output:
	// NE 1 (3, 2, -5)
	// SW
	// NE
	// true
}
