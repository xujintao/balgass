package object

import (
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/effect"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
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
