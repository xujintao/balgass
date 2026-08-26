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
		masterSkillValueTable: make(map[int]masterSkillValue),
	}
	value := masterSkillValue{valueType: valueTypeNormal}
	value.values[1] = 1.25
	value.values[2] = 2.5
	SkillManager.masterSkillValueTable[1] = value
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

func TestMasterSkillValue(t *testing.T) {
	for _, tt := range []struct {
		name  string
		index int
		level int
		want  float32
	}{
		{name: "STID 1", index: SkillIndexMaceMastery, level: 10, want: 6.68},
		{name: "STID 38", index: SkillIndexRecoverManaFully, level: 10, want: 2.64},
		{name: "zero level", index: SkillIndexMaceMastery},
		{name: "level above maximum", index: SkillIndexMaceMastery, level: 21},
		{name: "STID zero", index: 418, level: 10},
		{name: "unknown skill", index: 9999, level: 10},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := SkillManager.GetMasterSkillValue(tt.index, tt.level); got != tt.want {
				t.Fatalf("GetMasterSkillValue(%d, %d) = %.2f, want %.2f", tt.index, tt.level, got, tt.want)
			}
		})
	}
}

func TestMasterSkillCurrentAndNextValues(t *testing.T) {
	for _, tt := range []struct {
		name string
		load func(*testing.T, Skills)
	}{
		{
			name: "loaded skill",
			load: func(_ *testing.T, skills Skills) {
				skills[300] = &Skill{SkillBase: SkillManager.skillTable[300], Index: 300, Level: 1}
				skills.FillSkillData(int(class.Knight))
			},
		},
		{
			name: "new skill",
			load: func(t *testing.T, skills Skills) {
				if !skills.GetMaster(int(class.Knight), 300, 1, func(int, int, int, int, float32, float32) {}) {
					t.Fatal("GetMaster failed")
				}
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			withTestSkillManager(t)
			skills := make(Skills)
			tt.load(t, skills)
			s := skills[300]
			if s.CurValue != 1.25 || s.NextValue != 2.5 {
				t.Fatalf("master values = %.2f/%.2f, want 1.25/2.50", s.CurValue, s.NextValue)
			}
		})
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
