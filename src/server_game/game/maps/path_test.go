package maps

import (
	"fmt"
	"testing"
)

type testData struct {
	ok     bool
	width  int
	height int
	begin  Pot
	end    Pot
	blocks []Pot
}

var testDatas = []*testData{
	{true, 256, 256, Pot{1, 1}, Pot{6, 1}, []Pot{{2, 0}, {2, 1}, {2, 2}, {1, 2}}},
	{false, 256, 256, Pot{1, 1}, Pot{6, 1}, []Pot{{2, 0}, {2, 1}, {2, 2}, {1, 2}, {0, 2}}},
	{false, 8, 8, Pot{1, 2}, Pot{5, 7}, []Pot{{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 6}, {5, 6}, {6, 6}, {6, 7}}},
	{true, 256, 256, Pot{1, 2}, Pot{5, 7}, []Pot{{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 6}, {5, 6}, {6, 6}, {6, 7}}},
}

/*
go test -timeout 30s -run ^TestFindPath$ ~/github.com/xujintao/balgass/src/server_game/game/maps/path_test.go ~/github.com/xujintao/balgass/src/server_game/game/maps/path.go
ok      command-line-arguments  0.002s
*/
func TestFindPath(t *testing.T) {
	for i, td := range testDatas {
		width := td.width
		height := td.height
		blocks := make(map[Pot]struct{})
		for _, pot := range td.blocks {
			blocks[pot] = struct{}{}
		}
		path := _path{
			validator: func(x, y int) bool {
				if x < 0 || x >= width || y < 0 || y >= height {
					return false
				}
				_, ok := blocks[Pot{x, y}]
				return !ok
			},
			hits: make(map[Pot]struct{}),
		}
		p, ok := path.findPath(td.begin.X, td.begin.Y, td.end.X, td.end.Y)
		if p != nil {
			for _, pot := range p {
				fmt.Printf("(%d,%d)", pot.X, pot.Y)
			}
			fmt.Println()
		}
		if ok != td.ok {
			t.Errorf("test data %d failed", i)
		}
	}
}

/*
go test -benchmem -run=^$ -bench ^BenchmarkFindPath$ ~/github.com/xujintao/balgass/s
rc/server_game/game/maps/path_test.go ~/github.com/xujintao/balgass/src/server_game/game/maps/path.go
goos: linux
goarch: amd64
cpu: 13th Gen Intel(R) Core(TM) i5-13500H
BenchmarkFindPath-16               47414             23979 ns/op            5229 B/op         87 allocs/op
PASS
ok      command-line-arguments  1.397s
*/
func BenchmarkFindPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, td := range testDatas {
			width := td.width
			height := td.height
			blocks := make(map[Pot]struct{})
			for _, pot := range td.blocks {
				blocks[pot] = struct{}{}
			}
			path := _path{
				validator: func(x, y int) bool {
					if x < 0 || x >= width || y < 0 || y >= height {
						return false
					}
					_, ok := blocks[Pot{x, y}]
					return !ok
				},
				hits: make(map[Pot]struct{}),
			}
			path.findPath(td.begin.X, td.begin.Y, td.end.X, td.end.Y)
		}
	}
}
