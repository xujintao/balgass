package bot

import (
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/model"
)

func TestWorldMoveReplyDoesNotTeleportSelf(t *testing.T) {
	world := newWorld(&resources{})
	world.setSelfIndex(7)
	world.self.Alive = true
	world.self.X, world.self.Y = 1, 1
	now := time.Now()
	world.startMove([]Position{{X: 2, Y: 1}, {X: 3, Y: 1}}, now)

	world.handle(&model.MsgMoveReply{Number: 7, X: 3, Y: 1, Dir: 3 << 4})
	if world.self.X != 1 || world.self.Y != 1 {
		t.Fatalf("self position = (%d,%d), want (1,1)", world.self.X, world.self.Y)
	}
	world.advance(now.Add(399 * time.Millisecond))
	if world.self.X != 1 || world.self.Y != 1 {
		t.Fatalf("self position before step = (%d,%d), want (1,1)", world.self.X, world.self.Y)
	}
	world.advance(now.Add(400 * time.Millisecond))
	if world.self.X != 2 || world.self.Y != 1 {
		t.Fatalf("self position after first step = (%d,%d), want (2,1)", world.self.X, world.self.Y)
	}
	world.advance(now.Add(800 * time.Millisecond))
	if world.self.X != 3 || world.self.Y != 1 || world.moving() {
		t.Fatalf("self final state = (%d,%d) moving=%v, want (3,1) false", world.self.X, world.self.Y, world.moving())
	}
}

func TestWorldReloadClearsSightAndRevivesSelf(t *testing.T) {
	world := newWorld(&resources{})
	world.setSelfIndex(7)
	world.self.Alive = false
	world.objects[1] = Actor{Index: 1, Alive: true}

	world.handle(&model.MsgReloadCharacterReply{MapNumber: 0, X: 10, Y: 20, Dir: 3, HP: 100})

	if !world.self.Alive || world.self.X != 10 || world.self.Y != 20 {
		t.Fatalf("self = %#v, want revived at (10,20)", world.self)
	}
	if len(world.objects) != 0 {
		t.Fatalf("object count = %d, want 0", len(world.objects))
	}
}
