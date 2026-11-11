// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import (
	"fmt"
	"slices"
	"strings"
)

// Bearing names a compass direction independently of hex orientation.
// Its values are not the direction indices accepted by Hex.Neighbor;
// use Layout.BearingToDirection to convert them.
// North is toward negative screen Y and east is toward positive screen X
// for layouts with positive pixel sizes. Pixel size and origin do not affect
// bearing conversions.
type Bearing int

const (
	North Bearing = iota
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
)

var bearingNames = [...]string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}

// String returns the compass abbreviation, or "Bearing(n)" for an invalid value.
func (b Bearing) String() string {
	if b < North || b > NorthWest {
		return fmt.Sprintf("Bearing(%d)", int(b))
	}
	return bearingNames[b]
}

// ParseBearing parses N, NE, E, SE, S, SW, W, or NW, ignoring case and surrounding
// whitespace. It does not check whether a layout supports the bearing.
func ParseBearing(s string) (Bearing, error) {
	name := strings.ToUpper(strings.TrimSpace(s))
	if b := slices.Index(bearingNames[:], name); b >= 0 {
		return Bearing(b), nil
	}
	return 0, fmt.Errorf("invalid bearing %q", s)
}

// Opposite returns the bearing 180 degrees from b. It panics if b is invalid.
func (b Bearing) Opposite() Bearing {
	if b < North || b > NorthWest {
		panic("invalid bearing")
	}
	return (b + 4) % 8
}

// DirectionToBearing converts an existing direction index in 0..5 to its
// compass bearing for l. It returns an error for an out-of-range index rather
// than coercing it. Odd/even offset parity does not change the mapping.
func (l Layout) DirectionToBearing(direction int) (Bearing, error) {
	if direction < 0 || direction >= 6 {
		return 0, fmt.Errorf("invalid direction %d: must be in 0..5", direction)
	}
	if l.IsFlat() {
		return [...]Bearing{SouthEast, NorthEast, North, NorthWest, SouthWest, South}[direction], nil
	}
	return [...]Bearing{East, NorthEast, NorthWest, West, SouthWest, SouthEast}[direction], nil
}

// BearingToDirection converts b to an existing direction index in 0..5 for l.
// Flat-top layouts reject E and W; pointy-top layouts reject N and S.
// Invalid Bearing values also return an error.
func (l Layout) BearingToDirection(b Bearing) (int, error) {
	for direction := range 6 {
		candidate, _ := l.DirectionToBearing(direction)
		if candidate == b {
			return direction, nil
		}
	}
	return 0, fmt.Errorf("bearing %s is not supported by layout offset %d", b, l.offset)
}

// CompassVectors returns neighbor vectors in clockwise compass order:
// N, NE, SE, S, SW, NW for flat-top layouts, and
// NE, E, SE, SW, W, NW for pointy-top layouts.
// These array positions are not the indices accepted by Hex.Neighbor.
func (l Layout) CompassVectors() [6]Hex {
	start := 1 // NE for pointy-top layouts
	if l.IsFlat() {
		start = 2 // N for flat-top layouts
	}
	var vectors [6]Hex
	for i := range vectors {
		vectors[i] = DirectionVector((start - i + 6) % 6)
	}
	return vectors
}

// BearingBetween returns the bearing from a to an adjacent hex b.
// It returns an error for identical or nonadjacent hexes. Adjacency is planar:
// this method does not normalize coordinates or recognize cylinder seams.
func (l Layout) BearingBetween(a, b Hex) (Bearing, error) {
	delta := b.Subtract(a)
	if direction := slices.Index(directionVectors[:], delta); direction >= 0 {
		return l.DirectionToBearing(direction)
	}
	return 0, fmt.Errorf("hexes %s and %s are not adjacent", a, b)
}
