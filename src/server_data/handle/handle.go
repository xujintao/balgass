package handle

import (
	"encoding/binary"
	"log"

	"github.com/xujintao/balgass/cmd/server_data/service"
	"github.com/xujintao/balgass/network"
)

type handleBase struct {
	cmds map[int]func(index interface{}, req *network.Request)
}

// OnConn implements network.Handler.OnConn
func (*handleBase) OnConn(addr string, conn network.ConnWriter) (interface{}, error) {
	index, err := service.ServerManager.ServerAdd(addr, conn)
	if err != nil {
		log.Printf("[%s] add server failed, %v", addr, err)
		return nil, err
	}
	log.Printf("[%s] connected", addr)
	return index, nil
}

// OnClose implements network.Handler.OnConn
func (*handleBase) OnClose(v interface{}) {
	index := v.(string)
	addr := service.ServerManager.ServerGetAddr(index)
	log.Printf("[%s] disconnected", addr)
	service.ServerManager.ServerExit(index)
}

// Handle *CMDHandle implements network.Handler
func (h *handleBase) Handle(index interface{}, req *network.Request) {
	code := req.Code
	if h, ok := h.cmds[int(code)]; ok {
		h(index, req)
		return
	}
	subcode := req.Body[0]
	codes := []byte{code, subcode}
	code16 := binary.BigEndian.Uint16(codes)
	if h, ok := h.cmds[int(code16)]; ok {
		req.Body = req.Body[1:]
		h(index, req)
		return
	}
	addr := service.ServerManager.ServerGetAddr(index)
	log.Printf("[%s], invalid cmd, code:[%02x], body:[%v]", addr, code, network.Hex2string(req.Body))
}

var (
	HandleJoin   *handleJoin
	HandleData   *handleData
	HandleExData *handleExData
)

func init() {
	HandleJoin = &handleJoin{handleBase{cmds: cmdsJoin}}
	HandleData = &handleData{handleBase{cmds: cmdsData}, nil, nil}
	HandleJoin = &handleJoin{handleBase{cmds: cmdsExData}}
}