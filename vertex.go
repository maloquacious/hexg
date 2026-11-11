// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import (
	"cmp"
	"slices"
)

// Vertex identifies a corner shared by three hexes. It is comparable and
// independent of layout orientation. The planar owner is the smallest incident
// hex under Hex.Compare; its local corner is 4 or 5 (wire slot 0 or 1).
// The zero value is the origin's corner-4 vertex.
// Use Cylinder methods for wrapped identity and incidence.
type Vertex struct {
	owner  Hex
	corner uint8
}

// Vertex returns the corner shared by h, h.Neighbor(corner), and
// h.Neighbor(corner+1). Corner is coerced to 0..5. This numbering is based on
// neighbor directions, not Layout.PolygonCorners indices.
func (h Hex) Vertex(corner int) Vertex {
	corner = (corner%6 + 6) % 6
	hexes := [3]Hex{h, h.Neighbor(corner), h.Neighbor(corner + 1)}
	owner := slices.MinFunc(hexes[:], Hex.Compare)
	// A minimum-r, then minimum-q owner has neighbors (4,5) or (5,0).
	var slot uint8
	if slices.Contains(hexes[:], owner.Neighbor(0)) {
		slot = 1
	}
	return Vertex{owner: owner, corner: slot}
}

// Vertices returns h's vertices in local corner order 0..5.
func (h Hex) Vertices() [6]Vertex {
	var vertices [6]Vertex
	for corner := range vertices {
		vertices[corner] = h.Vertex(corner)
	}
	return vertices
}

// Owner returns the canonical owning hex. On a cylinder its offset column is
// canonical, but it need not sort before the other normalized incident hexes.
func (v Vertex) Owner() Hex { return v.owner }

// Corner returns the local corner from Owner, always 4 or 5.
func (v Vertex) Corner() int { return int(v.corner) + 4 }

// Compare orders vertices by Owner.Compare, then Corner. This order is stable.
func (v Vertex) Compare(other Vertex) int {
	return cmp.Or(v.owner.Compare(other.owner), cmp.Compare(v.corner, other.corner))
}

// Hexes returns the three planar incident hexes in Hex.Compare order.
func (v Vertex) Hexes() [3]Hex {
	hexes := [3]Hex{v.owner, v.owner.Neighbor(v.Corner()), v.owner.Neighbor(v.Corner() + 1)}
	slices.SortFunc(hexes[:], Hex.Compare)
	return hexes
}

// Edges returns the three planar incident edges in Edge.Compare order.
func (v Vertex) Edges() [3]Edge {
	h := v.Hexes()
	a, _ := EdgeBetween(h[0], h[1])
	b, _ := EdgeBetween(h[1], h[2])
	c, _ := EdgeBetween(h[2], h[0])
	edges := [3]Edge{a, b, c}
	slices.SortFunc(edges[:], Edge.Compare)
	return edges
}
