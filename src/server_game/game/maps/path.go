package maps

import "math"

type validator interface {
	canMoveForward(int) bool
}

type Pot struct {
	X int
	Y int
}

func dist(x1, y1, x2, y2 int) int {
	x := x2 - x1
	y := y2 - y1
	return int(math.Pow(float64(x), 2) + math.Pow(float64(y), 2))
}

type Path []*Pot

var Dirs = [8]*Pot{
	{-1, -1},
	{0, -1},
	{1, -1},
	{1, 0},
	{1, 1},
	{0, 1},
	{-1, 1},
	{-1, 0},
}

func CalcDir(sx, sy, tx, ty int) int {
	x := tx - sx
	y := ty - sy
	dir := 0
	switch {
	case x <= -1 && y <= -1:
		dir = 0
	case x == 0 && y <= -1:
		dir = 1
	case x >= 1 && y <= -1:
		dir = 2
	case x >= 1 && y == 0:
		dir = 3
	case x >= 1 && y >= 1:
		dir = 4
	case x == 0 && y >= 1:
		dir = 5
	case x <= -1 && y >= 1:
		dir = 6
	case x <= -1 && y == 0:
		dir = 7
	}
	return dir
}

type _path struct {
	validator validator
	width     int
	height    int
	path      []*Pot
	hits      []bool
}

func (p *_path) valid(x, y int) bool {
	return x >= 0 && x < p.width && y >= 0 && y < p.height
}

func (p *_path) posOK(x, y int) bool {
	if !p.valid(x, y) {
		return false
	}
	pos := x + y*p.width
	if !p.validator.canMoveForward(pos) {
		p.hits[pos] = true
	}
	return !p.hits[pos]
}

func (p *_path) findNexDir(x1, y1, x2, y2 int) int {
	mindist := 100000000
	nextdir := 0
	for i, pot := range Dirs {
		x := x1 + pot.X
		y := y1 + pot.Y
		if p.posOK(x, y) {
			dist := dist(x, y, x2, y2)
			if dist < mindist {
				mindist = dist
				nextdir = i
			}
		}
	}
	p.hits[x1+y1*p.width] = true
	if mindist == 100000000 {
		return -1
	}
	return nextdir
}

func (p *_path) findPath(x1, y1, x2, y2 int) (Path, bool) {
	cnt := 10
	for !(x1 == x2 && y1 == y2) {
		dir := p.findNexDir(x1, y1, x2, y2)
		if dir >= 0 {
			// forward
			if len(p.path) >= 15 {
				return nil, false
			}
			p.path = append(p.path, &Pot{x1, y1})
			x1 += Dirs[dir].X
			y1 += Dirs[dir].Y
		} else {
			// backward
			cnt--
			if cnt < 0 {
				return nil, false
			}
			l := len(p.path)
			if l < 1 {
				return nil, false
			}
			x1 = p.path[l-1].X
			y1 = p.path[l-1].Y
			p.path = p.path[:l-1]
		}
	}
	return p.path, true
}
