package object

import "github.com/xujintao/balgass/src/server-game/game/model"

func (obj *Object) Chat(msg *model.MsgChat) {
	l := len(msg.Msg)
	if l == 0 {
		return
	}
	switch {
	case msg.Msg[0] == '!' && l > 2: // global announcement
		return
	case msg.Msg[0] == '/' && l > 1: // command
		return
	case msg.Msg[0] == '~' || msg.Msg[0] == ']': // party
		return
	case msg.Msg[0] == '$': // gens
		return
	case msg.Msg[0] == '@': //guild
		return
	default:
		reply := model.MsgChatReply{MsgChat: *msg}
		obj.PushViewport(&reply)
	}
}

func (obj *Object) Whisper(msg *model.MsgWhisper) {
	if len(msg.Name) == 0 {
		return
	}
	if obj.Name == msg.Name {
		return
	}
	tobj := ObjectManager.GetPlayerByName(msg.Name)
	if tobj == nil {
		reply := model.MsgWhisperReplyFailed{
			Flag: 0,
		}
		obj.Push(&reply)
		return
	}
	reply := model.MsgWhisperReply{}
	reply.Name = obj.Name
	reply.Msg = msg.Msg
	tobj.Push(&reply)
}
