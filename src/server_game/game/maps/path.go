package maps

import (
	"math"
)

type Pot struct {
	X int `json:"x"`
	Y int `json:"y"`
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
	validator func(int, int) bool
	hits      map[Pot]struct{} // use set instead to aviod mass memory gc
}

func (p *_path) valid(x, y int) bool {
	if !p.validator(x, y) {
		return false
	}
	_, ok := p.hits[Pot{x, y}]
	return !ok
}

func (p *_path) findNexDir(x1, y1, x2, y2 int) int {
	mindist := 100000000
	nextdir := 0
	for i, pot := range Dirs {
		x := x1 + pot.X
		y := y1 + pot.Y
		if p.valid(x, y) {
			dist := dist(x, y, x2, y2)
			if dist < mindist {
				mindist = dist
				nextdir = i
			}
		}
	}
	p.hits[Pot{x1, y1}] = struct{}{}
	if mindist == 100000000 {
		return -1
	}
	return nextdir
}

func (p *_path) findPath(x1, y1, x2, y2 int) (Path, bool) {
	var path []*Pot
	cnt := 10
	for !(x1 == x2 && y1 == y2) {
		dir := p.findNexDir(x1, y1, x2, y2)
		if dir >= 0 {
			// forward
			if len(path) > 15 {
				return nil, false
			}
			path = append(path, &Pot{x1, y1})
			x1 += Dirs[dir].X
			y1 += Dirs[dir].Y
		} else {
			// backward
			cnt--
			if cnt < 0 {
				return nil, false
			}
			l := len(path)
			if l < 1 {
				return nil, false
			}
			x1 = path[l-1].X
			y1 = path[l-1].Y
			path = path[:l-1]
		}
	}
	path = append(path, &Pot{x2, y2})
	return path[1:], true
}

func (p *_path) findPathBFS(x1, y1, x2, y2 int) (Path, bool) {
	p.hits[Pot{x1, y1}] = struct{}{}
	path := []*Pot{{x1, y1}}
	bfs := [][]*Pot{path}
	for len(bfs) > 0 {
		path = bfs[0]
		bfs = bfs[1:]
		x1 = path[len(path)-1].X
		y1 = path[len(path)-1].Y
		if x1 == x2 && y1 == y2 {
			return path[1:], true
		}
		if len(path) > 14 {
			return nil, false
		}
		for _, dir := range Dirs {
			x := x1 + dir.X
			y := y1 + dir.Y
			if p.valid(x, y) {
				p.hits[Pot{x, y}] = struct{}{}
				newPath := make([]*Pot, len(path), len(path)+1)
				copy(newPath, path)
				newPath = append(newPath, &Pot{x, y})
				bfs = append(bfs, newPath)
			}
		}
	}
	return nil, false
}
