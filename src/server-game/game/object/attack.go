package object

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/effect"
	"github.com/xujintao/balgass/src/server-game/game/exp"
	"github.com/xujintao/balgass/src/server-game/game/formula"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

const (
	defaultAttackSpeedTimeLimit   = 800
	defaultDecTimePerAttackSpeed  = 5.33
	defaultMinimumAttackSpeedTime = 200
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
	reduction := obj.effects.DefenseReduction()
	defense -= defense * reduction / 100
	if defense < 0 {
		defense = 0
	}
	return defense
}

func (obj *Object) getDamage(s *skill.Skill, t int, tobj *Object) int {
	// get (physical/magic/curse/special)damage from skill type
	damageMin := 0
	damageMax := 0
	switch s.Index {
	case 0: // normal attack skill
		damageMin = obj.AttackMin
		damageMax = obj.AttackMax
	case skill.SkillIndexPoison, // 1毒咒
		skill.SkillIndexMeteorite,      // 2陨石
		skill.SkillIndexLightning,      // 3掌心雷
		skill.SkillIndexFireBall,       // 4火球
		skill.SkillIndexFlame,          // 5火龙
		skill.SkillIndexIce,            // 7冰封
		skill.SkillIndexTwister,        // 8龙卷风
		skill.SkillIndexEvilSpirit,     // 9黑龙波
		skill.SkillIndexHellFire,       // 10地狱火
		skill.SkillIndexPowerWave,      // 11真空波
		skill.SkillIndexAquaBeam,       // 12极光
		skill.SkillIndexCometFall,      // 13爆炎
		skill.SkillIndexInferno,        // 14毁灭烈焰
		skill.SkillIndexEnergyBall,     // 17能量球
		skill.SkillIndexDecay,          // 38单毒炎
		skill.SkillIndexIceStorm,       // 39暴风雪
		skill.SkillIndexLance,          // 45回旋刃(攻城)
		skill.SkillIndexDrainLife,      // 214摄魂咒
		skill.SkillIndexLightningShock, // 230烈光闪
		skill.SkillIndexGiganticStorm:  // 237闪电轰顶
		damageMin = obj.GetMagicAttackMin() + obj.effects.MagicMin() + s.DamageMin
		damageMax = obj.GetMagicAttackMax() + obj.effects.MagicMax() + s.DamageMax
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
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		damageMin = int(float64(damageMin) * obj.GetKnightGladiatorCalcSkillBonus())
		damageMax = int(float64(damageMax) * obj.GetKnightGladiatorCalcSkillBonus())
	case skill.SkillIndexTripleShot, // 24多重箭(武器)
		skill.SkillIndexStarfall,    // 46天堂之箭(攻城)
		skill.SkillIndexIceArrow,    // 51冰封箭
		skill.SkillIndexPenetration, // 52穿透箭
		skill.SkillIndexMultiShot:   // 235五重箭
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		formula.Elf_CalcSkillBonus(damageMin, obj.GetEnergy(), &damageMin)
		formula.Elf_CalcSkillBonus(damageMax, obj.GetEnergy(), &damageMax)
	case skill.SkillIndexForce, // 60冲击(初始)
		skill.SkillIndexFireBurst,     // 61星云火链
		skill.SkillIndexEarthshake,    // 62地裂(黑王马)
		skill.SkillIndexElectricSpike, // 65圣极光
		skill.SkillIndexForceWave,     // 66冲击波
		skill.SkillIndexFireScream:    // 78火舞旋风
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		formula.Lord_CalcSkillBonus(damageMin, obj.GetEnergy(), &damageMin)
		formula.Lord_CalcSkillBonus(damageMax, obj.GetEnergy(), &damageMax)
	case skill.SkillIndexSummonerExplosion, // 223爆裂
		skill.SkillIndexRequiem,   // 224刺袭
		skill.SkillIndexPollution: // 225污染
		damageMin = obj.GetCurseAttackMin() + obj.effects.CurseMin() + s.DamageMin
		damageMax = obj.GetCurseAttackMax() + obj.effects.CurseMax() + s.DamageMax
	case skill.SkillIndexNova: // 40星辰一怒
		chargeTable := [...]int{0, 20, 50, 99, 160, 225, 325, 425, 550, 700, 880, 1090, 1320}
		count := obj.novaSkill.count
		if count < 0 || count >= len(chargeTable) {
			count = 0
		}
		base := obj.GetStrength()/2 + chargeTable[count]
		damageMin = base + obj.GetMagicAttackMin()
		damageMax = base + obj.GetMagicAttackMax()
	case skill.SkillIndexImpale: // 47钻云枪
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		damageMin = int(float64(damageMin) * obj.GetImpaleSkillCalc())
		damageMax = int(float64(damageMax) * obj.GetImpaleSkillCalc())
	case skill.SkillIndexPowerSlash: // 56天雷闪(武器)
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		formula.GladiatorPowerSlash(damageMin, obj.GetEnergy(), &damageMin)
		formula.GladiatorPowerSlash(damageMax, obj.GetEnergy(), &damageMax)
	case skill.SkillIndexPlasmaStorm: // 76闪电链(炎狼兽)
		switch obj.Class {
		case 0, 5:
			damageMin = obj.GetStrength()/5 + obj.GetDexterity()/5 + obj.GetVitality()/7 + obj.GetEnergy()/3
		case 1, 3, 6:
			damageMin = obj.GetStrength()/3 + obj.GetDexterity()/5 + obj.GetVitality()/5 + obj.GetEnergy()/7
		case 2:
			damageMin = obj.GetStrength()/5 + obj.GetDexterity()/3 + obj.GetVitality()/7 + obj.GetEnergy()/5
		default:
			damageMin = obj.GetStrength()/5 + obj.GetDexterity()/5 + obj.GetVitality()/7 + obj.GetEnergy()/3 + obj.GetLeadership()/3
		}
		damageMax = damageMin
		damageMin += s.DamageMin
		damageMax += s.DamageMax
		formula.FenrirSkillCalc(damageMin, obj.Level, obj.GetMasterLevel(), &damageMin)
		formula.FenrirSkillCalc(damageMax, obj.Level, obj.GetMasterLevel(), &damageMax)
	case skill.SkillIndexChainLightning: // 215链雷咒
		damageMin = obj.GetMagicAttackMin() + s.DamageMin
		damageMax = obj.GetMagicAttackMax() + s.DamageMax
		targetNumber := obj.chainLightningTarget
		if targetNumber < 1 || targetNumber > 3 {
			targetNumber = 1
		}
		formula.ChainLightningCalc(damageMin, targetNumber, &damageMin)
		formula.ChainLightningCalc(damageMax, targetNumber, &damageMax)
	case skill.SkillIndexStrikeDestruction: // 232破坏一击
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		formula.StrikeOfDestructionCalc(damageMin, obj.GetEnergy(), &damageMin)
		formula.StrikeOfDestructionCalc(damageMax, obj.GetEnergy(), &damageMax)
	case skill.SkillIndexFlameStrike: // 236火剑袭
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		formula.FlameStrikeCalc(damageMin, &damageMin)
		formula.FlameStrikeCalc(damageMax, &damageMax)
	case skill.SkillIndexChaoticDiseier: // 238黑暗之力
		damageMin = obj.GetMagicAttackMin() + s.DamageMin
		damageMax = obj.GetMagicAttackMax() + s.DamageMax
		formula.ChaoticDiseierCalc(damageMin, obj.GetEnergy(), &damageMin)
		formula.ChaoticDiseierCalc(damageMax, obj.GetEnergy(), &damageMax)
	case skill.SkillIndexKillingBlow: // 260幽冥青狼拳
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		formula.RageFighterKillingBlow(damageMin, obj.GetVitality(), &damageMin)
		formula.RageFighterKillingBlow(damageMax, obj.GetVitality(), &damageMax)
	case skill.SkillIndexBeastUppercut: // 261斗气爆裂拳
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		formula.RageFighterBeastUppercut(damageMin, obj.GetVitality(), &damageMin)
		formula.RageFighterBeastUppercut(damageMax, obj.GetVitality(), &damageMax)
	case skill.SkillIndexChainDrive: // 262回旋踢
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		formula.RageFighterChainDrive(damageMin, obj.GetVitality(), &damageMin)
		formula.RageFighterChainDrive(damageMax, obj.GetVitality(), &damageMax)
	case skill.SkillIndexDarkSide: // 263幽冥光速拳
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		formula.RageFighterDarkSideIncDamage(damageMin, obj.GetDexterity(), obj.GetEnergy(), &damageMin)
		formula.RageFighterDarkSideIncDamage(damageMax, obj.GetDexterity(), obj.GetEnergy(), &damageMax)
	case skill.SkillIndexDragonRoar: // 264炎龙拳
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		formula.RageFighterDragonRoar(damageMin, obj.GetEnergy(), &damageMin)
		formula.RageFighterDragonRoar(damageMax, obj.GetEnergy(), &damageMax)
	case skill.SkillIndexDragonSlasher: // 265噬血之龙
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		targetType := int(ObjectTypeMonster)
		if tobj != nil {
			targetType = int(tobj.Type)
		}
		formula.RageFighterDragonSlasher(damageMin, obj.GetEnergy(), targetType, &damageMin)
		formula.RageFighterDragonSlasher(damageMax, obj.GetEnergy(), targetType, &damageMax)
	case skill.SkillIndexCharge: // 269冲锋(攻城)
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		formula.RageFighterCharge(damageMin, obj.GetVitality(), &damageMin)
		formula.RageFighterCharge(damageMax, obj.GetVitality(), &damageMax)
	case skill.SkillIndexPhoenixShot: // 270神圣气旋
		damageMin = obj.AttackMin + s.DamageMin
		damageMax = obj.AttackMax + s.DamageMax
		formula.RageFighterPhoenixShot(damageMin, obj.GetVitality(), &damageMin)
		formula.RageFighterPhoenixShot(damageMax, obj.GetVitality(), &damageMax)
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
		damage += obj.effects.CriticalDamage()
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

type attackMode uint8

const (
	attackModeInvalid attackMode = iota
	attackModeCalculated
	attackModeFixed
	attackModeReflected
	attackModeReturned
	attackModeDOT
)

type attackRequest struct {
	mode   attackMode
	skill  *skill.Skill
	damage int
}

func (obj *Object) attack(tobj *Object, req attackRequest) int {
	needCalculate := false
	allowReflect := false
	allowReturn := false
	switch req.mode {
	case attackModeCalculated:
		needCalculate = true
		allowReflect = true
		allowReturn = true
	case attackModeFixed:
		allowReflect = true
		allowReturn = true
	case attackModeReflected:
		allowReturn = true
	case attackModeReturned:
		allowReflect = true
		allowReturn = true
	case attackModeDOT:
	default:
		return 0
	}

	s := req.skill
	damage := req.damage
	damageType := 0
	if needCalculate && !obj.CheckMiss(tobj) {
		// treat normal attack as skill 0
		if s == nil {
			s = skill.Skill0
		}

		// 1. calc target defense
		// rand ignore target defense and get target defense
		ignoreDefenseRate := obj.GetIgnoreDefenseRate() + obj.effects.IgnoreDefense()
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
		damage = obj.getDamage(s, damageType, tobj)
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
		// 8. angel reduce damage
		damage -= damage * tobj.GetHelperReduceDamage() / 100
		// 9. pet increase/reduce damage
		damage += damage * obj.GetPetIncreaseDamage() / 100
		damage -= damage * tobj.GetPetReduceDamage() / 100
		// 10. effect reduce damage
		reduction := obj.effects.AttackReduction()
		damage -= damage * reduction / 100
		if barrier := tobj.effect(effect.BuffSoulBarrier); barrier != nil && damage > 0 {
			mana := tobj.MP * barrier.ManaRate / 1000
			if mana < tobj.MP {
				tobj.MP -= mana
				damage -= damage * barrier.DamageReduction / 100
				tobj.PushMPAG(tobj.MP, tobj.AG)
			}
		}
		if damage <= 0 {
			damage = 0
		}
	}
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
	if damage > 0 {
		tobj.removeSleepEffect()
		obj.applySkillEffect(tobj, s, damage)
		if allowReflect {
			reflect := tobj.effects.Reflect()
			if reflect > 0 && obj.Live {
				reflectedDamage := damage * reflect / 100
				if reflectedDamage > 0 {
					tobj.AddDelayMsg(9, reflectedDamage, 10, obj.Index)
				}
			}
		}
		if allowReturn && tobj.Type == ObjectTypePlayer {
			returnRate := tobj.GetReturnDamageRate()
			if returnRate > 0 && rand.Intn(10000) < returnRate*100 {
				returnedDamage := 0
				switch obj.Type {
				case ObjectTypePlayer:
					returnedDamage = damage
				case ObjectTypeMonster:
					returnedDamage = obj.AttackMax
				}
				if returnedDamage > 0 {
					tobj.AddDelayMsg(12, returnedDamage, 10, obj.Index)
				}
			}
		}
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
		tobj.removeAllEffects()
		tobj.durationSkill = durationSkillState{}
		tobj.novaSkill = novaSkillState{}
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
		if tobj.IsSummon() {
			ObjectManager.DeleteCallMonster(tobj.Index)
		} else if tobj.summonIndex >= 0 {
			ObjectManager.DeleteCallMonster(tobj.summonIndex)
		}
	}
	// debug
	if conf.ServerEnv.Debug {
		slog.Debug("attack",
			"index", obj.Index, "name", obj.Name,
			"target", tobj.Index, "name", tobj.Name,
			"hp", tobj.HP)
	}
	return damage
}

func (obj *Object) basicAttackTarget(msg *model.MsgAttack) *Object {
	if msg == nil ||
		obj.ConnectState != ConnectStatePlaying ||
		!obj.Live ||
		obj.State == 4 ||
		obj.State == 8 ||
		obj.cannotAct() ||
		msg.Target < 0 ||
		msg.Target >= len(ObjectManager.objects) {
		return nil
	}
	tobj := ObjectManager.objects[msg.Target]
	if tobj == nil ||
		tobj == obj ||
		!tobj.Live ||
		tobj.State == 4 ||
		tobj.State == 8 ||
		tobj.MapNumber != obj.MapNumber {
		return nil
	}
	for _, vp := range obj.Viewports {
		if vp.State != 0 && vp.Type == int(tobj.Type) && vp.Number == tobj.Index {
			return tobj
		}
	}
	return nil
}

func basicAttackDelay(speed int) time.Duration {
	timeLimit := conf.CommonServer.GameServerInfo.AttackSpeedTimeLimit
	if timeLimit <= 0 {
		timeLimit = defaultAttackSpeedTimeLimit
	}
	decrement := conf.CommonServer.GameServerInfo.DecTimePerAttackSpeed
	if decrement <= 0 {
		decrement = defaultDecTimePerAttackSpeed
	}
	minimum := conf.CommonServer.GameServerInfo.MinimumAttackSpeedTime
	if minimum <= 0 {
		minimum = defaultMinimumAttackSpeedTime
	}
	if speed < 0 {
		speed = 0
	}
	delay := float64(timeLimit) - float64(speed)*decrement
	if delay < float64(minimum) {
		delay = float64(minimum)
	}
	return time.Duration(delay * float64(time.Millisecond))
}

func (obj *Object) basicAttackBlockedInSafeZone(tobj *Object) bool {
	if !conf.Common.AntiHack.EnableBlockAttackInSafeZone {
		return false
	}
	return maps.MapManager.GetMapAttr(obj.MapNumber, obj.X, obj.Y)&1 != 0 ||
		maps.MapManager.GetMapAttr(tobj.MapNumber, tobj.X, tobj.Y)&1 != 0
}

func (obj *Object) consumeBasicAttackAmmo() bool {
	if obj.Type != ObjectTypePlayer || obj.HasBuff(effect.BuffInfinityArrow) {
		return true
	}
	leftHand := obj.GetInventoryItem(0)
	rightHand := obj.GetInventoryItem(1)
	position := -1
	switch {
	case rightHand != nil && rightHand.KindB == item.KindBCrossbow:
		if leftHand == nil || leftHand.Code != item.Code(4, 7) || leftHand.Durability <= 0 {
			return false
		}
		position = 0
	case leftHand != nil && leftHand.KindB == item.KindBBow:
		if rightHand == nil || rightHand.Code != item.Code(4, 15) || rightHand.Durability <= 0 {
			return false
		}
		position = 1
	default:
		return true
	}
	ammo := obj.GetInventoryItem(position)
	ammo.Durability--
	obj.Push(&model.MsgItemDurabilityReply{
		Position:   position,
		Durability: ammo.Durability,
		Flag:       0,
	})
	return true
}

func (obj *Object) Attack(msg *model.MsgAttack) {
	tobj := obj.basicAttackTarget(msg)
	if tobj == nil || obj.basicAttackBlockedInSafeZone(tobj) {
		return
	}
	now := time.Now()
	if !obj.lastBasicAttackTime.IsZero() &&
		now.Before(obj.lastBasicAttackTime.Add(basicAttackDelay(obj.GetAttackSpeedForDelay()))) {
		return
	}
	obj.lastBasicAttackTime = now
	// Push attack action to viewport
	reply := model.MsgActionReply{
		Index:  obj.Index,
		Action: msg.Action,
		Dir:    msg.Dir,
		Target: tobj.Index,
	}
	obj.PushViewport(&reply)
	if !obj.consumeBasicAttackAmmo() {
		return
	}
	obj.attack(tobj, attackRequest{mode: attackModeCalculated})
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
