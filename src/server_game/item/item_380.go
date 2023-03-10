package item

import (
	"path"

	"github.com/xujintao/balgass/src/server_game/conf"
)

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

type item380 struct {
	effects [2]item380EffectPair
}

type item380Mix struct {
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
}

type item380Manager struct {
	mix   item380Mix
	items map[int]map[int]*item380
}

func (o *item380Manager) Is380Item(section, index int) bool {
	_, ok := o.items[section][index]
	return ok
}

func (o *item380Manager) Apply380ItemEffect(section, index int, effect *Item380Effect) {
	pairs := o.items[section][index].effects[:]
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

var Item380Manager item380Manager

func init() {
	type Item380Xml struct {
		Mix        item380Mix `xml:"Mix"`
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
	var item380Xml Item380Xml
	conf.XML(path.Join(conf.PathCommon, "Items/IGC_Item380Option.xml"), &item380Xml)
	// convert
	Item380Manager.items = make(map[int]map[int]*item380)
	for _, item := range item380Xml.ItemOption.Items {
		_, ok := Item380Manager.items[item.Cat]
		if !ok {
			Item380Manager.items[item.Cat] = make(map[int]*item380)
		}
		_, ok = Item380Manager.items[item.Cat][item.Index]
		if !ok {
			i380 := new(item380)
			i380.effects[0] = item380EffectPair{item.Option1, item.Value1}
			i380.effects[1] = item380EffectPair{item.Option2, item.Value2}
			Item380Manager.items[item.Cat][item.Index] = i380
		}
	}
	Item380Manager.mix = item380Xml.Mix
}
