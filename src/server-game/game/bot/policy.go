package bot

import (
	"encoding/binary"
	"hash/fnv"
	"sort"
	"time"
)

const defaultAttackInterval = 800 * time.Millisecond
const minimumAttackInterval = 200 * time.Millisecond

type ActionKind int

const (
	ActionNone ActionKind = iota
	ActionConnect
	ActionLogin
	ActionLoadCharacter
	ActionStop
	ActionLearnSkill
	ActionMove
	ActionContinueMove
	ActionSyncPosition
	ActionCancel
	ActionAttack
	ActionUseSkill
	ActionChat
	ActionWhisper
)

type Action struct {
	Kind            ActionKind
	Target          int
	Skill           int
	Position        int
	Path            []Position
	Name            string
	Text            string
	SelfPosition    Position
	Dir             int
	ReadyAt         time.Time
	SentAt          time.Time
	NextStepAt      time.Time
	PathNext        int
	PositionVersion uint64
	CancelCurrent   bool
}

type Policy interface {
	Decide(time.Time, WorldSnapshot, ExecutorSnapshot) Action
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

func (p *rulePolicy) Decide(now time.Time, world WorldSnapshot, execution ExecutorSnapshot) Action {
	switch world.Phase {
	case PhaseDisconnected:
		if execution.CurrentAction.Kind == ActionConnect {
			return Action{}
		}
		return Action{Kind: ActionConnect}
	case PhaseConnected:
		if execution.CurrentAction.Kind == ActionLogin {
			return Action{}
		}
		return Action{Kind: ActionLogin}
	case PhaseLoggedIn:
		if execution.CurrentAction.Kind == ActionLoadCharacter {
			return Action{}
		}
		return Action{Kind: ActionLoadCharacter}
	case PhaseFailed:
		if execution.CurrentAction.Kind == ActionStop {
			return Action{}
		}
		return Action{Kind: ActionStop, Text: world.Failure}
	case PhasePlaying:
	default:
		return Action{}
	}
	if execution.PositionVersion != world.PositionVersion {
		return Action{
			Kind:            ActionSyncPosition,
			SelfPosition:    world.Self.position(),
			Dir:             world.Self.Dir,
			PositionVersion: world.PositionVersion,
			CancelCurrent:   execution.Move.Active,
		}
	}
	if execution.CurrentAction.Kind == ActionLearnSkill {
		if !world.learnedSkill(execution.CurrentAction.Skill) {
			return Action{}
		}
		return Action{Kind: ActionCancel}
	}
	if !world.Self.Alive {
		if execution.CurrentAction.Kind != ActionNone || execution.Move.Active {
			return Action{Kind: ActionCancel}
		}
		return Action{}
	}
	if execution.Move.Active {
		if !now.Before(execution.Move.NextStepAt) {
			return continueMoveAction(execution)
		}
		return Action{}
	}
	if (execution.CurrentAction.Kind == ActionAttack ||
		execution.CurrentAction.Kind == ActionUseSkill) &&
		now.Before(execution.ReadyAt) {
		return Action{}
	}

	snapshot := world
	snapshot.Self.X = execution.Position.X
	snapshot.Self.Y = execution.Position.Y
	snapshot.Self.Dir = execution.Dir
	combatSkill, hasSkill := bestCombatSkill(snapshot)
	targets := p.visibleTargets(snapshot)
	for _, target := range targets {
		attackRange := 1
		if hasSkill {
			attackRange = combatSkill.Distance
		}
		if pathDistance(snapshot.Self.position(), target.position()) <= attackRange {
			dir := calcDir(snapshot.Self.position(), target.position())
			if hasSkill {
				return Action{
					Kind:    ActionUseSkill,
					Target:  target.Index,
					Skill:   combatSkill.Index,
					Dir:     dir,
					ReadyAt: now.Add(skillInterval(combatSkill, world)),
				}
			}
			return Action{
				Kind:    ActionAttack,
				Target:  target.Index,
				Dir:     dir,
				ReadyAt: now.Add(speedInterval(world.AttackSpeed)),
			}
		}
		if path := p.pathToActor(snapshot, target, attackRange); len(path) > 0 {
			return moveAction(now, execution, world.PositionVersion, limitPath(path))
		}
	}
	if len(snapshot.LearnSkills) > 0 {
		learn := snapshot.LearnSkills[0]
		return Action{Kind: ActionLearnSkill, Skill: learn.Index, Position: learn.Position}
	}
	if path := p.pathToSpawnArea(snapshot); len(path) > 0 {
		return moveAction(now, execution, world.PositionVersion, limitPath(path))
	}
	if execution.CurrentAction.Kind != ActionNone {
		return Action{Kind: ActionCancel}
	}
	return Action{}
}

func moveAction(now time.Time, execution ExecutorSnapshot, positionVersion uint64, path []Position) Action {
	return Action{
		Kind:            ActionMove,
		Path:            path,
		SelfPosition:    execution.Position,
		Dir:             execution.Dir,
		SentAt:          now,
		NextStepAt:      now.Add(stepDuration(execution.Position, path[0])),
		PositionVersion: positionVersion,
	}
}

func continueMoveAction(execution ExecutorSnapshot) Action {
	next := execution.Move.Path[execution.Move.PathNext]
	pathNext := execution.Move.PathNext + 1
	action := Action{
		Kind:         ActionContinueMove,
		SelfPosition: next,
		Dir:          calcDir(execution.Position, next),
		PathNext:     pathNext,
	}
	if pathNext < len(execution.Move.Path) {
		action.NextStepAt = execution.Move.NextStepAt.Add(
			stepDuration(next, execution.Move.Path[pathNext]),
		)
	}
	return action
}

func speedInterval(speed int) time.Duration {
	interval := defaultAttackInterval -
		time.Duration(float64(speed)*5.33*float64(time.Millisecond))
	if interval < minimumAttackInterval {
		return minimumAttackInterval
	}
	return interval
}

func skillInterval(combatSkill CombatSkill, world WorldSnapshot) time.Duration {
	var interval time.Duration
	switch combatSkill.Type {
	case skillTypePhysical:
		interval = speedInterval(world.AttackSpeed)
	case skillTypeMagic:
		interval = speedInterval(world.MagicSpeed)
	default:
		return defaultAttackInterval
	}
	delay := time.Duration(combatSkill.Delay) * time.Millisecond
	if delay > interval {
		return delay
	}
	return interval
}

func stepDuration(from, to Position) time.Duration {
	if from.X != to.X && from.Y != to.Y {
		return 520 * time.Millisecond
	}
	return 400 * time.Millisecond
}

func (p *rulePolicy) visibleTargets(snapshot WorldSnapshot) []Actor {
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

func bestCombatSkill(snapshot WorldSnapshot) (CombatSkill, bool) {
	var best CombatSkill
	found := false
	for _, s := range snapshot.Skills {
		if snapshot.MP < s.MP || snapshot.AG < s.AG {
			continue
		}
		if !found || s.Damage > best.Damage || s.Damage == best.Damage && s.Index < best.Index {
			best = s
			found = true
		}
	}
	return best, found
}

func (p *rulePolicy) pathToActor(snapshot WorldSnapshot, target Actor, attackRange int) []Position {
	t := p.resources.terrain(snapshot.Self.MapNumber)
	if t == nil {
		return nil
	}
	if path := p.pathToActorWithBlockers(t, snapshot, target, attackRange, snapshot.blockers(target.Index)); len(path) > 0 {
		return path
	}
	return p.pathToActorWithBlockers(t, snapshot, target, attackRange, nil)
}

func (p *rulePolicy) pathToActorWithBlockers(t *terrain, snapshot WorldSnapshot, target Actor, attackRange int, blockers map[Position]struct{}) []Position {
	var best []Position
	bestScore := 0
	for y := target.Y - attackRange; y <= target.Y+attackRange; y++ {
		for x := target.X - attackRange; x <= target.X+attackRange; x++ {
			goal := Position{X: x, Y: y}
			if goal == target.position() ||
				pathDistance(goal, target.position()) > attackRange ||
				!t.walkable(goal) {
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
	}
	return best
}

func (p *rulePolicy) pathToSpawnArea(snapshot WorldSnapshot) []Position {
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
