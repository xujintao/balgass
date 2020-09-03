package item

import (
	"math/rand"
	"path"

	"github.com/xujintao/balgass/cmd/server_game/conf"
	"github.com/xujintao/balgass/cmd/server_game/game/random"
	"github.com/xujintao/balgass/cmd/server_game/lang"
)

type JewelHarmonyItemEffect struct {
}

type harmonyItemKind int

const (
	harmonyItemNull harmonyItemKind = iota - 1
	harmonyItemWeapon
	harmonyItemStaff
	harmonyItemDefense
)

type harmonyEffectKind int

const harmonyEffectNull harmonyEffectKind = -1
const (
	harmonyEffectWeaponIncDamageMin harmonyEffectKind = iota
	harmonyEffectWeaponIncDamageMax
	harmonyEffectWeaponDecReqStrength
	harmonyEffectWeaponDecReqAgility
	harmonyEffectWeaponIncDamage
	harmonyEffectWeaponIncCriticalDamage // 加重
	harmonyEffectWeaponIncSkillDamage
	harmonyEffectWeaponIncAttackRate
	harmonyEffectWeaponDecSDRate
	harmonyEffectWeaponIgnoreSDRate
)
const (
	harmonyEffectStaffIncMagicDamage harmonyEffectKind = iota
	harmonyEffectStaffDecReqStrength
	harmonyEffectStaffDecReqAgility
	harmonyEffectStaffIncSkillDamage
	harmonyEffectStaffIncCriticalDamage
	harmonyEffectStaffDecSDRate
	harmonyEffectStaffIncAttackRate
	harmonyEffectStaffIgnoreSDRate
)
const (
	harmonyEffectDefenseIncDefense harmonyEffectKind = iota
	harmonyEffectDefenseIncMaxAG
	harmonyEffectDefenseIncMaxHP
	harmonyEffectDefenseIncAutoRecoveryHP
	harmonyEffectDefenseIncAutoRecoveryMP
	harmonyEffectDefenseIncDefenseRate
	harmonyEffectDefenseDecDamage
	harmonyEffectDefenseIncSDRate
)

type HarmonyEffect struct {
	HarmonyEffectIncDamageMin      int
	HarmonyEffectIncDamageMax      int
	HarmonyEffectDecReqStrength    int
	HarmonyEffectDecReqAgility     int
	HarmonyEffectIncDamage         int
	HarmonyEffectIncCriticalDamage int
	HarmonyEffectIncSkillDamage    int
	HarmonyEffectIncAttackRate     int
	HarmonyEffectDecSDRate         int
	HarmonyEffectIgnoreSDRate      int

	HarmonyEffectIncMagicDamage int

	HarmonyEffectIncDefense        int
	HarmonyEffectIncMaxAG          int
	HarmonyEffectIncMaxHP          int
	HarmonyEffectIncAutoRecoveryHP int
	HarmonyEffectIncAutoRecoveryMP int
	HarmonyEffectIncDefenseRate    int
	HarmonyEffectDecDamage         int
	HarmonyEffectIncSDRate         int
}

type harmonyPair struct {
	value    int
	reqMoney int
}

type harmony struct {
	weight   int
	reqLevel int
	pairs    [16]harmonyPair
}

type harmonyManager struct {
	// harmonyOriginCode     int // 再生原石
	// harmonyJewelCode      int // 再生宝石
	// refineStoneLowerCode  int // 低级进化宝石
	// refineStoneHigherCode int // 高级进化宝石
	Config struct {
		EnableHarmonyJewelMix bool `ini:"HarmonyJewelMix"`       // 合成再生宝石
		EnableRefineStoneMix  bool `ini:"RefiningStoneMix"`      // 合成进化宝石
		EnableStrengthenItem  bool `ini:"StrengthenItem"`        // 强化道具
		EnableRestoreItem     bool `ini:"RestoreStrengthenItem"` // 还原道具
		EnableRefineItem      bool `ini:"SmeltItem"`             // 进化道具

		HarmonyJewelMixSuccessRate      int `ini:"JewelOfHarmonyMixSuccessRate"` // 合成再生宝石
		HarmonyJewelMixReqMoney         int `ini:"JewelOfHarmonyMixReqMoney"`
		RefineStoneLowerMixSuccessRate  int `ini:"LowerRefiningStoneMixSuccessRate"` // 合成进化宝石
		RefineStoneHigherMixSuccessRate int `ini:"HigherRefiningStoneMixSuccessRate"`
		RefineStoneMinxReqMoney         int `ini:"RefiningStoneMixReqMoney"`
		StrengthenItemSuccessRate       int `ini:"StrengthenItemSuccessRate"`      // 强化道具
		RefineItemLowerSuccessRate      int `ini:"SmeltingItemSuccessRate_Normal"` // 进化道具
		RefineItemHigherSuccessRate     int `ini:"SmeltingItemSuccessRate_Enhanced"`
	}
	harmonys [][]harmony
	items    []map[int]int
}

func (o *harmonyManager) Out() {}

func (o *harmonyManager) getHarmonyItemKind(item *Item) harmonyItemKind {
	switch item.Section {
	case 0, 1, 2, 3:
		return harmonyItemWeapon
	case 4:
		if item.Index == 7 || item.Index == 15 {
			return harmonyItemNull
		}
		return harmonyItemWeapon
	case 5:
		return harmonyItemStaff
	case 6, 7, 8, 9, 10, 11:
		return harmonyItemDefense
	}

	return harmonyItemNull
}

func (o *harmonyManager) addRandEffect(item *Item, itemKind harmonyItemKind) {
	rm := random.RandManager{}
	for i, h := range o.harmonys[itemKind] {
		if item.Level < h.reqLevel {
			continue
		}
		rm.Put(i, h.weight)
	}
	effect, _ := rm.GetWithWeight().(int)
	item.HarmonyEffect = harmonyEffectKind(effect)
	item.HarmonyLevel = o.harmonys[itemKind][effect].reqLevel
}

func (o *harmonyManager) StrengthenItem(item *Item) (bool, error) {
	if item.HarmonyEffect > 0 {
		return false, lang.MsgRestrengthen // already strengthened
	}
	if item.Set > 0 && conf.Common.General.EnableUseSetHarmonyItem == false {
		return false, lang.MsgStrengthenSet // set can not be strengthened
	}
	// socket item

	itemKind := o.getHarmonyItemKind(item)
	if itemKind == harmonyItemNull {
		return false, lang.MsgStrengthenFailed // invalid harmony item
	}
	if rand.Intn(100)+1 > o.Config.StrengthenItemSuccessRate {
		return true, lang.MsgStrengthenFailed // strengthen failed
	}
	o.addRandEffect(item, itemKind)
	// can not find persistence
	return true, lang.MsgStrengthenSuccess
}

var HarmonyManager harmonyManager

func init() {
	conf.INISection(path.Join(conf.PathCommon, "IGC_HarmonySystem.ini"), "HarmonySystem", &HarmonyManager.Config)
	conf.INISection(path.Join(conf.PathCommon, "IGC_HarmonySystem.ini"), "HarmonyMix", &HarmonyManager.Config)
	// conf.INI(path.Join(conf.PathCommon, "IGC_HarmonySystem.ini"), &HarmonyManager.config)
	type HarmonySystem struct {
		Type []struct {
			// ID          int    `xml:"ID,attr"`
			Description string `xml:"Description,attr"`
			Option      []struct {
				// Index        int    `xml:"Index,attr"`
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
	var harmonySystem HarmonySystem
	conf.XML(path.Join(conf.PathCommon, "Items/IGC_HarmonyItem_Option.xml"), &harmonySystem)
	// convert
	for _, harmonyType := range harmonySystem.Type {
		var harmonys []harmony
		for _, option := range harmonyType.Option {
			var h harmony
			h.weight = option.RandomWeight
			h.reqLevel = option.ReqLevel
			for i, effect := range option.Effect {
				h.pairs[i] = harmonyPair{effect.Value, effect.ReqMoney}
			}
			harmonys = append(harmonys, h)
		}
		HarmonyManager.harmonys = append(HarmonyManager.harmonys, harmonys)
	}
	type SmeltItem struct {
		Section []struct {
			ID   int `xml:"ID,attr"`
			Item []struct {
				Index    int `xml:"Index,attr"`
				ReqLevel int `xml:"ReqLevel,attr"`
			} `xml:"Item"`
		} `xml:"Section"`
	}
	var smeltItem SmeltItem
	conf.XML(path.Join(conf.PathCommon, "Items/IGC_HarmonyItem_Smelt.xml"), &smeltItem)
	// convert
	HarmonyManager.items = make([]map[int]int, len(smeltItem.Section))
	for _, section := range smeltItem.Section {
		items := make(map[int]int)
		for _, item := range section.Item {
			items[item.Index] = item.ReqLevel
		}
		HarmonyManager.items[section.ID] = items
	}
}
