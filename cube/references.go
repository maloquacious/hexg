// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cube

// This code is based on the Reb Blob Games guide to hex grids:
// * https://www.redblobgames.com/grids/hexagons/
// * https://www.redblobgames.com/grids/hexagons/implementation.html
//
// Many comments are lifted from those page and are copyright by Red Blob Games.

/////////////////////////////////////////////////////////////////////////////
// references
// * https://www.redblobgames.com/grids/hexagons/#references

/////////////////////////////////////////////////////////////////////////////
// geometry.go:
// geometry
//   * https://www.redblobgames.com/grids/hexagons/#basics
// spacing
//   * https://www.redblobgames.com/grids/hexagons/#spacing
// angles
//   * https://www.redblobgames.com/grids/hexagons/#angles

/////////////////////////////////////////////////////////////////////////////
// coords.go:
// coordinate systems
//   * https://www.redblobgames.com/grids/hexagons/#coordinates
// offset coordinates
//   * https://www.redblobgames.com/grids/hexagons/#coordinates-offset
// cube coordinates
//   * https://www.redblobgames.com/grids/hexagons/#coordinates-cube
// axial coordinates
//   * https://www.redblobgames.com/grids/hexagons/#coordinates-axial
// doubled coordinates
//   * https://www.redblobgames.com/grids/hexagons/#coordinates-doubled
// others
//   * https://www.redblobgames.com/grids/hexagons/#coordinates-other

/////////////////////////////////////////////////////////////////////////////
// conversions.go:
// coordinate conversions
//   * https://www.redblobgames.com/grids/hexagons/#conversions
// axial coordinates
//   * https://www.redblobgames.com/grids/hexagons/#conversions-axial
// offset coordinates
//   * https://www.redblobgames.com/grids/hexagons/#conversions-offset
// doubled coordinates
//   * https://www.redblobgames.com/grids/hexagons/#conversions-doubled

/////////////////////////////////////////////////////////////////////////////
// neighbors.go:
// neighbors
//   * https://www.redblobgames.com/grids/hexagons/#neighbors
// cube coordinates
//   * https://www.redblobgames.com/grids/hexagons/#neighbors-cube
// axial coordinates
//   * https://www.redblobgames.com/grids/hexagons/#neighbors-axial
// offset coordinates
//   * https://www.redblobgames.com/grids/hexagons/#neighbors-offset
// doubled coordinates
//   * https://www.redblobgames.com/grids/hexagons/#neighbors-doubled
// diagonals
//   * https://www.redblobgames.com/grids/hexagons/#neighbors-diagonal

/////////////////////////////////////////////////////////////////////////////
// distances.go:
// distances
//   * https://www.redblobgames.com/grids/hexagons/#distances
// cube coordinates
//   * https://www.redblobgames.com/grids/hexagons/#distances-cube
// axial coordinates
//   * https://www.redblobgames.com/grids/hexagons/#distances-axial
// offset coordinates
//   * https://www.redblobgames.com/grids/hexagons/#distances-offset
// doubled coordinates
//   * https://www.redblobgames.com/grids/hexagons/#distances-doubled

/////////////////////////////////////////////////////////////////////////////
// lines.go:
// line drawing
//   * https://www.redblobgames.com/grids/hexagons/#line-drawing

/////////////////////////////////////////////////////////////////////////////
// movement.go:
// movement range
//   * https://www.redblobgames.com/grids/hexagons/#range
// coordinate range
//   * https://www.redblobgames.com/grids/hexagons/#range-coordinate
// intersecting ranges
//   * https://www.redblobgames.com/grids/hexagons/#range-intersection
// obstacles
//   * https://www.redblobgames.com/grids/hexagons/#range-obstacles

/////////////////////////////////////////////////////////////////////////////
// rotation.go:
// rotation
//   * https://www.redblobgames.com/grids/hexagons/#rotation

/////////////////////////////////////////////////////////////////////////////
// reflection.go:
// reflection
//   * https://www.redblobgames.com/grids/hexagons/#reflection

/////////////////////////////////////////////////////////////////////////////
// rings.go:
// rings
//   * https://www.redblobgames.com/grids/hexagons/#rings
// single ring
//   * https://www.redblobgames.com/grids/hexagons/#rings-single
// spiral rings
//   * https://www.redblobgames.com/grids/hexagons/#rings-spiral
// spiral coordinates
//   * https://www.redblobgames.com/grids/hexagons/#rings-spiral-coordinates

/////////////////////////////////////////////////////////////////////////////
// field_of_view.go:
// field of view
//   * https://www.redblobgames.com/grids/hexagons/#field-of-view
// todo: implement

/////////////////////////////////////////////////////////////////////////////
// hex_to_pixel.go:
// hex to pixel
//   * https://www.redblobgames.com/grids/hexagons/#hex-to-pixel
// axial coordinates
//   * https://www.redblobgames.com/grids/hexagons/#hex-to-pixel-axial
// offset coordinates
//   * https://www.redblobgames.com/grids/hexagons/#hex-to-pixel-offset
// doubled coordinates
//   * https://www.redblobgames.com/grids/hexagons/#hex-to-pixel-doubled
// mod: non-zero origin
//   * https://www.redblobgames.com/grids/hexagons/#hex-to-pixel-mod-origin
// mod: pixel sizes
//   * https://www.redblobgames.com/grids/hexagons/#hex-to-pixel-mod-pixelsize

/////////////////////////////////////////////////////////////////////////////
// pixel_to_hex.go:
// pixel to hex
//   * https://www.redblobgames.com/grids/hexagons/#pixel-to-hex
// axial coordinates
//   * https://www.redblobgames.com/grids/hexagons/#pixel-to-hex-axial
// offset coordinates
//   * https://www.redblobgames.com/grids/hexagons/#pixel-to-hex-offset
// doubled coordinates
//   * https://www.redblobgames.com/grids/hexagons/#pixel-to-hex-doubled
// mod:non-zero origin
//   * https://www.redblobgames.com/grids/hexagons/#pixel-to-hex-mod-origin
// mod: pixel sizes
//   * https://www.redblobgames.com/grids/hexagons/#pixel-to-hex-mod-pixelsize

/////////////////////////////////////////////////////////////////////////////
// rounding.go:
// rounding to nearest hex
//   * https://www.redblobgames.com/grids/hexagons/#rounding

/////////////////////////////////////////////////////////////////////////////
// storage.go:
// map storage in axial coordinates
//   * https://www.redblobgames.com/grids/hexagons/#map-storage

/////////////////////////////////////////////////////////////////////////////
// wrapping.go:
// wraparound maps
//   * https://www.redblobgames.com/grids/hexagons/#wraparound

/////////////////////////////////////////////////////////////////////////////
// paths.go:
// pathfinding
//   * https://www.redblobgames.com/grids/hexagons/#pathfinding
