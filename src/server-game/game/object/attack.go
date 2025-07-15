package object

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/exp"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

func (obj *Object) CheckMiss(tobj *Object) bool {

	if obj.Type == ObjectTypePlayer && tobj.Type == ObjectTypePlayer {
		// pvp
		attackLevel := obj.Level + obj.GetMasterLevel()
		defenseLevel := tobj.Level + tobj.GetMasterLevel()
		attackRate := obj.GetAttackRatePVP()
		defenseRate := tobj.GetDefenseRatePVP()
		expressionA := attackRate * 100 / (attackRate + defenseRate)
		expressionB := attackLevel * 100 / (attackLevel + defenseLevel)
		rate := expressionA * expressionB / 100
		switch {
		case defenseLevel-attackLevel >= 100:
			rate -= 5
		case defenseLevel-attackLevel >= 200:
			rate -= 10
		case defenseLevel-attackLevel >= 300:
			rate -= 15
		}
		if rand.Intn(100) > rate {
			return true
		}
	} else {
		// pve
		attackRate := obj.AttackRate
		defenseRate := tobj.DefenseRate
		if attackRate <= 0 {
			attackRate = 1
		}
		if defenseRate <= 0 {
			defenseRate = 1
		}
		if attackRate < defenseRate {
			if rand.Intn(100) >= 5 {
				return true
			}
		} else {
			if rand.Intn(attackRate) < defenseRate {
				return true
			}
		}
	}
	return false
}

func (obj *Object) getDefense(attacker *Object, t int) int {
	defense := 0
	switch t {
	case 1:
		defense = 0
	default:
		defense = obj.Defense
		if obj.Type == ObjectTypePlayer {
			if attacker.Type == ObjectTypePlayer {
				// pvp
			} else {
				defense /= 2 // pve
			}
		}
	}
	return defense
}

func (obj *Object) getDamage(s *skill.Skill, t int) int {
	// get (physical/magic/curse/special)damage from skill type
	damageMin := 0
	damageMax := 0
	switch s.Index {
	case 0: // normal attack skill
		damageMin = obj.AttackMin
		damageMax = obj.AttackMax
	case skill.SkillIndexFallingSlash, // 19地裂斩(武器)
		skill.SkillIndexLunge,             // 20牙突刺(武器)
		skill.SkillIndexUppercut,          // 21升龙击(武器)
		skill.SkillIndexCyclone,           // 22旋风斩(武器)
		skill.SkillIndexSlash,             // 23天地十字剑(武器)
		skill.SkillIndexTwistingSlash,     // 41霹雳回旋斩
		skill.SkillIndexRagefulBlow,       // 42雷霆裂闪
		skill.SkillIndexDeathStab,         // 43袭风刺
		skill.SkillIndexCrescentMoonSlash, // 44半月斩(攻城)
		skill.SkillIndexFireBreath,        // 49流星焰(彩云兽)
		skill.SkillIndexFireSlash,         // 55玄月斩
		skill.SkillIndexSpiralSlash:       // 57风舞回旋斩(攻城)
		damageMin = obj.AttackMin
		damageMax = obj.AttackMax
		damageMin = int(float64(damageMin) * obj.GetKnightGladiatorCalcSkillBonus())
		damageMax = int(float64(damageMax) * obj.GetKnightGladiatorCalcSkillBonus())
	case skill.SkillIndexImpale: // 47钻云枪
		damageMin = obj.AttackMin
		damageMax = obj.AttackMax
		damageMin = int(float64(damageMin) * obj.GetImpaleSkillCalc())
		damageMax = int(float64(damageMax) * obj.GetImpaleSkillCalc())
	default:
		damageMin = obj.AttackMin
		damageMax = obj.AttackMax
	}

	// get damage from damage type(normal/critical/excellent)
	damage := 0
	switch t {
	case 3:
		damage = damageMax
		damage += obj.GetCriticalAttackDamage()
	case 2:
		damage = damageMax
		damage += damage * 20 / 100
		damage += obj.GetExcellentAttackDamage()
	default:
		sub := damageMax - damageMin
		if sub < 0 {
			return 0
		}
		damage = damageMin + rand.Intn(sub+1)
	}
	return damage
}

func (obj *Object) attack(tobj *Object, s *skill.Skill, damage int) {
	damageType := 0
	if damage == 0 && !obj.CheckMiss(tobj) {
		if s == nil {
			s = skill.Skill0
		}

		// 1. calc target defense
		// rand ignore target defense and get target defense
		ignoreDefenseRate := obj.GetIgnoreDefenseRate()
		if rand.Intn(10000) < ignoreDefenseRate*100 {
			damageType = 1
		}
		defense := tobj.getDefense(obj, damageType)
		// 2. calc object skill damage
		// rand normal/critical/excel and get object attack panel or skill attack
		criticalAttackRate := obj.GetCriticalAttackRate()
		if rand.Intn(10000) < criticalAttackRate*100 {
			damageType = 3
		}
		excellentAttackRate := obj.GetExcellentAttackRate()
		if rand.Intn(10000) < excellentAttackRate*100 {
			damageType = 2
		}
		// normal attack --> physical attack
		// skill attack --> physical/magic/curse attack
		damage = obj.getDamage(s, damageType)
		// 3. calc attack damage
		damage = damage - defense
		if damage < 0 {
			damage = 0
		}
		// 4. add damage
		damage += obj.GetAddDamage()
		// 5. premium scroll damage

		// 6. armor reduce damage
		damage -= damage * tobj.GetArmorReduceDamage() / 100
		// 7. wing increase/reduce damage
		damage += damage * obj.GetWingIncreaseDamage() / 100
		damage -= damage * tobj.GetWingReduceDamage() / 100
		// 8. helper reduce damage
		damage -= damage * tobj.GetHelperReduceDamage() / 100
		// 9. pet reduce damage
		damage += damage * obj.GetPetIncreaseDamage() / 100
		damage -= damage * tobj.GetPetReduceDamage() / 100
		if damage <= 0 {
			damage = 0
		}
	}
	// 9. reflect damage
	// 10. return damage
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
		MaxHP:        tobj.MaxHP,
		Level:        tobj.Level,
		IceEffect:    0,
		PoisonEffect: 0,
	}
	tobj.PushViewport(&attackEffectReply)

	// Push attack hp reply
	attackHPReply := model.MsgAttackHPReply{
		Target: tobj.Index,
		MaxHP:  tobj.MaxHP,
		HP:     tobj.HP,
	}
	tobj.PushViewport(&attackHPReply)

	// handle target die
	if tobj.HP == 0 {
		tobj.Live = false
		tobj.State = 4
		tobj.dieTime = time.Now()
		tobj.Die(obj, damage)
		maps.MapManager.ClearMapAttrStand(tobj.MapNumber, tobj.TX, tobj.TY)
		tobj.dieRegen = true

		// Push attack die reply
		attackDieReply := model.MsgAttackDieReply{
			Target: tobj.Index,
			Skill:  0,
			Killer: obj.Index,
		}
		tobj.PushViewport(&attackDieReply)
	}
	// debug
	if conf.ServerEnv.Debug {
		slog.Debug("attack",
			"index", obj.Index, "name", obj.Name,
			"target", tobj.Index, "name", tobj.Name,
			"hp", tobj.HP)
	}
}

func (obj *Object) Attack(msg *model.MsgAttack) {
	tobj := ObjectManager.objects[msg.Target]
	if tobj == nil {
		slog.Error("Attack target is nil", "index", obj.Index, "target", msg.Target)
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
	obj.attack(tobj, nil, 0)
}

func (obj *Object) DieGiveExperience(tobj *Object, damage int) {
	targetLevel := tobj.Level + tobj.GetMasterLevel()
	level := (obj.Level + 25) * obj.Level / 3
	if obj.Level+10 < targetLevel {
		level = level * (obj.Level + 10) / targetLevel
	}
	if obj.Level >= 65 {
		level += (obj.Level - 64) * obj.Level / 4
	}
	baseExp := 0
	maxExp := 0
	if level > 0 {
		maxExp = level / 2
	} else {
		level = 0
	}
	if maxExp < 1 {
		baseExp = level
	} else {
		baseExp = maxExp/2 + level
	}
	if baseExp <= 0 {
		return
	}
	obj.MoneyDrop = baseExp
	var mapBonus float64
	var baseBonus float64
	if !tobj.IsMasterLevel() {
		mapBonus = maps.MapManager.GetExpBonus(obj.MapNumber)
		baseBonus = exp.ExpManager.Normal
	} else {
		mapBonus = maps.MapManager.GetMasterExpBonus(obj.MapNumber)
		baseBonus = exp.ExpManager.Master
	}
	addexp := int(float64(baseExp) * (1 + mapBonus) * baseBonus)
	if !tobj.LevelUp(addexp) {
		reply := model.MsgExperienceReply{
			Number:     obj.Index,
			Experience: addexp,
			Damage:     damage,
		}
		tobj.Push(&reply)
	}
}

func (obj *Object) DieDropExperience() {
	slog.Debug("object DieDropExperience placeholder")
}

func (obj *Object) DieRecoverHPMP(tobj *Object) {
	if tobj.GetMonsterDieGetHP() != 0 {
		tobj.HP += int(float64(tobj.MaxHP) * tobj.GetMonsterDieGetHP())
		if tobj.HP >= tobj.MaxHP {
			tobj.HP = tobj.MaxHP
		}
		tobj.PushHPSD(tobj.HP, tobj.SD)
	}
	if tobj.GetMonsterDieGetMP() != 0 {
		tobj.MP += int(float64(tobj.MaxMP) * tobj.GetMonsterDieGetMP())
		if tobj.MP >= tobj.MaxMP {
			tobj.MP = tobj.MaxMP
		}
		tobj.PushMPAG(tobj.MP, tobj.AG)
	}
}
