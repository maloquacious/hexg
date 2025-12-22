// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"testing"

	"github.com/maloquacious/hexg"
)

func TestHex_String(t *testing.T) {
	for _, tc := range []struct {
		h    hexg.Hex
		want string
	}{
		{hexg.NewHex(0, 0), "(0, 0, 0)"},
		{hexg.NewHex(1, 2), "(1, 2, -3)"},
		{hexg.NewHex(-1, -2), "(-1, -2, 3)"},
		{hexg.NewHex(10, -5), "(10, -5, -5)"},
	} {
		if got := tc.h.String(); got != tc.want {
			t.Errorf("Hex%v.String() = %q, want %q", tc.h, got, tc.want)
		}
	}
}

func TestHex_ConciseString(t *testing.T) {
	for _, tc := range []struct {
		h    hexg.Hex
		want string
	}{
		{hexg.NewHex(0, 0), "+0+0+0"},
		{hexg.NewHex(1, 2), "+1+2-3"},
		{hexg.NewHex(-1, -2), "-1-2+3"},
		{hexg.NewHex(10, -5), "+10-5-5"},
	} {
		if got := tc.h.ConciseString(); got != tc.want {
			t.Errorf("Hex%v.ConciseString() = %q, want %q", tc.h, got, tc.want)
		}
	}
}
