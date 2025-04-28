package maps

import (
	"fmt"
	"testing"
)

/*
go test -v -timeout 30s -bench . -benchmem ~/github.com/xujintao/balgass/src/server_game/game/maps/path_test.go ~/github.com/xujintao/balgass/src/server_game/game/maps/path.go
=== RUN   TestFindPath
(1,0)(0,1)(0,2)(1,3)(2,3)(3,2)(4,1)(5,1)(6,1)
(2,3)(3,3)(4,3)(5,4)(5,5)(6,5)(7,6)(7,7)(6,8)(5,7)
--- PASS: TestFindPath (0.00s)
=== RUN   TestFindPathBFS
(0,2)(1,3)(2,3)(3,2)(4,1)(5,0)(6,1)
(2,1)(3,2)(4,3)(5,4)(4,5)(3,6)(4,7)(5,7)
(2,1)(3,2)(4,3)(5,4)(4,5)(3,6)(4,7)(5,7)
--- PASS: TestFindPathBFS (0.00s)
goos: linux
goarch: amd64
cpu: 13th Gen Intel(R) Core(TM) i5-13500H
BenchmarkFindPath
BenchmarkFindPath-16               49614             22418 ns/op            3566 B/op         67 allocs/op
BenchmarkFindPathBFS
BenchmarkFindPathBFS-16            11943             98247 ns/op           42590 B/op        515 allocs/op
PASS
ok      command-line-arguments  3.403s
*/
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
	{true, 8, 8, Pot{1, 2}, Pot{5, 7}, []Pot{{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 6}, {5, 6}, {6, 6}, {6, 7}}},
	{true, 256, 256, Pot{1, 2}, Pot{5, 7}, []Pot{{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 6}, {5, 6}, {6, 6}, {6, 7}}},
}

/*
go test -v -timeout 30s -bench . -benchmem ~/github.com/xujintao/balgass/src/server_game/game/maps/path_test.go ~/github.com/xujintao/balgass/src/server_game/game/maps/path.go
=== RUN   TestFindPath
(1,0)(0,1)(0,2)(1,3)(2,3)(3,2)(4,1)(5,1)(6,1)
(2,3)(3,3)(4,3)(5,4)(5,5)(6,5)(7,6)(7,7)(6,8)(5,7)
--- PASS: TestFindPath (0.00s)
=== RUN   TestFindPathBFS
(0,2)(1,3)(2,3)(3,2)(4,1)(5,0)(6,1)
(2,1)(3,2)(4,3)(5,4)(4,5)(3,6)(4,7)(5,7)
(2,1)(3,2)(4,3)(5,4)(4,5)(3,6)(4,7)(5,7)
--- PASS: TestFindPathBFS (0.00s)
PASS
ok      command-line-arguments  0.002s
*/
func TestFindPath(t *testing.T) {
	var testDatas = []*testData{
		{true, 256, 256, Pot{1, 1}, Pot{6, 1}, []Pot{{2, 0}, {2, 1}, {2, 2}, {1, 2}}},
		{false, 256, 256, Pot{1, 1}, Pot{6, 1}, []Pot{{2, 0}, {2, 1}, {2, 2}, {1, 2}, {0, 2}}},
		{false, 8, 8, Pot{1, 2}, Pot{5, 7}, []Pot{{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 6}, {5, 6}, {6, 6}, {6, 7}}},
		{true, 256, 256, Pot{1, 2}, Pot{5, 7}, []Pot{{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 6}, {5, 6}, {6, 6}, {6, 7}}},
	}
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

func TestFindPathBFS(t *testing.T) {
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
		p, ok := path.findPathBFS(td.begin.X, td.begin.Y, td.end.X, td.end.Y)
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

func TestFindPathAStar(t *testing.T) {
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
		p, ok := path.findPathAStar(td.begin.X, td.begin.Y, td.end.X, td.end.Y)
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
go test -v -timeout 30s -bench . -benchmem -run=^$ ~/github.com/xujintao/balgass/src/server_game/game/maps/path_test.go ~/github.com/xujintao/balgass/src/server_game/game/maps/path.go
goos: linux
goarch: amd64
cpu: 13th Gen Intel(R) Core(TM) i5-13500H
BenchmarkFindPath
BenchmarkFindPath-16               49992             21972 ns/op            3565 B/op         67 allocs/op
BenchmarkFindPathBFS
BenchmarkFindPathBFS-16            12277             95807 ns/op           42584 B/op        515 allocs/op
PASS
ok      command-line-arguments  3.511s
*/
func BenchmarkFindPath(b *testing.B) {
	var testDatas = []*testData{
		{true, 256, 256, Pot{1, 1}, Pot{6, 1}, []Pot{{2, 0}, {2, 1}, {2, 2}, {1, 2}}},
		{false, 256, 256, Pot{1, 1}, Pot{6, 1}, []Pot{{2, 0}, {2, 1}, {2, 2}, {1, 2}, {0, 2}}},
		{false, 8, 8, Pot{1, 2}, Pot{5, 7}, []Pot{{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 6}, {5, 6}, {6, 6}, {6, 7}}},
		{true, 256, 256, Pot{1, 2}, Pot{5, 7}, []Pot{{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 6}, {5, 6}, {6, 6}, {6, 7}}},
	}
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

func BenchmarkFindPathBFS(b *testing.B) {
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
			path.findPathBFS(td.begin.X, td.begin.Y, td.end.X, td.end.Y)
		}
	}
}

func BenchmarkFindPathAStar(b *testing.B) {
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
			path.findPathAStar(td.begin.X, td.begin.Y, td.end.X, td.end.Y)
		}
	}
}
