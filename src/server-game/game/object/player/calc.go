package player

import (
	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/buff"
	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/formula"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/model"
)

func (p *Player) Calc() {
	p.addStrength = 0
	p.addDexterity = 0
	p.addVitality = 0
	p.addEnergy = 0
	p.addLeadership = 0
	p.leftAttackMin = 0
	p.leftAttackMax = 0
	p.rightAttackMin = 0
	p.rightAttackMax = 0
	p.AttackMin = 0
	p.AttackMax = 0
	p.magic = 0
	p.magicAttackMin = 0
	p.magicAttackMax = 0
	p.curse = 0
	p.curseAttackMin = 0
	p.curseAttackMax = 0
	p.Defense = 0
	p.MagicDefense = 0
	p.AttackRate = 0
	p.attackRatePVP = 0
	p.attackBonusPVP = 0
	p.DefenseRate = 0
	p.defenseRatePVP = 0
	p.defenseBonusPVP = 0
	p.harmonyIncSDRate = 0
	p.harmonyDecSDRate = 0
	p.harmonyIgnoreSDRate = 0
	p.AttackSpeed = 0
	p.magicSpeed = 0
	p.MaxHP = 0
	p.MaxMP = 0
	p.MaxSD = 0
	p.MaxAG = 0
	p.knightGladiatorCalcSkillBonus = 0
	p.impaleSkillCalc = 0

	p.addSkillAttack = 0
	p.setAddDamage = 0
	p.setIncreaseTwoHandWeaponDamage = 0
	p.criticalAttackRate = 0
	p.criticalAttackDamage = 0
	p.excellentAttackRate = 0
	p.excellentAttackDamage = 0
	p.doubleDamageRate = 0
	p.ignoreDefenseRate = 0
	p.monsterDieGetHP = 0
	p.monsterDieGetMP = 0
	p.monsterDieGetMoney = 0
	p.recoverHP = 0
	p.recoverMP = 0
	p.recoverAG = 0
	p.agReduce = 0
	p.armorReduceDamage = 0
	p.armorReflectDamage = 0
	p.returnDamage = 0
	p.recoverMaxHP = 0
	p.recoverMaxMP = 0
	p.autoRecoverSD = false
	p.autoRecoverSDRate = 0
	p.BuffDefenseReduction = 0
	p.BuffAttackReduction = 0
	p.BuffReflect = 0

	p.wingIncreaseDamage = 0
	p.wingReduceDamage = 0
	p.helperReduceDamage = 0
	p.petIncreaseDamage = 0
	p.petReduceDamage = 0
	p.calcBaseStats()
	p.calcNormalItem()
	p.calcExcellentItem()
	p.calcHarmonyItem()
	p.calc380Item()
	p.calcSocketItem()
	p.calcSetItem(false)
	p.convertAttack()
	p.calcHelperItem()
	p.calcBuff()
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
	if p.MP > p.MaxMP {
		p.MP = p.MaxMP
	}

	// Push
	p.Push(&model.MsgAttackSpeedReply{
		AttackSpeed: p.AttackSpeed,
		MagicSpeed:  p.magicSpeed,
	})
	p.PushMaxHPSD(p.MaxHP, p.MaxSD)
	p.PushMaxMPAG(p.MaxMP, p.MaxAG)
	p.PushHPSD(p.HP, p.SD)
	p.PushMPAG(p.MP, p.AG)
}

func (p *Player) calcBaseStats() {
	leftHand := p.inventory.Items[0]
	rightHand := p.inventory.Items[1]
	wing := p.inventory.Items[7]

	// Earth seeds modify attributes used by the base character formulas.
	item.SocketManager.ApplySocketItemOption(&p.inventory, func(option item.ItemSocketOption, _ int) {
		p.addVitality += option.Vitality
	})
	// wing item contribution
	if wing != nil && wing.Wing2ExcellentLeadership {
		p.addLeadership += 10 + 5*wing.Level
	}
	// set item contribution
	p.calcSetItem(true)
	// master skill contribution
	// todo: implement master skill contribution
	strength := p.strength + p.addStrength
	dexterity := p.dexterity + p.addDexterity
	vitality := p.vitality + p.addVitality
	energy := p.energy + p.addEnergy
	leadership := p.leadership + p.addLeadership
	level := p.Level + p.masterLevel

	// base attack
	switch class.Class(p.Class) {
	case class.Wizard:
		formula.WizardDamageCalc(strength, dexterity, vitality, energy, &p.leftAttackMin, &p.rightAttackMin, &p.leftAttackMax, &p.rightAttackMax)
		formula.WizardMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.Knight:
		formula.KnightDamageCalc(strength, dexterity, vitality, energy, &p.leftAttackMin, &p.rightAttackMin, &p.leftAttackMax, &p.rightAttackMax)
		formula.KnightMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.Elf:
		if leftHand != nil && leftHand.KindB == item.KindBCrossbow ||
			rightHand != nil && rightHand.KindB == item.KindBBow {
			formula.ElfWithBowDamageCalc(strength, dexterity, vitality, energy, &p.leftAttackMin, &p.rightAttackMin, &p.leftAttackMax, &p.rightAttackMax)
		} else {
			formula.ElfWithoutBowDamageCalc(strength, dexterity, vitality, energy, &p.leftAttackMin, &p.rightAttackMin, &p.leftAttackMax, &p.rightAttackMax)
		}
		formula.ElfMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.Magumsa:
		formula.GladiatorDamageCalc(strength, dexterity, vitality, energy, &p.leftAttackMin, &p.rightAttackMin, &p.leftAttackMax, &p.rightAttackMax)
		formula.GladiatorMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.DarkLord:
		formula.LordDamageCalc(strength, dexterity, vitality, energy, leadership, &p.leftAttackMin, &p.rightAttackMin, &p.leftAttackMax, &p.rightAttackMax)
		formula.LordMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.Summoner:
		formula.SummonerDamageCalc(strength, dexterity, vitality, energy, &p.leftAttackMin, &p.rightAttackMin, &p.leftAttackMax, &p.rightAttackMax)
		formula.SummonerMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax, &p.curseAttackMin, &p.curseAttackMax)
	case class.RageFighter:
		formula.RageFighterDamageCalc(strength, dexterity, vitality, energy, &p.leftAttackMin, &p.rightAttackMin, &p.leftAttackMax, &p.rightAttackMax)
		formula.RageFighterMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	}

	// base defense
	formula.CalcDefense(p.Class, dexterity, &p.Defense)

	// base attack/defense success rate
	formula.CalcAttackSuccessRate_PvM(p.Class, strength, dexterity, leadership, p.Level, &p.AttackRate)
	formula.CalcAttackSuccessRate_PvP(p.Class, dexterity, p.Level, &p.attackRatePVP)
	formula.CalcDefenseSuccessRate_PvM(p.Class, dexterity, &p.DefenseRate)
	formula.CalcDefenseSuccessRate_PvP(p.Class, dexterity, p.Level, &p.defenseRatePVP)

	// base speed
	formula.CalcAttackSpeed(p.Class, dexterity, &p.AttackSpeed, &p.magicSpeed)

	// Stat Specialization
	// calc Stat Specialization: increase attack power
	// http://muonline.webzen.com/guides/219/1976/season-9/season-9-character-renewal
	// - (Bonus stat generated by equipping items like weapon, armor, wing, or master skill,
	// etc is not applied to specialization calculation.)
	var options []*model.MsgStatSpec
	var percent float64
	// STAT_OPTION_INC_ATTACK_POWER
	formula.StatSpec_GetPercent(p.Class, 1, strength, dexterity, vitality, energy, leadership, &percent)
	min := p.leftAttackMin * int(percent) / 100
	max := p.leftAttackMax * int(percent) / 100
	p.leftAttackMin += min
	p.leftAttackMax += max
	p.rightAttackMin += min
	p.rightAttackMax += max
	options = append(options, &model.MsgStatSpec{ID: 1, Min: min, Max: max})
	// STAT_OPTION_INC_MAGIC_DAMAGE
	formula.StatSpec_GetPercent(p.Class, 9, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.magicAttackMin * int(percent) / 100
	max = p.magicAttackMax * int(percent) / 100
	p.magicAttackMin += min
	p.magicAttackMax += max
	options = append(options, &model.MsgStatSpec{ID: 9, Min: min, Max: max})
	// STAT_OPTION_INC_CURSE_DAMAGE
	formula.StatSpec_GetPercent(p.Class, 10, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.curseAttackMin * int(percent) / 100
	max = p.curseAttackMax * int(percent) / 100
	p.curseAttackMin += min
	p.curseAttackMax += max
	options = append(options, &model.MsgStatSpec{ID: 10, Min: min, Max: max})
	// STAT_OPTION_INC_DEFENSE
	formula.StatSpec_GetPercent(p.Class, 4, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.Defense * int(percent) / 100
	p.Defense += min
	options = append(options, &model.MsgStatSpec{ID: 4, Min: min})
	// STAT_OPTION_INC_ATTACK_RATE
	formula.StatSpec_GetPercent(p.Class, 2, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.AttackRate * int(percent) / 100
	p.AttackRate += min
	options = append(options, &model.MsgStatSpec{ID: 2, Min: min})
	// STAT_OPTION_INC_ATTACK_RATE_PVP
	formula.StatSpec_GetPercent(p.Class, 3, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.attackRatePVP * int(percent) / 100
	p.attackRatePVP += min
	options = append(options, &model.MsgStatSpec{ID: 3, Min: min})
	// STAT_OPTION_INC_DEFENSE_RATE
	formula.StatSpec_GetPercent(p.Class, 6, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.DefenseRate * int(percent) / 100
	p.DefenseRate += min
	options = append(options, &model.MsgStatSpec{ID: 6, Min: min})
	// STAT_OPTION_INC_DEFENSE_RATE_PVP
	formula.StatSpec_GetPercent(p.Class, 7, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.defenseRatePVP * int(percent) / 100
	p.defenseRatePVP += min
	options = append(options, &model.MsgStatSpec{ID: 7, Min: min})
	p.Push(&model.MsgStatSpecReply{
		Options: options,
	})

	// hp/mp
	c := CharacterTable[p.Class]
	p.MaxHP = c.HP + int(float32(level-1)*c.LevelHP) + int(float32(vitality-c.Vitality)*c.VitalityHP)
	p.MaxMP = c.MP + int(float32(level-1)*c.LevelMP) + int(float32(energy-c.Energy)*c.EnergyMP)
	// sd
	sdGageConstA := conf.CommonServer.GameServerInfo.SDGageConstA
	sdGageConstB := conf.CommonServer.GameServerInfo.SDGageConstB
	expressionA := strength + dexterity + vitality + energy
	if p.Class == int(class.DarkLord) {
		expressionA += leadership
	}
	expressionB := level * level / sdGageConstB
	p.MaxSD = expressionA*sdGageConstA/10 + expressionB + p.Defense
	// ag
	f := func(s, d, v, e, l float32) int {
		return int(float32(strength)*s + float32(dexterity)*d + float32(vitality)*v + float32(energy)*e + float32(leadership)*l)
	}
	switch class.Class(p.Class) {
	case class.Wizard:
		p.MaxAG = f(0.2, 0.4, 0.3, 0.2, 0)
	case class.Knight:
		p.MaxAG = f(0.15, 0.2, 0.3, 1.0, 0)
	case class.Elf:
		p.MaxAG = f(0.3, 0.2, 0.3, 0.2, 0)
	case class.Magumsa:
		p.MaxAG = f(0.2, 0.25, 0.3, 0.15, 0)
	case class.DarkLord:
		p.MaxAG = f(0.3, 0.2, 0.1, 0.15, 0.3)
	case class.Summoner:
		p.MaxAG = f(0.2, 0.25, 0.3, 0.15, 0)
	case class.RageFighter:
		p.MaxAG = f(0.15, 0.2, 0.3, 1.0, 0)
	}

	// skill bonus
	formula.Knight_Gladiator_CalcSkillBonus(p.Class, energy, &p.knightGladiatorCalcSkillBonus)
	formula.ImpaleSkillCalc(p.Class, energy, &p.impaleSkillCalc)
}

func (p *Player) calcNormalItem() {
	boot := p.inventory.Items[6]
	sameCount := 0
	level10Count := 0
	level11Count := 0
	level12Count := 0
	level13Count := 0
	level14Count := 0
	level15Count := 0
	for i, it := range p.inventory.WearingItems() {
		if it == nil {
			continue
		}
		switch i {
		case 0: // weapon attack
			p.leftAttackMin += it.AttackMin
			p.leftAttackMax += it.AttackMax
			p.curse = it.Magic
		case 1: // weapon attack
			p.rightAttackMin += it.AttackMin
			p.rightAttackMax += it.AttackMax
			p.magic = it.Magic
		case 5: // glove speed
			p.AttackSpeed += it.AttackSpeed
			p.magicSpeed += it.AttackSpeed
		}
		// addition attack
		p.leftAttackMin += it.AdditionAttack
		p.leftAttackMax += it.AdditionAttack
		p.rightAttackMin += it.AdditionAttack
		p.rightAttackMax += it.AdditionAttack
		p.magicAttackMin += it.AdditionMagicAttack
		p.magicAttackMax += it.AdditionMagicAttack
		p.curseAttackMin += it.AdditionCurseAttack
		p.curseAttackMax += it.AdditionCurseAttack
		// addition defense and defense rate
		p.Defense += it.Defense
		p.Defense += it.AdditionDefense
		p.DefenseRate += it.SuccessfulBlocking
		p.DefenseRate += it.AdditionDefenseRate
		if boot != nil && i >= 2 && i <= 6 &&
			it.Code%item.MaxItemIndex == boot.Code%item.MaxItemIndex {
			sameCount++
			if it.Level > 9 {
				level10Count++
			}
			if it.Level > 10 {
				level11Count++
			}
			if it.Level > 11 {
				level12Count++
			}
			if it.Level > 12 {
				level13Count++
			}
			if it.Level > 13 {
				level14Count++
			}
			if it.Level > 14 {
				level15Count++
			}
		}
	}
	// armor bonus defense and defense rate
	// defense level>=10 bonus of item the same type contributed
	// defense success rate bonus of item the same type contributed
	if boot != nil {
		if p.Class == int(class.Magumsa) || p.Class == int(class.RageFighter) {
			sameCount++
			level10Count++
			level11Count++
			level12Count++
			level13Count++
			level14Count++
			level15Count++
		}
		if sameCount == 5 {
			n := 0
			switch {
			case level15Count == 5:
				n = 30
			case level14Count == 5:
				n = 25
			case level13Count == 5:
				n = 20
			case level12Count == 5:
				n = 15
			case level11Count == 5:
				n = 10
			case level10Count == 5:
				n = 5
			}
			p.Defense += p.Defense * n / 100
			p.DefenseRate += p.DefenseRate / 10
		}
	}
	// lucky item bonus
	for _, it := range p.inventory.WearingItems() {
		if it != nil && it.Lucky {
			p.criticalAttackRate += 5
		}
	}
}

func (p *Player) calcExcellentItem() {
	item.ExcellentManager.ApplyExcellentItemOption(&p.inventory, func(option item.ItemExcellentOption, position int) {
		p.excellentAttackRate += option.ExcellentAttackRate
		if option.AttackLevel != 0 {
			value := (p.Level + p.masterLevel) / option.AttackLevel
			switch position {
			case 0:
				p.leftAttackMin += value
				p.leftAttackMax += value
			case 1:
				p.rightAttackMin += value
				p.rightAttackMax += value
			}
		}
		if option.MagicAttackLevel != 0 {
			value := (p.Level + p.masterLevel) / option.MagicAttackLevel
			p.magicAttackMin += value
			p.magicAttackMax += value
		}
		switch position {
		case 0:
			p.leftAttackMin += p.leftAttackMin * option.AttackPercent / 100
			p.leftAttackMax += p.leftAttackMax * option.AttackPercent / 100
		case 1:
			p.rightAttackMin += p.rightAttackMin * option.AttackPercent / 100
			p.rightAttackMax += p.rightAttackMax * option.AttackPercent / 100
		}
		p.magicAttackMin += p.magicAttackMin * option.MagicAttackPercent / 100
		p.magicAttackMax += p.magicAttackMax * option.MagicAttackPercent / 100
		p.AttackSpeed += option.AttackSpeed
		p.magicSpeed += option.MagicSpeed
		p.monsterDieGetHP += option.MonsterDieGetHP
		p.monsterDieGetMP += option.MonsterDieGetMP
		p.MaxHP += p.MaxHP * option.MaxHPPercent / 100
		p.MaxMP += p.MaxMP * option.MaxMPPercent / 100
		p.armorReduceDamage += option.ReduceDamage
		p.armorReflectDamage += option.ReflectDamage
		p.DefenseRate += p.DefenseRate * option.DefenseRatePercent / 100
		p.monsterDieGetMoney += option.MonsterDieGetMoney
		p.MaxHP += option.MaxHP
		p.MaxMP += option.MaxMP
		p.MaxAG += option.MaxAG
		p.ignoreDefenseRate += option.IgnoreDefenseRate
		p.returnDamage += option.ReturnDamage
		p.recoverMaxHP += option.RecoverMaxHP
		p.recoverMaxMP += option.RecoverMaxMP
	})
}

func (p *Player) calcHarmonyItem() {
	item.HarmonyManager.ApplyHarmonyItemOption(&p.inventory, func(option item.ItemHarmonyOption, position int) {
		switch position {
		case 0:
			p.leftAttackMin += option.IncDamageMin + option.IncDamage
			p.leftAttackMax += option.IncDamageMax + option.IncDamage
		case 1:
			p.rightAttackMin += option.IncDamageMin + option.IncDamage
			p.rightAttackMax += option.IncDamageMax + option.IncDamage
		}
		p.magicAttackMin += option.IncMagicDamage + option.IncSkillDamage
		p.magicAttackMax += option.IncMagicDamage + option.IncSkillDamage
		p.criticalAttackDamage += option.IncCriticalDamage
		p.attackRatePVP += option.IncAttackRate
		p.defenseRatePVP += option.IncDefenseRate
		p.Defense += option.IncDefense * 10 / 20
		p.MaxAG += option.IncMaxAG
		p.MaxHP += option.IncMaxHP
		p.recoverHP += option.IncAutoRecoveryHP
		p.recoverMP += option.IncAutoRecoveryMP
		p.armorReduceDamage += option.DecDamage
		p.harmonyIncSDRate += option.IncSDRate
		p.harmonyDecSDRate += option.DecSDRate
		p.harmonyIgnoreSDRate += option.IgnoreSDRate
	})
}

func (p *Player) calc380Item() {
	item.Item380Manager.Apply380ItemOption(&p.inventory, func(option item.Item380Option) {
		p.attackRatePVP += option.IncPVPAttackRate
		p.attackBonusPVP += option.IncPVPAttack
		p.defenseRatePVP += option.IncPVPDefenseRate
		p.defenseBonusPVP += option.IncPVPDefense
		p.MaxHP += option.IncMaxHP
		p.MaxSD += option.IncMaxSD
		p.autoRecoverSD = option.AutoRecoverySD
		p.autoRecoverSDRate += option.IncRecoverySDRate
	})
}
func (p *Player) calcSocketItem() {
	baseMaxHP, baseMaxMP := p.MaxHP, p.MaxMP
	shield := p.inventory.Items[1]
	hasShield := shield != nil && shield.KindB == item.KindBShield && shield.Durability > 0
	item.SocketManager.ApplySocketItemOption(&p.inventory, func(option item.ItemSocketOption, position int) {
		if option.LevelAttack != 0 {
			value := (p.Level + p.masterLevel) / option.LevelAttack
			switch position {
			case 0:
				p.leftAttackMin += value
				p.leftAttackMax += value
			case 1:
				p.rightAttackMin += value
				p.rightAttackMax += value
			}
			p.magicAttackMin += value
			p.magicAttackMax += value
		}
		p.AttackSpeed += option.AttackSpeed
		p.magicSpeed += option.AttackSpeed
		switch position {
		case 0:
			p.leftAttackMax += option.AttackMax
		case 1:
			p.rightAttackMax += option.AttackMax
		}
		p.magicAttackMax += option.AttackMax
		switch position {
		case 0:
			p.leftAttackMin += option.AttackMin
		case 1:
			p.rightAttackMin += option.AttackMin
		}
		p.magicAttackMin += option.AttackMin
		switch position {
		case 0:
			p.leftAttackMin += option.Attack
			p.leftAttackMax += option.Attack
		case 1:
			p.rightAttackMin += option.Attack
			p.rightAttackMax += option.Attack
		}
		p.magicAttackMin += option.Attack
		p.magicAttackMax += option.Attack
		p.agReduce += option.ReduceAGCost
		if p.agReduce > 100 {
			p.agReduce = 100
		}
		p.DefenseRate += p.DefenseRate * option.DefenseRate / 100
		p.Defense += option.Defense
		p.MagicDefense += option.Defense
		if hasShield {
			p.Defense += p.Defense * option.ShieldDefense / 100
		}
		p.armorReduceDamage += option.ReduceDamage
		p.armorReflectDamage += option.ReflectDamage
		if option.KillHP != 0 {
			p.monsterDieGetHP += 1 / float64(option.KillHP)
		}
		if option.KillMP != 0 {
			p.monsterDieGetMP += 1 / float64(option.KillMP)
		}
		p.addSkillAttack += option.SkillAttack
		p.AttackRate += option.AttackRate
		p.recoverHP += option.RecoverHP
		p.recoverMP += option.RecoverMP
		p.MaxHP += baseMaxHP * option.MaxHP / 100
		p.MaxMP += baseMaxMP * option.MaxMP / 100
		p.recoverAG += option.RecoverAG
		p.MaxAG += option.MaxAG
		p.excellentAttackDamage += option.ExcellentDamage
		p.excellentAttackRate += option.ExcellentRate
		p.criticalAttackDamage += option.CriticalDamage
		p.criticalAttackRate += option.CriticalRate
		p.doubleDamageRate += option.DoubleDamage
		p.ignoreDefenseRate += option.IgnoreDefense
		switch position {
		case 0:
			p.leftAttackMin += option.BonusAttack
			p.leftAttackMax += option.BonusAttack
		case 1:
			p.rightAttackMin += option.BonusAttack
			p.rightAttackMax += option.BonusAttack
		}
		p.magicAttackMin += option.BonusMagicAttack
		p.magicAttackMax += option.BonusMagicAttack
		p.addSkillAttack += option.BonusSkillAttack
		p.Defense += option.BonusDefense
		p.MagicDefense += option.BonusDefense
		p.MaxHP += option.BonusMaxHP
	})
}

func (p *Player) calcSetItem(base bool) {
	item.SetManager.ApplySetItemOption(&p.inventory, func(option item.ItemSetOption) {
		if base {
			p.addStrength += option.Strength
			p.addDexterity += option.Dexterity
			p.addEnergy += option.Energy
			p.addVitality += option.Vitality
			p.addLeadership += option.Leadership
		} else {
			p.leftAttackMin += option.AttackMin
			p.rightAttackMin += option.AttackMin
			p.leftAttackMax += option.AttackMax
			p.rightAttackMax += option.AttackMax
			p.magicAttackMin += option.MagicAttack
			p.magicAttackMax += option.MagicAttack
			p.setAddDamage += option.Damage
			p.AttackRate += option.AttackRate
			p.Defense += option.Defense
			p.MaxHP += option.MaxHP
			p.MaxMP += option.MaxMP
			p.MaxAG += option.MaxAG
			p.criticalAttackRate += option.CriticalAttackRate
			p.criticalAttackDamage += option.CriticalAttackDamage
			p.excellentAttackRate += option.ExcellentAttackRate
			p.excellentAttackDamage += option.ExcellentAttackDamage
			p.addSkillAttack += option.SkillAttack
			p.doubleDamageRate += option.DoubleDamageRate
			p.ignoreDefenseRate += option.IgnoreDefenseRate
			p.MaxSD += option.MaxSD
			p.setIncreaseTwoHandWeaponDamage += option.TwoHandWeaponDamage
		}
	})
}

func (p *Player) convertAttack() {
	leftHand := p.inventory.Items[0]
	rightHand := p.inventory.Items[1]
	// convert left/right attack to attack
	left, right := false, false
	if leftHand != nil &&
		leftHand.KindA == item.KindAWeapon &&
		leftHand.Code != item.Code(4, 7) &&
		leftHand.Code != item.Code(4, 15) {
		left = true
	}
	if rightHand != nil &&
		rightHand.KindA == item.KindAWeapon &&
		rightHand.Code != item.Code(4, 7) &&
		rightHand.Code != item.Code(4, 15) {
		right = true
	}
	switch {
	case left && right:
		switch class.Class(p.Class) {
		case class.Knight, class.Magumsa:
			if leftHand.Code == rightHand.Code {
				formula.CalcTwoSameWeaponBonus(
					p.leftAttackMin, p.leftAttackMax, p.rightAttackMin, p.rightAttackMax,
					&p.leftAttackMin, &p.rightAttackMin, &p.leftAttackMax, &p.rightAttackMax)
			} else {
				formula.CalcTwoDifferentWeaponBonus(
					p.leftAttackMin, p.leftAttackMax, p.rightAttackMin, p.rightAttackMax,
					&p.leftAttackMin, &p.rightAttackMin, &p.leftAttackMax, &p.rightAttackMax)
			}
		case class.RageFighter:
			formula.CalcRageFighterTwoWeaponBonus(
				p.leftAttackMin, p.leftAttackMax, p.rightAttackMin, p.rightAttackMax,
				&p.leftAttackMin, &p.rightAttackMin, &p.leftAttackMax, &p.rightAttackMax)
		}
		p.AttackMin = p.leftAttackMin + p.rightAttackMin
		p.AttackMax = p.leftAttackMax + p.rightAttackMax
		p.AttackSpeed += (leftHand.AttackSpeed + rightHand.AttackSpeed) / 2
		p.magicSpeed += (leftHand.AttackSpeed + rightHand.AttackSpeed) / 2
	case left:
		p.AttackMin = p.leftAttackMin
		p.AttackMax = p.leftAttackMax
		p.AttackSpeed += leftHand.AttackSpeed
		p.magicSpeed += leftHand.AttackSpeed
	case right:
		p.AttackMin = p.rightAttackMin
		p.AttackMax = p.rightAttackMax
		p.AttackSpeed += rightHand.AttackSpeed
		p.magicSpeed += rightHand.AttackSpeed
	default:
		p.AttackMin = (p.leftAttackMin + p.rightAttackMin) / 2
		p.AttackMax = (p.leftAttackMax + p.rightAttackMax) / 2
	}
	// convert magic/curse to magic/curse attack
	p.magicAttackMin += p.magicAttackMax * p.magic / 100
	p.magicAttackMax += p.magicAttackMax * p.magic / 100
	p.curseAttackMin += p.curseAttackMax * p.curse / 100
	p.curseAttackMax += p.curseAttackMax * p.curse / 100
}

func (p *Player) calcHelperItem() {
	// wing
	wing := p.inventory.Items[7]
	if wing != nil {
		wingIncreaseDamage := 0
		formula.Wings_CalcIncAttack(wing.Code, wing.Level, &wingIncreaseDamage)
		p.wingIncreaseDamage = wingIncreaseDamage - 100
		wingReduceDamage := 0
		formula.Wings_CalcAbsorb(wing.Code, wing.Level, &wingReduceDamage)
		p.wingReduceDamage = 100 - wingReduceDamage
	}
	// helper item contribution
	helper := p.inventory.Items[8]
	switch {
	case p.equippedItem(helper, item.Code(13, 1)): // 小恶魔
		percent := conf.PetRing.Pets.Imp.AddAttackPercent
		p.AttackMin += p.AttackMin * percent / 100
		p.AttackMax += p.AttackMax * percent / 100
		p.magicAttackMin += p.magicAttackMin * percent / 100
		p.magicAttackMax += p.magicAttackMax * percent / 100
		p.curseAttackMin += p.curseAttackMin * percent / 100
		p.curseAttackMax += p.curseAttackMax * percent / 100
	case p.equippedItem(helper, item.Code(13, 64)): // 强化恶魔
		percent := conf.PetRing.Pets.Demon.AddAttackPercent
		p.AttackMin += p.AttackMin * percent / 100
		p.AttackMax += p.AttackMax * percent / 100
		p.magicAttackMin += p.magicAttackMin * percent / 100
		p.magicAttackMax += p.magicAttackMax * percent / 100
		p.curseAttackMin += p.curseAttackMin * percent / 100
		p.curseAttackMax += p.curseAttackMax * percent / 100
		p.AttackSpeed += conf.PetRing.Pets.Demon.AddAttackSpeed
		p.magicSpeed += conf.PetRing.Pets.Demon.AddAttackSpeed
	case p.equippedItem(helper, item.Code(13, 0)): // 守护天使
		p.MaxHP += conf.PetRing.Pets.Angel.AddHP
		p.helperReduceDamage = conf.PetRing.Pets.Angel.ReduceDamagePercent
	case p.equippedItem(helper, item.Code(13, 65)): // 强化天使
		p.MaxHP += conf.PetRing.Pets.SpiritAngel.AddHP
		p.helperReduceDamage = conf.PetRing.Pets.SpiritAngel.ReduceDamagePercent
	case p.equippedItem(helper, item.Code(13, 80)): // 熊猫
		p.Defense += conf.PetRing.Pets.Panda.AddDefenseValue
	}
	// pet item contribution
	if p.pet != nil && p.pet.Durability > 0 {
		switch p.pet.Code {
		case item.Code(13, 3): // Horn of Dinorant 彩云兽
			p.petIncreaseDamage = 15
			p.petReduceDamage = 10
		case item.Code(13, 4): // Dark Horse 黑王马之角
			p.petReduceDamage = (30 + p.pet.Level) / 2
		case item.Code(13, 37): // Horn of Fenrir 炎狼兽之角
			p.petIncreaseDamage = 33
			p.petReduceDamage = 10
		}
	}
	// ...
	if p.SD > p.MaxSD {
		p.SD = p.MaxSD
	}
	if p.AG > p.MaxAG {
		p.AG = p.MaxAG
	}
}

func (p *Player) equippedItem(it *item.Item, code int) bool {
	if it == nil || it.Durability == 0 {
		return false
	}
	return it.Code == code
}

func (p *Player) calcBuff() {
	delayLevel := 0
	p.Buffs().ApplyBuffEffect(func(effect buff.BuffEffect) {
		p.AttackMin += effect.Attack + effect.AttackMin
		p.AttackMax += effect.Attack + effect.AttackMax
		p.magicAttackMin += effect.MagicMin
		p.magicAttackMax += effect.MagicMax
		p.curseAttackMin += effect.CurseMin
		p.curseAttackMax += effect.CurseMax
		p.Defense += effect.Defense
		p.DefenseRate += effect.DefenseRate
		p.MaxHP += int(float64(p.MaxHP) * effect.IncMaxHP / 100)
		p.MaxMP += int(float64(p.MaxMP) * effect.IncMaxMP / 100)
		p.criticalAttackDamage += effect.CriticalDamage
		p.ignoreDefenseRate += effect.IgnoreDefense
		p.BuffDefenseReduction += effect.DefenseReduction
		p.BuffAttackReduction += effect.AttackReduction
		p.BuffReflect += effect.Reflect
	})
	for _, active := range p.Buffs() {
		if active.Effect.Slow {
			delayLevel++
		}
	}
	p.SetEffectDelayLevel(delayLevel)
}
