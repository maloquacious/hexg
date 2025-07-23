// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import "fmt"

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
