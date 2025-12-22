// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"log"

	"github.com/maloquacious/hexg"
)

func LayoutPointExample() {
	for q := -4; q <= 4; q++ {
		for r := -2; r <= 2; r++ {
			h := hexg.NewHex(q, r)
			log.Printf("Hex{q: %d, r: %d, s: %d}\n", h.Q(), h.R(), h.S())
		}
	}
}

func LayoutFlatExample() {
	for q := -4; q <= 4; q++ {
		for r := -2; r <= 2; r++ {
			h := hexg.NewHex(q, r)
			log.Printf("Hex{q: %d, r: %d, s: %d}\n", h.Q(), h.R(), h.S())
		}
	}
}
