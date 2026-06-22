package skill

import (
	"testing"

	"github.com/xujintao/balgass/src/server-game/game/class"
)

func withTestSkillManager(t *testing.T) {
	t.Helper()

	old := SkillManager
	t.Cleanup(func() {
		SkillManager = old
	})

	SkillManager = skillManager{
		skillTable: map[int]*SkillBase{
			SkillIndexFireBall: {Index: SkillIndexFireBall, Damage: 8},
			300:                {Index: 300, STID: 1},
			305:                {Index: 305, STID: 2},
			306:                {Index: 306, STID: 3},
		},
	}
	SkillManager.masterSkillTable[class.Knight][0][0][0] = &MasterSkillBase{
		Index:       1,
		ReqMinPoint: 1,
		MaxPoint:    2,
		SkillID:     300,
	}
	SkillManager.masterSkillTable[class.Knight][0][0][1] = &MasterSkillBase{
		Index:        2,
		ReqMinPoint:  1,
		MaxPoint:     2,
		ParentSkill1: 300,
		SkillID:      305,
	}
	SkillManager.masterSkillTable[class.Knight][0][0][2] = &MasterSkillBase{
		Index:        3,
		ReqMinPoint:  1,
		MaxPoint:     1,
		ParentSkill1: 300,
		ParentSkill2: 305,
		SkillID:      306,
	}
}

func TestGetMasterRequiresParentSkills(t *testing.T) {
	withTestSkillManager(t)
	skills := make(Skills)

	if skills.GetMaster(int(class.Knight), 305, 1, func(int, int, int, int, float32, float32) {}) {
		t.Fatal("GetMaster learned child without parent")
	}
	if !skills.GetMaster(int(class.Knight), 300, 1, func(int, int, int, int, float32, float32) {}) {
		t.Fatal("GetMaster failed to learn parent")
	}
	if skills.GetMaster(int(class.Knight), 306, 1, func(int, int, int, int, float32, float32) {}) {
		t.Fatal("GetMaster learned two-parent skill without second parent")
	}
	if !skills.GetMaster(int(class.Knight), 305, 1, func(int, int, int, int, float32, float32) {}) {
		t.Fatal("GetMaster failed after parent was learned")
	}
}

func TestGetMasterHonorsMaxPoint(t *testing.T) {
	withTestSkillManager(t)
	skills := make(Skills)

	if !skills.GetMaster(int(class.Knight), 300, 1, func(int, int, int, int, float32, float32) {}) {
		t.Fatal("first point failed")
	}
	if !skills.GetMaster(int(class.Knight), 300, 1, func(int, int, int, int, float32, float32) {}) {
		t.Fatal("second point failed")
	}
	if skills.GetMaster(int(class.Knight), 300, 1, func(int, int, int, int, float32, float32) {}) {
		t.Fatal("third point exceeded MaxPoint")
	}
}

func TestForEachMasterSkillSkipsNormalSkills(t *testing.T) {
	withTestSkillManager(t)
	skills := Skills{
		SkillIndexFireBall: {SkillBase: SkillManager.skillTable[SkillIndexFireBall], Index: SkillIndexFireBall},
		300:                {SkillBase: SkillManager.skillTable[300], Index: 300, UIIndex: 1, Level: 1},
		305:                {SkillBase: SkillManager.skillTable[305], Index: 305, UIIndex: 0, Level: 1},
	}

	var got []int
	skills.ForEachMasterSkill(func(index, level int, curValue, nextValue float32) {
		got = append(got, index)
	})

	if len(got) != 1 || got[0] != 1 {
		t.Fatalf("master skill UI indexes = %v, want [1]", got)
	}
}
