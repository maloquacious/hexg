// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"testing"

	"github.com/maloquacious/hexg"
)

func TestTribeNetRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string // expected output after round-trip
		wantErr  bool
	}{
		// Valid round-trip test cases
		{"Top-left", "AA 0101", "AA 0101", false},
		{"Mid-grid", "BC 0812", "BC 0812", false},
		{"Lower-right", "ZZ 3021", "ZZ 3021", false},
		{"Random valid", "JK 0609", "JK 0609", false},

		// Valid bounds of grids
		{"Grid AA upper-left", "AA 0101", "AA 0101", false},
		{"Grid AA lower-right", "AA 3021", "AA 3021", false},
		{"Grid AZ upper-left", "AZ 0101", "AZ 0101", false},
		{"Grid AZ lower-right", "AZ 3021", "AZ 3021", false},
		{"Grid ZA upper-left", "ZA 0101", "ZA 0101", false},
		{"Grid ZA lower-right", "ZA 3021", "ZA 3021", false},
		{"Grid ZZ upper-left", "ZZ 0101", "ZZ 0101", false},
		{"Grid ZZ lower-right", "ZZ 3021", "ZZ 3021", false},

		// invalid row or column
		{"BC 0021", "BC 0021", "", true},
		{"BC 0800", "BC 0800", "", true},
		{"BC 0824", "BC 0824", "", true},
		{"BC 3112", "BC 3112", "", true},

		// Edge cases (invalid formats)
		{"Too short", "A 0102", "", true},
		{"No space", "AA0102", "", true},
		{"Bad grid row", "1A 0102", "", true},
		{"Bad grid col", "A1 0102", "", true},
		{"Bad subcol", "AA 0001", "", true},
		{"Bad subrow", "AA 0100", "", true},
		{"Subcol too big", "AA 3101", "", true},
		{"Subrow too big", "AA 0122", "", true},

		// Out of grid bounds (ZZ + 1)
		{"Grid row overflow", "Z[ 0101", "", true},
		{"Grid col overflow", "[Z 0101", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oc, err := hexg.NewTribeNetOffsetCoord(tt.input)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("%s: NewTribeNetOffsetCoord(%q) error = %v, wantErr %v", tt.name, tt.input, err, tt.wantErr)
				}
				return
			}

			// Round-trip
			got, err := oc.ToTribeNetCoord()
			if err != nil {
				t.Errorf("%s: ToTribeNetCoord() error = %v", tt.name, err)
				return
			}

			if got != tt.expected {
				t.Errorf("%s: Round-trip mismatch: got = %q, want = %q", tt.name, got, tt.expected)
			}
		})
	}
}
