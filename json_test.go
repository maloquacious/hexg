// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"encoding/json"
	"testing"

	"github.com/maloquacious/hexg"
)

func TestHex_MarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		hex  hexg.Hex
		want string
	}{
		{hexg.NewHex(0, 0), `"+0+0+0"`},
		{hexg.NewHex(1, 2), `"+1+2-3"`},
		{hexg.NewHex(-1, -2), `"-1-2+3"`},
		{hexg.NewHex(10, -5), `"+10-5-5"`},
	} {
		data, err := json.Marshal(tc.hex)
		if err != nil {
			t.Errorf("Marshal(%v) error: %v", tc.hex, err)
			continue
		}
		if string(data) != tc.want {
			t.Errorf("Marshal(%v) = %s, want %s", tc.hex, data, tc.want)
		}
	}
}

func TestHex_UnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  hexg.Hex
	}{
		{`"+0+0+0"`, hexg.NewHex(0, 0)},
		{`"+1+2-3"`, hexg.NewHex(1, 2)},
		{`"-1-2+3"`, hexg.NewHex(-1, -2)},
		{`"+10-5-5"`, hexg.NewHex(10, -5)},
	} {
		var got hexg.Hex
		err := json.Unmarshal([]byte(tc.input), &got)
		if err != nil {
			t.Errorf("Unmarshal(%s) error: %v", tc.input, err)
			continue
		}
		if !got.Equals(tc.want) {
			t.Errorf("Unmarshal(%s) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestHex_UnmarshalJSON_Invalid(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
	}{
		{"not a string", `123`},
		{"missing s coord", `"+1+2"`},
		{"invalid constraint", `"+1+2+3"`},
		{"empty string", `""`},
		{"no signs", `"123"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var h hexg.Hex
			err := json.Unmarshal([]byte(tc.input), &h)
			if err == nil {
				t.Errorf("Unmarshal(%s) should have failed", tc.input)
			}
		})
	}
}

func TestHex_JSONRoundTrip(t *testing.T) {
	hexes := []hexg.Hex{
		hexg.NewHex(0, 0),
		hexg.NewHex(1, 2),
		hexg.NewHex(-3, 5),
		hexg.NewHex(100, -50),
	}

	for _, h := range hexes {
		data, err := json.Marshal(h)
		if err != nil {
			t.Errorf("Marshal(%v) error: %v", h, err)
			continue
		}

		var got hexg.Hex
		err = json.Unmarshal(data, &got)
		if err != nil {
			t.Errorf("Unmarshal(%s) error: %v", data, err)
			continue
		}

		if !got.Equals(h) {
			t.Errorf("round trip failed: %v -> %s -> %v", h, data, got)
		}
	}
}

func TestHex_Value(t *testing.T) {
	h := hexg.NewHex(1, 2)
	val, err := h.Value()
	if err != nil {
		t.Errorf("Value() error: %v", err)
	}
	if val != "+1+2-3" {
		t.Errorf("Value() = %v, want +1+2-3", val)
	}
}

func TestHex_Scan(t *testing.T) {
	var h hexg.Hex
	err := h.Scan("+1+2-3")
	if err != nil {
		t.Errorf("Scan() error: %v", err)
	}
	if !h.Equals(hexg.NewHex(1, 2)) {
		t.Errorf("Scan(+1+2-3) = %v, want (1, 2, -3)", h)
	}
}

func TestHex_Scan_InvalidType(t *testing.T) {
	var h hexg.Hex
	err := h.Scan(123)
	if err == nil {
		t.Error("Scan(int) should have failed")
	}
}

func TestHex_Scan_InvalidString(t *testing.T) {
	var h hexg.Hex
	err := h.Scan("+1+2+3")
	if err == nil {
		t.Error("Scan(invalid constraint) should have failed")
	}
}

func TestEdgeAndVertex_Serialization(t *testing.T) {
	for _, tc := range []struct {
		id                   int
		h                    hexg.Hex
		direction, corner    int
		edgeWire, vertexWire string
	}{
		{id: 1, edgeWire: "E:+0+0+0:0", corner: 4, vertexWire: "V:+0+0+0:0"},
		{id: 2, h: hexg.NewHex(2, -3), direction: 2, corner: 5, edgeWire: "E:+2-3+1:2", vertexWire: "V:+2-3+1:1"},
		{id: 3, h: hexg.NewHex(-2, 5), direction: 3, corner: 0, edgeWire: "E:-3+5-2:0", vertexWire: "V:-1+4-3:0"},
		{id: 4, h: hexg.NewHex(-2, 5), direction: 4, corner: 3, edgeWire: "E:-3+6-3:1", vertexWire: "V:-3+5-2:1"},
	} {
		e, v := tc.h.Edge(tc.direction), tc.h.Vertex(tc.corner)
		if e.ConciseString() != tc.edgeWire || e.String() != tc.edgeWire || v.ConciseString() != tc.vertexWire || v.String() != tc.vertexWire {
			t.Fatalf("%d: wrong wire identifiers: %s, %s", tc.id, e, v)
		}
		if got, err := hexg.ParseEdge(tc.edgeWire); err != nil || got != e {
			t.Fatalf("%d: parse edge: got %v, %v", tc.id, got, err)
		}
		if got, err := hexg.ParseVertex(tc.vertexWire); err != nil || got != v {
			t.Fatalf("%d: parse vertex: got %v, %v", tc.id, got, err)
		}
		data, err := json.Marshal(e)
		if err != nil || string(data) != `"`+tc.edgeWire+`"` {
			t.Fatalf("%d: marshal edge: got %s, %v", tc.id, data, err)
		}
		var decodedEdge hexg.Edge
		if err := json.Unmarshal(data, &decodedEdge); err != nil || decodedEdge != e {
			t.Fatalf("%d: unmarshal edge: got %v, %v", tc.id, decodedEdge, err)
		}
		data, err = json.Marshal(v)
		if err != nil || string(data) != `"`+tc.vertexWire+`"` {
			t.Fatalf("%d: marshal vertex: got %s, %v", tc.id, data, err)
		}
		var decodedVertex hexg.Vertex
		if err := json.Unmarshal(data, &decodedVertex); err != nil || decodedVertex != v {
			t.Fatalf("%d: unmarshal vertex: got %v, %v", tc.id, decodedVertex, err)
		}
		if value, err := e.Value(); err != nil || value != tc.edgeWire {
			t.Fatalf("%d: edge SQL value: got %v, %v", tc.id, value, err)
		}
		if value, err := v.Value(); err != nil || value != tc.vertexWire {
			t.Fatalf("%d: vertex SQL value: got %v, %v", tc.id, value, err)
		}
		for _, source := range []any{tc.edgeWire, []byte(tc.edgeWire)} {
			decodedEdge = tc.h.Edge(tc.direction + 1)
			if err := decodedEdge.Scan(source); err != nil || decodedEdge != e {
				t.Fatalf("%d: scan edge: got %v, %v", tc.id, decodedEdge, err)
			}
		}
		for _, source := range []any{tc.vertexWire, []byte(tc.vertexWire)} {
			decodedVertex = tc.h.Vertex(tc.corner + 1)
			if err := decodedVertex.Scan(source); err != nil || decodedVertex != v {
				t.Fatalf("%d: scan vertex: got %v, %v", tc.id, decodedVertex, err)
			}
		}
	}
}

func TestEdgeAndVertex_InvalidSerialization(t *testing.T) {
	wantEdge := hexg.NewHex(4, -2).Edge(1)
	wantVertex := hexg.NewHex(4, -2).Vertex(5)
	maxInt := int(^uint(0) >> 1)
	for _, suffix := range []string{
		"", ":0", "+0+0:0", "+1+2+3:0", "+0+0+0:3", "+0+0+0:-1",
		"+0+0+0:00", "+0+0+0:0:1", "+0+0+0:0 ", "+00+0+0:0",
		"-0+0+0:0", "junk+0+0+0:0", "+0+0+0+0:0", "123:0",
		"+999999999999999999999999+0-999999999999999999999999:0",
		hexg.NewHex(maxInt/2+1, 0).ConciseString() + ":0",
		// These coordinates sum to zero only through integer overflow.
		hexg.NewHex(maxInt, maxInt).ConciseString() + ":0",
	} {
		e, v := wantEdge, wantVertex
		edgeWire, vertexWire := "E:"+suffix, "V:"+suffix
		if _, err := hexg.ParseEdge(edgeWire); err == nil {
			t.Errorf("accepted invalid edge %q", edgeWire)
		}
		if _, err := hexg.ParseVertex(vertexWire); err == nil {
			t.Errorf("accepted invalid vertex %q", vertexWire)
		}
		if err := e.Scan(edgeWire); err == nil || e != wantEdge {
			t.Errorf("scan invalid edge %q succeeded or mutated receiver", edgeWire)
		}
		if err := v.Scan([]byte(vertexWire)); err == nil || v != wantVertex {
			t.Errorf("scan invalid vertex %q succeeded or mutated receiver", vertexWire)
		}
		data, _ := json.Marshal(edgeWire)
		if err := json.Unmarshal(data, &e); err == nil || e != wantEdge {
			t.Errorf("invalid edge JSON %s succeeded or mutated receiver", data)
		}
		data, _ = json.Marshal(vertexWire)
		if err := json.Unmarshal(data, &v); err == nil || v != wantVertex {
			t.Errorf("invalid vertex JSON %s succeeded or mutated receiver", data)
		}
	}
	for _, source := range []any{nil, 12, true, []byte{}, "", "V:+0+0+0:0"} {
		e := wantEdge
		if err := e.Scan(source); err == nil || e != wantEdge {
			t.Errorf("invalid edge SQL source %v succeeded or mutated receiver", source)
		}
	}
	for _, source := range []any{nil, 12, true, []byte{}, "", "E:+0+0+0:0", "V:+0+0+0:2"} {
		v := wantVertex
		if err := v.Scan(source); err == nil || v != wantVertex {
			t.Errorf("invalid vertex SQL source %v succeeded or mutated receiver", source)
		}
	}
	for _, input := range []string{`null`, `12`, `{}`, `[]`, `""`, `"X:+0+0+0:0"`, `"E:+0+0+0:0`} {
		e, v := wantEdge, wantVertex
		if err := json.Unmarshal([]byte(input), &e); err == nil || e != wantEdge {
			t.Errorf("invalid edge JSON %s succeeded or mutated receiver", input)
		}
		if err := json.Unmarshal([]byte(input), &v); err == nil || v != wantVertex {
			t.Errorf("invalid vertex JSON %s succeeded or mutated receiver", input)
		}
	}
}

func TestEdgeAndVertex_ParseAtCoordinateBoundary(t *testing.T) {
	bound := int(^uint(0)>>1) / 2
	for _, owner := range []hexg.Hex{hexg.NewHex(bound, -bound), hexg.NewHex(-bound, bound)} {
		if e, err := hexg.ParseEdge("E:" + owner.ConciseString() + ":2"); err != nil || e != owner.Edge(2) {
			t.Errorf("boundary edge: got %v, %v", e, err)
		}
		if v, err := hexg.ParseVertex("V:" + owner.ConciseString() + ":1"); err != nil || v != owner.Vertex(5) {
			t.Errorf("boundary vertex: got %v, %v", v, err)
		}
	}
}
