package model

type MsgConnectResult struct {
	CBIAlign int      `cbi:"1"`
	Result   int      `cbi:"byte"`
	ID       int      `cbi:"uint16,big"`
	Version  [5]uint8 `cbi:"bytes"`
}

type MsgItemUse struct {
	InventoryPos       int `cbi:"byte"`
	InventoryPosTarget int `cbi:"byte"`
	ItemUseType        int `cbi:"byte"`
}

type MsgMasterSkillLearn struct {
	SkillIndex int `cbi:"dword"`
}

type MsgSkillList struct {
}
