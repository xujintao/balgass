package item

import (
	"database/sql/driver"
	"encoding/json"
	"log/slog"
)

// Item represents a item
type Item struct {
	*ItemBase                        `json:"-"`
	Position                         int    `json:"position"`
	ID                               int    `json:"id"`      // serial
	Section                          int    `json:"section"` // 0 ~ 15
	Index                            int    `json:"index"`   // 0 ~ 511
	Code                             int    `json:"-"`       // section*512 + index
	Level                            int    `json:"level"`
	Durability                       int    `json:"durability"`
	Lucky                            bool   `json:"lucky,omitempty"`
	Skill                            bool   `json:"skill,omitempty"`
	Addition                         int    `json:"addition,omitempty"`                  // 0/4/8/12/16
	ExcellentAttackRate              bool   `json:"excellent_attack_rate,omitempty"`     // bit5:卓越攻击几率10%
	ExcellentAttackLevel             bool   `json:"excellent_attack_level,omitempty"`    // bit4:攻击力增加等级/20
	ExcellentAttackPercent           bool   `json:"excellent_attack_percent,omitempty"`  // bit3:攻击力增加2%
	ExcellentAttackSpeed             bool   `json:"excellent_attack_speed,omitempty"`    // bit2:攻击(魔法)速度增加7
	ExcellentAttackHP                bool   `json:"excellent_attack_hp,omitempty"`       // bit1:杀死怪物时所获生命值增加生命值/8
	ExcellentAttackMP                bool   `json:"excellent_attack_mp,omitempty"`       // bit0:杀死怪物时所获魔法值增加魔法值/8
	ExcellentDefenseHP               bool   `json:"excellent_defense_hp,omitempty"`      // bit5:最大生命值增加4%
	ExcellentDefenseMP               bool   `json:"excellent_defense_mp,omitempty"`      // bit4:最大魔法值增加4%
	ExcellentDefenseReduce           bool   `json:"excellent_defense_reduce,omitempty"`  // bit3:伤害减少4%
	ExcellentDefenseReflect          bool   `json:"excellent_defense_reflect,omitempty"` // bit2:伤害反射5%
	ExcellentDefenseRate             bool   `json:"excellent_defense_rate,omitempty"`    // bit1:防御成功率10%
	ExcellentDefenseMoney            bool   `json:"excellent_defense_money,omitempty"`   // bit0:杀死怪物时所获金币增加30%
	ExcellentWingAdditionAttack      bool   `json:"excellent_wing_addition_attack,omitempty"`
	ExcellentWingAdditionMagicAttack bool   `json:"excellent_wing_addition_magic_attack,omitempty"`
	ExcellentWingAdditionCurseAttack bool   `json:"excellent_wing_addition_curse_attack,omitempty"`
	ExcellentWingAdditionDefense     bool   `json:"excellent_wing_addition_defense,omitempty"`
	ExcellentWingAdditionRecoverHP   bool   `json:"excellent_wing_addition_recover_hp,omitempty"`
	ExcellentWing2Speed              bool   `json:"excellent_wing2_speed,omitempty"`      // 2D:攻击(魔法)速度增加5
	ExcellentWing2AG                 bool   `json:"excellent_wing2_ag,omitempty"`         // 2D:最大AG增加50
	ExcellentWing2Leadership         bool   `json:"excellent_wing2_leadership,omitempty"` // 2D:声望增加10+5*level
	ExcellentWing2Ignore             bool   `json:"excellent_wing2_ignore,omitempty"`     // 2D:无视防御伤害几率3%
	ExcellentWing2MP                 bool   `json:"excellent_wing2_mp,omitempty"`         // 2D:魔法值增加50+5*level
	ExcellentWing2HP                 bool   `json:"excellent_wing2_hp,omitempty"`         // 2D:生命值增加50+5*level
	ExcellentWing3MP                 bool   `json:"excellent_wing3_mp,omitempty"`         // 3D:魔法值完全恢复几率5%
	ExcellentWing3HP                 bool   `json:"excellent_wing3_hp,omitempty"`         // 3D:生命值完全恢复几率5%
	ExcellentWing3Return             bool   `json:"excellent_wing3_return,omitempty"`     // 3D:反弹攻击力几率5%
	ExcellentWing3Ignore             bool   `json:"excellent_wing3_ignore,omitempty"`     // 3D:无视防御伤害几率5%
	ExcellentWing25HP                bool   `json:"excellent_wing25_hp,omitempty"`        // 2.5D:生命值增加5%
	ExcellentWing25Ignore            bool   `json:"excellent_wing25_ignore,omitempty"`    // 2.5D:无视防御伤害几率3%
	Set                              int    `json:"set,omitempty"`
	Option380                        bool   `json:"option380,omitempty"`
	Period                           int    `json:"period,omitempty"`
	HarmonyEffect                    int    `json:"harmony_effect,omitempty"`
	HarmonyLevel                     int    `json:"harmony_level,omitempty"`
	HarmonyOption                    int    `json:"harmony_option,omitempty"`
	PentagramBonus                   int    `json:"pentagram_bonus,omitempty"`
	MuunRank                         int    `json:"muun_rank,omitempty"`
	SocketBonus                      int    `json:"socket_bonus,omitempty"`
	SocketSlots                      [5]int `json:"-"` // slot array
	SocketSlot1                      int    `json:"socket_slot1,omitempty"`
	SocketSlot2                      int    `json:"socket_slot2,omitempty"`
	SocketSlot3                      int    `json:"socket_slot3,omitempty"`
	SocketSlot4                      int    `json:"socket_slot4,omitempty"`
	SocketSlot5                      int    `json:"socket_slot5,omitempty"`
	MaxDurability                    int    `json:"-"`
	AttackMin                        int    `json:"-"`
	AttackMax                        int    `json:"-"`
	Magic                            int    `json:"-"`
	Defense                          int    `json:"-"`
	DefenseRate                      int    `json:"-"`
	AdditionAttack                   int    `json:"-"`
	AdditionMagicAttack              int    `json:"-"`
	AdditionCurseAttack              int    `json:"-"`
	AdditionDefense                  int    `json:"-"`
	AdditionDefenseRate              int    `json:"-"`
	AdditionRecoverHP                int    `json:"-"`
	AdditionAbsorbDamagePercent5     int    `json:"-"`
	AdditionAG50                     int    `json:"-"`
	AdditionSpeed5                   int    `json:"-"`
	Money                            int    `json:"-"`
}

// NewItem construct a item with section and index
func NewItem(section, index int) *Item {
	return &Item{
		ID:       0,
		Section:  section,
		Index:    index,
		Code:     Code(section, index),
		ItemBase: ItemTable.GetItemBaseMust(section, index),
	}
}

func (i *Item) Copy() *Item {
	var copyItem Item = *i
	return &copyItem
}

func (i *Item) IsExcellent() bool {
	return i.ExcellentAttackRate ||
		i.ExcellentAttackLevel ||
		i.ExcellentAttackPercent ||
		i.ExcellentAttackSpeed ||
		i.ExcellentAttackHP ||
		i.ExcellentAttackMP ||
		i.ExcellentDefenseHP ||
		i.ExcellentDefenseMP ||
		i.ExcellentDefenseReduce ||
		i.ExcellentDefenseReflect ||
		i.ExcellentDefenseRate ||
		i.ExcellentDefenseMoney
}

func (i *Item) IsSet() bool {
	return i.Set > 0
}

func (it *Item) GetSkillIndex() int {
	if it.Skill {
		if it.Code == Code(12, 11) { // 召唤之石
			return it.SkillIndex + it.Level
		}
		return it.SkillIndex
	}
	return 0
}

func (it *Item) Calc() {
	// calc dur
	dur := it.ItemBase.Durability
	if dur > 0 && it.KindA != KindAPentagram {
		switch it.Level {
		case 0, 1, 2, 3, 4:
			dur += it.Level
		case 5, 6, 7, 8, 9:
			dur += it.Level*2 - 4
		case 10:
			dur += it.Level*2 - 3
		case 11:
			dur += it.Level*2 - 1
		case 12:
			dur += it.Level*2 + 2
		case 13:
			dur += it.Level*2 + 6
		case 14:
			dur += it.Level*2 + 11
		case 15:
			dur += it.Level*2 + 17
		}
		if it.KindA != KindAWing && it.Type != TypeArchangel {
			switch {
			case it.IsExcellent():
				dur += 15
			case it.IsSet():
				dur += 20
			}
		}
		if dur > 255 {
			dur = 255
		}
		it.MaxDurability = dur
	}

	level := it.ItemBase.DropLevel
	// calc required strength/dexterity/vitality/energy/leadership/level
	if it.IsExcellent() || it.IsSet() {
		level = it.ItemBase.DropLevel + 25
	}
	// ...

	// calc attack
	if it.IsSet() {
		level = it.ItemBase.DropLevel + 30
	}
	attackMin := it.ItemBase.DamageMin
	attackMax := it.ItemBase.DamageMax
	if attackMin > 0 && attackMax > 0 {
		delta := 0
		switch {
		case it.IsExcellent():
			delta += attackMin*25/it.ItemBase.DropLevel + 5
		case it.IsSet():
			delta += attackMin*25/it.ItemBase.DropLevel + 5
			delta += level/40 + 5
		}
		delta += it.Level * 3
		if it.Level >= 10 {
			delta += (it.Level - 9) * (it.Level - 8) / 2
		}
		it.AttackMin = attackMin + delta
		it.AttackMax = attackMax + delta
	}

	// calc magic
	magic := it.ItemBase.MagicPower
	if magic > 0 {
		delta := 0
		switch {
		case it.IsExcellent():
			delta += magic*25/it.ItemBase.DropLevel + 5
		case it.IsSet():
			delta += magic*25/it.ItemBase.DropLevel + 5
			delta += level/60 + 2
		}
		delta += it.Level * 3
		if it.Level >= 10 {
			delta += (it.Level - 9) * (it.Level - 8) / 2
		}
		it.Magic = (magic+delta)/2 + it.Level*2
	}

	// calc defense
	defense := it.ItemBase.Defense
	if defense > 0 {
		delta := 0
		switch it.KindB {
		case KindBShield:
			delta += it.Level
			if it.IsSet() || it.IsExcellent() {
				delta += defense*20/level + 2
			}
			delta += it.Level * 3
			if it.Level >= 10 {
				delta += (it.Level - 9) * (it.Level - 8) / 2
			}
		case KindBWing1st, KindBWingMonster:
			delta += it.Level * 3
			if it.Level >= 10 {
				delta += (it.Level - 9) * (it.Level - 8) / 2
			}
		case KindBWing2nd, KindBCapeFighter:
			delta += it.Level * 2
			if it.Level >= 10 {
				delta += (it.Level-9)*(it.Level-8)/2 + it.Level - 9
			}
		case KindBCapeLord:
			delta += it.Level*2 + 15
			if it.Level >= 10 {
				delta += (it.Level-9)*(it.Level-8)/2 + it.Level - 9
			}
		case KindBWing3rd:
			delta += it.Level * 4
			if it.Level >= 10 {
				delta += (it.Level - 9) * (it.Level - 8) / 2
			}
		default: // armor
			switch {
			case it.IsExcellent():
				delta += defense*12/it.ItemBase.DropLevel + it.ItemBase.DropLevel/5 + 4
			case it.IsSet():
				delta += defense*12/it.ItemBase.DropLevel + it.ItemBase.DropLevel/5 + 4
				delta += defense*3/level + level/30 + 2
			}
			delta += it.Level * 3
			if it.Level >= 10 {
				delta += (it.Level - 9) * (it.Level - 8) / 2
			}
		}
		it.Defense = defense + delta
	}

	// calc defense rate
	defenseRate := it.ItemBase.SuccessfulBlocking
	if defenseRate > 0 {
		delta := 0
		switch {
		case it.IsExcellent():
			delta += defenseRate*25/it.ItemBase.DropLevel + 5
		case it.IsSet():
			delta += defenseRate*25/it.ItemBase.DropLevel + 5
			delta += level/40 + 5
		}
		delta += it.Level * 3
		if it.Level >= 10 {
			delta += (it.Level - 9) * (it.Level - 8) / 2
		}
		it.DefenseRate = defenseRate + delta
	}

	// calc addition
	switch {
	case it.Code >= Code(0, 0) && it.Code < Code(5, 0): // weapon
		it.AdditionAttack = it.Addition
	case it.Code >= Code(5, 0) && it.Code < Code(6, 0): // staff
		it.AdditionMagicAttack = it.Addition
	case it.Code >= Code(6, 0) && it.Code < Code(7, 0): // shield
		it.AdditionDefenseRate = it.Addition
	case it.Code >= Code(7, 0) && it.Code < Code(12, 0): // armor
		it.AdditionDefense = it.Addition
	case it.KindA == KindAPendant || it.KindA == KindARing: // pendant/ring
		switch it.Code {
		case Code(13, 24):
		case Code(13, 28):
		default:
		}
		it.AdditionRecoverHP = it.Addition / 4
	case it.Code == Code(13, 3): // Horn of Dinorant 彩云兽
		if it.Addition&0x01 != 0 {
			it.AdditionAbsorbDamagePercent5 = 5
		}
		if it.Addition&0x02 != 0 {
			it.AdditionAG50 = 50
		}
		if it.Addition&0x04 != 0 {
			it.AdditionSpeed5 = 5
		}
	case it.KindA == KindAWing: // wings
		switch {
		case it.ExcellentWingAdditionAttack:
			it.AdditionAttack = it.Addition
		case it.ExcellentWingAdditionMagicAttack:
			it.AdditionMagicAttack = it.Addition
		case it.ExcellentWingAdditionCurseAttack:
			it.AdditionCurseAttack = it.Addition
		case it.ExcellentWingAdditionDefense:
			it.AdditionDefense = it.Addition
		case it.ExcellentWingAdditionRecoverHP:
			it.AdditionRecoverHP = it.Addition / 4
		}
	}
}

// func (i *Item) GetExcelItem() int {
// 	return i.Excel & 0x3F
// }

// Marshal marshal item struct to [32]byte variable
// ----------------------------------------------
// field[0]
// bit0~bit7: itembase index, 0~255
// ----------------------------------------------
// field[1]
// bit0~bit1: addition attack/defense
// bit2: lucky flag
// bit3~bit6: level, 0~15
// bit7: skill flag
// ----------------------------------------------
// field[2]
// bit0~bit7: durability, 0~255
// ----------------------------------------------
// field[3]~field[6]
// serial number
// ----------------------------------------------
// field[7]
// bit0~bit5: excellent option
// bit6: addition attack/defense 16 flag
// now field[7].bit6, field[1].bit1 and field[1].bit0 may range as follow:
// 001: addition 4
// 010: addition 8
// 011: addition 12
// 100: addition 16
// bit7: extension flag of itembase index , act as the bit8 of field[0]
// now itembase index is 9 bit and range is 0~511
// ----------------------------------------------
// field[8]
// set or ancient
// ----------------------------------------------
// field[9]
// bit0: period
// bit1: period expire
// bit3: option380 flag
// bit4~bit7: itembase section, 0~15
// ----------------------------------------------
// field[10] ~ field[15]
// socket index and socket slots
// ----------------------------------------------
// field[16] ~ field[19]
// extension of serial number, now serial number is 8 bytes
// ----------------------------------------------
// field[20] ~ field[31]
// total 12 bytes, reserved

// ----------------------------------------------
// ----------------------------------------------
// Marshal marshal item struct to [12]byte variable
// ----------------------------------------------
// field[0]
// bit0~bit7: itembase index, 0~255
// ----------------------------------------------
// field[1]
// bit0~bit1: addition attack/defense
// bit2: lucky flag
// bit3~bit6: level, 0~15
// bit7: skill flag
// ----------------------------------------------
// field[2]
// bit0~bit7: durability, 0~255
// ----------------------------------------------
// field[3]
// bit0~bit5: excellent option
// bit6: addition attack/defense 16 flag
// now field[3].bit6, field[1].bit1 and field[1].bit0 may range as follow:
// 001: addition 4
// 010: addition 8
// 011: addition 12
// 100: addition 16
// bit7: extension flag of itembase index , act as the bit8 of field[0]
// now itembase index is 9 bit and range is 0~511
// ----------------------------------------------
// field[4]
// set or ancient
// ----------------------------------------------
// field[5]
// bit0: period
// bit1: period expire
// bit3: option380 flag
// bit4~bit7: itembase section, 0~15
// ----------------------------------------------
// field[6] ~ field[11]
// socket index and socket slots
func (item *Item) Marshal() ([]byte, error) {
	var data [12]byte
	if item == nil {
		return data[:], nil
	}
	if item.Code == Code(14, 15) {
		money := item.Durability
		data[0] = byte(item.Index)
		data[1] = byte(money >> 16)
		data[2] = byte(money >> 8)
		data[4] = byte(money >> 0)
		data[5] = byte(item.Section << 4)
		return data[:], nil
	}
	data[0] = byte(item.Index)
	data[1] = byte(item.Addition & 0x0C >> 2)
	if item.Lucky {
		data[1] |= byte(1 << 2)
	}
	data[1] |= byte(item.Level << 3)
	if item.Skill {
		data[1] |= byte(1 << 7)
	}
	data[2] = byte(item.Durability)
	if item.ExcellentAttackRate || item.ExcellentDefenseHP {
		data[3] |= 1 << 5
	}
	if item.ExcellentAttackLevel || item.ExcellentDefenseMP || item.ExcellentWing2Speed {
		data[3] |= 1 << 4
	}
	switch {
	case item.Code == Code(12, 3) && item.ExcellentWingAdditionRecoverHP, // Wings of Spirits 圣灵之翼
		item.Code == Code(12, 4) && item.ExcellentWingAdditionMagicAttack,   // Wings of Soul 魔魂之翼
		item.Code == Code(12, 5) && item.ExcellentWingAdditionAttack,        // Wings of Dragon 飞龙之翼
		item.Code == Code(12, 6) && item.ExcellentWingAdditionAttack,        // Wings of Darkness 暗黑之翼
		item.Code == Code(13, 30) && item.ExcellentWingAdditionAttack,       // Wings of Darkness 王者披风
		item.Code == Code(12, 42) && item.ExcellentWingAdditionMagicAttack,  // Wing of Despair 绝望之翼
		item.Code == Code(12, 49) && item.ExcellentWingAdditionAttack,       // Cape of Fighter 武者披风
		item.Code == Code(12, 36) && item.ExcellentWingAdditionDefense,      // Wing of Storm 暴风之翼
		item.Code == Code(12, 37) && item.ExcellentWingAdditionDefense,      // Wing of Eternal 时空之翼
		item.Code == Code(12, 38) && item.ExcellentWingAdditionDefense,      // Wing of Illusion 幻影之翼
		item.Code == Code(12, 39) && item.ExcellentWingAdditionMagicAttack,  // Wing of Ruin 破灭之翼
		item.Code == Code(12, 40) && item.ExcellentWingAdditionDefense,      // Cape of Emperor 帝王披风
		item.Code == Code(12, 43) && item.ExcellentWingAdditionCurseAttack,  // Wing of Dimension 次元之翼
		item.Code == Code(12, 50) && item.ExcellentWingAdditionDefense,      // Cape of Overrule 斗皇披风
		item.Code == Code(12, 262) && item.ExcellentWingAdditionRecoverHP,   // Cloak of Death 死亡披风
		item.Code == Code(12, 263) && item.ExcellentWingAdditionRecoverHP,   // Wings of Chaos 混沌之翼
		item.Code == Code(12, 264) && item.ExcellentWingAdditionCurseAttack, // Wings of Magic 魔力之翼
		item.Code == Code(12, 265) && item.ExcellentWingAdditionRecoverHP:   // Wings of Life 生命之翼
		data[3] |= 1 << 5
	case item.Code == Code(12, 36) && item.ExcellentWingAdditionAttack, // Wing of Storm 暴风之翼
		item.Code == Code(12, 37) && item.ExcellentWingAdditionMagicAttack,  // Wing of Eternal 时空之翼
		item.Code == Code(12, 38) && item.ExcellentWingAdditionAttack,       // Wing of Illusion 幻影之翼
		item.Code == Code(12, 39) && item.ExcellentWingAdditionAttack,       // Wing of Ruin 破灭之翼
		item.Code == Code(12, 40) && item.ExcellentWingAdditionAttack,       // Cape of Emperor 帝王披风
		item.Code == Code(12, 43) && item.ExcellentWingAdditionMagicAttack,  // Wing of Dimension 次元之翼
		item.Code == Code(12, 50) && item.ExcellentWingAdditionAttack,       // Cape of Overrule 斗皇披风
		item.Code == Code(12, 262) && item.ExcellentWingAdditionAttack,      // Cloak of Death 死亡披风
		item.Code == Code(12, 263) && item.ExcellentWingAdditionAttack,      // Wings of Chaos 混沌之翼
		item.Code == Code(12, 264) && item.ExcellentWingAdditionMagicAttack, // Wings of Magic 魔力之翼
		item.Code == Code(12, 265) && item.ExcellentWingAdditionAttack:      // Wings of Life 生命之翼
		data[3] |= 1 << 4
	}
	if item.ExcellentAttackPercent || item.ExcellentDefenseReduce ||
		item.ExcellentWing2AG || item.ExcellentWing2Leadership || item.ExcellentWing3MP {
		data[3] |= 1 << 3
	}
	if item.ExcellentAttackSpeed || item.ExcellentDefenseReflect ||
		item.ExcellentWing2Ignore || item.ExcellentWing3HP {
		data[3] |= 1 << 2
	}
	if item.ExcellentAttackHP || item.ExcellentDefenseRate ||
		item.ExcellentWing2MP || item.ExcellentWing3Return || item.ExcellentWing25HP {
		data[3] |= 1 << 1
	}
	if item.ExcellentAttackMP || item.ExcellentDefenseMoney ||
		item.ExcellentWing2HP || item.ExcellentWing3Ignore || item.ExcellentWing25Ignore {
		data[3] |= 1 << 0
	}
	data[3] |= byte(item.Addition & 0x10 << 2)
	data[3] |= byte(item.Index & 0x100 >> 1)
	data[4] = byte(SetManager.GetTierIndex(item.Section, item.Index, item.Set))
	data[5] = byte(item.Period << 1)
	if item.Option380 {
		data[5] |= byte(1 << 3)
	}
	data[5] |= byte(item.Section << 4)
	// harmony/socket/pentagram/muun system
	if item.Type == TypeSocket ||
		item.KindA == KindAPentagram {
		data[6] = byte(item.SocketBonus)
		data[7] = byte(item.SocketSlot1)
		data[8] = byte(item.SocketSlot2)
		data[9] = byte(item.SocketSlot3)
		data[10] = byte(item.SocketSlot4)
		data[11] = byte(item.SocketSlot5)
	} else {
		data[6] = byte(item.HarmonyOption)
		data[7] = 0xFF
		data[8] = 0xFF
		data[9] = 0xFF
		data[10] = 0xFF
		data[11] = 0xFF
	}
	return data[:], nil
}

func (item *Item) Unmarshal(buf []byte) error {
	return nil
}

type PositionedItems struct {
	Size              int
	Items             []*Item
	Flags             []bool
	CheckFlagsForItem func(int, *Item) bool
	SetFlagsForItem   func(int, *Item)
}

func (pi PositionedItems) MarshalJSON() ([]byte, error) {
	var items []*Item
	for i, v := range pi.Items {
		if v == nil {
			continue
		}
		v.Position = i
		items = append(items, v)
	}
	data, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (pi *PositionedItems) UnmarshalJSON(buf []byte) error {
	var items []*Item
	err := json.Unmarshal(buf, &items)
	if err != nil {
		return err
	}
	pi.Items = make([]*Item, pi.Size)
	pi.Flags = make([]bool, pi.Size)
	for _, v := range items {
		v.Code = Code(v.Section, v.Index)
		itemBase, err := ItemTable.GetItemBase(v.Section, v.Index)
		if err != nil {
			return err
		}
		v.ItemBase = itemBase
		v.Calc()
		ok := pi.CheckFlagsForItem(v.Position, v)
		if !ok {
			slog.Error("pi.CheckFlagsForItem",
				"position", v.Position, "name", v.Name, "annotation", v.Annotation)
			continue
		}
		pi.SetFlagsForItem(v.Position, v)
		pi.Items[v.Position] = v
	}
	return nil
}

func (pi PositionedItems) Value() (driver.Value, error) {
	return pi.MarshalJSON()
}
