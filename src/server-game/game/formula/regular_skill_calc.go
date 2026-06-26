package formula

func Knight_Gladiator_CalcSkillBonus(class, energy int, rate *float64) {
	call(f.RegularSkillCacl, "Knight_Gladiator_CalcSkillBonus", "iii>d", class, 1, energy, rate)
}

func ImpaleSkillCalc(class, energy int, rate *float64) {
	call(f.RegularSkillCacl, "ImpaleSkillCalc", "iii>d", class, 1, energy, rate)
}

func ElfHeal(class, index, targetIndex, energy int, addLife *int) {
	call(f.RegularSkillCacl, "ElfHeal", "iiii>i", class, index, targetIndex, energy, addLife)
}

func ElfAttack(class, index, targetIndex, energy int, attack, duration *float64) {
	call(f.RegularSkillCacl, "ElfAttack", "iiii>dd", class, index, targetIndex, energy, attack, duration)
}

func ElfDefense(class, index, targetIndex, energy int, defense, duration *float64) {
	call(f.RegularSkillCacl, "ElfDefense", "iiii>dd", class, index, targetIndex, energy, defense, duration)
}

func KnightSkillAddLife(vitality, energy, partyBonus int, addLifeRate *float64, duration *int) {
	call(f.RegularSkillCacl, "KnightSkillAddLife", "iii>di", vitality, energy, partyBonus, addLifeRate, duration)
}
