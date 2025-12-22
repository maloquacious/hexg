// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/maloquacious/hexg"
)

func main() {
	var (
		orientation = flag.String("orientation", "pointy", "hex orientation: pointy or flat")
		sizeW       = flag.Float64("size-w", 30, "hex width (before art-sprite adjustment)")
		sizeH       = flag.Float64("size-h", 30, "hex height (before art-sprite adjustment)")
		originX     = flag.Float64("origin-x", 200, "origin X coordinate")
		originY     = flag.Float64("origin-y", 200, "origin Y coordinate")
		artSprite   = flag.Bool("art-sprite", false, "use art-friendly sprite sizing")
		radius      = flag.Int("radius", 2, "hex grid radius")
		output      = flag.String("output", "testdata/hexgrid.png", "output filename")
	)
	flag.Parse()

	var offset hexg.LayoutOffset
	switch *orientation {
	case "pointy":
		offset = hexg.OddR
	case "flat":
		offset = hexg.OddQ
	default:
		log.Fatalf("invalid orientation: %s (use 'pointy' or 'flat')", *orientation)
	}

	w, h := *sizeW, *sizeH
	if *artSprite {
		if offset == hexg.OddR || offset == hexg.EvenR {
			// pointy-top: w/2, h/sqrt(3)
			w = *sizeW / 2
			h = *sizeH / math.Sqrt(3)
		} else {
			// flat-top: w/sqrt(3), h/2
			w = *sizeW / math.Sqrt(3)
			h = *sizeH / 2
		}
	}

	size := hexg.Point{X: w, Y: h}
	origin := hexg.Point{X: *originX, Y: *originY}
	layout := hexg.NewLayout(offset, size, origin)

	hexes := hexg.NewHex(0, 0).Spiral(*radius)

	if err := renderGrid(layout, hexes, *output); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %s (%s, size=%.1f,%.1f, origin=%.1f,%.1f, radius=%d)\n",
		*output, *orientation, w, h, *originX, *originY, *radius)
}

func renderGrid(layout hexg.Layout, hexes []hexg.Hex, filename string) error {
	imgWidth, imgHeight := 400, 400

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// fill background
	bgColor := color.RGBA{R: 240, G: 240, B: 240, A: 255}
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			img.Set(x, y, bgColor)
		}
	}

	hexFill := color.RGBA{R: 200, G: 220, B: 255, A: 255}
	hexStroke := color.RGBA{R: 60, G: 60, B: 100, A: 255}
	textColor := color.RGBA{R: 40, G: 40, B: 80, A: 255}

	for _, h := range hexes {
		corners := layout.PolygonCorners(h)
		fillPolygon(img, corners[:], hexFill)
		drawPolygon(img, corners[:], hexStroke)

		center := layout.HexToPixel(h)
		label := fmt.Sprintf("%d,%d,%d", h.Q(), h.R(), h.S())
		drawText(img, int(center.X), int(center.Y), label, textColor)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func fillPolygon(img *image.RGBA, corners []hexg.Point, c color.RGBA) {
	minY, maxY := int(corners[0].Y), int(corners[0].Y)
	for _, p := range corners {
		if int(p.Y) < minY {
			minY = int(p.Y)
		}
		if int(p.Y) > maxY {
			maxY = int(p.Y)
		}
	}

	for y := minY; y <= maxY; y++ {
		var intersections []float64
		n := len(corners)
		for i := 0; i < n; i++ {
			j := (i + 1) % n
			y1, y2 := corners[i].Y, corners[j].Y
			if (y1 <= float64(y) && y2 > float64(y)) || (y2 <= float64(y) && y1 > float64(y)) {
				x := corners[i].X + (float64(y)-y1)/(y2-y1)*(corners[j].X-corners[i].X)
				intersections = append(intersections, x)
			}
		}
		for i := 0; i < len(intersections)-1; i++ {
			for j := i + 1; j < len(intersections); j++ {
				if intersections[i] > intersections[j] {
					intersections[i], intersections[j] = intersections[j], intersections[i]
				}
			}
		}
		for i := 0; i+1 < len(intersections); i += 2 {
			for x := int(intersections[i]); x <= int(intersections[i+1]); x++ {
				if x >= 0 && x < img.Bounds().Max.X && y >= 0 && y < img.Bounds().Max.Y {
					img.Set(x, y, c)
				}
			}
		}
	}
}

func drawPolygon(img *image.RGBA, corners []hexg.Point, c color.RGBA) {
	n := len(corners)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		drawLine(img, int(corners[i].X), int(corners[i].Y), int(corners[j].X), int(corners[j].Y), c)
	}
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx, sy := 1, 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}
	err := dx + dy
	for {
		if x0 >= 0 && x0 < img.Bounds().Max.X && y0 >= 0 && y0 < img.Bounds().Max.Y {
			img.Set(x0, y0, c)
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

var font3x5 = map[rune][]string{
	'0': {"###", "# #", "# #", "# #", "###"},
	'1': {" # ", "## ", " # ", " # ", "###"},
	'2': {"###", "  #", "###", "#  ", "###"},
	'3': {"###", "  #", "###", "  #", "###"},
	'4': {"# #", "# #", "###", "  #", "  #"},
	'5': {"###", "#  ", "###", "  #", "###"},
	'6': {"###", "#  ", "###", "# #", "###"},
	'7': {"###", "  #", " # ", " # ", " # "},
	'8': {"###", "# #", "###", "# #", "###"},
	'9': {"###", "# #", "###", "  #", "###"},
	',': {"   ", "   ", "   ", " # ", "#  "},
	'-': {"   ", "   ", "###", "   ", "   "},
}

func drawText(img *image.RGBA, cx, cy int, text string, c color.RGBA) {
	charWidth := 4
	totalWidth := len(text) * charWidth
	startX := cx - totalWidth/2
	startY := cy - 2

	for i, ch := range text {
		pattern, ok := font3x5[ch]
		if !ok {
			continue
		}
		ox := startX + i*charWidth
		for row, line := range pattern {
			for col, pixel := range line {
				if pixel == '#' {
					px, py := ox+col, startY+row
					if px >= 0 && px < img.Bounds().Max.X && py >= 0 && py < img.Bounds().Max.Y {
						img.Set(px, py, c)
					}
				}
			}
		}
	}
}
