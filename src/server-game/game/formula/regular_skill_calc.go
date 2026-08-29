package formula

func Knight_Gladiator_CalcSkillBonus(class, energy int, rate *float64) {
	call(f.RegularSkillCalc, "Knight_Gladiator_CalcSkillBonus", "iii>d", class, 1, energy, rate)
}

func ImpaleSkillCalc(class, energy int, rate *float64) {
	call(f.RegularSkillCalc, "ImpaleSkillCalc", "iii>d", class, 1, energy, rate)
}

func Elf_CalcSkillBonus(damage, energy int, out *int) {
	call(f.RegularSkillCalc, "Elf_CalcSkillBonus", "ii>i", damage, energy, out)
}

func GladiatorPowerSlash(damage, energy int, out *int) {
	call(f.RegularSkillCalc, "GladiatorPowerSlash", "ii>i", damage, energy, out)
}

func Lord_CalcSkillBonus(damage, energy int, out *int) {
	call(f.RegularSkillCalc, "Lord_CalcSkillBonus", "ii>i", damage, energy, out)
}

func ChainLightningCalc(damage, targetNumber int, out *int) {
	call(f.RegularSkillCalc, "ChainLightningCalc", "ii>i", damage, targetNumber, out)
}

func StrikeOfDestructionCalc(damage, energy int, out *int) {
	call(f.RegularSkillCalc, "StrikeOfDestructionCalc", "ii>i", damage, energy, out)
}

func FlameStrikeCalc(damage int, out *int) {
	call(f.RegularSkillCalc, "FlameStrikeCalc", "i>i", damage, out)
}

func ChaoticDiseierCalc(damage, energy int, out *int) {
	call(f.RegularSkillCalc, "ChaoticDiseierCalc", "ii>i", damage, energy, out)
}

func RageFighterKillingBlow(damage, vitality int, out *int) {
	call(f.RegularSkillCalc, "RageFighterKillingBlow", "ii>i", damage, vitality, out)
}

func RageFighterBeastUppercut(damage, vitality int, out *int) {
	call(f.RegularSkillCalc, "RageFighterBeastUppercut", "ii>i", damage, vitality, out)
}

func RageFighterChainDrive(damage, vitality int, out *int) {
	call(f.RegularSkillCalc, "RageFighterChainDrive", "ii>i", damage, vitality, out)
}

func RageFighterDarkSideIncDamage(damage, dexterity, energy int, out *int) {
	call(f.RegularSkillCalc, "RageFighterDarkSideIncDamage", "iii>i", damage, dexterity, energy, out)
}

func RageFighterDragonRoar(damage, energy int, out *int) {
	call(f.RegularSkillCalc, "RageFighterDragonRoar", "ii>i", damage, energy, out)
}

func RageFighterDragonSlasher(damage, energy, targetType int, out *int) {
	call(f.RegularSkillCalc, "RageFighterDragonSlasher", "iii>i", damage, energy, targetType, out)
}

func RageFighterCharge(damage, vitality int, out *int) {
	call(f.RegularSkillCalc, "RageFighterCharge", "ii>i", damage, vitality, out)
}

func RageFighterPhoenixShot(damage, vitality int, out *int) {
	call(f.RegularSkillCalc, "RageFighterPhoenixShot", "ii>i", damage, vitality, out)
}

func ElfHeal(class, index, targetIndex, energy int, addLife *int) {
	call(f.RegularSkillCalc, "ElfHeal", "iiii>i", class, index, targetIndex, energy, addLife)
}

func ElfAttack(class, index, targetIndex, energy int, attack, duration *float64) {
	call(f.RegularSkillCalc, "ElfAttack", "iiii>dd", class, index, targetIndex, energy, attack, duration)
}

func ElfDefense(class, index, targetIndex, energy int, defense, duration *float64) {
	call(f.RegularSkillCalc, "ElfDefense", "iiii>dd", class, index, targetIndex, energy, defense, duration)
}

func KnightSkillAddLife(vitality, energy, partyBonus int, addLifeRate *float64, duration *int) {
	call(f.RegularSkillCalc, "KnightSkillAddLife", "iii>di", vitality, energy, partyBonus, addLifeRate, duration)
}

func WizardMagicDefense(index, targetIndex, dexterity, energy int, effect *float64, duration *int) {
	call(f.RegularSkillCalc, "WizardMagicDefense", "iiii>di", index, targetIndex, dexterity, energy, effect, duration)
}

func DarkLordCriticalDamage(leadership, energy int, effect, duration *int) {
	call(f.RegularSkillCalc, "DarkLordCriticalDamage", "ii>ii", leadership, energy, effect, duration)
}

func ElfShieldRecovery(energy, level int, effect *float64) {
	call(f.RegularSkillCalc, "ElfShieldRecovery", "ii>d", energy, level, effect)
}

func SummonerDrainLifeMonster(energy, monsterLevel int, addHP *int) {
	call(f.RegularSkillCalc, "SummonerDrainLife_Monster", "ii>i", energy, monsterLevel, addHP)
}

func SummonerDamageReflect(energy int, effect, duration *int) {
	call(f.RegularSkillCalc, "SummonerDamageReflect", "i>ii", energy, effect, duration)
}

func SummonerBerserker(energy int, up, down, duration *int) {
	call(f.RegularSkillCalc, "SummonerBerserker", "i>iii", energy, up, down, duration)
}

func SummonerBerserkerAttackDamage(strength, dexterity int, min, max *int) {
	call(f.RegularSkillCalc, "SummonerBerserkerAttackDamage", "ii>ii", strength, dexterity, min, max)
}

func SummonerBerserkerMagicDamage(effect, energy int, min, max *float64) {
	call(f.RegularSkillCalc, "SummonerBerserkerMagicDamage", "ii>dd", effect, energy, min, max)
}

func SummonerBerserkerCurseDamage(effect, energy int, min, max *float64) {
	call(f.RegularSkillCalc, "SummonerBerserkerCurseDamage", "ii>dd", effect, energy, min, max)
}

func SleepMonster(energy, curse, monsterLevel int, rate, duration *int) {
	call(f.RegularSkillCalc, "Sleep_Monster", "iii>ii", energy, curse, monsterLevel, rate, duration)
}

func SummonerWeaknessMonster(energy, curse, monsterLevel int, rate, effect, duration *int) {
	call(f.RegularSkillCalc, "SummonerWeakness_Monster", "iii>iii", energy, curse, monsterLevel, rate, effect, duration)
}

func SummonerInnovationMonster(energy, curse, monsterLevel int, rate, effect, duration *int) {
	call(f.RegularSkillCalc, "SummonerInnovation_Monster", "iii>iii", energy, curse, monsterLevel, rate, effect, duration)
}

func ExplosionDotDamage(damage, masterEffect int, dot, duration *int) {
	call(f.RegularSkillCalc, "ExplosionDotDamage", "ii>ii", damage, masterEffect, dot, duration)
}

func RequiemDotDamage(damage int, dot, duration *int) {
	call(f.RegularSkillCalc, "RequiemDotDamage", "i>ii", damage, dot, duration)
}

func FenrirSkillCalc(damage, level, masterLevel int, out *int) {
	call(f.RegularSkillCalc, "FenrirSkillCalc", "iii>i", damage, level, masterLevel, out)
}
