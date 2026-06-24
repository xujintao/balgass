package bot

import (
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

type policyCall struct {
	now      time.Time
	world    WorldSnapshot
	executor ExecutorSnapshot
}

type stubPolicy struct {
	actions []Action
	calls   []policyCall
}

func (p *stubPolicy) Decide(now time.Time, world WorldSnapshot, executor ExecutorSnapshot) Action {
	p.calls = append(p.calls, policyCall{now: now, world: world, executor: executor})
	call := len(p.calls) - 1
	if call >= len(p.actions) {
		return Action{}
	}
	return p.actions[call]
}

func newUnitBot(g game, policy Policy) *bot {
	b := &bot{
		account:  "account1",
		password: "password",
		name:     "char1",
		game:     g,
		world:    newWorld(&resources{}),
		policy:   policy,
		msgChan:  make(chan any, 1),
	}
	b.executor = newExecutor(b)
	b.id.Store(-1)
	return b
}

func skillBook(index int) *item.Item {
	base := &item.ItemBase{
		KindA:      item.KindASkill,
		SkillIndex: index,
	}
	base.ReqClass[class.Wizard] = 1
	return &item.Item{
		ItemBase: base,
		Code:     item.Code(15, 0),
		ID:       123,
	}
}

func TestExecutorTranslatesMoveAndAttackActions(t *testing.T) {
	g := newTestGame()
	b := newUnitBot(g, &stubPolicy{})
	b.id.Store(7)
	b.world.phase = PhasePlaying
	b.world.self = Actor{Index: 7, X: 1, Y: 1, Alive: true}
	now := time.Unix(0, 0)

	b.executor.Execute(Action{
		Kind:         ActionMove,
		Path:         []Position{{X: 2, Y: 1}, {X: 3, Y: 2}},
		SelfPosition: Position{X: 1, Y: 1},
		SentAt:       now,
		NextStepAt:   now.Add(400 * time.Millisecond),
	})

	moveAction := waitAction(t, g)
	if moveAction.id != 7 || moveAction.action != "Move" {
		t.Fatalf("move action = %#v, want player 7 Move", moveAction)
	}
	move := moveAction.msg.(*model.MsgMove)
	if move.X != 1 || move.Y != 1 || len(move.Path) != 2 ||
		move.Path[1].X != 3 || move.Path[1].Y != 2 {
		t.Fatalf("move msg = %#v, want path from (1,1) to (3,2)", move)
	}
	if len(move.Dirs) != 2 || move.Dirs[0] != 3 || move.Dirs[1] != 4 {
		t.Fatalf("move dirs = %#v, want [3 4]", move.Dirs)
	}

	b.executor.Execute(Action{
		Kind:    ActionAttack,
		Target:  9,
		Dir:     4,
		ReadyAt: now.Add(time.Second),
	})

	attackAction := waitAction(t, g)
	attack := attackAction.msg.(*model.MsgAttack)
	if attackAction.action != "Attack" || attack.Target != 9 || attack.Action != 120 || attack.Dir != 4 {
		t.Fatalf("attack action = %#v msg=%#v", attackAction, attack)
	}
}

func TestExecutorLifecyclePendingDoesNotChangeWorldPhase(t *testing.T) {
	g := newTestGame()
	b := newUnitBot(g, &stubPolicy{})

	b.executor.Execute(Action{Kind: ActionConnect})
	_ = waitConn(t, g)
	if b.world.phase != PhaseDisconnected {
		t.Fatalf("phase after connect action = %v, want disconnected", b.world.phase)
	}

	b.world.Handle(&model.MsgConnectReply{Result: 1, ID: 7})
	if b.world.phase != PhaseConnected {
		t.Fatalf("phase after reply = %v, want connected", b.world.phase)
	}
}

func TestExecutorSnapshotIsPure(t *testing.T) {
	b := newUnitBot(newTestGame(), &stubPolicy{})
	b.executor.current = Action{
		Kind: ActionMove,
		Path: []Position{{X: 2, Y: 1}},
	}
	b.executor.position = Position{X: 1, Y: 1}
	b.executor.move = moveExecution{
		path:       []Position{{X: 2, Y: 1}},
		nextStepAt: time.Unix(1, 0),
	}

	first := b.executor.Snapshot()
	first.CurrentAction.Path[0] = Position{X: 9, Y: 9}
	first.Move.Path[0] = Position{X: 9, Y: 9}
	second := b.executor.Snapshot()

	if second.CurrentAction.Path[0] != (Position{X: 2, Y: 1}) ||
		second.Move.Path[0] != (Position{X: 2, Y: 1}) ||
		second.Position != (Position{X: 1, Y: 1}) {
		t.Fatalf("snapshot mutated executor state: %#v", second)
	}
}

func TestBotTickDecidesOncePerSecond(t *testing.T) {
	g := newTestGame()
	policy := &stubPolicy{actions: []Action{
		{Kind: ActionConnect},
		{},
	}}
	b := newUnitBot(g, policy)
	now := time.Unix(0, 0)

	b.tick(now)
	_ = waitConn(t, g)
	b.tick(now.Add(100 * time.Millisecond))
	b.tick(now.Add(999 * time.Millisecond))
	if len(policy.calls) != 1 {
		t.Fatalf("policy calls before decision interval = %d, want 1", len(policy.calls))
	}
	assertNoAction(t, g)
	b.tick(now.Add(time.Second))

	if len(policy.calls) != 2 {
		t.Fatalf("policy calls = %d, want 2", len(policy.calls))
	}
	select {
	case conn := <-g.conns:
		t.Fatalf("duplicate connection = %#v", conn)
	case <-time.After(50 * time.Millisecond):
	}
}

func TestBotTickSkipsDecisionWhileContinuationExists(t *testing.T) {
	now := time.Unix(0, 0)
	tests := []struct {
		name  string
		setup func(*bot)
	}{
		{
			name: "move waiting",
			setup: func(b *bot) {
				b.executor.move = moveExecution{
					path:       []Position{{X: 2, Y: 1}},
					nextStepAt: now.Add(time.Millisecond),
				}
			},
		},
		{
			name: "attack waiting",
			setup: func(b *bot) {
				b.executor.current = Action{Kind: ActionAttack, Target: 9, Dir: 3}
				b.executor.readyAt = now.Add(time.Millisecond)
			},
		},
		{
			name: "skill waiting",
			setup: func(b *bot) {
				b.executor.current = Action{Kind: ActionUseSkill, Target: 9, Skill: skill.SkillIndexFireBall}
				b.executor.readyAt = now.Add(time.Millisecond)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGame()
			policy := &stubPolicy{actions: []Action{{Kind: ActionConnect}}}
			b := newUnitBot(g, policy)
			b.world.phase = PhasePlaying
			b.world.self = Actor{Alive: true}
			tt.setup(b)

			b.tick(now)

			if len(policy.calls) != 0 {
				t.Fatalf("policy calls = %d, want 0", len(policy.calls))
			}
			assertNoAction(t, g)
		})
	}
}

func TestBotTickExecutesDueContinuationBeforeDecision(t *testing.T) {
	g := newTestGame()
	policy := &stubPolicy{actions: []Action{{Kind: ActionConnect}}}
	b := newUnitBot(g, policy)
	b.id.Store(7)
	b.world.phase = PhasePlaying
	b.world.self = Actor{Index: 7, MapNumber: 0, X: 1, Y: 1, Alive: true}
	b.world.objects[9] = Actor{Index: 9, MapNumber: 0, X: 2, Y: 1, Alive: true, Attackable: true}
	b.executor.current = Action{Kind: ActionAttack, Target: 9, Dir: 3}
	b.executor.position = Position{X: 1, Y: 1}
	now := time.Unix(0, 0)

	b.tick(now)

	if len(policy.calls) != 0 {
		t.Fatalf("policy calls = %d, want 0", len(policy.calls))
	}
	if action := waitAction(t, g); action.action != "Attack" {
		t.Fatalf("action = %#v, want Attack", action)
	}
}

func TestMoveStartsImmediatelyAndAdvancesBetweenFullDecisions(t *testing.T) {
	g := newTestGame()
	policy := &stubPolicy{actions: []Action{
		{
			Kind:         ActionMove,
			Path:         []Position{{X: 2, Y: 1}, {X: 3, Y: 1}},
			SelfPosition: Position{X: 1, Y: 1},
			SentAt:       time.Unix(0, 0),
			NextStepAt:   time.Unix(0, 0).Add(400 * time.Millisecond),
		},
	}}
	b := newUnitBot(g, policy)
	b.id.Store(7)
	b.world.phase = PhasePlaying
	b.world.self = Actor{Index: 7, X: 1, Y: 1, Alive: true}
	now := time.Unix(0, 0)

	b.tick(now)
	_ = waitAction(t, g)
	first := b.executor.Snapshot()
	if !first.Move.Active || first.Move.NextStepAt != now.Add(400*time.Millisecond) {
		t.Fatalf("move snapshot = %#v", first.Move)
	}

	b.tick(now.Add(500 * time.Millisecond))
	if b.executor.position != (Position{X: 2, Y: 1}) {
		t.Fatalf("position after first continue = %+v, want (2,1)", b.executor.position)
	}

	b.tick(now.Add(600 * time.Millisecond))
	if b.executor.position != (Position{X: 2, Y: 1}) {
		t.Fatalf("position before second continue = %+v, want (2,1)", b.executor.position)
	}

	b.tick(now.Add(800 * time.Millisecond))
	if b.executor.position != (Position{X: 3, Y: 1}) || b.executor.Snapshot().Move.Active {
		t.Fatalf("final movement state = position %+v active=%v", b.executor.position, b.executor.Snapshot().Move.Active)
	}
	snapshot := b.executor.Snapshot()
	if snapshot.Position != (Position{X: 3, Y: 1}) {
		t.Fatalf("completed move position = %+v, want (3,1)", snapshot.Position)
	}

	b.world.HandleSetPositionReply(&model.MsgSetPositionReply{Number: 7, X: 4, Y: 1})
	b.executor.Execute(Action{
		Kind:            ActionSyncPosition,
		SelfPosition:    Position{X: 4, Y: 1},
		PositionVersion: b.world.positionVersion,
		CancelCurrent:   true,
	})
	snapshot = b.executor.Snapshot()
	if snapshot.Position != (Position{X: 4, Y: 1}) {
		t.Fatalf("authoritative position = %+v, want (4,1)", snapshot.Position)
	}
}

func TestContinueActionUsesUnifiedDueTimeForMoveAndCombat(t *testing.T) {
	now := time.Unix(1, 0)
	world := WorldSnapshot{
		Self:        Actor{MapNumber: 0},
		MP:          20,
		MagicSpeed:  1000,
		AttackSpeed: 100,
		Objects: []Actor{
			{Index: 9, MapNumber: 0, X: 2, Y: 1, Alive: true, Attackable: true},
		},
		Skills: []CombatSkill{
			{Index: skill.SkillIndexFireBall, Distance: 3, MP: 10, Type: skillTypeMagic, Delay: 700},
		},
	}
	moveExecution := ExecutorSnapshot{
		CurrentAction: Action{Kind: ActionMove},
		Position:      Position{X: 1, Y: 1},
		Move: MoveSnapshot{
			Active:     true,
			Path:       []Position{{X: 2, Y: 1}},
			NextStepAt: now.Add(time.Millisecond),
		},
	}
	if action, ok := continueAction(now, world, moveExecution); !ok || action.Kind != ActionNone {
		t.Fatalf("move before due action = %#v ok=%v, want none with continuation", action, ok)
	}
	moveExecution.Move.NextStepAt = now
	if action, ok := continueAction(now, world, moveExecution); !ok || action.Kind != ActionContinueMove {
		t.Fatalf("move due action = %#v ok=%v, want continue move", action, ok)
	}

	attackExecution := ExecutorSnapshot{
		CurrentAction: Action{Kind: ActionAttack, Target: 9, Dir: 3},
		Position:      Position{X: 1, Y: 1},
		ReadyAt:       now.Add(time.Millisecond),
	}
	if action, ok := continueAction(now, world, attackExecution); !ok || action.Kind != ActionNone {
		t.Fatalf("attack before due action = %#v ok=%v, want none with continuation", action, ok)
	}
	attackExecution.ReadyAt = now
	if action, ok := continueAction(now, world, attackExecution); !ok || action.Kind != ActionContinueAttack {
		t.Fatalf("attack due action = %#v ok=%v, want continue attack", action, ok)
	}

	skillExecution := ExecutorSnapshot{
		CurrentAction: Action{Kind: ActionUseSkill, Target: 9, Skill: skill.SkillIndexFireBall},
		Position:      Position{X: 1, Y: 1},
		ReadyAt:       now,
	}
	if action, ok := continueAction(now, world, skillExecution); !ok || action.Kind != ActionContinueUseSkill {
		t.Fatalf("skill due action = %#v ok=%v, want continue use skill", action, ok)
	}

	if action, ok := continueAction(now, world, ExecutorSnapshot{}); ok || action.Kind != ActionNone {
		t.Fatalf("no continuation action = %#v ok=%v, want none without continuation", action, ok)
	}
}

func TestMoveContinuationCancelsInvalidPathIndex(t *testing.T) {
	now := time.Unix(1, 0)
	action, ok := continueAction(now, WorldSnapshot{}, ExecutorSnapshot{
		Move: MoveSnapshot{
			Active:     true,
			Path:       []Position{{X: 2, Y: 1}},
			PathNext:   1,
			NextStepAt: now,
		},
	})
	if !ok || action.Kind != ActionCancel {
		t.Fatalf("action = %#v ok=%v, want cancel with continuation", action, ok)
	}
}

func TestAttackContinuesAfterCooldownWithoutFullDecision(t *testing.T) {
	g := newTestGame()
	policy := newRulePolicy(&resources{}, "account1:char1")
	b := newUnitBot(g, policy)
	b.id.Store(7)
	b.world.phase = PhasePlaying
	b.world.self = Actor{Index: 7, X: 1, Y: 1, Alive: true}
	b.world.objects[9] = Actor{Index: 9, X: 2, Y: 1, Alive: true, Attackable: true}
	b.world.Handle(&model.MsgAttackSpeedReply{AttackSpeed: 100})
	now := time.Unix(0, 0)
	b.executor.Execute(Action{
		Kind:            ActionSyncPosition,
		SelfPosition:    Position{X: 1, Y: 1},
		PositionVersion: b.world.positionVersion,
	})

	b.tick(now)
	_ = waitAction(t, g)
	b.tick(now.Add(100 * time.Millisecond))
	assertNoAction(t, g)
	b.tick(now.Add(267 * time.Millisecond))
	if action := waitAction(t, g); action.action != "Attack" {
		t.Fatalf("action = %#v, want Attack", action)
	}
	if b.executor.current.Kind != ActionContinueAttack {
		t.Fatalf("current action = %v, want continue attack", b.executor.current.Kind)
	}

	combatSkill := CombatSkill{Type: skillTypeMagic, Delay: 700}
	if got := skillInterval(combatSkill, WorldSnapshot{MagicSpeed: 1000}); got != 700*time.Millisecond {
		t.Fatalf("magic skill interval = %v, want 700ms", got)
	}
	combatSkill = CombatSkill{Type: skillTypePhysical}
	if got := skillInterval(combatSkill, WorldSnapshot{AttackSpeed: 100}); got != 267*time.Millisecond {
		t.Fatalf("physical skill interval = %v, want 267ms", got)
	}
	combatSkill = CombatSkill{Type: skillTypeSiege}
	if got := skillInterval(combatSkill, WorldSnapshot{AttackSpeed: 100}); got != 267*time.Millisecond {
		t.Fatalf("siege skill interval = %v, want 267ms", got)
	}
}

func TestSkillContinuesAfterCooldownWithoutFullDecision(t *testing.T) {
	g := newTestGame()
	policy := newRulePolicy(&resources{}, "account1:char1")
	b := newUnitBot(g, policy)
	b.id.Store(7)
	b.world.phase = PhasePlaying
	b.world.self = Actor{Index: 7, X: 1, Y: 1, Alive: true}
	b.world.objects[9] = Actor{Index: 9, X: 3, Y: 1, Alive: true, Attackable: true}
	b.world.mp = 20
	b.world.magicSpeed = 1000
	b.world.skills[skill.SkillIndexFireBall] = &skill.Skill{
		SkillBase: &skill.SkillBase{
			Index:     skill.SkillIndexFireBall,
			ManaUsage: 10,
			Distance:  3,
			Delay:     700,
			Type:      skillTypeMagic,
		},
		Index:     skill.SkillIndexFireBall,
		DamageMax: 10,
	}
	now := time.Unix(0, 0)
	b.executor.Execute(Action{
		Kind:            ActionSyncPosition,
		SelfPosition:    Position{X: 1, Y: 1},
		PositionVersion: b.world.positionVersion,
	})

	b.tick(now)
	first := waitAction(t, g)
	if first.action != "UseSkill" {
		t.Fatalf("action = %#v, want UseSkill", first)
	}
	b.tick(now.Add(699 * time.Millisecond))
	assertNoAction(t, g)
	b.tick(now.Add(700 * time.Millisecond))
	continued := waitAction(t, g)
	if continued.action != "UseSkill" {
		t.Fatalf("action = %#v, want UseSkill", continued)
	}
	if b.executor.current.Kind != ActionContinueUseSkill {
		t.Fatalf("current action = %v, want continue use skill", b.executor.current.Kind)
	}
}

func TestAttackContinuationCancelsWhenTargetIsInvalid(t *testing.T) {
	g := newTestGame()
	policy := newRulePolicy(&resources{}, "account1:char1")
	b := newUnitBot(g, policy)
	b.id.Store(7)
	b.world.phase = PhasePlaying
	b.world.self = Actor{Index: 7, X: 1, Y: 1, Alive: true}
	b.world.objects[9] = Actor{Index: 9, X: 2, Y: 1, Alive: true, Attackable: true}
	b.world.Handle(&model.MsgAttackSpeedReply{AttackSpeed: 100})
	now := time.Unix(0, 0)
	b.executor.Execute(Action{
		Kind:            ActionSyncPosition,
		SelfPosition:    Position{X: 1, Y: 1},
		PositionVersion: b.world.positionVersion,
	})

	b.tick(now)
	_ = waitAction(t, g)
	target := b.world.objects[9]
	target.Alive = false
	target.HP = 0
	b.world.objects[9] = target

	b.tick(now.Add(267 * time.Millisecond))
	assertNoAction(t, g)
	if b.executor.current.Kind != ActionNone {
		t.Fatalf("current action = %v, want none", b.executor.current.Kind)
	}
}

func TestAttackContinuationCancelsWhenTargetLeavesRange(t *testing.T) {
	g := newTestGame()
	policy := newRulePolicy(&resources{}, "account1:char1")
	b := newUnitBot(g, policy)
	b.id.Store(7)
	b.world.phase = PhasePlaying
	b.world.self = Actor{Index: 7, X: 1, Y: 1, Alive: true}
	b.world.objects[9] = Actor{Index: 9, X: 2, Y: 1, Alive: true, Attackable: true}
	b.world.Handle(&model.MsgAttackSpeedReply{AttackSpeed: 100})
	now := time.Unix(0, 0)
	b.executor.Execute(Action{
		Kind:            ActionSyncPosition,
		SelfPosition:    Position{X: 1, Y: 1},
		PositionVersion: b.world.positionVersion,
	})

	b.tick(now)
	_ = waitAction(t, g)
	target := b.world.objects[9]
	target.X = 5
	b.world.objects[9] = target

	b.tick(now.Add(267 * time.Millisecond))
	assertNoAction(t, g)
	if b.executor.current.Kind != ActionNone {
		t.Fatalf("current action = %v, want none", b.executor.current.Kind)
	}
}

func TestSkillContinuationCancelsWhenResourcesAreMissing(t *testing.T) {
	g := newTestGame()
	policy := newRulePolicy(&resources{}, "account1:char1")
	b := newUnitBot(g, policy)
	b.id.Store(7)
	b.world.phase = PhasePlaying
	b.world.self = Actor{Index: 7, X: 1, Y: 1, Alive: true}
	b.world.objects[9] = Actor{Index: 9, X: 3, Y: 1, Alive: true, Attackable: true}
	b.world.mp = 20
	b.world.magicSpeed = 1000
	b.world.skills[skill.SkillIndexFireBall] = &skill.Skill{
		SkillBase: &skill.SkillBase{
			Index:     skill.SkillIndexFireBall,
			ManaUsage: 10,
			Distance:  3,
			Delay:     700,
			Type:      skillTypeMagic,
		},
		Index:     skill.SkillIndexFireBall,
		DamageMax: 10,
	}
	now := time.Unix(0, 0)
	b.executor.Execute(Action{
		Kind:            ActionSyncPosition,
		SelfPosition:    Position{X: 1, Y: 1},
		PositionVersion: b.world.positionVersion,
	})

	b.tick(now)
	_ = waitAction(t, g)
	b.world.mp = 0

	b.tick(now.Add(700 * time.Millisecond))
	assertNoAction(t, g)
	if b.executor.current.Kind != ActionNone {
		t.Fatalf("current action = %v, want none", b.executor.current.Kind)
	}
}

func TestLearnSkillRemainsPendingUntilWorldContainsSkill(t *testing.T) {
	g := newTestGame()
	policy := newRulePolicy(&resources{}, "account1:char1")
	b := newUnitBot(g, policy)
	b.id.Store(7)
	b.world.phase = PhasePlaying
	b.world.self = Actor{Index: 7, Class: int(class.Wizard), Level: 1, Alive: true}
	b.world.itemsSet = true
	b.world.skillsSet = true
	b.world.characterSet = true
	b.world.selfClassSet = true
	b.world.inventory = make([]*item.Item, 13)
	b.world.inventory[12] = skillBook(skill.SkillIndexFireBall)
	now := time.Unix(0, 0)

	b.tick(now)
	if action := waitAction(t, g); action.action != "UseItem" {
		t.Fatalf("action = %#v, want UseItem", action)
	}
	b.tick(now.Add(100 * time.Millisecond))
	assertNoAction(t, g)
	b.world.HandleSetPositionReply(&model.MsgSetPositionReply{Number: 7, X: 1, Y: 1})
	b.tick(now.Add(200 * time.Millisecond))
	assertNoAction(t, g)
	if b.executor.current.Kind != ActionLearnSkill {
		t.Fatal("position correction released learn skill pending")
	}

	b.world.Handle(&model.MsgSkillOneReply{
		Flag: -2,
		Skill: &skill.Skill{
			SkillBase: &skill.SkillBase{Index: skill.SkillIndexFireBall},
			Index:     skill.SkillIndexFireBall,
		},
	})
	b.tick(now.Add(300 * time.Millisecond))
	assertNoAction(t, g)
	if b.executor.current.Kind != ActionLearnSkill {
		t.Fatal("learn skill action released before next full decision")
	}
	b.tick(now.Add(time.Second))
	assertNoAction(t, g)
	if b.executor.current.Kind != ActionLearnSkill {
		t.Fatal("position sync released learn skill pending")
	}
	b.tick(now.Add(2 * time.Second))
	assertNoAction(t, g)
	if b.executor.current.Kind == ActionLearnSkill {
		t.Fatal("learn skill action remained pending after skill reply")
	}
}
