// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"testing"

	"github.com/maloquacious/hexg"
)

func TestNewFractionalHex(t *testing.T) {
	fh := hexg.NewFractionalHex(1.5, 2.5, -4.0)
	rounded := fh.Round()
	q, r, s := rounded.QRS()
	if q+r+s != 0 {
		t.Errorf("NewFractionalHex(1.5, 2.5, -4.0).Round() = (%d, %d, %d), constraint q+r+s=0 violated", q, r, s)
	}
}

func TestFractionalHex_Round(t *testing.T) {
	for _, tc := range []struct {
		name    string
		q, r, s float64
		wantQ   int
		wantR   int
		wantS   int
	}{
		{"origin", 0.0, 0.0, 0.0, 0, 0, 0},
		{"exact integer", 1.0, 2.0, -3.0, 1, 2, -3},
		{"round q", 0.6, 0.2, -0.8, 1, 0, -1},
		{"round r", 0.2, 0.6, -0.8, 0, 1, -1},
		{"round s", 0.2, -0.8, 0.6, 0, -1, 1},
		{"negative coords", -1.4, -0.3, 1.7, -2, 0, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fh := hexg.NewFractionalHex(tc.q, tc.r, tc.s)
			got := fh.Round()
			gotQ, gotR, gotS := got.QRS()
			if gotQ != tc.wantQ || gotR != tc.wantR || gotS != tc.wantS {
				t.Errorf("FractionalHex{%g, %g, %g}.Round() = (%d, %d, %d), want (%d, %d, %d)",
					tc.q, tc.r, tc.s, gotQ, gotR, gotS, tc.wantQ, tc.wantR, tc.wantS)
			}
			if gotQ+gotR+gotS != 0 {
				t.Errorf("constraint q+r+s=0 violated: %d+%d+%d=%d", gotQ, gotR, gotS, gotQ+gotR+gotS)
			}
		})
	}
}

func TestFractionalHex_Lerp(t *testing.T) {
	a := hexg.NewFractionalHex(0.0, 0.0, 0.0)
	b := hexg.NewFractionalHex(4.0, -2.0, -2.0)

	for _, tc := range []struct {
		t     float64
		wantQ float64
		wantR float64
		wantS float64
	}{
		{0.0, 0.0, 0.0, 0.0},
		{0.5, 2.0, -1.0, -1.0},
		{1.0, 4.0, -2.0, -2.0},
		{0.25, 1.0, -0.5, -0.5},
	} {
		got := a.Lerp(b, tc.t)
		rounded := got.Round()
		q, r, s := rounded.QRS()
		if q+r+s != 0 {
			t.Errorf("Lerp result at t=%g violates constraint: %d+%d+%d != 0", tc.t, q, r, s)
		}
	}
}

func TestHex_Lerp(t *testing.T) {
	a := hexg.NewHex(0, 0)
	b := hexg.NewHex(4, -2)

	mid := a.Lerp(b, 0.5)
	rounded := mid.Round()
	q, r, s := rounded.QRS()
	if q+r+s != 0 {
		t.Errorf("Hex.Lerp constraint violated: %d+%d+%d != 0", q, r, s)
	}
	if q != 2 || r != -1 || s != -1 {
		t.Errorf("Hex.Lerp(0.5) = (%d, %d, %d), want (2, -1, -1)", q, r, s)
	}
}
