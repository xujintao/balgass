package bot

import (
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/model"
)

type stubPolicy struct {
	actions []Action
	calls   int
}

func (p *stubPolicy) Decide(Snapshot) Action {
	call := p.calls
	p.calls++
	if call >= len(p.actions) {
		return Action{}
	}
	return p.actions[call]
}

func TestBotExecuteTranslatesMoveAndAttackActions(t *testing.T) {
	g := newTestGame()
	b := &bot{
		game:  g,
		world: newWorld(&resources{}),
	}
	b.id.Store(7)
	b.world.self = Actor{Index: 7, X: 1, Y: 1, Alive: true}

	b.execute(Action{
		Kind: ActionMove,
		Path: []Position{{X: 2, Y: 1}, {X: 3, Y: 2}},
	}, time.Now())

	moveAction := waitAction(t, g)
	if moveAction.id != 7 || moveAction.action != "Move" {
		t.Fatalf("move action = %#v, want player 7 Move", moveAction)
	}
	move := moveAction.msg.(*model.MsgMove)
	if move.X != 1 || move.Y != 1 || len(move.Path) != 2 || move.Path[1].X != 3 || move.Path[1].Y != 2 {
		t.Fatalf("move msg = %#v, want path from (1,1) to (3,2)", move)
	}
	if len(move.Dirs) != 2 || move.Dirs[0] != 3 || move.Dirs[1] != 4 {
		t.Fatalf("move dirs = %#v, want [3 4]", move.Dirs)
	}

	b.world.movement = movement{}
	b.world.objects[9] = Actor{Index: 9, X: 2, Y: 2, Alive: true, Attackable: true}
	b.execute(Action{Kind: ActionAttack, Target: 9}, time.Now())

	attackAction := waitAction(t, g)
	if attackAction.id != 7 || attackAction.action != "Attack" {
		t.Fatalf("attack action = %#v, want player 7 Attack", attackAction)
	}
	attack := attackAction.msg.(*model.MsgAttack)
	if attack.Target != 9 || attack.Action != 120 || attack.Dir != 4 {
		t.Fatalf("attack msg = %#v, want target 9 action 120 dir 4", attack)
	}
}

func TestBotTickDoesNotRedecideDuringMove(t *testing.T) {
	g := newTestGame()
	policy := &stubPolicy{
		actions: []Action{
			{
				Kind: ActionMove,
				Path: []Position{{X: 2, Y: 1}, {X: 3, Y: 1}, {X: 4, Y: 1}, {X: 5, Y: 1}},
			},
		},
	}
	b := &bot{
		game:   g,
		world:  newWorld(&resources{}),
		policy: policy,
	}
	b.id.Store(7)
	b.world.playing = true
	b.world.self = Actor{Index: 7, X: 1, Y: 1, Alive: true}
	now := time.Unix(0, 0)

	b.tick(now)
	moveAction := waitAction(t, g)
	if moveAction.action != "Move" {
		t.Fatalf("action = %#v, want Move", moveAction)
	}
	if policy.calls != 1 {
		t.Fatalf("policy calls = %d, want 1", policy.calls)
	}

	b.tick(now.Add(500 * time.Millisecond))
	assertNoAction(t, g)
	if policy.calls != 1 {
		t.Fatalf("policy calls while moving = %d, want 1", policy.calls)
	}

	b.tick(now.Add(1600 * time.Millisecond))
	assertNoAction(t, g)
	if policy.calls != 1 || b.action.kind != actionExecutionIdle {
		t.Fatalf("after move complete policy calls=%d action=%v, want 1 idle", policy.calls, b.action.kind)
	}

	b.tick(now.Add(1700 * time.Millisecond))
	assertNoAction(t, g)
	if policy.calls != 2 {
		t.Fatalf("policy calls after move complete = %d, want 2", policy.calls)
	}
}

func TestBotTickRepeatsAttackUntilInterrupted(t *testing.T) {
	g := newTestGame()
	policy := &stubPolicy{
		actions: []Action{{Kind: ActionAttack, Target: 9}},
	}
	b := &bot{
		game:   g,
		world:  newWorld(&resources{}),
		policy: policy,
	}
	b.id.Store(7)
	b.world.playing = true
	b.world.self = Actor{Index: 7, X: 1, Y: 1, Alive: true}
	b.world.objects[9] = Actor{Index: 9, X: 2, Y: 1, Alive: true, Attackable: true}
	now := time.Unix(0, 0)

	b.tick(now)
	first := waitAction(t, g)
	if first.action != "Attack" {
		t.Fatalf("first action = %#v, want Attack", first)
	}
	if policy.calls != 1 {
		t.Fatalf("policy calls = %d, want 1", policy.calls)
	}

	b.tick(now.Add(500 * time.Millisecond))
	assertNoAction(t, g)
	if policy.calls != 1 {
		t.Fatalf("policy calls during attack = %d, want 1", policy.calls)
	}

	b.tick(now.Add(defaultAttackInterval))
	second := waitAction(t, g)
	if second.action != "Attack" {
		t.Fatalf("second action = %#v, want Attack", second)
	}

	target := b.world.objects[9]
	target.Alive = false
	b.world.objects[9] = target
	b.tick(now.Add(2 * defaultAttackInterval))
	assertNoAction(t, g)
	if b.action.kind != actionExecutionIdle {
		t.Fatalf("action after target death = %v, want idle", b.action.kind)
	}
}

func TestBotTickDoesNotAdvanceMovementDuringAttack(t *testing.T) {
	g := newTestGame()
	b := &bot{
		game:   g,
		world:  newWorld(&resources{}),
		policy: &stubPolicy{},
		action: actionExecution{
			kind:       actionExecutionAttacking,
			target:     9,
			nextAttack: time.Unix(0, 0).Add(time.Hour),
		},
	}
	b.id.Store(7)
	b.world.playing = true
	b.world.self = Actor{Index: 7, X: 1, Y: 1, Alive: true}
	b.world.objects[9] = Actor{Index: 9, X: 2, Y: 1, Alive: true, Attackable: true}
	now := time.Unix(0, 0)
	b.world.startMove([]Position{{X: 2, Y: 1}}, now)

	b.tick(now.Add(time.Second))

	assertNoAction(t, g)
	if b.world.self.X != 1 || b.world.self.Y != 1 {
		t.Fatalf("self position = (%d,%d), want (1,1)", b.world.self.X, b.world.self.Y)
	}
	if !b.world.moving() {
		t.Fatal("world moving = false, want residual movement untouched during attack")
	}
}

func TestBotTickStopsAttackWhenTargetLeavesRange(t *testing.T) {
	g := newTestGame()
	policy := &stubPolicy{
		actions: []Action{
			{Kind: ActionAttack, Target: 9},
			{Kind: ActionMove, Path: []Position{{X: 2, Y: 1}}},
		},
	}
	b := &bot{
		game:   g,
		world:  newWorld(&resources{}),
		policy: policy,
	}
	b.id.Store(7)
	b.world.playing = true
	b.world.self = Actor{Index: 7, X: 1, Y: 1, Alive: true}
	b.world.objects[9] = Actor{Index: 9, X: 2, Y: 1, Alive: true, Attackable: true}
	now := time.Unix(0, 0)

	b.tick(now)
	if action := waitAction(t, g); action.action != "Attack" {
		t.Fatalf("action = %#v, want Attack", action)
	}

	target := b.world.objects[9]
	target.X, target.Y = 5, 1
	b.world.objects[9] = target
	b.tick(now.Add(defaultAttackInterval))
	assertNoAction(t, g)
	if b.action.kind != actionExecutionIdle {
		t.Fatalf("action after target leaves range = %v, want idle", b.action.kind)
	}

	b.tick(now.Add(defaultAttackInterval + 100*time.Millisecond))
	if action := waitAction(t, g); action.action != "Move" {
		t.Fatalf("action after attack interrupted = %#v, want Move", action)
	}
}
