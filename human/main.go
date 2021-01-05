package main

import (
	"image/color"
	"math/rand"
	"strconv"

	"github.com/fogleman/gg"
)

// Direction is where a cell is facing
type Direction int

const (
	dirN    Direction = 0
	dirE    Direction = 90
	dirS    Direction = 180
	dirW    Direction = 270
	dirNone Direction = -1
)

func (d Direction) isPerpendicular(dirs map[Direction]bool) bool {
	for dir := range dirs {
		if dir == d || dir.opposite() == d {
			return false
		}
	}

	return true
}

func (d Direction) String() string {
	switch d {
	case dirN:
		return "N"
	case dirS:
		return "S"
	case dirW:
		return "W"
	default:
		return "E"
	}
}

func (d Direction) opposite() Direction {
	return (d + 180) % 360
}

func (d Direction) coordinatesDelta() (int, int) {
	switch d {
	case dirN:
		return 0, -1
	case dirS:
		return 0, 1
	case dirW:
		return -1, 0
	default: // dirE
		return 1, 0
	}
}

// Position is whether the cell is under or above
type Position uint

const (
	posAbove Position = iota
	posBelow
)

var cellsSoFar = 0

// Cell represents a maze cell
type Cell struct {
	dirs    map[Direction]bool
	pos     Position
	visited bool
	weaved  bool
	num     int
}

func (c Cell) drawOn(dc *gg.Context, x, y, size float64) {
	offset := size / 5.0
	walls := map[Direction]bool{dirN: true, dirS: true, dirW: true, dirE: true}

	// Exit locations and dimensions
	exits := map[Direction][2][2]float64{
		dirN: {{x + offset, y}, {size - offset*2, offset}},
		dirS: {{x + offset, y + size - offset}, {size - offset*2, offset}},
		dirW: {{x, y + offset}, {offset, size - offset*2}},
		dirE: {{x + size - offset, y + offset}, {offset, size - offset*2}},
	}

	// Cell core
	core := [2][2]float64{{x + offset, y + offset}, {size - offset*2, size - offset*2}}

	cellColor := color.RGBA{249, 205, 173, 255}
	patt := gg.NewSolidPattern(cellColor)
	dc.SetFillStyle(patt)

	// Draw all our exits
	for dir := range c.dirs {
		if _, ok := exits[dir]; ok {
			dc.DrawRectangle(exits[dir][0][0], exits[dir][0][1], exits[dir][1][0], exits[dir][1][1])
			dc.Fill()
		}

		// And also remove the wall
		delete(walls, dir)
	}

	// Generate a horizontal or vertical gradient
	// based on coordinates
	gradientForCoords := func(coords [2][2]float64, dir Direction) gg.Gradient {
		x0, y0, w, h := coords[0][0], coords[0][1], coords[1][0], coords[1][1]
		x1, y1 := x0+w, y0+h

		// Gradients facing north or west need to be reversed
		if dir == dirN {
			y0, y1 = y1, y0
		}
		if dir == dirW {
			x0, x1 = x1, x0
		}

		// Horizontal gradient
		grad := gg.NewLinearGradient(x0, y0, x1, y0)
		// Vertical gradient
		if dir == dirN || dir == dirS || (h > w && dir == dirNone) {
			grad = gg.NewLinearGradient(x0, y0, x0, y1)
		}
		return grad
	}

	// Draw the core rectangle.
	dc.DrawRectangle(core[0][0], core[0][1], core[1][0], core[1][0])

	// Solid fill by default.
	dc.SetFillStyle(gg.NewSolidPattern(cellColor))

	// Weaved cells get special effects,
	// and need to draw the passages underneath.
	if c.weaved {
		grad := gradientForCoords(core, dirNone)
		if !dirN.isPerpendicular(c.dirs) {
			grad = gradientForCoords(core, dirN)
		}
		grad.AddColorStop(0, cellColor)
		grad.AddColorStop(0.4, color.RGBA{247, 217, 195, 255})
		grad.AddColorStop(0.6, color.RGBA{212, 164, 131, 255})
		grad.AddColorStop(1, cellColor)
		dc.SetFillStyle(grad)
		dc.Fill()
		dc.SetColor(color.Black)

		// Our walls are "exits" for the cell underneath,
		// so we can draw them and style them differently.
		for wall := range walls {
			if _, ok := exits[wall]; ok {
				grad := gradientForCoords(exits[wall], wall.opposite())
				grad.AddColorStop(0, cellColor)
				grad.AddColorStop(1, color.RGBA{184, 144, 116, 255})
				dc.SetFillStyle(grad)
				dc.DrawRectangle(exits[wall][0][0], exits[wall][0][1], exits[wall][1][0], exits[wall][1][1])
				dc.Fill()
			}
		}

	} else {
		dc.Fill()
		dc.SetColor(color.Black)
	}

	// dc.SetColor(color.Black)
	// dc.DrawString(c.String(), x+size/2-offset, y+size/2-offset)
}

func (c Cell) String() string {
	text := strconv.Itoa(c.num)
	for dir := range c.dirs {
		text = text + dir.String()
	}
	return text
}

// Maze represents our whole maze
type Maze struct {
	cells         [][]*Cell
	width, height int
}

func newMaze(width, height int) *Maze {
	m := &Maze{}
	m.width = width
	m.height = height
	m.cells = make([][]*Cell, width, width)
	for x := 0; x < width; x++ {
		m.cells[x] = make([]*Cell, height)
		for y := 0; y < height; y++ {
			cell := &Cell{}
			cell.dirs = make(map[Direction]bool, 4)
			m.cells[x][y] = cell
		}
	}
	m.visitCell(width/2, height/2, dirNone)
	return m
}

func (m *Maze) visitCell(x, y int, dir Direction) {
	cellsSoFar++
	cell := m.cells[x][y]
	if dir != dirNone {
		cell.dirs[dir] = true
	}
	cell.visited = true
	cell.num = cellsSoFar

	// Get all the possible directions and randomize them
	directions := [4]Direction{dirE, dirW, dirS, dirN}
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})
	// Then, return the first available cell
	for _, newDir := range directions {
		dX, dY := newDir.coordinatesDelta()
		newX, newY := x+dX, y+dY
		if m.canVisit(newX, newY) {
			nextCell := m.cells[newX][newY]

			// Weave when possible
			if nextCell.visited {
				if newDir.isPerpendicular(nextCell.dirs) && m.canVisit(newX+dX, newY+dY) && !m.cells[newX+dX][newY+dY].visited {
					cell.dirs[newDir] = true
					nextCell.weaved = true
					m.visitCell(newX+dX, newY+dY, newDir.opposite())
					continue
				}
			}

			// Or if it's free to go
			if !nextCell.visited {
				cell.dirs[newDir] = true
				m.visitCell(newX, newY, newDir.opposite())
			}
		}
	}
}

func (m *Maze) canVisit(x, y int) bool {
	// Is it out of bounds?
	if x < 0 || x >= len(m.cells) || y < 0 || y >= len(m.cells[0]) {
		return false
	}

	return true
}

func (m *Maze) canWeave(dir Direction, x, y int) bool {
	return false
	cell := m.cells[x][y]
	// If both cells are horizontal or vertical,
	// then we can't weave
	if dir.isPerpendicular(cell.dirs) {
		return false
	}

	// Where would we end up if we tried to move
	// past x and y?
	newX, newY := dir.coordinatesDelta()
	newX = newX*2 + x // Twice, to get out on the other side
	newY = newY*2 + y

	return m.canVisit(newX, newY)
}

func (m *Maze) drawOn(dc *gg.Context) {
	scale := float64(dc.Width()) / float64(m.width)

	for x := 0; x < m.width; x++ {
		for y := 0; y < m.height; y++ {
			m.cells[x][y].drawOn(dc, float64(x)*scale, float64(y)*scale, scale)
		}
	}
}

func main() {
	const (
		s        = 5000
		mazeSize = 7
	)
	dc := gg.NewContext(int(s), int(s))

	rand.Seed(1073 / (1337 + 42))

	// Set a background color
	dc.SetColor(color.RGBA{131, 175, 155, 255})
	dc.Clear()

	// MAZE WITH A HEART IN THE CENTER
	maze := newMaze(mazeSize, mazeSize)
	dc.SetColor(color.Black)
	maze.drawOn(dc)

	// Draw a heart
	dc.SetColor(color.RGBA{254, 67, 101, 255})
	dc.LoadFontFace("truetype/freefont/FreeSans.ttf", 296)
	dc.DrawString("♥", s*0.478, s*0.52)

	// Save the output
	dc.SavePNG("output.png")
}
