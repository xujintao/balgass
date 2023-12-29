package object

import (
	"log"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
)

func (obj *Object) getDefense(t int) int {
	defense := 0
	switch t {
	case 1:
		defense = 0
	default:
		defense = obj.Defense
	}
	return defense
}

func (obj *Object) getDamage(t int) int {
	damage := 0
	switch t {
	case 3:
		damage = obj.AttackMax
		damage += obj.GetCriticalAttackDamage()
	case 2:
		damage = obj.AttackMax
		damage += damage * 20 / 100
		damage += obj.GetExcellentAttackDamage()
	default:
		sub := obj.AttackMax - obj.AttackMin
		if sub < 0 {
			return 0
		}
		damage = obj.AttackMin + rand.Intn(sub+1)
	}
	return damage
}

func (obj *Object) attack(tobj *Object, damage int) {
	damageType := 0
	if damage == 0 {
		// 1. think about miss

		// 2. rand ignore target defense and get target defense
		ignoreDefenseRate := obj.GetIgnoreDefenseRate()
		if rand.Intn(10000) < ignoreDefenseRate*100 {
			damageType = 1
		}
		defense := tobj.getDefense(damageType)
		// 3. rand normal/critical/excel and get object attack panel or skill attack
		criticalAttackRate := obj.GetCriticalAttackRate()
		if rand.Intn(10000) < criticalAttackRate*100 {
			damageType = 3
		}
		excellentAttackRate := obj.GetExcellentAttackRate()
		if rand.Intn(10000) < excellentAttackRate*100 {
			damageType = 2
		}
		damage = obj.getDamage(damageType)
		// 4. calc attack damage
		damage = damage - defense
		// 5. add damage
		damage += obj.GetIncreaseDamage()
		// 6. decrease damage
		damage += damage * obj.GetWingIncreaseDamage() / 100
		// 7. absorb damage
		damage -= damage * tobj.GetWingReduceDamage() / 100
		// 8. combo damage
	}
	// 9. reflect damage
	// 10. rebound damage
	// 11. rand double damage
	doubleDamageRate := obj.GetDoubleDamageRate()
	if rand.Intn(10000) < doubleDamageRate*100 {
		damage *= 2
	}
	// 12. target recover all hp/mp/sd
	// 13. mace stun
	// 14. decrease target hp
	// 15. check target hp

	// limit attack damage min
	// attackDamageMin := tobj.Level / 10
	// if attackDamageMin <= 0 {
	// 	attackDamageMin = 1
	// }
	// if attackDamage < attackDamageMin {
	// 	attackDamage = attackDamageMin
	// }
	if damage < 0 {
		damage = 1
	}
	tobj.HP -= damage
	if tobj.HP <= 0 {
		tobj.HP = 0
	}

	// Push attack damage reply
	attackDamageReply := model.MsgAttackDamageReply{
		Target:     tobj.Index,
		Damage:     damage,
		DamageType: damageType,
		SDDamage:   0,
	}
	obj.Push(&attackDamageReply)
	tobj.Push(&attackDamageReply)

	// Push attack effect reply
	attackEffectReply := model.MsgAttackEffectReply{
		Target:       tobj.Index,
		HP:           tobj.HP,
		MaxHP:        tobj.MaxHP + tobj.AddHP,
		Level:        tobj.Level,
		IceEffect:    0,
		PoisonEffect: 0,
	}
	tobj.PushViewport(&attackEffectReply)

	// Push attack hp reply
	attackHPReply := model.MsgAttackHPReply{
		Target: tobj.Index,
		MaxHP:  tobj.MaxHP + tobj.AddHP,
		HP:     tobj.HP,
	}
	tobj.PushViewport(&attackHPReply)

	// handle target die
	if tobj.HP == 0 {
		tobj.Live = false
		tobj.State = 4
		tobj.dieTime = time.Now()
		tobj.Die(obj)
		maps.MapManager.ClearMapAttrStand(tobj.MapNumber, tobj.X, tobj.Y)
		tobj.dieRegen = true

		// Push attack die reply
		attackDieReply := model.MsgAttackDieReply{
			Target: tobj.Index,
			Skill:  0,
			Killer: obj.Index,
		}
		tobj.PushViewport(&attackDieReply)
	}
	// log.Printf("attack [%d][%s]->[%d][%s] hp[%d]\n",
	// 	obj.index, obj.Annotation, tobj.index, tobj.Annotation, tobj.HP)
}

func (obj *Object) Attack(msg *model.MsgAttack) {
	tobj := ObjectManager.objects[msg.Target]
	if tobj == nil {
		log.Printf("Attack target is invalid [index]%d->[index]%d\n",
			obj.Index, msg.Target)
		return
	}
	// Push attack action to viewport
	reply := model.MsgActionReply{
		Index:  obj.Index,
		Action: msg.Action,
		Dir:    msg.Dir,
		Target: tobj.Index,
	}
	obj.PushViewport(&reply)
	obj.attack(tobj, 0)
}
