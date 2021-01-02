package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/fogleman/ease"
	"github.com/fogleman/gg"
)

type circle struct {
	x, y, r, step float64
	c             color.Color
}

func main() {
	// config
	const (
		s             = 1000.0   // Size of final image
		padding       = s * 0.15 // Padding around the circle
		animationStep = 0.05

		maxCircleRadius = s / 30
	)

	rand.Seed(time.Hour.Microseconds())

	clrs := [][]uint8{{204, 178, 76}, {247, 214, 131}, {255, 253, 192}, {255, 255, 253}, {69, 125, 151}}

	// Times goes from 0 to 1 and back to 0
	times := []float64{}
	for i := 0.0; i < 1.01; i += animationStep {
		times = append(times, i)
	}
	for i := 0.9; i > 0.01; i -= animationStep {
		times = append(times, i)
	}

	circles := []*circle{}

	// Generate all the circles
	for x := 0.0; x < s; x += maxCircleRadius + 2 {
		for y := 0.0; y < s; y += maxCircleRadius + 2 {
			for c := 0; c < len(clrs); c++ {
				r := rand.Float64() * maxCircleRadius
				step := rand.Float64()
				c := color.RGBA{clrs[c][0], clrs[c][1], clrs[c][2], 255}

				circles = append(circles, &circle{x, y, r, step, c})
			}
		}
	}

	// Go through the times and draw the circles
	for i := 0; i < len(times); i++ {

		dc := gg.NewContext(int(s), int(s))

		// Set a background color
		dc.SetColor(color.RGBA{clrs[1][0], clrs[1][1], clrs[1][2], 255})
		dc.Clear()

		for o := 0; o < len(circles); o++ {
			time := (times[i] + circles[o].step)
			if time > 1 {
				time = 1 - (time - 1)
			}
			r := ease.InOutCubic(time)*circles[o].r*0.7 + 0.3
			makeCircle(dc, circles[o].x, circles[o].y, r, circles[o].c)
		}

		// Save the output
		dc.SavePNG(fmt.Sprintf("i-%d.png", i))
	}
}

func makeCircle(dc *gg.Context, x, y, r float64, color color.Color) {
	dc.SetColor(color)
	dc.DrawCircle(x, y, r)
	dc.Fill()
}
