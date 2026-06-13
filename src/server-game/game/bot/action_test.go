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

func TestBotTickDecidesEveryTick(t *testing.T) {
	g := newTestGame()
	policy := &stubPolicy{actions: []Action{
		{Kind: ActionConnect},
		{},
		{},
	}}
	b := newUnitBot(g, policy)
	now := time.Unix(0, 0)

	b.tick(now)
	_ = waitConn(t, g)
	b.tick(now.Add(100 * time.Millisecond))
	b.tick(now.Add(200 * time.Millisecond))

	if len(policy.calls) != 3 {
		t.Fatalf("policy calls = %d, want 3", len(policy.calls))
	}
	select {
	case conn := <-g.conns:
		t.Fatalf("duplicate connection = %#v", conn)
	case <-time.After(50 * time.Millisecond):
	}
}

func TestMoveStartsImmediatelyAndAdvancesOnlyByPolicyAction(t *testing.T) {
	g := newTestGame()
	policy := &stubPolicy{actions: []Action{
		{
			Kind:         ActionMove,
			Path:         []Position{{X: 2, Y: 1}, {X: 3, Y: 1}},
			SelfPosition: Position{X: 1, Y: 1},
			SentAt:       time.Unix(0, 0),
			NextStepAt:   time.Unix(0, 0).Add(400 * time.Millisecond),
		},
		{},
		{
			Kind:         ActionContinueMove,
			SelfPosition: Position{X: 2, Y: 1},
			Dir:          3,
			PathNext:     1,
			NextStepAt:   time.Unix(0, 0).Add(800 * time.Millisecond),
		},
		{
			Kind:         ActionContinueMove,
			SelfPosition: Position{X: 3, Y: 1},
			Dir:          3,
			PathNext:     2,
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
	if b.executor.position != (Position{X: 1, Y: 1}) {
		t.Fatalf("position without continue = %+v, want (1,1)", b.executor.position)
	}

	b.tick(now.Add(600 * time.Millisecond))
	if b.executor.position != (Position{X: 2, Y: 1}) {
		t.Fatalf("position after first continue = %+v, want (2,1)", b.executor.position)
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

func TestAttackAndSkillAreOneShotActionsWithCooldown(t *testing.T) {
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

	combatSkill := CombatSkill{Type: skillTypeMagic, Delay: 700}
	if got := skillInterval(combatSkill, WorldSnapshot{MagicSpeed: 1000}); got != 700*time.Millisecond {
		t.Fatalf("magic skill interval = %v, want 700ms", got)
	}
	combatSkill = CombatSkill{Type: skillTypePhysical}
	if got := skillInterval(combatSkill, WorldSnapshot{AttackSpeed: 100}); got != 267*time.Millisecond {
		t.Fatalf("physical skill interval = %v, want 267ms", got)
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
	if b.executor.current.Kind == ActionLearnSkill {
		t.Fatal("learn skill action remained pending after skill reply")
	}
}
