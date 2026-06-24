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
	ActionContinueAttack
	ActionUseSkill
	ActionContinueUseSkill
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
	traceWorld := world
	finish := func(reason string, action Action, extra map[string]interface{}) Action {
		tracePolicyDecision(p.seed, traceWorld, execution, reason, action, extra)
		return action
	}
	switch world.Phase {
	case PhaseDisconnected:
		if execution.CurrentAction.Kind == ActionConnect {
			return finish("connect_pending", Action{}, nil)
		}
		return finish("connect", Action{Kind: ActionConnect}, nil)
	case PhaseConnected:
		if execution.CurrentAction.Kind == ActionLogin {
			return finish("login_pending", Action{}, nil)
		}
		return finish("login", Action{Kind: ActionLogin}, nil)
	case PhaseLoggedIn:
		if execution.CurrentAction.Kind == ActionLoadCharacter {
			return finish("load_character_pending", Action{}, nil)
		}
		return finish("load_character", Action{Kind: ActionLoadCharacter}, nil)
	case PhaseFailed:
		if execution.CurrentAction.Kind == ActionStop {
			return finish("stop_pending", Action{}, nil)
		}
		return finish("stop", Action{Kind: ActionStop, Text: world.Failure}, map[string]interface{}{
			"failure": world.Failure,
		})
	case PhasePlaying:
	default:
		return finish("unsupported_phase", Action{}, nil)
	}
	if execution.PositionVersion != world.PositionVersion {
		return finish("sync_position", Action{
			Kind:            ActionSyncPosition,
			SelfPosition:    world.Self.position(),
			Dir:             world.Self.Dir,
			PositionVersion: world.PositionVersion,
			CancelCurrent:   execution.Move.Active,
		}, map[string]interface{}{
			"world_position_version":     world.PositionVersion,
			"execution_position_version": execution.PositionVersion,
		})
	}
	if execution.CurrentAction.Kind == ActionLearnSkill {
		if !world.learnedSkill(execution.CurrentAction.Skill) {
			return finish("learn_skill_pending", Action{}, map[string]interface{}{
				"skill": execution.CurrentAction.Skill,
			})
		}
		return finish("learn_skill_done_cancel", Action{Kind: ActionCancel}, map[string]interface{}{
			"skill": execution.CurrentAction.Skill,
		})
	}
	if !world.Self.Alive {
		if execution.CurrentAction.Kind != ActionNone || execution.Move.Active {
			return finish("dead_cancel", Action{Kind: ActionCancel}, nil)
		}
		return finish("dead_idle", Action{}, nil)
	}

	snapshot := world
	snapshot.Self.X = execution.Position.X
	snapshot.Self.Y = execution.Position.Y
	snapshot.Self.Dir = execution.Dir
	traceWorld = snapshot
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
				return finish("use_skill", useSkillAction(now, world, combatSkill, target, dir), map[string]interface{}{
					"target":       traceActor(target),
					"targets":      len(targets),
					"attack_range": attackRange,
					"skill":        combatSkill.Index,
				})
			}
			return finish("attack", attackAction(now, world, target, dir), map[string]interface{}{
				"target":       traceActor(target),
				"targets":      len(targets),
				"attack_range": attackRange,
			})
		}
		if path := p.pathToActor(snapshot, target, attackRange); len(path) > 0 {
			limitedPath := limitPath(path)
			return finish("move_to_target", moveAction(now, execution, world.PositionVersion, limitedPath), map[string]interface{}{
				"target":       traceActor(target),
				"targets":      len(targets),
				"attack_range": attackRange,
				"path_len":     len(path),
				"path_end":     path[len(path)-1],
			})
		}
	}
	if len(snapshot.LearnSkills) > 0 {
		learn := snapshot.LearnSkills[0]
		return finish("learn_skill", Action{Kind: ActionLearnSkill, Skill: learn.Index, Position: learn.Position}, map[string]interface{}{
			"learn_skills": len(snapshot.LearnSkills),
			"skill":        learn.Index,
			"position":     learn.Position,
		})
	}
	if path := p.pathToSpawnArea(snapshot); len(path) > 0 {
		limitedPath := limitPath(path)
		return finish("move_to_spawn_area", moveAction(now, execution, world.PositionVersion, limitedPath), map[string]interface{}{
			"path_len": len(path),
			"path_end": path[len(path)-1],
		})
	}
	if execution.CurrentAction.Kind != ActionNone {
		return finish("cancel_current", Action{Kind: ActionCancel}, nil)
	}
	return finish("none", Action{}, nil)
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

func attackAction(now time.Time, world WorldSnapshot, target Actor, dir int) Action {
	return Action{
		Kind:    ActionAttack,
		Target:  target.Index,
		Dir:     dir,
		ReadyAt: now.Add(speedInterval(world.AttackSpeed)),
	}
}

func useSkillAction(now time.Time, world WorldSnapshot, combatSkill CombatSkill, target Actor, dir int) Action {
	return Action{
		Kind:    ActionUseSkill,
		Target:  target.Index,
		Skill:   combatSkill.Index,
		Dir:     dir,
		ReadyAt: now.Add(skillInterval(combatSkill, world)),
	}
}

type continuation struct {
	kind  ActionKind
	dueAt time.Time
}

func continueAction(now time.Time, world WorldSnapshot, execution ExecutorSnapshot) (Action, bool) {
	currentKind := execution.CurrentAction.Kind
	if execution.Move.Active {
		currentKind = ActionMove
	}
	var cont continuation
	switch currentKind {
	case ActionMove, ActionContinueMove:
		cont = continuation{
			kind:  ActionContinueMove,
			dueAt: execution.Move.NextStepAt,
		}
	case ActionAttack, ActionContinueAttack:
		cont = continuation{
			kind:  ActionContinueAttack,
			dueAt: execution.ReadyAt,
		}
	case ActionUseSkill, ActionContinueUseSkill:
		cont = continuation{
			kind:  ActionContinueUseSkill,
			dueAt: execution.ReadyAt,
		}
	default:
		return Action{}, false
	}
	if now.Before(cont.dueAt) {
		return Action{}, true
	}
	switch cont.kind {
	case ActionContinueMove:
		if !execution.Move.Active ||
			execution.Move.PathNext < 0 ||
			execution.Move.PathNext >= len(execution.Move.Path) {
			return Action{Kind: ActionCancel}, true
		}
		return continueMoveAction(execution), true
	case ActionContinueAttack:
		target, ok := world.object(execution.CurrentAction.Target)
		if !ok || !validAttackTarget(world, execution, target, 1) {
			return Action{Kind: ActionCancel}, true
		}
		return continueAttackAction(now, world, execution), true
	case ActionContinueUseSkill:
		combatSkill, ok := world.combatSkill(execution.CurrentAction.Skill)
		if !ok || world.MP < combatSkill.MP || world.AG < combatSkill.AG {
			return Action{Kind: ActionCancel}, true
		}
		target, ok := world.object(execution.CurrentAction.Target)
		if !ok || !validAttackTarget(world, execution, target, combatSkill.Distance) {
			return Action{Kind: ActionCancel}, true
		}
		return continueUseSkillAction(now, world, execution), true
	default:
		return Action{Kind: ActionCancel}, true
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

func continueAttackAction(now time.Time, world WorldSnapshot, execution ExecutorSnapshot) Action {
	return Action{
		Kind:    ActionContinueAttack,
		Target:  execution.CurrentAction.Target,
		Dir:     execution.CurrentAction.Dir,
		ReadyAt: now.Add(speedInterval(world.AttackSpeed)),
	}
}

func continueUseSkillAction(now time.Time, world WorldSnapshot, execution ExecutorSnapshot) Action {
	combatSkill, _ := world.combatSkill(execution.CurrentAction.Skill)
	return Action{
		Kind:    ActionContinueUseSkill,
		Target:  execution.CurrentAction.Target,
		Skill:   execution.CurrentAction.Skill,
		Dir:     execution.CurrentAction.Dir,
		ReadyAt: now.Add(skillInterval(combatSkill, world)),
	}
}

func validAttackTarget(world WorldSnapshot, execution ExecutorSnapshot, target Actor, attackRange int) bool {
	if attackRange < 1 {
		attackRange = 1
	}
	return target.Alive &&
		target.Attackable &&
		target.MapNumber == world.Self.MapNumber &&
		pathDistance(execution.Position, target.position()) <= attackRange
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
	case skillTypeSiege, skillTypePhysical:
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
