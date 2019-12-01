package handle

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"

	"github.com/xujintao/balgass/cmd/server_connect/model"
	"github.com/xujintao/balgass/cmd/server_connect/service"
	"github.com/xujintao/balgass/protocol"
)

// CMDHandle tcp cmd handle
type CMDHandle struct{}

// Handle *CMDHandle implements protocol.Handler
func (CMDHandle) Handle(req *protocol.Request, res *protocol.Response) bool {
	code := req.Code
	if h, ok := cmds[int(code)]; ok {
		return h(req, res)
	}
	subcode := req.Body[0]
	codes := []byte{code, subcode}
	code16 := binary.BigEndian.Uint16(codes)
	if h, ok := cmds[int(code16)]; ok {
		req.Body = req.Body[1:]
		return h(req, res)
	}
	log.Printf("invalid cmd, code:%02dx, body: %v", code, req.Body)
	return false
}

// HandleUDP implements protocol.UDPHandler
func (CMDHandle) HandleUDP(req *protocol.Request, res *protocol.Response) bool {
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

var udpcmds = map[int]func(req *protocol.Request, res *protocol.Response) bool{
	0x0100: registerServer,
}

func registerServer(req *protocol.Request, res *protocol.Response) bool {
	// unmarshal
	sri := model.ServerRegisterInfo{}
	if err := sri.Unmarshal(req.Body); err != nil {
		log.Println(err)
		return false
	}

	// service
	service.Server.RegisterServer(&sri)
	return false
}

// OnConn implements protocol.Handler.OnConn
func (CMDHandle) OnConn(req *protocol.Request, res *protocol.Response) bool {
	res.WriteHead2(0xC1, 0x00, 0x01)
	return true
}

// TrackConnState track the connect state
func (CMDHandle) TrackConnState(c net.Conn, state protocol.ConnState) {
	switch state {
	case protocol.StateNew:
		log.Printf("%s connected", c.RemoteAddr().String())
	case protocol.StateClosed:
		log.Printf("%s disconnected", c.RemoteAddr().String())
	}
}

var cmds = map[int]func(req *protocol.Request, res *protocol.Response) bool{
	0x05:   checkVersion,
	0xf406: getServerList,
	0xf403: getServerInfo,
}

func checkVersion(req *protocol.Request, res *protocol.Response) bool {
	// unmarshal
	ver := model.Version{}
	if err := binary.Read(bytes.NewReader(req.Body), binary.LittleEndian, &ver); err != nil {
		log.Println(err)
		return false
	}

	// service
	result := service.Update.CheckVersion(&ver)

	// marshal
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, result); err != nil {
		log.Println(err)
		return false
	}

	// write
	res.WriteHead(0xC1, 0x05).Write(buf.Bytes())
	return true
}

func getServerList(req *protocol.Request, res *protocol.Response) bool {
	// service
	slis := service.Server.GetServerList()

	// marshal body and write
	slr := model.ServerListRes{
		Size: int16(len(slis)),
		Slis: slis,
	}
	buf, err := slr.Marshal()
	if err != nil {
		log.Println("Marshal failed", err)
		return false
	}

	// write
	res.WriteHead2(0xc2, 0xf4, 0x06).Write(buf)
	return true
}

func getServerInfo(req *protocol.Request, res *protocol.Response) bool {
	// get param
	code := binary.LittleEndian.Uint16(req.Body)

	// service
	sci := service.Server.GetServerInfo(int(code))
	if sci == nil {
		log.Println("GetServerInfo failed, code: %u", code)
		return false
	}

	// marshal body and write
	buf, err := sci.Marshal()
	if err != nil {
		log.Println("Marshal failed", err)
		return false
	}
	// write head
	res.WriteHead2(0xc1, 0xf4, 0x03).Write(buf)
	return true
}
