package item

import (
	"fmt"
	"math/rand"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/class"
)

func init() {
	ItemTable.init()
}

func Code(section, index int) int {
	return section*512 + index
}

// ItemTable a map table
var ItemTable itemTable

type itemTable []map[int]*ItemBase

func (table *itemTable) init() {
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
	t := make(itemTable, len(itemList.Sections))
	for i, section := range itemList.Sections {
		t[i] = make(map[int]*ItemBase)
		for _, v := range section.Items {
			v.ReqClass[class.Wizard] = v.DarkWizard
			v.ReqClass[class.Knight] = v.DarkKnight
			v.ReqClass[class.Elf] = v.FairyElf
			v.ReqClass[class.Magumsa] = v.MagicGladiator
			v.ReqClass[class.DarkLord] = v.DarkLord
			v.ReqClass[class.RageFighter] = v.RageFighter
			// v.ReqClass[class.GrowLancer] = v.GrowLancer
			t[i][v.Index] = v
		}
	}
	*table = t
}

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

func (table itemTable) GetItemLevel(i, j, level int) int {
	itBase := table[i][j]
	if !itBase.Drop {
		return -1
	}
	itLevel := itBase.DropLevel
	if i == 13 {
		itLevel = itBase.ReqLevel
	}
	// Orb of Summoning 召唤之石
	if i == 12 && j == 11 {
		if rand.Intn(10) == 0 {
			itLevel = level / 10
			if itLevel > 0 {
				itLevel--
			}
			if itLevel > 6 {
				itLevel = 6
			}
			return itLevel
		}
		return -1
	}
	// Transformation Ring 变身戒指
	if i == 13 && j == 10 {
		if rand.Intn(10) == 0 {
			itLevel = level / 10
			if itLevel > 0 {
				itLevel--
			}
			if itLevel > 5 {
				itLevel = 5
			}
			return itLevel
		}
		return -1
	}
	if i == 14 {
		if j == 15 { // Zen 金
			return -1
		}
		if itLevel >= level-8 && itLevel <= level {
			return 0
		}
		return -1
	}
	if itLevel >= level-18 && itLevel <= level {
		if itBase.KindA == KindACommon {
			return 0
		}
		itLevel = (level - itLevel) / 3
		if itBase.KindA == KindAPendant || itBase.KindA == KindARing {
			if itLevel > 4 {
				itLevel = 4
			}
		}
		return itLevel
	}
	return -1
}

const MaxItemIndex int = 512

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
	KindBSwordKnight                                  // 剑战士用
	KindBSwordMagicGladiator                          // 剑魔剑士用
	KindBSwordRageFighter                             // 剑格斗用
	KindBAxe                                          // 斧
	KindBMace                                         // 槌
	KindBScepter                                      // 权杖
	KindBSpear                                        // 矛
	KindBBow                                          // 弓
	KindBCrossbow                                     // 弩
	KindBArrow                                        // 弓箭
	KindBBolt                                         // 弩箭
	KindBStaffWizard                                  // 法师杖
	KindBStaffSummoner                                // 召唤杖
	KindBBookSummoner                                 // 召唤副手书
	KindBShield                                       // 盾类
	KindBHelmet                                       // 头盔
	KindBArmor                                        // 铠甲
	KindBPants                                        // 护腿
	KindBGloves                                       // 护手
	KindBBoots                                        // 靴子
	KindBWingTalisman                                 // 翅膀合成符咒
	KindBSmallWing                                    // 小翅膀
	KindBWing1st                                      // 1代翅膀
	KindBWing2nd                                      // 2代翅膀
	KindBWing3rd                                      // 3代翅膀
	KindBCapeLord                                     // 王者披风
	KindBCapeFighter                                  // 舞者披风
	KindBWingMonster                                  // 2.5代翅膀
	KindBPendantAttackPhysical                        // 项链物理系(火/风/紫晶/愁怀)
	KindBPendantAttackMagic                           // 项链魔法系(雷/冰/水/坚硬)
	KindBRing                                         // 指环
	KindBRingEvent                                    // 指环活动变身用
	KindBRingTransform                                // 变身指环
	KindBPremiumPendant                               // 会员项链
	KindBPremiumRing                                  // 会员指环
	KindBSkillScroll                                  // 技能书
	KindBSkillParchment                               // 技能羊皮纸
	KindBSkillOrb                                     // 技能石
	KindBSkillDarkLord                                // 技能卷轴
	KindBSocketSeed                                   // 荧之石
	KindBSocketSphere                                 // 光之石
	KindBSocketSeedSphere                             // 荧光宝石
	KindBPentagram                                    // 元素卷轴
	KindBPentagramJewel                               // 艾尔特
	KindBHelper                                       // 装备辅助道具(小恶魔/守护天使等)
	KindBEquipment               itemKindB = 1 + iota // 背包宠物
	KindBPremiumScroll                                // 会员卷轴(加速/防御/愤怒等卷轴)
	KindBPremiumElixir                                // 会员圣水(力量/敏捷/体力等圣水)
	KindBPremiumPotion                                // 会员活力秘药
	KindBPotion                                       // 药水(红/蓝/酒等)
	KindBRibbonBox                                    // 圣诞宝箱
	KindBEventTicket                                  // 活动门票
	KindBTicket                                       // 传承道具兑换券
	KindBQuest                                        // 任务物品
	KindBJewel                                        // 宝石
	KindBJewelBundle                                  // 宝石组合
	KindBJewelHarmonyStoneGems                        // 再生原石
	KindBJewelHarmonyStoneRefine                      // 进化宝石
	KindBWingConqueror                                // 征服者的翅膀
	KindBWingAngelDevil          itemKindB = 2 + iota // 善恶的翅膀
	KindBMuun                                         // 宠物道具
	KindBMuunExchange                                 // 宠物道具
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
