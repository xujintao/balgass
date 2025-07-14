package model

type MsgAddBot struct {
	Name string `json:"name"`
}

type MsgAddBotReply struct {
}

type MsgDeleteBot struct {
	Name string `json:"name"`
}

type MsgDeleteBotReply struct {
}

type MsgOfflineAllObjects struct {
}

type MsgOfflineAllObjectsReply struct {
}

type MsgGetOnlineObjectNumber struct {
}

type MsgGetOnlineObjectNumberReply struct {
	PlayerNumber int `json:"player_number"`
	UserNumber   int `json:"user_number"`
}

type MsgHandleErrorReply struct {
	Err string `json:"err"`
}

type MsgSubscribeMap struct {
	Name string `json:"name"`
}

type MsgSubscribeMapReply struct {
	Name string `json:"name"`
	Data any    `json:"data"`
	Err  string `json:"err"`
}
