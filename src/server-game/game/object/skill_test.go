package object

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

type skillTestActor struct {
	messages []any
}

func (*skillTestActor) Addr() string                              { return "test" }
func (*skillTestActor) Offline()                                  {}
func (a *skillTestActor) Push(msg any)                            { a.messages = append(a.messages, msg) }
func (*skillTestActor) ProcessAction()                            {}
func (*skillTestActor) Process1000ms()                            {}
func (*skillTestActor) SpawnPosition()                            {}
func (*skillTestActor) Die(*Object, int)                          {}
func (*skillTestActor) LevelUp(int) bool                          { return false }
func (*skillTestActor) DieDropItem(*Object)                       {}
func (*skillTestActor) Regen()                                    {}
func (*skillTestActor) EquipmentChanged()                         {}
func (*skillTestActor) GetPKLevel() int                           { return 0 }
func (*skillTestActor) GetMasterLevel() int                       { return 0 }
func (*skillTestActor) IsMasterLevel() bool                       { return false }
func (*skillTestActor) GetSkillMPAG(s *skill.Skill) (int, int)    { return s.ManaUsage, s.BPUsage }
func (*skillTestActor) GetChangeUp() int                          { return 0 }
func (*skillTestActor) CanUseItem(*item.Item) bool                { return true }
func (*skillTestActor) GetInventory() *item.Inventory             { return nil }
func (*skillTestActor) GetInventoryItem(int) *item.Item           { return nil }
func (*skillTestActor) GetWarehouse() *item.Warehouse             { return nil }
func (*skillTestActor) SetDelayRecoverHP(int, int)                {}
func (*skillTestActor) SetDelayRecoverSD(int, int)                {}
func (*skillTestActor) GetAttackRatePVP() int                     { return 1000 }
func (*skillTestActor) GetDefenseRatePVP() int                    { return 1 }
func (*skillTestActor) GetIgnoreDefenseRate() int                 { return 0 }
func (*skillTestActor) GetCriticalAttackRate() int                { return 0 }
func (*skillTestActor) GetCriticalAttackDamage() int              { return 0 }
func (*skillTestActor) GetExcellentAttackRate() int               { return 0 }
func (*skillTestActor) GetExcellentAttackDamage() int             { return 0 }
func (*skillTestActor) GetMonsterDieGetHP() float64               { return 0 }
func (*skillTestActor) GetMonsterDieGetMP() float64               { return 0 }
func (*skillTestActor) GetAddDamage() int                         { return 0 }
func (*skillTestActor) GetArmorReduceDamage() int                 { return 0 }
func (*skillTestActor) GetWingIncreaseDamage() int                { return 0 }
func (*skillTestActor) GetWingReduceDamage() int                  { return 0 }
func (*skillTestActor) GetHelperReduceDamage() int                { return 0 }
func (*skillTestActor) GetPetIncreaseDamage() int                 { return 0 }
func (*skillTestActor) GetPetReduceDamage() int                   { return 0 }
func (*skillTestActor) GetDoubleDamageRate() int                  { return 0 }
func (*skillTestActor) GetMonsterDieGetMoney() float64            { return 0 }
func (*skillTestActor) GetKnightGladiatorCalcSkillBonus() float64 { return 1 }
func (*skillTestActor) GetImpaleSkillCalc() float64               { return 1 }
func (*skillTestActor) PickItem(*model.MsgPickItem)               {}
func (*skillTestActor) DropItem(*model.MsgDropItem)               {}
func (*skillTestActor) BuyItem(*model.MsgBuyItem)                 {}
func (*skillTestActor) SellItem(*model.MsgSellItem)               {}
func (*skillTestActor) MoveItem(*model.MsgMoveItem)               {}
func (*skillTestActor) UseItem(*model.MsgUseItem)                 {}
func (*skillTestActor) RepairItem(*model.MsgRepairItem)           {}
func (*skillTestActor) Move(*model.MsgMove)                       {}
func (*skillTestActor) Teleport(*model.MsgTeleport)               {}
func (*skillTestActor) MapMove(*model.MsgMapMove)                 {}
func (*skillTestActor) SetPosition(*model.MsgSetPosition)         {}
func (*skillTestActor) Action(*model.MsgAction)                   {}
func (*skillTestActor) UseSkill(*model.MsgUseSkill)               {}
func (*skillTestActor) Attack(*model.MsgAttack)                   {}
func (*skillTestActor) Chat(*model.MsgChat)                       {}
func (*skillTestActor) Whisper(*model.MsgWhisper)                 {}
func (*skillTestActor) Login(*model.MsgLogin)                     {}
func (*skillTestActor) Logout(*model.MsgLogout)                   {}
func (*skillTestActor) GetCharacterList(*model.MsgEmpty)          {}
func (*skillTestActor) CreateCharacter(*model.MsgCreateCharacter) {}
func (*skillTestActor) DeleteCharacter(*model.MsgDeleteCharacter) {}
func (*skillTestActor) CheckCharacter(*model.MsgCheckCharacter)   {}
func (*skillTestActor) LoadCharacter(*model.MsgLoadCharacter)     {}
func (*skillTestActor) Talk(*model.MsgTalk)                       {}
func (*skillTestActor) CloseTalkWindow(*model.MsgEmpty)           {}
func (*skillTestActor) CloseWarehouseWindow(*model.MsgEmpty)      {}
func (*skillTestActor) KeepLive(*model.MsgKeepLive)               {}
func (*skillTestActor) Hack(*model.MsgHack)                       {}
func (*skillTestActor) BattleCoreNotice(*model.MsgEmpty)          {}
func (*skillTestActor) MapDataLoadingOK(*model.MsgEmpty)          {}
func (*skillTestActor) AddLevelPoint(*model.MsgAddLevelPoint)     {}
func (*skillTestActor) LearnMasterSkill(*model.MsgLearnMasterSkill) {
}
func (*skillTestActor) DefineMuKey(*model.MsgDefineMuKey)        {}
func (*skillTestActor) DefineMuBot(*model.MsgDefineMuBot)        {}
func (*skillTestActor) EnableMuBot(*model.MsgEnableMuBot)        {}
func (*skillTestActor) UsePet(*model.MsgUsePet)                  {}
func (*skillTestActor) MuunSystem(*model.MsgMuunSystem)          {}
func (*skillTestActor) StartPartyNumberPosition(*model.MsgEmpty) {}
func (*skillTestActor) StopPartyNumberPosition(*model.MsgEmpty)  {}

func newSkillTestObject(index int, typ ObjectType) (*Object, *skillTestActor) {
	actor := &skillTestActor{}
	obj := &Object{
		Objecter:     actor,
		Index:        index,
		Type:         typ,
		MapNumber:    0,
		X:            10,
		Y:            10,
		Level:        100,
		ConnectState: ConnectStatePlaying,
		Live:         true,
		State:        2,
		HP:           100,
		MaxHP:        100,
		MP:           100,
		MaxMP:        100,
		AG:           100,
		MaxAG:        100,
		AttackMin:    30,
		AttackMax:    30,
		AttackRate:   1000,
		DefenseRate:  1,
	}
	obj.Init()
	return obj, actor
}

func learnSkillForTest(t *testing.T, obj *Object, index int) *skill.Skill {
	t.Helper()
	s, ok := obj.LearnSkill(index)
	if !ok {
		t.Fatalf("LearnSkill(%d) = false", index)
	}
	return s
}

func assertResourceUnchanged(t *testing.T, obj *Object, mp, ag int) {
	t.Helper()
	if obj.MP != mp || obj.AG != ag {
		t.Fatalf("resources = %d/%d, want %d/%d", obj.MP, obj.AG, mp, ag)
	}
}

func hasMessage[T any](messages []any) bool {
	for _, msg := range messages {
		if _, ok := msg.(*T); ok {
			return true
		}
	}
	return false
}

func TestUseSkillRejectsInvalidRequestsWithoutResourceCost(t *testing.T) {
	for _, tt := range []struct {
		name  string
		setup func(caster, target *Object)
		msg   model.MsgUseSkill
	}{
		{
			name: "unlearned skill",
			msg:  model.MsgUseSkill{Target: 2, Skill: skill.SkillIndexFireBall},
		},
		{
			name: "missing target",
			setup: func(caster, target *Object) {
				learnSkillForTest(t, caster, skill.SkillIndexFireBall)
			},
			msg: model.MsgUseSkill{Target: 99, Skill: skill.SkillIndexFireBall},
		},
		{
			name: "dead target",
			setup: func(caster, target *Object) {
				learnSkillForTest(t, caster, skill.SkillIndexFireBall)
				target.Live = false
			},
			msg: model.MsgUseSkill{Target: 2, Skill: skill.SkillIndexFireBall},
		},
		{
			name: "cross map",
			setup: func(caster, target *Object) {
				learnSkillForTest(t, caster, skill.SkillIndexFireBall)
				target.MapNumber = caster.MapNumber + 1
			},
			msg: model.MsgUseSkill{Target: 2, Skill: skill.SkillIndexFireBall},
		},
		{
			name: "out of distance",
			setup: func(caster, target *Object) {
				learnSkillForTest(t, caster, skill.SkillIndexFireBall)
				target.X = caster.X + 20
			},
			msg: model.MsgUseSkill{Target: 2, Skill: skill.SkillIndexFireBall},
		},
		{
			name: "insufficient resources",
			setup: func(caster, target *Object) {
				learnSkillForTest(t, caster, skill.SkillIndexFireBall)
				caster.MP = 0
			},
			msg: model.MsgUseSkill{Target: 2, Skill: skill.SkillIndexFireBall},
		},
		{
			name: "delay not elapsed",
			setup: func(caster, target *Object) {
				s := learnSkillForTest(t, caster, skill.SkillIndexFireBall)
				s.Delay = 1000
				caster.skillUseTimes[s.Index] = time.Now()
			},
			msg: model.MsgUseSkill{Target: 2, Skill: skill.SkillIndexFireBall},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			caster, actor := newSkillTestObject(1, ObjectTypePlayer)
			target, _ := newSkillTestObject(2, ObjectTypeMonster)
			withTestObjectManager(t, caster, target)
			if tt.setup != nil {
				tt.setup(caster, target)
			}
			mp, ag := caster.MP, caster.AG

			caster.UseSkill(&tt.msg)

			assertResourceUnchanged(t, caster, mp, ag)
			if hasMessage[model.MsgMPReply](actor.messages) {
				t.Fatal("resource reply was sent for rejected skill use")
			}
		})
	}
}

func TestUseSkillUnknownImplementedSkillDoesNotCostResources(t *testing.T) {
	caster, actor := newSkillTestObject(1, ObjectTypePlayer)
	target, _ := newSkillTestObject(2, ObjectTypeMonster)
	withTestObjectManager(t, caster, target)
	learnSkillForTest(t, caster, skill.SkillIndexHellFire)
	mp, ag := caster.MP, caster.AG

	caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: skill.SkillIndexHellFire})

	assertResourceUnchanged(t, caster, mp, ag)
	if hasMessage[model.MsgMPReply](actor.messages) {
		t.Fatal("resource reply was sent for unimplemented skill")
	}
}

func TestUseSkillExpandedAttackSkills(t *testing.T) {
	for _, index := range []int{
		skill.SkillIndexTwister,
		skill.SkillIndexEvilSpirit,
		skill.SkillIndexPowerWave,
		skill.SkillIndexAquaBeam,
		skill.SkillIndexCometFall,
		skill.SkillIndexDecay,
		skill.SkillIndexIceStorm,
		skill.SkillIndexIceArrow,
		skill.SkillIndexFireSlash,
		skill.SkillIndexForce,
		skill.SkillIndexFireBurst,
		skill.SkillIndexElectricSpike,
		skill.SkillIndexForceWave,
	} {
		t.Run(fmt.Sprintf("skill_%d", index), func(t *testing.T) {
			rand.Seed(1)
			caster, actor := newSkillTestObject(1, ObjectTypePlayer)
			target, _ := newSkillTestObject(2, ObjectTypeMonster)
			withTestObjectManager(t, caster, target)
			caster.MP = 1000
			caster.MaxMP = 1000
			caster.AG = 1000
			caster.MaxAG = 1000
			caster.AttackRate = 1000000
			target.HP = 10000
			target.MaxHP = 10000
			s := learnSkillForTest(t, caster, index)
			targetHP := target.HP

			caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: index})

			if caster.MP != 1000-s.ManaUsage || caster.AG != 1000-s.BPUsage {
				t.Fatalf("resources = %d/%d, want %d/%d",
					caster.MP, caster.AG, 1000-s.ManaUsage, 1000-s.BPUsage)
			}
			if target.HP >= targetHP {
				t.Fatalf("target HP = %d, want below %d", target.HP, targetHP)
			}
			if !hasMessage[model.MsgUseSkillReply](actor.messages) {
				t.Fatal("skill success reply was not sent")
			}
			if !hasMessage[model.MsgMPReply](actor.messages) {
				t.Fatal("resource reply was not sent")
			}
		})
	}
}

func TestExpandedSkillDamageSource(t *testing.T) {
	for _, tt := range []struct {
		name       string
		index      int
		wantDamage int
	}{
		{name: "table damage", index: skill.SkillIndexEvilSpirit, wantDamage: 70},
		{name: "ice arrow physical damage", index: skill.SkillIndexIceArrow, wantDamage: 30},
		{name: "fire slash physical damage", index: skill.SkillIndexFireSlash, wantDamage: 30},
		{name: "dark lord physical damage", index: skill.SkillIndexFireBurst, wantDamage: 30},
	} {
		t.Run(tt.name, func(t *testing.T) {
			caster, _ := newSkillTestObject(1, ObjectTypePlayer)
			caster.AttackMin = 30
			caster.AttackMax = 30
			s := learnSkillForTest(t, caster, tt.index)
			s.DamageMin = 70
			s.DamageMax = 70

			if damage := caster.getDamage(s, 0); damage != tt.wantDamage {
				t.Fatalf("damage = %d, want %d", damage, tt.wantDamage)
			}
		})
	}
}

func TestUseSkillAttackSuccessCostsResourcesAndDamagesTarget(t *testing.T) {
	caster, actor := newSkillTestObject(1, ObjectTypePlayer)
	target, _ := newSkillTestObject(2, ObjectTypeMonster)
	withTestObjectManager(t, caster, target)
	s := learnSkillForTest(t, caster, skill.SkillIndexFireBall)
	targetHP := target.HP

	caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: skill.SkillIndexFireBall})

	if caster.MP != 100-s.ManaUsage || caster.AG != 100-s.BPUsage {
		t.Fatalf("resources = %d/%d, want %d/%d", caster.MP, caster.AG, 100-s.ManaUsage, 100-s.BPUsage)
	}
	if target.HP >= targetHP {
		t.Fatalf("target HP = %d, want below %d", target.HP, targetHP)
	}
	if !hasMessage[model.MsgUseSkillReply](actor.messages) {
		t.Fatal("skill success reply was not sent")
	}
	if !hasMessage[model.MsgMPReply](actor.messages) {
		t.Fatal("resource reply was not sent")
	}
}
