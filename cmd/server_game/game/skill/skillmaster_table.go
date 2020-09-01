package skill

import (
	"path"

	"github.com/xujintao/balgass/cmd/server_game/conf"
)

type SkillMasterBase struct {
	Index        int    `xml:"Index,attr"`
	ReqMinPoint  int    `xml:"ReqMinPoint,attr"`
	MaxPoint     int    `xml:"MaxPoint,attr"`
	ParentSkill1 int    `xml:"ParentSkill1,attr"`
	ParentSkill2 int    `xml:"ParentSkill2,attr"`
	SkillID      int    `xml:"MagicNumber,attr"`
	Name         string `xml:"Name,attr"`
}

func init() {
	type skillMasterTreeConfig struct {
		Class []struct {
			ID   int `xml:"ID,attr"`
			Tree []struct {
				Type   int                `xml:"Type,attr"`
				Skills []*SkillMasterBase `xml:"Skill"`
			} `xml:"Tree"`
		} `xml:"Class"`
	}
	var skillMasterTree skillMasterTreeConfig
	conf.XML(path.Join(conf.PathCommon, "IGC_MasterSkillTree.xml"), &skillMasterTree)

	for _, v := range skillMasterTree.Class {
		for _, v := range v.Tree {
			for _, v := range v.Skills {
				SkillMasterTable[v.SkillID] = v
			}
		}
	}
}

// var charIDTable = map[int]game.CharType{
// 	2:   game.CharTypeWizard,
// 	1:   game.CharTypeKnight,
// 	4:   game.CharTypeElf,
// 	16:  game.CharTypeMagumsa,
// 	8:   game.CharTypeSummoner,
// 	32:  game.CharTypeDarkLord,
// 	64:  game.CharTypeRageFighter,
// 	128: game.CharTypeGrowLancer,
// }

var SkillMasterTable skillMasterTable

// type skillMasterTable [game.MaxCharType][3][9][4]*SkillMasterBase
type skillMasterTable map[int]*SkillMasterBase

// func (t *skillMasterTable)
