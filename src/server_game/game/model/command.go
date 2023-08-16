package model

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
