package bot

import "testing"

func TestFindPathUsesStaticTerrainAndOptionalBlockers(t *testing.T) {
	terrain := newTestTerrain(5, 5)
	terrain.attrs[2+2*terrain.width] = 4
	start := Position{X: 1, Y: 2}
	goal := Position{X: 3, Y: 2}

	path, ok := findPath(terrain, start, goal, map[Position]struct{}{{X: 2, Y: 1}: {}})
	if !ok {
		t.Fatal("findPath() ok = false, want true")
	}
	for _, pos := range path {
		if pos == (Position{X: 2, Y: 2}) || pos == (Position{X: 2, Y: 1}) {
			t.Fatalf("findPath() crossed blocked position %+v: %#v", pos, path)
		}
	}
}

func newTestTerrain(width, height int) *terrain {
	return &terrain{
		width:  width,
		height: height,
		attrs:  make([]byte, width*height),
	}
}
