package bot

import (
	"testing"
	"time"
)

func TestRulePolicyAttacksAdjacentMonsterAndIgnoresNPC(t *testing.T) {
	resources := &resources{
		terrains:          map[int]*terrain{0: newTestTerrain(8, 8)},
		attackableClasses: map[int]struct{}{1: {}},
	}
	policy := newRulePolicy(resources, "account1:char1")
	now := time.Now()
	action := policy.Decide(Snapshot{
		Now:  now,
		Self: Actor{Index: 10, MapNumber: 0, X: 1, Y: 1, Alive: true},
		Objects: []Actor{
			{Index: 20, Class: 100, MapNumber: 0, X: 1, Y: 2, Alive: true},
			{Index: 21, Class: 1, MapNumber: 0, X: 2, Y: 1, Alive: true, Attackable: true},
		},
	})
	if action.Kind != ActionAttack || action.Target != 21 {
		t.Fatalf("action = %#v, want attack target 21", action)
	}
}

func TestRulePolicyMovesToMonsterAdjacentPosition(t *testing.T) {
	resources := &resources{
		terrains:          map[int]*terrain{0: newTestTerrain(8, 8)},
		attackableClasses: map[int]struct{}{1: {}},
	}
	policy := newRulePolicy(resources, "account1:char1")
	action := policy.Decide(Snapshot{
		Now:  time.Now(),
		Self: Actor{Index: 10, MapNumber: 0, X: 1, Y: 1, Alive: true},
		Objects: []Actor{
			{Index: 21, Class: 1, MapNumber: 0, X: 5, Y: 1, Alive: true, Attackable: true},
		},
	})
	if action.Kind != ActionMove || len(action.Path) == 0 {
		t.Fatalf("action = %#v, want move", action)
	}
	end := action.Path[len(action.Path)-1]
	if pathDistance(end, Position{X: 5, Y: 1}) != 1 {
		t.Fatalf("path end = %+v, want adjacent to monster", end)
	}
}

func TestRulePolicyMovesTowardNearestSpawnAreaWithoutVisibleMonster(t *testing.T) {
	resources := &resources{
		terrains: map[int]*terrain{0: newTestTerrain(32, 8)},
		spawnAreas: map[int][]spawnArea{
			0: {{class: 1, min: Position{X: 25, Y: 1}, max: Position{X: 30, Y: 5}}},
		},
		attackableClasses: map[int]struct{}{1: {}},
	}
	policy := newRulePolicy(resources, "account1:char1")
	action := policy.Decide(Snapshot{
		Now:  time.Now(),
		Self: Actor{Index: 10, MapNumber: 0, X: 1, Y: 1, Alive: true},
	})
	if action.Kind != ActionMove {
		t.Fatalf("action = %#v, want move", action)
	}
	if len(action.Path) != 15 {
		t.Fatalf("path length = %d, want segmented path length 15", len(action.Path))
	}
}

func TestRulePolicyStableKeyDispersesActorApproach(t *testing.T) {
	resources := &resources{
		terrains:          map[int]*terrain{0: newTestTerrain(12, 12)},
		attackableClasses: map[int]struct{}{1: {}},
	}
	snapshot := Snapshot{
		Now:  time.Now(),
		Self: Actor{Index: 10, MapNumber: 0, X: 1, Y: 5, Alive: true},
		Objects: []Actor{
			{Index: 21, Class: 1, MapNumber: 0, X: 5, Y: 5, Alive: true, Attackable: true},
		},
	}
	first := newRulePolicy(resources, "account1:char1").Decide(snapshot)
	again := newRulePolicy(resources, "account1:char1").Decide(snapshot)
	if first.Kind != ActionMove || again.Kind != ActionMove || len(first.Path) == 0 || len(again.Path) == 0 {
		t.Fatalf("actions = %#v %#v, want move actions", first, again)
	}
	firstEnd := first.Path[len(first.Path)-1]
	againEnd := again.Path[len(again.Path)-1]
	if firstEnd != againEnd {
		t.Fatalf("same key endpoint = %+v then %+v, want stable", firstEnd, againEnd)
	}

	endpoints := map[Position]struct{}{firstEnd: {}}
	for _, key := range []string{
		"account2:char2",
		"account3:char3",
		"account4:char4",
		"account5:char5",
		"account6:char6",
		"account7:char7",
		"account8:char8",
		"account9:char9",
		"account10:char10",
		"account11:char11",
		"account12:char12",
		"account13:char13",
		"account14:char14",
		"account15:char15",
		"account16:char16",
	} {
		action := newRulePolicy(resources, key).Decide(snapshot)
		if action.Kind != ActionMove || len(action.Path) == 0 {
			t.Fatalf("key %q action = %#v, want move", key, action)
		}
		endpoints[action.Path[len(action.Path)-1]] = struct{}{}
	}
	if len(endpoints) < 2 {
		t.Fatalf("endpoints = %#v, want different keys to disperse actor approach", endpoints)
	}
}

func TestRulePolicyStableKeyDispersesSpawnAreaDestination(t *testing.T) {
	resources := &resources{
		terrains: map[int]*terrain{0: newTestTerrain(32, 16)},
		spawnAreas: map[int][]spawnArea{
			0: {{class: 1, min: Position{X: 20, Y: 4}, max: Position{X: 30, Y: 12}}},
		},
		attackableClasses: map[int]struct{}{1: {}},
	}
	snapshot := Snapshot{
		Now:  time.Now(),
		Self: Actor{Index: 10, MapNumber: 0, X: 21, Y: 6, Alive: true},
	}
	first := newRulePolicy(resources, "account1:char1").Decide(snapshot)
	again := newRulePolicy(resources, "account1:char1").Decide(snapshot)
	if first.Kind != ActionMove || again.Kind != ActionMove || len(first.Path) == 0 || len(again.Path) == 0 {
		t.Fatalf("actions = %#v %#v, want move actions", first, again)
	}
	firstEnd := first.Path[len(first.Path)-1]
	againEnd := again.Path[len(again.Path)-1]
	if firstEnd != againEnd {
		t.Fatalf("same key endpoint = %+v then %+v, want stable", firstEnd, againEnd)
	}

	endpoints := map[Position]struct{}{firstEnd: {}}
	for _, key := range []string{
		"account2:char2",
		"account3:char3",
		"account4:char4",
		"account5:char5",
		"account6:char6",
		"account7:char7",
		"account8:char8",
		"account9:char9",
		"account10:char10",
		"account11:char11",
		"account12:char12",
		"account13:char13",
		"account14:char14",
		"account15:char15",
		"account16:char16",
	} {
		action := newRulePolicy(resources, key).Decide(snapshot)
		if action.Kind != ActionMove || len(action.Path) == 0 {
			t.Fatalf("key %q action = %#v, want move", key, action)
		}
		endpoints[action.Path[len(action.Path)-1]] = struct{}{}
	}
	if len(endpoints) < 2 {
		t.Fatalf("endpoints = %#v, want different keys to disperse spawn destinations", endpoints)
	}
}
