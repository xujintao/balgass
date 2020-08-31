package item

type itemType int

const (
	TypeCommon itemType = iota
	TypeRegular
	TypeSocket
	Type380
	TypeLucky
	TypeEvent
	TypeArchangel
	TypeChaos // 玛雅龙斧 玛雅神弓 玛雅雷杖
)

type itemKindA int

const (
	KindACommon itemKindA = iota
	KindAWeapon
	KindAPendant
	KindAArmor
	KindARing
	KindAHelper // imp etc.
	KindAWing
	KindAPremiumBuff
	KindAPentagram
	KindAJewel
	KindASkill
	KindAEvent
	KindAMuun
)

type itemKindB int

const (
	KindBCommon                  itemKindB = iota
	kindBSwordKnight                                  // 剑战士用
	kindBSwordMagicGladiator                          // 剑魔剑士用
	kindBSwordRageFighter                             // 剑格斗用
	kindBAxe                                          // 斧
	kindBMace                                         // 槌
	kindBScepter                                      // 权杖
	kindBSpear                                        // 矛
	kindBBow                                          // 弓
	kindBCrossbow                                     // 弩
	kindBArrow                                        // 弓箭
	kindBBolt                                         // 弩箭
	kindBStaffWizard                                  // 法师杖
	kindBStaffSummoner                                // 召唤杖
	kindBBookSummoner                                 // 召唤副手书
	kindBShield                                       // 盾类
	kindBHelmet                                       // 头盔
	kindBArmor                                        // 铠甲
	kindBPants                                        // 护腿
	kindBGloves                                       // 护手
	kindBBoots                                        // 靴子
	kindBWingTalisman                                 // 翅膀合成符咒
	kindBSmallWing                                    // 小翅膀
	kindBWing1st                                      // 1代翅膀
	kindBWing2nd                                      // 2代翅膀
	kindBWing3rd                                      // 3代翅膀
	kindBCapeLord                                     // 王者披风
	kindBCapeFighter                                  // 舞者披风
	kindBWingMonster                                  // 2.5代翅膀
	kindBPendantAttackPhysical                        // 项链物理系(火/风/紫晶/愁怀)
	kindBPendantAttackMagic                           // 项链魔法系(雷/冰/水/坚硬)
	kindBRing                                         // 指环
	kindBRingEvent                                    // 指环活动变身用
	kindBRingTransform                                // 变身指环
	kindBPremiumPendant                               // 会员项链
	kindBPremiumRing                                  // 会员指环
	kindBSkillScroll                                  // 技能书
	kindBSkillParchment                               // 技能羊皮纸
	kindBSkillOrb                                     // 技能石
	kindBSkillDarkLord                                // 技能卷轴
	kindBSocketSeed                                   // 荧之石
	kindBSocketSphere                                 // 光之石
	kindBSocketSeedSphere                             // 荧光宝石
	kindBPentagram                                    // 元素卷轴
	kindBPentagramJewel                               // 艾尔特
	kindBHelper                                       // 装备辅助道具(小恶魔/守护天使等)
	kindBEquipment               itemKindB = 1 + iota // 背包宠物
	kindBPremiumScroll                                // 会员卷轴(加速/防御/愤怒等卷轴)
	kindBPremiumElixir                                // 会员圣水(力量/敏捷/体力等圣水)
	kindBPremiumPotion                                // 会员活力秘药
	kindBPotion                                       // 药水(红/蓝/酒等)
	kindBRibbonBox                                    // 圣诞宝箱
	kindBEventTicket                                  // 活动门票
	kindBTicket                                       // 传承道具兑换券
	kindBQuest                                        // 任务物品
	kindBJewel                                        // 宝石
	kindBJewelBundle                                  // 宝石组合
	kindBJewelHarmonyStoneGems                        // 再生原石
	kindBJewelHarmonyStoneRefine                      // 进化宝石
	kindBWingConqueror                                // 征服者的翅膀
	kindBWingAngelDevil          itemKindB = 2 + iota // 善恶的翅膀
	kindBMuun                                         // 宠物道具
	kindBMuunExchange                                 // 宠物道具
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
	SellToNPC          int       `xml:"SellToNPC,attr"`
	Repair             bool      `xml:"Repair,attr"`
	KindA              itemKindA `xml:"KindA,attr"`
	KindB              itemKindB `xml:"KindB,attr"`
	Overlap            int       `xml:"Overlap,attr"`
	Name               string    `xml:"Name,attr"`
	Annotation         string    `xml:"annotation,attr"`
	ModelPath          string    `xml:"ModelPath,attr"`
	ModelFile          string    `xml:"ModelFile,attr"`
}

// Item represents a item
type Item struct {
	*ItemBase
	ID             int               `db:"item_id"`           // serial
	Section        int               `db:"item_base_section"` // 0 ~ 15
	Index          int               `db:"item_base_index"`   // 0 ~ 511
	Code           int               `db:"-"`                 // section*512 + index
	Level          int               `db:"item_level"`
	Durability     int               `db:"item_durability"`
	Lucky          bool              `db:"item_lucky"`
	Skill          bool              `db:"item_skill"`
	Addition       int               `db:"item_addition"`
	Excel          int               `db:"item_excel"`
	Set            int               `db:"item_set"`
	Option380      bool              `db:"item_option380"`
	Period         int               `db:"item_period"`
	HarmonyEffect  harmonyEffectKind `db:"item_harmony_effect"`
	HarmonyLevel   int               `db:"item_harmony_level"`
	PentagramBonus int               `db:"item_pentagram_bonus"`
	MuunRank       int               `db:"item_muun_rank"`
	SocketBonus    int               `db:"item_socket_bonus"`
	SocketSlots    [5]int            `db:"-"` // slot array
	SocketSlot1    int               `db:"item_socket_slot1"`
	SocketSlot2    int               `db:"item_socket_slot2"`
	SocketSlot3    int               `db:"item_socket_slot3"`
	SocketSlot4    int               `db:"item_socket_slot4"`
	SocketSlot5    int               `db:"item_socket_slot5"`
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
func (item *Item) Marshal() ([]byte, error) {
	return nil, nil
}
