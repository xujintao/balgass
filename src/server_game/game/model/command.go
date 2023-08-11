package model

type MsgGetObjectsByMapNumber struct {
	Number int `json:"number"`
}

type MsgGetObjectsByMapNumberReply struct {
	Name string `json:"name"`
	Data any    `json:"data"`
	Err  string `json:"err"`
}
