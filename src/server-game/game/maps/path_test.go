package maps

import (
	"fmt"
	"testing"
)

/*
go test -v -timeout 30s -bench . -benchmem ~/github.com/xujintao/balgass/src/server-game/game/maps/path_test.go ~/github.com/xujintao/balgass/src/server-game/game/maps/path.go
=== RUN   TestFindPath
(1,0)(0,1)(0,2)(1,3)(2,3)(3,2)(4,1)(5,1)(6,1)
(2,3)(3,3)(4,3)(5,4)(5,5)(6,5)(7,6)(7,7)(6,8)(5,7)
--- PASS: TestFindPath (0.00s)
=== RUN   TestFindPathBFS
(0,2)(1,3)(2,3)(3,2)(4,1)(5,0)(6,1)
(2,1)(3,2)(4,3)(5,4)(4,5)(3,6)(4,7)(5,7)
--- PASS: TestFindPathBFS (0.00s)
=== RUN   TestFindPathAStar
(0,2)(1,3)(2,3)(3,4)(4,3)(5,2)(6,1)
(2,3)(3,3)(4,3)(5,4)(4,5)(3,6)(4,7)(5,7)
--- PASS: TestFindPathAStar (0.00s)
goos: linux
goarch: amd64
cpu: 13th Gen Intel(R) Core(TM) i5-13500H
BenchmarkFindPath
BenchmarkFindPath-16               92775             12452 ns/op            3121 B/op         59 allocs/op
BenchmarkFindPathBFS
BenchmarkFindPathBFS-16            16488             64960 ns/op           32271 B/op        396 allocs/op
BenchmarkFindPathAStar
BenchmarkFindPathAStar-16          35486             30248 ns/op           15458 B/op        318 allocs/op
PASS
ok      command-line-arguments  4.515s
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
	// {true, 8, 8, Pot{1, 2}, Pot{5, 7}, []Pot{{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 6}, {5, 6}, {6, 6}, {6, 7}}},
	{true, 256, 256, Pot{1, 2}, Pot{5, 7}, []Pot{{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 6}, {5, 6}, {6, 6}, {6, 7}}},
}

func doFindPathFunc(td *testData, findPathFunc func(*_path, int, int, int, int) (Path, bool)) (Path, bool) {
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
	return findPathFunc(&path, td.begin.X, td.begin.Y, td.end.X, td.end.Y)
}

/*
go test -v -timeout 30s ~/github.com/xujintao/balgass/src/server-game/game/maps/path_test.go ~/github.com/xujintao/balgass/src/server-game/game/maps/path.go
=== RUN   TestFindPath
(1,0)(0,1)(0,2)(1,3)(2,3)(3,2)(4,1)(5,1)(6,1)
(2,3)(3,3)(4,3)(5,4)(5,5)(6,5)(7,6)(7,7)(6,8)(5,7)
--- PASS: TestFindPath (0.00s)
=== RUN   TestFindPathBFS
(0,2)(1,3)(2,3)(3,2)(4,1)(5,0)(6,1)
(2,1)(3,2)(4,3)(5,4)(4,5)(3,6)(4,7)(5,7)
--- PASS: TestFindPathBFS (0.00s)
=== RUN   TestFindPathAStar
(0,2)(1,3)(2,3)(3,4)(4,3)(5,2)(6,1)
(2,3)(3,3)(4,3)(5,4)(4,5)(3,6)(4,7)(5,7)
--- PASS: TestFindPathAStar (0.00s)
PASS
ok      command-line-arguments  0.002s
*/
func runTestFindPath(t *testing.T, findPathFunc func(*_path, int, int, int, int) (Path, bool)) {
	for _, td := range testDatas {
		p, ok := doFindPathFunc(td, findPathFunc)
		if p != nil {
			for _, pot := range p {
				fmt.Printf("(%d,%d)", pot.X, pot.Y)
			}
			fmt.Println()
		}
		if ok != td.ok {
			t.Errorf("test data failed: %v", td)
		}
	}
}

func TestFindPath(t *testing.T) {
	runTestFindPath(t, (*_path).findPath)
}

func TestFindPathBFS(t *testing.T) {
	runTestFindPath(t, (*_path).findPathBFS)
}

func TestFindPathAStar(t *testing.T) {
	runTestFindPath(t, (*_path).findPathAStar)
}

/*
go test -v -timeout 30s -bench . -benchmem -run=^$ ~/github.com/xujintao/balgass/src/server-game/game/maps/path_test.go ~/github.com/xujintao/balgass/src/server-game/game/maps/path.go
goos: linux
goarch: amd64
cpu: 13th Gen Intel(R) Core(TM) i5-13500H
BenchmarkFindPath
BenchmarkFindPath-16               82720             12196 ns/op            3121 B/op         59 allocs/op
BenchmarkFindPathBFS
BenchmarkFindPathBFS-16            16951             65484 ns/op           32274 B/op        396 allocs/op
BenchmarkFindPathAStar
BenchmarkFindPathAStar-16          35779             29971 ns/op           15462 B/op        318 allocs/op
PASS
ok      command-line-arguments  4.404s
*/
func runBenchmarkFindPath(b *testing.B, findPathFunc func(*_path, int, int, int, int) (Path, bool)) {
	for i := 0; i < b.N; i++ {
		for _, td := range testDatas {
			doFindPathFunc(td, findPathFunc)
		}
	}
}

func BenchmarkFindPath(b *testing.B) {
	runBenchmarkFindPath(b, (*_path).findPath)
}

func BenchmarkFindPathBFS(b *testing.B) {
	runBenchmarkFindPath(b, (*_path).findPathBFS)
}

func BenchmarkFindPathAStar(b *testing.B) {
	runBenchmarkFindPath(b, (*_path).findPathAStar)
}
