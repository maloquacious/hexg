// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// MarshalJSON implements json.Marshaler for Hex.
func (h Hex) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.ConciseString())
}

// UnmarshalJSON implements json.Unmarshaler for Hex.
func (h *Hex) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := scanConciseString(s)
	if err != nil {
		return err
	}
	*h = parsed
	return nil
}

// Value implements driver.Valuer for database storage.
func (h Hex) Value() (driver.Value, error) {
	return h.ConciseString(), nil
}

// Scan implements sql.Scanner for database retrieval.
func (h *Hex) Scan(src any) (err error) {
	if s, ok := src.(string); ok {
		*h, err = scanConciseString(s)
		return err
	}
	return fmt.Errorf("unknown type")
}

func scanConciseString(s string) (h Hex, err error) {
	qStart, rStart, sStart := -1, -1, -1
	for i := 0; i < len(s); i = i + 1 {
		if s[i] == '-' || s[i] == '+' {
			if qStart == -1 {
				qStart = i
			} else if rStart == -1 {
				rStart = i
			} else {
				sStart = i
			}
		}
	}
	if sStart == -1 {
		return Hex{}, fmt.Errorf("invalid hexg")
	}
	h.q, err = strconv.Atoi(s[qStart:rStart])
	if err != nil {
		return Hex{}, err
	}
	h.r, err = strconv.Atoi(s[rStart:sStart])
	if err != nil {
		return Hex{}, err
	}
	h.s, err = strconv.Atoi(s[sStart:])
	if err != nil {
		return Hex{}, err
	}
	if h.q+h.r+h.s != 0 {
		return Hex{}, fmt.Errorf("invalid hexg")
	}
	return h, nil
}

// ConciseString returns the stable edge identifier "E:<owner>:<direction>",
// where owner is a Hex.ConciseString and direction is in 0..2.
func (e Edge) ConciseString() string {
	return fmt.Sprintf("E:%s:%d", e.owner.ConciseString(), e.direction)
}

// String returns the canonical edge identifier.
func (e Edge) String() string { return e.ConciseString() }

// ConciseString returns the stable vertex identifier "V:<owner>:<slot>",
// where owner is a Hex.ConciseString and slot is Corner()-4, in 0..1.
func (v Vertex) ConciseString() string {
	return fmt.Sprintf("V:%s:%d", v.owner.ConciseString(), v.corner)
}

// String returns the canonical vertex identifier.
func (v Vertex) String() string { return v.ConciseString() }

// ParseEdge parses a canonical edge identifier. The owner must be in Hex's
// supported coordinate domain. Parsing does not normalize cylinder coordinates.
func ParseEdge(s string) (Edge, error) {
	owner, slot, err := scanFeature(s, 'E', 2)
	if err != nil {
		return Edge{}, err
	}
	return Edge{owner: owner, direction: slot}, nil
}

// ParseVertex parses a canonical vertex identifier. The owner must be in Hex's
// supported coordinate domain. Parsing does not normalize cylinder coordinates.
func ParseVertex(s string) (Vertex, error) {
	owner, slot, err := scanFeature(s, 'V', 1)
	if err != nil {
		return Vertex{}, err
	}
	return Vertex{owner: owner, corner: slot}, nil
}

// MarshalJSON implements json.Marshaler for Edge using its canonical identifier.
func (e Edge) MarshalJSON() ([]byte, error) { return json.Marshal(e.ConciseString()) }

// MarshalJSON implements json.Marshaler for Vertex using its canonical identifier.
func (v Vertex) MarshalJSON() ([]byte, error) { return json.Marshal(v.ConciseString()) }

// UnmarshalJSON implements json.Unmarshaler for Edge. Errors leave e unchanged.
func (e *Edge) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return e.Scan(s)
}

// UnmarshalJSON implements json.Unmarshaler for Vertex. Errors leave v unchanged.
func (v *Vertex) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return v.Scan(s)
}

// Value implements driver.Valuer for Edge using its canonical identifier.
func (e Edge) Value() (driver.Value, error) { return e.ConciseString(), nil }

// Value implements driver.Valuer for Vertex using its canonical identifier.
func (v Vertex) Value() (driver.Value, error) { return v.ConciseString(), nil }

// Scan implements sql.Scanner for Edge, accepting string and []byte identifiers.
// NULL and other source types return errors. Errors leave e unchanged.
func (e *Edge) Scan(src any) error {
	owner, slot, err := scanFeature(src, 'E', 2)
	if err != nil {
		return err
	}
	*e = Edge{owner: owner, direction: slot}
	return nil
}

// Scan implements sql.Scanner for Vertex, accepting string and []byte identifiers.
// NULL and other source types return errors. Errors leave v unchanged.
func (v *Vertex) Scan(src any) error {
	owner, slot, err := scanFeature(src, 'V', 1)
	if err != nil {
		return err
	}
	*v = Vertex{owner: owner, corner: slot}
	return nil
}

func scanFeature(src any, kind byte, maxSlot byte) (Hex, uint8, error) {
	var s string
	switch value := src.(type) {
	case string:
		s = value
	case []byte:
		s = string(value)
	default:
		return Hex{}, 0, fmt.Errorf("cannot scan %c identifier from %T", kind, src)
	}
	rest, ok := strings.CutPrefix(s, string(kind)+":")
	if !ok {
		return Hex{}, 0, fmt.Errorf("invalid %c identifier %q", kind, s)
	}
	encoded, slot, ok := strings.Cut(rest, ":")
	if !ok || len(slot) != 1 || slot[0] < '0' || slot[0] > '0'+maxSlot {
		return Hex{}, 0, fmt.Errorf("invalid %c identifier %q", kind, s)
	}
	owner, err := scanConciseString(encoded)
	if err != nil || owner.ConciseString() != encoded {
		return Hex{}, 0, fmt.Errorf("invalid %c identifier %q", kind, s)
	}
	bound := int(^uint(0)>>1) / 2
	if owner.q < -bound || owner.q > bound || owner.r < -bound || owner.r > bound || owner.s < -bound || owner.s > bound {
		return Hex{}, 0, fmt.Errorf("invalid %c identifier %q: owner outside supported coordinate domain", kind, s)
	}
	return owner, slot[0] - '0', nil
}
