package formula

// wizard
func WizardDamageCalc(strengh, dexterity, vitality, energy int, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax *int) {
	call(f.CalcCharacter, "WizardDamageCalc", "iiii>iiii", strengh, dexterity, vitality, energy, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax)
}

func WizardMagicDamageCalc(energy int, skillAttackMin, skillAttackMax *int) {
	call(f.CalcCharacter, "WizardDamageCalc", "i>ii", energy, skillAttackMin, skillAttackMax)
}

// knight
func KnightDamageCalc(strengh, dexterity, vitality, energy int, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax *int) {
	call(f.CalcCharacter, "KnightDamageCalc", "iiii>iiii", strengh, dexterity, vitality, energy, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax)
}

func KnightMagicDamageCalc(energy int, skillAttackMin, skillAttackMax *int) {
	call(f.CalcCharacter, "KnightMagicDamageCalc", "i>ii", energy, skillAttackMin, skillAttackMax)
}

// elf
func ElfWithBowDamageCalc(strengh, dexterity, vitality, energy int, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax *int) {
	call(f.CalcCharacter, "ElfWithBowDamageCalc", "iiii>iiii", strengh, dexterity, vitality, energy, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax)
}

func ElfWithoutBowDamageCalc(strengh, dexterity, vitality, energy int, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax *int) {
	call(f.CalcCharacter, "ElfWithoutBowDamageCalc", "iiii>iiii", strengh, dexterity, vitality, energy, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax)
}

func ElfMagicDamageCalc(energy int, skillAttackMin, skillAttackMax *int) {
	call(f.CalcCharacter, "ElfMagicDamageCalc", "i>ii", energy, skillAttackMin, skillAttackMax)
}

// magumsa
func GladiatorDamageCalc(strengh, dexterity, vitality, energy int, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax *int) {
	call(f.CalcCharacter, "GladiatorDamageCalc", "iiii>iiii", strengh, dexterity, vitality, energy, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax)
}

func GladiatorMagicDamageCalc(energy int, skillAttackMin, skillAttackMax *int) {
	call(f.CalcCharacter, "GladiatorMagicDamageCalc", "i>ii", energy, skillAttackMin, skillAttackMax)
}

// darkload
func LordDamageCalc(strengh, dexterity, vitality, energy, leadership int, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax *int) {
	call(f.CalcCharacter, "LordDamageCalc", "iiiii>iiii", strengh, dexterity, vitality, energy, leadership, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax)
}

func LordMagicDamageCalc(energy int, skillAttackMin, skillAttackMax *int) {
	call(f.CalcCharacter, "LordMagicDamageCalc", "i>ii", energy, skillAttackMin, skillAttackMax)
}

// summoner
func SummonerDamageCalc(strengh, dexterity, vitality, energy int, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax *int) {
	call(f.CalcCharacter, "SummonerDamageCalc", "iiii>iiii", strengh, dexterity, vitality, energy, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax)
}

func SummonerMagicDamageCalc(energy int, skillAttackMin, skillAttackMax, curseAttackMin, curseAttackMax *int) {
	call(f.CalcCharacter, "SummonerMagicDamageCalc", "i>iiii", energy, skillAttackMin, skillAttackMax, curseAttackMin, curseAttackMax)
}

// ragefighter
func RageFighterDamageCalc(strengh, dexterity, vitality, energy int, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax *int) {
	call(f.CalcCharacter, "RageFighterDamageCalc", "iiii>iiii", strengh, dexterity, vitality, energy, leftAttackMin, rightAttackMin, leftAttackMax, rightAttackMax)
}

func RageFighterMagicDamageCalc(energy int, skillAttackMin, skillAttackMax *int) {
	call(f.CalcCharacter, "RageFighterMagicDamageCalc", "i>ii", energy, skillAttackMin, skillAttackMax)
}

// defense
func CalcDefense(class, dexterity int, defense *int) {
	call(f.CalcCharacter, "CalcDefense", "ii>i", class, dexterity, defense)
}

// attack rate
func CalcAttackSuccessRate_PvM(class, strengh, dexterity, leadership, level int, attackRate *int) {
	call(f.CalcCharacter, "CalcAttackSuccessRate_PvM", "iiiii>i", class, strengh, dexterity, leadership, level, attackRate)
}

func CalcAttackSuccessRate_PvP(class, dexterity, level int, attackRatePVP *int) {
	call(f.CalcCharacter, "CalcAttackSuccessRate_PvP", "iii>i", class, dexterity, level, attackRatePVP)
}

// defense rate
func CalcDefenseSuccessRate_PvM(class, dexterity int, defenseRate *int) {
	call(f.CalcCharacter, "CalcDefenseSuccessRate_PvM", "ii>i", class, dexterity, defenseRate)
}

func CalcDefenseSuccessRate_PvP(class, dexterity, level int, defenseRatePVP *int) {
	call(f.CalcCharacter, "CalcDefenseSuccessRate_PvP", "iii>i", class, dexterity, level, defenseRatePVP)
}

// speed
func CalcAttackSpeed(class, dexterity int, attackSpeed, magicSpeed *int) {
	call(f.CalcCharacter, "CalcAttackSpeed", "ii>ii", class, dexterity, attackSpeed, magicSpeed)
}
