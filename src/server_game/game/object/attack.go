package object

import (
	"log"
	"math/rand"

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
	tobj.HP -= attackDamage
}

func (obj *object) Attack(msg *model.MsgAttack) {
	tobj := obj.objectManager.objects[msg.Target]
	if tobj == nil {
		log.Printf("Attack target is invalid [index]%d->[index]%d\n",
			obj.index, msg.Target)
		return
	}
	// push attack action to viewport
	obj.attack(tobj)
}

func (obj *object) SkillAttack(msg *model.MsgSkillAttack) {

}
