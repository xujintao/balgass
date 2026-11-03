package item

import (
	"fmt"

	"github.com/xujintao/balgass/src/server-game/conf"
)

type item380OptionType int

const (
	item380OptionNone item380OptionType = iota
	item380OptionIncPVPAttackRate
	item380OptionIncPVPAttack
	item380OptionIncPVPDefenseRate
	item380OptionIncPVPDefense
	item380OptionIncMaxHP
	item380OptionIncMaxSD
	item380OptionAutoRecoverySD
	item380OptionIncRecoverySDRate
)

type Item380Option struct {
	IncPVPAttackRate  int
	IncPVPAttack      int
	IncPVPDefenseRate int
	IncPVPDefense     int
	IncMaxHP          int
	IncMaxSD          int
	AutoRecoverySD    bool
	IncRecoverySDRate int
}

var Item380Manager item380Manager

type item380Manager struct {
	mix   Item380Mix
	items map[int]map[int]*item380
}

type Item380Mix struct {
	Enable             bool `xml:"Enable,attr"`
	HarmonyJewelCount  int  `xml:"JewelOfHarmonyCount,attr"`
	GuardianJewelCount int  `xml:"JewelOfGuardianCount,attr"`
	ReqMoney           int  `xml:"ReqMoney,attr"`
	SuccessRates       [4]int
}

func (mix Item380Mix) SuccessRate(level, charm int) int {
	grade := 0
	switch {
	case level < 7:
		grade = 0
	case level < 10:
		grade = 1
	case level <= 13:
		grade = 2
	default:
		grade = 3
	}
	return mix.SuccessRates[grade] + charm
}

type item380 struct {
	options [2]item380Option
}

type item380Option struct {
	Key   item380OptionType
	Value int
}

func (m *item380Manager) init() {
	const file = "Items/IGC_Item380Option.xml"
	type config struct {
		Mix struct {
			Enable             bool `xml:"Enable,attr"`
			HarmonyJewelCount  int  `xml:"JewelOfHarmonyCount,attr"`
			GuardianJewelCount int  `xml:"JewelOfGuardianCount,attr"`
			ReqMoney           int  `xml:"ReqMoney,attr"`
			SuccessRate        struct {
				Grade1 int `xml:"Grade1,attr"`
				Grade2 int `xml:"Grade2,attr"`
				Grade3 int `xml:"Grade3,attr"`
				Grade4 int `xml:"Grade4,attr"`
			} `xml:"SuccessRate"`
		} `xml:"Mix"`
		ItemOption struct {
			Items []struct {
				Cat     int               `xml:"Cat,attr"`
				Index   int               `xml:"Index,attr"`
				Option1 item380OptionType `xml:"Option1,attr"`
				Value1  int               `xml:"Value1,attr"`
				Option2 item380OptionType `xml:"Option2,attr"`
				Value2  int               `xml:"Value2,attr"`
				Time    int               `xml:"Time,attr"`
			} `xml:"Item"`
		} `xml:"ItemOption"`
	}
	var cfg config
	conf.XML(conf.PathCommon, file, &cfg)
	mix := Item380Mix{
		Enable:             cfg.Mix.Enable,
		HarmonyJewelCount:  cfg.Mix.HarmonyJewelCount,
		GuardianJewelCount: cfg.Mix.GuardianJewelCount,
		ReqMoney:           cfg.Mix.ReqMoney,
	}
	mix.SuccessRates = [4]int{
		cfg.Mix.SuccessRate.Grade1,
		cfg.Mix.SuccessRate.Grade2,
		cfg.Mix.SuccessRate.Grade3,
		cfg.Mix.SuccessRate.Grade4,
	}
	if mix.HarmonyJewelCount < 0 || mix.GuardianJewelCount < 0 || mix.ReqMoney < 0 {
		panic(fmt.Errorf("%s: negative mix requirement", file))
	}
	for _, rate := range mix.SuccessRates {
		if rate < 0 || rate > 100 {
			panic(fmt.Errorf("%s: invalid mix success rate %d", file, rate))
		}
	}
	// convert
	m.items = make(map[int]map[int]*item380)
	for _, item := range cfg.ItemOption.Items {
		if _, err := ItemTable.GetItemBase(item.Cat, item.Index); err != nil {
			panic(fmt.Errorf("%s: section %d index %d: %w", file, item.Cat, item.Index, err))
		}
		if item.Option1 < item380OptionNone || item.Option1 > item380OptionIncRecoverySDRate ||
			item.Option2 < item380OptionNone || item.Option2 > item380OptionIncRecoverySDRate ||
			item.Value1 < 0 || item.Value2 < 0 {
			panic(fmt.Errorf("%s: invalid options for section %d index %d", file, item.Cat, item.Index))
		}
		if m.items[item.Cat] == nil {
			m.items[item.Cat] = make(map[int]*item380)
		}
		if m.items[item.Cat][item.Index] != nil {
			panic(fmt.Errorf("%s: duplicate item section %d index %d", file, item.Cat, item.Index))
		}
		i380 := new(item380)
		i380.options[0] = item380Option{item.Option1, item.Value1}
		i380.options[1] = item380Option{item.Option2, item.Value2}
		m.items[item.Cat][item.Index] = i380
	}
	m.mix = mix
}

func (m *item380Manager) MixConfig() Item380Mix {
	return m.mix
}

func (m *item380Manager) Is380Item(section, index int) bool {
	_, ok := m.items[section][index]
	return ok
}

func (m *item380Manager) Apply380ItemOption(inventory *Inventory, apply func(Item380Option)) {
	if inventory == nil || apply == nil {
		return
	}
	var option Item380Option
	for _, it := range inventory.WearingItems() {
		m.applyItem380Option(it, &option)
	}
	apply(option)
}

func (m *item380Manager) applyItem380Option(it *Item, result *Item380Option) {
	if it == nil || it.Durability <= 0 || !it.Option380 || result == nil {
		return
	}
	entry := m.items[it.Section][it.Index]
	if entry == nil {
		return
	}
	for _, option := range entry.options {
		key, value := option.Key, option.Value
		switch key {
		case item380OptionIncPVPAttackRate:
			result.IncPVPAttackRate += value
		case item380OptionIncPVPAttack:
			result.IncPVPAttack += value
		case item380OptionIncPVPDefenseRate:
			result.IncPVPDefenseRate += value
		case item380OptionIncPVPDefense:
			result.IncPVPDefense += value
		case item380OptionIncMaxHP:
			result.IncMaxHP += value
		case item380OptionIncMaxSD:
			result.IncMaxSD += value
		case item380OptionAutoRecoverySD:
			result.AutoRecoverySD = true
		case item380OptionIncRecoverySDRate:
			result.IncRecoverySDRate += value
		}
	}
}
