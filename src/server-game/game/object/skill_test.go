package object

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

type skillTestActor struct {
	messages       []any
	magicAttackMin int
	magicAttackMax int
	strength       int
	dexterity      int
	energy         int
	vitality       int
	leadership     int
	knightRate     float64
	impaleRate     float64
}

func (*skillTestActor) Addr() string                           { return "test" }
func (*skillTestActor) Offline()                               {}
func (a *skillTestActor) Push(msg any)                         { a.messages = append(a.messages, msg) }
func (*skillTestActor) ProcessAction()                         {}
func (*skillTestActor) Process1000ms()                         {}
func (*skillTestActor) SpawnPosition()                         {}
func (*skillTestActor) Die(*Object, int)                       {}
func (*skillTestActor) LevelUp(int) bool                       { return false }
func (*skillTestActor) DieDropItem(*Object)                    {}
func (*skillTestActor) Regen()                                 {}
func (*skillTestActor) EquipmentChanged()                      {}
func (*skillTestActor) GetPKLevel() int                        { return 0 }
func (*skillTestActor) GetMasterLevel() int                    { return 0 }
func (*skillTestActor) IsMasterLevel() bool                    { return false }
func (*skillTestActor) GetSkillMPAG(s *skill.Skill) (int, int) { return s.ManaUsage, s.BPUsage }
func (a *skillTestActor) GetMagicAttackMin() int               { return a.magicAttackMin }
func (a *skillTestActor) GetMagicAttackMax() int               { return a.magicAttackMax }
func (a *skillTestActor) GetStrength() int                     { return a.strength }
func (a *skillTestActor) GetDexterity() int                    { return a.dexterity }
func (a *skillTestActor) GetEnergy() int                       { return a.energy }
func (a *skillTestActor) GetVitality() int                     { return a.vitality }
func (a *skillTestActor) GetLeadership() int                   { return a.leadership }
func (*skillTestActor) GetChangeUp() int                       { return 0 }
func (*skillTestActor) CanUseItem(*item.Item) bool             { return true }
func (*skillTestActor) GetInventory() *item.Inventory          { return nil }
func (*skillTestActor) GetInventoryItem(int) *item.Item        { return nil }
func (*skillTestActor) GetWarehouse() *item.Warehouse          { return nil }
func (*skillTestActor) SetDelayRecoverHP(int, int)             {}
func (*skillTestActor) SetDelayRecoverSD(int, int)             {}
func (*skillTestActor) GetAttackRatePVP() int                  { return 1000 }
func (*skillTestActor) GetDefenseRatePVP() int                 { return 1 }
func (*skillTestActor) GetIgnoreDefenseRate() int              { return 0 }
func (*skillTestActor) GetCriticalAttackRate() int             { return 0 }
func (*skillTestActor) GetCriticalAttackDamage() int           { return 0 }
func (*skillTestActor) GetExcellentAttackRate() int            { return 0 }
func (*skillTestActor) GetExcellentAttackDamage() int          { return 0 }
func (*skillTestActor) GetMonsterDieGetHP() float64            { return 0 }
func (*skillTestActor) GetMonsterDieGetMP() float64            { return 0 }
func (*skillTestActor) GetAddDamage() int                      { return 0 }
func (*skillTestActor) GetArmorReduceDamage() int              { return 0 }
func (*skillTestActor) GetWingIncreaseDamage() int             { return 0 }
func (*skillTestActor) GetWingReduceDamage() int               { return 0 }
func (*skillTestActor) GetHelperReduceDamage() int             { return 0 }
func (*skillTestActor) GetPetIncreaseDamage() int              { return 0 }
func (*skillTestActor) GetPetReduceDamage() int                { return 0 }
func (*skillTestActor) GetDoubleDamageRate() int               { return 0 }
func (*skillTestActor) GetMonsterDieGetMoney() float64         { return 0 }
func (a *skillTestActor) GetKnightGladiatorCalcSkillBonus() float64 {
	if a.knightRate != 0 {
		return a.knightRate
	}
	return 1
}
func (a *skillTestActor) GetImpaleSkillCalc() float64 {
	if a.impaleRate != 0 {
		return a.impaleRate
	}
	return 1
}
func (*skillTestActor) PickItem(*model.MsgPickItem)       {}
func (*skillTestActor) DropItem(*model.MsgDropItem)       {}
func (*skillTestActor) BuyItem(*model.MsgBuyItem)         {}
func (*skillTestActor) SellItem(*model.MsgSellItem)       {}
func (*skillTestActor) MoveItem(*model.MsgMoveItem)       {}
func (*skillTestActor) UseItem(*model.MsgUseItem)         {}
func (*skillTestActor) RepairItem(*model.MsgRepairItem)   {}
func (*skillTestActor) Move(*model.MsgMove)               {}
func (*skillTestActor) Teleport(*model.MsgTeleport)       {}
func (*skillTestActor) MapMove(*model.MsgMapMove)         {}
func (*skillTestActor) SetPosition(*model.MsgSetPosition) {}
func (*skillTestActor) Action(*model.MsgAction)           {}
func (*skillTestActor) UseSkill(*model.MsgUseSkill)       {}
func (*skillTestActor) UseSkillDuration(*model.MsgUseSkillDuration) {
}
func (*skillTestActor) UseSkillAttackMultiTarget(*model.MsgUseSkillAttackMultiTarget) {
}
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

func countMessages[T any](messages []any) int {
	count := 0
	for _, msg := range messages {
		if _, ok := msg.(*T); ok {
			count++
		}
	}
	return count
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

func TestUseSkillDeferredSummonerSkillsDoNotCostResources(t *testing.T) {
	for _, index := range []int{
		skill.SkillIndexDrainLife,
		skill.SkillIndexSummonerExplosion,
	} {
		t.Run(fmt.Sprintf("skill_%d", index), func(t *testing.T) {
			caster, actor := newSkillTestObject(1, ObjectTypePlayer)
			target, _ := newSkillTestObject(2, ObjectTypeMonster)
			withTestObjectManager(t, caster, target)
			learnSkillForTest(t, caster, index)
			mp, ag := caster.MP, caster.AG

			caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: index})

			assertResourceUnchanged(t, caster, mp, ag)
			if hasMessage[model.MsgMPReply](actor.messages) {
				t.Fatal("resource reply was sent for deferred summoner skill")
			}
		})
	}
}

func TestUseSkillHealRestoresHPAndCostsResources(t *testing.T) {
	caster, actor := newSkillTestObject(1, ObjectTypePlayer)
	target, targetActor := newSkillTestObject(2, ObjectTypePlayer)
	withTestObjectManager(t, caster, target)
	caster.Class = int(class.Elf)
	actor.energy = 100
	target.HP = 80
	target.MaxHP = 100
	s := learnSkillForTest(t, caster, skill.SkillIndexHeal)

	caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: skill.SkillIndexHeal})

	if target.HP != target.MaxHP {
		t.Fatalf("target HP = %d, want %d", target.HP, target.MaxHP)
	}
	if caster.MP != 100-s.ManaUsage || caster.AG != 100-s.BPUsage {
		t.Fatalf("resources = %d/%d, want %d/%d", caster.MP, caster.AG, 100-s.ManaUsage, 100-s.BPUsage)
	}
	if !hasMessage[model.MsgUseSkillReply](actor.messages) {
		t.Fatal("skill success reply was not sent")
	}
	if !hasMessage[model.MsgHPReply](targetActor.messages) {
		t.Fatal("target HP reply was not sent")
	}
}

func TestUseSkillSupportSkillFailureDoesNotCostResources(t *testing.T) {
	caster, actor := newSkillTestObject(1, ObjectTypePlayer)
	target, _ := newSkillTestObject(2, ObjectTypePlayer)
	withTestObjectManager(t, caster, target)
	caster.Class = int(class.Wizard)
	learnSkillForTest(t, caster, skill.SkillIndexHeal)
	mp, ag := caster.MP, caster.AG

	caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: skill.SkillIndexHeal})

	assertResourceUnchanged(t, caster, mp, ag)
	if hasMessage[model.MsgMPReply](actor.messages) {
		t.Fatal("resource reply was sent for rejected support skill")
	}
}

func TestUseSkillGreaterAttackExpires(t *testing.T) {
	caster, actor := newSkillTestObject(1, ObjectTypePlayer)
	target, _ := newSkillTestObject(2, ObjectTypePlayer)
	withTestObjectManager(t, caster, target)
	caster.Class = int(class.Elf)
	actor.energy = 100
	target.AttackMin = 10
	target.AttackMax = 10
	learnSkillForTest(t, caster, skill.SkillIndexGreaterAttack)

	caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: skill.SkillIndexGreaterAttack})

	if damage := target.getDamage(skill.Skill0, 0, nil); damage <= 10 {
		t.Fatalf("damage = %d, want above 10", damage)
	}
	target.skillEffects[skill.SkillIndexGreaterAttack].expire = time.Now().Add(-time.Second)
	target.processSkillEffect()
	if damage := target.getDamage(skill.Skill0, 0, nil); damage != 10 {
		t.Fatalf("expired damage = %d, want 10", damage)
	}
}

func TestUseSkillGreaterDefenseExpires(t *testing.T) {
	caster, actor := newSkillTestObject(1, ObjectTypePlayer)
	target, _ := newSkillTestObject(2, ObjectTypePlayer)
	withTestObjectManager(t, caster, target)
	caster.Class = int(class.Elf)
	actor.energy = 100
	target.Defense = 10
	learnSkillForTest(t, caster, skill.SkillIndexGreaterDefense)

	caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: skill.SkillIndexGreaterDefense})

	if target.Defense <= 10 {
		t.Fatalf("target defense = %d, want above 10", target.Defense)
	}
	target.skillEffects[skill.SkillIndexGreaterDefense].expire = time.Now().Add(-time.Second)
	target.processSkillEffect()
	if target.Defense != 10 {
		t.Fatalf("expired target defense = %d, want 10", target.Defense)
	}
}

func TestUseSkillSwellHPExpiresAndClampsHP(t *testing.T) {
	caster, actor := newSkillTestObject(1, ObjectTypePlayer)
	target, targetActor := newSkillTestObject(2, ObjectTypePlayer)
	withTestObjectManager(t, caster, target)
	caster.Class = int(class.Knight)
	actor.energy = 100
	actor.vitality = 200
	target.HP = 100
	target.MaxHP = 100
	learnSkillForTest(t, caster, skill.SkillIndexSwellHP)

	caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: skill.SkillIndexSwellHP})

	if target.MaxHP <= 100 {
		t.Fatalf("target MaxHP = %d, want above 100", target.MaxHP)
	}
	target.HP = target.MaxHP
	target.skillEffects[skill.SkillIndexSwellHP].expire = time.Now().Add(-time.Second)
	target.processSkillEffect()
	if target.MaxHP != 100 || target.HP != 100 {
		t.Fatalf("expired HP/MaxHP = %d/%d, want 100/100", target.HP, target.MaxHP)
	}
	if !hasMessage[model.MsgHPReply](targetActor.messages) {
		t.Fatal("target HP reply was not sent")
	}
}

func TestUseSkillExpandedAttackSkills(t *testing.T) {
	for _, index := range []int{
		skill.SkillIndexPowerWave,
		skill.SkillIndexDecay,
		skill.SkillIndexIceStorm,
		skill.SkillIndexIceArrow,
		skill.SkillIndexForce,
		skill.SkillIndexFireBurst,
		skill.SkillIndexElectricSpike,
		skill.SkillIndexForceWave,
		skill.SkillIndexChainLightning,
		skill.SkillIndexLightningShock,
		skill.SkillIndexMultiShot,
		skill.SkillIndexKillingBlow,
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
		setup      func(*Object, *skillTestActor)
		targetType ObjectType
		wantDamage int
	}{
		{
			name:  "magic skill adds magic attack",
			index: skill.SkillIndexEvilSpirit,
			setup: func(caster *Object, actor *skillTestActor) {
				actor.magicAttackMin = 20
				actor.magicAttackMax = 20
			},
			wantDamage: 90,
		},
		{
			name:  "knight skill adds skill damage before bonus",
			index: skill.SkillIndexFallingSlash,
			setup: func(caster *Object, actor *skillTestActor) {
				actor.knightRate = 2
			},
			wantDamage: 200,
		},
		{
			name:       "power slash uses gladiator formula",
			index:      skill.SkillIndexPowerSlash,
			wantDamage: 200,
		},
		{
			name:       "elf skill uses elf formula",
			index:      skill.SkillIndexIceArrow,
			wantDamage: 200,
		},
		{
			name:       "triple shot uses elf formula",
			index:      skill.SkillIndexTripleShot,
			wantDamage: 200,
		},
		{
			name:       "multi shot uses elf formula",
			index:      skill.SkillIndexMultiShot,
			wantDamage: 200,
		},
		{
			name:  "dark lord skill uses lord formula",
			index: skill.SkillIndexFireBurst,
			setup: func(caster *Object, actor *skillTestActor) {
				actor.energy = 100
			},
			wantDamage: 205,
		},
		{
			name:  "force wave uses lord formula",
			index: skill.SkillIndexForceWave,
			setup: func(caster *Object, actor *skillTestActor) {
				actor.energy = 100
			},
			wantDamage: 205,
		},
		{
			name:  "summoner magic skill adds magic attack",
			index: skill.SkillIndexLightningShock,
			setup: func(caster *Object, actor *skillTestActor) {
				actor.magicAttackMin = 20
				actor.magicAttackMax = 20
			},
			wantDamage: 90,
		},
		{
			name:  "strike of destruction uses strike formula",
			index: skill.SkillIndexStrikeDestruction,
			setup: func(caster *Object, actor *skillTestActor) {
				actor.energy = 100
			},
			wantDamage: 210,
		},
		{
			name:       "flame strike uses flame formula",
			index:      skill.SkillIndexFlameStrike,
			wantDamage: 200,
		},
		{
			name:  "chain lightning uses single target formula",
			index: skill.SkillIndexChainLightning,
			setup: func(caster *Object, actor *skillTestActor) {
				actor.magicAttackMin = 20
				actor.magicAttackMax = 20
			},
			wantDamage: 90,
		},
		{
			name:  "rage fighter vitality skill uses vitality formula",
			index: skill.SkillIndexKillingBlow,
			setup: func(caster *Object, actor *skillTestActor) {
				actor.vitality = 100
			},
			wantDamage: 60,
		},
		{
			name:  "dark side uses dexterity and energy formula",
			index: skill.SkillIndexDarkSide,
			setup: func(caster *Object, actor *skillTestActor) {
				actor.dexterity = 80
				actor.energy = 100
			},
			wantDamage: 144,
		},
		{
			name:       "dragon slasher uses player target formula",
			index:      skill.SkillIndexDragonSlasher,
			setup:      func(caster *Object, actor *skillTestActor) { actor.energy = 100 },
			targetType: ObjectTypePlayer,
			wantDamage: 60,
		},
		{
			name:       "dragon slasher uses monster target formula",
			index:      skill.SkillIndexDragonSlasher,
			setup:      func(caster *Object, actor *skillTestActor) { actor.energy = 100 },
			targetType: ObjectTypeMonster,
			wantDamage: 480,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			caster, actor := newSkillTestObject(1, ObjectTypePlayer)
			targetType := tt.targetType
			if targetType == ObjectTypeEmpty {
				targetType = ObjectTypeMonster
			}
			target, _ := newSkillTestObject(2, targetType)
			caster.AttackMin = 30
			caster.AttackMax = 30
			s := learnSkillForTest(t, caster, tt.index)
			s.DamageMin = 70
			s.DamageMax = 70
			if tt.setup != nil {
				tt.setup(caster, actor)
			}

			if damage := caster.getDamage(s, 0, target); damage != tt.wantDamage {
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

func TestUseSkillAttackAreaTargetsNearbyMonsters(t *testing.T) {
	caster, actor := newSkillTestObject(1, ObjectTypePlayer)
	target, _ := newSkillTestObject(2, ObjectTypeMonster)
	nearby, _ := newSkillTestObject(3, ObjectTypeMonster)
	far, _ := newSkillTestObject(4, ObjectTypeMonster)
	player, _ := newSkillTestObject(5, ObjectTypePlayer)
	npc, _ := newSkillTestObject(6, ObjectTypeNPC)
	withTestObjectManager(t, caster, target, nearby, far, player, npc)

	caster.MP, caster.MaxMP = 1000, 1000
	caster.AG, caster.MaxAG = 1000, 1000
	caster.AttackRate = 1_000_000
	target.X, target.Y = 11, 10
	nearby.X, nearby.Y = 12, 10
	far.X, far.Y = 16, 10
	player.X, player.Y = 12, 11
	npc.X, npc.Y = 12, 9
	for _, obj := range []*Object{target, nearby, far, player, npc} {
		obj.HP = 1000
		obj.MaxHP = 1000
		if !caster.addViewportObject(obj) {
			t.Fatalf("addViewportObject(%d) = false", obj.Index)
		}
	}
	s := learnSkillForTest(t, caster, skill.SkillIndexDecay)

	caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: s.Index})

	if target.HP >= target.MaxHP || nearby.HP >= nearby.MaxHP {
		t.Fatalf("area targets HP = %d/%d, want both damaged", target.HP, nearby.HP)
	}
	if far.HP != far.MaxHP || player.HP != player.MaxHP || npc.HP != npc.MaxHP {
		t.Fatalf("excluded targets HP = %d/%d/%d, want unchanged", far.HP, player.HP, npc.HP)
	}
	if caster.MP != 1000-s.ManaUsage || caster.AG != 1000-s.BPUsage {
		t.Fatalf("resources = %d/%d, want %d/%d", caster.MP, caster.AG, 1000-s.ManaUsage, 1000-s.BPUsage)
	}
	if got := countMessages[model.MsgUseSkillReply](actor.messages); got != 1 {
		t.Fatalf("skill success replies = %d, want 1", got)
	}
	if got := countMessages[model.MsgMPReply](actor.messages); got != 1 {
		t.Fatalf("resource replies = %d, want 1", got)
	}
}

func TestUseSkillAttackFrustumTargetsFrontMonsters(t *testing.T) {
	caster, _ := newSkillTestObject(1, ObjectTypePlayer)
	target, _ := newSkillTestObject(2, ObjectTypeMonster)
	front, _ := newSkillTestObject(3, ObjectTypeMonster)
	behind, _ := newSkillTestObject(4, ObjectTypeMonster)
	withTestObjectManager(t, caster, target, front, behind)

	caster.AttackRate = 1_000_000
	target.X, target.Y = 10, 13
	front.X, front.Y = 11, 12
	behind.X, behind.Y = 10, 8
	for _, obj := range []*Object{target, front, behind} {
		obj.HP = 1000
		obj.MaxHP = 1000
		if !caster.addViewportObject(obj) {
			t.Fatalf("addViewportObject(%d) = false", obj.Index)
		}
	}
	s := learnSkillForTest(t, caster, skill.SkillIndexPowerSlash)

	caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: s.Index})

	if target.HP >= target.MaxHP || front.HP >= front.MaxHP {
		t.Fatalf("front targets HP = %d/%d, want both damaged", target.HP, front.HP)
	}
	if behind.HP != behind.MaxHP {
		t.Fatalf("behind target HP = %d, want %d", behind.HP, behind.MaxHP)
	}
}

func TestUseSkillAttackAreaSelfUsesCasterPosition(t *testing.T) {
	caster, _ := newSkillTestObject(1, ObjectTypePlayer)
	target, _ := newSkillTestObject(2, ObjectTypeMonster)
	nearCaster, _ := newSkillTestObject(3, ObjectTypeMonster)
	nearTarget, _ := newSkillTestObject(4, ObjectTypeMonster)
	withTestObjectManager(t, caster, target, nearCaster, nearTarget)

	caster.AttackRate = 1_000_000
	target.X, target.Y = 12, 10
	nearCaster.X, nearCaster.Y = 10, 12
	nearTarget.X, nearTarget.Y = 14, 10
	for _, obj := range []*Object{target, nearCaster, nearTarget} {
		obj.HP, obj.MaxHP = 1000, 1000
		if !caster.addViewportObject(obj) {
			t.Fatalf("addViewportObject(%d) = false", obj.Index)
		}
	}
	s := learnSkillForTest(t, caster, skill.SkillIndexTwistingSlash)

	caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: s.Index})

	if target.HP >= target.MaxHP || nearCaster.HP >= nearCaster.MaxHP {
		t.Fatalf("caster area targets HP = %d/%d, want both damaged", target.HP, nearCaster.HP)
	}
	if nearTarget.HP != nearTarget.MaxHP {
		t.Fatalf("target-centered monster HP = %d, want %d", nearTarget.HP, nearTarget.MaxHP)
	}
}

func TestUseSkillAttackHitBoxes(t *testing.T) {
	for _, tt := range []struct {
		name  string
		index int
	}{
		{name: "spear", index: skill.SkillIndexForce},
		{name: "electric", index: skill.SkillIndexElectricSpike},
	} {
		t.Run(tt.name, func(t *testing.T) {
			caster, _ := newSkillTestObject(1, ObjectTypePlayer)
			target, _ := newSkillTestObject(2, ObjectTypeMonster)
			front, _ := newSkillTestObject(3, ObjectTypeMonster)
			behind, _ := newSkillTestObject(4, ObjectTypeMonster)
			withTestObjectManager(t, caster, target, front, behind)

			caster.AttackRate = 1_000_000
			target.X, target.Y = 10, 12
			front.X, front.Y = 10, 13
			behind.X, behind.Y = 10, 7
			for _, obj := range []*Object{target, front, behind} {
				obj.HP, obj.MaxHP = 1000, 1000
				if !caster.addViewportObject(obj) {
					t.Fatalf("addViewportObject(%d) = false", obj.Index)
				}
			}
			s := learnSkillForTest(t, caster, tt.index)

			caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: s.Index})

			if target.HP >= target.MaxHP || front.HP >= front.MaxHP {
				t.Fatalf("hitbox targets HP = %d/%d, want both damaged", target.HP, front.HP)
			}
			if behind.HP != behind.MaxHP {
				t.Fatalf("behind target HP = %d, want %d", behind.HP, behind.MaxHP)
			}
		})
	}
}

func TestUseSkillChainLightningSelectsThreeTargets(t *testing.T) {
	caster, _ := newSkillTestObject(1, ObjectTypePlayer)
	target, _ := newSkillTestObject(2, ObjectTypeMonster)
	near1, _ := newSkillTestObject(3, ObjectTypeMonster)
	near2, _ := newSkillTestObject(4, ObjectTypeMonster)
	extra, _ := newSkillTestObject(5, ObjectTypeMonster)
	withTestObjectManager(t, caster, target, near1, near2, extra)

	caster.AttackRate = 1_000_000
	target.X, target.Y = 11, 10
	near1.X, near1.Y = 12, 10
	near2.X, near2.Y = 13, 10
	extra.X, extra.Y = 14, 10
	for _, obj := range []*Object{target, near1, near2, extra} {
		obj.HP, obj.MaxHP = 1000, 1000
		if !caster.addViewportObject(obj) {
			t.Fatalf("addViewportObject(%d) = false", obj.Index)
		}
	}
	s := learnSkillForTest(t, caster, skill.SkillIndexChainLightning)

	caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: s.Index})

	if target.HP >= target.MaxHP || near1.HP >= near1.MaxHP || near2.HP >= near2.MaxHP {
		t.Fatalf("chain targets HP = %d/%d/%d, want all damaged", target.HP, near1.HP, near2.HP)
	}
	if extra.HP != extra.MaxHP {
		t.Fatalf("fourth target HP = %d, want %d", extra.HP, extra.MaxHP)
	}
}

func TestUseSkillDarkSideSelectsFiveTargets(t *testing.T) {
	caster, actor := newSkillTestObject(1, ObjectTypePlayer)
	target, _ := newSkillTestObject(2, ObjectTypeMonster)
	withTestObjectManager(t, caster, target)

	caster.AttackRate = 1_000_000
	actor.dexterity = 80
	actor.energy = 100
	target.X, target.Y = 11, 10
	target.HP, target.MaxHP = 1000, 1000
	targets := []*Object{target}
	for i, position := range [][2]int{{10, 11}, {11, 11}, {9, 11}, {10, 12}, {10, 13}} {
		obj, _ := newSkillTestObject(i+3, ObjectTypeMonster)
		obj.X, obj.Y = position[0], position[1]
		obj.HP, obj.MaxHP = 1000, 1000
		ObjectManager.objects[obj.Index] = obj
		targets = append(targets, obj)
	}
	for _, obj := range targets {
		if !caster.addViewportObject(obj) {
			t.Fatalf("addViewportObject(%d) = false", obj.Index)
		}
	}
	s := learnSkillForTest(t, caster, skill.SkillIndexDarkSide)

	caster.UseSkill(&model.MsgUseSkill{Target: target.Index, Skill: s.Index})

	for _, obj := range targets[:5] {
		if obj.HP >= obj.MaxHP {
			t.Fatalf("target %d HP = %d, want damaged", obj.Index, obj.HP)
		}
	}
	if extra := targets[5]; extra.HP != extra.MaxHP {
		t.Fatalf("sixth target HP = %d, want %d", extra.HP, extra.MaxHP)
	}
}

func TestUseSkillDurationMultiTargetCostsOnceAndDamagesTargets(t *testing.T) {
	caster, actor := newSkillTestObject(1, ObjectTypePlayer)
	target1, _ := newSkillTestObject(2, ObjectTypeMonster)
	target2, _ := newSkillTestObject(3, ObjectTypeMonster)
	withTestObjectManager(t, caster, target1, target2)

	caster.MP, caster.MaxMP = 1000, 1000
	caster.AG, caster.MaxAG = 1000, 1000
	caster.AttackRate = 1_000_000
	target1.X, target1.Y = 11, 10
	target2.X, target2.Y = 12, 10
	for _, target := range []*Object{target1, target2} {
		target.HP, target.MaxHP = 1000, 1000
		if !caster.addViewportObject(target) {
			t.Fatalf("addViewportObject(%d) = false", target.Index)
		}
	}
	s := learnSkillForTest(t, caster, skill.SkillIndexTwister)

	caster.UseSkill(&model.MsgUseSkill{Target: target1.Index, Skill: s.Index})
	assertResourceUnchanged(t, caster, 1000, 1000)
	if target1.HP != target1.MaxHP {
		t.Fatalf("legacy use target HP = %d, want %d", target1.HP, target1.MaxHP)
	}

	caster.UseSkillDuration(&model.MsgUseSkillDuration{
		Target: target1.Index,
		Skill:  s.Index,
		Dir:    3,
	})
	wantMP, wantAG := 1000-s.ManaUsage, 1000-s.BPUsage
	if caster.MP != wantMP || caster.AG != wantAG {
		t.Fatalf("resources = %d/%d, want %d/%d", caster.MP, caster.AG, wantMP, wantAG)
	}
	if target1.HP != target1.MaxHP || target2.HP != target2.MaxHP {
		t.Fatalf("duration cast damaged targets: %d/%d", target1.HP, target2.HP)
	}
	if got := countMessages[model.MsgUseSkillDurationReply](actor.messages); got != 1 {
		t.Fatalf("duration replies = %d, want 1", got)
	}

	caster.UseSkillAttackMultiTarget(&model.MsgUseSkillAttackMultiTarget{
		Skill: s.Index,
		Targets: []model.MultiTarget{
			{Target: target1.Index},
			{Target: target2.Index},
		},
	})

	if target1.HP >= target1.MaxHP || target2.HP >= target2.MaxHP {
		t.Fatalf("multi targets HP = %d/%d, want both damaged", target1.HP, target2.HP)
	}
	if caster.MP != wantMP || caster.AG != wantAG {
		t.Fatalf("resources after hits = %d/%d, want %d/%d", caster.MP, caster.AG, wantMP, wantAG)
	}
	if got := countMessages[model.MsgMPReply](actor.messages); got != 1 {
		t.Fatalf("resource replies = %d, want 1", got)
	}
	if got := countMessages[model.MsgUseSkillReply](actor.messages); got != 0 {
		t.Fatalf("regular skill replies = %d, want 0", got)
	}
}

func TestUseSkillAttackMultiTargetRejectsInvalidState(t *testing.T) {
	for _, tt := range []struct {
		name       string
		skillIndex int
		setup      func(caster, target *Object)
		msg        model.MsgUseSkillAttackMultiTarget
	}{
		{
			name:       "not cast",
			skillIndex: skill.SkillIndexTwister,
			msg: model.MsgUseSkillAttackMultiTarget{
				Skill: skill.SkillIndexTwister,
				Targets: []model.MultiTarget{{
					Target: 2,
				}},
			},
		},
		{
			name:       "skill mismatch",
			skillIndex: skill.SkillIndexTwister,
			setup: func(caster, target *Object) {
				caster.UseSkillDuration(&model.MsgUseSkillDuration{
					Target: target.Index,
					Skill:  skill.SkillIndexTwister,
				})
			},
			msg: model.MsgUseSkillAttackMultiTarget{
				Skill: skill.SkillIndexFlame,
				Targets: []model.MultiTarget{{
					Target: 2,
				}},
			},
		},
		{
			name:       "expired cast",
			skillIndex: skill.SkillIndexTwister,
			setup: func(caster, target *Object) {
				caster.UseSkillDuration(&model.MsgUseSkillDuration{
					Target: target.Index,
					Skill:  skill.SkillIndexTwister,
				})
				caster.durationSkill.startedAt = time.Now().Add(-9 * time.Second)
			},
			msg: model.MsgUseSkillAttackMultiTarget{
				Skill: skill.SkillIndexTwister,
				Targets: []model.MultiTarget{{
					Target: 2,
				}},
			},
		},
		{
			name:       "invalid evil spirit key",
			skillIndex: skill.SkillIndexEvilSpirit,
			setup: func(caster, target *Object) {
				caster.UseSkillDuration(&model.MsgUseSkillDuration{
					Target:   target.Index,
					Skill:    skill.SkillIndexEvilSpirit,
					MagicKey: 1,
				})
			},
			msg: model.MsgUseSkillAttackMultiTarget{
				Skill: skill.SkillIndexEvilSpirit,
				Targets: []model.MultiTarget{{
					Target:   2,
					MagicKey: 2,
				}},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			caster, _ := newSkillTestObject(1, ObjectTypePlayer)
			target, _ := newSkillTestObject(2, ObjectTypeMonster)
			withTestObjectManager(t, caster, target)
			caster.MP, caster.MaxMP = 1000, 1000
			caster.AG, caster.MaxAG = 1000, 1000
			caster.AttackRate = 1_000_000
			target.X = caster.X + 1
			target.HP, target.MaxHP = 1000, 1000
			if !caster.addViewportObject(target) {
				t.Fatal("addViewportObject() = false")
			}
			learnSkillForTest(t, caster, tt.skillIndex)
			if tt.setup != nil {
				tt.setup(caster, target)
			}
			mp, ag := caster.MP, caster.AG

			caster.UseSkillAttackMultiTarget(&tt.msg)

			if target.HP != target.MaxHP {
				t.Fatalf("target HP = %d, want %d", target.HP, target.MaxHP)
			}
			assertResourceUnchanged(t, caster, mp, ag)
		})
	}
}
