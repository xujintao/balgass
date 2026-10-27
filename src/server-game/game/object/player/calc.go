package player

import (
	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/formula"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/model"
)

func (p *Player) calc() {
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
	p.DefenseRate = 0
	p.defenseRatePVP = 0
	p.AttackSpeed = 0
	p.magicSpeed = 0
	p.MaxHP = 0
	p.MaxMP = 0
	p.MaxSD = 0
	p.MaxAG = 0
	p.knightGladiatorCalcSkillBonus = 0
	p.impaleSkillCalc = 0

	p.setAddDamage = 0
	p.setIncreaseTwoHandWeaponDamage = 0
	p.criticalAttackRate = 0
	p.criticalAttackDamage = 0
	p.excellentAttackRate = 0
	p.excellentAttackDamage = 0
	p.addSkillAttack = 0
	p.doubleDamageRate = 0
	p.ignoreDefenseRate = 0
	p.monsterDieGetHP = 0
	p.monsterDieGetMP = 0
	p.monsterDieGetMoney = 0
	p.recoverHP = 0
	p.recoverMP = 0
	p.recoverAG = 0
	p.agReduce = 0
	p.shieldDefense = 0
	p.armorReduceDamage = 0
	p.armorReflectDamage = 0
	p.returnDamage = 0
	p.recoverMaxHP = 0
	p.recoverMaxMP = 0

	p.wingIncreaseDamage = 0
	p.wingReduceDamage = 0
	p.helperReduceDamage = 0
	p.petIncreaseDamage = 0
	p.petReduceDamage = 0

	// calc all stats
	p.calcBaseStats()
	p.calcNormalItem()
	p.calcSetItem(false)
	p.calcSocketItem()
	p.calcExcellentItem()
	p.convertAttack()
	p.calcHelperItem()

	// Push
	p.Push(&model.MsgAttackSpeedReply{
		AttackSpeed: p.AttackSpeed,
		MagicSpeed:  p.magicSpeed,
	})
	p.PushMaxHPSD(p.MaxHP+p.AddHP, p.MaxSD)
	p.PushMaxMPAG(p.MaxMP+p.AddMP, p.MaxAG)
	p.PushHPSD(p.HP, p.SD)
	p.PushMPAG(p.MP, p.AG)
}

func (p *Player) calcBaseStats() {
	leftHand := p.inventory.Items[0]
	rightHand := p.inventory.Items[1]
	wing := p.inventory.Items[7]

	// Earth seeds modify attributes used by the base character formulas.
	for i := 0; i < 12; i++ {
		it := p.inventory.Items[i]
		if it == nil || it.ItemBase == nil || it.Type != item.TypeSocket || it.Durability <= 0 {
			continue
		}
		for _, slot := range []int{it.SocketSlot1, it.SocketSlot2, it.SocketSlot3, it.SocketSlot4, it.SocketSlot5} {
			effect, ok := item.SocketManager.SeedEffect(slot)
			if !ok {
				continue
			}
			switch effect.Index {
			case item.SocketEffectVitality:
				p.addVitality += socketEffectValue(p.vitality+p.addVitality, effect.ValueType, effect.Value)
			}
		}
	}
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
	leftHand := p.inventory.Items[0]
	rightHand := p.inventory.Items[1]
	glove := p.inventory.Items[5]
	boot := p.inventory.Items[6]
	wing := p.inventory.Items[7]

	// weapon and weapon addition attack
	if leftHand != nil {
		p.leftAttackMin += leftHand.AttackMin + leftHand.AdditionAttack
		p.leftAttackMax += leftHand.AttackMax + leftHand.AdditionAttack
		p.magicAttackMin += leftHand.AdditionMagicAttack
		p.magicAttackMax += leftHand.AdditionMagicAttack
		p.curseAttackMin += leftHand.AdditionCurseAttack
		p.curseAttackMax += leftHand.AdditionCurseAttack
		p.curse = leftHand.Magic
	}
	if rightHand != nil {
		p.rightAttackMin += rightHand.AttackMin + rightHand.AdditionAttack
		p.rightAttackMax += rightHand.AttackMax + rightHand.AdditionAttack
		p.magicAttackMin += rightHand.AdditionMagicAttack
		p.magicAttackMax += rightHand.AdditionMagicAttack
		p.curseAttackMin += rightHand.AdditionCurseAttack
		p.curseAttackMax += rightHand.AdditionCurseAttack
		p.magic = rightHand.Magic
	}

	// wing addition attack
	if wing != nil {
		p.leftAttackMin += wing.AdditionAttack
		p.leftAttackMax += wing.AdditionAttack
		p.rightAttackMin += wing.AdditionAttack
		p.rightAttackMax += wing.AdditionAttack
		p.magicAttackMin += wing.AdditionMagicAttack
		p.magicAttackMax += wing.AdditionMagicAttack
		p.curseAttackMin += wing.AdditionCurseAttack
		p.curseAttackMax += wing.AdditionCurseAttack
	}

	// glove speed
	if glove != nil {
		p.AttackSpeed += glove.AttackSpeed
		p.magicSpeed += glove.AttackSpeed
	}

	// armor(shield|armor|wing) addition defense(rate)
	for i := 1; i <= 7; i++ {
		it := p.inventory.Items[i]
		if it != nil {
			p.Defense += it.Defense
			p.Defense += it.AdditionDefense
			p.DefenseRate += it.SuccessfulBlocking
			p.DefenseRate += it.AdditionDefenseRate
		}
	}

	// armor bonus defense and defense rate
	// defense level>=10 bonus of item the same type contributed
	// defense success rate bonus of item the same type contributed
	sameCount := 0
	level10Count := 0
	level11Count := 0
	level12Count := 0
	level13Count := 0
	level14Count := 0
	level15Count := 0
	if boot != nil {
		refer := boot.Code % item.MaxItemIndex
		for i := 2; i <= 6; i++ {
			it := p.inventory.Items[i]
			if it != nil && it.Code%item.MaxItemIndex == refer {
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
}

func (p *Player) calcSetItem(base bool) {
	type set struct {
		index int
		count int
	}
	var sets []set
	sameWeapon := 0
	sameRing := 0
	for i, it := range p.inventory.WearingItems() {
		if it == nil {
			continue
		}
		if it.Durability == 0 {
			continue
		}
		index := it.Set
		if index <= 0 {
			continue
		}
		if i == 0 {
			sameWeapon = index
		}
		if i == 1 && sameWeapon > 0 {
			continue
		}
		if i == 10 {
			sameRing = index
		}
		if i == 11 && sameRing > 0 {
			continue
		}
		ok := false
		for i := range sets {
			if sets[i].index == index {
				sets[i].count++
				ok = true
				break
			}
		}
		if !ok {
			sets = append(sets, set{index, 1})
		}
	}

	for _, v := range sets {
		index := v.index
		count := v.count
		if count >= 2 {
			for i := 0; i < count-1; i++ {
				setEffect := item.SetManager.GetSetEffect(index, i)
				p.applySetEffect(setEffect.Index, setEffect.Value, base)
			}

			if count > item.SetManager.GetSetEffectCount(index) {
				setEffects := item.SetManager.GetSetFull(index)
				for _, setEffect := range setEffects {
					p.applySetEffect(setEffect.Index, setEffect.Value, base)
				}
			}
		}
	}
}

func (p *Player) applySetEffect(index item.SetEffectType, value int, base bool) {
	if base {
		switch index {
		case item.SetEffectIncStrength:
			p.addStrength += value
		case item.SetEffectIncAgility:
			p.addDexterity += value
		case item.SetEffectIncEnergy:
			p.addEnergy += value
		case item.SetEffectIncVitality:
			p.addVitality += value
		case item.SetEffectIncLeadership:
			p.addLeadership += value
		}
	} else {
		switch index {
		case item.SetEffectIncAttackMin:
			p.leftAttackMin += value
			p.rightAttackMin += value
		case item.SetEffectIncAttackMax:
			p.leftAttackMax += value
			p.rightAttackMax += value
		case item.SetEffectIncMagicAttack:
			p.magicAttackMin += value
			p.magicAttackMax += value
		case item.SetEffectIncDamage:
			p.setAddDamage += value
		case item.SetEffectIncAttackRate:
			p.AttackRate += value
		case item.SetEffectIncDefense:
			p.Defense += value
		case item.SetEffectIncMaxHP:
			p.MaxHP += value
		case item.SetEffectIncMaxMP:
			p.MaxMP += value
		case item.SetEffectIncMaxAG:
			p.MaxAG += value
		case item.SetEffectIncAG:
			p.MaxAG += value
		case item.SetEffectIncCritiDamageRate:
			p.criticalAttackRate += value
		case item.SetEffectIncCritiDamage:
			p.criticalAttackDamage += value
		case item.SetEffectIncExcelDamageRate:
			p.excellentAttackRate += value
		case item.SetEffectIncExcelDamage:
			p.excellentAttackDamage += value
		case item.SetEffectIncSkillAttack:
			p.addSkillAttack += value
		case item.SetEffectDoubleDamage:
			p.doubleDamageRate += value
		case item.SetEffectIgnoreDefense:
			p.ignoreDefenseRate += value
		case item.SetEffectIncShieldDefense:
			p.MaxSD += value
		case item.SetEffectIncTwoHandSwordDamage:
			p.setIncreaseTwoHandWeaponDamage += value
		}
	}
}

func (p *Player) calcSocketItem() {
	var counts [6]int
	var bonuses []item.SocketBonusEffect
	baseHP, baseMP, baseAG := p.MaxHP, p.MaxMP, p.MaxAG
	for position := 0; position < 12; position++ {
		it := p.inventory.Items[position]
		if it == nil || it.ItemBase == nil || it.Type != item.TypeSocket || it.Durability <= 0 {
			continue
		}
		installed := false
		for _, slot := range []int{it.SocketSlot1, it.SocketSlot2, it.SocketSlot3, it.SocketSlot4, it.SocketSlot5} {
			effect, ok := item.SocketManager.SeedEffect(slot)
			if !ok {
				continue
			}
			installed = true
			counts[effect.Element-1]++
			p.applySocketEffect(effect, baseHP, baseMP, baseAG)
		}
		if installed {
			if effect, ok := item.SocketManager.BonusEffect(it.SocketBonus); ok &&
				it.Section >= effect.MinSection && it.Section <= effect.MaxSection {
				bonuses = append(bonuses, effect)
			}
		}
	}
	for _, effect := range bonuses {
		p.applySocketBonusEffect(effect, baseHP)
	}
	for _, effect := range item.SocketManager.SpecialEffects(counts) {
		p.applySocketEffect(effect, baseHP, baseMP, baseAG)
	}
	shield := p.inventory.Items[1]
	if p.shieldDefense != 0 &&
		shield != nil &&
		shield.KindB == item.KindBShield &&
		shield.Durability > 0 {
		p.Defense += p.Defense * p.shieldDefense / 100
	}
}

func (p *Player) applySocketEffect(effect item.SocketEffect, baseHP, baseMP, baseAG int) {
	value := func(base int) int {
		return socketEffectValue(base, effect.ValueType, effect.Value)
	}
	switch effect.Index {
	case item.SocketEffectLevelAttack:
		v := value(p.Level + p.masterLevel)
		p.leftAttackMin += v
		p.leftAttackMax += v
		p.rightAttackMin += v
		p.rightAttackMax += v
		p.magicAttackMin += v
		p.magicAttackMax += v
	case item.SocketEffectAttackSpeed:
		p.AttackSpeed += value(p.AttackSpeed)
		p.magicSpeed += value(p.magicSpeed)
	case item.SocketEffectAttackMax:
		p.leftAttackMax += value(p.leftAttackMax)
		p.rightAttackMax += value(p.rightAttackMax)
		p.magicAttackMax += value(p.magicAttackMax)
	case item.SocketEffectAttackMin:
		p.leftAttackMin += value(p.leftAttackMin)
		p.rightAttackMin += value(p.rightAttackMin)
		p.magicAttackMin += value(p.magicAttackMin)
	case item.SocketEffectAttack:
		p.leftAttackMin += value(p.leftAttackMin)
		p.leftAttackMax += value(p.leftAttackMax)
		p.rightAttackMin += value(p.rightAttackMin)
		p.rightAttackMax += value(p.rightAttackMax)
		p.magicAttackMin += value(p.magicAttackMin)
		p.magicAttackMax += value(p.magicAttackMax)
	case item.SocketEffectReduceAGCost:
		p.agReduce += effect.Value
		if p.agReduce > 100 {
			p.agReduce = 100
		}
	case item.SocketEffectDefenseRate:
		p.DefenseRate += value(p.DefenseRate)
	case item.SocketEffectDefense:
		p.Defense += value(p.Defense)
		p.MagicDefense += value(p.MagicDefense)
	case item.SocketEffectShieldDefense:
		p.shieldDefense += effect.Value
	case item.SocketEffectReduceDamage:
		p.armorReduceDamage += effect.Value
	case item.SocketEffectReflectDamage:
		p.armorReflectDamage += effect.Value
	case item.SocketEffectKillHP:
		if effect.Value > 0 {
			p.monsterDieGetHP += 1 / float64(effect.Value)
		}
	case item.SocketEffectKillMP:
		if effect.Value > 0 {
			p.monsterDieGetMP += 1 / float64(effect.Value)
		}
	case item.SocketEffectSkillAttack:
		p.addSkillAttack += value(p.addSkillAttack)
	case item.SocketEffectAttackRate:
		p.AttackRate += effect.Value
	case item.SocketEffectRecoverHP:
		p.recoverHP += value(baseHP)
	case item.SocketEffectMaxHP:
		p.MaxHP += value(baseHP)
	case item.SocketEffectMaxMP:
		p.MaxMP += value(baseMP)
	case item.SocketEffectRecoverMP:
		p.recoverMP += value(baseMP)
	case item.SocketEffectMaxAG:
		p.MaxAG += value(baseAG)
	case item.SocketEffectRecoverAG:
		p.recoverAG += value(p.recoverAG)
	case item.SocketEffectExcellentDamage:
		p.excellentAttackDamage += value(p.excellentAttackDamage)
	case item.SocketEffectExcellentRate:
		p.excellentAttackRate += effect.Value
	case item.SocketEffectCriticalDamage:
		p.criticalAttackDamage += value(p.criticalAttackDamage)
	case item.SocketEffectCriticalRate:
		p.criticalAttackRate += effect.Value
	case item.SocketEffectDoubleDamage:
		p.doubleDamageRate += effect.Value
	case item.SocketEffectIgnoreDefense:
		p.ignoreDefenseRate += effect.Value
	}
}

func (p *Player) applySocketBonusEffect(bonus item.SocketBonusEffect, baseHP int) {
	value := func(base int) int { return socketEffectValue(base, bonus.ValueType, bonus.Value) }
	switch bonus.Index {
	case item.SocketBonusWeaponAttack, item.SocketBonusStaffAttack,
		item.SocketBonusWeaponAttack2, item.SocketBonusStaffAttack2:
		v := value(p.Level + p.masterLevel)
		p.leftAttackMin += v
		p.leftAttackMax += v
		p.rightAttackMin += v
		p.rightAttackMax += v
		p.magicAttackMin += v
		p.magicAttackMax += v
	case item.SocketBonusWeaponSkillAttack, item.SocketBonusStaffSkillAttack,
		item.SocketBonusWeaponSkillAttack2, item.SocketBonusStaffSkillAttack2:
		p.addSkillAttack += value(p.addSkillAttack)
	case item.SocketBonusArmorDefense, item.SocketBonusArmorDefense2:
		p.Defense += value(p.Defense)
		p.MagicDefense += value(p.MagicDefense)
	case item.SocketBonusArmorMaxHP, item.SocketBonusArmorMaxHP2:
		p.MaxHP += value(baseHP)
	}
}

func socketEffectValue(base, valueType, effectValue int) int {
	if effectValue <= 0 {
		return 0
	}
	switch valueType {
	case 1:
		return effectValue
	case 2:
		return base * effectValue / 100
	case 3, 4, 5:
		return base / effectValue
	}
	return 0
}

func (p *Player) calcExcellentItem() {
	level := p.Level + p.masterLevel
	for i := 0; i < 12; i++ {
		it := p.inventory.Items[i]
		if it != nil {
			if it.Lucky {
				p.criticalAttackRate += 5
			}
			if it.ExcellentAttackRate {
				p.excellentAttackRate += 10
			}
			if it.ExcellentAttackLevel {
				value := level / 20
				switch i {
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
			if it.ExcellentAttackPercent {
				switch i {
				case 0:
					p.leftAttackMin += p.leftAttackMin * 2 / 100
					p.leftAttackMax += p.leftAttackMax * 2 / 100
				case 1:
					p.rightAttackMin += p.rightAttackMin * 2 / 100
					p.rightAttackMax += p.rightAttackMax * 2 / 100
				}
				p.magicAttackMin += p.magicAttackMin * 2 / 100
				p.magicAttackMax += p.magicAttackMax * 2 / 100
			}
			if it.ExcellentAttackSpeed {
				p.AttackSpeed += 7
				p.magicSpeed += 7
			}
			if it.ExcellentAttackHP {
				p.monsterDieGetHP += 1.0 / 8
			}
			if it.ExcellentAttackMP {
				p.monsterDieGetMP += 1.0 / 8
			}
			if it.ExcellentDefenseHP {
				p.MaxHP += p.MaxHP * 4 / 100
			}
			if it.ExcellentDefenseMP {
				p.MaxMP += p.MaxMP * 4 / 100
			}
			if it.ExcellentDefenseReduce {
				p.armorReduceDamage += 4
			}
			if it.ExcellentDefenseReflect {
				p.armorReflectDamage += 5
			}
			if it.ExcellentDefenseRate {
				p.DefenseRate += p.DefenseRate * 10 / 100
			}
			if it.ExcellentDefenseMoney {
				p.monsterDieGetMoney += float64(30) / 100
			}
			if it.Wing2ExcellentSpeed {
				p.AttackSpeed += 5
				p.magicSpeed += 5
			}
			if it.Wing2ExcellentAG {
				p.MaxAG += 50
			}
			if it.Wing2ExcellentIgnore {
				p.ignoreDefenseRate += 3
			}
			if it.Wing2ExcellentMP {
				p.MaxMP += 50 + 5*it.Level
			}
			if it.Wing2ExcellentHP {
				p.MaxHP += 50 + 5*it.Level
			}
			if it.Wing3ExcellentIgnore {
				p.ignoreDefenseRate += 5
			}
			if it.Wing3ExcellentReturn {
				p.returnDamage += 5
			}
			if it.Wing3ExcellentHP {
				p.recoverMaxHP += 5
			}
			if it.Wing3ExcellentMP {
				p.recoverMaxMP += 5
			}
		}
	}
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

	if p.HP > p.MaxHP+p.AddHP {
		p.HP = p.MaxHP + p.AddHP
	}
	if p.MP > p.MaxMP+p.AddMP {
		p.MP = p.MaxMP + p.AddMP
	}
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

func (player *Player) calc380Item() {
	for _, wItem := range player.inventory.WearingItems() {
		if wItem.Durability == 0 {
			continue
		}
		if !wItem.Option380 || !item.Item380Manager.Is380Item(wItem.Section, wItem.Index) {
			continue
		}
		item.Item380Manager.Apply380ItemEffect(wItem.Section, wItem.Index, &player.item380Effect)
	}
	player.MaxHP += player.item380Effect.Item380EffectIncMaxHP
	player.MaxSD += player.item380Effect.Item380EffectIncMaxSD
}
