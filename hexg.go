// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import (
	"fmt"
	"strconv"
	"unicode"
)

// this file implements the public methods

// 1.0 Hex coordinates

func NewHex(q, r, s int) Hex {
	return new_hex(q, r, s)
}

func (h Hex) String() string {
	//return fmt.Sprintf("%+d%+d%+d", h.q, h.r, h.s)
	return fmt.Sprintf("%d,%d,%d", h.q, h.r, h.s)
}

// 1.1 Equality

// 1.2 Coordinate arithmetic

// 1.3 Distance

// 1.3.1 Neighbors

func (h Hex) Neighbor(direction int) Hex {
	return hex_neighbor(h, direction)
}

// 2.0 Layout

func NewLayoutPointy(size, origin Point, shoveOddRowsRight bool) Layout {
	if shoveOddRowsRight {
		return new_layout(layout_pointy, size, origin, odd_r)
	}
	return new_layout(layout_pointy, size, origin, even_r)
}

func NewLayoutFlat(size, origin Point, shoveOddColumnsDown bool) Layout {
	if shoveOddColumnsDown {
		return new_layout(layout_flat, size, origin, odd_q)
	}
	return new_layout(layout_flat, size, origin, even_q)
}

func (layout Layout) IsFlatTop() bool {
	return layout.offset == odd_q || layout.offset == even_q
}

func (layout Layout) IsPointyTop() bool {
	return layout.offset == odd_r || layout.offset == even_r
}

func NewPoint(x, y float64) Point {
	return new_point(x, y)
}

func (p Point) String() string {
	return fmt.Sprintf("%g,%g", p.x, p.y)
}

// 2.1 Hex to screen

func (h Hex) ToPixel(layout Layout) Point {
	return hex_to_pixel(layout, h)
}

// 2.2 Screen to hex

func (p Point) ToFractionalHex(layout Layout) FractionalHex {
	return pixel_to_hex_fractional(layout, p)
}

// 2.3 Drawing a hex

func (h Hex) CornerOffset(layout Layout, corner int) Point {
	return hex_corner_offset(layout, corner)
}

func (layout Layout) HexCornerOffset(corner int) Point {
	return hex_corner_offset(layout, corner)
}

func (layout Layout) PolygonCorners() [6]Point {
	return polygon_corners(layout, Hex{})
}

func (h Hex) PolygonCorners(layout Layout) [6]Point {
	return polygon_corners(layout, h)
}

// 2.4 Layout examples

// 3.0 Fractional Hex

// 9.0 Other Systems

// 9.1 TribeNet Coordinates

// TribeNet coordinates are in the form "AB 0102":
// - "A" (grid row) and "B" (grid column) identify a sub-map.
// - "0102" is the in-submap position: column 01 (1-based) and row 02 (1-based).
// Each sub-map is 30 columns wide and 21 rows tall,
// with "0101" as the upper-left and "3021" as the lower-right corner.
//
// The global map origin is (1,1) at the upper-left.
// Even-numbered columns are vertically offset (odd-q layout), so (2,1) is southeast of (1,1).
//
// TribeNet coordinates are converted to OffsetCoord using "odd-q" layout,
// with the origin translated by (-1, -1) so "AA 0101" becomes (0,0).

const (
	tnRowsPerGrid  = 21
	tnColsPerGrid  = 30
	tnMaxGridIndex = 26 // A ... Z -> 1 ... 26
)

// NewTribeNetCoord parses a Tribenet coordinate (eg, "AB 0102") and returns
// an OffsetCoord. All TribeNet maps are "odd-q,", origin is (1,1) and will
// be translated to (0,0), sub-maps are 21 rows x 30 columns.
//
// Invalid inputs will return an error.
//
// TribeNet coordinate "AA 0101" corresponds to OffsetCoord (  0,  0).
// TribeNet coordinate "BC 0824" corresponds to OffsetCoord ( 61, 43).
// TribeNet coordinate "JK 0609" corresponds to OffsetCoord (163,104).
// TribeNet coordinate "ZZ 3021" corresponds to OffsetCoord (779,545).

func NewTribeNetCoord(input string) (OffsetCoord, error) {
	if len(input) != 7 || input[2] != ' ' {
		return OffsetCoord{}, fmt.Errorf("invalid format: expected 'AB 0102'")
	}
	gridRowRune, gridColRune := rune(input[0]), rune(input[1])
	subColStr, subRowStr := input[3:5], input[5:7]

	// convert letters to 0-based grid index (A=0, B=1, ...)
	if !(unicode.IsUpper(gridRowRune) && unicode.IsLetter(gridRowRune)) {
		return OffsetCoord{}, fmt.Errorf("invalid grid row: must be uppercase A-Z")
	} else if !(unicode.IsUpper(gridColRune) && unicode.IsLetter(gridColRune)) {
		return OffsetCoord{}, fmt.Errorf("invalid grid column: must be uppercase A-Z")
	}
	gridRowOffset := int(gridRowRune-'A') * tnRowsPerGrid
	gridColOffset := int(gridColRune-'A') * tnColsPerGrid

	// subCol is 1-based in the input; converted to 0-based for OffsetCoord
	subCol, err := strconv.Atoi(subColStr)
	if !(err == nil && 1 <= subCol && subCol <= tnColsPerGrid) {
		return OffsetCoord{}, fmt.Errorf("invalid sub-map column: %s", subColStr)
	}
	subCol -= 1

	// subRow is 1-based in the input; converted to 0-based for OffsetCoord
	subRow, err := strconv.Atoi(subRowStr)
	if !(err == nil && 1 <= subRow && subRow <= tnRowsPerGrid) {
		return OffsetCoord{}, fmt.Errorf("invalid sub-map row: %s", subRowStr)
	}
	subRow -= 1

	// convert global row and column to an offset coordinate
	return OffsetCoord{
		col: gridColOffset + subCol,
		row: gridRowOffset + subRow,
	}, nil
}

// ToTribeNetCoord converts an OffsetCoord to a TribeNet coordinate string ("AB 0102").
//
// The conversion is based on:
// - A 0-based global offset grid with "odd-q" layout and origin (0,0) = "AA 0101".
// - Each 30×21 cell in the global grid is a sub-map labeled by row and column letters (A–Z).
// - The function converts the offset col/row to sub-map ID and in-map sub-coordinates.
//
// Returns an error if the coordinate falls outside the supported 26×26 letter grid.
func (oc OffsetCoord) ToTribeNetCoord() (string, error) {
	if oc.col < 0 || oc.row < 0 {
		return "", fmt.Errorf("invalid offset coordinates: %v", oc)
	}

	gridRow, gridCol := oc.row/tnRowsPerGrid, oc.col/tnColsPerGrid
	if gridRow >= tnMaxGridIndex || gridCol >= tnMaxGridIndex {
		return "", fmt.Errorf("coordinates out of range for A-Z grid system")
	}
	gridRowChar, gridColChar := 'A'+rune(gridRow), 'A'+rune(gridCol)
	subCol, subRow := (oc.col%tnColsPerGrid)+1, (oc.row%tnRowsPerGrid)+1

	return fmt.Sprintf("%c%c %02d%02d", gridRowChar, gridColChar, subCol, subRow), nil
}
