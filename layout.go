package hexg

import "math"

// Layout constants for hexagonal grid orientations
const (
	EVEN_Q = iota // vertical layout shoves even columns down
	ODD_Q         // vertical layout shoves odd columns down
	EVEN_R        // horizontal layout shoves even rows right
	ODD_R         // horizontal layout shoves odd rows right
)

// Layout represents the orientation of a hexagonal grid
type Layout struct {
	Orientation int
	Size        float64
	Width       float64
	Height      float64
}

// NewLayout creates a new Layout with default EVEN_Q orientation and size 1.0
func NewLayout() Layout {
	return NewLayoutWithSize(EVEN_Q, 1.0)
}

// NewLayoutWithOrientation creates a new Layout with specified orientation and size 1.0
func NewLayoutWithOrientation(orientation int) Layout {
	return NewLayoutWithSize(orientation, 1.0)
}

// NewLayoutWithSize creates a new Layout with specified orientation and size
func NewLayoutWithSize(orientation int, size float64) Layout {
	var width, height float64
	
	// Flat top orientations (EVEN_Q, EVEN_R)
	if orientation == EVEN_Q || orientation == EVEN_R {
		height = size * math.Sqrt(3.0)
		width = size * 2
	} else {
		// Pointy top orientations (ODD_Q, ODD_R)
		height = size * 2
		width = size * math.Sqrt(3.0)
	}
	
	return Layout{
		Orientation: orientation,
		Size:        size,
		Width:       width,
		Height:      height,
	}
}
