package main

import (
	"encoding/binary"

	"github.com/xujintao/balgass/protocol"
)

type cmdHandle struct{}

func (cmdHandle) Handle(req *protocol.Message) *protocol.Message {
	code := req.Code
	if h, ok := cmds[int(code)]; ok {
		return h(req)
	}
	subcode := req.Body[0]
	codes := []byte{code, subcode}
	code16 := binary.BigEndian.Uint16(codes)
	if h, ok := cmds[int(code16)]; ok {
		req.Body = req.Body[1:]
		return h(req)
	}
	return nil
}

var cmds = map[int]func(msg *protocol.Message) *protocol.Message{
	0xf3:   handleAlive,
	0xf406: handleServerList,
	0xf403: handleServerInfo,
}

func handleAlive(req *protocol.Message) *protocol.Message {
	return nil
}

func handleServerList(req *protocol.Message) *protocol.Message {
	return nil
}

func handleServerInfo(req *protocol.Message) *protocol.Message {
	return nil
}
