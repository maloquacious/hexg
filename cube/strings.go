// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cube

import "fmt"

// ConciseString returns the coordinates with signs.
// It returns the coordinates formatted as (+q+r+s).
func (h Cube) ConciseString() string {
	return fmt.Sprintf("%+d%+d%+d", h.q, h.r, h.s)
}
