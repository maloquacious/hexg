// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"testing"

	"github.com/maloquacious/hexg"
)

func TestPoint_String(t *testing.T) {
	for _, tc := range []struct {
		x, y float64
		want string
	}{
		{0, 0, "0,0"},
		{10.5, 20.5, "10.5,20.5"},
		{-5, 3.14159, "-5,3.14159"},
		{100, -200, "100,-200"},
	} {
		p := hexg.Point{X: tc.x, Y: tc.y}
		if got := p.String(); got != tc.want {
			t.Errorf("Point{%g, %g}.String() = %q, want %q", tc.x, tc.y, got, tc.want)
		}
	}
}

func TestPoint_Fields(t *testing.T) {
	p := hexg.Point{X: 42.5, Y: 73.2}
	if p.X != 42.5 {
		t.Errorf("Point.X = %g, want 42.5", p.X)
	}
	if p.Y != 73.2 {
		t.Errorf("Point.Y = %g, want 73.2", p.Y)
	}
}
