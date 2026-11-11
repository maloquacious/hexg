// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/maloquacious/hexg"
)

func TestNewCylinder(t *testing.T) {
	for _, tc := range []struct {
		id            int
		offset        hexg.LayoutOffset
		columns, rows int
		valid         bool
	}{
		{id: 1, offset: hexg.OddQ, columns: 1920, rows: 1080, valid: true},
		{id: 2, offset: hexg.EvenQ, columns: 4, rows: 1, valid: true},
		{id: 3, offset: hexg.OddR, columns: 511, rows: 255, valid: true},
		{id: 4, offset: hexg.EvenR, columns: 3, rows: 1, valid: true},
		{id: 5, offset: hexg.OddQ, columns: 5, rows: 3},
		{id: 6, offset: hexg.EvenQ, columns: 2, rows: 3},
		{id: 7, offset: hexg.OddR, columns: 4, rows: 3},
		{id: 8, offset: hexg.EvenR, columns: 1, rows: 3},
		{id: 9, offset: hexg.OddQ, columns: 0, rows: 3},
		{id: 10, offset: hexg.OddR, columns: -3, rows: 3},
		{id: 11, offset: hexg.OddQ, columns: 4, rows: 0},
		{id: 12, offset: hexg.EvenR, columns: 3, rows: -1},
		{id: 13, offset: hexg.OddQ, columns: 4, rows: int(^uint(0) >> 1)},
	} {
		layout := hexg.NewLayout(tc.offset, hexg.Point{}, hexg.Point{})
		c, err := hexg.NewCylinder(layout, tc.columns, tc.rows)
		if got := err == nil; got != tc.valid {
			t.Fatalf("%d: valid: want %v, got error %v", tc.id, tc.valid, err)
		}
		if tc.valid && (c.Columns() != tc.columns || c.Rows() != tc.rows) {
			t.Errorf("%d: dimensions: got %dx%d", tc.id, c.Columns(), c.Rows())
		}
	}
}

func TestCylinder_OffsetAndIndex(t *testing.T) {
	for _, offset := range []hexg.LayoutOffset{hexg.OddQ, hexg.EvenQ, hexg.OddR, hexg.EvenR} {
		layout := hexg.NewLayout(offset, hexg.Point{}, hexg.Point{})
		columns := 5
		if layout.IsFlat() {
			columns = 6
		}
		c, err := hexg.NewCylinder(layout, columns, 7)
		if err != nil {
			t.Fatal(err)
		}
		for row := -1; row <= c.Rows(); row++ {
			for col := -3 * columns; col <= 3*columns; col++ {
				h := layout.OffsetToCube(hexg.OffsetCoord{Col: col, Row: row})
				wantCol := (col%columns + columns) % columns
				want := layout.OffsetToCube(hexg.OffsetCoord{Col: wantCol, Row: row})
				if got := c.Normalize(h); got != want || c.Normalize(got) != got {
					t.Fatalf("%v: normalize %v: want %v, got %v", offset, h, want, got)
				}
				inside := 0 <= row && row < c.Rows()
				if got := c.Contains(h); got != inside {
					t.Fatalf("%v: contains %v: want %v, got %v", offset, h, inside, got)
				}
				index, ok := c.Index(h)
				if ok != inside || (inside && index != row*columns+wantCol) {
					t.Fatalf("%v: index %v: got %d, %v", offset, h, index, ok)
				}
				if inside {
					if got, ok := c.HexAt(index); !ok || got != want {
						t.Fatalf("%v: inverse %d: want %v, got %v, %v", offset, index, want, got, ok)
					}
				} else if len(c.Neighbors(h)) != 0 || len(c.Ring(h, 0)) != 0 || len(c.Spiral(h, 2)) != 0 || len(c.LineDraw(h, hexg.Hex{})) != 0 {
					t.Fatalf("%v: walks from outside the walls must be empty", offset)
				}
			}
		}
		for _, index := range []int{-1, columns * c.Rows()} {
			if _, ok := c.HexAt(index); ok {
				t.Errorf("%v: accepted invalid index %d", offset, index)
			}
		}
	}
}

// Build the graph in offset space without calling any Cylinder methods. BFS
// is an independent oracle: it does not assume three copies suffice, that walls
// leave distances unchanged, or that normalizing a planar ring gives the ring.
func cylinderGraph(layout hexg.Layout, columns, rows int) [][]int {
	graph := make([][]int, columns*rows)
	for index := range graph {
		h := layout.OffsetToCube(hexg.OffsetCoord{Col: index % columns, Row: index / columns})
		for direction := range 6 {
			n := layout.CubeToOffset(h.Neighbor(direction))
			if 0 <= n.Row && n.Row < rows {
				col := (n.Col%columns + columns) % columns
				graph[index] = append(graph[index], n.Row*columns+col)
			}
		}
	}
	return graph
}

func TestCylinder_WalksAgainstGraph(t *testing.T) {
	for _, offset := range []hexg.LayoutOffset{hexg.OddQ, hexg.EvenQ, hexg.OddR, hexg.EvenR} {
		layout := hexg.NewLayout(offset, hexg.Point{}, hexg.Point{})
		for _, columns := range []int{3, 5, 7} {
			if layout.IsFlat() {
				columns++
			}
			for _, rows := range []int{1, 2, 7} {
				c, err := hexg.NewCylinder(layout, columns, rows)
				if err != nil {
					t.Fatal(err)
				}
				graph := cylinderGraph(layout, columns, rows)
				for start := range graph {
					a, _ := c.HexAt(start)
					distances := make([]int, len(graph))
					for i := range distances {
						distances[i] = -1
					}
					distances[start] = 0
					queue := []int{start}
					for head := 0; head < len(queue); head++ {
						for _, next := range graph[queue[head]] {
							if distances[next] < 0 {
								distances[next] = distances[queue[head]] + 1
								queue = append(queue, next)
							}
						}
					}
					wantNeighbors := make([]hexg.Hex, 0, 6)
					for _, index := range graph[start] {
						h, _ := c.HexAt(index)
						wantNeighbors = append(wantNeighbors, h)
					}
					if got := c.Neighbors(a); !slices.Equal(wantNeighbors, got) {
						t.Fatalf("%v %dx%d %v: neighbors: want %v, got %v", offset, columns, rows, a, wantNeighbors, got)
					}
					for direction := -6; direction < 12; direction++ {
						n, ok := c.Neighbor(a, direction)
						if n2, ok2 := c.Neighbor(a, (direction+6)%6); n != n2 || ok != ok2 {
							t.Fatal("neighbor direction coercion disagrees")
						}
						if ok {
							if back, valid := c.Neighbor(n, direction+3); !valid || back != a {
								t.Fatalf("%v: seam neighbor does not round trip", offset)
							}
						}
					}
					for end, want := range distances {
						b, _ := c.HexAt(end)
						if got := c.Distance(a, b); got != want {
							t.Fatalf("%v %dx%d %v -> %v: distance: want %d, got %d", offset, columns, rows, a, b, want, got)
						}
						line := c.LineDraw(a, b)
						if len(line) != want+1 || line[0] != a || line[len(line)-1] != b {
							t.Fatalf("%v %dx%d %v -> %v: line: want %d steps, got %v", offset, columns, rows, a, b, want, line)
						}
						previous := start
						for i, h := range line {
							index, ok := c.Index(h)
							if !ok || h != c.Normalize(h) || (i > 0 && !slices.Contains(graph[previous], index)) {
								t.Fatalf("%v: invalid line step %v", offset, line)
							}
							previous = index
						}
						// Changing the lift must not change any tie-breaking.
						oc := layout.CubeToOffset(b)
						oc.Col -= 3 * columns
						copy := layout.OffsetToCube(oc)
						if c.Distance(a, copy) != want || !slices.Equal(line, c.LineDraw(a, copy)) {
							t.Fatalf("%v: noncanonical endpoint changes route", offset)
						}
					}
					for radius := 0; radius <= 2*columns+rows; radius++ {
						ring, spiral := make(hexg.HexSet), make(hexg.HexSet)
						for index, d := range distances {
							h, _ := c.HexAt(index)
							if d == radius {
								ring[h] = struct{}{}
							}
							if d <= radius {
								spiral[h] = struct{}{}
							}
						}
						if got := c.Ring(a, radius); !slices.Equal(ring.Sorted(), got) {
							t.Fatalf("%v %dx%d %v radius %d: ring: want %v, got %v", offset, columns, rows, a, radius, ring.Sorted(), got)
						}
						if got := c.Spiral(a, radius); !slices.Equal(spiral.Sorted(), got) {
							t.Fatalf("%v %dx%d %v radius %d: spiral: want %v, got %v", offset, columns, rows, a, radius, spiral.Sorted(), got)
						}
					}
				}
			}
		}
	}
}

func TestCylinder_NegativeRadius(t *testing.T) {
	c, err := hexg.NewCylinder(hexg.Layout{}, 3, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, walk := range []func(hexg.Hex, int) []hexg.Hex{c.Ring, c.Spiral} {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("negative radius did not panic")
				}
			}()
			walk(hexg.Hex{}, -1)
		}()
	}
}

func TestCylinder_LineDrawTieAndWall(t *testing.T) {
	layout := hexg.NewLayout(hexg.OddQ, hexg.Point{}, hexg.Point{})
	c, err := hexg.NewCylinder(layout, 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	a := layout.OffsetToCube(hexg.OffsetCoord{Col: 0, Row: 0})
	b := layout.OffsetToCube(hexg.OffsetCoord{Col: 2, Row: 0})
	// Both directions are equally short. Prefer the canonical endpoint's copy,
	// and stay on row zero: the ordinary planar line rounds outside the wall.
	want := []hexg.Hex{a, layout.OffsetToCube(hexg.OffsetCoord{Col: 1, Row: 0}), b}
	if got := c.LineDraw(a, b); !slices.Equal(want, got) {
		t.Fatalf("line: want %v, got %v", want, got)
	}
}

func ExampleCylinder() {
	layout := hexg.NewLayout(hexg.OddQ, hexg.Point{X: 32, Y: 32}, hexg.Point{})
	world, err := hexg.NewCylinder(layout, 1920, 1080)
	if err != nil {
		panic(err)
	}
	west := layout.OffsetToCube(hexg.OffsetCoord{Col: 0, Row: 540})
	east := layout.OffsetToCube(hexg.OffsetCoord{Col: 1919, Row: 540})
	fmt.Println(world.Distance(west, east))
	index, ok := world.Index(east)
	fmt.Println(index, ok)
	h, ok := world.HexAt(index)
	fmt.Println(h == east, ok)
	fmt.Println(len(world.LineDraw(west, east)))
	// Output:
	// 1
	// 1038719 true
	// true true
	// 2
}
