package bot

import (
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/model"
)

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
