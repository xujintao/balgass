package object

import (
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/effect"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

type basicAttackTestActor struct {
	*skillTestActor
	inventory                  item.Inventory
	attackSpeed                int
	hitChecks                  int
	observeAmmoPosition        int
	observedAmmoDurability     int
	ammoConsumedBeforeHitCheck bool
}

func newBasicAttackTestObject(index int, typ ObjectType) (*Object, *basicAttackTestActor) {
	obj, base := newSkillTestObject(index, typ)
	actor := &basicAttackTestActor{
		skillTestActor:      base,
		observeAmmoPosition: -1,
	}
	actor.inventory.Items = make([]*item.Item, item.INVENTORY_WEAR_SIZE)
	actor.inventory.Flags = make([]bool, item.INVENTORY_WEAR_SIZE)
	obj.Objecter = actor
	return obj, actor
}

func (a *basicAttackTestActor) GetInventory() *item.Inventory {
	return &a.inventory
}

func (a *basicAttackTestActor) GetInventoryItem(position int) *item.Item {
	return a.inventory.Items[position]
}

func (a *basicAttackTestActor) GetAttackSpeedForDelay() int {
	return a.attackSpeed
}

func (a *basicAttackTestActor) GetAttackRatePVP() int {
	a.hitChecks++
	if a.observeAmmoPosition >= 0 {
		ammo := a.inventory.Items[a.observeAmmoPosition]
		a.ammoConsumedBeforeHitCheck = ammo != nil &&
			ammo.Durability == a.observedAmmoDurability
	}
	return 1000
}

func countBasicAttackMessages[T any](messages []any) int {
	count := 0
	for _, message := range messages {
		if _, ok := message.(T); ok {
			count++
		}
	}
	return count
}

func prepareBasicAttackTest(t *testing.T) (*Object, *basicAttackTestActor, *Object, *basicAttackTestActor) {
	t.Helper()
	attacker, attackerActor := newBasicAttackTestObject(1, ObjectTypePlayer)
	target, targetActor := newBasicAttackTestObject(2, ObjectTypePlayer)
	withTestObjectManager(t, attacker, target)
	if !attacker.addViewportObject(target) {
		t.Fatal("failed to add target to attacker viewport")
	}
	if !target.addViewportObject(attacker) {
		t.Fatal("failed to add attacker to target viewport")
	}
	return attacker, attackerActor, target, targetActor
}

func findBasicAttackMapPosition(t *testing.T, safe bool) (int, int) {
	t.Helper()
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			if maps.MapManager.GetMapAttr(0, x, y)&1 != 0 == safe {
				return x, y
			}
		}
	}
	t.Fatalf("map 0 has no position with safe=%t", safe)
	return 0, 0
}

func TestBasicAttackRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Object, *Object, **model.MsgAttack)
	}{
		{name: "nil message", mutate: func(_ *Object, _ *Object, msg **model.MsgAttack) {
			*msg = nil
		}},
		{name: "attacker not playing", mutate: func(attacker, _ *Object, _ **model.MsgAttack) {
			attacker.ConnectState = ConnectStateLogged
		}},
		{name: "attacker dead", mutate: func(attacker, _ *Object, _ **model.MsgAttack) {
			attacker.Live = false
		}},
		{name: "attacker death state", mutate: func(attacker, _ *Object, _ **model.MsgAttack) {
			attacker.State = 4
		}},
		{name: "attacker cleanup state", mutate: func(attacker, _ *Object, _ **model.MsgAttack) {
			attacker.State = 8
		}},
		{name: "attacker cannot act", mutate: func(attacker, _ *Object, _ **model.MsgAttack) {
			attacker.effects[effect.BuffSleep] = &effect.Effect{Sleep: true}
		}},
		{name: "negative index", mutate: func(_ *Object, _ *Object, msg **model.MsgAttack) {
			(*msg).Target = -1
		}},
		{name: "index past object table", mutate: func(_ *Object, _ *Object, msg **model.MsgAttack) {
			(*msg).Target = len(ObjectManager.objects)
		}},
		{name: "empty object slot", mutate: func(_ *Object, _ *Object, msg **model.MsgAttack) {
			(*msg).Target = 3
		}},
		{name: "self target", mutate: func(attacker, _ *Object, msg **model.MsgAttack) {
			(*msg).Target = attacker.Index
		}},
		{name: "target dead", mutate: func(_ *Object, target *Object, _ **model.MsgAttack) {
			target.Live = false
		}},
		{name: "target death state", mutate: func(_ *Object, target *Object, _ **model.MsgAttack) {
			target.State = 4
		}},
		{name: "target cleanup state", mutate: func(_ *Object, target *Object, _ **model.MsgAttack) {
			target.State = 8
		}},
		{name: "different map", mutate: func(_ *Object, target *Object, _ **model.MsgAttack) {
			target.MapNumber = 1
		}},
		{name: "outside active viewport", mutate: func(attacker, target *Object, _ **model.MsgAttack) {
			for _, vp := range attacker.Viewports {
				if vp.Number == target.Index {
					vp.reset()
					break
				}
			}
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attacker, attackerActor, target, targetActor := prepareBasicAttackTest(t)
			msg := &model.MsgAttack{Target: target.Index}
			tt.mutate(attacker, target, &msg)

			attacker.Attack(msg)

			if attackerActor.hitChecks != 0 {
				t.Fatalf("hit checks = %d, want 0", attackerActor.hitChecks)
			}
			if got := countBasicAttackMessages[*model.MsgActionReply](targetActor.messages); got != 0 {
				t.Fatalf("action messages = %d, want 0", got)
			}
			if !attacker.lastBasicAttackTime.IsZero() {
				t.Fatalf("last attack time = %v, want zero", attacker.lastBasicAttackTime)
			}
			if target.HP != target.MaxHP {
				t.Fatalf("target HP = %d, want %d", target.HP, target.MaxHP)
			}
		})
	}
}

func TestBasicAttackRejectsSafeZones(t *testing.T) {
	old := conf.Common.AntiHack.EnableBlockAttackInSafeZone
	conf.Common.AntiHack.EnableBlockAttackInSafeZone = true
	t.Cleanup(func() {
		conf.Common.AntiHack.EnableBlockAttackInSafeZone = old
	})
	safeX, safeY := findBasicAttackMapPosition(t, true)
	unsafeX, unsafeY := findBasicAttackMapPosition(t, false)

	for _, tt := range []struct {
		name         string
		attackerSafe bool
		targetSafe   bool
	}{
		{name: "attacker in safe zone", attackerSafe: true},
		{name: "target in safe zone", targetSafe: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			attacker, attackerActor, target, targetActor := prepareBasicAttackTest(t)
			attacker.X, attacker.Y = unsafeX, unsafeY
			target.X, target.Y = unsafeX, unsafeY
			if tt.attackerSafe {
				attacker.X, attacker.Y = safeX, safeY
			}
			if tt.targetSafe {
				target.X, target.Y = safeX, safeY
			}

			attacker.Attack(&model.MsgAttack{Target: target.Index})

			if attackerActor.hitChecks != 0 {
				t.Fatalf("hit checks = %d, want 0", attackerActor.hitChecks)
			}
			if got := countBasicAttackMessages[*model.MsgActionReply](targetActor.messages); got != 0 {
				t.Fatalf("action messages = %d, want 0", got)
			}
			if !attacker.lastBasicAttackTime.IsZero() {
				t.Fatalf("last attack time = %v, want zero", attacker.lastBasicAttackTime)
			}
		})
	}
}

func TestBasicAttackCooldown(t *testing.T) {
	attacker, attackerActor, target, targetActor := prepareBasicAttackTest(t)
	msg := &model.MsgAttack{Target: target.Index}

	attacker.Attack(msg)
	firstAttackTime := attacker.lastBasicAttackTime
	attacker.Attack(msg)

	if attackerActor.hitChecks != 1 {
		t.Fatalf("hit checks during cooldown = %d, want 1", attackerActor.hitChecks)
	}
	if got := countBasicAttackMessages[*model.MsgActionReply](targetActor.messages); got != 1 {
		t.Fatalf("action messages during cooldown = %d, want 1", got)
	}
	if attacker.lastBasicAttackTime != firstAttackTime {
		t.Fatalf("cooldown rejection changed last attack time")
	}

	attacker.lastBasicAttackTime = time.Now().Add(-basicAttackDelay(attackerActor.attackSpeed))
	attacker.Attack(msg)
	if attackerActor.hitChecks != 2 {
		t.Fatalf("hit checks at cooldown boundary = %d, want 2", attackerActor.hitChecks)
	}
}

func TestBasicAttackWithoutAmmoUsesCooldown(t *testing.T) {
	attacker, attackerActor, target, targetActor := prepareBasicAttackTest(t)
	attackerActor.inventory.Items[1] = item.NewItem(4, 8)
	msg := &model.MsgAttack{Target: target.Index}

	attacker.Attack(msg)
	if attacker.lastBasicAttackTime.IsZero() {
		t.Fatal("attack without ammo did not record cooldown")
	}
	if attackerActor.hitChecks != 0 {
		t.Fatalf("hit checks without ammo = %d, want 0", attackerActor.hitChecks)
	}
	if got := countBasicAttackMessages[*model.MsgActionReply](targetActor.messages); got != 1 {
		t.Fatalf("action messages without ammo = %d, want 1", got)
	}

	bolt := item.NewItem(4, 7)
	bolt.Durability = 2
	attackerActor.inventory.Items[0] = bolt
	attacker.Attack(msg)
	if bolt.Durability != 2 {
		t.Fatalf("ammo durability during cooldown = %d, want 2", bolt.Durability)
	}
	if attackerActor.hitChecks != 0 {
		t.Fatalf("hit checks during cooldown = %d, want 0", attackerActor.hitChecks)
	}
}

func TestBasicAttackAmmo(t *testing.T) {
	tests := []struct {
		name              string
		setup             func(*testing.T, *Object, *basicAttackTestActor)
		ammoPosition      int
		wantDurability    int
		wantHitCheck      bool
		wantDurabilityMsg bool
	}{
		{
			name: "bow consumes right hand arrow",
			setup: func(_ *testing.T, _ *Object, actor *basicAttackTestActor) {
				actor.inventory.Items[0] = item.NewItem(4, 0)
				actor.inventory.Items[1] = item.NewItem(4, 15)
				actor.inventory.Items[1].Durability = 2
			},
			ammoPosition:      1,
			wantDurability:    1,
			wantHitCheck:      true,
			wantDurabilityMsg: true,
		},
		{
			name: "crossbow consumes left hand bolt",
			setup: func(_ *testing.T, _ *Object, actor *basicAttackTestActor) {
				actor.inventory.Items[1] = item.NewItem(4, 8)
				actor.inventory.Items[0] = item.NewItem(4, 7)
				actor.inventory.Items[0].Durability = 2
			},
			ammoPosition:      0,
			wantDurability:    1,
			wantHitCheck:      true,
			wantDurabilityMsg: true,
		},
		{
			name: "wrong ammo",
			setup: func(_ *testing.T, _ *Object, actor *basicAttackTestActor) {
				actor.inventory.Items[0] = item.NewItem(4, 0)
				actor.inventory.Items[1] = item.NewItem(4, 7)
				actor.inventory.Items[1].Durability = 2
			},
			ammoPosition:   1,
			wantDurability: 2,
		},
		{
			name: "zero durability",
			setup: func(_ *testing.T, _ *Object, actor *basicAttackTestActor) {
				actor.inventory.Items[0] = item.NewItem(4, 0)
				actor.inventory.Items[1] = item.NewItem(4, 15)
			},
			ammoPosition:   1,
			wantDurability: 0,
		},
		{
			name: "infinity arrow bypasses ammo",
			setup: func(t *testing.T, attacker *Object, actor *basicAttackTestActor) {
				actor.inventory.Items[0] = item.NewItem(4, 0)
				if !attacker.addEffect(&effect.Effect{
					BuffIndex: effect.BuffInfinityArrow,
					Expire:    time.Now().Add(time.Minute),
				}) {
					t.Fatal("failed to add infinity arrow")
				}
			},
			ammoPosition: -1,
			wantHitCheck: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attacker, attackerActor, target, targetActor := prepareBasicAttackTest(t)
			tt.setup(t, attacker, attackerActor)
			attackerActor.observeAmmoPosition = tt.ammoPosition
			attackerActor.observedAmmoDurability = tt.wantDurability

			attacker.Attack(&model.MsgAttack{Target: target.Index})

			gotHitCheck := attackerActor.hitChecks == 1
			if gotHitCheck != tt.wantHitCheck {
				t.Fatalf("hit check = %t, want %t", gotHitCheck, tt.wantHitCheck)
			}
			if tt.ammoPosition >= 0 {
				if got := attackerActor.inventory.Items[tt.ammoPosition].Durability; got != tt.wantDurability {
					t.Fatalf("ammo durability = %d, want %d", got, tt.wantDurability)
				}
			}
			if tt.wantHitCheck && tt.ammoPosition >= 0 && !attackerActor.ammoConsumedBeforeHitCheck {
				t.Fatal("ammo was not consumed before hit check")
			}
			gotDurabilityMsg := countBasicAttackMessages[*model.MsgItemDurabilityReply](attackerActor.messages) == 1
			if gotDurabilityMsg != tt.wantDurabilityMsg {
				t.Fatalf("durability message = %t, want %t", gotDurabilityMsg, tt.wantDurabilityMsg)
			}
			if got := countBasicAttackMessages[*model.MsgActionReply](targetActor.messages); got != 1 {
				t.Fatalf("action messages = %d, want 1", got)
			}
			if attacker.lastBasicAttackTime.IsZero() {
				t.Fatal("legal attack did not record cooldown")
			}
		})
	}
}

func TestBasicAttackDelay(t *testing.T) {
	info := &conf.CommonServer.GameServerInfo
	oldTimeLimit := info.AttackSpeedTimeLimit
	oldDecrement := info.DecTimePerAttackSpeed
	oldMinimum := info.MinimumAttackSpeedTime
	t.Cleanup(func() {
		info.AttackSpeedTimeLimit = oldTimeLimit
		info.DecTimePerAttackSpeed = oldDecrement
		info.MinimumAttackSpeedTime = oldMinimum
	})

	for _, tt := range []struct {
		name      string
		timeLimit int
		decrement float64
		minimum   int
		speed     int
		want      time.Duration
	}{
		{name: "defaults", want: 800 * time.Millisecond},
		{name: "configured", timeLimit: 1000, decrement: 10, minimum: 250, speed: 40, want: 600 * time.Millisecond},
		{name: "minimum", timeLimit: 1000, decrement: 10, minimum: 250, speed: 100, want: 250 * time.Millisecond},
		{name: "negative speed", timeLimit: 1000, decrement: 10, minimum: 250, speed: -1, want: time.Second},
	} {
		t.Run(tt.name, func(t *testing.T) {
			info.AttackSpeedTimeLimit = tt.timeLimit
			info.DecTimePerAttackSpeed = tt.decrement
			info.MinimumAttackSpeedTime = tt.minimum
			if got := basicAttackDelay(tt.speed); got != tt.want {
				t.Fatalf("delay = %v, want %v", got, tt.want)
			}
		})
	}
}

func withShieldSystemTest(t *testing.T, enabled bool, rate int) {
	t.Helper()
	info := &conf.CommonServer.GameServerInfo
	oldEnabled := info.ShieldSystemEnable
	oldRate := info.DamageDivideSDRate
	info.ShieldSystemEnable = enabled
	info.DamageDivideSDRate = rate
	t.Cleanup(func() {
		info.ShieldSystemEnable = oldEnabled
		info.DamageDivideSDRate = oldRate
	})
}

func findAttackDamageReply(messages []any) *model.MsgAttackDamageReply {
	for i := len(messages) - 1; i >= 0; i-- {
		if reply, ok := messages[i].(*model.MsgAttackDamageReply); ok {
			return reply
		}
	}
	return nil
}

type damageFormulaTestActor struct {
	*skillTestActor
	masterLevel       int
	addDamage         int
	armorReduceDamage int
}

type attackProcTestActor struct {
	*skillTestActor
	inventory      item.Inventory
	fullHPRecovery float64
	fullMPRecovery float64
	fullSDRecovery float64
	maceStun       float64
}

func newAttackProcTestObject(index int, typ ObjectType) (*Object, *attackProcTestActor) {
	obj, base := newSkillTestObject(index, typ)
	actor := &attackProcTestActor{skillTestActor: base}
	actor.inventory.Items = make([]*item.Item, item.INVENTORY_WEAR_SIZE)
	actor.inventory.Flags = make([]bool, item.INVENTORY_WEAR_SIZE)
	obj.Objecter = actor
	return obj, actor
}

func (a *attackProcTestActor) GetInventory() *item.Inventory {
	return &a.inventory
}

func (a *attackProcTestActor) GetInventoryItem(position int) *item.Item {
	return a.inventory.Items[position]
}

func (a *attackProcTestActor) GetFullHPRecoveryRate() float64 {
	return a.fullHPRecovery
}

func (a *attackProcTestActor) GetFullMPRecoveryRate() float64 {
	return a.fullMPRecovery
}

func (a *attackProcTestActor) GetFullSDRecoveryRate() float64 {
	return a.fullSDRecovery
}

func (a *attackProcTestActor) GetMaceStunRate() float64 {
	return a.maceStun
}

func newDamageFormulaTestObject(index int, typ ObjectType) (*Object, *damageFormulaTestActor) {
	obj, base := newSkillTestObject(index, typ)
	actor := &damageFormulaTestActor{skillTestActor: base}
	obj.Objecter = actor
	return obj, actor
}

func (a *damageFormulaTestActor) GetMasterLevel() int {
	return a.masterLevel
}

func (a *damageFormulaTestActor) GetAddDamage() int {
	return a.addDamage
}

func (a *damageFormulaTestActor) GetArmorReduceDamage() int {
	return a.armorReduceDamage
}

func withDamageRateTest(t *testing.T, pvpRate, pveRate int) {
	t.Helper()
	oldPVP := conf.CalcChar.DamageRate.PVP.DarkWizard.DarkKnight
	oldPVE := conf.CalcChar.DamageRate.PVE.DarkWizard
	conf.CalcChar.DamageRate.PVP.DarkWizard.DarkKnight = pvpRate
	conf.CalcChar.DamageRate.PVE.DarkWizard = pveRate
	t.Cleanup(func() {
		conf.CalcChar.DamageRate.PVP.DarkWizard.DarkKnight = oldPVP
		conf.CalcChar.DamageRate.PVE.DarkWizard = oldPVE
	})
}

func TestCalculatedDamageStageOrder(t *testing.T) {
	withShieldSystemTest(t, false, 0)
	withDamageRateTest(t, 100, 100)
	attacker, attackerActor := newDamageFormulaTestObject(1, ObjectTypePlayer)
	target, targetActor := newDamageFormulaTestObject(2, ObjectTypePlayer)
	attacker.Class = 0
	target.Class = 1
	attacker.Level = 90
	attackerActor.masterLevel = 10
	attacker.AttackMin, attacker.AttackMax = 10, 10
	target.Level = 0
	target.Defense = 100
	targetActor.armorReduceDamage = 50
	attacker.effects[9001] = &effect.Effect{AttackReduction: 50}
	withTestObjectManager(t, attacker, target)

	damage := attacker.attack(target, attackRequest{mode: attackModeCalculated})

	if damage != 5 {
		t.Fatalf("attack damage = %d, want 5", damage)
	}
	if target.HP != 95 {
		t.Fatalf("target HP = %d, want 95", target.HP)
	}
}

func TestCalculatedAddDamageAfterSoulBarrier(t *testing.T) {
	withShieldSystemTest(t, false, 0)
	withDamageRateTest(t, 100, 100)
	attacker, attackerActor := newDamageFormulaTestObject(1, ObjectTypePlayer)
	target, _ := newDamageFormulaTestObject(2, ObjectTypePlayer)
	attacker.Class = 0
	target.Class = 1
	target.Level = 0
	attacker.AttackMin, attacker.AttackMax = 20, 20
	attackerActor.addDamage = 10
	target.effects[effect.BuffSoulBarrier] = &effect.Effect{
		DamageReduction: 50,
		ManaRate:        100,
	}
	withTestObjectManager(t, attacker, target)

	damage := attacker.attack(target, attackRequest{mode: attackModeCalculated})

	if damage != 20 {
		t.Fatalf("attack damage = %d, want 20", damage)
	}
	if target.HP != 80 || target.MP != 90 {
		t.Fatalf("target resources = HP:%d MP:%d, want HP:80 MP:90", target.HP, target.MP)
	}
}

func TestDamageRateSelection(t *testing.T) {
	withDamageRateTest(t, 50, 80)
	oldZeroRate := conf.CalcChar.DamageRate.PVP.FairyElf.DarkKnight
	conf.CalcChar.DamageRate.PVP.FairyElf.DarkKnight = 0
	t.Cleanup(func() {
		conf.CalcChar.DamageRate.PVP.FairyElf.DarkKnight = oldZeroRate
	})

	for _, tt := range []struct {
		name          string
		attackerType  ObjectType
		attackerClass int
		targetType    ObjectType
		targetClass   int
		want          int
	}{
		{name: "pvp matrix", attackerType: ObjectTypePlayer, attackerClass: 0, targetType: ObjectTypePlayer, targetClass: 1, want: 10},
		{name: "pve rate", attackerType: ObjectTypePlayer, attackerClass: 0, targetType: ObjectTypeMonster, want: 16},
		{name: "invalid attacker class", attackerType: ObjectTypePlayer, attackerClass: 99, targetType: ObjectTypePlayer, targetClass: 1, want: 20},
		{name: "invalid target class", attackerType: ObjectTypePlayer, attackerClass: 0, targetType: ObjectTypePlayer, targetClass: 99, want: 20},
		{name: "zero rate", attackerType: ObjectTypePlayer, attackerClass: 2, targetType: ObjectTypePlayer, targetClass: 1, want: 0},
		{name: "monster attacker", attackerType: ObjectTypeMonster, attackerClass: 0, targetType: ObjectTypePlayer, targetClass: 1, want: 20},
	} {
		t.Run(tt.name, func(t *testing.T) {
			attacker, _ := newSkillTestObject(1, tt.attackerType)
			target, _ := newSkillTestObject(2, tt.targetType)
			attacker.Class = tt.attackerClass
			target.Class = tt.targetClass

			if got := attacker.applyDamageRate(target, 20); got != tt.want {
				t.Fatalf("damage = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFixedAndDOTIgnoreCalculatedDamageStages(t *testing.T) {
	withShieldSystemTest(t, false, 0)
	withDamageRateTest(t, 0, 100)
	attacker, _ := newSkillTestObject(1, ObjectTypePlayer)
	target, _ := newSkillTestObject(2, ObjectTypePlayer)
	attacker.Class = 0
	target.Class = 1
	target.Level = 0
	withTestObjectManager(t, attacker, target)

	if damage := attacker.attack(target, attackRequest{mode: attackModeCalculated}); damage != 0 {
		t.Fatalf("calculated damage = %d, want 0", damage)
	}
	if damage := attacker.attack(target, attackRequest{mode: attackModeFixed, damage: 2}); damage != 2 {
		t.Fatalf("fixed damage = %d, want 2", damage)
	}
	if damage := attacker.attack(target, attackRequest{mode: attackModeDOT, damage: 2}); damage != 2 {
		t.Fatalf("DOT damage = %d, want 2", damage)
	}
	if target.HP != 96 {
		t.Fatalf("target HP = %d, want 96", target.HP)
	}
}

func TestFullRecoveryRequiresPositiveDamage(t *testing.T) {
	withShieldSystemTest(t, false, 0)
	for _, tt := range []struct {
		name         string
		damage       int
		wantRecovery bool
	}{
		{name: "positive damage", damage: 10, wantRecovery: true},
		{name: "zero damage"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			attacker, _ := newAttackProcTestObject(1, ObjectTypePlayer)
			target, targetActor := newAttackProcTestObject(2, ObjectTypePlayer)
			attacker.Class = 0
			target.Class = 1
			target.Level = 0
			attacker.SD, attacker.MaxSD = 10, 100
			target.HP, target.MP = 50, 40
			target.SD, target.MaxSD = 10, 100
			targetActor.fullHPRecovery = 100
			targetActor.fullMPRecovery = 100
			targetActor.fullSDRecovery = 100
			withTestObjectManager(t, attacker, target)

			attacker.attack(target, attackRequest{mode: attackModeFixed, damage: tt.damage})

			queued := 0
			for _, msg := range target.msgs {
				if msg.code >= 13 && msg.code <= 15 {
					queued++
					msg.time = time.Now().Add(-time.Millisecond)
				}
			}
			if tt.wantRecovery {
				if queued != 3 {
					t.Fatalf("recovery messages = %d, want 3", queued)
				}
				target.processDelayMsg()
				if target.HP != target.MaxHP || target.MP != target.MaxMP || target.SD != target.MaxSD {
					t.Fatalf("target resources after delay = HP:%d MP:%d SD:%d", target.HP, target.MP, target.SD)
				}
			} else if queued != 0 {
				t.Fatalf("recovery messages = %d, want 0", queued)
			}
			if attacker.SD != 10 {
				t.Fatalf("attacker SD = %d, want 10", attacker.SD)
			}
		})
	}
}

func TestMaceStunAppliesToFixedZeroDamage(t *testing.T) {
	withShieldSystemTest(t, false, 0)
	attacker, attackerActor := newAttackProcTestObject(1, ObjectTypePlayer)
	target, _ := newAttackProcTestObject(2, ObjectTypePlayer)
	attackerActor.inventory.Items[0] = item.NewItem(2, 0)
	attackerActor.inventory.Items[0].Durability = 1
	attackerActor.maceStun = 100
	target.PathMoving = true
	withTestObjectManager(t, attacker, target)

	attacker.attack(target, attackRequest{mode: attackModeFixed})

	if !target.HasBuff(effect.BuffStun) || !target.cannotAct() {
		t.Fatal("mace attack did not stun target")
	}
	if target.PathMoving {
		t.Fatal("stunned target continued moving")
	}
	x, y := target.X, target.Y
	target.Move(&model.MsgMove{X: x + 1, Y: y + 1})
	if target.X != x || target.Y != y {
		t.Fatal("stunned target accepted a new movement request")
	}
}

func TestReflectedAndReturnedDamageCanPassivelyStun(t *testing.T) {
	withShieldSystemTest(t, false, 0)
	for _, tt := range []struct {
		name    string
		msgCode int
		setup   func(*Object, *attackProcTestActor)
	}{
		{
			name:    "reflected",
			msgCode: 9,
			setup: func(defender *Object, _ *attackProcTestActor) {
				defender.effects[effect.BuffDamageReflection] = &effect.Effect{Reflect: 100}
			},
		},
		{
			name:    "returned",
			msgCode: 12,
			setup: func(_ *Object, defenderActor *attackProcTestActor) {
				defenderActor.returnDamage = 100
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			attacker, _ := newAttackProcTestObject(1, ObjectTypePlayer)
			defender, defenderActor := newAttackProcTestObject(2, ObjectTypePlayer)
			defenderActor.inventory.Items[0] = item.NewItem(2, 0)
			defenderActor.inventory.Items[0].Durability = 1
			defenderActor.maceStun = 100
			tt.setup(defender, defenderActor)
			withTestObjectManager(t, attacker, defender)

			attacker.attack(defender, attackRequest{mode: attackModeFixed, damage: 20})
			found := false
			for _, msg := range defender.msgs {
				if msg.code == tt.msgCode {
					found = true
					msg.time = time.Now().Add(-time.Millisecond)
				}
			}
			if !found {
				t.Fatalf("delay message %d was not queued", tt.msgCode)
			}
			defender.processDelayMsg()

			if !attacker.HasBuff(effect.BuffStun) {
				t.Fatal("passive damage did not stun original attacker")
			}
		})
	}
}

func TestSplitDamage(t *testing.T) {
	for _, tt := range []struct {
		name         string
		enabled      bool
		rate         int
		attackerType ObjectType
		targetType   ObjectType
		targetSD     int
		wantHP       int
		wantSD       int
	}{
		{name: "shield disabled", rate: 90, attackerType: ObjectTypePlayer, targetType: ObjectTypePlayer, targetSD: 200, wantHP: 100},
		{name: "non pvp", enabled: true, rate: 90, attackerType: ObjectTypePlayer, targetType: ObjectTypeMonster, targetSD: 200, wantHP: 100},
		{name: "normal split", enabled: true, rate: 90, attackerType: ObjectTypePlayer, targetType: ObjectTypePlayer, targetSD: 200, wantHP: 10, wantSD: 90},
		{name: "shield overflow", enabled: true, rate: 90, attackerType: ObjectTypePlayer, targetType: ObjectTypePlayer, targetSD: 50, wantHP: 50, wantSD: 50},
		{name: "zero shield", enabled: true, rate: 90, attackerType: ObjectTypePlayer, targetType: ObjectTypePlayer, wantHP: 100},
		{name: "negative rate", enabled: true, rate: -1, attackerType: ObjectTypePlayer, targetType: ObjectTypePlayer, targetSD: 200, wantHP: 100},
		{name: "rate over one hundred", enabled: true, rate: 101, attackerType: ObjectTypePlayer, targetType: ObjectTypePlayer, targetSD: 200, wantSD: 100},
	} {
		t.Run(tt.name, func(t *testing.T) {
			withShieldSystemTest(t, tt.enabled, tt.rate)
			attacker, _ := newSkillTestObject(1, tt.attackerType)
			target, _ := newSkillTestObject(2, tt.targetType)
			target.SD = tt.targetSD

			hpDamage, sdDamage := attacker.splitDamage(target, 100)

			if hpDamage != tt.wantHP || sdDamage != tt.wantSD {
				t.Fatalf("damage split = HP:%d SD:%d, want HP:%d SD:%d", hpDamage, sdDamage, tt.wantHP, tt.wantSD)
			}
			if target.SD != tt.targetSD {
				t.Fatalf("splitDamage changed target SD to %d, want %d", target.SD, tt.targetSD)
			}
		})
	}
}

func TestAttackAppliesUnifiedShieldDamage(t *testing.T) {
	withShieldSystemTest(t, true, 90)
	for _, tt := range []struct {
		name string
		mode attackMode
	}{
		{name: "fixed", mode: attackModeFixed},
		{name: "dot", mode: attackModeDOT},
	} {
		t.Run(tt.name, func(t *testing.T) {
			attacker, attackerActor := newSkillTestObject(1, ObjectTypePlayer)
			target, _ := newSkillTestObject(2, ObjectTypePlayer)
			target.SD, target.MaxSD = 100, 100
			withTestObjectManager(t, attacker, target)

			damage := attacker.attack(target, attackRequest{mode: tt.mode, damage: 20})

			if damage != 20 {
				t.Fatalf("attack damage = %d, want 20", damage)
			}
			if target.HP != 98 || target.SD != 82 {
				t.Fatalf("target resources = HP:%d SD:%d, want HP:98 SD:82", target.HP, target.SD)
			}
			reply := findAttackDamageReply(attackerActor.messages)
			if reply == nil {
				t.Fatal("attack damage reply was not sent")
			}
			if reply.Damage != 2 || reply.SDDamage != 18 {
				t.Fatalf("reply damage = HP:%d SD:%d, want HP:2 SD:18", reply.Damage, reply.SDDamage)
			}
		})
	}
}

func TestAttackReportsAllocatedOverkillDamage(t *testing.T) {
	withShieldSystemTest(t, true, 0)
	attacker, attackerActor := newSkillTestObject(1, ObjectTypePlayer)
	target, _ := newSkillTestObject(2, ObjectTypePlayer)
	target.HP = 5
	target.TX, target.TY = target.X, target.Y
	withTestObjectManager(t, attacker, target)

	attacker.attack(target, attackRequest{mode: attackModeFixed, damage: 20})

	if target.HP != 0 {
		t.Fatalf("target HP = %d, want 0", target.HP)
	}
	reply := findAttackDamageReply(attackerActor.messages)
	if reply == nil {
		t.Fatal("attack damage reply was not sent")
	}
	if reply.Damage != 20 || reply.SDDamage != 0 {
		t.Fatalf("reply damage = HP:%d SD:%d, want HP:20 SD:0", reply.Damage, reply.SDDamage)
	}
}

type deathTestActor struct {
	*skillTestActor
	dieCalls  int
	dieDamage int
}

func newDeathTestObject(index int, typ ObjectType) (*Object, *deathTestActor) {
	obj, actor := newSkillTestObject(index, typ)
	deathActor := &deathTestActor{skillTestActor: actor}
	obj.Objecter = deathActor
	return obj, deathActor
}

func (a *deathTestActor) Die(_ *Object, damage int) {
	a.dieCalls++
	a.dieDamage = damage
}

func findAttackDieReply(messages []any) *model.MsgAttackDieReply {
	for i := len(messages) - 1; i >= 0; i-- {
		if reply, ok := messages[i].(*model.MsgAttackDieReply); ok {
			return reply
		}
	}
	return nil
}

func TestAttackSettlesDeathOnce(t *testing.T) {
	withShieldSystemTest(t, false, 0)
	attacker, attackerActor := newAttackProcTestObject(1, ObjectTypePlayer)
	target, targetActor := newDeathTestObject(2, ObjectTypePlayer)
	attackerActor.inventory.Items[0] = item.NewItem(2, 0)
	attackerActor.inventory.Items[0].Durability = 1
	attackerActor.maceStun = 100
	target.HP = 5
	target.TX, target.TY = target.X, target.Y
	withTestObjectManager(t, attacker, target)

	firstDamage := attacker.attack(target, attackRequest{mode: attackModeFixed, damage: 20})
	attackerMessages := len(attackerActor.messages)
	targetMessages := len(targetActor.messages)
	secondDamage := attacker.attack(target, attackRequest{mode: attackModeFixed, damage: 20})

	if firstDamage != 20 || secondDamage != 0 {
		t.Fatalf("attack damage = first:%d second:%d, want first:20 second:0", firstDamage, secondDamage)
	}
	if target.HP != 0 || target.Live {
		t.Fatalf("target state = HP:%d Live:%t, want HP:0 Live:false", target.HP, target.Live)
	}
	if targetActor.dieCalls != 1 || targetActor.dieDamage != 20 {
		t.Fatalf("Die calls = %d with damage %d, want 1 with damage 20", targetActor.dieCalls, targetActor.dieDamage)
	}
	if len(attackerActor.messages) != attackerMessages || len(targetActor.messages) != targetMessages {
		t.Fatal("second attack sent messages for an already dead target")
	}
	if got := countMessages[model.MsgAttackDieReply](targetActor.messages); got != 1 {
		t.Fatalf("death replies = %d, want 1", got)
	}
	if target.HasBuff(effect.BuffStun) {
		t.Fatal("killed target received mace stun")
	}
}

func TestAttackRejectsDeadTargets(t *testing.T) {
	withShieldSystemTest(t, false, 0)
	for _, tt := range []struct {
		name  string
		mode  attackMode
		setup func(target *Object)
	}{
		{
			name: "fixed with non-live target",
			mode: attackModeFixed,
			setup: func(target *Object) {
				target.Live = false
			},
		},
		{
			name: "calculated with zero-HP target",
			mode: attackModeCalculated,
			setup: func(target *Object) {
				target.HP = 0
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			attacker, attackerActor := newSkillTestObject(1, ObjectTypePlayer)
			target, targetActor := newSkillTestObject(2, ObjectTypePlayer)
			tt.setup(target)
			withTestObjectManager(t, attacker, target)
			targetHP, targetSD := target.HP, target.SD

			damage := attacker.attack(target, attackRequest{mode: tt.mode, damage: 20})

			if damage != 0 {
				t.Fatalf("attack damage = %d, want 0", damage)
			}
			if target.HP != targetHP || target.SD != targetSD {
				t.Fatal("rejected attack changed participant resources")
			}
			if len(attackerActor.messages) != 0 || len(targetActor.messages) != 0 {
				t.Fatal("rejected attack sent messages")
			}
		})
	}
}

func TestDOTContinuesAfterSourceDeath(t *testing.T) {
	withShieldSystemTest(t, false, 0)
	source, sourceActor := newAttackProcTestObject(1, ObjectTypePlayer)
	target, _ := newSkillTestObject(2, ObjectTypePlayer)
	sourceActor.inventory.Items[0] = item.NewItem(2, 0)
	sourceActor.inventory.Items[0].Durability = 1
	sourceActor.maceStun = 100
	withTestObjectManager(t, source, target)
	if !target.addEffect(&effect.Effect{
		BuffIndex: 9001,
		Dot:       20,
		Source:    source.Index,
		Expire:    time.Now().Add(time.Minute),
		NextTick:  time.Now().Add(-time.Second),
	}) {
		t.Fatal("addEffect() = false")
	}
	source.Live = false
	source.HP = 0

	target.processEffect()

	if target.HP != 80 {
		t.Fatalf("target HP = %d, want 80", target.HP)
	}
	if !target.HasBuff(effect.BuffStun) {
		t.Fatal("dead DOT source did not stun target")
	}
}

func TestFireScreamExplosionStopsAfterSourceDeath(t *testing.T) {
	withShieldSystemTest(t, false, 0)
	source, sourceActor := newSkillTestObject(1, ObjectTypePlayer)
	center, centerActor := newSkillTestObject(2, ObjectTypeMonster)
	withTestObjectManager(t, source, center)
	if !source.addViewportObject(center) {
		t.Fatal("failed to add explosion center to source viewport")
	}
	s := learnSkillForTest(t, source, skill.SkillIndexFireScream)
	source.AddDelayMsg(8, s.Index, 0, center.Index)
	for _, msg := range source.msgs {
		if msg.code == 8 {
			msg.time = time.Now().Add(-time.Second)
			break
		}
	}
	source.Live = false
	source.HP = 0
	sourceMessages := len(sourceActor.messages)
	centerMessages := len(centerActor.messages)

	source.processDelayMsg()

	if center.HP != center.MaxHP {
		t.Fatalf("center HP = %d, want %d", center.HP, center.MaxHP)
	}
	if len(sourceActor.messages) != sourceMessages || len(centerActor.messages) != centerMessages {
		t.Fatal("dead source explosion sent attack messages")
	}
}

func TestAttackDeathReplySkill(t *testing.T) {
	withShieldSystemTest(t, false, 0)
	for _, tt := range []struct {
		name      string
		mode      attackMode
		withSkill bool
		wantSkill int
	}{
		{name: "calculated skill", mode: attackModeCalculated, withSkill: true, wantSkill: skill.SkillIndexFireBall},
		{name: "normal attack", mode: attackModeCalculated},
		{name: "fixed damage", mode: attackModeFixed},
		{name: "dot damage", mode: attackModeDOT},
	} {
		t.Run(tt.name, func(t *testing.T) {
			attacker, attackerActor := newSkillTestObject(1, ObjectTypePlayer)
			target, targetActor := newSkillTestObject(2, ObjectTypePlayer)
			target.Level = 0
			target.HP = 1
			target.TX, target.TY = target.X, target.Y
			withTestObjectManager(t, attacker, target)
			req := attackRequest{mode: tt.mode, damage: 20}
			if tt.withSkill {
				attackerActor.magicAttackMin = 100
				attackerActor.magicAttackMax = 100
				req.skill = learnSkillForTest(t, attacker, skill.SkillIndexFireBall)
			}

			attacker.attack(target, req)

			reply := findAttackDieReply(targetActor.messages)
			if reply == nil {
				t.Fatal("attack death reply was not sent")
			}
			if reply.Skill != tt.wantSkill {
				t.Fatalf("death skill = %d, want %d", reply.Skill, tt.wantSkill)
			}
		})
	}
}

func TestDamageReturnUsesExpectedDamage(t *testing.T) {
	for _, tt := range []struct {
		name         string
		attackerType ObjectType
		wantDamage   int
	}{
		{name: "player returns final damage", attackerType: ObjectTypePlayer, wantDamage: 20},
		{name: "monster returns attack max", attackerType: ObjectTypeMonster, wantDamage: 30},
	} {
		t.Run(tt.name, func(t *testing.T) {
			attacker, _ := newSkillTestObject(1, tt.attackerType)
			target, targetActor := newSkillTestObject(2, ObjectTypePlayer)
			withTestObjectManager(t, attacker, target)
			targetActor.returnDamage = 100

			attacker.attack(target, attackRequest{
				mode:   attackModeFixed,
				damage: 20,
			})

			if attacker.HP != attacker.MaxHP {
				t.Fatalf("attacker HP before delay = %d, want %d", attacker.HP, attacker.MaxHP)
			}
			var returned *DelayMsg
			for _, msg := range target.msgs {
				if msg.code == 12 {
					returned = msg
					break
				}
			}
			if returned == nil {
				t.Fatal("returned damage was not queued")
			}
			if returned.subcode != tt.wantDamage {
				t.Fatalf("queued damage = %d, want %d", returned.subcode, tt.wantDamage)
			}

			returned.time = time.Now().Add(-time.Millisecond)
			target.processDelayMsg()
			if attacker.HP != attacker.MaxHP-tt.wantDamage {
				t.Fatalf("attacker HP after delay = %d, want %d", attacker.HP, attacker.MaxHP-tt.wantDamage)
			}
		})
	}
}

func TestDamageReturnModePolicy(t *testing.T) {
	for _, tt := range []struct {
		name        string
		mode        attackMode
		wantReflect bool
		wantReturn  bool
	}{
		{name: "reflected can return", mode: attackModeReflected, wantReturn: true},
		{name: "returned can reflect and return", mode: attackModeReturned, wantReflect: true, wantReturn: true},
		{name: "dot cannot reflect or return", mode: attackModeDOT},
	} {
		t.Run(tt.name, func(t *testing.T) {
			attacker, _ := newSkillTestObject(1, ObjectTypePlayer)
			target, targetActor := newSkillTestObject(2, ObjectTypePlayer)
			withTestObjectManager(t, attacker, target)
			targetActor.returnDamage = 100
			if !target.addEffect(&effect.Effect{
				BuffIndex: effect.BuffDamageReflection,
				Reflect:   100,
				Expire:    time.Now().Add(time.Minute),
			}) {
				t.Fatal("addEffect() = false")
			}

			attacker.attack(target, attackRequest{
				mode:   tt.mode,
				damage: 20,
			})

			gotReflect := false
			gotReturn := false
			for _, msg := range target.msgs {
				switch msg.code {
				case 9:
					gotReflect = true
				case 12:
					gotReturn = true
				}
			}
			if gotReflect != tt.wantReflect {
				t.Fatalf("queued reflection = %t, want %t", gotReflect, tt.wantReflect)
			}
			if gotReturn != tt.wantReturn {
				t.Fatalf("queued return = %t, want %t", gotReturn, tt.wantReturn)
			}
		})
	}
}
