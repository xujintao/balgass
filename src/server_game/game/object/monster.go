package object

import (
	"encoding/xml"
	"fmt"

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
	Level                  string `xml:"Level,attr"`
	HP                     string `xml:"HP,attr"`
	MP                     string `xml:"MP,attr"`
	DamageMin              string `xml:"DamageMin,attr"`
	DamageMax              string `xml:"DamageMax,attr"`
	Defense                string `xml:"Defense,attr"`
	MagicDefense           string `xml:"MagicDefense,attr"`
	AttackRate             string `xml:"AttackRate,attr"`
	BlockRate              string `xml:"BlockRate,attr"`
	MoveRange              string `xml:"MoveRange,attr"`
	AttackType             string `xml:"AttackType,attr"`
	AttackRange            string `xml:"AttackRange,attr"`
	ViewRange              string `xml:"ViewRange,attr"`
	MoveSpeed              string `xml:"MoveSpeed,attr"`
	AttackSpeed            string `xml:"AttackSpeed,attr"`
	RegenTime              string `xml:"RegenTime,attr"`
	Attribute              string `xml:"Attribute,attr"`
	ItemDropRate           string `xml:"ItemDropRate,attr"`
	MoneyDropRate          string `xml:"MoneyDropRate,attr"`
	MaxItemLevel           string `xml:"MaxItemLevel,attr"`
	MonsterSkill           string `xml:"MonsterSkill,attr"`
	IceRes                 string `xml:"IceRes,attr"`
	PoisonRes              string `xml:"PoisonRes,attr"`
	LightRes               string `xml:"LightRes,attr"`
	FireRes                string `xml:"FireRes,attr"`
	PentagramMainAttrib    int    `xml:"PentagramMainAttrib,attr"`
	PentagramAttribPattern string `xml:"PentagramAttribPattern,attr"`
	PentagramDamageMin     string `xml:"PentagramDamageMin,attr"`
	PentagramDamageMax     string `xml:"PentagramDamageMax,attr"`
	PentagramAttackRate    string `xml:"PentagramAttackRate,attr"`
	PentagramDefenseRate   string `xml:"PentagramDefenseRate,attr"`
	PentagramDefense       string `xml:"PentagramDefense,attr"`
	Name                   string `xml:"Name,attr"`
	Annotation             string `xml:"annotation,attr"`
}

func NewMonster(kind int) *Monster {
	mc, ok := MonsterTable[kind]
	if !ok {
		panic(fmt.Sprintf("monster invalid [kind]%d", kind))
	}
	return &Monster{MonsterConfig: *mc}
}

type Monster struct {
	object
	MonsterConfig
}
