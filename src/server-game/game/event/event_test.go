package event

import (
	"errors"
	"testing"
	"time"
)

// NewTestManager creates an isolated manager for external integration tests.
func NewTestManager() *eventManager { return &eventManager{} }

// RegisteredEventsForTest returns a snapshot without exposing mutable registrations.
func RegisteredEventsForTest() map[string]Event {
	result := make(map[string]Event)
	for _, entry := range EventManager.activities {
		result[entry.name] = entry.activity
	}
	return result
}

func TestManagerStartupBaseline(t *testing.T) {
	if !EventManager.last.IsZero() {
		t.Fatal("scheduling started before the game loop")
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
	m := &eventManager{}
	a, failed, monday := &testActivity{}, &testActivity{fail: true}, &testActivity{}
	registrations := []struct {
		name string
		a    *testActivity
		s    []Schedule
	}{
		{"dragon", a, []Schedule{{DayOfWeek: -1, AppearanceRate: 100}, {DayOfWeek: -1, Minute: 1, AppearanceRate: 100}}},
		{"skeleton", failed, []Schedule{{DayOfWeek: 0, AppearanceRate: 100}}},
		{"group", monday, []Schedule{{DayOfWeek: 1, AppearanceRate: 100}}},
	}
	for _, r := range registrations {
		if err := m.Register(r.name, r.a, r.s); err != nil {
			t.Fatal(err)
		}
	}
	m.Tick(base)
	m.Tick(base.Add(time.Second))
	if a.starts != 0 {
		t.Fatal("early start")
	}
	m.Tick(base.Add(7 * time.Second))
	if a.starts != 1 || failed.starts != 1 || monday.starts != 0 {
		t.Fatal("weekday/catch-up failed")
	}
	m.Tick(base.Add(20 * time.Second))
	m.Tick(base.Add(-time.Minute))
	m.Tick(base.Add(25 * time.Second))
	if a.starts != 1 || failed.starts != 1 {
		t.Fatal("repeated occurrence")
	}
	m.Tick(base.Add(62 * time.Second))
	if a.starts != 2 {
		t.Fatal("running event was not restarted")
	}
	m.Tick(base.Add(63 * time.Second))
	if a.starts != 2 {
		t.Fatal("repeated restart in the same occurrence")
	}
	m.Tick(base.Add(24*time.Hour + 7*time.Second))
	if a.starts != 3 || monday.starts != 1 || monday.ticks == 0 {
		t.Fatal("next day dispatch failed")
	}
}
func TestManagerBaselineProbabilityAndRefresh(t *testing.T) {
	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.Local)
	m := &eventManager{}
	group, never := &testActivity{}, &testActivity{}
	schedule := []Schedule{{DayOfWeek: -1, Hour: 12, AppearanceRate: 100}}
	if err := m.Register("monstergroup/0", group, schedule); err != nil {
		t.Fatal(err)
	}
	if err := m.Register("zero-rate", never, []Schedule{{DayOfWeek: -1, Hour: 12, AppearanceRate: 0}}); err != nil {
		t.Fatal(err)
	}
	m.Tick(now)
	if !m.last.Equal(now) || group.starts != 0 || group.ticks != 0 || m.activities[0].activity != group {
		t.Fatal("first Tick changed state")
	}
	m.Tick(now.Add(time.Second))
	if group.starts != 0 {
		t.Fatal("replayed pre-baseline occurrence")
	}
	m.Tick(now.AddDate(0, 0, 1))
	m.Tick(now.AddDate(0, 0, 1).Add(time.Second))
	m.Tick(now.AddDate(0, 0, 2))
	if group.starts != 2 || never.starts != 0 {
		t.Fatal("rate/refresh policy failed")
	}
	restarted := &eventManager{}
	fresh := &testActivity{}
	if err := restarted.Register("group", fresh, schedule); err != nil {
		t.Fatal(err)
	}
	restarted.Tick(now.AddDate(0, 0, 2))
	restarted.Tick(now.AddDate(0, 0, 2).Add(time.Second))
	if fresh.starts != 0 {
		t.Fatal("restart replayed history")
	}
}
func TestManagerRejectsInvalidDefinitions(t *testing.T) {
	valid := Schedule{DayOfWeek: -1, Hour: 12, AppearanceRate: 100}
	for _, s := range [][]Schedule{{{DayOfWeek: 7}}, {{DayOfWeek: -1, Hour: 24}}, {{DayOfWeek: -1, Minute: 60}}, {{DayOfWeek: -1, AppearanceRate: 101}}, {valid, valid}} {
		if err := (&eventManager{}).Register("test", &testActivity{}, s); err == nil {
			t.Fatal("accepted invalid schedule", s)
		}
	}
	m := &eventManager{}
	if err := m.Register("same", &testActivity{}, nil); err != nil {
		t.Fatal(err)
	}
	if err := m.Register("same", &testActivity{}, nil); err == nil {
		t.Fatal("duplicate registration accepted")
	}
}
