package skill

import (
	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/class"
)

func init() {
	SkillManager.init()
}

var SkillManager skillManager

type SkillBase struct {
	Index          int    `xml:"Index,attr"`
	Name           string `xml:"Name,attr"`
	ReqLevel       int    `xml:"ReqLevel,attr"`
	Damage         int    `xml:"Damage,attr"`
	STID           int    `xml:"STID,attr"`
	ManaUsage      int    `xml:"ManaUsage,attr"`
	BPUsage        int    `xml:"BPUsage,attr"`
	Distance       int    `xml:"Distance,attr"`
	Delay          int    `xml:"Delay,attr"`
	ReqStrength    int    `xml:"ReqStrength,attr"`
	ReqDexterity   int    `xml:"ReqDexterity,attr"`
	ReqEnergy      int    `xml:"ReqEnergy,attr"`
	ReqCommand     int    `xml:"ReqCommand,attr"`
	ReqMLPoint     int    `xml:"ReqMLPoint,attr"`
	Attribute      int    `xml:"Attribute,attr"`
	Type           int    `xml:"Type,attr"`
	UseType        int    `xml:"UseType,attr"`
	Brand          int    `xml:"Brand,attr"`
	KillCount      int    `xml:"KillCount,attr"`
	ReqStatus0     int    `xml:"ReqStatus0,attr"`
	ReqStatus1     int    `xml:"ReqStatus1,attr"`
	ReqStatus2     int    `xml:"ReqStatus2,attr"`
	DarkWizard     int    `xml:"DarkWizard,attr"`
	DarkKnight     int    `xml:"DarkKnight,attr"`
	FairyElf       int    `xml:"FairyElf,attr"`
	MagicGladiator int    `xml:"MagicGladiator,attr"`
	DarkLord       int    `xml:"DarkLord,attr"`
	Summoner       int    `xml:"Summoner,attr"`
	RageFighter    int    `xml:"RageFighter,attr"`
	ReqClass       [8]int `xml:"-"`
	Rank           int    `xml:"Rank,attr"`
	Group          int    `xml:"Group,attr"`
	HP             int    `xml:"HP,attr"`
	SD             int    `xml:"SD,attr"`
	Duration       int    `xml:"Duration,attr"`
	IconNumber     int    `xml:"IconNumber,attr"`
	ItemSkill      bool   `xml:"ItemSkill,attr"`
	IsDamage       int    `xml:"isDamage,attr"`
	BuffIndex      int    `xml:"BuffIndex,attr"`
}

type MasterSkillBase struct {
	Index        int    `xml:"Index,attr"`
	ReqMinPoint  int    `xml:"ReqMinPoint,attr"`
	MaxPoint     int    `xml:"MaxPoint,attr"`
	ParentSkill1 int    `xml:"ParentSkill1,attr"`
	ParentSkill2 int    `xml:"ParentSkill2,attr"`
	SkillID      int    `xml:"MagicNumber,attr"`
	Name         string `xml:"Name,attr"`
}

type valueType int

const (
	valueTypeNormal = iota
	valueTypeDamage
	valueTypeManaInc
)

type masterSkillValue struct {
	valueType valueType
	values    [21]float32
}

type skillManager struct {
	skillTable            map[int]*SkillBase
	masterSkillTable      [8][3][9][4]*MasterSkillBase
	masterSkillValueTable [30]masterSkillValue
}

func (m *skillManager) init() {
	type skillListConfig struct {
		Skills []*SkillBase `xml:"Skill"`
	}
	var skillList skillListConfig
	conf.XML(conf.PathCommon, "Skills/IGC_SkillList.xml", &skillList)

	// array -> map
	m.skillTable = make(map[int]*SkillBase)
	for _, v := range skillList.Skills {
		v.ReqClass[class.Wizard] = v.DarkWizard
		v.ReqClass[class.Knight] = v.DarkKnight
		v.ReqClass[class.Elf] = v.FairyElf
		v.ReqClass[class.Magumsa] = v.MagicGladiator
		v.ReqClass[class.DarkLord] = v.DarkLord
		v.ReqClass[class.Summoner] = v.Summoner
		v.ReqClass[class.RageFighter] = v.RageFighter
		// v.ReqClass[class.GrowLancer] = v.GrowLancer
		m.skillTable[v.Index] = v
	}

	type MasterSkillTree struct {
		Class []struct {
			ID   int `xml:"ID,attr"`
			Tree []struct {
				Type   int                `xml:"Type,attr"`
				Skills []*MasterSkillBase `xml:"Skill"`
			} `xml:"Tree"`
		} `xml:"Class"`
	}
	var masterSkillTree MasterSkillTree
	conf.XML(conf.PathCommon, "IGC_MasterSkillTree.xml", &masterSkillTree)
	id2class := map[int]class.Class{
		1:   class.Knight,
		2:   class.Wizard,
		4:   class.Elf,
		8:   class.Summoner,
		16:  class.Magumsa,
		32:  class.DarkLord,
		64:  class.RageFighter,
		128: class.GrowLancer,
	}
	for _, class := range masterSkillTree.Class {
		for _, tree := range class.Tree {
			for _, skill := range tree.Skills {
				index := skill.Index%36 - 1
				rank := index >> 2
				pos := index % 4
				m.masterSkillTable[id2class[class.ID]][tree.Type][rank][pos] = skill
			}
		}
	}
	// for t := 0; t < 3; t++ {
	// 	for rank := 0; rank < 9; rank++ {
	// 		for pos := 0; pos < 4; pos++ {
	// 			skill := m.table[class.Knight][t][rank][pos]
	// 			if skill == nil {
	// 				fmt.Print("[null]")
	// 				fmt.Print("\t")
	// 				continue
	// 			}
	// 			fmt.Print(skill.Name)
	// 			fmt.Print("\t")
	// 		}
	// 		fmt.Println()
	// 	}
	// }
	// fmt.Println(1)

	// fulfill masterSkillVauleTable by lua script
}
