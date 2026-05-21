package model

type MsgAddBot struct {
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type MsgAddBotReply struct {
	Key string `json:"key"`
}

type MsgDeleteBot struct {
	Account string `json:"account" validate:"required"`
	Name    string `json:"name" validate:"required"`
}

type MsgDeleteBotReply struct {
	Key string `json:"key"`
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
