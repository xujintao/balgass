package main

import (
	"encoding/binary"

	"github.com/xujintao/balgass/protocol"
	"github.com/xujintao/balgass/wzudp"
)

// tcp cmd
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

// udp cmd
type udpcmdHandle struct{}

func (udpcmdHandle) Handle(req *wzudp.Message) *wzudp.Message {
	code := req.Code
	if h, ok := udpcmds[int(code)]; ok {
		return h(req)
	}
	subcode := req.Body[0]
	codes := []byte{code, subcode}
	code16 := binary.BigEndian.Uint16(codes)
	if h, ok := udpcmds[int(code16)]; ok {
		req.Body = req.Body[1:]
		return h(req)
	}
	return nil
}

var udpcmds = map[int]func(msg *wzudp.Message) *wzudp.Message{
	0x01: handleServerState,
}

func handleServerState(msg *wzudp.Message) *wzudp.Message {
	return nil
}
