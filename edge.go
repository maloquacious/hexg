// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import (
	"cmp"
	"slices"
)

// Edge identifies an undirected edge between adjacent hexes. It is comparable
// and independent of layout orientation. Its canonical owner sees the edge in
// direction 0, 1, or 2. The zero value is the origin's direction-0 edge.
// Use Cylinder methods for wrapped identity and incidence.
type Edge struct {
	owner     Hex
	direction uint8
}

// Edge returns h's edge in direction, coerced to 0..5 as in Neighbor.
func (h Hex) Edge(direction int) Edge {
	direction = (direction%6 + 6) % 6
	if direction >= 3 {
		h = h.Neighbor(direction)
		direction -= 3
	}
	return Edge{owner: h, direction: uint8(direction)}
}

// Edges returns h's edges in direction order 0..5.
func (h Hex) Edges() [6]Edge {
	var edges [6]Edge
	for direction := range edges {
		edges[direction] = h.Edge(direction)
	}
	return edges
}

// EdgeBetween returns the shared edge of planar neighbors a and b.
// Identical and nonadjacent hexes return the zero Edge and false.
func EdgeBetween(a, b Hex) (Edge, bool) {
	if direction := slices.Index(directionVectors[:], b.Subtract(a)); direction >= 0 {
		return a.Edge(direction), true
	}
	return Edge{}, false
}

// Owner returns the canonical owning hex.
func (e Edge) Owner() Hex { return e.owner }

// Direction returns the edge's direction from Owner, always in 0..2.
func (e Edge) Direction() int { return int(e.direction) }

// Compare orders edges by Owner.Compare, then Direction. This order is stable.
func (e Edge) Compare(other Edge) int {
	return cmp.Or(e.owner.Compare(other.owner), cmp.Compare(e.direction, other.direction))
}

// Hexes returns the owner followed by the hex across the edge, in planar space.
func (e Edge) Hexes() [2]Hex {
	return [2]Hex{e.owner, e.owner.Neighbor(e.Direction())}
}

// Vertices returns the edge's endpoints in Vertex.Compare order, in planar space.
func (e Edge) Vertices() [2]Vertex {
	vertices := [2]Vertex{e.owner.Vertex(e.Direction() - 1), e.owner.Vertex(e.Direction())}
	slices.SortFunc(vertices[:], Vertex.Compare)
	return vertices
}

// Continuations returns the two other edges at each endpoint. Endpoints are in
// Vertices order; each pair is in Edge.Compare order. Incidence is planar.
func (e Edge) Continuations() [2][2]Edge {
	var result [2][2]Edge
	for i, vertex := range e.Vertices() {
		j := 0
		for _, edge := range vertex.Edges() {
			if edge != e {
				result[i][j] = edge
				j++
			}
		}
	}
	return result
}
