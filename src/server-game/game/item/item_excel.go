package item

import (
	"encoding/xml"
	"math/rand"

	"github.com/xujintao/balgass/src/server-game/conf"
)

func init() {
	ExcellentManager.init()
}

var ExcellentManager excellentManager

type excellentManager struct {
	regulars map[int][]int
	wings    map[int][]int
	multiple []int
}

func (m *excellentManager) init() {
	// ExcellentOptions was generated 2024-01-15 17:48:16 by https://xml-to-go.github.io/ in Ukraine.
	type ExcellentOptions struct {
		XMLName xml.Name `xml:"ExcellentOptions"`
		Text    string   `xml:",chardata"`
		Common  struct {
			Text   string `xml:",chardata"`
			Option []struct {
				Text       string `xml:",chardata"`
				ID         int    `xml:"ID,attr"`
				Number     int    `xml:"Number,attr"`
				Value      int    `xml:"Value,attr"`
				ItemKindA1 int    `xml:"ItemKindA_1,attr"`
				ItemKindA2 int    `xml:"ItemKindA_2,attr"`
				ItemKindA3 int    `xml:"ItemKindA_3,attr"`
				Rate       int    `xml:"Rate,attr"`
				Name       string `xml:"Name,attr"`
			} `xml:"Option"`
		} `xml:"Common"`
		Wings struct {
			Text   string `xml:",chardata"`
			Option []struct {
				Text      string `xml:",chardata"`
				ID        int    `xml:"ID,attr"`
				Number    int    `xml:"Number,attr"`
				Value     int    `xml:"Value,attr"`
				ItemKindA int    `xml:"ItemKindA,attr"`
				ItemKindB int    `xml:"ItemKindB,attr"`
				Name      string `xml:"Name,attr"`
			} `xml:"Option"`
		} `xml:"Wings"`
		OptionDropRate struct {
			Text   string `xml:",chardata"`
			Common struct {
				Text  string `xml:",chardata"`
				One   int    `xml:"One,attr"`
				Two   int    `xml:"Two,attr"`
				Three int    `xml:"Three,attr"`
				Four  int    `xml:"Four,attr"`
				Five  int    `xml:"Five,attr"`
				Six   int    `xml:"Six,attr"`
			} `xml:"Common"`
		} `xml:"OptionDropRate"`
	}
	var excellentOptions ExcellentOptions
	conf.XML(conf.PathCommon, "Items/IGC_ExcellentOptions.xml", &excellentOptions)
	// regular
	m.regulars = make(map[int][]int)
	for _, v := range excellentOptions.Common.Option {
		m.regulars[v.ItemKindA1] = append(m.regulars[v.ItemKindA1], v.Number)
		m.regulars[v.ItemKindA2] = append(m.regulars[v.ItemKindA2], v.Number)
	}
	// wing
	m.wings = make(map[int][]int)
	for _, v := range excellentOptions.Wings.Option {
		m.wings[v.ItemKindB] = append(m.wings[v.ItemKindB], v.Number)
	}
	// multiple
	rate := excellentOptions.OptionDropRate.Common
	m.multiple = append(m.multiple,
		rate.One,
		rate.Two,
		rate.Three,
		rate.Four,
		rate.Five,
		rate.Six,
	)
}

func (m *excellentManager) DropExcellent(kindA itemKindA, kindB itemKindB) int {
	var pool []int
	switch kindA {
	case KindAWeapon, KindAPendant, KindAArmor, KindARing:
		pool = m.regulars[int(kindA)]
	case KindAWing:
		pool = m.wings[int(kindB)]
	}
	n := len(pool)
	if n <= 0 {
		return 0
	}
	cnt := m.dropExcellentCount()
	s1 := make(map[int]struct{})
	for cnt > 0 {
		i := rand.Intn(n)
		v := pool[i]
		_, ok := s1[v]
		if !ok {
			s1[v] = struct{}{}
			cnt--
		}
	}
	excel := 0
	for k := range s1 {
		excel |= 1 << k
	}
	return excel
}

func (m *excellentManager) dropExcellentCount() int {
	num := rand.Intn(10000)
	offset := 0
	for i, v := range m.multiple {
		if num >= offset && num < v+offset {
			return i + 1
		}
		offset += v
	}
	return 0
}

// ItemExcellentEffect describes the fixed bonuses, percentages and level divisors
// contributed by one item.
type ItemExcellentEffect struct {
	ExcellentAttackRate int     // 卓越攻击几率10%
	AttackLevel         int     // 攻击力增加等级/20
	MagicAttackLevel    int     // 攻击力(魔法)增加等级/20
	AttackPercent       int     // 攻击力增加2%
	MagicAttackPercent  int     // 攻击力(魔法)增加2%
	AttackSpeed         int     // 攻击(魔法)速度增加7/5(2D)
	MagicSpeed          int     // 攻击(魔法)速度增加7/5(2D)
	MonsterDieGetHP     float64 // 杀死怪物时所获生命值增加生命值/8
	MonsterDieGetMP     float64 // 杀死怪物时所获魔法值增加魔法值/8
	MaxHPPercent        int     // 最大生命值增加4%/5%(2.5D)
	MaxMPPercent        int     // 最大魔法值增加4%
	ReduceDamage        int     // 伤害减少4%
	ReflectDamage       int     // 伤害反射5%
	DefenseRatePercent  int     // 防御成功率10%
	MonsterDieGetMoney  float64 // 杀死怪物时所获金币增加30%
	MaxHP               int     // 生命值增加50+5*level(2D)
	MaxMP               int     // 魔法值增加50+5*level(2D)
	MaxAG               int     // 最大AG增加50(2D)
	IgnoreDefenseRate   int     // 无视防御伤害几率5%(3D)/3%(2.5D)
	ReturnDamage        int     // 反弹攻击力几率5%(3D)
	RecoverMaxHP        int     // 生命值完全恢复几率5%(3D)
	RecoverMaxMP        int     // 魔法值完全恢复几率5%(3D)
}

func (m *excellentManager) ApplyExcellentItemEffect(inventory *Inventory, apply func(ItemExcellentEffect, int)) {
	if inventory == nil || apply == nil {
		return
	}
	for position, it := range inventory.WearingItems() {
		if it == nil {
			continue
		}
		var effect ItemExcellentEffect
		m.applyExcellentEffect(it, &effect)
		apply(effect, position)
	}
}

func (m *excellentManager) applyExcellentEffect(it *Item, effect *ItemExcellentEffect) {
	if it.ExcellentAttackRate {
		effect.ExcellentAttackRate += 10
	}
	if it.ExcellentAttackLevel {
		switch it.KindA {
		case KindAWeapon:
			switch it.KindB {
			case KindBSwordMagicGladiator:
				effect.AttackLevel += 20
				effect.MagicAttackLevel += 20
			case KindBStaffSummoner, KindBBookSummoner:
				effect.MagicAttackLevel += 20
			default:
				effect.AttackLevel += 20
			}
		case KindAPendant:
			switch it.KindB {
			case KindBPendantAttackPhysical:
				effect.AttackLevel += 20
			case KindBPendantAttackMagic:
				effect.MagicAttackLevel += 20
			}
		}
	}
	if it.ExcellentAttackPercent {
		switch it.KindA {
		case KindAWeapon:
			switch it.KindB {
			case KindBSwordMagicGladiator:
				effect.AttackPercent += 2
				effect.MagicAttackPercent += 2
			case KindBStaffSummoner, KindBBookSummoner:
				effect.MagicAttackPercent += 2
			default:
				effect.AttackPercent += 2
			}
		case KindAPendant:
			switch it.KindB {
			case KindBPendantAttackPhysical:
				effect.AttackPercent += 2
			case KindBPendantAttackMagic:
				effect.MagicAttackPercent += 2
			}
		}
	}
	if it.ExcellentAttackSpeed {
		effect.AttackSpeed += 7
		effect.MagicSpeed += 7
	}
	if it.ExcellentAttackHP {
		effect.MonsterDieGetHP += 1.0 / 8
	}
	if it.ExcellentAttackMP {
		effect.MonsterDieGetMP += 1.0 / 8
	}
	if it.ExcellentDefenseHP {
		effect.MaxHPPercent += 4
	}
	if it.ExcellentDefenseMP {
		effect.MaxMPPercent += 4
	}
	if it.ExcellentDefenseReduce {
		effect.ReduceDamage += 4
	}
	if it.ExcellentDefenseReflect {
		effect.ReflectDamage += 5
	}
	if it.ExcellentDefenseRate {
		effect.DefenseRatePercent += 10
	}
	if it.ExcellentDefenseMoney {
		effect.MonsterDieGetMoney += float64(30) / 100
	}
	if it.Wing2ExcellentSpeed {
		effect.AttackSpeed += 5
		effect.MagicSpeed += 5
	}
	if it.Wing2ExcellentHP {
		effect.MaxHP += 50 + 5*it.Level
	}
	if it.Wing2ExcellentMP {
		effect.MaxMP += 50 + 5*it.Level
	}
	if it.Wing2ExcellentAG {
		effect.MaxAG += 50
	}
	if it.Wing2ExcellentIgnore {
		effect.IgnoreDefenseRate += 3
	}
	if it.Wing3ExcellentIgnore {
		effect.IgnoreDefenseRate += 5
	}
	if it.Wing3ExcellentReturn {
		effect.ReturnDamage += 5
	}
	if it.Wing3ExcellentHP {
		effect.RecoverMaxHP += 5
	}
	if it.Wing3ExcellentMP {
		effect.RecoverMaxMP += 5
	}
}
