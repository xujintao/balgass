package formula

func Knight_Gladiator_CalcSkillBonus(class, energy int, rate *float64) {
	call(f.RegularSkillCacl, "Knight_Gladiator_CalcSkillBonus", "iii>d", class, 1, energy, rate)
}

func ImpaleSkillCalc(class, energy int, rate *float64) {
	call(f.RegularSkillCacl, "ImpaleSkillCalc", "iii>d", class, 1, energy, rate)
}
