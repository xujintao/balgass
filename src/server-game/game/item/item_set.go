package item

import (
	"github.com/xujintao/balgass/src/server-game/conf"
)

func init() {
	SetManager.init()
}

type setItem struct {
	tiers    [4]int
	mixLevel [2]int
}

type SetEffectType int

const (
	SetEffectIncStrength SetEffectType = iota
	SetEffectIncAgility
	SetEffectIncEnergy
	SetEffectIncVitality
	SetEffectIncLeadership
	SetEffectIncAttackMin // 最小攻击力增加
	SetEffectIncAttackMax // 最大攻击力增加
	SetEffectIncMagicAttack
	SetEffectIncDamage // 伤害增加
	SetEffectIncAttackRate
	SetEffectIncDefense // 防御力增加
	SetEffectIncMaxHP
	SetEffectIncMaxMP
	SetEffectIncMaxAG
	SetEffectIncAG
	SetEffectIncCritiDamageRate // 幸运一击伤害几率增加
	SetEffectIncCritiDamage     // 幸运一击伤害增加
	SetEffectIncExcelDamageRate // 卓越一击伤害几率增加
	SetEffectIncExcelDamage     // 卓越一击伤害增加
	SetEffectIncSkillAttack     // 技能攻击力增加
	SetEffectDoubleDamage       // 双倍伤害几率
	SetEffectIgnoreDefense      // 无视目标防御
	SetEffectIncShieldDefense
	SetEffectIncTwoHandSwordDamage
)

type SetEffect struct {
	Index SetEffectType
	Value int
}

type set struct {
	name        string
	count       int
	effects     [6]SetEffect
	effectsFull [5]SetEffect
	// class       [8]bool
}

var SetManager setManager

type setManager struct {
	items []map[int]*setItem
	sets  []*set
}

func (m *setManager) GetSetIndex(section, index, tierIndex int) int {
	if _, ok := m.items[section][index]; !ok {
		return 0
	}
	// 0x01=tier0=TierI
	// 0x02=tier1=TierII
	// 0x10=tier2=TierIII
	// 0x20=tier3=TierIV
	return m.items[section][index].tiers[tierIndex-1]
}

func (m *setManager) GetSet(setIndex int, effectIndex int) *SetEffect {
	set := m.sets[setIndex]
	if set == nil {
		return nil
	}
	return &set.effects[effectIndex]
}

func (m *setManager) GetSetEffectCount(setIndex int) int {
	set := m.sets[setIndex]
	if set == nil {
		return 0
	}
	return set.count
}

func (m *setManager) GetSetFull(setIndex int) []SetEffect {
	set := m.sets[setIndex]
	if set == nil {
		return nil
	}
	return set.effectsFull[:]
}

func (m *setManager) init() {
	type SetItemXml struct {
		DropRate struct {
			Sections []struct {
				Index    int `xml:"Index,attr"`
				DropRate int `xml:"DropRate,attr"`
			} `xml:"Section"`
		} `xml:"DropRate"`
		Sections []struct {
			Index int    `xml:"Index,attr"`
			Name  string `xml:"Name,attr"`
			Items []struct {
				Index     int `xml:"Index,attr"`
				TierI     int `xml:"TierI,attr"`
				TierII    int `xml:"TierII,attr"`
				TierIII   int `xml:"TierIII,attr"`
				TierIV    int `xml:"TierIV,attr"`
				MixLevelA int `xml:"MixLevelA,attr"`
				MixLevelB int `xml:"MixLevelB,attr"`
			} `xml:"Item"`
		} `xml:"Section"`
	}
	var setItemXml SetItemXml
	conf.XML(conf.PathCommon, "Items/IGC_ItemSetType.xml", &setItemXml)
	// convert
	m.items = make([]map[int]*setItem, len(setItemXml.Sections))
	for _, section := range setItemXml.Sections {
		items := make(map[int]*setItem)
		for _, item := range section.Items {
			var sItem setItem
			sItem.tiers[0] = item.TierI
			sItem.tiers[1] = item.TierII
			sItem.tiers[2] = item.TierIII
			sItem.tiers[3] = item.TierIV
			sItem.mixLevel[0] = item.MixLevelA
			sItem.mixLevel[1] = item.MixLevelB
			items[item.Index] = &sItem
		}
		m.items[section.Index] = items
	}

	type SetXml struct {
		Sets []struct {
			Index          int           `xml:"Index,attr"`
			Name           string        `xml:"Name,attr"`
			OptIdx11       SetEffectType `xml:"OptIdx1_1,attr"`
			OptVal11       int           `xml:"OptVal1_1,attr"`
			OptIdx21       SetEffectType `xml:"OptIdx2_1,attr"`
			OptVal21       int           `xml:"OptVal2_1,attr"`
			OptIdx12       SetEffectType `xml:"OptIdx1_2,attr"`
			OptVal12       int           `xml:"OptVal1_2,attr"`
			OptIdx22       SetEffectType `xml:"OptIdx2_2,attr"`
			OptVal22       int           `xml:"OptVal2_2,attr"`
			OptIdx13       SetEffectType `xml:"OptIdx1_3,attr"`
			OptVal13       int           `xml:"OptVal1_3,attr"`
			OptIdx23       SetEffectType `xml:"OptIdx2_3,attr"`
			OptVal23       int           `xml:"OptVal2_3,attr"`
			OptIdx14       SetEffectType `xml:"OptIdx1_4,attr"`
			OptVal14       int           `xml:"OptVal1_4,attr"`
			OptIdx24       SetEffectType `xml:"OptIdx2_4,attr"`
			OptVal24       int           `xml:"OptVal2_4,attr"`
			OptIdx15       SetEffectType `xml:"OptIdx1_5,attr"`
			OptVal15       int           `xml:"OptVal1_5,attr"`
			OptIdx25       SetEffectType `xml:"OptIdx2_5,attr"`
			OptVal25       int           `xml:"OptVal2_5,attr"`
			OptIdx16       SetEffectType `xml:"OptIdx1_6,attr"`
			OptVal16       int           `xml:"OptVal1_6,attr"`
			OptIdx26       SetEffectType `xml:"OptIdx2_6,attr"`
			OptVal26       int           `xml:"OptVal2_6,attr"`
			SpecialOptIdx1 SetEffectType `xml:"SpecialOptIdx1,attr"`
			SpecialOptVal1 int           `xml:"SpecialOptVal1,attr"`
			SpecialOptIdx2 SetEffectType `xml:"SpecialOptIdx2,attr"`
			SpecialOptVal2 int           `xml:"SpecialOptVal2,attr"`
			FullOptIdx1    SetEffectType `xml:"FullOptIdx1,attr"`
			FullOptVal1    int           `xml:"FullOptVal1,attr"`
			FullOptIdx2    SetEffectType `xml:"FullOptIdx2,attr"`
			FullOptVal2    int           `xml:"FullOptVal2,attr"`
			FullOptIdx3    SetEffectType `xml:"FullOptIdx3,attr"`
			FullOptVal3    int           `xml:"FullOptVal3,attr"`
			FullOptIdx4    SetEffectType `xml:"FullOptIdx4,attr"`
			FullOptVal4    int           `xml:"FullOptVal4,attr"`
			FullOptIdx5    SetEffectType `xml:"FullOptIdx5,attr"`
			FullOptVal5    int           `xml:"FullOptVal5,attr"`
			DarkWizard     bool          `xml:"DarkWizard,attr"`
			DarkKnight     bool          `xml:"DarkKnight,attr"`
			FairyElf       bool          `xml:"FairyElf,attr"`
			MagicGladiator bool          `xml:"MagicGladiator,attr"`
			DarkLord       bool          `xml:"DarkLord,attr"`
			Summoner       bool          `xml:"Summoner,attr"`
			RageFighter    bool          `xml:"RageFighter,attr"`
		} `xml:"SetItem"`
	}
	var setXml SetXml
	conf.XML(conf.PathCommon, "Items/IGC_ItemSetOption.xml", &setXml)
	// convert
	m.sets = make([]*set, 64)
	for _, s := range setXml.Sets {
		var set set
		set.name = s.Name
		set.effects[0].Index = s.OptIdx11
		set.effects[0].Value = s.OptVal11
		if s.OptIdx11 != -1 {
			set.count++
		}
		set.effects[1].Index = s.OptIdx12
		set.effects[1].Value = s.OptVal12
		if s.OptIdx12 != -1 {
			set.count++
		}
		set.effects[2].Index = s.OptIdx13
		set.effects[2].Value = s.OptVal13
		if s.OptIdx13 != -1 {
			set.count++
		}
		set.effects[3].Index = s.OptIdx14
		set.effects[3].Value = s.OptVal14
		if s.OptIdx14 != -1 {
			set.count++
		}
		set.effects[4].Index = s.OptIdx15
		set.effects[4].Value = s.OptVal15
		if s.OptIdx15 != -1 {
			set.count++
		}
		set.effects[5].Index = s.OptIdx16
		set.effects[5].Value = s.OptVal16
		if s.OptIdx16 != -1 {
			set.count++
		}

		set.effectsFull[0].Index = s.FullOptIdx1
		set.effectsFull[0].Value = s.FullOptVal1
		set.effectsFull[1].Index = s.FullOptIdx2
		set.effectsFull[1].Value = s.FullOptVal2
		set.effectsFull[2].Index = s.FullOptIdx3
		set.effectsFull[2].Value = s.FullOptVal3
		set.effectsFull[3].Index = s.FullOptIdx4
		set.effectsFull[3].Value = s.FullOptVal4
		set.effectsFull[4].Index = s.FullOptIdx5
		set.effectsFull[4].Value = s.FullOptVal5

		m.sets[s.Index] = &set
	}
}
