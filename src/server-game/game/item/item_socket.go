package item

import (
	"fmt"
	"math/rand"

	"github.com/xujintao/balgass/src/server-game/conf"
)

var SocketManager socketManager

type socketSeedOption struct {
	name        string
	element     int
	itemLevel   int
	valueType   int
	values      [5]int
	poolValue   int
	itemSection int
	itemIndex   int
}

type socketBonusOption struct {
	name       string
	minSection int
	maxSection int
	valueType  int
	value      int
	rate       int
	elements   [5]int
}

type socketSpecialOption struct {
	name       string
	kind       int
	valueType  int
	value      int
	rate       int
	seedCounts [6]int
}

type socketManager struct {
	sphereEnabled  bool
	sphereRates    [5]int
	sphereLevels   [5]int
	tetraEnabled   bool
	tetraRates     [3]int
	tetraMinLevels [3]int
	tetraMaxLevels [3]int
	seedMoney      [4]int
	rates          [12][5]int
	twoHand        [5]int
	seedOptions    map[int]socketSeedOption
	bonusOptions   map[int]socketBonusOption
	specialOptions map[int]socketSpecialOption
}

func (m *socketManager) init() {
	type sphereValues struct {
		Level1 *int `xml:"Level1,attr"`
		Level2 *int `xml:"Level2,attr"`
		Level3 *int `xml:"Level3,attr"`
		Level4 *int `xml:"Level4,attr"`
		Level5 *int `xml:"Level5,attr"`
	}
	type rateConfig struct {
		CatID  int `xml:"CatID,attr"`
		Count1 int `xml:"Count1,attr"`
		Count2 int `xml:"Count2,attr"`
		Count3 int `xml:"Count3,attr"`
		Count4 int `xml:"Count4,attr"`
		Count5 int `xml:"Count5,attr"`
	}
	type tetraConfig struct {
		Rate *int `xml:"Rate,attr"`
		Min  *int `xml:"LevelMin,attr"`
		Max  *int `xml:"LevelMax,attr"`
	}
	type socketSystemConfig struct {
		SphereDrop *int `xml:"SphereDrop,attr"`
		TetraDrop  *int `xml:"TetraDrop,attr"`
		Sphere     struct {
			Rate  *sphereValues `xml:"Rate"`
			Level *sphereValues `xml:"Level"`
		} `xml:"SphereDrop"`
		Tetra struct {
			Bronze *tetraConfig `xml:"Bronze"`
			Silver *tetraConfig `xml:"Silver"`
			Gold   *tetraConfig `xml:"Gold"`
		} `xml:"TetraDrop"`
		Seed *struct {
			ExtractMoney *int `xml:"ExtractMoney,attr"`
			Sphere       *struct {
				CompositeMoney *int `xml:"CompositeMoney,attr"`
				SetMoney       *int `xml:"SetMoney,attr"`
				RemoveMoney    *int `xml:"RemoveMoney,attr"`
			} `xml:"Sphere"`
		} `xml:"Seed"`
		SocketRateSettings struct {
			Sockets []rateConfig `xml:"Socket"`
			TwoHand *rateConfig  `xml:"TwoHandWeapons"`
		} `xml:"SocketRateSettings"`
	}
	var socketCfg socketSystemConfig
	const socketFile = "SocketSystem/IGC_SocketSystem.xml"
	conf.XML(conf.PathCommon, socketFile, &socketCfg)
	if socketCfg.SphereDrop == nil || (*socketCfg.SphereDrop != 0 && *socketCfg.SphereDrop != 1) {
		panic(fmt.Errorf("%s: invalid SphereDrop switch", socketFile))
	}
	if socketCfg.Sphere.Rate == nil || socketCfg.Sphere.Level == nil {
		panic(fmt.Errorf("%s: missing SphereDrop Rate or Level", socketFile))
	}
	m.sphereEnabled = *socketCfg.SphereDrop == 1
	rates := socketCfg.Sphere.Rate
	levels := socketCfg.Sphere.Level
	rateValues := [5]*int{rates.Level1, rates.Level2, rates.Level3, rates.Level4, rates.Level5}
	levelValues := [5]*int{levels.Level1, levels.Level2, levels.Level3, levels.Level4, levels.Level5}
	for i := range m.sphereRates {
		if rateValues[i] == nil || *rateValues[i] < 0 || *rateValues[i] > 10000 || levelValues[i] == nil || *levelValues[i] < 0 {
			panic(fmt.Errorf("%s: invalid SphereDrop level %d", socketFile, i+1))
		}
		base, err := ItemTable.GetItemBase(12, 70+i)
		if err != nil || base.KindB != KindBSocketSphere {
			panic(fmt.Errorf("%s: missing Sphere item 12/%d: %v", socketFile, 70+i, err))
		}
		m.sphereRates[i] = *rateValues[i]
		m.sphereLevels[i] = *levelValues[i]
	}
	if socketCfg.TetraDrop == nil || (*socketCfg.TetraDrop != 0 && *socketCfg.TetraDrop != 1) {
		panic(fmt.Errorf("%s: invalid TetraDrop switch", socketFile))
	}
	m.tetraEnabled = *socketCfg.TetraDrop == 1
	entries := [3]*tetraConfig{socketCfg.Tetra.Bronze, socketCfg.Tetra.Silver, socketCfg.Tetra.Gold}
	for i, entry := range entries {
		if entry == nil || entry.Rate == nil || entry.Min == nil || entry.Max == nil ||
			*entry.Rate < 0 || *entry.Rate > 10000 || *entry.Min < 0 || *entry.Max < *entry.Min {
			panic(fmt.Errorf("%s: invalid TetraDrop tier %d", socketFile, i+1))
		}
		if _, err := ItemTable.GetItemBase(14, 247-i); err != nil {
			panic(fmt.Errorf("%s: Tetra tier %d: %w", socketFile, i+1, err))
		}
		m.tetraRates[i], m.tetraMinLevels[i], m.tetraMaxLevels[i] = *entry.Rate, *entry.Min, *entry.Max
	}
	if socketCfg.Seed == nil || socketCfg.Seed.Sphere == nil ||
		socketCfg.Seed.ExtractMoney == nil || socketCfg.Seed.Sphere.CompositeMoney == nil ||
		socketCfg.Seed.Sphere.SetMoney == nil || socketCfg.Seed.Sphere.RemoveMoney == nil {
		panic(fmt.Errorf("%s: missing Seed fees", socketFile))
	}
	fees := [4]*int{socketCfg.Seed.ExtractMoney, socketCfg.Seed.Sphere.CompositeMoney,
		socketCfg.Seed.Sphere.SetMoney, socketCfg.Seed.Sphere.RemoveMoney}
	for i, fee := range fees {
		if *fee < 0 {
			panic(fmt.Errorf("%s: invalid Seed fee %d", socketFile, i))
		}
		m.seedMoney[i] = *fee
	}
	weights := func(label string, rate rateConfig) [5]int {
		values := [5]int{rate.Count1, rate.Count2, rate.Count3, rate.Count4, rate.Count5}
		total := 0
		for _, value := range values {
			if value < 0 || value > 10000 {
				panic(fmt.Errorf("%s: invalid %s slot weight %d", socketFile, label, value))
			}
			total += value
		}
		if total == 0 {
			panic(fmt.Errorf("%s: %s has no positive slot weights", socketFile, label))
		}
		return values
	}
	var seen [12]bool
	for _, rate := range socketCfg.SocketRateSettings.Sockets {
		if rate.CatID < 0 || rate.CatID >= len(m.rates) {
			panic(fmt.Errorf("%s: invalid section %d", socketFile, rate.CatID))
		}
		if seen[rate.CatID] {
			panic(fmt.Errorf("%s: duplicate section %d", socketFile, rate.CatID))
		}
		m.rates[rate.CatID] = weights(fmt.Sprintf("section %d", rate.CatID), rate)
		seen[rate.CatID] = true
	}
	for section, present := range seen {
		if !present {
			panic(fmt.Errorf("%s: missing section %d", socketFile, section))
		}
	}
	if socketCfg.SocketRateSettings.TwoHand == nil {
		panic(fmt.Errorf("%s: missing TwoHandWeapons", socketFile))
	}
	m.twoHand = weights("TwoHandWeapons", *socketCfg.SocketRateSettings.TwoHand)

	type seedConfig struct {
		Index       *int   `xml:"Index,attr"`
		ElementType *int   `xml:"ElementType,attr"`
		Level       *int   `xml:"Level,attr"`
		BonusType   *int   `xml:"BonusType,attr"`
		Value1      *int   `xml:"BonusValue1,attr"`
		Value2      *int   `xml:"BonusValue2,attr"`
		Value3      *int   `xml:"BonusValue3,attr"`
		Value4      *int   `xml:"BonusValue4,attr"`
		Value5      *int   `xml:"BonusValue5,attr"`
		PoolValue   *int   `xml:"PoolValue,attr"`
		Name        string `xml:"Name,attr"`
	}
	type bonusConfig struct {
		Index        *int   `xml:"Index,attr"`
		MinItemType  *int   `xml:"MinItemType,attr"`
		MaxItemType  *int   `xml:"MaxItemType,attr"`
		BonusType    *int   `xml:"BonusType,attr"`
		BonusValue   *int   `xml:"BonusValue,attr"`
		BonusRate    *int   `xml:"BonusRate,attr"`
		ElementType1 *int   `xml:"ElementType1,attr"`
		ElementType2 *int   `xml:"ElementType2,attr"`
		ElementType3 *int   `xml:"ElementType3,attr"`
		ElementType4 *int   `xml:"ElementType4,attr"`
		ElementType5 *int   `xml:"ElementType5,attr"`
		Name         string `xml:"Name,attr"`
	}
	type specialConfig struct {
		Index          *int   `xml:"Index,attr"`
		Type           *int   `xml:"Type,attr"`
		BonusType      *int   `xml:"BonusType,attr"`
		BonusValue     *int   `xml:"BonusValue,attr"`
		BonusRate      *int   `xml:"BonusRate,attr"`
		SeedCountFire  *int   `xml:"SeedCountFire,attr"`
		SeedCountWater *int   `xml:"SeedCountWater,attr"`
		SeedCountIce   *int   `xml:"SeedCountIce,attr"`
		SeedCountWind  *int   `xml:"SeedCountWind,attr"`
		SeedCountLight *int   `xml:"SeedCountLight,attr"`
		SeedCountEarth *int   `xml:"SeedCountEarth,attr"`
		Name           string `xml:"Name,attr"`
	}
	type optionSettingsConfig struct {
		Seeds    []seedConfig    `xml:"SocketItemOptionSettings>Option"`
		Bonuses  []bonusConfig   `xml:"SocketBonusSettings>Option"`
		Specials []specialConfig `xml:"SpecialOptionSettings>Option"`
	}
	var optionCfg optionSettingsConfig
	const file = "SocketSystem/IGC_SocketOption.xml"
	conf.XML(conf.PathCommon, file, &optionCfg)
	required := func(group string, index int, name string, value *int) int {
		if value == nil {
			panic(fmt.Errorf("%s: %s option %d missing %s", file, group, index, name))
		}
		return *value
	}
	optionIndex := func(group string, value *int, max int) int {
		index := required(group, -1, "Index", value)
		if index < 0 || index >= max {
			panic(fmt.Errorf("%s: %s option invalid index %d", file, group, index))
		}
		return index
	}

	if len(optionCfg.Seeds) == 0 || len(optionCfg.Bonuses) == 0 || len(optionCfg.Specials) == 0 {
		panic(fmt.Errorf("%s: missing Socket option settings", file))
	}
	m.seedOptions = make(map[int]socketSeedOption, len(optionCfg.Seeds))
	totalSeedWeight := 0
	for _, option := range optionCfg.Seeds {
		index := optionIndex("seed", option.Index, 50)
		if _, exists := m.seedOptions[index]; exists {
			panic(fmt.Errorf("%s: duplicate seed option %d", file, index))
		}
		element := required("seed", index, "ElementType", option.ElementType)
		itemLevel := required("seed", index, "Level", option.Level)
		valueType := required("seed", index, "BonusType", option.BonusType)
		values := [5]int{
			required("seed", index, "BonusValue1", option.Value1),
			required("seed", index, "BonusValue2", option.Value2),
			required("seed", index, "BonusValue3", option.Value3),
			required("seed", index, "BonusValue4", option.Value4),
			required("seed", index, "BonusValue5", option.Value5),
		}
		weight := required("seed", index, "PoolValue", option.PoolValue)
		if element < 1 || element > 6 || itemLevel < 0 || itemLevel > 5 || valueType < 0 || valueType > 5 || weight < 0 || option.Name == "" {
			panic(fmt.Errorf("%s: invalid seed option %d", file, index))
		}
		for _, value := range values {
			if value < 0 || value > 65535 || valueType != 0 && value == 0 {
				panic(fmt.Errorf("%s: invalid seed option %d value %d", file, index, value))
			}
		}
		itemSection, itemIndex := 12, 59+element
		if _, err := ItemTable.GetItemBase(itemSection, itemIndex); err != nil {
			panic(fmt.Errorf("%s: seed option %d: %w", file, index, err))
		}
		m.seedOptions[index] = socketSeedOption{
			name: option.Name, element: element, itemLevel: itemLevel, valueType: valueType,
			values: values, poolValue: weight, itemSection: itemSection, itemIndex: itemIndex,
		}
		totalSeedWeight += weight
	}
	if totalSeedWeight == 0 {
		panic(fmt.Errorf("%s: seed options have no positive PoolValue", file))
	}
	m.bonusOptions = make(map[int]socketBonusOption, len(optionCfg.Bonuses))
	for _, option := range optionCfg.Bonuses {
		index := optionIndex("bonus", option.Index, 12)
		if _, exists := m.bonusOptions[index]; exists {
			panic(fmt.Errorf("%s: duplicate bonus option %d", file, index))
		}
		minSection := required("bonus", index, "MinItemType", option.MinItemType)
		maxSection := required("bonus", index, "MaxItemType", option.MaxItemType)
		valueType := required("bonus", index, "BonusType", option.BonusType)
		value := required("bonus", index, "BonusValue", option.BonusValue)
		rate := required("bonus", index, "BonusRate", option.BonusRate)
		elements := [5]int{
			required("bonus", index, "ElementType1", option.ElementType1),
			required("bonus", index, "ElementType2", option.ElementType2),
			required("bonus", index, "ElementType3", option.ElementType3),
			required("bonus", index, "ElementType4", option.ElementType4),
			required("bonus", index, "ElementType5", option.ElementType5),
		}
		if minSection < 0 || maxSection < minSection || maxSection > 12 || valueType < 0 || valueType > 5 || value < 0 || value > 65535 || rate < 0 || rate > 100 || option.Name == "" {
			panic(fmt.Errorf("%s: invalid bonus option %d", file, index))
		}
		for _, element := range elements {
			if element < 0 || element > 6 {
				panic(fmt.Errorf("%s: bonus option %d has invalid element %d", file, index, element))
			}
		}
		m.bonusOptions[index] = socketBonusOption{
			name: option.Name, minSection: minSection, maxSection: maxSection,
			valueType: valueType, value: value, rate: rate, elements: elements,
		}
	}

	m.specialOptions = make(map[int]socketSpecialOption, len(optionCfg.Specials))
	for _, option := range optionCfg.Specials {
		index := optionIndex("special", option.Index, 2)
		if _, exists := m.specialOptions[index]; exists {
			panic(fmt.Errorf("%s: duplicate special option %d", file, index))
		}
		kind := required("special", index, "Type", option.Type)
		valueType := required("special", index, "BonusType", option.BonusType)
		value := required("special", index, "BonusValue", option.BonusValue)
		rate := required("special", index, "BonusRate", option.BonusRate)
		counts := [6]int{
			required("special", index, "SeedCountFire", option.SeedCountFire),
			required("special", index, "SeedCountWater", option.SeedCountWater),
			required("special", index, "SeedCountIce", option.SeedCountIce),
			required("special", index, "SeedCountWind", option.SeedCountWind),
			required("special", index, "SeedCountLight", option.SeedCountLight),
			required("special", index, "SeedCountEarth", option.SeedCountEarth),
		}
		if kind < 0 || kind > 255 || valueType < 0 || valueType > 5 || value < 0 || value > 65535 || rate < 0 || rate > 100 || option.Name == "" {
			panic(fmt.Errorf("%s: invalid special option %d", file, index))
		}
		for _, count := range counts {
			if count < 0 || count > 255 {
				panic(fmt.Errorf("%s: special option %d has invalid seed count %d", file, index, count))
			}
		}
		m.specialOptions[index] = socketSpecialOption{
			name: option.Name, kind: kind, valueType: valueType, value: value, rate: rate, seedCounts: counts,
		}
	}
}

func (m *socketManager) RandomSphereIndex(monsterLevel int) (int, bool) {
	if !m.sphereEnabled {
		return 0, false
	}
	highest := 0
	for i := len(m.sphereLevels) - 1; i >= 0; i-- {
		if m.sphereLevels[i] > 0 && monsterLevel >= m.sphereLevels[i] {
			highest = i + 1
			break
		}
	}
	index := -1
	for i := 0; i < highest; i++ {
		if rand.Intn(10000) < m.sphereRates[i] {
			index = 70 + i
		}
	}
	return index, index >= 0
}

func (m *socketManager) RandomTetraIndex(monsterLevel int) (int, bool) {
	if !m.tetraEnabled {
		return 0, false
	}
	highest := 0
	for i := len(m.tetraRates) - 1; i >= 0; i-- {
		if monsterLevel >= m.tetraMinLevels[i] && monsterLevel <= m.tetraMaxLevels[i] {
			highest = i + 1
			break
		}
	}
	index := -1
	for i := 0; i < highest; i++ {
		if rand.Intn(10000) < m.tetraRates[i] {
			index = 247 - i
		}
	}
	return index, index >= 0
}

// SeedFee returns the configured Zen cost for a Seed mix fee index (0 through 3).
func (m *socketManager) SeedFee(index int) int {
	return m.seedMoney[index]
}

func (m *socketManager) RandomSeed() (int, int, int) {
	total := 0
	for _, seed := range m.seedOptions {
		total += seed.poolValue
	}
	roll := rand.Intn(total)
	for index := 0; index < 50; index++ {
		seed, ok := m.seedOptions[index]
		if !ok {
			continue
		}
		roll -= seed.poolValue
		if roll < 0 {
			return seed.itemSection, seed.itemIndex, seed.itemLevel
		}
	}
	panic("RandomSeed: positive seed weight was not selected")
}

func (m *socketManager) SeedIndex(itemSection, itemIndex, level int) (int, bool) {
	for index, seed := range m.seedOptions {
		if seed.itemSection == itemSection && seed.itemIndex == itemIndex && seed.itemLevel == level {
			return index, true
		}
	}
	return 0, false
}

func (m *socketManager) SeedSphereIndex(itemSection, itemIndex, level int) (int, int, bool) {
	if itemSection != 12 || itemIndex < 100 || itemIndex > 129 {
		return 0, 0, false
	}
	sphereLevel := (itemIndex-100)/6 + 1
	index, ok := m.SeedIndex(12, 60+(itemIndex-100)%6, level)
	return index, sphereLevel, ok
}

func (m *socketManager) SeedSphereItem(index, sphereLevel int) *Item {
	seed, ok := m.seedOptions[index]
	if !ok || sphereLevel < 1 || sphereLevel > 5 {
		panic(fmt.Errorf("SeedSphereItem: invalid option %d or sphere level %d", index, sphereLevel))
	}
	it := NewItem(12, 100+(sphereLevel-1)*6+seed.element-1)
	it.Level = seed.itemLevel
	it.Durability = 1
	it.Calc()
	return it
}

func (m *socketManager) SeedElement(index int) (int, bool) {
	seed, ok := m.seedOptions[index]
	return seed.element, ok
}

func (m *socketManager) SeedItem(index int) (int, int, int, bool) {
	seed, ok := m.seedOptions[index]
	return seed.itemSection, seed.itemIndex, seed.itemLevel, ok
}

func (it *Item) SocketSlotAt(position int) (int, bool) {
	switch position {
	case 0:
		return it.SocketSlot1, true
	case 1:
		return it.SocketSlot2, true
	case 2:
		return it.SocketSlot3, true
	case 3:
		return it.SocketSlot4, true
	case 4:
		return it.SocketSlot5, true
	default:
		return 0, false
	}
}

func (it *Item) SetSocketSlotAt(position, value int) bool {
	switch position {
	case 0:
		it.SocketSlot1 = value
	case 1:
		it.SocketSlot2 = value
	case 2:
		it.SocketSlot3 = value
	case 3:
		it.SocketSlot4 = value
	case 4:
		it.SocketSlot5 = value
	default:
		return false
	}
	return true
}

func (m *socketManager) RollBonus(it *Item) {
	it.SocketBonus = 0
	var elements [5]int
	highLevelFirstThree := 0
	for i, slot := range [5]int{it.SocketSlot1, it.SocketSlot2, it.SocketSlot3, it.SocketSlot4, it.SocketSlot5} {
		if slot < 100 {
			continue
		}
		if i < 3 && slot/100 == 3 {
			highLevelFirstThree++
		}
		if seed, ok := m.seedOptions[slot%100]; ok {
			elements[i] = seed.element
		}
	}
	for index := 0; index < 12; index++ {
		bonus, ok := m.bonusOptions[index]
		if !ok || it.Section < bonus.minSection || it.Section > bonus.maxSection {
			continue
		}
		if index >= 6 && highLevelFirstThree != 3 {
			continue
		}
		matches := 0
		required := 0
		for i, element := range elements {
			if bonus.elements[i] != 0 {
				required++
			}
			if element != 0 && element == bonus.elements[i] {
				matches++
			}
		}
		if required > 0 && matches >= required {
			if rand.Intn(100) < bonus.rate {
				it.SocketBonus = index + 1
				return
			}
			if index >= 6 || highLevelFirstThree != 3 {
				return
			}
		}
	}
}

func (m *socketManager) RandomSlotCount(it *Item) int {
	if it == nil || it.ItemBase == nil || it.Type != TypeSocket {
		return 0
	}
	var rates [5]int
	if it.TwoHand {
		rates = m.twoHand
	} else if it.Section >= 0 && it.Section < len(m.rates) {
		rates = m.rates[it.Section]
	}
	total := 0
	for _, rate := range rates {
		total += rate
	}
	if total == 0 {
		return 0
	}
	n := rand.Intn(total)
	for i, rate := range rates {
		if n < rate {
			return i + 1
		}
		n -= rate
	}
	return 0
}

func (it *Item) MakeSocketSlots(count int) {
	if it == nil ||
		it.ItemBase == nil ||
		it.Type != TypeSocket ||
		count < 0 || count > 5 {
		return
	}
	it.SocketBonus = 0
	for i, slot := range []*int{
		&it.SocketSlot1,
		&it.SocketSlot2,
		&it.SocketSlot3,
		&it.SocketSlot4,
		&it.SocketSlot5} {
		*slot = 0
		if i < count {
			*slot = 1
		}
	}
}

type SocketEffectType int

const (
	SocketEffectLevelAttack     SocketEffectType = 0  // 等级火
	SocketEffectAttackSpeed     SocketEffectType = 1  // 速度火
	SocketEffectAttackMax       SocketEffectType = 2  // 最大火
	SocketEffectAttackMin       SocketEffectType = 3  // 最小火
	SocketEffectAttack          SocketEffectType = 4  // 攻击火
	SocketEffectReduceAGCost    SocketEffectType = 5  // 减少AG火
	SocketEffectDefenseRate     SocketEffectType = 10 // 防御率水
	SocketEffectDefense         SocketEffectType = 11 // 防御力水
	SocketEffectShieldDefense   SocketEffectType = 12 // 盾防水
	SocketEffectReduceDamage    SocketEffectType = 13 // 减伤水
	SocketEffectReflectDamage   SocketEffectType = 14 // 反伤水
	SocketEffectKillHP          SocketEffectType = 16 // 杀怪回HP冰
	SocketEffectKillMP          SocketEffectType = 17 // 杀怪回MP冰
	SocketEffectSkillAttack     SocketEffectType = 18 // 技能攻击冰
	SocketEffectAttackRate      SocketEffectType = 19 // 攻击率冰
	SocketEffectItemDurability  SocketEffectType = 20 // 道具耐久冰
	SocketEffectRecoverHP       SocketEffectType = 21 // HP恢复风
	SocketEffectMaxHP           SocketEffectType = 22 // 最大HP风
	SocketEffectMaxMP           SocketEffectType = 23 // 最大MP风
	SocketEffectRecoverMP       SocketEffectType = 24 // MP恢复风
	SocketEffectMaxAG           SocketEffectType = 25 // 最大AG风
	SocketEffectRecoverAG       SocketEffectType = 26 // AG恢复风
	SocketEffectExcellentDamage SocketEffectType = 29 // 卓越一击伤害雷
	SocketEffectExcellentRate   SocketEffectType = 30 // 卓越一击率雷
	SocketEffectCriticalDamage  SocketEffectType = 31 // 致命一击伤害雷
	SocketEffectCriticalRate    SocketEffectType = 32 // 致命一击率雷
	SocketEffectVitality        SocketEffectType = 36 // 体力土
	SocketEffectDoubleDamage    SocketEffectType = 44 // 双倍伤害
	SocketEffectIgnoreDefense   SocketEffectType = 45 // 忽视防御
)

type SocketBonusEffectType int

const (
	SocketBonusWeaponAttack       SocketBonusEffectType = 0
	SocketBonusWeaponSkillAttack  SocketBonusEffectType = 1
	SocketBonusStaffAttack        SocketBonusEffectType = 2
	SocketBonusStaffSkillAttack   SocketBonusEffectType = 3
	SocketBonusArmorDefense       SocketBonusEffectType = 4
	SocketBonusArmorMaxHP         SocketBonusEffectType = 5
	SocketBonusWeaponAttack2      SocketBonusEffectType = 6
	SocketBonusWeaponSkillAttack2 SocketBonusEffectType = 7
	SocketBonusStaffAttack2       SocketBonusEffectType = 8
	SocketBonusStaffSkillAttack2  SocketBonusEffectType = 9
	SocketBonusArmorDefense2      SocketBonusEffectType = 10
	SocketBonusArmorMaxHP2        SocketBonusEffectType = 11
)

// SocketEffect is a configured effect resolved to one value.
type SocketEffect struct {
	Index     SocketEffectType
	Element   int
	ValueType int
	Value     int
}

func (m *socketManager) SeedEffect(slot int) (SocketEffect, bool) {
	if slot < 100 {
		return SocketEffect{}, false
	}
	level, index := slot/100, slot%100
	if level < 1 || level > 5 || index >= 50 {
		return SocketEffect{}, false
	}
	option, ok := m.seedOptions[index]
	if !ok {
		return SocketEffect{}, false
	}
	return SocketEffect{
		Index:     SocketEffectType(index),
		Element:   option.element,
		ValueType: option.valueType,
		Value:     option.values[level-1]}, true
}

type SocketBonusEffect struct {
	Index      SocketBonusEffectType
	MinSection int
	MaxSection int
	ValueType  int
	Value      int
}

func (m *socketManager) BonusEffect(bonus int) (SocketBonusEffect, bool) {
	if bonus <= 0 {
		return SocketBonusEffect{}, false
	}
	index := bonus - 1
	option, ok := m.bonusOptions[index]
	if !ok {
		return SocketBonusEffect{}, false
	}
	return SocketBonusEffect{
		Index:      SocketBonusEffectType(index),
		MinSection: option.minSection,
		MaxSection: option.maxSection,
		ValueType:  option.valueType,
		Value:      option.value}, true
}

func (m *socketManager) SpecialEffects(counts [6]int) []SocketEffect {
	var effects []SocketEffect
	for index, effectType := range []SocketEffectType{SocketEffectDoubleDamage, SocketEffectIgnoreDefense} {
		option, ok := m.specialOptions[index]
		if !ok {
			continue
		}
		matched := true
		for element, required := range option.seedCounts {
			if counts[element] < required {
				matched = false
				break
			}
		}
		if matched {
			effects = append(effects, SocketEffect{Index: effectType, ValueType: option.valueType, Value: option.value})
		}
	}
	return effects
}
