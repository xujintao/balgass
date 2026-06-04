package bot

import "container/heap"

var directions = [8]Position{
	{X: -1, Y: -1},
	{X: 0, Y: -1},
	{X: 1, Y: -1},
	{X: 1, Y: 0},
	{X: 1, Y: 1},
	{X: 0, Y: 1},
	{X: -1, Y: 1},
	{X: -1, Y: 0},
}

type pathNode struct {
	Position
	cost     int
	priority int
	index    int
}

type pathQueue []*pathNode

func (q pathQueue) Len() int           { return len(q) }
func (q pathQueue) Less(i, j int) bool { return q[i].priority < q[j].priority }
func (q pathQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}
func (q *pathQueue) Push(v any) {
	n := v.(*pathNode)
	n.index = len(*q)
	*q = append(*q, n)
}
func (q *pathQueue) Pop() any {
	old := *q
	n := old[len(old)-1]
	n.index = -1
	*q = old[:len(old)-1]
	return n
}

func findPath(t *terrain, start, goal Position, blockers map[Position]struct{}) ([]Position, bool) {
	if t == nil || !t.walkable(start) || !t.walkable(goal) {
		return nil, false
	}
	if start == goal {
		return nil, true
	}
	if _, ok := blockers[goal]; ok {
		return nil, false
	}

	frontier := &pathQueue{}
	heap.Init(frontier)
	heap.Push(frontier, &pathNode{Position: start, priority: pathDistance(start, goal)})
	costs := map[Position]int{start: 0}
	parents := make(map[Position]Position)

	for frontier.Len() > 0 {
		current := heap.Pop(frontier).(*pathNode)
		if current.Position == goal {
			return buildPath(parents, start, goal), true
		}
		if cost, ok := costs[current.Position]; ok && current.cost > cost {
			continue
		}
		for _, direction := range directions {
			next := Position{X: current.X + direction.X, Y: current.Y + direction.Y}
			if !t.walkable(next) {
				continue
			}
			if _, ok := blockers[next]; ok {
				continue
			}
			cost := current.cost + 1
			if old, ok := costs[next]; ok && old <= cost {
				continue
			}
			costs[next] = cost
			parents[next] = current.Position
			heap.Push(frontier, &pathNode{
				Position: next,
				cost:     cost,
				priority: cost + pathDistance(next, goal),
			})
		}
	}
	return nil, false
}

func buildPath(parents map[Position]Position, start, goal Position) []Position {
	path := []Position{goal}
	for current := goal; current != start; {
		current = parents[current]
		if current != start {
			path = append(path, current)
		}
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func pathDistance(a, b Position) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	if dx > dy {
		return dx
	}
	return dy
}

func calcDir(from, to Position) int {
	dx := to.X - from.X
	dy := to.Y - from.Y
	switch {
	case dx < 0 && dy < 0:
		return 0
	case dx == 0 && dy < 0:
		return 1
	case dx > 0 && dy < 0:
		return 2
	case dx > 0 && dy == 0:
		return 3
	case dx > 0 && dy > 0:
		return 4
	case dx == 0 && dy > 0:
		return 5
	case dx < 0 && dy > 0:
		return 6
	case dx < 0 && dy == 0:
		return 7
	default:
		return 0
	}
}
