// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg_test

import (
	"testing"

	"github.com/maloquacious/hexg"
)

func TestOffsetEvenQ(t *testing.T) {
	l := hexg.NewLayoutEvenQ(hexg.NewPoint(1, 1), hexg.NewPoint(0, 0))
	if !l.IsEvenQ() {
		t.Fatalf("oc: even-q: isEvenQ: got %v, want %v\n", l.IsEvenQ(), true)
	}
	for _, tc := range []struct {
		oc     hexg.OffsetCoord
		expect string
	}{
		{hexg.OffsetCoordFromColRow(0, 0), "+0+0+0"},
		{hexg.OffsetCoordFromColRow(1, 1), "+1+0-1"},
	} {
		h := l.HexFromOffsetCoord(tc.oc)
		if h.ConciseString() != tc.expect {
			t.Errorf("oc %q: hex: got %q, wanted %q\n", tc.oc.String(), h.ConciseString(), tc.expect)
			continue
		}
	}
}

func TestOffsetOddQ(t *testing.T) {
	l := hexg.NewLayoutOddQ(hexg.NewPoint(1, 1), hexg.NewPoint(0, 0))
	if !l.IsOddQ() {
		t.Fatalf("oc: odd-q: isOddQ: got %v, want %v\n", l.IsOddQ(), true)
	}
	for _, tc := range []struct {
		oc     hexg.OffsetCoord
		expect string
	}{
		{hexg.OffsetCoordFromColRow(0, 0), "+0+0+0"},
		{hexg.OffsetCoordFromColRow(1, 0), "+1+0-1"},
		{hexg.OffsetCoordFromColRow(2, 0), "+2-1-1"},
		{hexg.OffsetCoordFromColRow(3, 0), "+3-1-2"},
	} {
		h := l.HexFromOffsetCoord(tc.oc)
		if h.ConciseString() != tc.expect {
			t.Errorf("oc %q: hex: got %q, wanted %q\n", tc.oc.String(), h.ConciseString(), tc.expect)
			continue
		}
	}
}
