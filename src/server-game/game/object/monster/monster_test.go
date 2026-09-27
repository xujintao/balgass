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

func TestNewMonsterSetsDefaultEventBagID(t *testing.T) {
	for _, tt := range []struct {
		class, eventBagID int
	}{{275, 26}, {673, 46}} {
		m := newMonster(tt.class, 0, 0, 0, 1, 1, 0, 0, 0)
		if m.EventBagID == nil || *m.EventBagID != tt.eventBagID {
			t.Fatalf("monster %d EventBagID = %v, want %d", tt.class, m.EventBagID, tt.eventBagID)
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
