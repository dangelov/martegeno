package main

import (
	"io/ioutil"
	"math"
	"math/rand"
	"strings"
	"unicode"

	"github.com/fogleman/gg"
)

func main() {
	s := 1000.0

	rand.Seed(1337 * 537531)

	b, _ := ioutil.ReadFile("m.go")
	t := string(b)

	dc := gg.NewContext(int(s), int(s))
	dc.SetRGB255(84, 92, 88)
	dc.Clear()

	sb := strings.Builder{}
	rs := 0
	for _, c := range t {
		if !unicode.IsSpace(c) {
			sb.WriteRune(c)
			rs++
		}
	}
	f := sb.String()

	l := math.Ceil(math.Sqrt(float64(rs)))
	sc := s / l

	dc.LoadFontFace("h.ttf", sc*0.8)

	for y := 0.0; y < l; y++ {
		for x := 0.0; x < l; x++ {
			i := y*l + x
			if i < float64(len(f)) {
				c := [][]int{{222, 232, 196}, {172, 196, 172}, {202, 219, 214}, {150, 255, 225}, {155, 213, 197}, {230, 230, 230}}[rand.Intn(6)]
				dc.SetRGB255(c[0], c[1], c[2])
				d1 := math.Pow(float64(x-l/2), 2)
				d2 := math.Pow(float64(y-l/2), 2)
				d := math.Sqrt(d1 + d2)
				if d > 10 && d < 12 {
					dc.SetRGB255(254, 138, 108)
				}
				dc.DrawStringAnchored(string(f[int(i)]), x*sc+sc/2, y*sc+sc*0.4, 0.5, 0.5)
			}

		}
	}

	dc.SavePNG("o.png")
}
