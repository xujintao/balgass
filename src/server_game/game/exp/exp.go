package exp

import (
	"encoding/xml"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/formula"
)

func init() {
	ExperienceTable = make([]int, conf.Common.General.MaxLevelNormal+1)
	MasterExperienceTable = make([]int, conf.Common.General.MaxLevelMaster+1)
	for i := range ExperienceTable {
		formula.SetExpTable_Normal(i, &ExperienceTable[i])
	}
	for i := range MasterExperienceTable {
		formula.SetExpTable_Master(i+1, conf.Common.General.MaxLevelNormal, &MasterExperienceTable[i])
	}
	ExpManager.init()
}

var ExperienceTable []int
var MasterExperienceTable []int

var ExpManager expManager

type expManager struct {
	Normal float64
	Master float64
	Event  float64
}

func (m *expManager) init() {
	// ExpSystem was generated 2024-01-10 19:09:46 by https://xml-to-go.github.io/ in Ukraine.
	type ExpSystem struct {
		XMLName   xml.Name `xml:"ExpSystem"`
		Text      string   `xml:",chardata"`
		CalcType  string   `xml:"CalcType,attr"`
		DebugMode string   `xml:"DebugMode,attr"`
		StaticExp struct {
			Text   string  `xml:",chardata"`
			Normal float64 `xml:"Normal,attr"`
			Master float64 `xml:"Master,attr"`
			Event  float64 `xml:"Event,attr"`
			Quest  float64 `xml:"Quest,attr"`
		} `xml:"StaticExp"`
		DynamicExpRangeList struct {
			Text  string `xml:",chardata"`
			Range []struct {
				Text      string `xml:",chardata"`
				MinReset  string `xml:"MinReset,attr"`
				MaxReset  string `xml:"MaxReset,attr"`
				MinLevel  string `xml:"MinLevel,attr"`
				MaxLevel  string `xml:"MaxLevel,attr"`
				NormalExp string `xml:"NormalExp,attr"`
				MasterExp string `xml:"MasterExp,attr"`
			} `xml:"Range"`
		} `xml:"DynamicExpRangeList"`
	}

	var expSystem ExpSystem
	conf.XML(conf.PathCommon, "IGC_ExpSystem.xml", &expSystem)
	m.Normal = expSystem.StaticExp.Normal
	m.Master = expSystem.StaticExp.Master
	m.Event = expSystem.StaticExp.Event
}
