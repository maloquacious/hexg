// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/maloquacious/hexg"
)

func TestEdgeAndVertex_IdentityAndIncidence(t *testing.T) {
	// Doubled edge midpoints and tripled vertex centers are independent,
	// integer-only identifiers used as an oracle for both sharing and aliasing.
	edgeAt := make(map[hexg.Hex]hexg.Edge)
	vertexAt := make(map[hexg.Hex]hexg.Vertex)
	edges := make(map[hexg.Edge]struct{})
	vertices := make(map[hexg.Vertex]struct{})
	for q := -3; q <= 3; q++ {
		for r := -3; r <= 3; r++ {
			h := hexg.NewHex(q, r)
			for d := range 6 {
				n := h.Neighbor(d)
				e := h.Edge(d)
				if e != n.Edge(d+3) || e != h.Edges()[d] || e != e.Owner().Edge(e.Direction()) {
					t.Fatalf("edge identity failed at %v direction %d", h, d)
				}
				if got, ok := hexg.EdgeBetween(h, n); !ok || got != e {
					t.Fatalf("edge between %v and %v: got %v, %v", h, n, got, ok)
				}
				hexes := e.Hexes()
				if !slices.Contains(hexes[:], h) || !slices.Contains(hexes[:], n) || hexes[0] != e.Owner() {
					t.Fatalf("edge %v: wrong incident hexes %v", e, hexes)
				}
				key := h.Add(n)
				if prior, ok := edgeAt[key]; ok && prior != e {
					t.Fatalf("same midpoint has different edges: %v and %v", prior, e)
				}
				edgeAt[key], edges[e] = e, struct{}{}
				endpoints := e.Vertices()
				if endpoints[0] == endpoints[1] || endpoints[0].Compare(endpoints[1]) >= 0 {
					t.Fatalf("edge %v: invalid endpoint order %v", e, endpoints)
				}
				for i, pair := range e.Continuations() {
					incident := endpoints[i].Edges()
					if pair[0] == e || pair[1] == e || pair[0].Compare(pair[1]) >= 0 || !slices.Contains(incident[:], pair[0]) || !slices.Contains(incident[:], pair[1]) {
						t.Fatalf("edge %v: wrong continuations %v", e, pair)
					}
				}
				v := h.Vertex(d)
				want := [3]hexg.Hex{h, n, h.Neighbor(d + 1)}
				key = want[0].Add(want[1]).Add(want[2])
				slices.SortFunc(want[:], hexg.Hex.Compare)
				if got := v.Hexes(); got != want || v != h.Vertices()[d] || v != v.Owner().Vertex(v.Corner()) || v.Owner() != want[0] {
					t.Fatalf("vertex %v: want hexes %v, got %v", v, want, got)
				}
				for _, incident := range want {
					corners := incident.Vertices()
					if !slices.Contains(corners[:], v) {
						t.Fatalf("vertex %v missing from incident hex %v", v, incident)
					}
				}
				if prior, ok := vertexAt[key]; ok && prior != v {
					t.Fatalf("same center has different vertices: %v and %v", prior, v)
				}
				vertexAt[key], vertices[v] = v, struct{}{}
				incident := v.Edges()
				if incident[0].Compare(incident[1]) >= 0 || incident[1].Compare(incident[2]) >= 0 {
					t.Fatalf("vertex %v: invalid edge order %v", v, incident)
				}
				for _, edge := range incident {
					ends := edge.Vertices()
					if !slices.Contains(ends[:], v) {
						t.Fatalf("vertex %v missing from edge %v", v, edge)
					}
				}
				if !slices.Contains(incident[:], h.Edge(d)) || !slices.Contains(incident[:], h.Edge(d+1)) {
					t.Fatalf("vertex %v does not join expected local edges", v)
				}
			}
		}
	}
	if len(edgeAt) != len(edges) || len(vertexAt) != len(vertices) {
		t.Fatal("distinct physical features alias to the same identifier")
	}
	for _, h := range []hexg.Hex{{}, hexg.NewHex(2, 0), hexg.NewHex(-1, 2)} {
		if _, ok := hexg.EdgeBetween(hexg.Hex{}, h); ok {
			t.Errorf("nonadjacent hex %v accepted", h)
		}
	}
	for _, d := range []int{-13, -1, 6, 19} {
		h := hexg.NewHex(-2, 5)
		normal := (d%6 + 6) % 6
		if h.Edge(d) != h.Edge(normal) || h.Vertex(d) != h.Vertex(normal) {
			t.Errorf("direction %d was not coerced", d)
		}
	}
}

func TestEdgeAndVertex_Compare(t *testing.T) {
	for _, tc := range []struct {
		id   int
		a, b hexg.Hex
		sign int
	}{
		{id: 1, a: hexg.NewHex(1, -2), b: hexg.NewHex(-5, 3), sign: -1},
		{id: 2, a: hexg.NewHex(-2, 1), b: hexg.NewHex(3, 1), sign: -1},
		{id: 3, a: hexg.NewHex(4, -2), b: hexg.NewHex(4, -2), sign: 0},
	} {
		for _, pair := range [][2]int{{0, 0}, {0, 1}, {1, 0}} {
			want := tc.sign
			if want == 0 {
				want = pair[0] - pair[1]
			}
			a, b := tc.a.Edge(pair[0]), tc.b.Edge(pair[1])
			if got := a.Compare(b); got != want || b.Compare(a) != -want || (got == 0) != (a == b) {
				t.Errorf("%d: edge compare: want %d, got %d", tc.id, want, got)
			}
			v, w := tc.a.Vertex(4+pair[0]), tc.b.Vertex(4+pair[1])
			if got := v.Compare(w); got != want || w.Compare(v) != -want || (got == 0) != (v == w) {
				t.Errorf("%d: vertex compare: want %d, got %d", tc.id, want, got)
			}
		}
	}
}

func ExampleEdge() {
	h := hexg.NewHex(0, 0)
	river := h.Edge(0)
	rivers := map[hexg.Edge]bool{river: true}
	across := h.Neighbor(0)
	fmt.Println(rivers[across.Edge(3)])
	fmt.Println(river.ConciseString())
	fmt.Println(h.Vertex(0) == across.Vertex(2))
	// Output:
	// true
	// E:+0+0+0:0
	// true
}
