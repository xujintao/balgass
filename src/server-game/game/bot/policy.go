package bot

import (
	"sort"
	"time"
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
	resources  *resources
	nextAttack time.Time
}

func newRulePolicy(resources *resources) Policy {
	return &rulePolicy{
		resources: resources,
	}
}

func (p *rulePolicy) Decide(snapshot Snapshot) Action {
	if !snapshot.Self.Alive || snapshot.Moving {
		return Action{}
	}
	targets := visibleTargets(snapshot)
	for _, target := range targets {
		if pathDistance(snapshot.Self.position(), target.position()) <= 1 {
			if snapshot.Now.Before(p.nextAttack) {
				return Action{}
			}
			p.nextAttack = snapshot.Now.Add(p.attackDelay())
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

func visibleTargets(snapshot Snapshot) []Actor {
	targets := make([]Actor, 0, len(snapshot.Objects))
	for _, actor := range snapshot.Objects {
		if actor.Alive && actor.Attackable && actor.MapNumber == snapshot.Self.MapNumber {
			targets = append(targets, actor)
		}
	}
	sort.Slice(targets, func(i, j int) bool {
		return pathDistance(snapshot.Self.position(), targets[i].position()) <
			pathDistance(snapshot.Self.position(), targets[j].position())
	})
	return targets
}

func (p *rulePolicy) pathToActor(snapshot Snapshot, target Actor) []Position {
	t := p.resources.terrain(snapshot.Self.MapNumber)
	if t == nil {
		return nil
	}
	var best []Position
	for _, direction := range directions {
		goal := Position{X: target.X + direction.X, Y: target.Y + direction.Y}
		if !t.walkable(goal) {
			continue
		}
		path, ok := findPath(t, snapshot.Self.position(), goal, snapshot.blockers(target.Index))
		if !ok {
			path, ok = findPath(t, snapshot.Self.position(), goal, nil)
		}
		if ok && len(path) > 0 && (best == nil || len(path) < len(best)) {
			best = path
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
		distance int
	}
	var destinations []destination
	seen := make(map[Position]struct{})
	add := func(pos Position) {
		if pos == snapshot.Self.position() || !t.walkable(pos) {
			return
		}
		if _, ok := seen[pos]; ok {
			return
		}
		seen[pos] = struct{}{}
		destinations = append(destinations, destination{
			Position: pos,
			distance: pathDistance(snapshot.Self.position(), pos),
		})
	}
	for _, area := range p.resources.spawnAreas[snapshot.Self.MapNumber] {
		if !p.resources.attackable(area.class) {
			continue
		}
		add(area.nearest(snapshot.Self.position()))
		add(area.center())
		if area.contains(snapshot.Self.position()) {
			add(Position{X: area.min.X, Y: area.min.Y})
			add(Position{X: area.max.X, Y: area.max.Y})
		}
	}
	sort.Slice(destinations, func(i, j int) bool {
		return destinations[i].distance < destinations[j].distance
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

func (p *rulePolicy) attackDelay() time.Duration {
	return 800 * time.Millisecond
}

func limitPath(path []Position) []Position {
	if len(path) > 15 {
		path = path[:15]
	}
	return append([]Position(nil), path...)
}
