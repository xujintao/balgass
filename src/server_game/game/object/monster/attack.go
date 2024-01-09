package monster

import "github.com/xujintao/balgass/src/server_game/game/object"

func (*Monster) GetAttackRatePVP() int {
	return 0
}

func (*Monster) GetDefenseRatePVP() int {
	return 0
}

func (m *Monster) GetIgnoreDefenseRate() int {
	return 0
}

func (m *Monster) GetCriticalAttackRate() int {
	return 0
}

func (m *Monster) GetCriticalAttackDamage() int {
	return 0
}

func (m *Monster) GetExcellentAttackRate() int {
	return 0
}

func (m *Monster) GetExcellentAttackDamage() int {
	return 0
}

func (m *Monster) GetAddDamage() int {
	return 0
}

func (*Monster) GetArmorReduceDamage() int {
	return 0
}

func (m *Monster) GetWingIncreaseDamage() int {
	return 0
}

func (m *Monster) GetWingReduceDamage() int {
	return 0
}

func (*Monster) GetHelperReduceDamage() int {
	return 0
}

func (*Monster) GetPetIncreaseDamage() int {
	return 0
}

func (*Monster) GetPetReduceDamage() int {
	return 0
}

func (m *Monster) GetDoubleDamageRate() int {
	return 0
}

func (*Monster) GetKnightGladiatorCalcSkillBonus() float64 {
	return 1.0
}

func (*Monster) GetImpaleSkillCalc() float64 {
	return 1.0
}

func (m *Monster) Die(obj *object.Object) {
	// give experience
	obj.MonsterDieGetExperience(&m.Object)
	// give item
	// obj delay recover hp/mp/sd
	obj.AddDelayMsg(2, 0, 2000, m.Index)
}

func (*Monster) MonsterDieGetExperience(*object.Object) {}

func (*Monster) MonsterDieRecoverHP() {}
