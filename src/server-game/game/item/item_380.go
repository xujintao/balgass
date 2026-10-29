package item

import (
	"fmt"

	"github.com/xujintao/balgass/src/server-game/conf"
)

var Item380Manager item380Manager

type item380EffectKind int

const (
	item380EffectNone item380EffectKind = iota
	item380EffectIncPVPAttackRate
	item380EffectIncPVPAttack
	item380EffectIncPVPDefenseRate
	item380EffectIncPVPDefense
	item380EffectIncMaxHP
	item380EffectIncMaxSD
	item380EffectAutoRecoverySD
	item380EffectIncRecoverySDRate
)

type item380EffectPair struct {
	Key   item380EffectKind
	Value int
}

type Item380Effect struct {
	Item380EffectIncPVPAttackRate  int
	Item380EffectIncPVPAttack      int
	Item380EffectIncPVPDefenseRate int
	Item380EffectIncPVPDefense     int
	Item380EffectIncMaxHP          int
	Item380EffectIncMaxSD          int
	Item380EffectAutoRecoverySD    bool
	Item380EffectIncRecoverySDRate int
}

type Item380Mix struct {
	Enable               bool `xml:"Enable,attr"`
	JewelOfHarmonyCount  int  `xml:"JewelOfHarmonyCount,attr"`
	JewelOfGuardianCount int  `xml:"JewelOfGuardianCount,attr"`
	ReqMoney             int  `xml:"ReqMoney,attr"`
	SuccessRates         [4]int
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
	effects [2]item380EffectPair
}

type item380Manager struct {
	mix   Item380Mix
	items map[int]map[int]*item380
}

func (o *item380Manager) init() {
	const file = "Items/IGC_Item380Option.xml"
	type config struct {
		Mix struct {
			Enable               bool `xml:"Enable,attr"`
			JewelOfHarmonyCount  int  `xml:"JewelOfHarmonyCount,attr"`
			JewelOfGuardianCount int  `xml:"JewelOfGuardianCount,attr"`
			ReqMoney             int  `xml:"ReqMoney,attr"`
			SuccessRate          struct {
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
				Option1 item380EffectKind `xml:"Option1,attr"`
				Value1  int               `xml:"Value1,attr"`
				Option2 item380EffectKind `xml:"Option2,attr"`
				Value2  int               `xml:"Value2,attr"`
				Time    int               `xml:"Time,attr"`
			} `xml:"Item"`
		} `xml:"ItemOption"`
	}
	var cfg config
	conf.XML(conf.PathCommon, file, &cfg)
	mix := Item380Mix{
		Enable:               cfg.Mix.Enable,
		JewelOfHarmonyCount:  cfg.Mix.JewelOfHarmonyCount,
		JewelOfGuardianCount: cfg.Mix.JewelOfGuardianCount,
		ReqMoney:             cfg.Mix.ReqMoney,
	}
	mix.SuccessRates = [4]int{
		cfg.Mix.SuccessRate.Grade1,
		cfg.Mix.SuccessRate.Grade2,
		cfg.Mix.SuccessRate.Grade3,
		cfg.Mix.SuccessRate.Grade4,
	}
	if mix.JewelOfHarmonyCount < 0 || mix.JewelOfGuardianCount < 0 || mix.ReqMoney < 0 {
		panic(fmt.Errorf("%s: negative mix requirement", file))
	}
	for _, rate := range mix.SuccessRates {
		if rate < 0 || rate > 100 {
			panic(fmt.Errorf("%s: invalid mix success rate %d", file, rate))
		}
	}
	// convert
	o.items = make(map[int]map[int]*item380)
	for _, item := range cfg.ItemOption.Items {
		if _, err := ItemTable.GetItemBase(item.Cat, item.Index); err != nil {
			panic(fmt.Errorf("%s: section %d index %d: %w", file, item.Cat, item.Index, err))
		}
		if item.Option1 < item380EffectNone || item.Option1 > item380EffectIncRecoverySDRate ||
			item.Option2 < item380EffectNone || item.Option2 > item380EffectIncRecoverySDRate ||
			item.Value1 < 0 || item.Value2 < 0 {
			panic(fmt.Errorf("%s: invalid options for section %d index %d", file, item.Cat, item.Index))
		}
		if o.items[item.Cat] == nil {
			o.items[item.Cat] = make(map[int]*item380)
		}
		if o.items[item.Cat][item.Index] != nil {
			panic(fmt.Errorf("%s: duplicate item section %d index %d", file, item.Cat, item.Index))
		}
		i380 := new(item380)
		i380.effects[0] = item380EffectPair{item.Option1, item.Value1}
		i380.effects[1] = item380EffectPair{item.Option2, item.Value2}
		o.items[item.Cat][item.Index] = i380
	}
	o.mix = mix
}

func (o *item380Manager) MixConfig() Item380Mix {
	return o.mix
}

func (o *item380Manager) Is380Item(section, index int) bool {
	_, ok := o.items[section][index]
	return ok
}

func (o *item380Manager) Apply380ItemEffect(section, index int, effect *Item380Effect) {
	entry := o.items[section][index]
	if entry == nil || effect == nil {
		return
	}
	pairs := entry.effects[:]
	for _, pair := range pairs {
		key := pair.Key
		value := pair.Value
		switch key {
		case item380EffectIncPVPAttackRate:
			effect.Item380EffectIncPVPAttackRate += value
		case item380EffectIncPVPAttack:
			effect.Item380EffectIncPVPAttack += value
		case item380EffectIncPVPDefenseRate:
			effect.Item380EffectIncPVPDefenseRate += value
		case item380EffectIncPVPDefense:
			effect.Item380EffectIncPVPDefense += value
		case item380EffectIncMaxHP:
			effect.Item380EffectIncMaxHP += value
		case item380EffectIncMaxSD:
			effect.Item380EffectIncMaxSD += value
		case item380EffectAutoRecoverySD:
			effect.Item380EffectAutoRecoverySD = true
		case item380EffectIncRecoverySDRate:
			effect.Item380EffectIncRecoverySDRate += value
		}
	}
}
