// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

// number is a constraint that permits any integer or floating-point type.
type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// abs returns the absolute value of x.
func abs[T number](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
