package handle

import (
	"encoding/binary"
	"log"

	"github.com/xujintao/balgass/cmd/server_data/model"
	"github.com/xujintao/balgass/cmd/server_data/service"
	"github.com/xujintao/balgass/network"
)

type handleBase struct {
	cmds map[int]func(server *model.Server, req *network.Request)
}

// OnConn implements network.Handler.OnConn
func (*handleBase) OnConn(addr string, conn network.ConnWriter) (interface{}, error) {
	if err := service.ServerManager.ServerAdd(addr, conn); err != nil {
		log.Printf("[%s] add server failed, %v", addr, err)
		return nil, err
	}
	log.Printf("[%s] connected", addr)
	return addr, nil
}

// OnClose implements network.Handler.OnConn
func (*handleBase) OnClose(v interface{}) {
	addr := v.(string)
	service.ServerManager.ServerDel(addr)
	log.Printf("[%s] disconnected", addr)
}

// Handle *CMDHandle implements network.Handler
func (h *handleBase) Handle(id interface{}, req *network.Request) {
	sid := id.(string)
	user := service.ServerManager.ServerGet(sid)
	code := req.Code
	if h, ok := h.cmds[int(code)]; ok {
		h(user, req)
		return
	}
	subcode := req.Body[0]
	codes := []byte{code, subcode}
	code16 := binary.BigEndian.Uint16(codes)
	if h, ok := h.cmds[int(code16)]; ok {
		req.Body = req.Body[1:]
		h(user, req)
		return
	}
	log.Printf("[%s], invalid cmd, code:[%02x], body:[%v]", sid, code, network.Hex2string(req.Body))
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
