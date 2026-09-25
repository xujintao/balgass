package event_test

import (
	"fmt"
	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/event"
	"github.com/xujintao/balgass/src/server-game/game/event/invasion"
	"github.com/xujintao/balgass/src/server-game/game/event/monstergroup"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/object/monster"
	"testing"
	"time"
)

func TestEventPackagesRegisterAtStartup(t *testing.T) {
	registered := event.RegisteredEventsForTest()
	want := len(monstergroup.MonsterGroupManager.Groups)
	if invasion.InvasionManager.Enabled {
		want += 2
	}
	if len(registered) != want {
		t.Fatalf("registered %d events, want %d", len(registered), want)
	}
	for _, name := range []string{"invasion/dragon", "invasion/skeleton"} {
		if (registered[name] != nil) != invasion.InvasionManager.Enabled {
			t.Fatalf("incorrect invasion registration: %s", name)
		}
	}
	if invasion.InvasionManager.Enabled {
		if _, ok := registered["invasion/dragon"].(*invasion.Dragon); !ok {
			t.Fatal("wrong dragon event")
		}
		if _, ok := registered["invasion/skeleton"].(*invasion.Skeleton); !ok {
			t.Fatal("wrong skeleton event")
		}
	}
	for _, group := range monstergroup.MonsterGroupManager.Groups {
		if registered[fmt.Sprintf("monstergroup/%d", group.Index)] != group {
			t.Fatal("group was not independently registered")
		}
	}
	for name, activity := range registered {
		if activity.Running() {
			t.Fatalf("%s started during registration", name)
		}
	}
}

func TestScheduledDragonWithWorldObjects(t *testing.T) {
	base := time.Date(2026, 9, 5, 11, 59, 59, 0, time.Local)
	m := event.NewTestManager()
	dragon := invasion.NewDragon()
	if err := m.Register("dragon", dragon, []event.Schedule{{DayOfWeek: -1, Hour: 12, AppearanceRate: 100}}); err != nil {
		t.Fatal(err)
	}
	// The event package does not spawn the static world at initialization.
	// Seed an unrelated red dragon so ownership cleanup remains exercised.
	unrelated, err := monster.SpawnEventMonster(monster.EventSpawn{
		Class: 44, MapNumber: 0, StartX: 135, StartY: 61, EndX: 146, EndY: 70, Direction: -1, Distance: 30,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { object.ObjectManager.DeleteEventMonster(unrelated, base.Add(time.Hour)) })
	before := make(map[*object.Object]bool)
	for i := 0; i < conf.Server.GameServerInfo.MaxMonsterCount; i++ {
		if obj := object.ObjectManager.GetObject(i); obj != nil {
			before[obj] = true
		}
	}
	var spawned []*object.Object
	t.Cleanup(func() {
		for _, obj := range spawned {
			object.ObjectManager.DeleteEventMonster(obj, base.Add(time.Hour))
		}
	})
	m.Tick(base)
	m.Tick(base.Add(time.Second))
	if !dragon.Running() {
		t.Fatal("scheduled invasion not started")
	}
	m.Tick(base.Add(4 * time.Second))
	for i := 0; i < conf.Server.GameServerInfo.MaxMonsterCount; i++ {
		if obj := object.ObjectManager.GetObject(i); obj != nil && !before[obj] {
			spawned = append(spawned, obj)
		}
	}
	if len(spawned) != 2 {
		t.Fatalf("got %d dragons", len(spawned))
	}
	for _, obj := range spawned {
		if obj.Class != 44 || !obj.NoRegen || !obj.Live || obj.MapNumber != spawned[0].MapNumber {
			t.Fatal("invalid event monster")
		}
	}
	m.Tick(base.Add(4*time.Second + 5*time.Minute))
	if dragon.Running() {
		t.Fatal("invasion did not end")
	}
	for _, obj := range spawned {
		if object.ObjectManager.GetObject(obj.Index) == obj {
			t.Fatal("event monster still registered")
		}
	}
	for obj := range before {
		if object.ObjectManager.GetObject(obj.Index) != obj {
			t.Fatal("removed an unrelated object")
		}
	}
}

func TestScheduledMonsterGroupWithWorldObjects(t *testing.T) {
	if len(monstergroup.MonsterGroupManager.Groups) == 0 {
		t.Fatal("missing fixture group")
	}
	group := *monstergroup.MonsterGroupManager.Groups[0]
	base := time.Date(2026, 9, 5, 11, 59, 59, 0, time.Local)
	m := event.NewTestManager()
	schedules := []event.Schedule{{DayOfWeek: -1, Hour: 12, AppearanceRate: 100}}
	if err := m.Register("monstergroup/0", &group, schedules); err != nil {
		t.Fatal(err)
	}
	before := map[*object.Object]bool{}
	for i := 0; i < conf.Server.GameServerInfo.MaxMonsterCount; i++ {
		if obj := object.ObjectManager.GetObject(i); obj != nil {
			before[obj] = true
		}
	}
	t.Cleanup(func() {
		for i := 0; i < conf.Server.GameServerInfo.MaxMonsterCount; i++ {
			if obj := object.ObjectManager.GetObject(i); obj != nil && !before[obj] {
				object.ObjectManager.DeleteEventMonster(obj, base.Add(time.Hour))
			}
		}
	})
	m.Tick(base)
	m.Tick(base.Add(time.Second))
	if !group.Running() {
		t.Fatal("invasion switch prevented actual group spawn")
	}
	first := map[*object.Object]bool{}
	for i := 0; i < conf.Server.GameServerInfo.MaxMonsterCount; i++ {
		if obj := object.ObjectManager.GetObject(i); obj != nil && !before[obj] {
			first[obj] = true
		}
	}
	m.Tick(base.Add(2 * time.Second))
	for obj := range first {
		if object.ObjectManager.GetObject(obj.Index) != obj {
			t.Fatal("group refreshed twice in same minute")
		}
	}
	m.Tick(base.AddDate(0, 0, 1).Add(time.Second))
	for obj := range first {
		if object.ObjectManager.GetObject(obj.Index) == obj {
			t.Fatal("next group occurrence did not replace old monsters")
		}
	}
}
