package object

import (
	"log"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
)

func (obj *Object) getAttackPanel() int {
	sub := obj.AttackMax - obj.AttackMin
	if sub < 0 {
		log.Printf("attack panel is invalid [index]%d [class]%d\n",
			obj.Index, obj.Class)
		return 0
	}
	attackDamage := obj.AttackMin + rand.Intn(sub+1)
	return attackDamage
}

func (obj *Object) attack(tobj *Object) {
	// if attackDamage == 0 {
	// 	1. think about miss
	// 	2. rand normal/critical/excel and get object attack panel or skill attack
	// 	3. rand ignore target defense and get target defense
	// 	4. calc attack damage
	// 	5. add damage
	// 	6. decrease damage
	// 	7. absorb damage
	// 	8. combo damage
	// }
	// 9. reflect damage
	// 10. rebound damage
	// 11. rand double damage
	// 12. target recover all hp/mp/sd
	// 13. mace stun
	// 14. decrease target hp
	// 15. check target hp

	attackPanel := obj.getAttackPanel()
	defense := tobj.Defense
	attackDamage := attackPanel - defense

	// limit attack damage min
	// attackDamageMin := tobj.Level / 10
	// if attackDamageMin <= 0 {
	// 	attackDamageMin = 1
	// }
	// if attackDamage < attackDamageMin {
	// 	attackDamage = attackDamageMin
	// }
	if attackDamage < 0 {
		attackDamage = 1
	}
	tobj.HP -= attackDamage
	if tobj.HP <= 0 {
		tobj.HP = 0
	}

	// Push attack damage reply
	attackDamageReply := model.MsgAttackDamageReply{
		Target:     tobj.Index,
		Damage:     attackDamage,
		DamageType: 0,
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
	obj.attack(tobj)
}
