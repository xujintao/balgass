package object

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/xujintao/balgass/src/server_game/conf"
)

func init() {
	// MonsterList was generated 2023-07-17 11:34:17 by https://xml-to-go.github.io/ in Ukraine.
	type MonsterList struct {
		XMLName  xml.Name         `xml:"MonsterList"`
		Text     string           `xml:",chardata"`
		Monsters []*MonsterConfig `xml:"Monster"`
	}
	var monsterList MonsterList
	conf.XML(conf.PathCommon, "Monsters/IGC_MonsterList.xml", &monsterList)
	MonsterTable = make(monsterTable)
	for _, monster := range monsterList.Monsters {
		MonsterTable[monster.Index] = monster
	}
}

var MonsterTable monsterTable

type monsterTable map[int]*MonsterConfig

type MonsterConfig struct {
	Text                   string `xml:",chardata"`
	Index                  int    `xml:"Index,attr"`
	ExpType                string `xml:"ExpType,attr"`
	Level                  int    `xml:"Level,attr"`
	HP                     int    `xml:"HP,attr"`
	MP                     int    `xml:"MP,attr"`
	DamageMin              int    `xml:"DamageMin,attr"`
	DamageMax              int    `xml:"DamageMax,attr"`
	Defense                int    `xml:"Defense,attr"`
	MagicDefense           int    `xml:"MagicDefense,attr"`
	AttackRate             int    `xml:"AttackRate,attr"`
	BlockRate              int    `xml:"BlockRate,attr"`
	MoveRange              int    `xml:"MoveRange,attr"`
	AttackType             int    `xml:"AttackType,attr"`
	AttackRange            int    `xml:"AttackRange,attr"`
	ViewRange              int    `xml:"ViewRange,attr"`
	MoveSpeed              int    `xml:"MoveSpeed,attr"`
	AttackSpeed            int    `xml:"AttackSpeed,attr"`
	RegenTime              int    `xml:"RegenTime,attr"`
	Attribute              int    `xml:"Attribute,attr"`
	ItemDropRate           int    `xml:"ItemDropRate,attr"`
	MoneyDropRate          int    `xml:"MoneyDropRate,attr"`
	MaxItemLevel           string `xml:"MaxItemLevel,attr"`
	MonsterSkill           string `xml:"MonsterSkill,attr"`
	IceRes                 int    `xml:"IceRes,attr"`
	PoisonRes              int    `xml:"PoisonRes,attr"`
	LightRes               int    `xml:"LightRes,attr"`
	FireRes                int    `xml:"FireRes,attr"`
	PentagramMainAttrib    int    `xml:"PentagramMainAttrib,attr"`
	PentagramAttribPattern int    `xml:"PentagramAttribPattern,attr"`
	PentagramDamageMin     int    `xml:"PentagramDamageMin,attr"`
	PentagramDamageMax     int    `xml:"PentagramDamageMax,attr"`
	PentagramAttackRate    int    `xml:"PentagramAttackRate,attr"`
	PentagramDefenseRate   int    `xml:"PentagramDefenseRate,attr"`
	PentagramDefense       int    `xml:"PentagramDefense,attr"`
	Name                   string `xml:"Name,attr"`
	Annotation             string `xml:"annotation,attr"`
}

func NewMonster(class int) *Monster {
	mc, ok := MonsterTable[class]
	if !ok {
		panic(fmt.Sprintf("monster invalid [class]%d", class))
	}
	monster := Monster{}
	monster.init()
	monster.ConnectState = ConnectStatePlaying
	monster.Live = true
	monster.State = 1
	switch class {
	case 240: // 仓库使者塞弗特
		monster.NpcType = NpcTypeWarehouse
	case 238, 368, 369, 370, 452, 453, 478, 450:
		monster.NpcType = NpcTypeChaosMix
	case 236:
		monster.NpcType = NpcTypeGoldarcher
	case 582:
		monster.NpcType = NpcTypePentagramMix
	default:
		monster.NpcType = NpcTypeNone
	}
	switch {
	case class >= 204 && class <= 259 ||
		class >= 367 && class <= 385 ||
		class >= 406 && class <= 408 ||
		class >= 414 && class <= 417 ||
		class >= 450 && class <= 453 ||
		class >= 464 && class <= 475 && class != 466 ||
		class == 478 || class == 479 ||
		class == 492 ||
		class == 522 ||
		class >= 540 && class <= 547 ||
		class >= 566 && class <= 568 ||
		class >= 577 && class <= 584 ||
		class == 603 || class == 604 ||
		class == 643 ||
		class == 651 ||
		class >= 658 && class <= 668 ||
		class >= 682 && class <= 688:
		monster.Type = ObjectNPC
	default:
		monster.Type = ObjectMonster
	}
	monster.Class = class
	monster.Name = mc.Name
	monster.Annotation = mc.Annotation
	monster.Level = mc.Level
	monster.attackDamageMin = mc.DamageMin
	monster.attackDamageMax = mc.DamageMax
	monster.attackRate = mc.AttackRate
	monster.attackSpeed = mc.AttackSpeed
	monster.defense = mc.Defense
	monster.magicDefense = mc.MagicDefense
	monster.defenseRate = mc.BlockRate
	monster.HP = mc.HP
	monster.MaxHP = mc.HP
	monster.MP = mc.MP
	monster.MaxMP = mc.MP
	monster.moveRange = mc.MoveRange
	monster.moveSpeed = mc.MoveSpeed
	monster.attackRange = mc.AttackRange
	monster.attackType = mc.AttackType
	monster.viewRange = mc.ViewRange
	monster.attribute = mc.Attribute
	monster.itemDropRate = mc.ItemDropRate
	monster.moneyDropRate = mc.MoneyDropRate
	monster.maxRegenTime = time.Duration(mc.RegenTime)
	monster.pentagramAttributePattern = mc.PentagramAttribPattern
	monster.pentagramAttackMin = mc.PentagramDamageMin
	monster.pentagramAttackMax = mc.PentagramDamageMax
	monster.pentagramAttackRate = mc.PentagramAttackRate
	monster.pentagramDefense = mc.PentagramDefense
	switch {
	case monster.attackType >= 100:
		monster.addSkill(monster.attackType-100, 1)
	case monster.attackType >= 1:
		monster.addSkill(monster.attackType, 1)
	}
	switch class {
	case 161, 181, 189, 197, 267, 275: // 昆顿
		monster.addSkill(1, 1)   // 毒咒
		monster.addSkill(17, 1)  // 能量球
		monster.addSkill(55, 1)  // 玄月斩
		monster.addSkill(200, 1) // 召唤怪
		monster.addSkill(201, 1) // 免疫魔攻
		monster.addSkill(202, 1) // 免疫物攻
	case 149, 179, 187, 195, 265, 273, 335: // 暗黑巫师
		monster.addSkill(1, 1)  // 毒咒
		monster.addSkill(17, 1) // 能量球
	case 66, 73, 77: // 诅咒之王 蓝魔龙 天魔菲尼斯
		// 163, 165, 167, 171, 173, 427: // 赤色要塞
		monster.addSkill(17, 1) // 能量球
	case 89, 95, 112, 118, 124, 130, 143: // 骷灵巫师
		monster.addSkill(3, 1)  // 掌心雷
		monster.addSkill(17, 1) // 能量球
	case 433: // 骷髅法师
		monster.addSkill(3, 1) // 掌心雷
	case 561: // 美杜莎
		monster.addSkill(9, 1)   // 黑龙波
		monster.addSkill(38, 1)  // 单毒炎
		monster.addSkill(237, 1) // 闪电轰顶
		monster.addSkill(238, 1) // 黑暗之力
	case 673: // 辛维斯特
		// monster.addSkill(622, 1) // ?
	}
	return &monster
}

type NpcType int

const (
	NpcTypeNone = iota
	NpcTypeShop
	NpcTypeWarehouse
	NpcTypeChaosMix
	NpcTypeGoldarcher
	NpcTypePentagramMix
)

type Monster struct {
	object
	NpcType int
}

func (m *Monster) processRegen() {
	if m.ConnectState < ConnectStatePlaying {
		return
	}
	if time.Now().Unix()-int64(m.regenTime) < int64(m.maxRegenTime) {
		return
	}

	m.dieRegen = false
	m.State = 1
}
