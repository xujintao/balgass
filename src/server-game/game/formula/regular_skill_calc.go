package formula

func Knight_Gladiator_CalcSkillBonus(class, energy int, rate *float64) {
	call(f.RegularSkillCacl, "Knight_Gladiator_CalcSkillBonus", "iii>d", class, 1, energy, rate)
}

func ImpaleSkillCalc(class, energy int, rate *float64) {
	call(f.RegularSkillCacl, "ImpaleSkillCalc", "iii>d", class, 1, energy, rate)
}

func Elf_CalcSkillBonus(damage, energy int, out *int) {
	call(f.RegularSkillCacl, "Elf_CalcSkillBonus", "ii>i", damage, energy, out)
}

func GladiatorPowerSlash(damage, energy int, out *int) {
	call(f.RegularSkillCacl, "GladiatorPowerSlash", "ii>i", damage, energy, out)
}

func Lord_CalcSkillBonus(damage, energy int, out *int) {
	call(f.RegularSkillCacl, "Lord_CalcSkillBonus", "ii>i", damage, energy, out)
}

func ChainLightningCalc(damage, targetNumber int, out *int) {
	call(f.RegularSkillCacl, "ChainLightningCalc", "ii>i", damage, targetNumber, out)
}

func StrikeOfDestructionCalc(damage, energy int, out *int) {
	call(f.RegularSkillCacl, "StrikeOfDestructionCalc", "ii>i", damage, energy, out)
}

func FlameStrikeCalc(damage int, out *int) {
	call(f.RegularSkillCacl, "FlameStrikeCalc", "i>i", damage, out)
}

func ChaoticDiseierCalc(damage, energy int, out *int) {
	call(f.RegularSkillCacl, "ChaoticDiseierCalc", "ii>i", damage, energy, out)
}

func RageFighterKillingBlow(damage, vitality int, out *int) {
	call(f.RegularSkillCacl, "RageFighterKillingBlow", "ii>i", damage, vitality, out)
}

func RageFighterBeastUppercut(damage, vitality int, out *int) {
	call(f.RegularSkillCacl, "RageFighterBeastUppercut", "ii>i", damage, vitality, out)
}

func RageFighterChainDrive(damage, vitality int, out *int) {
	call(f.RegularSkillCacl, "RageFighterChainDrive", "ii>i", damage, vitality, out)
}

func RageFighterDarkSideIncDamage(damage, dexterity, energy int, out *int) {
	call(f.RegularSkillCacl, "RageFighterDarkSideIncDamage", "iii>i", damage, dexterity, energy, out)
}

func RageFighterDragonRoar(damage, energy int, out *int) {
	call(f.RegularSkillCacl, "RageFighterDragonRoar", "ii>i", damage, energy, out)
}

func RageFighterDragonSlasher(damage, energy, targetType int, out *int) {
	call(f.RegularSkillCacl, "RageFighterDragonSlasher", "iii>i", damage, energy, targetType, out)
}

func RageFighterCharge(damage, vitality int, out *int) {
	call(f.RegularSkillCacl, "RageFighterCharge", "ii>i", damage, vitality, out)
}

func RageFighterPhoenixShot(damage, vitality int, out *int) {
	call(f.RegularSkillCacl, "RageFighterPhoenixShot", "ii>i", damage, vitality, out)
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

func WizardMagicDefense(index, targetIndex, dexterity, energy int, effect *float64, duration *int) {
	call(f.RegularSkillCacl, "WizardMagicDefense", "iiii>di", index, targetIndex, dexterity, energy, effect, duration)
}

func DarkLordCriticalDamage(leadership, energy int, effect, duration *int) {
	call(f.RegularSkillCacl, "DarkLordCriticalDamage", "ii>ii", leadership, energy, effect, duration)
}

func ElfShieldRecovery(energy, level int, effect *float64) {
	call(f.RegularSkillCacl, "ElfShieldRecovery", "ii>d", energy, level, effect)
}

func SummonerDrainLifeMonster(energy, monsterLevel int, addHP *int) {
	call(f.RegularSkillCacl, "SummonerDrainLife_Monster", "ii>i", energy, monsterLevel, addHP)
}

func SummonerDamageReflect(energy int, effect, duration *int) {
	call(f.RegularSkillCacl, "SummonerDamageReflect", "i>ii", energy, effect, duration)
}

func SummonerBerserker(energy int, up, down, duration *int) {
	call(f.RegularSkillCacl, "SummonerBerserker", "i>iii", energy, up, down, duration)
}

func SummonerBerserkerAttackDamage(strength, dexterity int, min, max *int) {
	call(f.RegularSkillCacl, "SummonerBerserkerAttackDamage", "ii>ii", strength, dexterity, min, max)
}

func SummonerBerserkerMagicDamage(effect, energy int, min, max *float64) {
	call(f.RegularSkillCacl, "SummonerBerserkerMagicDamage", "ii>dd", effect, energy, min, max)
}

func SummonerBerserkerCurseDamage(effect, energy int, min, max *float64) {
	call(f.RegularSkillCacl, "SummonerBerserkerCurseDamage", "ii>dd", effect, energy, min, max)
}

func SleepMonster(energy, curse, monsterLevel int, rate, duration *int) {
	call(f.RegularSkillCacl, "Sleep_Monster", "iii>ii", energy, curse, monsterLevel, rate, duration)
}

func SummonerWeaknessMonster(energy, curse, monsterLevel int, rate, effect, duration *int) {
	call(f.RegularSkillCacl, "SummonerWeakness_Monster", "iii>iii", energy, curse, monsterLevel, rate, effect, duration)
}

func SummonerInnovationMonster(energy, curse, monsterLevel int, rate, effect, duration *int) {
	call(f.RegularSkillCacl, "SummonerInnovation_Monster", "iii>iii", energy, curse, monsterLevel, rate, effect, duration)
}

func ExplosionDotDamage(damage, masterEffect int, dot, duration *int) {
	call(f.RegularSkillCacl, "ExplosionDotDamage", "ii>ii", damage, masterEffect, dot, duration)
}

func RequiemDotDamage(damage int, dot, duration *int) {
	call(f.RegularSkillCacl, "RequiemDotDamage", "i>ii", damage, dot, duration)
}

func FenrirSkillCalc(damage, level, masterLevel int, out *int) {
	call(f.RegularSkillCacl, "FenrirSkillCalc", "iii>i", damage, level, masterLevel, out)
}
