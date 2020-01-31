package item

import (
	"path"

	"github.com/xujintao/balgass/cmd/server_game/conf"
)

type JewelHarmonyItemEffect struct {
}

type harmonyEffectKind int

const (
	harmonyEffectWeaponIncDamageMin harmonyEffectKind = iota + 1
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
	harmonyEffectStaffIncMagicDamage harmonyEffectKind = iota + 1
	harmonyEffectStaffDecReqStrength
	harmonyEffectStaffDecReqAgility
	harmonyEffectStaffIncSkillDamage
	harmonyEffectStaffIncCriticalDamage
	harmonyEffectStaffDecSDRate
	harmonyEffectStaffIncAttackRate
	harmonyEffectStaffIgnoreSDRate
)
const (
	harmonyEffectDefenseIncDefense harmonyEffectKind = iota + 1
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
	reqLevel int
	pairs    [16]harmonyPair
}

type harmonyManager struct {
	harmonyOriginCode     int
	harmonyJewelCode      int
	refineStoneLowerCode  int
	refineStoneHigherCode int
	config                struct {
		EnableHarmonyJewelMix               bool `ini:"HarmonyJewelMix"`
		EnableRefineStoneMix                bool `ini:"RefiningStoneMix"`
		EnableStrengthenItem                bool `ini:"StrengthenItem"`
		EnableRestoreItem                   bool `ini:"RestoreStrengthenItem"`
		EnableExtractItem                   bool `ini:"SmeltItem"`
		HarmonyJewelMixSuccessRate          int  `ini:"JewelOfHarmonyMixSuccessRate"`
		HarmonyJewelMixReqMoney             int  `ini:"JewelOfHarmonyMixReqMoney"`
		RefineStoneLowerMixSuccessRate      int  `ini:"LowerRefiningStoneMixSuccessRate"`
		RefineStoneHigherMixSuccessRate     int  `ini:"HigherRefiningStoneMixSuccessRate"`
		RefineStoneMinxReqMoney             int  `ini:"RefiningStoneMixReqMoney"`
		StrengthenItemSuccessRate           int  `ini:"StrengthenItemSuccessRate"`
		ExtractRefineStoneLowerSuccessRate  int  `ini:"SmeltingItemSuccessRate_Normal"`
		ExtractRefineStoneHigherSuccessRate int  `ini:"SmeltingItemSuccessRate_Enhanced"`
	}
	harmonys [][]harmony
}

func (o *harmonyManager) Out() {

}

var HarmonyManager harmonyManager

func init() {
	conf.INISection(path.Join(conf.PathCommon, "IGC_HarmonySystem.ini"), "HarmonySystem", &HarmonyManager.config)
	conf.INISection(path.Join(conf.PathCommon, "IGC_HarmonySystem.ini"), "HarmonyMix", &HarmonyManager.config)
	// conf.INI(path.Join(conf.PathCommon, "IGC_HarmonySystem.ini"), &HarmonyManager.config)
	type HarmonySystem struct {
		Type []struct {
			ID          int    `xml:"ID,attr"`
			Description string `xml:"Description,attr"`
			Option      []struct {
				Index        int    `xml:"Index,attr"`
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
			h.reqLevel = option.ReqLevel
			for i, effect := range option.Effect {
				h.pairs[i] = harmonyPair{effect.Value, effect.ReqMoney}
			}
			harmonys = append(harmonys, h)
		}
		HarmonyManager.harmonys = append(HarmonyManager.harmonys, harmonys)
	}
}
