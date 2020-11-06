package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/anthonynsimon/bild/paint"
	"github.com/fogleman/gg"
)

func drawCurve(dc *gg.Context) {
	dc.SetRGBA(0, 0, 0, 0)
	dc.FillPreserve()
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(10)
	dc.Stroke()
}

func drawPoints(dc *gg.Context) {
	dc.SetRGBA(1, 0, 0, 0.5)
	dc.SetLineWidth(2)
	dc.Stroke()
}

func getColor(theme string) (uint8, uint8, uint8) {
	palletes := map[string][][]uint8{
		"zen":             {{63, 63, 63}, {143, 175, 159}, {220, 163, 163}, {240, 223, 175}, {239, 239, 239}},
		"monokai":         {{39, 40, 34}, {249, 38, 114}, {102, 217, 239}, {166, 226, 46}, {253, 151, 31}},
		"goldfish":        {{105, 210, 231}, {167, 219, 216}, {224, 228, 204}, {243, 134, 48}, {250, 105, 0}},
		"???":             {{254, 67, 101}, {252, 157, 154}, {249, 205, 173}, {200, 200, 169}, {131, 175, 155}},
		"thought":         {{236, 208, 120}, {217, 91, 67}, {192, 41, 66}, {84, 36, 55}, {83, 119, 122}},
		"adrift":          {{207, 240, 158}, {168, 219, 168}, {121, 189, 154}, {59, 134, 134}, {11, 72, 107}},
		"cheer-emo":       {{85, 98, 112}, {78, 205, 196}, {199, 244, 100}, {255, 107, 107}, {196, 77, 88}},
		"cake":            {{119, 79, 56}, {224, 142, 121}, {241, 212, 175}, {236, 229, 206}, {197, 224, 220}},
		"terra":           {{232, 221, 203}, {205, 179, 128}, {3, 101, 100}, {3, 54, 73}, {3, 22, 52}},
		"melon":           {{209, 242, 165}, {239, 250, 180}, {255, 196, 140}, {255, 159, 128}, {245, 105, 145}},
		"curious":         {{73, 10, 61}, {189, 21, 80}, {233, 127, 2}, {248, 202, 0}, {138, 155, 15}},
		"pancake":         {{89, 79, 79}, {84, 121, 128}, {69, 173, 168}, {157, 224, 173}, {229, 252, 194}},
		"fire-ocean":      {{0, 160, 176}, {106, 74, 60}, {204, 51, 63}, {235, 104, 65}, {237, 201, 81}},
		"japanese-lovers": {{233, 78, 119}, {214, 129, 137}, {198, 164, 154}, {198, 229, 217}, {244, 234, 213}},
		"compatible":      {{63, 184, 175}, {127, 199, 175}, {218, 216, 167}, {255, 158, 157}, {255, 61, 127}},
		"friends":         {{217, 206, 178}, {148, 140, 117}, {213, 222, 217}, {122, 106, 83}, {153, 178, 183}},
	}

	color := palletes[theme][rand.Intn(len(palletes[theme]))]

	return color[0], color[1], color[2]
}

func main() {
	// config
	const (
		s       = 6000.0   // Size of final image
		padding = s * 0.15 // Padding around the circle

		lines            = 40.0     // Number of lines to draw
		stretchX         = s * 0.01 // How to stretch (or in this case compress) the mid points's X of the bezier curves
		stretchY         = s * 0.4  // Same but for Y
		debug            = false    // If enabled, draws the start, end and control points of the bezier curves
		angleOffsetStart = -20.0    // Controls the angle offset of where the line begins
		angleOffsetEnd   = -20.0    // Controls the angle offset of where the line ends

		angleStart = 110.0 // Start angle
		angleEnd   = 190.0 // End angle
	)

	// A list of themes to use for the composite
	// Must have N or more items, where N is the product of X*Y from the composite code
	makeThemes := []string{"japanese-lovers", "cake", "compatible", "goldfish", "cheer-emo", "melon"}
	for _, theme := range makeThemes {
		angle := angleStart
		stepAngle := (angle - angleEnd) / lines
		stepX := stretchX / lines
		r := s/2.0 - padding

		// Draw the circle
		dc := gg.NewContext(int(s), int(s))
		dc.DrawCircle(s/2.0, s/2.0, r)
		dc.SetRGB(1, 1, 1)
		dc.Fill()

		// Draw each line
		for i := 0; i < int(lines); i++ {
			angle += stepAngle
			// Bezier start point
			sin, cos := math.Sincos(gg.Degrees(angle*-1 + angleOffsetStart))
			startX := s/2 + r*cos
			startY := s/2 + r*sin
			// Bezier mid point 1
			midX1 := s/2 + (lines/2*-1*stepX + stepX*float64(i))
			midY1 := s/2 - stretchY/2
			// Bezier mid point 2
			midX2 := midX1
			midY2 := s/2 + stretchY/2
			// Bezier end point
			sin, cos = math.Sincos(gg.Degrees(angle + angleOffsetEnd))
			endX := s/2 + r*cos
			endY := s/2 + r*sin
			// draw the curve
			dc.MoveTo(startX, startY)
			dc.CubicTo(midX1, midY1, midX2, midY2, endX, endY)
			drawCurve(dc)

			// If debug is on, draw lines and circles marking the bezier curve
			if debug {
				dc.MoveTo(startX, startY)
				dc.LineTo(midX1, midY1)
				dc.LineTo(midX2, midY2)
				dc.LineTo(endX, endY)
				drawPoints(dc)
			}
		}

		// Flood fill points randomly, if the point is white
		// This leaves rough edges, but it doesn't matter at the resolution we're using
		img := dc.Image()
		for x := int(padding); x < s; x++ {
			for y := int(padding); y < s; y++ {
				r, g, b, a := img.At(x, y).RGBA()
				if r == 65535 && g == 65535 && b == 65535 && a == 65535 {
					newR, newG, newB := getColor(theme)
					img = paint.FloodFill(img, image.Point{x, y}, color.RGBA{newR, newG, newB, 255}, 15)
				}
			}
		}

		gg.SavePNG(fmt.Sprintf("out-%s.png", theme), img)
	}

	// Now, create the composite, choosing the matrix size
	const (
		sizeX = 2
		sizeY = 3
	)
	dc := gg.NewContext(int(s)*sizeX, int(s)*sizeY)

	// Paint the background in the composite
	dc.SetHexColor("#cbf1e9")
	dc.DrawRectangle(0, 0, s*sizeX, s*sizeY)
	dc.Fill()

	// Create an X * Y matrix of the pieces
	for y := 0; y < sizeY; y++ {
		for x := 0; x < sizeX; x++ {
			theme := makeThemes[x+y*2%len(makeThemes)]
			im, _ := gg.LoadImage(fmt.Sprintf("out-%s.png", theme))

			dc.DrawImage(im, x*s, y*s)
		}
	}
	dc.SavePNG("composite.png")
}
