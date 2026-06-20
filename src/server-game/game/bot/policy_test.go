package bot

import (
	"testing"
	"time"
)

func decide(policy Policy, snapshot WorldSnapshot) Action {
	return policy.Decide(time.Time{}, snapshot, ExecutorSnapshot{
		Position:        snapshot.Self.position(),
		Dir:             snapshot.Self.Dir,
		PositionVersion: snapshot.PositionVersion,
	})
}

func TestRulePolicyChoosesLifecycleActions(t *testing.T) {
	policy := newRulePolicy(&resources{}, "account1:char1")
	tests := []struct {
		phase Phase
		want  ActionKind
	}{
		{phase: PhaseDisconnected, want: ActionConnect},
		{phase: PhaseConnected, want: ActionLogin},
		{phase: PhaseLoggedIn, want: ActionLoadCharacter},
		{phase: PhaseFailed, want: ActionStop},
	}
	for _, tt := range tests {
		action := decide(policy, WorldSnapshot{Phase: tt.phase, Failure: "failed"})
		if action.Kind != tt.want {
			t.Fatalf("phase %v action = %v, want %v", tt.phase, action.Kind, tt.want)
		}
	}
}

func TestRulePolicyUsesExecutorPendingAndMovementState(t *testing.T) {
	policy := newRulePolicy(&resources{}, "account1:char1")
	if action := policy.Decide(
		time.Time{},
		WorldSnapshot{Phase: PhaseDisconnected},
		ExecutorSnapshot{CurrentAction: Action{Kind: ActionConnect}},
	); action.Kind != ActionNone {
		t.Fatalf("pending connect action = %#v, want none", action)
	}
}

func TestRulePolicyExplicitlySyncsAuthoritativePosition(t *testing.T) {
	policy := newRulePolicy(&resources{}, "account1:char1")
	world := WorldSnapshot{
		Phase:           PhasePlaying,
		Self:            Actor{X: 4, Y: 5, Dir: 3, Alive: true},
		PositionVersion: 2,
	}
	execution := ExecutorSnapshot{
		CurrentAction:   Action{Kind: ActionMove},
		Position:        Position{X: 1, Y: 1},
		PositionVersion: 1,
		Move:            MoveSnapshot{Active: true},
	}

	action := policy.Decide(time.Time{}, world, execution)
	if action.Kind != ActionSyncPosition ||
		action.SelfPosition != (Position{X: 4, Y: 5}) ||
		action.PositionVersion != 2 ||
		!action.CancelCurrent {
		t.Fatalf("sync action = %#v", action)
	}
}

func TestRulePolicyKeepsLearnSkillPendingUntilWorldReply(t *testing.T) {
	policy := newRulePolicy(&resources{}, "account1:char1")
	world := WorldSnapshot{
		Phase: PhasePlaying,
		Self:  Actor{Alive: true},
	}
	execution := ExecutorSnapshot{
		CurrentAction: Action{Kind: ActionLearnSkill, Skill: 3},
	}
	if action := policy.Decide(time.Time{}, world, execution); action.Kind != ActionNone {
		t.Fatalf("pending learn action = %#v, want none", action)
	}

	world.LearnedSkills = []int{3}
	if action := policy.Decide(time.Time{}, world, execution); action.Kind == ActionLearnSkill {
		t.Fatalf("learned skill action = %#v, want pending released", action)
	}
}

func TestRulePolicyPrefersCombatThenLearningThenRoaming(t *testing.T) {
	resources := &resources{
		terrains: map[int]*terrain{0: newTestTerrain(32, 8)},
		spawnAreas: map[int][]spawnArea{
			0: {{class: 1, min: Position{X: 25, Y: 1}, max: Position{X: 30, Y: 5}}},
		},
		attackableClasses: map[int]struct{}{1: {}},
	}
	policy := newRulePolicy(resources, "account1:char1")
	snapshot := WorldSnapshot{
		Phase:       PhasePlaying,
		Self:        Actor{Index: 10, MapNumber: 0, X: 1, Y: 1, Alive: true},
		LearnSkills: []LearnSkill{{Index: 3, Position: 12}},
		Objects: []Actor{
			{Index: 21, Class: 1, MapNumber: 0, X: 2, Y: 1, Alive: true, Attackable: true},
		},
	}
	if action := decide(policy, snapshot); action.Kind != ActionAttack {
		t.Fatalf("action with target = %#v, want attack", action)
	}

	snapshot.Objects = nil
	action := decide(policy, snapshot)
	if action.Kind != ActionLearnSkill || action.Skill != 3 || action.Position != 12 {
		t.Fatalf("action without target = %#v, want learn skill", action)
	}

	snapshot.LearnSkills = nil
	if action := decide(policy, snapshot); action.Kind != ActionMove {
		t.Fatalf("action without skill = %#v, want roam", action)
	}
}

func TestRulePolicyAttacksAdjacentMonsterAndIgnoresNPC(t *testing.T) {
	resources := &resources{
		terrains:          map[int]*terrain{0: newTestTerrain(8, 8)},
		attackableClasses: map[int]struct{}{1: {}},
	}
	policy := newRulePolicy(resources, "account1:char1")
	action := decide(policy, WorldSnapshot{
		Phase: PhasePlaying,
		Self:  Actor{Index: 10, MapNumber: 0, X: 1, Y: 1, Alive: true},
		Objects: []Actor{
			{Index: 20, Class: 100, MapNumber: 0, X: 1, Y: 2, Alive: true},
			{Index: 21, Class: 1, MapNumber: 0, X: 2, Y: 1, Alive: true, Attackable: true},
		},
	})
	if action.Kind != ActionAttack || action.Target != 21 {
		t.Fatalf("action = %#v, want attack target 21", action)
	}
}

func TestRulePolicyUsesBestAffordableSkillAtRange(t *testing.T) {
	resources := &resources{
		terrains:          map[int]*terrain{0: newTestTerrain(8, 8)},
		attackableClasses: map[int]struct{}{1: {}},
	}
	policy := newRulePolicy(resources, "account1:char1")
	action := decide(policy, WorldSnapshot{
		Phase: PhasePlaying,
		Self:  Actor{Index: 10, MapNumber: 0, X: 1, Y: 1, Alive: true},
		MP:    20,
		AG:    5,
		Skills: []CombatSkill{
			{Index: 2, Damage: 30, Distance: 3, MP: 10},
			{Index: 4, Damage: 60, Distance: 5, MP: 30},
			{Index: 1, Damage: 30, Distance: 3, MP: 10},
		},
		Objects: []Actor{
			{Index: 21, Class: 1, MapNumber: 0, X: 4, Y: 1, Alive: true, Attackable: true},
		},
	})
	if action.Kind != ActionUseSkill || action.Target != 21 || action.Skill != 1 {
		t.Fatalf("action = %#v, want skill 1 against target 21", action)
	}
}

func TestRulePolicyFallsBackToNormalAttackWithoutSkillResources(t *testing.T) {
	resources := &resources{
		terrains:          map[int]*terrain{0: newTestTerrain(8, 8)},
		attackableClasses: map[int]struct{}{1: {}},
	}
	policy := newRulePolicy(resources, "account1:char1")
	action := decide(policy, WorldSnapshot{
		Phase:  PhasePlaying,
		Self:   Actor{Index: 10, MapNumber: 0, X: 1, Y: 1, Alive: true},
		Skills: []CombatSkill{{Index: 4, Damage: 60, Distance: 5, MP: 30}},
		Objects: []Actor{
			{Index: 21, Class: 1, MapNumber: 0, X: 2, Y: 1, Alive: true, Attackable: true},
		},
	})
	if action.Kind != ActionAttack || action.Target != 21 {
		t.Fatalf("action = %#v, want normal attack target 21", action)
	}
}

func TestRulePolicyMovesToMonsterAdjacentPosition(t *testing.T) {
	resources := &resources{
		terrains:          map[int]*terrain{0: newTestTerrain(8, 8)},
		attackableClasses: map[int]struct{}{1: {}},
	}
	policy := newRulePolicy(resources, "account1:char1")
	action := decide(policy, WorldSnapshot{
		Phase: PhasePlaying,
		Self:  Actor{Index: 10, MapNumber: 0, X: 1, Y: 1, Alive: true},
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
	action := decide(policy, WorldSnapshot{
		Phase: PhasePlaying,
		Self:  Actor{Index: 10, MapNumber: 0, X: 1, Y: 1, Alive: true},
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
	snapshot := WorldSnapshot{
		Phase: PhasePlaying,
		Self:  Actor{Index: 10, MapNumber: 0, X: 1, Y: 5, Alive: true},
		Objects: []Actor{
			{Index: 21, Class: 1, MapNumber: 0, X: 5, Y: 5, Alive: true, Attackable: true},
		},
	}
	first := decide(newRulePolicy(resources, "account1:char1"), snapshot)
	again := decide(newRulePolicy(resources, "account1:char1"), snapshot)
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
		action := decide(newRulePolicy(resources, key), snapshot)
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
	snapshot := WorldSnapshot{
		Phase: PhasePlaying,
		Self:  Actor{Index: 10, MapNumber: 0, X: 21, Y: 6, Alive: true},
	}
	first := decide(newRulePolicy(resources, "account1:char1"), snapshot)
	again := decide(newRulePolicy(resources, "account1:char1"), snapshot)
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
		action := decide(newRulePolicy(resources, key), snapshot)
		if action.Kind != ActionMove || len(action.Path) == 0 {
			t.Fatalf("key %q action = %#v, want move", key, action)
		}
		endpoints[action.Path[len(action.Path)-1]] = struct{}{}
	}
	if len(endpoints) < 2 {
		t.Fatalf("endpoints = %#v, want different keys to disperse spawn destinations", endpoints)
	}
}
