package model

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
	Number int `json:"number"`
}

type MsgSubscribeMapReply struct {
	Name string `json:"name"`
	Data any    `json:"data"`
	Err  string `json:"err"`
}
