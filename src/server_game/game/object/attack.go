package object

import (
	"log"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
)

func (obj *object) getAttackPanel() int {
	sub := obj.attackPanelMax - obj.attackPanelMin
	if sub < 0 {
		log.Printf("attack panel is invalid [index]%d [class]%d\n",
			obj.index, obj.Class)
		return 0
	}
	attackDamage := obj.attackPanelMin + rand.Intn(sub+1)
	return attackDamage
}

func (obj *object) attack(tobj *object) {
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
	defense := tobj.defense
	attackDamage := attackPanel - defense

	// limit attack damage min
	attackDamageMin := tobj.Level / 10
	if attackDamageMin <= 0 {
		attackDamageMin = 1
	}
	if attackDamage < attackDamageMin {
		attackDamage = attackDamageMin
	}
	tobj.HP -= attackDamage
	if tobj.HP <= 0 {
		tobj.HP = 0
		tobj.Live = false
		tobj.State = 4
		maps.MapManager.ClearMapAttrStand(tobj.MapNumber, tobj.X, tobj.Y)
		tobj.dieRegen = true
		tobj.regenTime = time.Duration(time.Now().Unix())

		// push attack die reply
		attackDieReply := model.MsgAttackDieReply{
			Target: tobj.index,
			Skill:  0,
			Killer: obj.index,
		}
		tobj.pushViewport(&attackDieReply)
	}

	// push attack reply
	attackReply := model.MsgAttackReply{
		Index:      tobj.index,
		Damage:     attackDamage,
		DamageType: 0,
		SDDamage:   0,
	}
	obj.push(&attackReply)
	tobj.push(&attackReply)

	// push attack effect reply
	attackEffectReply := model.MsgAttackEffectReply{
		Target:       tobj.index,
		HP:           tobj.HP,
		MaxHP:        tobj.MaxHP + tobj.AddHP,
		Level:        tobj.Level,
		IceEffect:    0,
		PoisonEffect: 0,
	}
	tobj.pushViewport(&attackEffectReply)

	// push attack hp reply
	attackHPReply := model.MsgAttackHPReply{
		Target: tobj.index,
		MaxHP:  tobj.MaxHP + tobj.AddHP,
		HP:     tobj.HP,
	}
	tobj.pushViewport(&attackHPReply)

	// log.Printf("attack [%d][%s]->[%d][%s] hp[%d]\n",
	// 	obj.index, obj.Annotation, tobj.index, tobj.Annotation, tobj.HP)
}

func (obj *object) Attack(msg *model.MsgAttack) {
	tobj := obj.objectManager.objects[msg.Target]
	if tobj == nil {
		log.Printf("Attack target is invalid [index]%d->[index]%d\n",
			obj.index, msg.Target)
		return
	}
	// push attack action to viewport
	reply := model.MsgActionReply{
		Index:  obj.index,
		Action: msg.Action,
		Dir:    msg.Dir,
		Target: tobj.index,
	}
	obj.pushViewport(&reply)
	obj.attack(tobj)
}

func (obj *object) SkillAttack(msg *model.MsgSkillAttack) {

}
