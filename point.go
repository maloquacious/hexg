package hexg

// Point represents a coordinate on the hexagonal grid
type Point struct {
	X float64
	Y float64
}

// NewPoint creates a new Point with specified coordinates
func NewPoint(x, y float64) Point {
	return Point{X: x, Y: y}
}
