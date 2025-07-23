// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements some tests for the hexg package
package main

import (
	"fmt"
	"github.com/maloquacious/hexg"
)

func main() {
	// print out neighbors
	h := hexg.NewHex(0, 0, 0)
	for _, direction := range []int{0, 1, 2, 3, 4, 5} {
		n := h.Neighbor(direction)
		fmt.Printf("%d: %q\n", direction, n)
	}

	// print out some pixels
	l := hexg.NewLayoutFlat(hexg.NewPoint(1, 1), hexg.NewPoint(0, 0), false)
	fmt.Printf("%s: pixel %q\n", h, h.ToPixel(l))
	for _, corner := range []int{0, 1, 2, 3, 4, 5} {
		fmt.Printf("%s: corner %d: pixel %q\n", h, corner, l.HexCornerOffset(corner))
	}
	for corner, point := range h.PolygonCorners(l) {
		fmt.Printf("%s: corner %d: pixel %q\n", h, corner, point)
	}
	for corner, point := range l.PolygonCorners() {
		fmt.Printf("layout: corner %d: pixel %q\n", corner, point)
	}
}
