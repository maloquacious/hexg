// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package hexg implements hex grids in the style of the Red Blob Games blog
// https://www.redblobgames.com/grids/hexagons/ and
// https://www.redblobgames.com/grids/hexagons/implementation.html.
//
// Block Comments are lifted directly from the Red Blog Games page.
// Those comments are copyright by them and used without permission.
// My intention is to remove them after I get this package working
// and tested.
package hexg

import (
	"fmt"
	"log"
	"math"
)

/*
The [main page](https://www.redblobgames.com/grids/hexagons/) covers the
theory for hex grid algorithms and math. Now let’s write a library to
handle hex grids. The first thing to think about is what the core concepts
will be.

* Since most of the algorithms work with cube coordinates, I’ll need a
data structure for cube coordinates, along with algorithms that work with
them. I’ll call this the **Hex** class.

* For some games I want to show coordinates to the player, and those will
probably not be cube, but instead axial or offset, so I’ll need a data
structure for the player-visible coordinate system, as well as functions
for converting back and forth. Cube and axial are basically the same so
I’m not going to bother implementing a separate axial system, and I’ll
reuse **Hex**. For offset coordinates, I’ll make a separate data
structure, **Offset**.

* A grid map will likely need additional storage for terrain, objects,
units, etc. A 2D array can be used, but it’s not always straightforward,
so I’ll create a **Map** class for this.

* To draw hexes on the screen, I need a way to convert hex coordinates
into screen space. I’ll call this the **Layout** class. The main article
doesn’t cover some of the additional features I want:
  * Support y-axis pointing down (common in 2D libraries) as well as
y-axis pointing up (common in 3D libraries). The main article only covers
y-axis pointing down.
  * Support stretched or squashed hexes, which are common with pixel
graphics. The main article only supports equilateral hexes.
  * Support the 0,0 hex being located on the screen anywhere. The main
article always places the 0,0 hex at x=0, y=0.

* I also need a way to convert mouse clicks and other pixel coordinates
back into hex coordinates. I will put this into the **Layout** class. The
same things I need to deal with for hex to screen (y-axis direction,
stretch/squash, origin) have to be dealt with for screen to hex, so it
makes sense to put them together.

* The main article doesn’t distinguish hexes that have integer
coordinates from those with fractional coordinates. I’ll define a second
class **FractionalHex** for the two algorithms where I want to have
floating point coordinates: linear interpolation and rounding.
* Once I have coordinates and the `neighbors` function implemented I can
use all graph algorithms including movement range and pathfinding. I
cover pathfinding for graphs on
[another page](https://www.redblobgames.com/pathfinding/a-star/introduction.html)
and won’t duplicate that code here.
*/

/*
1.0 Hex coordinates

On the main page, I treat Cube and Axial systems separately.
Cube coordinates derived from x,y,z cartesian coordinates, and use three
axes q,r,s 120° apart, where q+r+s = 0. Axial coordinates have two axes
q,r that are 120° apart.

type Hex struct { // Cube storage
	q, r, s int
}

func NewHex(q, r, s int) Hex { // Cube constructor
	if q+r+s != 0 {
		panic("assert(q + r + s == 0)")
	}
	return Hex{q: q, r: r, s: s}
}

Pretty simple. Here’s a class that stores axial coordinates internally,
but uses cube coordinates for the interface:

type Hex struct { // Axial storage
	q, r int
}

func NewHex(q, r, s int) Hex { // Cube constructor
	if q+r+s != 0 {
		panic("assert(q + r + s == 0)")
	}
	return Hex{q: q, r: r}
}

These two classes are effectively equivalent. The first one stores `s`
explicitly and the second one uses accessors and calculates `s` when
needed. **Cube and Axial are essentially the same system**, so I’m not
going to write a separate class for each. However, **the labels on this
page are different**. On the main page, the axial/cube relationship is
q→x, r→z, but here I am making q→q, r→r. And that means on the main page
cube coordinates are (q, -q-r, r) but on this page cube coordinates are
(q, r, -q-r). This makes my two pages inconsistent and I plan to update
the main page to match this page.

Yet another style is to calculate s in the constructor instead of passing
it in:

type Hex struct { // Cube storage
	q, r, s int
}

func NewHex(q, r int) Hex { // Axial constructor
	return Hex{q: q, r: r, s: -q - r}
}

An advantage of the axial constructor style is that more than half the
time, you’re doing this anyway at the call site. You’ll have `q` and `r`
and not `s`, so you’ll pass in `-q-r` for the third parameter. You can
also combine this with the second style (axial storage), and store only
`q` and `r`, and calculate `s` in an accessor.

Yet another style is to use an array instead of named fields:

type Hex [3]int // Vector storage

func NewHex(q, r, s int) Hex { // Cube constructor
	if q+r+s != 0 {
		panic("assert(q + r + s == 0)")
	}
	return Hex{q, r, s}
}

func (v Hex) Q() int { return v[0] }
func (v Hex) R() int { return v[1] }
func (v Hex) S() int { return v[2] }

An advantage of this style is that you start seeing patterns where `q`,
`r`, `s` are all treated the same way. You can write loops to handle them
uniformly instead of duplicating code. You might use SIMD instructions on
the CPU. You might use `vec3` on the GPU. When you read the equality,
`hex_add`, `hex_subtract`, `hex_scale`, `hex_length`, `hex_round`, and
`hex_lerp` functions below, you’ll see how it might be useful to treat
the coordinates uniformly. When you read `hex_to_pixel` and
`pixel_to_hex`, you’ll see that vector and matrix operations (CPU or GPU)
can be used with hex coordinates when expressed this way.

Programming is full of tradeoffs. For this page, I want to focus on
simplicity and readability, not performance or abstraction, so I’m going
to use the first style: cube storage, cube constructor. I find it easiest
to understand the algorithms in this style. However, I like all of these
styles, and wouldn’t hesitate to choose any of them, as long as things
are consistent in the project. In a language with multiple constructors,
I’d include both the axial and cube constructors for convenience. In C++,
the `int` could instead be a template parameter. In C or C++11, the
`int v[]` style and the `int q, r, s` style can be merged with a union.
A template parameter `w` can also be used to distinguish between positions
and vectors. Putting all of these together:

type Hex interface {}

type CubeHex struct { // Cube storage
	q, r, s int
}

type VectorHex struct { // Vector storage
	v [3]int;
}

func NewHexAxial(q, r int, w bool) Hex {
	if w {
		return VectorHex{q, r, -q - r}
	}
	return CubeHex{q:q, r:r, s:-q - r}
};

func NewHexCube(q, r, s int, w bool) Hex {
	if q+r+s != 0 {
		panic("assert(q + r + s == 0)")
	}
	if w {
		return VectorHex{q, r, s}
	}
	return CubeHex{q:q, r:r, s:s}
};

I didn’t use this style on this page because I want to make translation
to other languages straightforward.

Another design alternative is to use the `x`, `y`, `z` names so that you
can make hex coordinates and cartesian coordinates reuse the same data
structures. If you’re already using a vector library for geometry, you
can reuse that for hex coordinates, and then use a matrix library for
representing hex-to-pixel and pixel-to-hex operations.
*/

// NewHex returns a Hex initialized from axial coordinates.
func NewHex(q, r int) Hex { // Axial constructor
	return Hex{q: q, r: r, s: -q - r}
}

// Hex stores the `q`, `r`, and `s` coordinates.
type Hex struct {
	q, r, s int
}

// String returns the coordinates formatted as "(q, r, s)".
func (h Hex) String() string {
	return fmt.Sprintf("(%d, %d, %d)", h.q, h.r, h.s)
}

// ConciseString returns the coordinates with signs.
// It returns the coordinates formatted as "+q+r+s".
func (h Hex) ConciseString() string {
	return fmt.Sprintf("%+d%+d%+d", h.q, h.r, h.s)
}

/*
1.1 Equality

Equality and inequality is straightforward: two hexes are equal if their
coordinates are equal.

func hex_eq(a, b Hex) bool {
	return a.q == b.q && b.r == b.r && b.s == b.s
}

func hex_neq(a, b Hex) bool {
	return !(hex_eq(a, b))
}
*/

// Equals returns true if the two hexes have the same coordinates.
func (a Hex) Equals(b Hex) bool {
	return a.q == b.q && b.r == b.r && b.s == b.s
}

// NotEquals returns true if the two hexes have different coordinates.
func (a Hex) NotEquals(b Hex) bool {
	return !a.Equals(b)
}

/*
1.2 Coordinate arithmetic

Since cube coordinates come from 3D cartesian coordinates, I
automatically get things like addition, subtraction,
multiplication, and division. For example, you can have
`Hex(2, 0, -2)` that represents two steps northeast, and add that to
location `Hex(3, -5, 2)` the obvious way: `Hex(2 + 3, 0 + -5, -2 + -2)`.
With other coordinate systems like offset coordinates, you can’t do that
and get what you want. These operations are what you’d implement with 3D
cartesian vectors, but I am using `q`, `r`, `s` names in this class
instead of `x`, `y`, `z`:

func hex_add(a, b Hex) Hex {
	return Hex{q: a.q + b.q, r: a.r + b.r, s: a.s + b.s}
}

func hex_subtract(a, b Hex) Hex {
	return Hex{q: a.q - b.q, r: a.r - b.r, s: a.s - b.s}
}

func hex_multiply(a Hex, k int) Hex {
	return Hex{q: a.q * k, r: a.r * k, s: a.s * k}
}
*/

// Add returns the sum of two hexes.
func (a Hex) Add(b Hex) Hex {
	return Hex{q: a.q + b.q, r: a.r + b.r, s: a.s + b.s}
}

// Subtract returns the difference of two hexes.
func (a Hex) Subtract(b Hex) Hex {
	return Hex{q: a.q - b.q, r: a.r - b.r, s: a.s - b.s}
}

// Multiply returns a Hex scaled by an integer.
func (a Hex) Multiply(k int) Hex {
	return Hex{q: a.q * k, r: a.r * k, s: a.s * k}
}

/*
An alternate design is to reuse an existing vec3 class from your geometry
library to represent axial/cube coordinates, and in that case you won’t
have to write these functions.
*/

/*
1.3 Distance

The distance between two hexes is the length of the line between them.
Both the distance and length operations can come in handy. It looks like
[the distance function from the main article](https://www.redblobgames.com/grids/hexagons/#distances):

func hex_length(hex Hex) int {
	return (abs(hex.q) + abs(hex.r) + abs(hex.s)) / 2
}

func hex_distance(a, b Hex) int {
	return hex_length(hex_subtract(a, b))
}
*/

// Length is the distance from the origin to the hex.
func (hex Hex) Length() int {
	return (abs(hex.q) + abs(hex.r) + abs(hex.s)) / 2
}

// Distance between two hexes is the length of the line between them.
func (a Hex) Distance(b Hex) int {
	return a.Subtract(b).Length()
}

/*
1.3.1 Neighbors

With distance, I defined two functions: length works on one argument and
distance works with two. The same is true with neighbors. The direction
function is with one argument and the neighbor function is with two. It
looks like
[the neighbors function from the main article](https://www.redblobgames.com/grids/hexagons/#neighbors):

var hex_directions = [6]Hex{
	{1, 0, -1}, {1, -1, 0}, {0, -1, 1},
	{-1, 0, 1}, {-1, 1, 0}, {0, 1, -1},
}

func hex_direction(direction int) Hex {
	if !(0 <= direction && direction < 6) {
		panic("assert (0 <= direction && direction < 6)")
	}
	return hex_directions[direction]
}

func hex_neighbor(hex Hex, direction int) Hex {
	return hex_add(hex, hex_direction(direction))
}

To make directions outside the range 0..5 work, use
`hex_directions[(6 + (direction % 6)) % 6].` Yeah, it’s ugly, but it will
work with negative directions too.
*/

var hex_directions = [6]Hex{
	{1, 0, -1}, {1, -1, 0}, {0, -1, 1},
	{-1, 0, 1}, {-1, 1, 0}, {0, 1, -1},
}

// HexDirection returns the vector (q, r, and s offsets) to use based on the
// direction, which must be in the range 0..5. It will panic of the direction
// causes an out of bounds on the hex_directions slice.
func HexDirection(direction int) Hex {
	return hex_directions[direction]
}

// Neighbor returns the hex that is one step away in the given direction.
// Direction is coerced to the range 0..5.
func (hex Hex) Neighbor(direction int) Hex {
	direction = (6 + (direction % 6)) % 6
	return hex.Add(HexDirection(direction))
}

/*
2.0 Layout

The next major piece of functionality I need is a way to convert between
hex coordinates and screen coordinates. There’s a pointy top layout and a
flat top hex layout. The conversion uses a matrix as well as the inverse
of the matrix, so I need a way to store those. Also, for drawing the
corners, pointy top starts at 30° and flat top starts at 0°, so I need a
place to store that too.

I’m going to define an **Orientation** helper class to store these:
the 2×2 forward matrix, the 2×2 inverse matrix, and the starting angle:
*/

type orientation struct {
	f0, f1, f2, f3 float64
	b0, b1, b2, b3 float64
	start_angle    float64 // in multiples of 60°
}

/*
There are only two orientations, so I’m going to make constants for them:
*/

type Orientation int

const (
	// LayoutPointy is a pointy top orientation:
	// * staggered columns
	// * horizontal rows
	LayoutPointy Orientation = iota

	// LayoutFlat is a flat top orientation:
	// * vertical columns
	// * staggered rows
	LayoutFlat
)

// LayoutOffset describes the layout (horizontal row or vertical columns) and
// which way odd rows and columns are pushed or pulled
type LayoutOffset int

const (
	// OddR is a pointy-top layout with horizontal rows that shoves odd rows right
	OddR LayoutOffset = iota
	// EvenR is a pointy-top layout with horizontal rows that shoves even rows right
	EvenR
	// OddQ is a flat-top layout with vertical columns that shoves odd columns down
	OddQ
	// EvenQ is a flat-top layout with vertical columns that shoves even columns down
	EvenQ
)

/*
Now I will use them for the **Layout** class:
*/

// Layout implements the conversion between hex (q,r,s) and screen (x,y)
// coordinates. The conversion uses a matrix as well as the inverse of
// the matrix for each orientation.
type Layout struct {
	orientation orientation
	offset      LayoutOffset

	// origin is center of the (0,0,0) hex. It is used in translate
	// transformation to move hexes on the screen. If you don’t need
	// this, set it to Point(0, 0).
	origin Point

	// size is used for non-uniform scaling, especially for matching
	// pixel sprite sizes. It’s a scale transform. If you need uniform
	// scaling, set it to Point(width, height), where width == height.
	size Point
}

// NewLayout returns a layout with either flat-top or pointy-top hexes.
//
// Flat-top layouts have vertical columns and staggered rows.
// Drawing corners start at  0°.
//
// Pointy-top layouts have staggered columns and horizontal rows.
// Drawing corners start at 30°.
func NewLayout(offset LayoutOffset, size, origin Point) Layout {
	switch offset {
	case OddR: // OddR is a pointy-top layout with horizontal rows that shoves odd rows right
		return Layout{
			offset: OddR,
			orientation: orientation{
				f0: math.Sqrt(3.0), f1: math.Sqrt(3.0) / 2.0, f2: 0.0, f3: 3.0 / 2.0,
				b0: math.Sqrt(3.0) / 3.0, b1: -1.0 / 3.0, b2: 0.0, b3: 2.0 / 3.0,
				start_angle: 0.5,
			},
			origin: origin,
			size:   size,
		}
	case EvenR: // EvenR is a pointy-top layout with horizontal rows that shoves even rows right
		return Layout{
			offset: EvenR,
			orientation: orientation{
				f0: math.Sqrt(3.0), f1: math.Sqrt(3.0) / 2.0, f2: 0.0, f3: 3.0 / 2.0,
				b0: math.Sqrt(3.0) / 3.0, b1: -1.0 / 3.0, b2: 0.0, b3: 2.0 / 3.0,
				start_angle: 0.5,
			},
			origin: origin,
			size:   size,
		}
	case OddQ: // OddQ is a flat-top layout with vertical columns that shoves odd columns down
		return Layout{
			offset: OddQ,
			orientation: orientation{
				f0: 3.0 / 2.0, f1: 0.0, f2: math.Sqrt(3.0) / 2.0, f3: math.Sqrt(3.0),
				b0: 2.0 / 3.0, b1: 0.0, b2: -1.0 / 3.0, b3: math.Sqrt(3.0) / 3.0,
				start_angle: 0.0,
			},
			origin: origin,
			size:   size,
		}
	case EvenQ: // EvenQ is a flat-top layout with vertical columns that shoves even columns down
		return Layout{
			offset: EvenQ,
			orientation: orientation{
				f0: 3.0 / 2.0, f1: 0.0, f2: math.Sqrt(3.0) / 2.0, f3: math.Sqrt(3.0),
				b0: 2.0 / 3.0, b1: 0.0, b2: -1.0 / 3.0, b3: math.Sqrt(3.0) / 3.0,
				start_angle: 0.0,
			},
			origin: origin,
			size:   size,
		}
	}
	panic("assert(offset in (OddR, EvenR, OddQ, EvenQ)")
}

func (layout Layout) IsEven() bool {
	return layout.offset == EvenQ || layout.offset == EvenR
}

func (layout Layout) IsFlat() bool {
	return layout.offset == OddQ || layout.offset == EvenQ
}

func (layout Layout) IsHorizontal() bool {
	return layout.offset == OddR || layout.offset == EvenR
}

func (layout Layout) IsOdd() bool {
	return layout.offset == OddQ || layout.offset == OddR
}

func (layout Layout) IsPointy() bool {
	return layout.offset == OddR || layout.offset == EvenR
}

func (layout Layout) IsVertical() bool {
	return layout.offset == OddQ || layout.offset == EvenQ
}

/*
Oh, hm, I guess I need a minimal Point class. If your graphics/geometry
library already provides one, use that.
*/

// Point represents a screen coordinate
type Point struct {
	X float64
	Y float64
}

/*
Side note: observe how many of these are arrays of numbers underneath.
* Hex is int[3].
* Orientation is an angle, a double, and two matrices, each double[4] or double[2][2].
* Point is double[2].
* Layout is an Orientation and two Points.
Later on the page,
* FractionalHex is double[3], and
* OffsetCoord is int[2].
I use structs instead of arrays of numbers because giving a name to
things helps me understand them, and also helps with type checking.
However, an alternate design choice is to reuse a standard vector library
for all of these types, and then use standard matrix multiply for the
layout. You can then use your library’s matrix inverse to calculate the
inverse matrix. Using arrays of numbers (or a numeric array class)
instead of separate structs with names will allow you to reuse more code.
I chose not to do that, but I think it’s a reasonable choice.

Ok, now I’m ready to write the layout algorithms.
*/

/*
2.1 Hex to screen

The main article has two versions of axial hex-to-pixel, one for each
orientation. The code is essentially the same except the numbers are
different, so for this implementation I’ve put the numbers from the
matrix into the Orientation class, as `f0` through `f3`:

func hex_to_pixel(layout Layout, h Hex) Point {
	M := layout.orientation
	x := (M.f0*float64(h.q) + M.f1*float64(h.r)) * layout.size.X
	y := (M.f2*float64(h.q) + M.f3*float64(h.r)) * layout.size.Y
	return Point{
		X: x + layout.origin.X,
		Y: y + layout.origin.Y,
	}
}
*/

// HexToPixel returns the origin of the hex on the grid as a Point.
func (layout Layout) HexToPixel(h Hex) Point {
	M := layout.orientation
	x := (M.f0*float64(h.q) + M.f1*float64(h.r)) * layout.size.X
	y := (M.f2*float64(h.q) + M.f3*float64(h.r)) * layout.size.Y
	return Point{
		X: x + layout.origin.X,
		Y: y + layout.origin.Y,
	}
}

/*
The main article has two optional steps:

* Non-zero origin representing the center of the q=0,r=0 hexagon.
I store this in `layout.origin`. It’s a translate transformation.
If you don’t need this, set it to Point(0, 0).

* Non-uniform scaling, especially for matching pixel sprite sizes.
I store this in `layout.size`. It’s a scale transform.
If you need uniform scaling, set it to Point(size, size).

I’ll show some uses of these in the 2.4 section below.
*/

/*
2.2 Screen to hex

The main article has two versions of axial pixel-to-hex, one for each
orientation. Again, the code is the same except for the numbers, which
are the inverse of the matrix. I put the matrix inverse into the
Orientation class, as `b0` through `b3`, and used it here. In the
forward direction, to go from hex coordinates to screen coordinates
I first multiply by the matrix, then multiply by the size, then add
the origin. To go in the reverse direction, I have to undo these.
First undo the origin by subtracting it, then undo the size by dividing
by it, then undo the matrix multiply by multiplying by the inverse:

func pixel_to_hex_fractional(layout Layout, p Point) FractionalHex {
	M := layout.orientation
	pt := Point{
		X: (p.X - layout.origin.X) / layout.size.X,
		Y: (p.Y - layout.origin.Y) / layout.size.Y,
	}
	q := M.b0*pt.X + M.b1*pt.Y
	r := M.b2*pt.X + M.b3*pt.Y
	return FractionalHex{q: q, r: r, s: -q - r}
}

There’s a complication here: I start with integer hex coordinates to go
to screen coordinates, but when going in reverse, I have no guarantee
that the screen location will be exactly at a hexagon center. Instead of
getting back an integer hex coordinate, I get back a floating point value
(type double), which means I return a FractionalHex instead of a Hex. To
get back to the integer, I need to round it to the nearest hex. I’ll
implement that in a bit.
*/

// PixelToFractionalHex returns the fractional hex that encloses the pixel.
// In theory, the origin of that fractional hex will be the pixel.
func (layout Layout) PixelToFractionalHex(p Point) FractionalHex {
	M := layout.orientation
	pt := Point{
		X: (p.X - layout.origin.X) / layout.size.X,
		Y: (p.Y - layout.origin.Y) / layout.size.Y,
	}
	q := M.b0*pt.X + M.b1*pt.Y
	r := M.b2*pt.X + M.b3*pt.Y
	return FractionalHex{q: q, r: r, s: -q - r}
}

/*
2.3 Drawing a hex

To draw a hex, I need to know where each corner is relative to the center
of the hex.

* With the flat top orientation, the corners are at 0°, 60°, 120°, 180°, 240°, 300°.
* With pointy top, they’re at 30°, 90°, 150°, 210°, 270°, 330°.

I encode that in the Orientation class’s start_angle value, either 0.0
for 0° or 0.5 for 30°.

Once I know where the corners are relative to the center, I can calculate
the corners in screen locations by adding the center to each corner, and
putting the coordinates into an array.

func hex_corner_offset(layout Layout, corner int) Point {
	size := layout.size
	angle := 2.0 * math.Pi *
		(layout.orientation.start_angle + float64(corner)) / 6
	return Point{
		X: size.X * math.Cos(angle),
		Y: size.Y * math.Sin(angle),
	}
}

func polygon_corners(layout Layout, h Hex) [6]Point {
	var corners [6]Point
	center := hex_to_pixel(layout, h)
	for i := 0; i < 6; i++ {
		offset := hex_corner_offset(layout, i)
		corners[i] = Point{
			X: center.X + offset.X,
			Y: center.Y + offset.Y,
		}
	}
	return corners
}
*/

// HexCornerOffset returns the screen location (pixel) of a corner of a hex on the grid.
func (layout Layout) HexCornerOffset(corner int) Point {
	size := layout.size
	angle := 2.0 * math.Pi *
		(layout.orientation.start_angle + float64(corner)) / 6
	return Point{
		X: size.X * math.Cos(angle),
		Y: size.Y * math.Sin(angle),
	}
}

// PolygonCorners returns the location of the six corners of the hex on the grid.
func (layout Layout) PolygonCorners(h Hex) [6]Point {
	var corners [6]Point
	center := layout.HexToPixel(h)

	for i := 0; i < 6; i++ {
		offset := layout.HexCornerOffset(i)
		corners[i] = Point{
			X: center.X + offset.X,
			Y: center.Y + offset.Y,
		}
	}
	return corners
}

/*
2.4 Layout examples

Ok, let’s try it out! I have written Hex, Orientation, Layout, and Point
and the functions that go with each. That’s enough for me to draw hexes.
I’m going to use the Go version of these functions to draw some hexes in
the browser.

Let’s try the two orientations, `layout_pointy` and `layout_flat`:

Let’s try uniform scaling of size: `Point(10, 10)`, `Point(25, 25)`, and
`Point(50, 50)`. When the two values are the same, we get regular hexagons.
The `size` value is half the height for pointy-top hexagons and half the
width for flat-top hexagons.

But sometimes we want to stretch the hexagons to **fit sprite assets**,
so my `size` has a separate `x` and `y` scaling. From the main page, we
can use these calculations:

* for flat top art sprites `W✕H`, set `size` to `Point(W/2, H/sqrt(3))`.
The example fits `100✕100` sprites that are slightly taller.
* for pointy top art sprites `W✕H`, set `size` to `Point(W/sqrt(3), H/2)`.
The example fits `100✕100` sprites that are slightly wider.

Another thing we can do with size is to **flip the r axis**. Compare
`size` set to `Point(25, 25)` and set to `Point(25, -25)`. This is also
useful if your y-axis grows upwards, as you can choose whether to make
the `r` coordinate grow upwards with `y` (positive `size.y`) or downwards
opposite of `y` (negative `size.y`).

The `origin` is occasionally useful too. I usually set it to
`Point(0, 0)`. That puts the center of the q=0,r=0 hexagon at x=0,y=0.
But if we want the **top left** of that hexagon to be at x=0,y=0 then:

* for flat top hexes, set origin to Point(size.x, size.y * sqrt(3)/2):

* for pointy top hexes, set origin to Point(size.x * sqrt(3)/2, size.y):

I think the above diagrams are a reasonable set of tests for the
orientation, size, and origin. It shows that the Layout class can handle
a wide variety of needs, without having to make different variants of the
Hex class.

An alternate (simpler) implementation would be to always set origin to
0,0 and set size to 1,1. Then chain simpler transforms together by using
vector operations on the cartesian coordinates:

* **hex→pixel**
  * first hex→cartesian, then scale the cartesian coordinate by multiplying by the desired scale, and then translate it to the desired origin.

* **pixel→hex**
  * first undo the translate by subtracting the origin, then undo the scale by dividing by the scale, then run cartesian→hex.
*/

func LayoutPointExample() {
	for q := -4; q <= 4; q++ {
		for r := -2; r <= 2; r++ {
			h := NewHex(q, r)
			log.Printf("Hex{q: %d, r: %d, s: %d}\n", h.q, h.r, h.s)
		}
	}
}

func LayoutFlatExample() {
	for q := -4; q <= 4; q++ {
		for r := -2; r <= 2; r++ {
			h := NewHex(q, r)
			log.Printf("Hex{q: %d, r: %d, s: %d}\n", h.q, h.r, h.s)
		}
	}
}

/*
3.0 Fractional Hex

For pixel-to-hex I need fractional hex coordinates. It looks like the
**Hex** class, but uses double instead of int:

type FractionalHex struct {
	q, r, s float64
}
*/

// NewFractionalHex returns a FractionalHex initialized from cube coordinates.
func NewFractionalHex(q, r, s float64) FractionalHex { // Axial constructor
	return FractionalHex{q: q, r: r, s: s}
}

// FractionalHex implements a hex with floating point coordinates.
// Used for linear interpolation and rounding.
type FractionalHex struct {
	q, r, s float64
}

/*
3.1 Hex rounding

Rounding turns a fractional hex coordinate into the nearest integer hex
coordinate. The algorithm is straight out of the main article:

func hex_round(h FractionalHex) Hex {
    q := int(math.Round(h.q));
    r := int(math.Round(h.r));
    s := int(math.Round(h.s));
    q_diff := math.Abs(float64(q) - h.q);
    r_diff := math.Abs(float64(r) - h.r);
    s_diff := math.Abs(float64(s) - h.s);
    if (q_diff > r_diff and q_diff > s_diff) {
        q = -r - s;
    } else if (r_diff > s_diff) {
        r = -q - s;
    } else {
        s = -q - r;
    }
    return Hex{q:q, r:r, s:s};
}
*/

// HexRound turns a fractional hex coordinate into the nearest integer
// hex coordinate.
func (h FractionalHex) HexRound() Hex {
	q := int(math.Round(h.q))
	r := int(math.Round(h.r))
	s := int(math.Round(h.s))
	q_diff := math.Abs(float64(q) - h.q)
	r_diff := math.Abs(float64(r) - h.r)
	s_diff := math.Abs(float64(s) - h.s)
	if q_diff > r_diff && q_diff > s_diff {
		q = -r - s
	} else if r_diff > s_diff {
		r = -q - s
	} else {
		s = -q - r
	}
	return Hex{q: q, r: r, s: s}
}

/*
In the Screen to hex section I wrote a function that turns a pixel
coordinate into a fractional hex coordinate. Rounding turns that
into a regular hex coordinate:

Hex pixel_to_hex_rounded(Layout layout, Point p) {
    return hex_round(pixel_to_hex_fractional(layout, p));
}
*/

// PixelToHexRounded turns a fractional hex into a regular hex coordinate:
func (layout Layout) PixelToHexRounded(p Point) Hex {
	return layout.PixelToFractionalHex(p).HexRound()
}

/*
3.2 Line drawing

To draw a line, I linearly interpolate between two hexes, and then round
it to the nearest hex. To linearly interpolate between hex coordinates I
linearly interpolate each of the components (q, r, s) independently:

float lerp(double a, double b, double t) {
    return a * (1-t) + b * t;
    // better for floating point precision than
    //   a + (b - a) * t
    // which is what I usually write
}

FractionalHex hex_lerp(Hex a, Hex b, double t) {
    return FractionalHex(lerp(a.q, b.q, t),
    lerp(a.r, b.r, t),
    lerp(a.s, b.s, t));
}
*/

// lerp returns the linear interpolation of points on the line between two hexes
func lerp(a, b, t float64) float64 {
	return a*(1-t) + b*t // better for floating point precision than a + (b - a) * t
}

// HexLerp returns the linear interpolation of points on the line between two hexes
func (h Hex) HexLerp(b Hex, t float64) FractionalHex {
	return FractionalHex{
		q: lerp(float64(h.q), float64(b.q), t),
		r: lerp(float64(h.r), float64(b.r), t),
		s: lerp(float64(h.s), float64(b.s), t),
	}
}

// HexLerp returns the linear interpolation of points on the line between two hexes
func (h FractionalHex) HexLerp(b FractionalHex, t float64) FractionalHex {
	return FractionalHex{
		q: lerp((h.q), (b.q), t),
		r: lerp((h.r), (b.r), t),
		s: lerp((h.s), (b.s), t),
	}
}

/*
Line drawing is not too bad once I have linear interpolation:

vector<Hex> hex_linedraw(Hex a, Hex b) {
    int N = hex_distance(a, b);
    vector<Hex> results = {};
    double step = 1.0 / max(N, 1);
    for (int i = 0; i <= N; i++) {
        results.push_back(hex_round(hex_lerp(a, b, step * i)));
    }
    return results;
}

I needed to stick that max(N, 1) bit in there to handle lines with
length 0 (when A == B).

Sometimes the hex_lerp will output a point that’s on an edge. On some
systems, the rounding code will push that to one side or the other,
somewhat unpredictably and inconsistently. To make it always push these
points in the same direction, add an “epsilon” value to a. This will
“nudge” things in the same direction when it’s on an edge, and leave
other points unaffected.

vector<Hex> hex_linedraw(Hex a, Hex b) {
    int N = hex_distance(a, b);
    FractionalHex a_nudge(a.q + 1e-6, a.r + 1e-6, a.s - 2e-6);
    FractionalHex b_nudge(b.q + 1e-6, b.r + 1e-6, b.s - 2e-6);
    vector<Hex> results = {};
    double step = 1.0 / max(N, 1);
    for (int i = 0; i <= N; i++) {
        results.push_back(
            hex_round(hex_lerp(a_nudge, b_nudge, step * i)));
    }
    return results;
}

The nudge is not always needed. You might try without it first.
*/

// HexLinedraw returns the set of hexes that are between two hexes.
// Enable nudging to push points on an edge in a consistent direction.
func (a Hex) HexLinedraw(b Hex, withNudge bool) []Hex {
	N := a.Distance(b)
	var results []Hex
	var step float64
	if N == 0 {
		step = 1.0
	} else {
		step = 1.0 / float64(N)
	}
	if withNudge {
		a_nudge := FractionalHex{q: float64(a.q) + 1e-6, r: float64(a.r) + 1e-6, s: float64(a.s) - 2e-6}
		b_nudge := FractionalHex{q: float64(b.q) + 1e-6, r: float64(b.r) + 1e-6, s: float64(b.s) - 2e-6}
		for i := 0; i <= N; i++ {
			results = append(results, a_nudge.HexLerp(b_nudge, step*float64(i)).HexRound())
		}
		return results
	}
	for i := 0; i <= N; i++ {
		results = append(results, a.HexLerp(b, step*float64(i)).HexRound())
	}
	return results
}

/*
4.0 Map storage

There are two related problems to solve:

* how to generate a shape and
* how to store map data.

Let’s start with storing map data.

4.1 Map storage

The simplest way to store a map is to use a hash table. In C++, in order
to use unordered_map<Hex,_> or unordered_set<Hex> I need to define a hash
function for Hex. It would’ve been nice if C++ made it easier to define
this, but it’s not too bad. I hash the q and r fields (I can skip s
because it’s redundant), and combine them using the algorithm from
Boost’s hash_combine:

namespace std {
    template <> struct hash<Hex> {
        size_t operator()(const Hex& h) const {
            hash<int> int_hash;
            size_t hq = int_hash(h.q);
            size_t hr = int_hash(h.r);
            return hq ^ (hr + 0x9e3779b9 + (hq << 6) + (hq >> 2));
        }
    };
}
*/

// HashTable is a map of Hex indexed by the hash of the Hex.
type HashTable map[uint64]Hex

/*
Here’s an example of making a map with a float height at each hex:

type cell struct{
    hex Hex
    height float64
}
heights := make(map[uint64]call)
heights[Hex{q:1, r:-2, s:3}] = cell{Hex{q:1, r:-2, s:3}, 4.3};
fmt.Println(heights[Hex{q:1, r:-2, s:3}]);

The hash table by itself isn’t that useful. I need to combine it with
something that creates a map shape. In graph terms, I need something
that creates the nodes.
*/

// Key returns a hashable uint64 value derived from the axial
// coordinates (q, r) using a variation of Boost's hash_combine
// and MurmurHash3 finalization constants.
//
// The s coordinate is omitted because it is redundant in cube
// coordinates (s = -q - r).
//
// Casting to int64 before converting to uint64 preserves the signed
// bit pattern of negative values, maintaining good distribution
// across the entire hex grid, including negative coordinates.
func Key(q, r int) uint64 {
	const c1 = 0x9E3779B97F4A7C15 // golden ratio
	const c2 = 0xBF58476D1CE4E5B9
	const c3 = 0x94D049BB133111EB

	q64 := uint64(int64(q))
	r64 := uint64(int64(r))

	z := q64 ^ (r64 + c1 + (q64 << 6) + (q64 >> 2))
	z = (z ^ (z >> 30)) * c2
	z = (z ^ (z >> 27)) * c3
	return z ^ (z >> 31)
}

// Hash returns a hashable value for a Hex, derived from the
// q and r coordinates. It delegates to Key(q, r) to avoid
// duplication.
func (h Hex) Hash() uint64 {
	return Key(h.q, h.r)
}

// example using hex_hash to create a map of floats keyed by hex
// var heights map[uint64]float64
// heights[new_hex(1, -2, 3).Hash()] = 4.3

/*
4.2 Map shapes

In this section I write some loops that will produce various shapes of
maps. You can use these loops to make a set of hex coordinates for your
map, or fill in a map data structure, or iterate over the locations in
the map. I’ll write sample code that fills in a set of hex coordinates.

4.2.1 Parallelograms

With axial/cube coordinates, a straightforward loop over coordinates will
produce a parallelogram map instead of a rectangular one.

unordered_set<Hex> map;
for (int q = q1; q <= q2; q++) {
    for (int r = r1; r <= r2; r++) {
        map.insert(Hex(q, r, -q-r)));
    }
}

There are three coordinates, and the loop requires you choose any two of
them: (q,r), (s,q), or (r,s) lead to these pointy top maps, respectively:

And these flat top maps:

*/

func (layout Layout) ParallelogramQR(q1, r1 int, q2, r2 int) HashTable {
	gs := make(map[uint64]Hex)
	for q := q1; q <= q2; q++ {
		for r := r1; r <= r2; r++ {
			hex := Hex{q: q, r: r, s: -q - r}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

func (layout Layout) ParallelogramQS(q1, s1 int, q2, s2 int) HashTable {
	gs := make(map[uint64]Hex)
	for q := q1; q <= q2; q++ {
		for s := s1; s <= s2; s++ {
			hex := Hex{q: q, r: -q - s, s: s}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

func (layout Layout) ParallelogramRS(r1, s1 int, r2, s2 int) HashTable {
	gs := make(map[uint64]Hex)
	for r := r1; r <= r2; r++ {
		for s := s1; s <= s2; s++ {
			hex := Hex{q: -r - s, r: r, s: s}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

/*
4.2.2 Triangles

There are two directions for triangles to face, and the loop depends on
which direction you use. Assuming the y axis points down, with pointy top
these triangles face south/northwest/northeast, and with flat top these
triangles face east/northwest/southwest.

unordered_set<Hex> map;
for (int q = 0; q <= map_size; q++) {
    for (int r = 0; r <= map_size - q; r++) {
        map.insert(Hex(q, r, -q-r));
    }
}
*/

// TriangleUpDown returns a grid originating at (0,0,0).
// `map_size` is the length of a side.
func (layout Layout) TriangleUpDown(map_size int) HashTable {
	gs := HashTable{}
	for q := 0; q <= map_size; q++ {
		for r := 0; r <= map_size-q; r++ {
			hex := Hex{q: q, r: r, s: -q - r}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

/*
With pointy top these triangles face north/southwest/southeast and with
flat top these triangles face west/northeast/southeast:

unordered_set<Hex> map;
for (int q = 0; q <= map_size; q++) {
    for (int r = map_size - q; r <= map_size; r++) {
        map.insert(Hex(q, r, -q-r));
    }
}
*/

// TriangleLeftRight returns a grid originating at (0,0,0).
// `map_size` is the length of a side.
func (layout Layout) TriangleLeftRight(map_size int) HashTable {
	gs := HashTable{}
	for q := 0; q <= map_size; q++ {
		for r := map_size - q; r <= map_size; r++ {
			hex := Hex{q: q, r: r, s: -q - r}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

/*
If you flip your y-axis, then it’ll switch north and south here, as you
might expect.
*/

/*
4.2.3 Hexagons

Generating a hexagonal shape map is described on the main page.

unordered_set<Hex> map;
for (int q = -N; q <= N; q++) {
    int r1 = max(-N, -q - N);
    int r2 = min( N, -q + N);
    for (int r = r1; r <= r2; r++) {
        map.insert(Hex(q, r, -q-r));
    }
}

Here’s what I get for pointy top and flat top orientations:
*/

// Hexagon returns a grid centered about (0,0,0).
// does not depend on the orientation of the grid.
func Hexagon(radius int) HashTable {
	gs := HashTable{}
	N := radius
	for q := -N; q <= N; q++ {
		r1 := max(-N, -q-N)
		r2 := min(N, -q+N)
		for r := r1; r <= r2; r++ {
			hex := Hex{q: q, r: r, s: -q - r}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

/*
4.2.4 Rectangles

With axial/cube coordinates, getting rectangular maps is a little
trickier! The main article gives a clue, but I don’t actually show
the code. The code depends on whether using flat top or pointy top
hexes. The trick is to loop over offset coordinates and then convert
those to axial. Let’s start with pointy top hexes:

unordered_set<Hex> map;
for (int r = top; r <= bottom; r++) { // pointy top
    int r_offset = floor(r/2.0); // or r>>1
    for (int q = left - r_offset; q <= right - r_offset; q++) {
        map.insert(Hex(q, r, -q-r));
    }
}

That loop can produce grids like these:

The left/right/top/bottom are essentially offset coordinates, as offset
coordinates are a more natural fit for rectangular maps.
*/

/*
How about flat top hexes?

unordered_set<Hex> map;
for (int q = left; q <= right; q++) { // flat top
    int q_offset = floor(q/2.0); // or q>>1
    for (int r = top - q_offset; r <= bottom - q_offset; r++) {
        map.insert(Hex(q, r, -q-r));
    }
}

You might also need to experiment to get exactly the map you want.
Try setting the offset to floor((q+1)/2.0) or floor((q-1)/2.0) instead
of floor(q/2.0) for example, and the boundary will change slightly.
*/

// Rectangle returns a grid centered about (0,0,0).
// the internal logic depends on the orientation of the grid.
func (layout Layout) Rectangle(left, right, top, bottom int) HashTable {
	gs := HashTable{}
	if layout.IsPointy() {
		for r := top; r <= bottom; r++ {
			r_offset := r >> 1 // or math.Floor(float64(r) / 2.0)
			for q := left - r_offset; q <= right-r_offset; q++ {
				hex := Hex{q: q, r: r, s: -q - r}
				gs[hex.Hash()] = hex
			}
		}
		return gs
	}
	for q := left; q <= right; q++ {
		q_offset := q >> 1 // or math.Floor(float64(q) / 2.0)
		for r := top - q_offset; r <= bottom-q_offset; r++ {
			hex := Hex{q: q, r: r, s: -q - r}
			gs[hex.Hash()] = hex
		}
	}
	return gs
}

/*
4.3 Optimized storage

The hash table approach is pretty generic and works with any shape of
map, including weird shapes and shapes with holes. You can view it as
a type of node-and-edge graph structure, storing the nodes but explicitly
but calculating the edges on the fly with the hex_neighbor function.
*/

/*
A different way to store the node-and-edge graph structure is to
calculate all the edges ahead of time and store them explicitly. Give
each node an integer id and then use an array of arrays to store
neighbors. Or make each node an object and use a field to store a
list of neighbors. These graph structures are also generic and work
with any shape of map. You can also use any graph algorithm on them,
such as movement range, distance map, or pathfinding. Storing the edges
implicitly works well when the map is regular or is being edited; storing
them explicitly can work well when the map is irregularly shaped
(boundary, walls, holes) and isn’t changing frequently.

Some map shapes also allow a compact 2D or 1D array. The main article
gives a visual explanation. Here, I’ll give an explanation based on code.
The main idea is that for all the map shapes, there is a nested loop of
the form

for (int a = a1; a < a2; a++) {
    for (int b = b1; b < b2; b++) {
        ...
    }
}
For compact map storage, I’ll make an array of arrays, and index it with
array[a-a1][b-b1]. I subtract where the loop starts so that the first
index will be 0. For example, here’s the code for a rectangular shape
with pointy top hexes: (for flat top hexes, the loop is different)

for (int r = top; r <= bottom; r++) {
    int r_offset = floor(r/2.0);
    for (int q = left - r_offset; q <= right - r_offset; q++) {
        map.insert(Hex(q, r, -q-r));
    }
}
For pointy top hexes, variable a is r, and b is q. Value a1 (where the
r loop starts) is top and b1 (where the q loop starts) is
left - floor(r/2.0). That means the array will be indexed
array[r-top][q-(left-floor(r/2.0))] which simplifies to
array[r-top][q-left+floor(r/2.0)].

Note that floor(r/2.0) can be written r>>1.

The second thing I need to know is the size of the arrays. I need
a2-a1 arrays, and the size of each should be b2-b1. Be sure to check
for off-by-1 errors: if the loop is written a <= a2 then you’ll want
a2-a1+1 arrays, and similarly for b <= b2. I can build these arrays
using C++ vectors using this pattern:

vector<vector<T>> map(a2-a1);
for (int a = a1; a < a2; a++) {
    map.emplace_back(b2-b1);
}

For the rectangle example, a2-a1 becomes bottom-top+1 and b2-b1 becomes right-left+1:

int height = bottom - top + 1;
vector<vector<T>> map(height);
for (int r = 0; r < height; r++) {
    int width = right - left + 1;
    map.emplace_back(width);
}
I can encapsulate all of this into a Map class:

template<class T> class RectangularPointyTopMap {
    vector<vector<T>> map;

    int left_, top_;
  public:
    RectangularPointyTopMap(int left, int top, int right, int bottom)
                 : left_(left), top_(top)
    {
        int height = bottom - top + 1;
        map.resize(height);
        for (int r = 0; r < height; r++) {
            int width = right - left + 1;
            map.emplace_back(width);
        }
    }

    inline T& at(int q, int r) {
        return map[r - top_][q - left_ + (r >> 1)];
    }
};

For the other map shapes, it’s only slightly more complicated, but the
same pattern applies: I have to study the loop that created the map in
order to figure out the size and array access for the map.

1D arrays are trickier and I won’t try to tackle them here. In practice,
I rarely use array storage for hex maps, except when the maps are large,
and my code is written in C++. Although it’s more compact, it almost
never makes a difference in practice in my projects. For most of my
projects, I use a hash table and/or graph representation. It gives me
the most flexibility and reusability. I only need the more compact
storage when storage size matters.
*/

/*
5.0 Rotation

There are two one-step rotation functions, but which is “left” and which
is “right” depends on your map orientation. You may have to swap these.

Hex hex_rotate_left(Hex a)
{
    return Hex(-a.s, -a.q, -a.r);
}

Hex hex_rotate_right(Hex a)
{
    return Hex(-a.r, -a.s, -a.q);
}

Note that these are slightly different from the main page because q,r,s
don’t quite line up with x,y,z.
*/

func (a Hex) HexRotateLeft() Hex {
	return Hex{q: -a.s, r: -a.q, s: -a.r}
}

func (a Hex) HexRotateRight() Hex {
	return Hex{q: -a.r, r: -a.s, s: -a.q}
}

/*
If you think of the coordinates v in vector format, these operations are
3x3 matrix multiplies, M times v, where M = [0 0 -1; -1 0 0; 0 -1 0].
The matrix inverse M-1 = [0 -1 0; 0 0 -1; -1 0 0] rotates in the opposite
direction. Raising the matrix to a power Mk rotates k times. You can
precomputate all the rotation matrices, or combine the matrix with other
operations such as translate, scale, etc.
*/

/*
6.0 Offset coordinates

I use the names q and r for cube/axial coordinates, and col and row for
offset coordinates:
*/

type OffsetCoord struct {
	Col, Row int
}

func NewOffsetCoord(col, row int) OffsetCoord {
	return OffsetCoord{Col: col, Row: row}
}

/*
I’m expecting that I’ll use the cube/axial Hex class everywhere, except
for displaying to the player. That’s where offset coordinates will be
useful. That means the only operations I need are converting Hex to
OffsetCoord and back.

There are four offset types: odd-r, even-r, odd-q, even-q. The “r” types
are used with pointy top hexagons and the “q” types are used with flat
top. Whether it’s even or odd can be encoded as an offset direction
+1 or -1. For pointy top, the offset direction tells us whether to slide
alternate rows right or left. For flat top, the offset direction tells
us whether to slide alternate columns up or down.

const (
	EVEN int = +1
	ODD  int = -1
)

func qoffset_from_cube(offset int, h Hex) OffsetCoord {
	if !(offset == EVEN || offset == ODD) {
		panic("assert(offset == EVEN || offset == ODD)")
	}
	col := h.q
	row := h.r + int((h.q+int(offset)*(h.q&1))/2)
	return OffsetCoord{col: col, row: row}
}

func qoffset_to_cube(offset int, h OffsetCoord) Hex {
	if !(offset == EVEN || offset == ODD) {
		panic("assert(offset == EVEN || offset == ODD)")
	}
	q := h.col
	r := h.row - int((h.col+int(offset)*(h.col&1))/2)
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

func roffset_from_cube(offset int, h Hex) OffsetCoord {
	if !(offset == EVEN || offset == ODD) {
		panic("assert(offset == EVEN || offset == ODD)")
	}
	col := h.q + int((h.r+int(offset)*(h.r&1))/2)
	row := h.r
	return OffsetCoord{col: col, row: row}
}

func roffset_to_cube(offset int, h OffsetCoord) Hex {
	if !(offset == EVEN || offset == ODD) {
		panic("assert(offset == EVEN || offset == ODD)")
	}
	q := h.col - int((h.row+int(offset)*(h.row&1))/2)
	r := h.row
	s := -q - r
	return Hex{q: q, r: r, s: s}
}

If you’re only using even or odd, you can hard-code the value of
offset into the code, making it simpler and faster. Alternatively,
offset can be a template parameter so that the compiler can inline
and optimize it.
*/

func (h Hex) CubeToQOffset(even bool) OffsetCoord {
	col := h.q
	var row int
	if even {
		row = h.r + (h.q+1*(h.q&1))/2
	} else {
		row = h.r + (h.q-1*(h.q&1))/2
	}
	return OffsetCoord{Col: col, Row: row}
}

func (h OffsetCoord) QOffsetToCube(even bool) Hex {
	q := h.Col
	var r int
	if even {
		r = h.Row - (h.Col+1*(h.Col&1))/2
	} else {
		r = h.Row - (h.Col-1*(h.Col&1))/2
	}
	return Hex{
		q: q,
		r: r,
		s: -q - r,
	}
}

func (h Hex) CubeToROffset(even bool) OffsetCoord {
	var col int
	if even {
		col = h.q + (h.r+1*(h.r&1))/2
	} else {
		col = h.q + (h.r+1*(h.r&1))/2
	}
	return OffsetCoord{
		Col: col,
		Row: h.r,
	}
}

func (h OffsetCoord) ROffsetToCube(even bool) Hex {
	var q int
	if even {
		q = h.Col - (h.Row+1*(h.Row&1))/2
	} else {
		q = h.Col - (h.Row-1*(h.Row&1))/2
	}
	r := h.Row
	return Hex{
		q: q,
		r: r,
		s: -q - r,
	}
}
