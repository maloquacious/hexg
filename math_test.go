// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import (
	"math"
	"testing"
)

func TestAbs(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   int
		want int
	}{
		{"positive", 5, 5},
		{"negative", -5, 5},
		{"zero", 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := abs(tc.in); got != tc.want {
				t.Errorf("abs(%d) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

func TestAbsFloat(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   float64
		want float64
	}{
		{"positive", 5.5, 5.5},
		{"negative", -5.5, 5.5},
		{"zero", 0.0, 0.0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := abs(tc.in); got != tc.want {
				t.Errorf("abs(%g) = %g, want %g", tc.in, got, tc.want)
			}
		})
	}
}

func TestLerp(t *testing.T) {
	for _, tc := range []struct {
		name string
		a, b float64
		t    float64
		want float64
	}{
		{"t=0", 0.0, 10.0, 0.0, 0.0},
		{"t=1", 0.0, 10.0, 1.0, 10.0},
		{"t=0.5", 0.0, 10.0, 0.5, 5.0},
		{"t=0.25", 0.0, 10.0, 0.25, 2.5},
		{"negative range", -10.0, 10.0, 0.5, 0.0},
		{"backwards", 10.0, 0.0, 0.5, 5.0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := lerp(tc.a, tc.b, tc.t)
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("lerp(%g, %g, %g) = %g, want %g", tc.a, tc.b, tc.t, got, tc.want)
			}
		})
	}
}

func TestLerpInt(t *testing.T) {
	got := lerp(0, 10, 0.5)
	if math.Abs(got-5.0) > 1e-9 {
		t.Errorf("lerp(0, 10, 0.5) = %g, want 5.0", got)
	}
}
