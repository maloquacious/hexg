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
