// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cube

import "math"

// helpers for math
//

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

// lerp is a generic linear interpolation function.
// Accepts any integer or floating-point for a and b, always returns float64.
func lerp[T number](a, b T, t float64) float64 {
	return float64(a) + (float64(b)-float64(a))*t
}

func round(f float64) int {
	return int(math.Round(f))
}
