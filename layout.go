// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

// Layout_i defines the interface for layouts.
//
// Orientation is important for offset coordinates and every layout
// that implements this interface is expected to implement that per
// the Red Blob Games guide.
type Layout_i interface {
	// IsHorizontal returns true if the layout has horizontal rows.
	// Horizontal layouts have pointy-top hexes, staggered columns, and horizontal rows.
	IsHorizontal() bool

	// IsVertical returns true if the layout has vertical columns.
	// Vertical layouts have flat-top hexes, vertical columns, and staggered rows.
	IsVertical() bool

	// OffsetType returns the type of offset used for columns and rows.
	OffsetType() LayoutOffset_e

	// DirectionToBearing returns the bearing of a direction in the layout
	DirectionToBearing(direction int) string

	// HexagonalGrid returns a grid centered about a hex.
	HexagonalGrid(center Hex, radius int) GridStore

	// HexCorner returns the screen coordinates of the hex corner.
	// We should define what "corner" means in this context.
	HexCorner(h Hex, corner int) Point

	// HexCorners returns the screen coordinates for every corner of the hex.
	HexCorners(h Hex) [6]Point

	// HexToOffsetCoord returns the offset coordinates of the hex.
	// Uses the offset from the layout to shift rows and columns correctly.
	HexToOffsetCoord(h Hex) OffsetCoord

	// HexToPixel returns the origin of the hex on the screen as a pixel.
	HexToPixel(h Hex) Point

	// OffsetColRowToHex returns a new Hex using offset column and row coordinates.
	OffsetColRowToHex(col, row int) Hex

	// OffsetCoordToHex returns a new Hex from the OffsetCoord.
	OffsetCoordToHex(oc OffsetCoord) Hex

	// ParallelogramGrid returns a grid originating at (0,0,0).
	// I don't understand the comment in the source about there
	// being three coordinates and the caller has to choose two.
	// does that mean the grid has three orientations?
	ParallelogramGrid(q1, r1, q2, r2 int) GridStore

	// PixelToHexRounded turns a fractional hex into a regular hex coordinate:
	PixelToHexRounded(p Point) Hex

	// PixelToFractionalHex returns the fractional hex that encloses the pixel.
	// In theory, the origin of that fractional hex will be the pixel.
	PixelToFractionalHex(p Point) FractionalHex

	// PolygonCornerOffset returns the offset from the center of a hex to a corner.
	// We should define what the parameter "corner" means. Which corner?
	PolygonCornerOffset(corner int) Point

	// PolygonCornerOffsets returns the offset for every corner of a hex.
	PolygonCornerOffsets() [6]Point

	// RectangularGrid returns a grid centered about a hex.
	RectangularGrid(center Hex, left, right, top, bottom int) GridStore

	// TriagonalGrid returns a grid originating at (0,0,0).
	// there's a comment in the source about flipping the y-axis to
	// change the direction of the triangle, but I don't understand
	// how to implement that.
	TriagonalGrid(side_length int) GridStore
}
