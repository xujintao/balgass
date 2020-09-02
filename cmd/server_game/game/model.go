package game

type MsgItemUse struct {
	InventoryPos       int `cbi:"byte"`
	InventoryPosTarget int `cbi:"byte"`
	ItemUseType        int `cbi:"byte"`
}

type MsgMasterSkillLearn struct {
	SkillIndex int `cbi:"dword"`
}
