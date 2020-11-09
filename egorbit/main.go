package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type ego struct {
	X, Y, Radius, Angle, Speed float64
}

func (v *ego) init(maxRadius float64, speed float64) {
	v.Radius = float64(rand.Int31n(int32(maxRadius)))
	v.Angle = float64(rand.Int31n(180))
	v.Speed = speed
}

func (v *ego) travel() {
	// Move the angle forward
	v.Angle += v.Speed
	// Adjust the X, Y coordinates
	sin, cos := math.Sincos(gg.Radians(v.Angle))
	v.X = v.Radius * cos
	v.Y = v.Radius * sin
}

func main() {
	// config
	const (
		s         = 2000.0               // Size of final image
		maxRadius = (s - (s * 0.15)) / 2 // Max radius of the circle around which egos travel
	)

	// How many egos will be travelling?
	egos := [30]ego{}

	// Initialize all of them
	for i := range egos {
		egos[i] = ego{}
		egos[i].init(maxRadius, 0.5+rand.Float64()*6.0)
	}

	// Rotate all the egos until we complete a circle,
	// and then, for each rotation...
	for i := 0; i < 360; i++ {
		// Start a drawing context
		dc := gg.NewContext(int(s), int(s))

		// Set a background color
		dc.SetColor(color.RGBA{0, 0, 0, 255})
		dc.Clear()

		// Draw all the egos
		for n := range egos {
			ego := &egos[n]
			ego.travel()

			// Draw the orbit
			dc.SetColor(color.NRGBA{255, 255, 255, 60})
			dc.DrawCircle(s/2, s/2, ego.Radius)
			dc.SetLineWidth(3)
			dc.Stroke()

			// Draw the position
			dc.SetColor(color.RGBA{255, 0, 0, 255})
			dc.DrawCircle(ego.X+s/2, ego.Y+s/2, 12)
			dc.Fill()
		}

		// Save the output
		dc.SavePNG(fmt.Sprintf("i-%d.png", i))
	}
}
