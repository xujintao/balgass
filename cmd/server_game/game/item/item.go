package item

import (
	"path"

	"github.com/xujintao/balgass/cmd/server_game/conf"
)

type itemType int

const (
	itemTypeCommon itemType = iota
	itemTypeRegular
	itemTypeSocket
	itemType380
	itemTypeLucky
	itemTypeEvent
	itemTypeArchangel
	itemTypeChaos
)

type itemKindA int

const (
	itemKindACommon itemKindA = iota
	itemKindAWeapon
	itemKindAPendant
	itemKindAArmor
	itemKindARing
	itemKindAHelper // imp etc.
	itemKindAWing
	itemKindAPremiumBuff
	itemKindAPentagram
	itemKindAJewel
	itemKindASkill
	itemKindAEvent
	itemKindAMuun
)

type itemKindB int

const (
	itemKindBCommon                  itemKindB = iota
	itemKindBSwordKnight                                  // 剑战士用
	itemKindBSwordMagicGladiator                          // 剑魔剑士用
	itemKindBSwordRageFighter                             // 剑格斗用
	itemKindBAxe                                          // 斧
	itemKindBMace                                         // 槌
	itemKindBScepter                                      // 权杖
	itemKindBSpear                                        // 矛
	itemKindBBow                                          // 弓
	itemKindBCrossbow                                     // 弩
	itemKindBArrow                                        // 弓箭
	itemKindBBolt                                         // 弩箭
	itemKindBStaffWizard                                  // 法师杖
	itemKindBStaffSummoner                                // 召唤杖
	itemKindBBookSummoner                                 // 召唤副手书
	itemKindBShield                                       // 盾类
	itemKindBHelmet                                       // 头盔
	itemKindBArmor                                        // 铠甲
	itemKindBPants                                        // 护腿
	itemKindBGloves                                       // 护手
	itemKindBBoots                                        // 靴子
	itemKindBWingTalisman                                 // 翅膀合成符咒
	itemKindBSmallWing                                    // 小翅膀
	itemKindBWing1st                                      // 1代翅膀
	itemKindBWing2nd                                      // 2代翅膀
	itemKindBWing3rd                                      // 3代翅膀
	itemKindBCapeLord                                     // 王者披风
	itemKindBCapeFighter                                  // 舞者披风
	itemKindBWingMonster                                  // 2.5代翅膀
	itemKindBPendantAttackPhysical                        // 项链物理系(火/风/紫晶/愁怀)
	itemKindBPendantAttackMagic                           // 项链魔法系(雷/冰/水/坚硬)
	itemKindBRing                                         // 指环
	itemKindBRingEvent                                    // 指环活动变身用
	itemKindBRingTransform                                // 变身指环
	itemKindBPremiumPendant                               // 会员项链
	itemKindBPremiumRing                                  // 会员指环
	itemKindBSkillScroll                                  // 技能书
	itemKindBSkillParchment                               // 技能羊皮纸
	itemKindBSkillOrb                                     // 技能石
	itemKindBSkillDarkLord                                // 技能卷轴
	itemKindBSocketSeed                                   // 荧之石
	itemKindBSocketSphere                                 // 光之石
	itemKindBSocketSeedSphere                             // 荧光宝石
	itemKindBPentagram                                    // 元素卷轴
	itemKindBPentagramJewel                               // 艾尔特
	itemKindBHelper                                       // 装备辅助道具(小恶魔/守护天使等)
	itemKindBEquipment               itemKindB = 1 + iota // 背包宠物
	itemKindBPremiumScroll                                // 会员卷轴(加速/防御/愤怒等卷轴)
	itemKindBPremiumElixir                                // 会员圣水(力量/敏捷/体力等圣水)
	itemKindBPremiumPotion                                // 会员活力秘药
	itemKindBPotion                                       // 药水(红/蓝/酒等)
	itemKindBRibbonBox                                    // 圣诞宝箱
	itemKindBEventTicket                                  // 活动门票
	itemKindBTicket                                       // 传承道具兑换券
	itemKindBQuest                                        // 任务物品
	itemKindBJewel                                        // 宝石
	itemKindBJewelBundle                                  // 宝石组合
	itemKindBJewelHarmonyStoneGems                        // 再生原石
	itemKindBJewelHarmonyStoneRefine                      // 进化宝石
	itemKindBWingConqueror                                // 征服者的翅膀
	itemKindBWingAngelDevil          itemKindB = 2 + iota // 善恶的翅膀
	itemKindBMuun                                         // 宠物道具
	itemKindBMuunExchange                                 // 宠物道具
)

type itemWeaponKind int

const (
	itemWeaponKindSwordHandOne itemWeaponKind = iota
	itemWeaponKindSwordHandTwo
	itemWeaponKindMace
	itemWeaponKindSpear
	itemWeaponKindStaffHandOne
	itemWeaponKindStaffHandTwo
	itemWeaponKindShield
	itemWeaponKindBow
	itemWeaponKindCrossbow
	itemWeaponKindSummonerStaff
	itemWeaponKindSummonerBook
	itemWeaponKindSceper
	itemWeaponKindRageFighter
)

type ItemBase struct {
	Index              int       `xml:"Index,attr"`
	Slot               int       `xml:"Slot,attr"`
	SkillIndex         int       `xml:"SkillIndex,attr"`
	TwoHand            bool      `xml:"TwoHand,attr"`
	Width              int       `xml:"Width,attr"`
	Height             int       `xml:"Height,attr"`
	Serial             int       `xml:"Serial,attr"`
	Option             bool      `xml:"Option,attr"`
	Drop               bool      `xml:"Drop,attr"`
	DropLevel          int       `xml:"DropLevel,attr"`
	DamageMin          int       `xml:"DamageMin,attr"`
	DamageMax          int       `xml:"DamageMax,attr"`
	AttackSpeed        int       `xml:"AttackSpeed,attr"`
	WalkSpeed          int       `xml:"WalkSpeed,attr"`
	Defense            int       `xml:"Defense,attr"`
	MagicDefense       int       `xml:"MagicDefense,attr"`
	SuccessfulBlocking int       `xml:"SuccessfulBlocking,attr"`
	Durability         int       `xml:"Durability,attr"`
	MagicDurability    int       `xml:"MagicDurability,attr"`
	IceRes             int       `xml:"IceRes,attr"`
	PoisonRes          int       `xml:"PoisonRes,attr"`
	LightRes           int       `xml:"LightRes,attr"`
	FireRes            int       `xml:"FireRes,attr"`
	EarthRes           int       `xml:"EarthRes,attr"`
	WindRes            int       `xml:"WindRes,attr"`
	WaterRes           int       `xml:"WaterRes,attr"`
	MagicPower         int       `xml:"MagicPower,attr"`
	ReqLevel           int       `xml:"ReqLevel,attr"`
	ReqStrength        int       `xml:"ReqStrength,attr"`
	ReqDexterity       int       `xml:"ReqDexterity,attr"`
	ReqVitality        int       `xml:"ReqVitality,attr"`
	ReqEnergy          int       `xml:"ReqEnergy,attr"`
	ReqCommand         int       `xml:"ReqCommand,attr"`
	Money              int       `xml:"Money,attr"`
	SetAttrib          int       `xml:"SetAttrib,attr"`
	DarkWizard         int       `xml:"DarkWizard,attr"`
	DarkKnight         int       `xml:"DarkKnight,attr"`
	FairyElf           int       `xml:"FairyElf,attr"`
	MagicGladiator     int       `xml:"MagicGladiator,attr"`
	DarkLord           int       `xml:"DarkLord,attr"`
	Summoner           int       `xml:"Summoner,attr"`
	RageFighter        int       `xml:"RageFighter,attr"`
	Type               itemType  `xml:"Type,attr"`
	Dump               bool      `xml:"Dump,attr"`
	Transaction        bool      `xml:"Transaction,attr"`
	PersonalStore      bool      `xml:"PersonalStore,attr"`
	StoreWarehouse     bool      `xml:"StoreWarehouse,attr"`
	SellToNPC          bool      `xml:"SellToNPC,attr"`
	Repair             bool      `xml:"Repair,attr"`
	KindA              itemKindA `xml:"KindA,attr"`
	KindB              itemKindB `xml:"KindB,attr"`
	Overlap            int       `xml:"Overlap,attr"`
	Name               string    `xml:"Name,attr"`
	Annotation         string    `xml:"annotation,attr"`
	ModelPath          string    `xml:"ModelPath,attr"`
	ModelFile          string    `xml:"ModelFile,attr"`
}

type itemList struct {
	Sections []*struct {
		Index string      `xml:"Index,attr"`
		Name  string      `xml:"Name,attr"`
		Items []*ItemBase `xml:"Item"`
	} `xml:"Section"`
}

type itemBaseTable []map[int]*ItemBase

func (table itemBaseTable) GetItemBase(i, j int) *ItemBase {
	return table[i][j]
}

// ItemBaseTable represents item profile table
var ItemBaseTable itemBaseTable

func init() {
	var itemList itemList
	conf.XML(path.Join(conf.PathCommon, "Items/IGC_ItemList.xml"), &itemList)

	// [][]array -> []map
	ItemBaseTable = make(itemBaseTable, len(itemList.Sections))
	for i, section := range itemList.Sections {
		ItemBaseTable[i] = make(map[int]*ItemBase)
		for _, v := range section.Items {
			ItemBaseTable[i][v.Index] = v
		}
	}
}

// Item represents a item
type Item struct {
	*ItemBase
	ID          int  `db:"item_id"`
	BaseSection int  `db:"item_base_section"` // 0 ~ 15
	BaseIndex   int  `db:"item_base_index"`   // 0 ~ 511
	Level       int  `db:"item_level"`
	Durability  int  `db:"item_durability"`
	Lucky       bool `db:"item_lucky"`
	Skill       bool `db:"item_skill"`
	Append      int  `db:"item_append"`
	Excel       int  `db:"item_excel"`
	Set         int  `db:"item_set"`
	Option380   bool `db:"item_option380"`
	Period      int  `db:"item_period"`
	SocketIndex int  `db:"item_socket_index"`
	SocketSlots [5]int
	SocketSlot1 int `db:"item_socket_slot1"`
	SocketSlot2 int `db:"item_socket_slot2"`
	SocketSlot3 int `db:"item_socket_slot3"`
	SocketSlot4 int `db:"item_socket_slot4"`
	SocketSlot5 int `db:"item_socket_slot5"`
}

// NewItem construct a item with section and index
func NewItem(section, index int) *Item {
	return &Item{
		ID:          0,
		BaseSection: section,
		BaseIndex:   index,
		ItemBase:    ItemBaseTable.GetItemBase(section, index),
	}
}

// Marshal marshal item struct to [32]byte variable
// ----------------------------------------------
// field[0]
// bit0~bit7: itembase index, 0~255
// ----------------------------------------------
// field[1]
// bit0~bit1: append attack/defense
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
// bit6: append attack/defense 16 flag
// now field[7].bit6, field[1].bit1 and field[1].bit0 may range as follow:
// 001: append 4
// 010: append 8
// 011: append 12
// 100: append 16
// bit7: extension flag of itembase index , act as the bit8 of field[0]
// now itembase index is 9 bit and range is 0~511
// ----------------------------------------------
// field[8]
// set or ancient
// ----------------------------------------------
// field[9]
// bit0~bit2: period
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
func (item *Item) Marshal() ([]byte, error) {
	return nil, nil
}
