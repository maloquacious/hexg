// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import (
	"fmt"
	"strconv"
	"unicode"
)

// TribeNet coordinates are in the form "AB 0102":
// - "A" (grid row) and "B" (grid column) identify a sub-map.
// - "0102" is the in-submap position: column 01 (1-based) and row 02 (1-based).
// Each sub-map is 30 columns wide and 21 rows tall,
// with "0101" as the upper-left and "3021" as the lower-right corner.
//
// The global map origin is 0101 and is at the upper-left. On the TribeNet map,
// even-numbered columns are pushed down, so 0201 is southeast of 0101.
//
//
// TribeNet coordinates are converted to OffsetCoord using "odd-q" layout,
// with the origin translated by (-1, -1) so "AA 0101" becomes (0,0).

// NewTribeNetLayout returns an initialized layout for TribeNet maps.
// It uses the VerticalOddQLayout.
func NewTribeNetLayout() TribeNetLayout {
	return TribeNetLayout{
		VerticalOddQLayout{
			size:   Point{1, 1},
			origin: Point{0, 0},
		},
	}
}

const (
	tnRowsPerGrid  = 21
	tnColsPerGrid  = 30
	tnMaxGridIndex = 26 // A ... Z -> 1 ... 26
)

// TribeNetLayout returns an odd-q layout with some adaptations for TribeNet maps.
type TribeNetLayout struct {
	VerticalOddQLayout
}

// ColRowToTribeNetCoord converts a col, row value to a TribeNet coordinate string ("AB 0102").
//
// The conversion is based on:
// - A 0-based global offset grid with "odd-q" layout and origin (0,0) = "AA 0101".
// - Each 30 wide by 21 high cell in the global grid is a sub-map labeled by row and column letters (A–Z).
// - The function converts the offset col/row to sub-map ID and in-map sub-coordinates.
//
// Returns an error if the coordinate falls outside the supported 26×26 letter grid.
func (l TribeNetLayout) ColRowToTribeNetCoord(col, row int) (string, error) {
	if col < 0 || row < 0 {
		return "", fmt.Errorf("invalid col, row: %d, %d", col, row)
	}

	gridRow, gridCol := row/tnRowsPerGrid, col/tnColsPerGrid
	if gridRow >= tnMaxGridIndex || gridCol >= tnMaxGridIndex {
		return "", fmt.Errorf("coordinates out of range for A-Z grid system")
	}
	gridRowChar, gridColChar := 'A'+rune(gridRow), 'A'+rune(gridCol)
	subCol, subRow := (col % tnColsPerGrid), (row % tnRowsPerGrid)

	// be sure to translate the coordinates by (+1,+1) to shift the origin back
	return fmt.Sprintf("%c%c %02d%02d", gridRowChar, gridColChar, subCol+1, subRow+1), nil
}

// DirectionToBearing overrides the default to use the TribeNet bearings.
func (l TribeNetLayout) DirectionToBearing(direction int) string {
	return tnBearingNames[(6+(direction%6))%6]
}

func (l TribeNetLayout) HexToTribeNetCoord(h Hex) (string, error) {
	oc := l.HexToOffsetCoord(h)
	return l.ColRowToTribeNetCoord(oc.Col, oc.Row)
}

// TribeNetCoordToColRow parses a Tribenet coordinate (eg, "AB 0102") and returns
// the column and row. The result is translated to use the correct origin for the map.
//
// Invalid inputs will return an error.
//
// TribeNet coordinate "AA 0101" corresponds to (  0,  0).
// TribeNet coordinate "BC 0824" corresponds to ( 61, 43).
// TribeNet coordinate "JK 0609" corresponds to (163,104).
// TribeNet coordinate "ZZ 3021" corresponds to (779,545).
func (l TribeNetLayout) TribeNetCoordToColRow(input string) (col, row int, err error) {
	if len(input) != 7 || input[2] != ' ' {
		return 0, 0, fmt.Errorf("invalid format: expected 'AB 0102'")
	}
	gridRowRune, gridColRune := rune(input[0]), rune(input[1])
	subColStr, subRowStr := input[3:5], input[5:7]

	// convert letters to 0-based grid index (A=0, B=1, ...)
	if !(unicode.IsUpper(gridRowRune) && unicode.IsLetter(gridRowRune)) {
		return 0, 0, fmt.Errorf("invalid grid row: must be uppercase A-Z")
	} else if !(unicode.IsUpper(gridColRune) && unicode.IsLetter(gridColRune)) {
		return 0, 0, fmt.Errorf("invalid grid column: must be uppercase A-Z")
	}
	gridRowOffset := int(gridRowRune-'A') * tnRowsPerGrid
	gridColOffset := int(gridColRune-'A') * tnColsPerGrid

	subCol, err := strconv.Atoi(subColStr)
	if !(err == nil && 1 <= subCol && subCol <= tnColsPerGrid) {
		return 0, 0, fmt.Errorf("invalid sub-map column: %s", subColStr)
	}

	subRow, err := strconv.Atoi(subRowStr)
	if !(err == nil && 1 <= subRow && subRow <= tnRowsPerGrid) {
		return 0, 0, fmt.Errorf("invalid sub-map row: %s", subRowStr)
	}

	// convert the coordinates to the grid and then translate by (-1,-1) to shift the origin
	return gridColOffset + subCol - 1, gridRowOffset + subRow - 1, nil
}

// TribeNetCoordToHex converts the TribeNet coordinates to Hex.
//
// Invalid inputs will return an error.
func (l TribeNetLayout) TribeNetCoordToHex(input string) (Hex, error) {
	col, row, err := l.TribeNetCoordToColRow(input)
	if err != nil {
		return Hex{}, err
	}
	return l.OffsetColRowToHex(col, row), nil
}

// define convenient names for directions on a TribeNet grid.
const (
	TNNorth         = N
	TNNorthEast     = ENE
	TNSouthEast int = ESE
	TNSouth         = S
	TNSouthWest     = WSW
	TNNorthWest     = WNW
)

var tnBearingNames = []string{"SE", "NE", "N", "NW", "SW", "S"}
