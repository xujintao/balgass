package item

import (
	"fmt"
	"math/rand"

	"github.com/xujintao/balgass/src/server-game/conf"
)

var SocketManager socketManager

type socketSeedOption struct {
	name         string
	element      int
	itemLevel    int
	valueType    int
	values       [5]int
	poolValue    int
	seedItemCode int
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
	rates          [12][5]int
	twoHand        [5]int
	seedOptions    map[int]socketSeedOption
	bonusOptions   map[int]socketBonusOption
	specialOptions map[int]socketSpecialOption
}

func (m *socketManager) init() {
	type rateConfig struct {
		CatID  int `xml:"CatID,attr"`
		Count1 int `xml:"Count1,attr"`
		Count2 int `xml:"Count2,attr"`
		Count3 int `xml:"Count3,attr"`
		Count4 int `xml:"Count4,attr"`
		Count5 int `xml:"Count5,attr"`
	}
	type rateSettingsConfig struct {
		SocketRateSettings struct {
			Sockets []rateConfig `xml:"Socket"`
			TwoHand *rateConfig  `xml:"TwoHandWeapons"`
		} `xml:"SocketRateSettings"`
	}
	var rateCfg rateSettingsConfig
	conf.XML(conf.PathCommon, "SocketSystem/IGC_SocketSystem.xml", &rateCfg)
	weights := func(label string, rate rateConfig) [5]int {
		values := [5]int{rate.Count1, rate.Count2, rate.Count3, rate.Count4, rate.Count5}
		total := 0
		for _, value := range values {
			if value < 0 || value > 10000 {
				panic(fmt.Errorf("SocketSystem/IGC_SocketSystem.xml: invalid %s slot weight %d", label, value))
			}
			total += value
		}
		if total == 0 {
			panic(fmt.Errorf("SocketSystem/IGC_SocketSystem.xml: %s has no positive slot weights", label))
		}
		return values
	}
	var seen [12]bool
	for _, rate := range rateCfg.SocketRateSettings.Sockets {
		if rate.CatID < 0 || rate.CatID >= len(m.rates) {
			panic(fmt.Errorf("SocketSystem/IGC_SocketSystem.xml: invalid section %d", rate.CatID))
		}
		if seen[rate.CatID] {
			panic(fmt.Errorf("SocketSystem/IGC_SocketSystem.xml: duplicate section %d", rate.CatID))
		}
		m.rates[rate.CatID] = weights(fmt.Sprintf("section %d", rate.CatID), rate)
		seen[rate.CatID] = true
	}
	for section, present := range seen {
		if !present {
			panic(fmt.Errorf("SocketSystem/IGC_SocketSystem.xml: missing section %d", section))
		}
	}
	if rateCfg.SocketRateSettings.TwoHand == nil {
		panic(fmt.Errorf("SocketSystem/IGC_SocketSystem.xml: missing TwoHandWeapons"))
	}
	m.twoHand = weights("TwoHandWeapons", *rateCfg.SocketRateSettings.TwoHand)

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
		seedItemCode := Code(12, 59+element)
		if _, err := ItemTable.GetItemBase(12, 59+element); err != nil {
			panic(fmt.Errorf("%s: seed option %d: %w", file, index, err))
		}
		m.seedOptions[index] = socketSeedOption{
			name: option.Name, element: element, itemLevel: itemLevel, valueType: valueType,
			values: values, poolValue: weight, seedItemCode: seedItemCode,
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
	SocketEffectLevelAttack     SocketEffectType = 0
	SocketEffectAttackSpeed     SocketEffectType = 1
	SocketEffectAttackMax       SocketEffectType = 2
	SocketEffectAttackMin       SocketEffectType = 3
	SocketEffectAttack          SocketEffectType = 4
	SocketEffectReduceAGCost    SocketEffectType = 5
	SocketEffectDefenseRate     SocketEffectType = 10
	SocketEffectDefense         SocketEffectType = 11
	SocketEffectShieldDefense   SocketEffectType = 12
	SocketEffectReduceDamage    SocketEffectType = 13
	SocketEffectReflectDamage   SocketEffectType = 14
	SocketEffectKillHP          SocketEffectType = 16
	SocketEffectKillMP          SocketEffectType = 17
	SocketEffectSkillAttack     SocketEffectType = 18
	SocketEffectAttackRate      SocketEffectType = 19
	SocketEffectRecoverHP       SocketEffectType = 21
	SocketEffectMaxHP           SocketEffectType = 22
	SocketEffectMaxMP           SocketEffectType = 23
	SocketEffectRecoverMP       SocketEffectType = 24
	SocketEffectMaxAG           SocketEffectType = 25
	SocketEffectRecoverAG       SocketEffectType = 26
	SocketEffectKillMoney       SocketEffectType = 27
	SocketEffectExcellentDamage SocketEffectType = 29
	SocketEffectExcellentRate   SocketEffectType = 30
	SocketEffectCriticalDamage  SocketEffectType = 31
	SocketEffectCriticalRate    SocketEffectType = 32
	SocketEffectVitality        SocketEffectType = 36
	SocketEffectDoubleDamage    SocketEffectType = 44
	SocketEffectIgnoreDefense   SocketEffectType = 45
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
