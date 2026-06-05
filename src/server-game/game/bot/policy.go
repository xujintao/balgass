package bot

import (
	"encoding/binary"
	"hash/fnv"
	"sort"
)

type ActionKind int

const (
	ActionNone ActionKind = iota
	ActionMove
	ActionAttack
	ActionChat
	ActionWhisper
)

type Action struct {
	Kind   ActionKind
	Target int
	Path   []Position
	Name   string
	Text   string
}

type Policy interface {
	Decide(Snapshot) Action
}

type rulePolicy struct {
	resources *resources
	seed      string
}

func newRulePolicy(resources *resources, seed string) Policy {
	return &rulePolicy{
		resources: resources,
		seed:      seed,
	}
}

func (p *rulePolicy) Decide(snapshot Snapshot) Action {
	if !snapshot.Self.Alive || snapshot.Moving {
		return Action{}
	}
	targets := p.visibleTargets(snapshot)
	for _, target := range targets {
		if pathDistance(snapshot.Self.position(), target.position()) <= 1 {
			return Action{Kind: ActionAttack, Target: target.Index}
		}
		if path := p.pathToActor(snapshot, target); len(path) > 0 {
			return Action{Kind: ActionMove, Path: limitPath(path)}
		}
	}
	if path := p.pathToSpawnArea(snapshot); len(path) > 0 {
		return Action{Kind: ActionMove, Path: limitPath(path)}
	}
	return Action{}
}

func (p *rulePolicy) visibleTargets(snapshot Snapshot) []Actor {
	targets := make([]Actor, 0, len(snapshot.Objects))
	for _, actor := range snapshot.Objects {
		if actor.Alive && actor.Attackable && actor.MapNumber == snapshot.Self.MapNumber {
			targets = append(targets, actor)
		}
	}
	sort.Slice(targets, func(i, j int) bool {
		return p.targetScore(snapshot.Self, targets[i]) < p.targetScore(snapshot.Self, targets[j])
	})
	return targets
}

func (p *rulePolicy) pathToActor(snapshot Snapshot, target Actor) []Position {
	t := p.resources.terrain(snapshot.Self.MapNumber)
	if t == nil {
		return nil
	}
	if path := p.pathToActorWithBlockers(t, snapshot, target, snapshot.blockers(target.Index)); len(path) > 0 {
		return path
	}
	return p.pathToActorWithBlockers(t, snapshot, target, nil)
}

func (p *rulePolicy) pathToActorWithBlockers(t *terrain, snapshot Snapshot, target Actor, blockers map[Position]struct{}) []Position {
	var best []Position
	bestScore := 0
	for _, direction := range directions {
		goal := Position{X: target.X + direction.X, Y: target.Y + direction.Y}
		if !t.walkable(goal) {
			continue
		}
		path, ok := findPath(t, snapshot.Self.position(), goal, blockers)
		if !ok || len(path) == 0 {
			continue
		}
		score := len(path)*100 + p.positionJitter(11, target.MapNumber, target.Index, goal)
		if best == nil || score < bestScore {
			best = path
			bestScore = score
		}
	}
	return best
}

func (p *rulePolicy) pathToSpawnArea(snapshot Snapshot) []Position {
	t := p.resources.terrain(snapshot.Self.MapNumber)
	if t == nil {
		return nil
	}
	type destination struct {
		Position
		score int
	}
	var destinations []destination
	seen := make(map[Position]struct{})
	add := func(pos Position, area spawnArea, salt int) {
		if pos == snapshot.Self.position() || !t.walkable(pos) {
			return
		}
		if _, ok := seen[pos]; ok {
			return
		}
		seen[pos] = struct{}{}
		destinations = append(destinations, destination{
			Position: pos,
			score:    pathDistance(snapshot.Self.position(), pos)*100 + p.positionJitter(23+uint32(salt), snapshot.Self.MapNumber, area.class, pos),
		})
	}
	for _, area := range p.resources.spawnAreas[snapshot.Self.MapNumber] {
		if !p.resources.attackable(area.class) {
			continue
		}
		add(area.nearest(snapshot.Self.position()), area, 0)
		add(area.center(), area, 1)
		for i := 0; i < 4; i++ {
			add(p.spawnPoint(area, i), area, i+2)
		}
		if area.contains(snapshot.Self.position()) {
			add(Position{X: area.min.X, Y: area.min.Y}, area, 6)
			add(Position{X: area.max.X, Y: area.max.Y}, area, 7)
		}
	}
	sort.Slice(destinations, func(i, j int) bool {
		return destinations[i].score < destinations[j].score
	})
	blockers := snapshot.blockers(-1)
	for _, destination := range destinations {
		path, ok := findPath(t, snapshot.Self.position(), destination.Position, blockers)
		if !ok {
			path, ok = findPath(t, snapshot.Self.position(), destination.Position, nil)
		}
		if ok && len(path) > 0 {
			return path
		}
	}
	return nil
}

func (p *rulePolicy) targetScore(self, target Actor) int {
	return pathDistance(self.position(), target.position())*100 +
		p.positionJitter(5, target.MapNumber, target.Index, target.position())
}

func (p *rulePolicy) positionJitter(kind uint32, mapNumber, id int, pos Position) int {
	return int(stableHash(p.seed, kind, mapNumber, id, pos.X, pos.Y) % 32)
}

func (p *rulePolicy) spawnPoint(area spawnArea, salt int) Position {
	width := area.max.X - area.min.X + 1
	height := area.max.Y - area.min.Y + 1
	return Position{
		X: area.min.X + int(stableHash(p.seed, 31+uint32(salt), area.class, area.min.X, area.max.X)%uint32(width)),
		Y: area.min.Y + int(stableHash(p.seed, 47+uint32(salt), area.class, area.min.Y, area.max.Y)%uint32(height)),
	}
}

func stableHash(seed string, kind uint32, values ...int) uint32 {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(seed))

	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], kind)
	_, _ = hash.Write(buf[:])

	for _, value := range values {
		binary.LittleEndian.PutUint32(buf[:], uint32(value))
		_, _ = hash.Write(buf[:])
	}

	sum := hash.Sum32()
	if sum == 0 {
		return 1
	}
	return sum
}

func limitPath(path []Position) []Position {
	if len(path) > 15 {
		path = path[:15]
	}
	return append([]Position(nil), path...)
}
