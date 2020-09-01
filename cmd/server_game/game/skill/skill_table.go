package skill

import (
	"path"

	"github.com/xujintao/balgass/cmd/server_game/conf"
	"github.com/xujintao/balgass/cmd/server_game/game"
)

var SkillTable map[int]*SkillBase

func init() {
	type skillListConfig struct {
		Skills []*SkillBase `xml:"Skill"`
	}
	var skillList skillListConfig
	conf.XML(path.Join(conf.PathCommon, "Skills/IGC_SkillList.xml"), &skillList)

	// array -> map
	SkillTable = make(map[int]*SkillBase)
	for _, v := range skillList.Skills {
		v.ReqClass[game.Wizard] = v.DarkWizard
		v.ReqClass[game.Knight] = v.DarkKnight
		v.ReqClass[game.Elf] = v.FairyElf
		v.ReqClass[game.Magumsa] = v.MagicGladiator
		v.ReqClass[game.DarkLord] = v.DarkLord
		v.ReqClass[game.RageFighter] = v.RageFighter
		// v.ReqClassChar[game.GrowLancer] = v.GrowLancer
		SkillTable[v.Index] = v
	}
}

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
	ItemSkill      int    `xml:"ItemSkill,attr"`
	IsDamage       int    `xml:"isDamage,attr"`
	BuffIndex      int    `xml:"BuffIndex,attr"`
}