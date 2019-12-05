package handle

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/xujintao/balgass/cmd/server_connect/model"
	"github.com/xujintao/balgass/cmd/server_connect/service"
	"github.com/xujintao/balgass/network"
)

// CMDHandle tcp cmd handle
type CMDHandle struct{}

// HandleUDP implements network.UDPHandler
func (CMDHandle) HandleUDP(req *network.Request, res *network.Response) bool {
	code := req.Code
	if h, ok := udpcmds[int(code)]; ok {
		return h(req, res)
	}
	subcode := req.Body[0]
	codes := []byte{code, subcode}
	code16 := binary.BigEndian.Uint16(codes)
	if h, ok := udpcmds[int(code16)]; ok {
		req.Body = req.Body[1:]
		return h(req, res)
	}
	log.Printf("invalid cmd, code:%02dx, body: %v", code, req.Body)
	return false
}

var udpcmds = map[int]func(req *network.Request, res *network.Response) bool{
	0x0100: registerServer,
}

func registerServer(req *network.Request, res *network.Response) bool {
	// unmarshal
	sri := model.ServerRegisterInfo{}
	if err := sri.Unmarshal(req.Body); err != nil {
		log.Println(err)
		return false
	}

	// service
	service.ServerManager.RegisterServer(&sri)
	return false
}

// OnConn implements network.Handler.OnConn
func (CMDHandle) OnConn(addr string, conn network.ConnWriter) (interface{}, error) {
	if err := service.UserManager.AddUser(addr, conn); err != nil {
		log.Printf("[%s] add user failed, %v", addr, err)
		return nil, err
	}
	log.Printf("[%s] connected", addr)
	res := &network.Response{}
	res.WriteHead2(0xC1, 0x00, 0x01)
	conn.Write(res)
	return addr, nil
}

// OnClose implements network.Handler.OnConn
func (CMDHandle) OnClose(v interface{}) {
	addr := v.(string)
	service.UserManager.DelUser(addr)
	log.Printf("[%s] disconnected", addr)
}

// Handle *CMDHandle implements network.Handler
func (CMDHandle) Handle(v interface{}, req *network.Request) {
	code := req.Code
	if h, ok := cmds[int(code)]; ok {
		h(v, req)
		return
	}
	subcode := req.Body[0]
	codes := []byte{code, subcode}
	code16 := binary.BigEndian.Uint16(codes)
	if h, ok := cmds[int(code16)]; ok {
		req.Body = req.Body[1:]
		h(v, req)
		return
	}
	log.Printf("[%s], invalid cmd, code:[%02x], body:[%v]", v.(string), code, network.Hex2string(req.Body))
}

var cmds = map[int]func(v interface{}, req *network.Request){
	0x05:   checkVersion,
	0xf406: getServerList,
	0xf403: getServerInfo,
}

func checkVersion(v interface{}, req *network.Request) {
	// unmarshal
	ver := model.Version{}
	if err := binary.Read(bytes.NewReader(req.Body), binary.LittleEndian, &ver); err != nil {
		log.Println(err)
		return
	}

	// service
	result := service.Update.CheckVersion(&ver)

	// marshal
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, result); err != nil {
		log.Println(err)
		return
	}

	// write
	res := &network.Response{}
	res.WriteHead(0xC1, 0x05).Write(buf.Bytes())
	service.UserManager.GetUser(v.(string)).Conn.Write(res)
}

func getServerList(v interface{}, req *network.Request) {
	// service
	slis := service.ServerManager.GetServerList()

	// marshal body and write
	slr := model.ServerListRes{
		Size: int16(len(slis)),
		Slis: slis,
	}
	buf, err := slr.Marshal()
	if err != nil {
		log.Println("Marshal failed", err)
		return
	}

	// write
	res := &network.Response{}
	res.WriteHead2(0xc2, 0xf4, 0x06).Write(buf)
	service.UserManager.GetUser(v.(string)).Conn.Write(res)
}

func getServerInfo(v interface{}, req *network.Request) {
	// get param
	code := binary.LittleEndian.Uint16(req.Body)

	// service
	sci := service.ServerManager.GetServerInfo(int(code))
	if sci == nil {
		log.Println("GetServerInfo failed, code: %u", code)
	}

	// marshal body and write
	buf, err := sci.Marshal()
	if err != nil {
		log.Println("Marshal failed", err)
	}
	// write head
	res := &network.Response{}
	res.WriteHead2(0xc1, 0xf4, 0x03).Write(buf)
	service.UserManager.GetUser(v.(string)).Conn.Write(res)
}
