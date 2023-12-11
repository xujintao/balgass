package item

// Item represents a item
type Item struct {
	*ItemBase                `json:"-"`
	Position                 int               `json:"position"`
	ID                       int               `json:"id"`      // serial
	Section                  int               `json:"section"` // 0 ~ 15
	Index                    int               `json:"index"`   // 0 ~ 511
	Code                     int               `json:"-"`       // section*512 + index
	Level                    int               `json:"level"`
	Durability               int               `json:"durability"`
	Lucky                    bool              `json:"lucky,omitempty"`
	Skill                    bool              `json:"skill,omitempty"`
	Addition                 int               `json:"addition,omitempty"` // 0/4/8/12/16
	RecoverHP                float64           `json:"-"`
	ExcellentAttackRate      bool              `json:"excellent_attack_rate,omitempty"`
	ExcellentAttackLevel     bool              `json:"excellent_attack_level,omitempty"`
	ExcellentAttackPercent   bool              `json:"excellent_attack_percent,omitempty"`
	ExcellentAttackSpeed     bool              `json:"excellent_attack_speed,omitempty"`
	ExcellentAttackHP        bool              `json:"excellent_attack_hp,omitempty"`
	ExcellentAttackMP        bool              `json:"excellent_attack_mp,omitempty"`
	ExcellentDefenseHP       bool              `json:"excellent_defense_hp,omitempty"`
	ExcellentDefenseMP       bool              `json:"excellent_defense_mp,omitempty"`
	ExcellentDefenseDecrease bool              `json:"excellent_defense_decrease,omitempty"`
	ExcellentDefenseReflect  bool              `json:"excellent_defense_reflect,omitempty"`
	ExcellentDefenseRate     bool              `json:"excellent_defense_rate,omitempty"`
	ExcellentDefenseMoney    bool              `json:"excellent_defense_money,omitempty"`
	Set                      int               `json:"set,omitempty"`
	Option380                bool              `json:"option380,omitempty"`
	Period                   int               `json:"period,omitempty"`
	HarmonyEffect            harmonyEffectKind `json:"harmony_effect,omitempty"`
	HarmonyLevel             int               `json:"harmony_level,omitempty"`
	PentagramBonus           int               `json:"pentagram_bonus,omitempty"`
	MuunRank                 int               `json:"muun_rank,omitempty"`
	SocketBonus              int               `json:"socket_bonus,omitempty"`
	SocketSlots              [5]int            `json:"-"` // slot array
	SocketSlot1              int               `json:"socket_slot1,omitempty"`
	SocketSlot2              int               `json:"socket_slot2,omitempty"`
	SocketSlot3              int               `json:"socket_slot3,omitempty"`
	SocketSlot4              int               `json:"socket_slot4,omitempty"`
	SocketSlot5              int               `json:"socket_slot5,omitempty"`
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
		i.ExcellentDefenseDecrease ||
		i.ExcellentDefenseReflect ||
		i.ExcellentDefenseRate ||
		i.ExcellentDefenseMoney
}

func (i *Item) IsSet() bool {
	return i.Set > 0
}

func (i *Item) CalcMaxDurability() int {
	dur := i.ItemBase.Durability
	if i.KindA == KindAPentagram {
		return dur
	}
	switch i.Level {
	case 0, 1, 2, 3, 4:
		dur += i.Level
	case 5, 6, 7, 8, 9:
		dur += i.Level*2 - 4
	case 10:
		dur += i.Level*2 - 3
	case 11:
		dur += i.Level*2 - 1
	case 12:
		dur += i.Level*2 + 2
	case 13:
		dur += i.Level*2 + 6
	case 14:
		dur += i.Level*2 + 11
	case 15:
		dur += i.Level*2 + 17
	}
	if i.KindA != KindAWing && i.Type != TypeArchangel {
		switch {
		case i.IsExcellent():
			dur += 15
		case i.IsSet():
			dur += 20
		}
	}
	if dur > 255 {
		dur = 255
	}
	return dur
}

func (i *Item) GetSkillIndex() int {
	if i.Skill {
		if i.Code == Code(12, 11) { // 召唤之石
			return i.ItemBase.SkillIndex + i.Level
		}
		return i.ItemBase.SkillIndex
	}
	return 0
}

func (i *Item) GetSetTierIndex() int {
	return i.Set & 3
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
	if item.ExcellentAttackLevel || item.ExcellentDefenseMP {
		data[3] |= 1 << 4
	}
	if item.ExcellentAttackPercent || item.ExcellentDefenseDecrease {
		data[3] |= 1 << 3
	}
	if item.ExcellentAttackSpeed || item.ExcellentDefenseReflect {
		data[3] |= 1 << 2
	}
	if item.ExcellentAttackHP || item.ExcellentDefenseRate {
		data[3] |= 1 << 1
	}
	if item.ExcellentAttackMP || item.ExcellentDefenseMoney {
		data[3] |= 1 << 0
	}
	data[3] |= byte(item.Addition & 0x10 << 2)
	data[3] |= byte(item.Index & 0x100 >> 1)
	data[4] = byte(item.Set)
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
		data[6] = 0
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
