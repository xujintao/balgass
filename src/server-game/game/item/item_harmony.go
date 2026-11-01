package item

import (
	"math/rand"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/lang"
	"github.com/xujintao/balgass/src/server-game/game/random"
)

type harmonyItemKind int

const (
	harmonyItemKindNull harmonyItemKind = iota
	harmonyItemKindWeapon
	harmonyItemKindStaff
	harmonyItemKindDefense
)

type harmonyOption int

const (
	harmonyOptionWeaponIncDamageMin harmonyOption = iota + 1
	harmonyOptionWeaponIncDamageMax
	harmonyOptionWeaponDecReqStrength
	harmonyOptionWeaponDecReqAgility
	harmonyOptionWeaponIncDamage
	harmonyOptionWeaponIncCriticalDamage
	harmonyOptionWeaponIncSkillDamage
	harmonyOptionWeaponIncAttackRate
	harmonyOptionWeaponDecSDRate
	harmonyOptionWeaponIgnoreSDRate
)

const (
	harmonyOptionStaffIncMagicDamage harmonyOption = iota + 1
	harmonyOptionStaffDecReqStrength
	harmonyOptionStaffDecReqAgility
	harmonyOptionStaffIncSkillDamage
	harmonyOptionStaffIncCriticalDamage
	harmonyOptionStaffDecSDRate
	harmonyOptionStaffIncAttackRate
	harmonyOptionStaffIgnoreSDRate
)

const (
	harmonyOptionDefenseIncDefense harmonyOption = iota + 1
	harmonyOptionDefenseIncMaxAG
	harmonyOptionDefenseIncMaxHP
	harmonyOptionDefenseIncAutoRecoveryHP
	harmonyOptionDefenseIncAutoRecoveryMP
	harmonyOptionDefenseIncDefenseRate
	harmonyOptionDefenseDecDamage
	harmonyOptionDefenseIncSDRate
)

var HarmonyManager harmonyManager

type ItemHarmonyOption struct {
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
type HarmonyConfig struct {
	HarmonySystem struct {
		EnableHarmonyJewelMix  bool `ini:"HarmonyJewelMix"`       // 使能再生宝石合成
		EnableRefiningStoneMix bool `ini:"RefiningStoneMix"`      // 使能进化宝石合成
		EnableStrengthenItem   bool `ini:"StrengthenItem"`        // 使能强化道具
		EnableSmeltItem        bool `ini:"SmeltItem"`             // 使能进化道具
		EnableRestoreItem      bool `ini:"RestoreStrengthenItem"` // 使能恢复道具
	} `ini:"HarmonySystem"`
	HarmonyMix struct {
		HarmonyJewelMixSuccessRate        int `ini:"JewelOfHarmonyMixSuccessRate"`      // 再生宝石合成成功率
		HarmonyJewelMixReqMoney           int `ini:"JewelOfHarmonyMixReqMoney"`         // 再生宝石合成费用
		RefiningStoneLowerMixSuccessRate  int `ini:"LowerRefiningStoneMixSuccessRate"`  // 初级进化宝石合成成功率
		RefiningStoneHigherMixSuccessRate int `ini:"HigherRefiningStoneMixSuccessRate"` // 高级进化宝石合成成功率
		RefiningStoneMinxReqMoney         int `ini:"RefiningStoneMixReqMoney"`          // 进化宝石合成费用
		StrengthenItemSuccessRate         int `ini:"StrengthenItemSuccessRate"`         // 强化道具成功率
		SmeltItemLowerSuccessRate         int `ini:"SmeltingItemSuccessRate_Normal"`    // 进化道具初级成功率
		SmeltItemHigherSuccessRate        int `ini:"SmeltingItemSuccessRate_Enhanced"`  // 进化道具高级成功率
	} `ini:"HarmonyMix"`
}
type HarmonyLevel struct {
	value    int
	reqMoney int
}

type harmony struct {
	weight   int
	reqLevel int
	levels   [16]HarmonyLevel
}

type harmonyManager struct {
	cfg      HarmonyConfig
	harmonys [][]harmony
	items    []map[int]int
}

func (m *harmonyManager) init() {
	conf.INI(conf.PathCommon, "IGC_HarmonySystem.ini", &m.cfg)
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
				h.levels[i] = HarmonyLevel{value: effect.Value, reqMoney: effect.ReqMoney}
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

func (m *harmonyManager) ApplyHarmonyItemOption(inv *Inventory, apply func(ItemHarmonyOption, int)) {
	if inv == nil || apply == nil {
		return
	}
	for position, it := range inv.WearingItems() {
		if it == nil {
			continue
		}
		var option ItemHarmonyOption
		m.applyItemHarmonyOption(it, &option)
		apply(option, position)
	}
}

func (m *harmonyManager) applyItemHarmonyOption(it *Item, option *ItemHarmonyOption) {
	itemKind, itemOption, ok := m.itemKindOption(it)
	if !ok {
		return
	}
	value := itemOption.levels[it.HarmonyLevel].value
	switch itemKind {
	case harmonyItemKindWeapon:
		switch harmonyOption(it.HarmonyOption) {
		case harmonyOptionWeaponIncDamageMin:
			option.IncDamageMin += value
		case harmonyOptionWeaponIncDamageMax:
			option.IncDamageMax += value
		case harmonyOptionWeaponDecReqStrength:
			option.DecReqStrength += value
		case harmonyOptionWeaponDecReqAgility:
			option.DecReqAgility += value
		case harmonyOptionWeaponIncDamage:
			option.IncDamage += value
		case harmonyOptionWeaponIncCriticalDamage:
			option.IncCriticalDamage += value
		case harmonyOptionWeaponIncSkillDamage:
			option.IncSkillDamage += value
		case harmonyOptionWeaponIncAttackRate:
			option.IncAttackRate += value
		case harmonyOptionWeaponDecSDRate:
			option.DecSDRate += value
		case harmonyOptionWeaponIgnoreSDRate:
			option.IgnoreSDRate += value
		}
	case harmonyItemKindStaff:
		switch harmonyOption(it.HarmonyOption) {
		case harmonyOptionStaffIncMagicDamage:
			option.IncMagicDamage += value
		case harmonyOptionStaffDecReqStrength:
			option.DecReqStrength += value
		case harmonyOptionStaffDecReqAgility:
			option.DecReqAgility += value
		case harmonyOptionStaffIncSkillDamage:
			option.IncSkillDamage += value
		case harmonyOptionStaffIncCriticalDamage:
			option.IncCriticalDamage += value
		case harmonyOptionStaffDecSDRate:
			option.DecSDRate += value
		case harmonyOptionStaffIncAttackRate:
			option.IncAttackRate += value
		case harmonyOptionStaffIgnoreSDRate:
			option.IgnoreSDRate += value
		}
	case harmonyItemKindDefense:
		switch harmonyOption(it.HarmonyOption) {
		case harmonyOptionDefenseIncDefense:
			option.IncDefense += value
		case harmonyOptionDefenseIncMaxAG:
			option.IncMaxAG += value
		case harmonyOptionDefenseIncMaxHP:
			option.IncMaxHP += value
		case harmonyOptionDefenseIncAutoRecoveryHP:
			option.IncAutoRecoveryHP += value
		case harmonyOptionDefenseIncAutoRecoveryMP:
			option.IncAutoRecoveryMP += value
		case harmonyOptionDefenseIncDefenseRate:
			option.IncDefenseRate += value
		case harmonyOptionDefenseDecDamage:
			option.DecDamage += value
		case harmonyOptionDefenseIncSDRate:
			option.IncSDRate += value
		}
	}
}

func (m *harmonyManager) itemKindOption(it *Item) (harmonyItemKind, harmony, bool) {
	if it == nil || it.HarmonyOption <= 0 || it.HarmonyLevel < 0 || it.HarmonyLevel >= 16 ||
		it.Level < it.HarmonyLevel {
		return harmonyItemKindNull, harmony{}, false
	}
	itemKind := m.itemKind(it)
	if itemKind == harmonyItemKindNull || int(itemKind) >= len(m.harmonys) {
		return harmonyItemKindNull, harmony{}, false
	}
	options := m.harmonys[itemKind]
	index := it.HarmonyOption
	if index >= len(options) {
		return harmonyItemKindNull, harmony{}, false
	}
	return itemKind, options[index], true
}

func (m *harmonyManager) itemKind(it *Item) harmonyItemKind {
	switch it.Section {
	case 0, 1, 2, 3:
		return harmonyItemKindWeapon
	case 4:
		if it.Index == 7 || it.Index == 15 {
			return harmonyItemKindNull
		}
		return harmonyItemKindWeapon
	case 5:
		return harmonyItemKindStaff
	case 6, 7, 8, 9, 10, 11:
		return harmonyItemKindDefense
	default:
		return harmonyItemKindNull
	}
}

func (m *harmonyManager) itemHarmoneyEffectDecReqStrength(it *Item) int {
	return m.itemEffect(it).DecReqStrength
}

func (m *harmonyManager) itemHarmoneyEffectDecReqDexterity(it *Item) int {
	return m.itemEffect(it).DecReqAgility
}

func (m *harmonyManager) itemEffect(it *Item) ItemHarmonyOption {
	var option ItemHarmonyOption
	m.applyItemHarmonyOption(it, &option)
	return option
}

func (m *harmonyManager) Config() HarmonyConfig {
	return m.cfg
}

func (m *harmonyManager) CanMakeRefiningStone(it *Item) (bool, bool) {
	if it == nil || it.Type == TypeSocket || it.HarmonyOption != 0 ||
		it.IsSet() && !conf.Common.General.EnableUseSetHarmonyItem {
		return false, false
	}
	eligible := false
	switch it.Section {
	case 0, 1, 2, 3, 5, 6, 7, 8, 9, 10, 11:
		eligible = true
	case 4:
		eligible = it.Index != 7 && it.Index != 15
	}
	if !eligible {
		return false, false
	}
	if it.Section >= 0 && it.Section < len(m.items) {
		if reqLevel, ok := m.items[it.Section][it.Index]; ok && it.Level < reqLevel {
			return false, false
		}
	}
	return true, it.IsExcellent() || it.IsSet()
}

func (m *harmonyManager) RestoreItemMoney(it *Item) (int, bool) {
	if it == nil || it.HarmonyOption <= 0 || it.HarmonyLevel < 0 || it.HarmonyLevel >= 16 {
		return 0, false
	}
	itemKind := m.itemKind(it)
	if itemKind == harmonyItemKindNull {
		return 0, false
	}
	options := m.harmonys[itemKind]
	if it.HarmonyOption >= len(options) {
		return 0, false
	}
	return options[it.HarmonyOption].levels[it.HarmonyLevel].reqMoney, true
}

func (m *harmonyManager) StrengthenItem(it *Item) (bool, error) {
	return m.strengthenItem(it, false)
}

func (m *harmonyManager) StrengthenLuckyItem(it *Item) (bool, error) {
	return m.strengthenItem(it, true)
}

func (m *harmonyManager) strengthenItem(it *Item, lucky bool) (bool, error) {
	if !m.cfg.HarmonySystem.EnableStrengthenItem {
		return false, lang.MsgStrengthenFailed
	}
	if it == nil || it.Type == TypeSocket || it.HarmonyOption != 0 {
		if it != nil && it.HarmonyOption != 0 {
			return false, lang.MsgRestrengthen
		}
		return false, lang.MsgStrengthenFailed
	}
	if lucky {
		if it.Type != TypeLucky {
			return false, lang.MsgStrengthenFailed
		}
	}
	if !lucky && it.IsSet() && !conf.Common.General.EnableUseSetHarmonyItem {
		return false, lang.MsgStrengthenSet // set can not be strengthened
	}
	itemKind := m.itemKind(it)
	if itemKind == harmonyItemKindNull {
		return false, lang.MsgStrengthenFailed // invalid harmony item
	}
	option, ok := m.randomOption(it, itemKind)
	if !ok {
		return false, lang.MsgStrengthenFailed
	}
	if rand.Intn(100) >= m.cfg.HarmonyMix.StrengthenItemSuccessRate {
		return true, lang.MsgStrengthenFailed // strengthen failed
	}
	it.HarmonyOption = option
	it.HarmonyLevel = m.harmonys[itemKind][option].reqLevel
	return true, lang.MsgStrengthenSuccess
}

func (m *harmonyManager) randomOption(it *Item, itemKind harmonyItemKind) (int, bool) {
	rm := random.RandManager{}
	if int(itemKind) >= len(m.harmonys) {
		return 0, false
	}
	options := m.harmonys[itemKind]
	for i := 1; i < len(options); i++ {
		option := options[i]
		if it.Level < option.reqLevel {
			continue
		}
		if (itemKind == harmonyItemKindWeapon && (i == int(harmonyOptionWeaponDecReqStrength) || i == int(harmonyOptionWeaponDecReqAgility))) ||
			(itemKind == harmonyItemKindStaff && (i == int(harmonyOptionStaffDecReqStrength) || i == int(harmonyOptionStaffDecReqAgility))) {
			if it.ReqStrength == 0 && (i == int(harmonyOptionWeaponDecReqStrength) || i == int(harmonyOptionStaffDecReqStrength)) ||
				it.ReqDexterity == 0 && (i == int(harmonyOptionWeaponDecReqAgility) || i == int(harmonyOptionStaffDecReqAgility)) {
				continue
			}
		}
		rm.Put(i, option.weight)
	}
	option, ok := rm.GetWithWeight().(int)
	return option, ok && option > 0
}

func (m *harmonyManager) SmeltItem(it *Item, high bool) (bool, bool) {
	if !m.cfg.HarmonySystem.EnableSmeltItem || it == nil || it.HarmonyOption <= 0 || it.HarmonyLevel >= it.Level || it.HarmonyLevel == 13 {
		return false, false
	}
	itemKind := m.itemKind(it)
	if itemKind == harmonyItemKindNull {
		return false, false
	}
	options := m.harmonys[itemKind]
	if it.HarmonyOption >= len(options) || it.HarmonyLevel < 0 || it.HarmonyLevel >= len(options[it.HarmonyOption].levels) {
		return false, false
	}
	option := options[it.HarmonyOption]
	if itemKind == harmonyItemKindWeapon && it.HarmonyOption == int(harmonyOptionWeaponIncDamageMin) {
		next := it.HarmonyLevel + 1
		if it.AttackMin+option.levels[next].value > it.AttackMax-1 {
			return false, false
		}
	}
	rate := m.cfg.HarmonyMix.SmeltItemLowerSuccessRate
	if high {
		rate = m.cfg.HarmonyMix.SmeltItemHigherSuccessRate
	}
	if rand.Intn(100) >= rate {
		it.HarmonyLevel = option.reqLevel
		return true, false
	}
	it.HarmonyLevel++
	return true, true
}
