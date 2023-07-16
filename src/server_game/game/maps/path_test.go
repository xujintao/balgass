package maps

import (
	"testing"
)

type testMap struct {
	buf    []bool
	width  int
	height int
}

func (m *testMap) canMoveForward(pos int) bool {
	return !m.buf[pos]
}

type testData struct {
	ok     bool
	width  int
	height int
	begin  Pot
	end    Pot
	pots   []*Pot
}

func TestFindPath(t *testing.T) {
	testData := []*testData{
		{true, 8, 8, Pot{1, 1}, Pot{6, 1}, []*Pot{{2, 0}, {2, 1}, {2, 2}, {1, 2}}},
		{false, 8, 8, Pot{1, 1}, Pot{6, 1}, []*Pot{{2, 0}, {2, 1}, {2, 2}, {1, 2}, {0, 2}}},
	}
	for i, td := range testData {
		width := td.width
		height := td.height
		m := testMap{
			buf:    make([]bool, width*height),
			width:  width,
			height: height,
		}
		for _, pot := range td.pots {
			m.buf[pot.X+pot.Y*width] = true
		}
		path := _path{
			validator: &m,
			width:     width,
			height:    height,
			hits:      make([]bool, width*height),
		}
		_, ok := path.findPath(td.begin.X, td.begin.Y, td.end.X, td.end.Y)
		if ok != td.ok {
			t.Errorf("test data %d failed", i)
		}
	}
}
