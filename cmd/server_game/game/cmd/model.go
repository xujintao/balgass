package cmd

type MsgItemUse struct {
	InventoryPos       int `cbi:"byte"`
	InventoryPosTarget int `cbi:"byte"`
	ItemUseType        int `cbi:"byte"`
}

type MsgSkillMasterLearn struct {
	SkillIndex int `cbi:"dword"`
}
