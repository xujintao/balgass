package item

// Item represents a item
type Item struct {
	*ItemBase      `json:"-"`
	Position       int               `json:"position"`
	ID             int               `json:"id"`      // serial
	Section        int               `json:"section"` // 0 ~ 15
	Index          int               `json:"index"`   // 0 ~ 511
	Code           int               `json:"-"`       // section*512 + index
	Level          int               `json:"level"`
	Durability     int               `json:"durability"`
	Lucky          bool              `json:"lucky,omitempty"`
	Skill          bool              `json:"skill,omitempty"`
	Addition       int               `json:"addition,omitempty"` // 0/4/8/12/16
	Excel          int               `json:"excel,omitempty"`
	Set            int               `json:"set,omitempty"`
	Option380      bool              `json:"option380,omitempty"`
	Period         int               `json:"period,omitempty"`
	HarmonyEffect  harmonyEffectKind `json:"harmony_effect,omitempty"`
	HarmonyLevel   int               `json:"harmony_level,omitempty"`
	PentagramBonus int               `json:"pentagram_bonus,omitempty"`
	MuunRank       int               `json:"muun_rank,omitempty"`
	SocketBonus    int               `json:"socket_bonus,omitempty"`
	SocketSlots    [5]int            `json:"-"` // slot array
	SocketSlot1    int               `json:"socket_slot1,omitempty"`
	SocketSlot2    int               `json:"socket_slot2,omitempty"`
	SocketSlot3    int               `json:"socket_slot3,omitempty"`
	SocketSlot4    int               `json:"socket_slot4,omitempty"`
	SocketSlot5    int               `json:"socket_slot5,omitempty"`
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

func (i *Item) GetExcelItem() int {
	return i.Excel & 0x3F
}

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
	data[0] = byte(item.Index)
	data[1] = byte(item.Addition & 0x0C >> 2)
	lucky := 0
	if item.Lucky {
		lucky = 1
	}
	data[1] |= byte(lucky << 2)
	data[1] |= byte(item.Level << 3)
	skill := 0
	if item.Skill {
		skill = 1
	}
	data[1] |= byte(skill << 7)
	data[2] = byte(item.Durability)
	data[3] = byte(item.Excel)
	data[3] |= byte(item.Addition & 0x10 << 2)
	data[3] |= byte(item.Index & 0x100 >> 1)
	data[4] = byte(item.Set)
	data[5] = byte(item.Period << 1)
	option380 := 0
	if item.Option380 {
		option380 = 1
	}
	data[5] |= byte(option380 << 3)
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
