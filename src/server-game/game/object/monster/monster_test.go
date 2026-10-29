package monster

import (
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/object"
)

func TestIsHiddenMonsterClass(t *testing.T) {
	for _, tt := range []struct {
		class int
		want  bool
	}{
		{class: 99, want: false},
		{class: 100, want: true},
		{class: 110, want: true},
		{class: 111, want: false},
		{class: 247, want: false},
		{class: 249, want: false},
		{class: 523, want: true},
		{class: 689, want: true},
	} {
		if got := isHiddenMonsterClass(tt.class); got != tt.want {
			t.Errorf("isHiddenMonsterClass(%d) = %t, want %t", tt.class, got, tt.want)
		}
	}
}

func TestEventSpawnAndExistingRedDragonBag(t *testing.T) {
	s := EventSpawn{Class: 44, MapNumber: 0, StartX: 135, StartY: 61, EndX: 146, EndY: 70, Direction: -1, Distance: 30}
	dragon, err := SpawnEventMonster(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { object.ObjectManager.DeleteEventMonster(dragon, time.Now()) })
	if dragon.Dir < 0 || dragon.Dir > 7 {
		t.Fatalf("invalid random direction: %d", dragon.Dir)
	}
	if dragon.Class != 44 || dragon.MaxHP <= 0 || !dragon.NoRegen || object.ObjectManager.GetObject(dragon.Index) != dragon {
		t.Fatal("incomplete spawn")
	}
	if dragon.X < s.StartX || dragon.X >= s.EndX || dragon.Y < s.StartY || dragon.Y >= s.EndY || maps.MapManager.GetMapAttr(0, dragon.X, dragon.Y)&2 == 0 {
		t.Fatal("invalid spawn tile")
	}
	oneTile := s
	oneTile.StartX, oneTile.StartY, oneTile.EndX, oneTile.EndY = dragon.X, dragon.Y, dragon.X+1, dragon.Y+1
	if _, err := SpawnEventMonster(oneTile); err == nil {
		t.Fatal("spawned on occupied tile")
	}
	before := map[int]bool{}
	maps.MapManager.MapEachItem(0, func(_ *item.Item, index, x, y int, new bool) { before[index] = true })
	t.Cleanup(func() {
		maps.MapManager.MapEachItem(0, func(_ *item.Item, index, x, y int, new bool) {
			if !before[index] {
				maps.MapManager.RemoveItem(0, index)
			}
		})
	})
	// Real Monster.DieDropItem -> DropManager -> configured BAG_MONSTER 44.
	killer := &object.Object{Objecter: dragon.Objecter, Class: 0, Level: 100, MapNumber: 0}
	dragon.DieDropItem(killer)
	rewards := 0
	maps.MapManager.MapEachItem(0, func(it *item.Item, index, x, y int, new bool) {
		if !before[index] {
			rewards++
			if it == nil {
				t.Error("nil reward")
			}
		}
	})
	if rewards != 1 {
		t.Fatalf("red dragon bag produced %d rewards", rewards)
	}
	// Deleting a living dragon does not generate an additional reward.
	if !object.ObjectManager.DeleteEventMonster(dragon, time.Now()) {
		t.Fatal("despawn failed")
	}
	after := 0
	maps.MapManager.MapEachItem(0, func(_ *item.Item, index, x, y int, new bool) {
		if !before[index] {
			after++
		}
	})
	if after != rewards {
		t.Fatal("despawn produced loot")
	}
}

func TestItemBagDropPosition(t *testing.T) {
	for i := 0; i < 20; i++ {
		x, y := itemBagDropPosition(0, 140, 65)
		if x < 136 || x > 144 || y < 61 || y > 69 || (x != 140 || y != 65) && maps.MapManager.GetMapAttr(0, x, y)&(4|8) != 0 {
			t.Fatalf("drop position (%d, %d) is outside the valid nearby area", x, y)
		}
	}
	if x, y := itemBagDropPosition(0, -10, -10); x != -10 || y != -10 {
		t.Fatalf("blocked nearby area returned (%d, %d), want source (-10, -10)", x, y)
	}
}

func TestEventSpawnRejectsInvalidOrBlockedArea(t *testing.T) {
	s := EventSpawn{Class: 44, MapNumber: 0, EndX: 1, EndY: 1, Direction: -1}
	for _, bad := range []EventSpawn{
		{Class: 44, MapNumber: 999, EndX: 1, EndY: 1},
		{Class: 99999, MapNumber: 0, EndX: 1, EndY: 1},
		{Class: 44, MapNumber: 0, StartX: 5, EndX: 1, EndY: 1},
	} {
		if _, err := SpawnEventMonster(bad); err == nil {
			t.Fatal("accepted invalid spawn", bad)
		}
	}
	found := false
	for y := 0; y < 256 && !found; y++ {
		for x := 0; x < 256; x++ {
			if maps.MapManager.GetMapAttr(0, x, y)&(1|4|8) != 0 {
				s.StartX, s.StartY, s.EndX, s.EndY = x, y, x+1, y+1
				found = true
				break
			}
		}
	}
	if !found {
		t.Fatal("no blocked fixture tile")
	}
	if _, err := SpawnEventMonster(s); err == nil {
		t.Fatal("spawned on blocked terrain")
	}
}

func TestRegenUpdatesStandAndFrustum(t *testing.T) {
	m := &Monster{
		Object:      object.Object{MapNumber: maps.Lorencia, X: 1, Y: 1},
		spawnStartX: 140, spawnStartY: 65, spawnEndX: 141, spawnEndY: 66,
	}
	wasOccupied := maps.MapManager.GetMapAttr(m.MapNumber, 140, 65)&2 != 0
	maps.MapManager.ClearMapAttrStand(m.MapNumber, 140, 65)
	t.Cleanup(func() {
		if wasOccupied {
			maps.MapManager.SetMapAttrStand(m.MapNumber, 140, 65)
		} else {
			maps.MapManager.ClearMapAttrStand(m.MapNumber, 140, 65)
		}
	})
	m.CreateFrustum()
	m.Regen()
	if m.X != 140 || m.Y != 65 || m.TX != m.X || m.TY != m.Y {
		t.Fatalf("unexpected regen position: (%d,%d), target (%d,%d)", m.X, m.Y, m.TX, m.TY)
	}
	if maps.MapManager.GetMapAttr(m.MapNumber, m.TX, m.TY)&2 == 0 {
		t.Fatal("regen position is not occupied")
	}
	for i := range m.FrustumX {
		if m.FrustumX[i] != object.FrustumX[i]+m.X || m.FrustumY[i] != object.FrustumY[i]+m.Y {
			t.Fatalf("regen frustum vertex %d does not match new position", i)
		}
	}
}
