package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	stackblur "github.com/esimov/stackblur-go"
	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
)

type ego struct {
	X, Y, xOffset, yOffset, Radius, Angle, Speed float64
	Flick                                        bool
}

func (e *ego) init(xOffset, yOffset, minRadius, maxRadius float64, speed float64) {
	// The radius determines the size of our orbit
	e.Radius = minRadius + float64(rand.Int31n(int32(maxRadius-minRadius)))

	// This is the ego's starting angle on the orbit
	e.Angle = float64(rand.Int31n(180))

	// Speed controls how far the ego moves
	// during each call to travel()
	e.Speed = speed

	// Offsets help us calculate our X & Y
	// positions on the real canvas
	e.xOffset = xOffset
	e.yOffset = yOffset

	// A random flicker phase prevents having
	// all egos flick at the same time like a strobe
	e.Flick = rand.Int31n(10) < 5
}

func (e *ego) travel() {
	// Move the angle forward
	e.Angle += e.Speed
	// Adjust the X, Y coordinates
	sin, cos := math.Sincos(gg.Radians(e.Angle))
	e.X = e.Radius*cos + e.xOffset
	e.Y = e.Radius*sin + e.yOffset
}

func (e *ego) distanceTo(e2 *ego) float64 {
	return math.Sqrt(math.Pow(e2.Y-e.Y, 2) + math.Pow(e2.X-e.X, 2))
}

func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}

	if v > max {
		return max
	}

	return v
}

func clampInt(v, min, max int64) int64 {
	if v < min {
		return min
	}

	if v > max {
		return max
	}

	return v
}

func main() {
	// config
	const (
		s             = 1000.0 // Size of final image
		minRadius     = s * 0.075
		maxRadius     = (s * 0.85) / 2 // Max radius of the circle around which egos travel
		minSpeed      = 0.05
		maxSpeed      = 0.6
		egoSize       = s / 200
		distanceM     = s * 0.2   // Below distance M, start increasing opacity
		distanceN     = s * 0.1   // Below distance N, start changing colors
		distanceB     = s * 0.05  // Below distance B, start adding glow
		distanceV     = s * 0.025 // Below distance V, start flickering
		strokeWidth   = s / 150
		glowThickness = strokeWidth * 4
		glowRadius    = glowThickness / 2
	)

	// String-based seed for the random generator
	// is easier to remember between experiments
	seed := "spartacus"
	var seedNumerical int64 = 0
	for i := range seed {
		seedNumerical *= int64(seed[i])
	}
	rand.Seed(seedNumerical)

	// How many egos will be travelling?
	egos := [30]ego{}

	// Initialize all of them
	for i := range egos {
		speed := minSpeed + rand.Float64()*(maxSpeed-minSpeed)
		egos[i] = ego{}
		egos[i].init(s/2, s/2, minRadius, maxRadius, speed)
	}

	// Rotate all the egos until we complete a circle,
	// and then, for each rotation...
	for i := 0; i < 360; i++ {
		// Drawing context for the connecting lines
		lc := gg.NewContext(int(s), int(s))
		// Drawing context for the glow
		gc := gg.NewContext(int(s), int(s))
		// Drawing context for the egos
		ec := gg.NewContext(int(s), int(s))
		// Composite context
		cc := gg.NewContext(int(s), int(s))

		clr, err := colorful.Hex("#F72F4D")
		if err != nil {
			log.Fatal(err)
		}
		clr2, err := colorful.Hex("#D6FF32")
		if err != nil {
			log.Fatal(err)
		}

		// Draw all the ego connections
		for n := range egos {
			ego := &egos[n]
			ego.travel()

			// Draw the line between this and all the other egos
			for m := range egos {
				if m == n { // Don't do anything with itself
					continue
				}

				ego2 := &egos[m]

				distance := ego.distanceTo(ego2)

				fmt.Printf("distance %f", distance)

				// Too far, don't draw a line
				if distance > distanceM {
					println(" not drawn")
					continue
				}

				// Set opacity based on distance between M and N
				opacity := 1.0 - ((distance - distanceN) / (distanceM - distanceN))
				opacity = clampFloat(opacity, 0, 1.0)
				fmt.Printf(" opacity: %f", opacity)

				// Blend two colors based on distance betwee N and B
				midpoint := 1 - ((distance - distanceB) / (distanceN - distanceB))
				midpoint = clampFloat(midpoint, 0, 1.0)
				fmt.Printf(" midpoint: %f", midpoint)

				blended := clr.BlendHcl(clr2, midpoint).Clamped()

				// Create a new color based on our blend, but this time
				// with a custom Alpha opacity
				r, g, b, _ := blended.RGBA()
				a := uint8(opacity * 255)
				fmt.Printf(" a: %d", a)
				finalColor := color.NRGBA{uint8(r), uint8(g), uint8(b), a}

				// Super short distances should start flicking
				if distance < distanceV {
					ego.Flick = !ego.Flick
					if ego.Flick {
						a = 127
					}
					finalColor = color.NRGBA{uint8(r), uint8(g), uint8(b), a}
				}

				if distance < distanceB {
					// Blend two colors based on distance betwee N and B
					thickness := 1 - ((distance - distanceV) / (distanceB - distanceV))
					thickness = clampFloat(thickness, 0, 1.0)
					fmt.Printf(" thickness: %f", midpoint)

					gc.DrawLine(ego.X, ego.Y, ego2.X, ego2.Y)
					gc.SetColor(finalColor)
					gc.SetLineWidth(glowThickness * thickness)
					gc.Stroke()
				}

				lc.DrawLine(ego.X, ego.Y, ego2.X, ego2.Y)
				lc.SetColor(finalColor)
				lc.SetLineWidth(strokeWidth)
				lc.Stroke()

				println("")
			}
		}

		// Draw all the egos
		for n := range egos {
			ego := &egos[n]

			// Draw the position
			ec.SetColor(color.RGBA{11, 201, 221, 255})
			ec.DrawCircle(ego.X, ego.Y, egoSize)
			ec.Fill()
		}

		// Set a background color
		cc.SetColor(color.RGBA{0, 0, 0, 255})
		cc.Clear()

		// Bottom layer: the lines
		cc.DrawImage(lc.Image(), 0, 0)
		// On top of that, the glow
		cc.DrawImage(stackblur.Process(gc.Image(), uint32(math.Floor(glowRadius))), 0, 0)
		// And finally, the egos
		cc.DrawImage(ec.Image(), 0, 0)

		// Save the output
		cc.SavePNG(fmt.Sprintf("i-%d.png", i))
	}
}
