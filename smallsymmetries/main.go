package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

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
		s         = 2000.0 // Size of final image
		blocks    = 8.0
		blockSize = s / blocks
		padding   = s / 100
		thinLine  = s / 1000
		thickLine = thinLine * 4
	)

	dc := gg.NewContext(int(s), int(s))

	rand.Seed(101 * 1337)

	// Set a background color
	dc.SetColor(color.RGBA{255, 255, 255, 255})
	dc.Clear()

	dc.SetColor(color.Black)

	drawRandomLines := func(x, y, s float64) {
		for i := 0; i < 2+int(rand.Int31n(6)); i++ {
			x0 := x + rand.Float64()*s
			x1 := x + rand.Float64()*s
			y0 := y + rand.Float64()*s
			y1 := y + rand.Float64()*s
			dc.DrawLine(x0, y0, x1, y1)
			dc.SetLineWidth(thinLine)
			dc.Stroke()
		}
	}

	drawSymmetricalLines := func(x, y, s float64) {
		steps := (2 + rand.Int31n(5))
		angleStep := float64(360 / steps)
		angle := 0.0

		for i := int32(0); i < steps; i++ {
			sin, cos := math.Sincos(gg.Radians(angle))
			x0 := x + s/2 + s/2*cos
			y0 := y + s/2 + s/2*sin
			sin, cos = math.Sincos(gg.Radians(angle + 180))
			x1 := x + s/2 + s/2*cos
			y1 := y + s/2 + s/2*sin
			dc.DrawLine(x0, y0, x1, y1)
			dc.SetLineWidth(thickLine)
			dc.Stroke()

			angle += angleStep
		}
	}

	for x := 0.0; x < blocks; x++ {
		for y := 0.0; y < blocks; y++ {
			if int(x+1)%3 == 0 && int(y+1)%3 == 0 {
				drawSymmetricalLines(x*blockSize+padding, y*blockSize+padding, blockSize-padding*2)
				continue
			}

			drawRandomLines(x*blockSize+padding, y*blockSize+padding, blockSize-padding*2)
		}
	}

	dc.SavePNG("output.png")
}
