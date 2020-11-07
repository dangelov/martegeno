package main

import (
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"
)

func getColor(theme string) (uint8, uint8, uint8) {
	palletes := map[string][][]uint8{
		"mars": {{214, 129, 111}, {161, 131, 119}, {163, 149, 145}, {152, 181, 165}, {219, 205, 182}},
	}

	color := palletes[theme][rand.Intn(len(palletes[theme]))]

	return color[0], color[1], color[2]
}

func main() {
	// config
	const (
		s       = 600.0    // Size of final image
		padding = s * 0.15 // Padding around the circle

		lines = 300.0  // Number of lines to draw
		angle = -120.0 //
	)

	// Calculate the radius based on the radius and padding
	r := s/2.0 - padding

	dc := gg.NewContext(int(s), int(s))

	// Set a background color
	dc.SetColor(color.RGBA{255, 255, 255, 255})
	dc.Clear()

	// Use a circle as a mask
	dc.DrawCircle(s/2.0, s/2.0, r)
	dc.Clip()

	// Make sure everything we draw from now on is rotated
	// Only affects the lines really
	dc.Rotate(0.2)

	// Draw each line
	for i := 0; i < int(lines); i++ {
		// Lines of random lengths and locations
		y := float64(rand.Int31n(int32(s)))
		x1 := float64(rand.Int31n(int32(s)))
		dc.DrawLine(x1, y, s, y)
		// Random color from the pallete
		r, g, b := getColor("mars")
		dc.SetColor(color.RGBA{r, g, b, 255})
		// Random line widths
		dc.SetLineWidth(float64(rand.Int31n(100)) / 10)
		dc.Stroke()
	}

	dc.SavePNG("output.png")
}
