// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import "fmt"

// Point represents a screen coordinate.
type Point struct {
	X float64
	Y float64
}

// String returns the coordinates formatted as "x,y".
func (p Point) String() string {
	return fmt.Sprintf("%g,%g", p.X, p.Y)
}
