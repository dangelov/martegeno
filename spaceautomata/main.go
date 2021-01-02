package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

type automata struct {
	cells  []bool
	rule   map[[3]bool]bool
	colors []color.Color
}

func (a *automata) advance() {
	newCells := make([]bool, len(a.cells))
	for i := 0; i < len(a.cells); i++ {
		prev := i - 1
		next := i + 1
		if next == len(a.cells) {
			next = 0
		}
		if prev == -1 {
			prev = len(a.cells) - 1
		}
		newCells[i] = a.rule[[3]bool{
			a.cells[prev],
			a.cells[i],
			a.cells[next],
		}]
	}
	a.cells = newCells
}

func (a *automata) initRandom(num int) {
	rand.Seed(time.Hour.Microseconds())
	a.cells = make([]bool, num)
	for i := 0; i < num; i++ {
		a.cells[i] = rand.Float64() < 0.5
	}
}

func (a *automata) draw(dc *gg.Context, y int, size float64) {
	for i := 0; i < len(a.cells); i++ {
		if a.cells[i] {
			dc.SetColor(a.colors[rand.Intn(len(a.colors))])
			dc.DrawRectangle(float64(i)*size, float64(y), size, size)
			dc.Fill()
		}
	}
}

func main() {
	// config
	const (
		s            = 1600.0 // Size of final image
		padding      = s / 20
		chunkSize    = 20
		borderRadius = padding * 1.5
	)

	rand.Seed(time.Hour.Microseconds())
	rule106 := map[[3]bool]bool{
		{true, true, true}:    false,
		{true, true, false}:   true,
		{true, false, true}:   true,
		{true, false, false}:  false,
		{false, true, true}:   true,
		{false, true, false}:  false,
		{false, false, true}:  true,
		{false, false, false}: false,
	}
	clrs := []color.Color{
		color.RGBA{68, 240, 210, 255},
		color.RGBA{12, 136, 222, 255},
		color.RGBA{38, 39, 112, 255},
		color.RGBA{35, 5, 69, 255},
	}
	auto := &automata{rule: rule106, colors: clrs}
	auto.initRandom(s / chunkSize)

	dc := gg.NewContext(int(s), int(s))

	// Set a background color
	dc.SetColor(color.RGBA{28, 0, 33, 255})
	dc.Clear()

	dc.DrawRoundedRectangle(padding, padding, s-padding*2, s-padding*2, borderRadius)
	dc.Clip()

	for i := 0; i < s; i += chunkSize {
		auto.draw(dc, i, float64(chunkSize))
		auto.advance()
	}

	// Save the output
	dc.SavePNG("output.png")
}
