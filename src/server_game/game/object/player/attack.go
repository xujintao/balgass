package player

func (p *Player) GetAttackRatePVP() int {
	return p.attackRatePVP
}

func (p *Player) GetDefenseRatePVP() int {
	return p.defenseRatePVP
}

func (p *Player) GetIgnoreDefenseRate() int {
	return p.IgnoreDefenseRate
}

func (p *Player) GetCriticalAttackRate() int {
	return p.CriticalAttackRate
}

func (p *Player) GetCriticalAttackDamage() int {
	return p.CriticalAttackDamage
}

func (p *Player) GetExcellentAttackRate() int {
	return p.ExcellentAttackRate
}

func (p *Player) GetExcellentAttackDamage() int {
	return p.ExcellentAttackDamage
}

func (p *Player) GetAddDamage() int {
	return p.SetAddDamage
}

func (p *Player) GetWingIncreaseDamage() int {
	return p.WingIncreaseDamage
}

func (p *Player) GetWingReduceDamage() int {
	return p.WingReduceDamage
}

func (p *Player) GetDoubleDamageRate() int {
	return p.DoubleDamageRate
}

func (p *Player) GetArmorReduceDamage() int {
	return p.ArmorReduceDamage
}

func (p *Player) GetHelperReduceDamage() int {
	return p.HelperReduceDamage
}

func (p *Player) GetKnightGladiatorCalcSkillBonus() float64 {
	return p.KnightGladiatorCalcSkillBonus
}

func (p *Player) GetImpaleSkillCalc() float64 {
	return p.ImpaleSkillCalc
}
