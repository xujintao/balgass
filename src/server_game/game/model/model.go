package model

type MsgConnectResult struct {
	CBIPack struct{} `cbi:"1"`
	Result  int      `cbi:"byte"`
	ID      int      `cbi:"word,big"`
	Version [5]uint8 `cbi:"byte"`
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
