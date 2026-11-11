# hexg

A Go package for hexagonal grids using cube coordinates.

This package implements the hexagonal grid system described in the excellent
[Red Blob Games Hexagonal Grids Guide](https://www.redblobgames.com/grids/hexagons/).

## Installation

```bash
go get github.com/maloquacious/hexg
```

## Usage

### Creating Hexes

```go
import "github.com/maloquacious/hexg"

// Create a hex using axial coordinates (q, r)
// The s coordinate is derived from the constraint q + r + s = 0
hex := hexg.NewHex(2, -1)

// Access coordinates
q, r, s := hex.QRS()
```

### Hex Operations

```go
a := hexg.NewHex(1, 0)
b := hexg.NewHex(2, -1)

// Arithmetic
sum := a.Add(b)
diff := a.Subtract(b)
scaled := a.Multiply(3)

// Distance
dist := a.Distance(b)

// Neighbors (direction 0-5)
neighbor := a.Neighbor(0)

// Diagonal neighbors (direction 0-5)
diagonal := a.DiagonalNeighbor(0)

// Reflections across axes
reflected := a.ReflectQ() // also ReflectR(), ReflectS()

// Rotations (60° increments)
left := a.RotateLeft()
right := a.RotateRight()
```

### Drawing Lines

```go
start := hexg.NewHex(0, 0)
end := hexg.NewHex(5, -3)

// Get all hexes along the line
line := start.LineDraw(end)
```

### Layouts and Screen Coordinates

```go
// Create a layout for pointy-top hexes
// size is the hex size, origin is the screen origin
layout := hexg.NewLayout(hexg.OddR, hexg.Point{X: 32, Y: 32}, hexg.Point{X: 0, Y: 0})

// Convert hex to screen pixel
hex := hexg.NewHex(2, 1)
pixel := layout.HexToPixel(hex)

// Convert screen pixel to hex
clickedHex := layout.PixelToHexRounded(hexg.Point{X: 100, Y: 75})

// Get polygon corners for rendering
corners := layout.PolygonCorners(hex)
```

### Layout Types

- `OddR` - Pointy-top, horizontal rows, odd rows shoved right
- `EvenR` - Pointy-top, horizontal rows, even rows shoved right
- `OddQ` - Flat-top, vertical columns, odd columns shoved down
- `EvenQ` - Flat-top, vertical columns, even columns shoved down

### Generating Hex Grids

```go
layout := hexg.NewLayout(hexg.OddR, hexg.Point{X: 32, Y: 32}, hexg.Point{X: 0, Y: 0})

// Hexagonal grid with radius 3
hexGrid := layout.Hexagon(3)

// Rectangular grid
rectGrid := layout.Rectangle(0, 10, 0, 8)

// Triangle grids
triUp := layout.TriangleUpDown(5)
triLR := layout.TriangleLeftRight(5)

// Ring of hexes at radius 3 from center
center := hexg.NewHex(0, 0)
ring := center.Ring(3)

// Spiral outward from center (all hexes within radius)
spiral := center.Spiral(3)
```

The generators return a `HexSet`, which is a map, so ranging one visits hexes in
a different order every run. Use `Sorted` when the order matters:

```go
// Deterministic order: by r, then q (row-major for pointy-top layouts)
for _, h := range layout.Hexagon(3).Sorted() {
	fmt.Println(h.ConciseString())
}

// The same ordering, for sorting a slice yourself
slices.SortFunc(hexes, hexg.Hex.Compare)
```

### Offset Coordinates

```go
hex := hexg.NewHex(2, -1)

// Convert to offset coordinates
offset := hex.CubeToROffset(false) // odd-r
col, row := offset.Col, offset.Row

// Convert back to cube coordinates
hexBack := offset.ROffsetToCube(false)
```

### Cylindrical Worlds

`Cylinder` wraps an offset rectangle east-west and walls its north and south
edges. `NewCylinder(layout, columns, rows)` returns an error for invalid
dimensions: Q layouts require even columns of at least 4, R layouts require odd
columns of at least 3, and rows must be positive. The total cell count must fit
in an `int`. Canonical offset columns and rows start at zero.

```go
layout := hexg.NewLayout(hexg.OddQ, hexg.Point{X: 32, Y: 32}, hexg.Point{})
world, err := hexg.NewCylinder(layout, 1920, 1080)
if err != nil {
	return err
}

west := layout.OffsetToCube(hexg.OffsetCoord{Col: 0, Row: 540})
east := layout.OffsetToCube(hexg.OffsetCoord{Col: 1919, Row: 540})
distance := world.Distance(west, east) // 1, across the seam
line := world.LineDraw(west, east)    // both endpoints, along a shortest route

// Dense storage: index = row*columns + col.
cells := make([]int, world.Columns()*world.Rows())
if index, ok := world.Index(east); ok {
	cells[index] = 1
	h, _ := world.HexAt(index) // inverse lookup; h == east
	_ = h
}
```

- `Normalize` wraps the offset column without changing the offset row.
- `Contains` checks the row bounds; noncanonical columns are accepted.
- `Neighbor` returns `(Hex, bool)`, with false for a step through a wall.
  `Neighbors` omits those steps and returns the others in direction order 0..5.
- `Ring` returns hexes at exactly the wrapped radius; `Spiral` returns hexes
  within it. Both clip at the walls, remove duplicates, and sort by `Hex.Compare`.
  Negative radii panic, as on `Hex`.
- `LineDraw` stays on the map. Equal-length routes prefer the canonical endpoint
  copy, then west, then east; step ties prefer direction order 0..5.
- Walks from an off-map hex return empty results. `Index` and `HexAt` return
  false for out-of-bounds rows or indices. All operations accept wrapped copies.

### JSON Serialization

Hex coordinates serialize to a compact format:

```go
hex := hexg.NewHex(2, -1)
data, _ := json.Marshal(hex) // "+2-1-1"
```

## Acknowledgements

This package is a Go implementation based on the comprehensive hexagonal grid
guide by [Red Blob Games](https://www.redblobgames.com/grids/hexagons/).
Amit Patel's detailed explanations and interactive examples made this
implementation possible.

Development assisted by [Amp](https://ampcode.com), an AI coding agent.

## License

MIT License - see [LICENSE](LICENSE) for details.
