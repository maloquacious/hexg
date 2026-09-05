// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"math"
	"testing"

	"github.com/maloquacious/hexg"
)

func TestNewLayout(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}

	for _, tc := range []struct {
		name     string
		offset   hexg.LayoutOffset
		isFlat   bool
		isPointy bool
		isOdd    bool
		isEven   bool
	}{
		{"OddR", hexg.OddR, false, true, true, false},
		{"EvenR", hexg.EvenR, false, true, false, true},
		{"OddQ", hexg.OddQ, true, false, true, false},
		{"EvenQ", hexg.EvenQ, true, false, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			layout := hexg.NewLayout(tc.offset, size, origin)
			if layout.IsFlat() != tc.isFlat {
				t.Errorf("IsFlat() = %v, want %v", layout.IsFlat(), tc.isFlat)
			}
			if layout.IsPointy() != tc.isPointy {
				t.Errorf("IsPointy() = %v, want %v", layout.IsPointy(), tc.isPointy)
			}
			if layout.IsOdd() != tc.isOdd {
				t.Errorf("IsOdd() = %v, want %v", layout.IsOdd(), tc.isOdd)
			}
			if layout.IsEven() != tc.isEven {
				t.Errorf("IsEven() = %v, want %v", layout.IsEven(), tc.isEven)
			}
		})
	}
}

func TestLayout_IsOddR(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.OddR, size, origin)
	if !layout.IsOddR() {
		t.Error("OddR layout: IsOddR() should be true")
	}
	if layout.IsEvenR() || layout.IsOddQ() || layout.IsEvenQ() {
		t.Error("OddR layout: other Is*() should be false")
	}
}

func TestLayout_IsEvenR(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.EvenR, size, origin)
	if !layout.IsEvenR() {
		t.Error("EvenR layout: IsEvenR() should be true")
	}
	if layout.IsOddR() || layout.IsOddQ() || layout.IsEvenQ() {
		t.Error("EvenR layout: other Is*() should be false")
	}
}

func TestLayout_IsOddQ(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.OddQ, size, origin)
	if !layout.IsOddQ() {
		t.Error("OddQ layout: IsOddQ() should be true")
	}
	if layout.IsOddR() || layout.IsEvenR() || layout.IsEvenQ() {
		t.Error("OddQ layout: other Is*() should be false")
	}
}

func TestLayout_IsEvenQ(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.EvenQ, size, origin)
	if !layout.IsEvenQ() {
		t.Error("EvenQ layout: IsEvenQ() should be true")
	}
	if layout.IsOddR() || layout.IsEvenR() || layout.IsOddQ() {
		t.Error("EvenQ layout: other Is*() should be false")
	}
}

func TestLayout_HexToPixel(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 100, Y: 100}

	layout := hexg.NewLayout(hexg.OddR, size, origin)
	hex := hexg.NewHex(0, 0)
	p := layout.HexToPixel(hex)

	if p.X != 100 || p.Y != 100 {
		t.Errorf("HexToPixel(0,0) = (%g, %g), want (100, 100)", p.X, p.Y)
	}
}

func TestLayout_PixelToHexRounded(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.OddR, size, origin)

	hex := hexg.NewHex(2, 3)
	pixel := layout.HexToPixel(hex)
	roundTrip := layout.PixelToHexRounded(pixel)

	if !hex.Equals(roundTrip) {
		t.Errorf("round trip failed: started with %v, got %v", hex, roundTrip)
	}
}

func TestLayout_PolygonCorners(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.OddR, size, origin)

	hex := hexg.NewHex(0, 0)
	corners := layout.PolygonCorners(hex)

	if len(corners) != 6 {
		t.Errorf("PolygonCorners returned %d corners, want 6", len(corners))
	}

	center := layout.HexToPixel(hex)
	for i, corner := range corners {
		dx := corner.X - center.X
		dy := corner.Y - center.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if math.Abs(dist-10) > 0.0001 {
			t.Errorf("corner %d distance from center = %g, want 10", i, dist)
		}
	}
}

func TestLayout_BoundingBox(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.OddR, size, origin)

	hexes := []hexg.Hex{
		hexg.NewHex(0, 0),
		hexg.NewHex(2, 1),
		hexg.NewHex(-1, 3),
	}

	ul, lr := layout.BoundingBox(hexes...)
	ulOC := layout.CubeToOffset(ul)
	lrOC := layout.CubeToOffset(lr)

	if ulOC.Col > lrOC.Col || ulOC.Row > lrOC.Row {
		t.Errorf("BoundingBox upper-left (%v) should be <= lower-right (%v)", ulOC, lrOC)
	}
}

func TestLayout_BoundingBoxEmpty(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.OddR, size, origin)

	ul, lr := layout.BoundingBox()
	if !ul.Equals(hexg.NewHex(0, 0)) || !lr.Equals(hexg.NewHex(0, 0)) {
		t.Errorf("empty BoundingBox should return origin hexes")
	}
}

func TestLayout_CubeToOffsetRoundTrip(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}

	for _, offset := range []hexg.LayoutOffset{hexg.OddR, hexg.EvenR, hexg.OddQ, hexg.EvenQ} {
		layout := hexg.NewLayout(offset, size, origin)
		for q := -3; q <= 3; q++ {
			for r := -3; r <= 3; r++ {
				hex := hexg.NewHex(q, r)
				oc := layout.CubeToOffset(hex)
				back := layout.OffsetToCube(oc)
				if !hex.Equals(back) {
					t.Errorf("offset %d: round trip failed for %v -> %v -> %v", offset, hex, oc, back)
				}
			}
		}
	}
}

func TestLayout_Hexagon(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.OddR, size, origin)

	gs := layout.Hexagon(2)
	expectedCount := 1 + 6 + 12
	if len(gs) != expectedCount {
		t.Errorf("Hexagon(2) has %d hexes, want %d", len(gs), expectedCount)
	}

	center := hexg.NewHex(0, 0)
	if _, ok := gs[center]; !ok {
		t.Error("Hexagon(2) should contain center hex")
	}
}

func TestLayout_Rectangle(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.OddR, size, origin)

	gs := layout.Rectangle(0, 3, 0, 2)
	expectedCount := 4 * 3
	if len(gs) != expectedCount {
		t.Errorf("Rectangle(0,3,0,2) has %d hexes, want %d", len(gs), expectedCount)
	}
}

func TestLayout_Rectangle_OffsetRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		id                       int
		offset                   hexg.LayoutOffset
		name                     string
		left, right, top, bottom int
	}{
		{id: 1, offset: hexg.OddR, name: "OddR", left: 0, right: 4, top: 0, bottom: 3},
		{id: 2, offset: hexg.EvenR, name: "EvenR", left: 0, right: 4, top: 0, bottom: 3},
		{id: 3, offset: hexg.OddQ, name: "OddQ", left: 0, right: 4, top: 0, bottom: 3},
		{id: 4, offset: hexg.EvenQ, name: "EvenQ", left: 0, right: 4, top: 0, bottom: 3},
		{id: 5, offset: hexg.OddR, name: "OddR", left: -3, right: 2, top: -2, bottom: 1},
		{id: 6, offset: hexg.EvenR, name: "EvenR", left: -3, right: 2, top: -2, bottom: 1},
		{id: 7, offset: hexg.OddQ, name: "OddQ", left: -3, right: 2, top: -2, bottom: 1},
		{id: 8, offset: hexg.EvenQ, name: "EvenQ", left: -3, right: 2, top: -2, bottom: 1},
		{id: 9, offset: hexg.EvenR, name: "EvenR", left: 1, right: 1, top: 1, bottom: 1},
		{id: 10, offset: hexg.EvenQ, name: "EvenQ", left: 1, right: 1, top: 1, bottom: 1},
	} {
		layout := hexg.NewLayout(tc.offset, hexg.Point{X: 10, Y: 10}, hexg.Point{X: 0, Y: 0})
		gs := layout.Rectangle(tc.left, tc.right, tc.top, tc.bottom)

		cols := tc.right - tc.left + 1
		rows := tc.bottom - tc.top + 1
		if want, got := cols*rows, len(gs); want != got {
			t.Errorf("%d: %s: count: want %d, got %d\n", tc.id, tc.name, want, got)
		}

		// every hex must land in the requested offset window, and every
		// cell of that window must be occupied exactly once.
		seen := make(map[hexg.OffsetCoord]bool)
		for h := range gs {
			oc := layout.CubeToOffset(h)
			if oc.Col < tc.left || oc.Col > tc.right || oc.Row < tc.top || oc.Row > tc.bottom {
				t.Errorf("%d: %s: %s is outside col %d..%d, row %d..%d\n",
					tc.id, tc.name, oc, tc.left, tc.right, tc.top, tc.bottom)
				continue
			}
			if seen[oc] {
				t.Errorf("%d: %s: %s occupied twice\n", tc.id, tc.name, oc)
			}
			seen[oc] = true
		}
		for row := tc.top; row <= tc.bottom; row++ {
			for col := tc.left; col <= tc.right; col++ {
				if oc := hexg.NewOffsetCoord(col, row); !seen[oc] {
					t.Errorf("%d: %s: %s missing\n", tc.id, tc.name, oc)
				}
			}
		}
	}
}

func TestLayout_TriangleUpDown(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.OddR, size, origin)

	gs := layout.TriangleUpDown(3)
	expectedCount := (3 + 1) * (3 + 2) / 2
	if len(gs) != expectedCount {
		t.Errorf("TriangleUpDown(3) has %d hexes, want %d", len(gs), expectedCount)
	}
}

func TestLayout_TriangleLeftRight(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.OddR, size, origin)

	gs := layout.TriangleLeftRight(3)
	expectedCount := (3 + 1) * (3 + 2) / 2
	if len(gs) != expectedCount {
		t.Errorf("TriangleLeftRight(3) has %d hexes, want %d", len(gs), expectedCount)
	}
}

func TestLayout_Parallelogram(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.OddR, size, origin)

	gsQR := layout.ParallelogramQR(0, 0, 2, 2)
	if len(gsQR) != 9 {
		t.Errorf("ParallelogramQR(0,0,2,2) has %d hexes, want 9", len(gsQR))
	}

	gsQS := layout.ParallelogramQS(0, 0, 2, 2)
	if len(gsQS) != 9 {
		t.Errorf("ParallelogramQS(0,0,2,2) has %d hexes, want 9", len(gsQS))
	}

	gsRS := layout.ParallelogramRS(0, 0, 2, 2)
	if len(gsRS) != 9 {
		t.Errorf("ParallelogramRS(0,0,2,2) has %d hexes, want 9", len(gsRS))
	}
}

func TestLayout_RotateLeftRight(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	layout := hexg.NewLayout(hexg.OddR, size, origin)

	hex := hexg.NewHex(1, 0)
	rotated := layout.RotateLeft(hex)
	back := layout.RotateRight(rotated)

	if !hex.Equals(back) {
		t.Errorf("RotateLeft then RotateRight: got %v, want %v", back, hex)
	}
}

func TestLayout_IsHorizontalVertical(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}

	oddR := hexg.NewLayout(hexg.OddR, size, origin)
	if !oddR.IsHorizontal() || oddR.IsVertical() {
		t.Error("OddR should be horizontal, not vertical")
	}

	oddQ := hexg.NewLayout(hexg.OddQ, size, origin)
	if oddQ.IsHorizontal() || !oddQ.IsVertical() {
		t.Error("OddQ should be vertical, not horizontal")
	}
}

// TestLayout_ValueReceivers pins the value-receiver method set: a Layout must be
// usable in unaddressable positions - the result of NewLayout, a map value, a
// by-value parameter - without the caller binding it to a variable first.
// These calls do not compile against pointer receivers.
func TestLayout_ValueReceivers(t *testing.T) {
	size := hexg.Point{X: 10, Y: 10}
	origin := hexg.Point{X: 0, Y: 0}
	hex := hexg.NewHex(1, -2)

	layouts := map[hexg.LayoutOffset]hexg.Layout{
		hexg.OddR:  hexg.NewLayout(hexg.OddR, size, origin),
		hexg.EvenR: hexg.NewLayout(hexg.EvenR, size, origin),
		hexg.OddQ:  hexg.NewLayout(hexg.OddQ, size, origin),
		hexg.EvenQ: hexg.NewLayout(hexg.EvenQ, size, origin),
	}

	byValue := func(l hexg.Layout) hexg.Point { return l.HexToPixel(hex) }

	for _, tc := range []struct {
		id     int
		offset hexg.LayoutOffset
	}{
		{1, hexg.OddR},
		{2, hexg.EvenR},
		{3, hexg.OddQ},
		{4, hexg.EvenQ},
	} {
		layout := hexg.NewLayout(tc.offset, size, origin)
		want := layout.HexToPixel(hex)

		// method call directly on the constructor's return value
		if got := hexg.NewLayout(tc.offset, size, origin).HexToPixel(hex); want != got {
			t.Errorf("%d: constructor: want %v, got %v", tc.id, want, got)
		}
		// method call on an unaddressable map value
		if got := layouts[tc.offset].HexToPixel(hex); want != got {
			t.Errorf("%d: map value: want %v, got %v", tc.id, want, got)
		}
		// Layout passed by value
		if got := byValue(layout); want != got {
			t.Errorf("%d: by value: want %v, got %v", tc.id, want, got)
		}
		// pointers keep working: a *Layout can call a value method
		if got := (&layout).HexToPixel(hex); want != got {
			t.Errorf("%d: pointer: want %v, got %v", tc.id, want, got)
		}
	}
}
