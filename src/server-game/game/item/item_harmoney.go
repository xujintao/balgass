package item

import (
	"math/rand"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/lang"
	"github.com/xujintao/balgass/src/server-game/game/random"
)

type harmonyItemKind int

const (
	harmonyItemNull harmonyItemKind = iota
	harmonyItemWeapon
	harmonyItemStaff
	harmonyItemDefense
)

type HarmonyEffectType int

const (
	HarmonyEffectWeaponIncDamageMin HarmonyEffectType = iota + 1
	HarmonyEffectWeaponIncDamageMax
	HarmonyEffectWeaponDecReqStrength
	HarmonyEffectWeaponDecReqAgility
	HarmonyEffectWeaponIncDamage
	HarmonyEffectWeaponIncCriticalDamage
	HarmonyEffectWeaponIncSkillDamage
	HarmonyEffectWeaponIncAttackRate
	HarmonyEffectWeaponDecSDRate
	HarmonyEffectWeaponIgnoreSDRate
)

const (
	HarmonyEffectStaffIncMagicDamage HarmonyEffectType = iota + 1
	HarmonyEffectStaffDecReqStrength
	HarmonyEffectStaffDecReqAgility
	HarmonyEffectStaffIncSkillDamage
	HarmonyEffectStaffIncCriticalDamage
	HarmonyEffectStaffDecSDRate
	HarmonyEffectStaffIncAttackRate
	HarmonyEffectStaffIgnoreSDRate
)

const (
	HarmonyEffectDefenseIncDefense HarmonyEffectType = iota + 1
	HarmonyEffectDefenseIncMaxAG
	HarmonyEffectDefenseIncMaxHP
	HarmonyEffectDefenseIncAutoRecoveryHP
	HarmonyEffectDefenseIncAutoRecoveryMP
	HarmonyEffectDefenseIncDefenseRate
	HarmonyEffectDefenseDecDamage
	HarmonyEffectDefenseIncSDRate
)

var HarmonyManager harmonyManager

type ItemHarmonyEffect struct {
	IncDamageMin      int
	IncDamageMax      int
	IncMagicDamage    int
	DecReqStrength    int
	DecReqAgility     int
	IncDamage         int
	IncCriticalDamage int
	IncSkillDamage    int
	IncAttackRate     int
	DecSDRate         int
	IgnoreSDRate      int
	IncDefense        int
	IncMaxAG          int
	IncMaxHP          int
	IncAutoRecoveryHP int
	IncAutoRecoveryMP int
	IncDefenseRate    int
	DecDamage         int
	IncSDRate         int
}

type HarmonyEffectLevel struct {
	value    int
	reqMoney int
}

type harmony struct {
	weight   int
	reqLevel int
	levels   [16]HarmonyEffectLevel
}

type harmonySettings struct {
	strengthenItemSuccessRate int
}

type harmonyManager struct {
	settings harmonySettings
	harmonys [][]harmony
	items    []map[int]int
}

func (m *harmonyManager) init() {
	type harmonyConfig struct {
		HarmonySystem struct {
			EnableHarmonyJewelMix bool `ini:"HarmonyJewelMix"`
			EnableRefineStoneMix  bool `ini:"RefiningStoneMix"`
			EnableStrengthenItem  bool `ini:"StrengthenItem"`
			EnableRestoreItem     bool `ini:"RestoreStrengthenItem"`
			EnableRefineItem      bool `ini:"SmeltItem"`
		} `ini:"HarmonySystem"`
		HarmonyMix struct {
			HarmonyJewelMixSuccessRate      int `ini:"JewelOfHarmonyMixSuccessRate"`
			HarmonyJewelMixReqMoney         int `ini:"JewelOfHarmonyMixReqMoney"`
			RefineStoneLowerMixSuccessRate  int `ini:"LowerRefiningStoneMixSuccessRate"`
			RefineStoneHigherMixSuccessRate int `ini:"HigherRefiningStoneMixSuccessRate"`
			RefineStoneMinxReqMoney         int `ini:"RefiningStoneMixReqMoney"`
			StrengthenItemSuccessRate       int `ini:"StrengthenItemSuccessRate"`
			RefineItemLowerSuccessRate      int `ini:"SmeltingItemSuccessRate_Normal"`
			RefineItemHigherSuccessRate     int `ini:"SmeltingItemSuccessRate_Enhanced"`
		} `ini:"HarmonyMix"`
	}
	var cfg harmonyConfig
	conf.INI(conf.PathCommon, "IGC_HarmonySystem.ini", &cfg)
	m.settings.strengthenItemSuccessRate = cfg.HarmonyMix.StrengthenItemSuccessRate
	type harmonyOptionConfig struct {
		Type []struct {
			Description string `xml:"Description,attr"`
			Option      []struct {
				Name         string `xml:"Name,attr"`
				RandomWeight int    `xml:"RandomWeight,attr"`
				ReqLevel     int    `xml:"ReqLevel,attr"`
				Effect       []struct {
					Level    int `xml:"level,attr"`
					Value    int `xml:"Value,attr"`
					ReqMoney int `xml:"ReqMoney,attr"`
				} `xml:"Effect"`
			} `xml:"Option"`
		} `xml:"Type"`
	}
	var optionCfg harmonyOptionConfig
	conf.XML(conf.PathCommon, "Items/IGC_HarmonyItem_Option.xml", &optionCfg)
	m.harmonys = make([][]harmony, 1)
	for _, harmonyType := range optionCfg.Type {
		options := make([]harmony, 1)
		for _, option := range harmonyType.Option {
			var h harmony
			h.weight = option.RandomWeight
			h.reqLevel = option.ReqLevel
			for i, effect := range option.Effect {
				h.levels[i] = HarmonyEffectLevel{value: effect.Value, reqMoney: effect.ReqMoney}
			}
			options = append(options, h)
		}
		m.harmonys = append(m.harmonys, options)
	}
	type smeltConfig struct {
		Section []struct {
			ID   int `xml:"ID,attr"`
			Item []struct {
				Index    int `xml:"Index,attr"`
				ReqLevel int `xml:"ReqLevel,attr"`
			} `xml:"Item"`
		} `xml:"Section"`
	}
	var smeltCfg smeltConfig
	conf.XML(conf.PathCommon, "Items/IGC_HarmonyItem_Smelt.xml", &smeltCfg)
	m.items = make([]map[int]int, len(smeltCfg.Section))
	for _, section := range smeltCfg.Section {
		items := make(map[int]int)
		for _, it := range section.Item {
			items[it.Index] = it.ReqLevel
		}
		m.items[section.ID] = items
	}
}

func (m *harmonyManager) itemHarmoneyEffectDecReqStrength(it *Item) int {
	return m.itemEffect(it).DecReqStrength
}

func (m *harmonyManager) itemHarmoneyEffectDecReqDexterity(it *Item) int {
	return m.itemEffect(it).DecReqAgility
}

func (m *harmonyManager) itemEffect(it *Item) ItemHarmonyEffect {
	var effect ItemHarmonyEffect
	m.applyItemHarmonyEffect(it, &effect)
	return effect
}

func (m *harmonyManager) ApplyHarmonyItemEffect(inv *Inventory, apply func(ItemHarmonyEffect, int)) {
	if inv == nil || apply == nil {
		return
	}
	for position, it := range inv.WearingItems() {
		if it == nil {
			continue
		}
		var effect ItemHarmonyEffect
		m.applyItemHarmonyEffect(it, &effect)
		apply(effect, position)
	}
}

func (m *harmonyManager) applyItemHarmonyEffect(it *Item, effect *ItemHarmonyEffect) {
	itemKind, option, ok := m.itemOption(it)
	if !ok {
		return
	}
	value := option.levels[it.Level].value
	switch itemKind {
	case harmonyItemWeapon:
		switch HarmonyEffectType(it.Harmony) {
		case HarmonyEffectWeaponIncDamageMin:
			effect.IncDamageMin += value
		case HarmonyEffectWeaponIncDamageMax:
			effect.IncDamageMax += value
		case HarmonyEffectWeaponDecReqStrength:
			effect.DecReqStrength += value
		case HarmonyEffectWeaponDecReqAgility:
			effect.DecReqAgility += value
		case HarmonyEffectWeaponIncDamage:
			effect.IncDamage += value
		case HarmonyEffectWeaponIncCriticalDamage:
			effect.IncCriticalDamage += value
		case HarmonyEffectWeaponIncSkillDamage:
			effect.IncSkillDamage += value
		case HarmonyEffectWeaponIncAttackRate:
			effect.IncAttackRate += value
		case HarmonyEffectWeaponDecSDRate:
			effect.DecSDRate += value
		case HarmonyEffectWeaponIgnoreSDRate:
			effect.IgnoreSDRate += value
		}
	case harmonyItemStaff:
		switch HarmonyEffectType(it.Harmony) {
		case HarmonyEffectStaffIncMagicDamage:
			effect.IncMagicDamage += value
		case HarmonyEffectStaffDecReqStrength:
			effect.DecReqStrength += value
		case HarmonyEffectStaffDecReqAgility:
			effect.DecReqAgility += value
		case HarmonyEffectStaffIncSkillDamage:
			effect.IncSkillDamage += value
		case HarmonyEffectStaffIncCriticalDamage:
			effect.IncCriticalDamage += value
		case HarmonyEffectStaffDecSDRate:
			effect.DecSDRate += value
		case HarmonyEffectStaffIncAttackRate:
			effect.IncAttackRate += value
		case HarmonyEffectStaffIgnoreSDRate:
			effect.IgnoreSDRate += value
		}
	case harmonyItemDefense:
		switch HarmonyEffectType(it.Harmony) {
		case HarmonyEffectDefenseIncDefense:
			effect.IncDefense += value
		case HarmonyEffectDefenseIncMaxAG:
			effect.IncMaxAG += value
		case HarmonyEffectDefenseIncMaxHP:
			effect.IncMaxHP += value
		case HarmonyEffectDefenseIncAutoRecoveryHP:
			effect.IncAutoRecoveryHP += value
		case HarmonyEffectDefenseIncAutoRecoveryMP:
			effect.IncAutoRecoveryMP += value
		case HarmonyEffectDefenseIncDefenseRate:
			effect.IncDefenseRate += value
		case HarmonyEffectDefenseDecDamage:
			effect.DecDamage += value
		case HarmonyEffectDefenseIncSDRate:
			effect.IncSDRate += value
		}
	}
}

func (m *harmonyManager) itemOption(it *Item) (harmonyItemKind, harmony, bool) {
	if it == nil || it.Harmony <= 0 || it.Level < 0 || it.Level >= 16 {
		return harmonyItemNull, harmony{}, false
	}
	itemKind := m.itemKind(it)
	if itemKind == harmonyItemNull {
		return harmonyItemNull, harmony{}, false
	}
	options := m.harmonys[itemKind]
	index := it.Harmony
	if index >= len(options) {
		return harmonyItemNull, harmony{}, false
	}
	return itemKind, options[index], true
}

func (m *harmonyManager) itemKind(it *Item) harmonyItemKind {
	switch it.Section {
	case 0, 1, 2, 3:
		return harmonyItemWeapon
	case 4:
		if it.Index == 7 || it.Index == 15 {
			return harmonyItemNull
		}
		return harmonyItemWeapon
	case 5:
		return harmonyItemStaff
	case 6, 7, 8, 9, 10, 11:
		return harmonyItemDefense
	default:
		return harmonyItemNull
	}
}

func (m *harmonyManager) StrengthenItem(it *Item) (bool, error) {
	if it.Harmony != 0 {
		return false, lang.MsgRestrengthen // already strengthened
	}
	if it.Set > 0 && conf.Common.General.EnableUseSetHarmonyItem == false {
		return false, lang.MsgStrengthenSet // set can not be strengthened
	}
	itemKind := m.itemKind(it)
	if itemKind == harmonyItemNull {
		return false, lang.MsgStrengthenFailed // invalid harmony item
	}
	if rand.Intn(100)+1 > m.settings.strengthenItemSuccessRate {
		return true, lang.MsgStrengthenFailed // strengthen failed
	}
	m.addRandEffect(it, itemKind)
	// can not find persistence
	return true, lang.MsgStrengthenSuccess
}

func (m *harmonyManager) addRandEffect(it *Item, itemKind harmonyItemKind) {
	rm := random.RandManager{}
	options := m.harmonys[itemKind]
	for i := 1; i < len(options); i++ {
		option := options[i]
		if it.Level < option.reqLevel {
			continue
		}
		rm.Put(i, option.weight)
	}
	i, _ := rm.GetWithWeight().(int)
	it.Harmony = i
}
