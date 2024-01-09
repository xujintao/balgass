package exp

import (
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
}

var ExperienceTable []int
var MasterExperienceTable []int
