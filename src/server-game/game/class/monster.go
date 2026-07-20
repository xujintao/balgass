package class

import (
	"encoding/xml"

	"github.com/xujintao/balgass/src/server-game/conf"
)

var MonsterTable monsterTable

type monsterTable map[int]*MonsterConfig

type MonsterConfig struct {
	Text                   string `xml:",chardata"`
	Index                  int    `xml:"Index,attr"`
	ExpType                int    `xml:"ExpType,attr"`
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
	MaxItemLevel           int    `xml:"MaxItemLevel,attr"`
	MonsterSkill           int    `xml:"MonsterSkill,attr"`
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

func init() {
	MonsterTable.init()
}

func (m *monsterTable) init() {
	// MonsterList was generated 2023-07-17 11:34:17 by https://xml-to-go.github.io/ in Ukraine.
	type MonsterList struct {
		XMLName  xml.Name         `xml:"MonsterList"`
		Text     string           `xml:",chardata"`
		Monsters []*MonsterConfig `xml:"Monster"`
	}
	var monsterList MonsterList
	conf.XML(conf.PathCommon, "Monsters/IGC_MonsterList.xml", &monsterList)
	tm := make(monsterTable)
	for _, monster := range monsterList.Monsters {
		tm[monster.Index] = monster
	}
	*m = tm
}
