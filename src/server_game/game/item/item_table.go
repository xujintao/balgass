package item

import (
	"fmt"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/class"
)

func init() {
	type itemListConfig struct {
		Sections []*struct {
			Index string      `xml:"Index,attr"`
			Name  string      `xml:"Name,attr"`
			Items []*ItemBase `xml:"Item"`
		} `xml:"Section"`
	}
	var itemList itemListConfig
	conf.XML(conf.PathCommon, "Items/IGC_ItemList.xml", &itemList)

	// [][]array -> []map
	ItemTable = make(itemTable, len(itemList.Sections))
	for i, section := range itemList.Sections {
		ItemTable[i] = make(map[int]*ItemBase)
		for _, v := range section.Items {
			v.ReqClass[class.Wizard] = v.DarkWizard
			v.ReqClass[class.Knight] = v.DarkKnight
			v.ReqClass[class.Elf] = v.FairyElf
			v.ReqClass[class.Magumsa] = v.MagicGladiator
			v.ReqClass[class.DarkLord] = v.DarkLord
			v.ReqClass[class.RageFighter] = v.RageFighter
			// v.ReqClass[class.GrowLancer] = v.GrowLancer
			ItemTable[i][v.Index] = v
		}
	}
}

func Code(section, index int) int {
	return section*512 + index
}

// ItemTable a map table
var ItemTable itemTable

type itemTable []map[int]*ItemBase

func (table itemTable) GetItemBase(i, j int) (*ItemBase, error) {
	if i >= len(table) {
		return nil, fmt.Errorf("item section over bound")
	}
	items := table[i]
	item, ok := items[j]
	if !ok {
		return nil, fmt.Errorf("item index over bound")
	}
	return item, nil
}

func (table itemTable) GetItemBaseMust(i, j int) *ItemBase {
	return table[i][j]
}

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
	ReqClass           [8]int    `xml:"-"`
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
