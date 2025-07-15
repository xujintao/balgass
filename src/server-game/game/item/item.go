package item

import (
	"database/sql/driver"
	"encoding/json"
	"log/slog"
	"math"

	"github.com/xujintao/balgass/src/server-game/conf"
)

// Item represents a item
type Item struct {
	*ItemBase                    `json:"-"`
	Position                     int    `json:"position"`
	ID                           int    `json:"id"`      // serial
	Section                      int    `json:"section"` // 0 ~ 15
	Index                        int    `json:"index"`   // 0 ~ 511
	Code                         int    `json:"-"`       // section*512 + index
	Level                        int    `json:"level"`
	Durability                   int    `json:"durability"`
	Lucky                        bool   `json:"lucky,omitempty"`
	Skill                        bool   `json:"skill,omitempty"`
	Addition                     int    `json:"addition,omitempty"`                  // 0/4/8/12/16
	ExcellentAttackRate          bool   `json:"excellent_attack_rate,omitempty"`     // bit5:卓越攻击几率10%
	ExcellentAttackLevel         bool   `json:"excellent_attack_level,omitempty"`    // bit4:攻击力增加等级/20
	ExcellentAttackPercent       bool   `json:"excellent_attack_percent,omitempty"`  // bit3:攻击力增加2%
	ExcellentAttackSpeed         bool   `json:"excellent_attack_speed,omitempty"`    // bit2:攻击(魔法)速度增加7
	ExcellentAttackHP            bool   `json:"excellent_attack_hp,omitempty"`       // bit1:杀死怪物时所获生命值增加生命值/8
	ExcellentAttackMP            bool   `json:"excellent_attack_mp,omitempty"`       // bit0:杀死怪物时所获魔法值增加魔法值/8
	ExcellentDefenseHP           bool   `json:"excellent_defense_hp,omitempty"`      // bit5:最大生命值增加4%
	ExcellentDefenseMP           bool   `json:"excellent_defense_mp,omitempty"`      // bit4:最大魔法值增加4%
	ExcellentDefenseReduce       bool   `json:"excellent_defense_reduce,omitempty"`  // bit3:伤害减少4%
	ExcellentDefenseReflect      bool   `json:"excellent_defense_reflect,omitempty"` // bit2:伤害反射5%
	ExcellentDefenseRate         bool   `json:"excellent_defense_rate,omitempty"`    // bit1:防御成功率10%
	ExcellentDefenseMoney        bool   `json:"excellent_defense_money,omitempty"`   // bit0:杀死怪物时所获金币增加30%
	WingAdditionAttack           bool   `json:"wing_addition_attack,omitempty"`
	WingAdditionMagicAttack      bool   `json:"wing_addition_magic_attack,omitempty"`
	WingAdditionCurseAttack      bool   `json:"wing_addition_curse_attack,omitempty"`
	WingAdditionDefense          bool   `json:"wing_addition_defense,omitempty"`
	WingAdditionRecoverHP        bool   `json:"wing_addition_recover_hp,omitempty"`
	ExcellentWing2Speed          bool   `json:"excellent_wing2_speed,omitempty"`      // 2D:攻击(魔法)速度增加5
	ExcellentWing2AG             bool   `json:"excellent_wing2_ag,omitempty"`         // 2D:最大AG增加50
	ExcellentWing2Leadership     bool   `json:"excellent_wing2_leadership,omitempty"` // 2D:声望增加10+5*level
	ExcellentWing2Ignore         bool   `json:"excellent_wing2_ignore,omitempty"`     // 2D:无视防御伤害几率3%
	ExcellentWing2MP             bool   `json:"excellent_wing2_mp,omitempty"`         // 2D:魔法值增加50+5*level
	ExcellentWing2HP             bool   `json:"excellent_wing2_hp,omitempty"`         // 2D:生命值增加50+5*level
	ExcellentWing3MP             bool   `json:"excellent_wing3_mp,omitempty"`         // 3D:魔法值完全恢复几率5%
	ExcellentWing3HP             bool   `json:"excellent_wing3_hp,omitempty"`         // 3D:生命值完全恢复几率5%
	ExcellentWing3Return         bool   `json:"excellent_wing3_return,omitempty"`     // 3D:反弹攻击力几率5%
	ExcellentWing3Ignore         bool   `json:"excellent_wing3_ignore,omitempty"`     // 3D:无视防御伤害几率5%
	ExcellentWing25HP            bool   `json:"excellent_wing25_hp,omitempty"`        // 2.5D:生命值增加5%
	ExcellentWing25Ignore        bool   `json:"excellent_wing25_ignore,omitempty"`    // 2.5D:无视防御伤害几率3%
	Set                          int    `json:"set,omitempty"`
	Option380                    bool   `json:"option380,omitempty"`
	Period                       int    `json:"period,omitempty"`
	HarmonyEffect                int    `json:"harmony_effect,omitempty"`
	HarmonyLevel                 int    `json:"harmony_level,omitempty"`
	HarmonyOption                int    `json:"harmony_option,omitempty"`
	PentagramBonus               int    `json:"pentagram_bonus,omitempty"`
	MuunRank                     int    `json:"muun_rank,omitempty"`
	SocketBonus                  int    `json:"socket_bonus,omitempty"`
	SocketSlots                  [5]int `json:"-"` // slot array
	SocketSlot1                  int    `json:"socket_slot1,omitempty"`
	SocketSlot2                  int    `json:"socket_slot2,omitempty"`
	SocketSlot3                  int    `json:"socket_slot3,omitempty"`
	SocketSlot4                  int    `json:"socket_slot4,omitempty"`
	SocketSlot5                  int    `json:"socket_slot5,omitempty"`
	MaxDurability                int    `json:"-"`
	AttackMin                    int    `json:"-"`
	AttackMax                    int    `json:"-"`
	Magic                        int    `json:"-"`
	Defense                      int    `json:"-"`
	DefenseRate                  int    `json:"-"`
	AdditionAttack               int    `json:"-"`
	AdditionMagicAttack          int    `json:"-"`
	AdditionCurseAttack          int    `json:"-"`
	AdditionDefense              int    `json:"-"`
	AdditionDefenseRate          int    `json:"-"`
	AdditionRecoverHP            int    `json:"-"`
	AdditionAbsorbDamagePercent5 bool   `json:"-"`
	AdditionAG50                 bool   `json:"-"`
	AdditionSpeed5               bool   `json:"-"`
	BuyMoney                     int    `json:"-"`
	SellMoney                    int    `json:"-"`
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

func (it *Item) IsExcellent() bool {
	return it.ExcellentAttackRate ||
		it.ExcellentAttackLevel ||
		it.ExcellentAttackPercent ||
		it.ExcellentAttackSpeed ||
		it.ExcellentAttackHP ||
		it.ExcellentAttackMP ||
		it.ExcellentDefenseHP ||
		it.ExcellentDefenseMP ||
		it.ExcellentDefenseReduce ||
		it.ExcellentDefenseReflect ||
		it.ExcellentDefenseRate ||
		it.ExcellentDefenseMoney ||
		it.ExcellentWing2Speed ||
		it.ExcellentWing2AG ||
		it.ExcellentWing2Leadership ||
		it.ExcellentWing2Ignore ||
		it.ExcellentWing2MP ||
		it.ExcellentWing2HP ||
		it.ExcellentWing3MP ||
		it.ExcellentWing3HP ||
		it.ExcellentWing3Return ||
		it.ExcellentWing3Ignore ||
		it.ExcellentWing25HP ||
		it.ExcellentWing25Ignore
}

func (it *Item) ExcellentCount() int {
	flags := []bool{
		it.ExcellentAttackRate,
		it.ExcellentAttackLevel,
		it.ExcellentAttackPercent,
		it.ExcellentAttackSpeed,
		it.ExcellentAttackHP,
		it.ExcellentAttackMP,
		it.ExcellentDefenseHP,
		it.ExcellentDefenseMP,
		it.ExcellentDefenseReduce,
		it.ExcellentDefenseReflect,
		it.ExcellentDefenseRate,
		it.ExcellentDefenseMoney,
		it.ExcellentWing2Speed,
		it.ExcellentWing2AG,
		it.ExcellentWing2Leadership,
		it.ExcellentWing2Ignore,
		it.ExcellentWing2MP,
		it.ExcellentWing2HP,
		it.ExcellentWing3MP,
		it.ExcellentWing3HP,
		it.ExcellentWing3Return,
		it.ExcellentWing3Ignore,
		it.ExcellentWing25HP,
		it.ExcellentWing25Ignore,
	}
	count := 0
	for _, flag := range flags {
		if flag {
			count++
		}
	}
	return count
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

func (it *Item) CalculateMoney() {
	money := 0
	switch {
	// use the money of simply items
	case (it.Section == 12 || it.Section == 15) && it.ItemBase.Money != 0:
		money = it.ItemBase.Money
	// calculate the money of special items
	case it.KindA == KindAWeapon && it.Code != Code(0, 41), // weapons exclude Pandora Pick (two-handed)
		it.KindA == KindAArmor,
		it.KindA == KindAWing,
		it.KindA == KindAPendant,
		it.KindA == KindARing:
		level := it.ItemBase.DropLevel + it.Level*3
		if it.IsExcellent() && it.KindA != KindAWing {
			level += 25
		}
		levelTable := []int{
			0,   // 0
			0,   // 1
			0,   // 2
			0,   // 3
			0,   // 4
			4,   // 5
			10,  // 6
			25,  // 7
			45,  // 8
			65,  // 9
			95,  // 10
			135, // 11
			185, // 12
			245, // 13
			305, // 14
			365, // 15
		}
		level += levelTable[it.Level]
		// 1. calculate level value
		switch it.KindA {
		case KindAWeapon, KindAArmor:
			money = 100 + level*level*(level+40)/8
			if it.KindA == KindAWeapon && !it.ItemBase.TwoHand {
				money = money * 80 / 100
			}
		case KindAPendant, KindARing:
			money = 100 + level*level*level
		case KindAWing:
			money = 40000000 + 11*level*level*(level+40)
		}
		// 2. calculate skill value
		if it.Skill {
			money += int(float64(money) * 1.5)
		}
		// 3. calculate lucky value
		if it.Lucky {
			money += int(float64(money) * 25.0 / 100.0)
		}
		// 4. calculate addition value
		switch it.KindA {
		case KindAWeapon, KindAArmor, KindAWing:
		switch it.Addition {
		case 4:
			money += int(float64(money) * 6.0 / 10.0)
		case 8:
			money += int(float64(money) * 14.0 / 10.0)
		case 12:
			money += int(float64(money) * 28.0 / 10.0)
		case 16:
			money += int(float64(money) * 56.0 / 10.0)
			}
		case KindAPendant, KindARing:
			money += money * (it.Addition / 4)
		}
		// 5. calculate excellent value
		if it.IsExcellent() {
			switch it.KindA {
			case KindAWeapon, KindAArmor:
				for i := it.ExcellentCount(); i > 0; i-- {
					money += money
				}
			case KindAPendant, KindARing:
				// nothing
			case KindAWing:
				for i := it.ExcellentCount(); i > 0; i-- {
					money += int(float64(money) * 25.0 / 100.0)
				}
			}
		}
		if it.Type == Type380 {
			money += int(float64(money) * 16.0 / 100.0)
		}
		if it.Type == TypeSocket {
			// todo
		}
	case it.KindA == KindAPentagram &&
		it.KindB == KindBPentagram: // Scroll 卷轴
		// todo
	case it.KindA == KindAPentagram &&
		it.KindB == KindBPentagramJewel: // Errtel 艾尔特
		// todo
	case it.Code == Code(0, 41): // Pandora Pick (two-handed) 潘多拉锄头(双手)
		money = 100000
	case it.Code == Code(4, 7): // Bolt 弩箭
		if it.Durability > 0 {
			levelVaule := 100
			switch it.Level {
			case 1:
				levelVaule = 1400
			case 2:
				levelVaule = 2200
			case 3:
				levelVaule = 3000
			}
			money = levelVaule * it.Durability / it.ItemBase.Durability
		}
	case it.Code == Code(4, 15): // Arrow 弓箭
		if it.Durability > 0 {
			levelVaule := 70
			switch it.Level {
			case 1:
				levelVaule = 1200
			case 2:
				levelVaule = 2000
			case 3:
				levelVaule = 2800
			}
			money = levelVaule * it.Durability / it.ItemBase.Durability
		}
	case it.Code == Code(14, 13): // Jewel of Bless 祝福宝石
		money = conf.Price.Value.JewelOfBlessPrice
	case it.Code == Code(12, 30): // Jewel of Bless Bundle 祝福宝石组合
		money = conf.Price.Value.JewelOfBlessPrice * 10 * (it.Level + 1)
	case it.Code == Code(14, 14): // Jewel of Soul 灵魂宝石
		money = conf.Price.Value.JewelOfSoulPrice
	case it.Code == Code(12, 31): // Jewel of Bless Bundle 灵魂宝石组合
		money = conf.Price.Value.JewelOfSoulPrice * 10 * (it.Level + 1)
	case it.Code == Code(14, 15): // Jewel of Chaos 玛雅之石
		money = conf.Price.Value.JewelOfChaosPrice
	case it.Code == Code(12, 141): // Jewel of Chaos Bundle 玛雅之石组合
		money = conf.Price.Value.JewelOfChaosPrice * 10 * (it.Level + 1)
	case it.Code == Code(14, 16): // Jewel of Life 生命宝石
		money = conf.Price.Value.JewelOfLifePrice
	case it.Code == Code(12, 136): // Jewel of Life Bundle 生命宝石组合
		money = conf.Price.Value.JewelOfLifePrice * 10 * (it.Level + 1)
	case it.Code == Code(14, 22): // Jewel of Creation 创造宝石
		money = conf.Price.Value.JewelOfCreationPrice
	case it.Code == Code(12, 137): // Jewel of Creation Bundle 创造宝石组合
		money = conf.Price.Value.JewelOfCreationPrice * 10 * (it.Level + 1)
	case it.Code == Code(14, 31): // Jewel of Guardian 守护宝石
		money = conf.Price.Value.JewelOfGuardianPrice
	case it.Code == Code(12, 138): // Jewel of Guardian Bundle 守护宝石组合
		money = conf.Price.Value.JewelOfGuardianPrice * 10 * (it.Level + 1)
	case it.Code == Code(14, 41), // Gemstone 再生原石
		it.Code == Code(14, 42), // Jewel of Harmony 再生宝石
		it.Code == Code(14, 43), // Lower refining stone 初级进化宝石
		it.Code == Code(14, 44): // Higher refining stone 高级进化宝石
		money = 18600
	case it.Code == Code(12, 139), // Gemstone Bundle 再生原石组合
		it.Code == Code(12, 140), // Jewel of Harmony Bundle 再生宝石组合
		it.Code == Code(12, 142), // Lower Refining Stone Bundle 低级进化宝石组合
		it.Code == Code(12, 143): // Higher Refining Stone Bundle 高级进化宝石组合
		money = 186000 * (it.Level + 1)
	// -----------------------------
	case it.Code == Code(12, 144), // Mithril Fragment 元素之心碎片
		it.Code == Code(12, 146): // Elixir Fragment 艾丽亚之环碎片
		money = 300000 * it.Durability
	case it.Code == Code(13, 3): // Horn of Dinorant 彩云兽
		money = 960000
		if it.AdditionAbsorbDamagePercent5 {
			money += 300000
		}
		if it.AdditionAG50 {
			money += 300000
		}
		if it.AdditionSpeed5 {
			money += 300000
		}
	case it.Code == Code(13, 7): // Contract (Summon) 召唤佣兵
		switch it.Level {
		case 0:
			money = 1500000
		case 1:
			money = 1200000
		}
	case it.Code == Code(13, 11): // Order (Guardian/Life Stone) 复活石
		if it.Level == 1 {
			money = 2400000
		}
	case it.Code == Code(13, 14): // Loch's Feather 洛克之羽
		switch it.Level {
		case 0:
			money = conf.Price.Value.LochFeatherPrice // 洛克之羽
		case 1:
			money = conf.Price.Value.CrestOfMonarchPrice // 国王卷轴
		}
	case it.Code == Code(13, 15): // Fruits 果实
		money = 9000000
	case it.Code == Code(13, 16), // Scroll of Archangel 血灵之书
		it.Code == Code(13, 17): // Blood Bone 血灵之骷
		levelMoney := [9]int{
			1,       // 0
			10000,   // 1
			50000,   // 2
			100000,  // 3
			300000,  // 4
			500000,  // 5
			800000,  // 6
			1000000, // 7
			1200000, // 8
		}
		money = levelMoney[it.Level]
	case it.Code == Code(13, 18): // Invisibility Cloak 透明披风
		levelMoney := [9]int{
			1,       // 0
			150000,  // 1
			660000,  // 2
			720000,  // 3
			780000,  // 4
			840000,  // 5
			900000,  // 6
			960000,  // 7
			1020000, // 8
		}
		money = levelMoney[it.Level]
	case it.Code == Code(13, 20): // Wizards Ring 魔法戒指
		money = 3000
	case it.Code == Code(13, 29): // Armor of Guardsman 卫兵铠甲
		money = 5000
	case it.Code == Code(13, 31): // Spirit 兽之灵魂
		money = 9000000
	case it.Code == Code(13, 32): // Splinter of Armor 破烂的铠甲片
		money = 150 * it.Durability
	case it.Code == Code(13, 33): // Bless of Guardian 女神的灵智
		money = 300 * it.Durability
	case it.Code == Code(13, 34): // Claw of Beast 猛兽的脚甲
		money = 3000 * it.Durability
	case it.Code == Code(13, 35): // Fragment of Horn 碎角片
		money = 30000
	case it.Code == Code(13, 36): // Broken Horn 折断的角
		money = 90000
	case it.Code == Code(13, 37): // Horn of Fenrir 炎狼兽之角
		money = 150000
	case it.Code == Code(13, 49), // Old Scroll 旧卷纸
		it.Code == Code(13, 50), // Illusion Sorcerer Covenant 幻影教药水
		it.Code == Code(13, 51): // Scroll of Blood 血的卷纸
		levelMoney := [7]int{
			1,       // 0
			500000,  // 1
			600000,  // 2
			800000,  // 3
			1000000, // 4
			1200000, // 5
			1400000, // 6
		}
		money = levelMoney[it.Level]
	case it.Code == Code(13, 52), // Condor Flame 神鹰火种
		it.Code == Code(13, 53): // Condor Feather 神鹰之羽
		money = 3000000
	case it.Code == Code(13, 64), // Demon 强化恶魔
		it.Code == Code(13, 65),  // Spirit of Guardian 强化天使
		it.Code == Code(13, 76),  // Panda Ring 熊猫变身指环
		it.Code == Code(13, 77),  // Brown Panda Ring 棕色熊变身指环
		it.Code == Code(13, 78),  // Pink Panda Ring 粉红色熊变身指环
		it.Code == Code(13, 80),  // Pet Panda 熊猫
		it.Code == Code(13, 109), // Sapphire Ring 蓝宝石戒指
		it.Code == Code(13, 110), // Ruby Ring 红宝石戒指
		it.Code == Code(13, 111), // Topaz Ring 黄宝石戒指
		it.Code == Code(13, 112), // Amethyst Ring 紫水晶戒指
		it.Code == Code(13, 113), // Ruby Necklace 红宝石项链
		it.Code == Code(13, 114), // Emerald Necklace 绿宝石项链
		it.Code == Code(13, 115): // Sapphire Necklace 蓝宝石项链
		money = 3000
	case it.Code == Code(13, 122), // Skeleton Transformation Ring 骷髅变身戒指
		it.Code == Code(13, 123): // Pet Skeleton 幼龙骨架
		money = 6000
	case it.Code == Code(13, 145): // Spirit Map Fragment 精灵地图碎片
		money = 30000 * it.Durability
	case it.Code == Code(13, 146): // Spirit Map 精灵地图
		money = 600000
	case it.Code == Code(13, 147): // Trophies of Battle 战场战利品
		money = 300000 * it.Durability
	case it.Code == Code(13, 163), // Robot Knight Ring 机器人变身指环
		it.Code == Code(13, 164), // Mini Robot Ring 迷你骑士变身指环
		it.Code == Code(13, 165): // Great Heavenly Mage Ring 齐天大圣变身指环
		money = 6000
	case it.Code == Code(13, 169), // Decoration Ring 新年指环
		it.Code == Code(13, 170): // Blessed Decoration Ring 祝福新年指环
		money = 90000
	case it.Code == Code(13, 239): // Muun Energy Converter 宠物智力变换器
		money = 900000
	case it.Code == Code(14, 0), // Apple 苹果
		it.Code == Code(14, 1), // Small Healing Potion 小瓶治疗药水
		it.Code == Code(14, 2), // Healing Potion 中瓶治疗药水
		it.Code == Code(14, 3), // Large Healing Potion 大瓶治疗药水
		it.Code == Code(14, 4), // Small Mana Potion 小瓶魔力药水
		it.Code == Code(14, 5), // Mana Potion 中瓶魔力药水
		it.Code == Code(14, 6), // Large Mana Potion 大瓶魔力药水
		it.Code == Code(14, 8): // Antidote 解毒剂
		money += it.ItemBase.Money * it.ItemBase.Money * 10 / 12
		if it.Level > 0 {
			money *= int(math.Pow(2.0, float64(it.Level)))
		}
		money = money / 10 * 10
		money *= it.Durability
	case it.Code == Code(14, 7): // Siege Potion 攻城药水
		money = 900000 * it.Durability
	case it.Code == Code(14, 9): // Ale 酒
		money = 750
	case it.Code == Code(14, 17): // Devil's Eye 恶魔之眼
		levelMoney := [8]int{
			1,       // 0
			10000,   // 1
			50000,   // 2
			100000,  // 3
			300000,  // 4
			500000,  // 5
			800000,  // 6
			1000000, // 7
		}
		money = levelMoney[it.Level]
	case it.Code == Code(14, 18): // Devil's Key 恶魔之钥
		levelMoney := [8]int{
			1,       // 0
			15000,   // 1
			75000,   // 2
			150000,  // 3
			450000,  // 4
			750000,  // 5
			1200000, // 6
			1500000, // 7
		}
		money = levelMoney[it.Level]
	case it.Code == Code(14, 19): // Devil's Invitation 恶魔广场通行证
		levelMoney := [8]int{
			1,      // 0
			60000,  // 1
			84000,  // 2
			120000, // 3
			180000, // 4
			240000, // 5
			300000, // 6
			360000, // 7
		}
		money = levelMoney[it.Level]
	case it.Code == Code(14, 20): // Remedy of Love 爱情的魔力
		money = 900
	case it.Code == Code(14, 23), // Scroll of the Emperor 帝王之书
		it.Code == Code(14, 24), // Broken Sword 断魂之剑
		it.Code == Code(14, 25), // Tear of Elf 精灵之泪
		it.Code == Code(14, 26): // Soul Shard of Wizard 先知之魂
		money = 9000
	case it.Code == Code(14, 28): // Lost Map 失落的地图
		money = 210000
	case it.Code == Code(14, 29): // Symbol of Kundun 昆顿印记
		money = 30000 * it.Durability
	case it.Code == Code(14, 35): // Small SD Potion 小防护药水
		money = 2000 * it.Durability
	case it.Code == Code(14, 36): // SD Potion 中防护药水
		money = 3900 * it.Durability
	case it.Code == Code(14, 37): // Large SD Potion 大防护药水
		money = 6000 * it.Durability
	case it.Code == Code(14, 38): // Small Complex Potion 小生命圣水
		money = 2500 * it.Durability
	case it.Code == Code(14, 39): // Complex Potion 中生命圣水
		money = 4800 * it.Durability
	case it.Code == Code(14, 40): // Large Complex Potion 大生命圣水
		money = 7500 * it.Durability
	case it.Code == Code(14, 45), // Pumpkin of Luck 幸运南瓜
		it.Code == Code(14, 46), // Jack O'Lantern Blessings 南瓜灯的祝福
		it.Code == Code(14, 47), // Jack O'Lantern Wrath 南瓜灯的愤怒
		it.Code == Code(14, 48), // Jack O'Lantern Cry 南瓜灯的呐喊
		it.Code == Code(14, 49), // Jack O'Lantern Food 南瓜灯的食物
		it.Code == Code(14, 50): // Jack O'Lantern Drink 南瓜灯的饮料
		money = 150 * it.Durability
	case it.Code == Code(14, 51), // Chistmas Star 圣诞之星
		it.Code == Code(14, 63): // Firecracker 爆竹
		money = 200000
	case it.Code == Code(14, 65), // Death-beam Knight Flame 炽炎魔的火种
		it.Code == Code(14, 66), // Hell-Miner Horn 丛林召唤者的角
		it.Code == Code(14, 67), // Dark Phoenix Feather 天魔菲尼斯的羽毛
		it.Code == Code(14, 68): // Abyssal Eye 深渊之眼
		money = 9000
	case it.Code == Code(14, 85), // Cherry Blossom Wine 樱花酒
		it.Code == Code(14, 86),  // Cherry Blossom Rice Cake 樱花饼
		it.Code == Code(14, 87),  // Cherry Blossom Flower Petal 樱花花瓣
		it.Code == Code(14, 90),  // Golden Cherry Blossom Branch 黄色樱花树枝
		it.Code == Code(14, 100): // Lucky Coin 幸运铜币
		money = 300 * it.Durability
	case it.Code == Code(14, 101), // Suspicious Scrap of Paper 奇怪的纸条
		it.Code == Code(14, 102), // Gaion's Order 凯文的指令
		it.Code == Code(14, 103), // First Secromicon Fragment 希克拉碎片1
		it.Code == Code(14, 104), // Second Secromicon Fragment 希克拉碎片2
		it.Code == Code(14, 105), // Third Secromicon Fragment 希克拉碎片3
		it.Code == Code(14, 106), // Fourth Secromicon Fragment 希克拉碎片4
		it.Code == Code(14, 107), // Fifth Secromicon Fragment 希克拉碎片5
		it.Code == Code(14, 108), // Sixth Secromicon Fragment 希克拉碎片6
		it.Code == Code(14, 109), // Complete Secromicon 完整的希克拉
		it.Code == Code(14, 110): // Sign of Dimensions 次元标识
		money = 30000 * it.Durability
	case it.Code == Code(14, 111): // Mirror of Dimensions 次元魔镜
		money = 600000 * it.Durability
	case it.Code == Code(14, 121), // Sealed Golden Box 封印的金箱子
		it.Code == Code(14, 122): // Sealed Silver Box 封印的银箱子
		money = 3000
	case it.Code == Code(14, 141): // Shining Jewellery Case 闪耀的宝石盒
		money = 672000
	case it.Code == Code(14, 142): // Elegant Jewellery Case 精炼的宝石盒
		money = 546000
	case it.Code == Code(14, 143): // Steel Jewellery Case 铁制宝石盒
		money = 471000
	case it.Code == Code(14, 144): // Old Jewellery Case 陈旧的宝石盒
		money = 363000
	case it.Code == Code(14, 217), // Titan's Anger 泰坦的愤怒
		it.Code == Code(14, 218), // Tantalose's Punishment 破坏骑士的刑罚
		it.Code == Code(14, 219), // Erohim's Nightmare 炼狱魔王的噩梦
		it.Code == Code(14, 220), // Hell Maine's Insanity 丛林召唤者的疯狂
		it.Code == Code(14, 221): // Kundun's Greed 昆顿的贪欲
		money = 300000
	case it.Code == Code(14, 234): // Monster Summoning Scroll 怪物召唤书
		money = 3000000
	}
	// buy
	buyMoney := money
	if buyMoney >= 1000 {
		buyMoney = buyMoney / 100 * 100
	} else if buyMoney >= 100 {
		buyMoney = buyMoney / 10 * 10
	}
	it.BuyMoney = buyMoney
	// sell
	sellMoney := money / conf.Price.Value.ItemSellPriceDivisor
	if (it.KindA == KindAWeapon ||
		it.KindA == KindAArmor ||
		it.KindA == KindAWing ||
		it.KindA == KindAHelper ||
		it.KindA == KindAPendant ||
		it.KindA == KindARing) &&
		it.MaxDurability != 0 {
		// discount
		dur := float64(it.Durability) / float64(it.MaxDurability)
		sellMoney = int(float64(sellMoney) * (0.4 + 0.6*dur))
	}
	if sellMoney >= 1000 {
		sellMoney = sellMoney / 100 * 100
	} else if sellMoney >= 100 {
		sellMoney = sellMoney / 10 * 10
	}
	it.SellMoney = sellMoney
}

func (it *Item) CalculateRepairMoney(fast bool) int {
	deltaDur := it.MaxDurability - it.Durability
	if deltaDur == 0 {
		return 0
	}
	money := it.BuyMoney / conf.Price.Value.ItemSellPriceDivisor
	if it.Code == Code(13, 4) || it.Code == Code(13, 5) {
		money = it.BuyMoney
	}
	if money > 400000000 {
		money = 400000000
	}
	if money >= 1000 {
		money = money / 100 * 100
	} else if money >= 100 {
		money = money / 10 * 10
	}
	repairMoney := math.Sqrt(float64(money))
	repairMoney = 3.0 * repairMoney * math.Sqrt(repairMoney)
	repairMoney = repairMoney * float64(deltaDur) / float64(it.MaxDurability)
	repairMoney += 1
	if it.Durability == 0 {
		if it.Code == Code(13, 4) || it.Code == Code(13, 5) {
			repairMoney *= 2
		} else {
			repairMoney *= 1.4
		}
	}
	if fast {
		repairMoney += repairMoney * 0.05
	}
	money = int(repairMoney)
	if money >= 1000 {
		money = money / 100 * 100
	} else if money >= 100 {
		money = money / 10 * 10
	}
	return money
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
		if it.Addition&4 != 0 {
			it.AdditionAbsorbDamagePercent5 = true
		}
		if it.Addition&8 != 0 {
			it.AdditionAG50 = true
		}
		if it.Addition&16 != 0 {
			it.AdditionSpeed5 = true
		}
	case it.KindA == KindAWing: // wings
		switch {
		case it.WingAdditionAttack:
			it.AdditionAttack = it.Addition
		case it.WingAdditionMagicAttack:
			it.AdditionMagicAttack = it.Addition
		case it.WingAdditionCurseAttack:
			it.AdditionCurseAttack = it.Addition
		case it.WingAdditionDefense:
			it.AdditionDefense = it.Addition
		case it.WingAdditionRecoverHP:
			it.AdditionRecoverHP = it.Addition / 4
		}
	}

	// calc money
	it.CalculateMoney()
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
	case item.Code == Code(12, 3) && item.WingAdditionRecoverHP, // Wings of Spirits 圣灵之翼
		item.Code == Code(12, 4) && item.WingAdditionMagicAttack,   // Wings of Soul 魔魂之翼
		item.Code == Code(12, 5) && item.WingAdditionAttack,        // Wings of Dragon 飞龙之翼
		item.Code == Code(12, 6) && item.WingAdditionAttack,        // Wings of Darkness 暗黑之翼
		item.Code == Code(13, 30) && item.WingAdditionAttack,       // Wings of Darkness 王者披风
		item.Code == Code(12, 42) && item.WingAdditionMagicAttack,  // Wing of Despair 绝望之翼
		item.Code == Code(12, 49) && item.WingAdditionAttack,       // Cape of Fighter 武者披风
		item.Code == Code(12, 36) && item.WingAdditionDefense,      // Wing of Storm 暴风之翼
		item.Code == Code(12, 37) && item.WingAdditionDefense,      // Wing of Eternal 时空之翼
		item.Code == Code(12, 38) && item.WingAdditionDefense,      // Wing of Illusion 幻影之翼
		item.Code == Code(12, 39) && item.WingAdditionMagicAttack,  // Wing of Ruin 破灭之翼
		item.Code == Code(12, 40) && item.WingAdditionDefense,      // Cape of Emperor 帝王披风
		item.Code == Code(12, 43) && item.WingAdditionCurseAttack,  // Wing of Dimension 次元之翼
		item.Code == Code(12, 50) && item.WingAdditionDefense,      // Cape of Overrule 斗皇披风
		item.Code == Code(12, 262) && item.WingAdditionRecoverHP,   // Cloak of Death 死亡披风
		item.Code == Code(12, 263) && item.WingAdditionRecoverHP,   // Wings of Chaos 混沌之翼
		item.Code == Code(12, 264) && item.WingAdditionCurseAttack, // Wings of Magic 魔力之翼
		item.Code == Code(12, 265) && item.WingAdditionRecoverHP:   // Wings of Life 生命之翼
		data[3] |= 1 << 5
	case item.Code == Code(12, 36) && item.WingAdditionAttack, // Wing of Storm 暴风之翼
		item.Code == Code(12, 37) && item.WingAdditionMagicAttack,  // Wing of Eternal 时空之翼
		item.Code == Code(12, 38) && item.WingAdditionAttack,       // Wing of Illusion 幻影之翼
		item.Code == Code(12, 39) && item.WingAdditionAttack,       // Wing of Ruin 破灭之翼
		item.Code == Code(12, 40) && item.WingAdditionAttack,       // Cape of Emperor 帝王披风
		item.Code == Code(12, 43) && item.WingAdditionMagicAttack,  // Wing of Dimension 次元之翼
		item.Code == Code(12, 50) && item.WingAdditionAttack,       // Cape of Overrule 斗皇披风
		item.Code == Code(12, 262) && item.WingAdditionAttack,      // Cloak of Death 死亡披风
		item.Code == Code(12, 263) && item.WingAdditionAttack,      // Wings of Chaos 混沌之翼
		item.Code == Code(12, 264) && item.WingAdditionMagicAttack, // Wings of Magic 魔力之翼
		item.Code == Code(12, 265) && item.WingAdditionAttack:      // Wings of Life 生命之翼
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
	Size  int
	Items []*Item
	Flags []bool
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

func (pi *PositionedItems) UnmarshalJSON(
	buf []byte,
	CheckFlagsForItem func(int, *Item) bool,
	SetFlagsForItem func(int, *Item),
) error {
	var items []*Item
	err := json.Unmarshal(buf, &items)
	if err != nil {
		return err
	}
	pi.Items = make([]*Item, pi.Size)
	pi.Flags = make([]bool, pi.Size)
	for _, it := range items {
		it.Code = Code(it.Section, it.Index)
		itemBase, err := ItemTable.GetItemBase(it.Section, it.Index)
		if err != nil {
			return err
		}
		it.ItemBase = itemBase
		it.Calc()
		ok := CheckFlagsForItem(it.Position, it)
		if !ok {
			slog.Error("pi.CheckFlagsForItem",
				"position", it.Position, "name", it.Name)
			continue
		}
		SetFlagsForItem(it.Position, it)
		pi.Items[it.Position] = it
	}
	return nil
}

func (pi PositionedItems) Value() (driver.Value, error) {
	return pi.MarshalJSON()
}
