package bot

import (
	"errors"
	"reflect"
	"testing"

	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

func TestWorldOwnsRegisteredAPIHandlers(t *testing.T) {
	worldType := reflect.TypeOf((*world)(nil))
	if len(worldAPIs) != len(apis) {
		t.Fatalf("world API count = %d, want %d", len(worldAPIs), len(apis))
	}
	for _, api := range apis {
		if _, ok := worldType.MethodByName(api.handle); !ok {
			t.Fatalf("world handler %q is missing", api.handle)
		}
	}
	if newWorld(&resources{}).Handle(struct{}{}) {
		t.Fatal("Handle(unregistered) = true, want false")
	}
}

func TestWorldPhaseChangesOnlyFromHandledEvents(t *testing.T) {
	world := newWorld(&resources{})
	if world.phase != PhaseDisconnected {
		t.Fatalf("initial phase = %v, want disconnected", world.phase)
	}

	world.Handle(&model.MsgConnectReply{Result: 1, ID: 7})
	if world.phase != PhaseConnected || world.self.Index != 7 {
		t.Fatalf("connect reply state = phase %v index %d", world.phase, world.self.Index)
	}

	world.Handle(&model.MsgLoginReply{Result: 1})
	if world.phase != PhaseLoggedIn {
		t.Fatalf("login reply phase = %v, want logged in", world.phase)
	}

	world.Handle(&model.MsgLoadCharacterReply{MapNumber: 0, X: 10, Y: 20, HP: 100})
	if world.phase != PhasePlaying || world.self.X != 10 || world.self.Y != 20 {
		t.Fatalf("load reply state = phase %v self %#v", world.phase, world.self)
	}
}

func TestWorldHandleFailureSetsFailedPhase(t *testing.T) {
	world := newWorld(&resources{})
	world.Handle(&connectFailed{Err: errors.New("test connect failure")})
	if world.phase != PhaseFailed || world.failure == "" {
		t.Fatalf("connect failure state = phase %v reason %q", world.phase, world.failure)
	}
}

func TestWorldMoveReplyOnlyRecordsServerTarget(t *testing.T) {
	world := newWorld(&resources{})
	world.setSelfIndex(7)
	world.self.Alive = true
	world.self.X, world.self.Y = 1, 1
	version := world.positionVersion

	world.HandleMoveReply(&model.MsgMoveReply{Number: 7, X: 3, Y: 1, Dir: 3 << 4})

	if world.self.X != 1 || world.self.Y != 1 {
		t.Fatalf("self position = (%d,%d), want (1,1)", world.self.X, world.self.Y)
	}
	if world.self.TX != 3 || world.self.TY != 1 || world.self.Dir != 3 {
		t.Fatalf("move target = (%d,%d) dir=%d, want (3,1) dir=3", world.self.TX, world.self.TY, world.self.Dir)
	}
	if world.positionVersion != version {
		t.Fatalf("position version = %d, want %d", world.positionVersion, version)
	}
}

func TestWorldAuthoritativePositionIncrementsVersion(t *testing.T) {
	world := newWorld(&resources{})
	world.setSelfIndex(7)
	before := world.positionVersion

	world.HandleSetPositionReply(&model.MsgSetPositionReply{Number: 7, X: 4, Y: 5})

	if world.positionVersion != before+1 {
		t.Fatalf("position version = %d, want %d", world.positionVersion, before+1)
	}
}

func TestWorldSnapshotIncludesAttackAndSkillTimingState(t *testing.T) {
	world := newWorld(&resources{})
	world.Handle(&model.MsgAttackSpeedReply{AttackSpeed: 100, MagicSpeed: 50})
	world.skills[skill.SkillIndexFireBall] = &skill.Skill{
		SkillBase: &skill.SkillBase{
			Index:    skill.SkillIndexFireBall,
			Type:     skillTypeMagic,
			Delay:    700,
			Distance: 4,
		},
		Index: skill.SkillIndexFireBall,
	}

	snapshot := world.Snapshot()
	if snapshot.AttackSpeed != 100 || snapshot.MagicSpeed != 50 {
		t.Fatalf("speed snapshot = %d/%d, want 100/50", snapshot.AttackSpeed, snapshot.MagicSpeed)
	}
	combatSkill, ok := snapshot.combatSkill(skill.SkillIndexFireBall)
	if !ok || combatSkill.Type != skillTypeMagic || combatSkill.Delay != 700 {
		t.Fatalf("combat skill = %#v, want magic delay 700", combatSkill)
	}
}

func TestWorldSnapshotRejectsUnsupportedSkillType(t *testing.T) {
	world := newWorld(&resources{})
	world.skills[skill.SkillIndexPowerSlash] = &skill.Skill{
		SkillBase: &skill.SkillBase{
			Index:    skill.SkillIndexPowerSlash,
			Type:     2,
			Distance: 5,
		},
		Index: skill.SkillIndexPowerSlash,
	}

	snapshot := world.Snapshot()
	if len(snapshot.Skills) != 0 {
		t.Fatal("unsupported skill type was included in combat skills")
	}
	if !snapshot.learnedSkill(skill.SkillIndexPowerSlash) {
		t.Fatal("learned skill was omitted from learned skill state")
	}
}

func TestWorldLearnSkillsWaitForCharacterClass(t *testing.T) {
	world := newWorld(&resources{})
	world.setSelfIndex(7)
	book := testSkillItem(skill.SkillIndexFireBall)
	world.HandleItemListReply(&model.MsgItemListReply{Items: []*item.Item{book}})
	world.HandleSkillListReply(&model.MsgSkillListReply{})
	world.setCharacter(&model.MsgLoadCharacterReply{
		HP:       100,
		Strength: 10,
	})

	if got := world.Snapshot().LearnSkills; len(got) != 0 {
		t.Fatalf("learn skills before self viewport = %#v, want none", got)
	}

	world.HandleCreateViewportPlayerReply(&model.MsgCreateViewportPlayerReply{
		Players: []*model.CreateViewportPlayer{{
			Index: 7,
			Class: int(class.Wizard),
			Level: 1,
			HP:    100,
		}},
	})

	assertLearnSkill(t, world.Snapshot(), skill.SkillIndexFireBall)
}

func TestWorldLearnSkillsHonorClassAndChangeUp(t *testing.T) {
	book := testSkillItem(skill.SkillIndexFireBall)
	world := readySkillWorld(book)

	assertLearnSkill(t, world.Snapshot(), skill.SkillIndexFireBall)

	world.self.Class = int(class.Knight)
	if got := world.Snapshot().LearnSkills; len(got) != 0 {
		t.Fatalf("learn skills for disallowed class = %#v, want none", got)
	}

	world.self.Class = int(class.Wizard)
	book.ReqClass[class.Wizard] = 2
	if got := world.Snapshot().LearnSkills; len(got) != 0 {
		t.Fatalf("learn skills before required change-up = %#v, want none", got)
	}

	world.self.ChangeUp = 1
	assertLearnSkill(t, world.Snapshot(), skill.SkillIndexFireBall)
}

func TestWorldLearnSkillsHonorBaseRequirements(t *testing.T) {
	book := testSkillItem(skill.SkillIndexFireBall)
	book.ReqLevel = 100
	book.ReqStrength = 20
	book.ReqDexterity = 21
	book.ReqVitality = 22
	book.ReqEnergy = 23
	book.ReqCommand = 24
	world := readySkillWorld(book)
	world.self.Level = 100
	world.strength = 20
	world.dexterity = 21
	world.vitality = 22
	world.energy = 23
	world.leadership = 24

	assertLearnSkill(t, world.Snapshot(), skill.SkillIndexFireBall)

	tests := []struct {
		name string
		set  func()
	}{
		{"level", func() { world.self.Level-- }},
		{"strength", func() { world.strength-- }},
		{"dexterity", func() { world.dexterity-- }},
		{"vitality", func() { world.vitality-- }},
		{"energy", func() { world.energy-- }},
		{"leadership", func() { world.leadership-- }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.set()
			if got := world.Snapshot().LearnSkills; len(got) != 0 {
				t.Fatalf("learn skills below %s requirement = %#v, want none", tt.name, got)
			}
			world.self.Level = 100
			world.strength = 20
			world.dexterity = 21
			world.vitality = 22
			world.energy = 23
			world.leadership = 24
		})
	}
}

func TestWorldSuccessfulLevelPointEnablesSkill(t *testing.T) {
	book := testSkillItem(skill.SkillIndexFireBall)
	book.ReqEnergy = 11
	world := readySkillWorld(book)
	world.energy = 10

	world.HandleAddLevelPointReply(&model.MsgAddLevelPointReply{})
	if got := world.Snapshot().LearnSkills; len(got) != 0 {
		t.Fatalf("learn skills after failed level point = %#v, want none", got)
	}

	world.HandleAddLevelPointReply(&model.MsgAddLevelPointReply{Type: 0x13})
	assertLearnSkill(t, world.Snapshot(), skill.SkillIndexFireBall)
}

func TestSkillIndexFromSummoningOrbLevel(t *testing.T) {
	book := testSkillItem(10)
	book.Code = item.Code(12, 11)
	book.Level = 3

	index, ok := skillIndexFromItem(book)
	if !ok || index != 13 {
		t.Fatalf("skill index = %d, %v, want 13, true", index, ok)
	}
}

func TestWorldReloadClearsSightAndRevivesSelf(t *testing.T) {
	world := newWorld(&resources{})
	world.setSelfIndex(7)
	world.self.Alive = false
	world.objects[1] = Actor{Index: 1, Alive: true}

	world.HandleReloadCharacterReply(&model.MsgReloadCharacterReply{
		MapNumber: 0,
		X:         10,
		Y:         20,
		Dir:       3,
		HP:        100,
	})

	if !world.self.Alive || world.self.X != 10 || world.self.Y != 20 {
		t.Fatalf("self = %#v, want revived at (10,20)", world.self)
	}
	if len(world.objects) != 0 {
		t.Fatalf("object count = %d, want 0", len(world.objects))
	}
}

func testSkillItem(index int) *item.Item {
	base := &item.ItemBase{
		KindA:      item.KindASkill,
		SkillIndex: index,
	}
	base.ReqClass[class.Wizard] = 1
	return &item.Item{
		ItemBase: base,
		Code:     item.Code(15, 0),
	}
}

func readySkillWorld(book *item.Item) *world {
	world := newWorld(&resources{})
	world.self = Actor{
		Class: int(class.Wizard),
		Level: 1,
		Alive: true,
	}
	world.inventory = []*item.Item{book}
	world.itemsSet = true
	world.skillsSet = true
	world.characterSet = true
	world.selfClassSet = true
	return world
}

func assertLearnSkill(t *testing.T, snapshot WorldSnapshot, index int) {
	t.Helper()
	if len(snapshot.LearnSkills) != 1 ||
		snapshot.LearnSkills[0].Index != index ||
		snapshot.LearnSkills[0].Position != 0 {
		t.Fatalf("learn skills = %#v, want index %d at position 0", snapshot.LearnSkills, index)
	}
}
