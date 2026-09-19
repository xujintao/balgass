package event

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/event/invasion"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/object/monster"
)

func TestScheduleConfigLoadsAndValidates(t *testing.T) {
	const times = `<InvasionManager Enable="0"><Invasion Type="0"><Start DayOfWeek="-1" Hour="12" Minute="0"/></Invasion><Invasion Type="1"><Start DayOfWeek="0" Hour="13" Minute="0"/></Invasion></InvasionManager>`
	for _, tt := range []struct {
		name, body string
		bad        bool
	}{
		{"valid", times, false},
		{"enum", strings.Replace(times, `Type="1"`, `Type="9"`, 1), true},
		{"required", strings.Replace(times, ` Minute="0"`, "", 1), true},
		{"time", strings.Replace(times, `Hour="12"`, `Hour="24"`, 1), true},
		{"duplicate", strings.Replace(times, `</Invasion>`, `<Start DayOfWeek="-1" Hour="12" Minute="0"/></Invasion>`, 1), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			base := t.TempDir()
			if err := os.Mkdir(filepath.Join(base, "Events"), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(base, "Events/IGC_InvasionManager.xml"), []byte(tt.body), 0600); err != nil {
				t.Fatal(err)
			}
			var m eventManager
			err := m.loadSchedules(base)
			if (err != nil) != tt.bad {
				t.Fatalf("load error=%v", err)
			}
			if !tt.bad && (m.enabled || len(m.schedules) != 2) {
				t.Fatal("lost disabled state or unsupported schedules")
			}
		})
	}
}

func TestEventManagerPackageInitialization(t *testing.T) {
	var configured eventManager
	if err := configured.loadSchedules(conf.PathCommon); err != nil {
		t.Fatal(err)
	}
	if EventManager.enabled != configured.enabled || !reflect.DeepEqual(EventManager.schedules, configured.schedules) {
		t.Fatal("package manager did not load the schedule configuration")
	}
	dragon, ok := EventManager.activities[invasion.RedDragon].(*invasion.Dragon)
	if !ok || dragon.Running() || len(EventManager.activities) != 1 {
		t.Fatal("package manager must register one idle red dragon activity")
	}
	if !EventManager.last.IsZero() {
		t.Fatal("scheduling started before the game loop")
	}
}

func TestScheduledDragonWithWorldObjects(t *testing.T) {
	base := time.Date(2026, 9, 5, 11, 59, 59, 0, time.Local)
	m, err := newManager(true, []Schedule{{ActivityID: invasion.RedDragon, DayOfWeek: -1, Hour: 12}})
	if err != nil {
		t.Fatal(err)
	}
	dragon := invasion.NewDragon()
	if err := m.register(invasion.RedDragon, dragon); err != nil {
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

type testActivity struct {
	starts, ticks int
	running       bool
	fail          bool
}

func (a *testActivity) Start(time.Time) error {
	a.starts++
	if a.fail {
		return errors.New("start failed")
	}
	a.running = true
	return nil
}
func (a *testActivity) Tick(time.Time) { a.ticks++ }
func (a *testActivity) Running() bool  { return a.running }

func TestManagerScheduling(t *testing.T) {
	base := time.Date(2026, 9, 5, 23, 59, 58, 0, time.Local)
	schedules := []Schedule{
		{ActivityID: 0, DayOfWeek: -1, Hour: 0, Minute: 0},
		{ActivityID: 0, DayOfWeek: -1, Hour: 0, Minute: 1},
		{ActivityID: 1, DayOfWeek: 0, Hour: 0, Minute: 0}, // Sunday
		{ActivityID: 2, DayOfWeek: 1, Hour: 0, Minute: 0},
		{ActivityID: 3, DayOfWeek: -1, Hour: 0, Minute: 0}, // not implemented
	}
	m, err := newManager(true, schedules)
	if err != nil {
		t.Fatal(err)
	}
	a, failed, monday := &testActivity{}, &testActivity{fail: true}, &testActivity{}
	for id, act := range map[int]*testActivity{0: a, 1: failed, 2: monday} {
		if err := m.register(id, act); err != nil {
			t.Fatal(err)
		}
	}
	m.Tick(base)
	m.Tick(base.Add(time.Second))
	if a.starts != 0 {
		t.Fatal("started early")
	}
	m.Tick(base.Add(7 * time.Second)) // crossed midnight during a short stall
	if a.starts != 1 || failed.starts != 1 || monday.starts != 0 {
		t.Fatal("weekday or catch-up dispatch failed")
	}
	m.Tick(base.Add(20 * time.Second))
	m.Tick(base.Add(-time.Minute)) // clock rollback
	m.Tick(base.Add(25 * time.Second))
	if a.starts != 1 || failed.starts != 1 {
		t.Fatal("repeated occurrence or retried failure")
	}
	m.Tick(base.Add(62 * time.Second))
	if a.starts != 1 {
		t.Fatal("overlapping event restarted")
	}
	a.running = false
	m.Tick(base.Add(63 * time.Second))
	if a.starts != 1 {
		t.Fatal("overlap occurrence retried")
	}
	m.Tick(base.Add(24*time.Hour + 7*time.Second))
	if a.starts != 2 || monday.starts != 1 {
		t.Fatal("next day not dispatched")
	}
	if monday.ticks == 0 {
		t.Fatal("idle activities were not ticked")
	}
}

func TestManagerDisabledAndRestartBaseline(t *testing.T) {
	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.Local)
	for _, enabled := range []bool{false, true} {
		m, err := newManager(enabled, []Schedule{{DayOfWeek: -1, Hour: 12}})
		if err != nil {
			t.Fatal(err)
		}
		a := &testActivity{}
		if err := m.register(0, a); err != nil {
			t.Fatal(err)
		}
		if !m.last.IsZero() {
			t.Fatal("constructor established scheduling baseline")
		}
		m.Tick(now)
		if !m.last.Equal(now) || a.ticks != 0 {
			t.Fatal("first Tick must only establish the baseline")
		}
		if m.activities[0] != a || a.starts != 0 || m.enabled != enabled || len(m.schedules) != 1 {
			t.Fatal("first Tick replaced configuration or started an activity")
		}
		m.Tick(now.Add(10 * time.Second))
		if a.starts != 0 {
			t.Fatal("replayed schedule at or before the first Tick")
		}
		m.Tick(now.Add(24 * time.Hour))
		want := 0
		if enabled {
			want = 1
		}
		if a.starts != want {
			t.Fatalf("enabled=%v starts=%d", enabled, a.starts)
		}
	}
}

func TestManagerRejectsInvalidDefinitions(t *testing.T) {
	valid := Schedule{DayOfWeek: -1, Hour: 12}
	for _, schedules := range [][]Schedule{{{DayOfWeek: 7}}, {{DayOfWeek: -1, Hour: 24}}, {{DayOfWeek: -1, Minute: 60}}, {valid, valid}} {
		if _, err := newManager(true, schedules); err == nil {
			t.Fatalf("accepted %v", schedules)
		}
	}
	m, _ := newManager(true, nil)
	a := &testActivity{}
	if err := m.register(0, a); err != nil {
		t.Fatal(err)
	}
	if err := m.register(0, a); err == nil {
		t.Fatal("duplicate registration accepted")
	}
}
