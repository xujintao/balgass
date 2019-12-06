package handle

import (
	"log"

	"github.com/xujintao/balgass/cmd/server_data/model"
	"github.com/xujintao/balgass/network"
)

type handleBase struct {
	cmds map[int]func(server *model.Server, req *network.Request)
}

// OnConn implements network.Handler.OnConn
func (*handleBase) OnConn(addr string, conn network.ConnWriter) (interface{}, error) {
	// if err := service.UserManager.AddUser(addr, conn); err != nil {
	// 	log.Printf("[%s] add user failed, %v", addr, err)
	// 	return nil, err
	// }
	log.Printf("[%s] connected", addr)
	res := &network.Response{}
	res.WriteHead2(0xC1, 0x00, 0x01)
	conn.Write(res)
	return addr, nil
}

// OnClose implements network.Handler.OnConn
func (*handleBase) OnClose(v interface{}) {
	addr := v.(string)
	// service.UserManager.DelUser(addr)
	log.Printf("[%s] disconnected", addr)
}

// Handle *CMDHandle implements network.Handler
func (h *handleBase) Handle(v interface{}, req *network.Request) {
	// id := v.(string)
	// user := service.ServerManager.GetServer(id)

	// code := req.Code
	// if h, ok := h.cmds[int(code)]; ok {
	// 	h(user, req)
	// 	return
	// }
	// subcode := req.Body[0]
	// codes := []byte{code, subcode}
	// code16 := binary.BigEndian.Uint16(codes)
	// if h, ok := h.cmds[int(code16)]; ok {
	// 	req.Body = req.Body[1:]
	// 	h(user, req)
	// 	return
	// }
	// log.Printf("[%s], invalid cmd, code:[%02x], body:[%v]", uid, code, network.Hex2string(req.Body))
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
