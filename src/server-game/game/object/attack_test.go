package object

import (
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/effect"
)

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
