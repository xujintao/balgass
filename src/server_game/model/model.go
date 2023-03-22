package model

type MsgChat struct {
	Name string
	Msg  string
}

type MsgWhisper struct {
	Name string
	Msg  string
}

type MsgLive struct {
	Time         int
	AttackSpeed  int
	Agility      int
	MagicSpeed   int
	Version      string
	ServerSeason int
}

type MsgConnectResult struct {
	CBIPack struct{} `cbi:"1"`
	Result  int      `cbi:"byte"`
	ID      int      `cbi:"word,big"`
	Version [5]uint8 `cbi:"byte"`
}

type MsgUseItem struct {
	InventoryPos       int `cbi:"byte"`
	InventoryPosTarget int `cbi:"byte"`
	ItemUseType        int `cbi:"byte"`
}

type MsgLearnMasterSkill struct {
	SkillIndex int `cbi:"dword"`
}

type MsgSkillList struct {
}
