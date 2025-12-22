// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import "fmt"

// String returns the coordinates formatted as "(q, r, s)".
func (h Hex) String() string {
	return fmt.Sprintf("(%d, %d, %d)", h.q, h.r, h.s)
}

// ConciseString returns the coordinates with signs.
// It returns the coordinates formatted as "+q+r+s".
func (h Hex) ConciseString() string {
	return fmt.Sprintf("%+d%+d%+d", h.q, h.r, h.s)
}
