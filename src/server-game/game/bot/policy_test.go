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
	policy := newRulePolicy(resources)
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
	policy := newRulePolicy(resources)
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
	policy := newRulePolicy(resources)
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
