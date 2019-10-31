package handle

import (
	"encoding/binary"
	"unsafe"

	"github.com/xujintao/balgass/cmd/server_connect/service"
	"github.com/xujintao/balgass/protocol"
	"github.com/xujintao/balgass/wzudp"
)

// CMDHandle tcp cmd handle
type CMDHandle struct{}

func (CMDHandle) Handle(req *protocol.Message) interface{} {
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

var cmds = map[int]func(msg *protocol.Message) interface{}{
	0xf3:   handleAlive,
	0xf406: getServerList,
	0xf403: getServerInfo,
}

func handleAlive(req *protocol.Message) interface{} {
	return nil
}

func getServerList(req *protocol.Message) interface{} {
	// handle
	sibs := service.Server.GetServerList()

	// write head
	msg := &protocol.Message2{Flag: 0xc1, Code: 0xf4, SubCode: 0x06}

	// write body and marshal
	var body []byte
	binary.BigEndian.PutUint16(body, uint16(len(sibs)))
	for _, sib := range sibs {
		s := (*[unsafe.Sizeof(*sib)]byte)(unsafe.Pointer(sib))[:]
		body = append(body, s...)
	}
	msg.Body = body
	return msg
}

func getServerInfo(req *protocol.Message) interface{} {
	code := int(req.Body[0])

	// handle
	sic := service.Server.GetServerInfo(code)

	// write head
	msg := &protocol.Message2{Flag: 0xc1, Code: 0xf4, SubCode: 0x03}

	// write body and marshal
	var body []byte
	var ip [16]byte
	copy(ip[:], sic.IP)
	body = append(body, ip[:]...)
	var port []byte
	binary.LittleEndian.PutUint16(port, sic.Port)
	body = append(body, port...)
	msg.Body = body
	return msg
}

// UDPCMDHandle udp cmd handle
type UDPCMDHandle struct{}

func (UDPCMDHandle) Handle(req *wzudp.Message) *wzudp.Message {
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
	0x01: updateServerInfo,
}

func updateServerInfo(msg *wzudp.Message) *wzudp.Message {
	si := *(**service.ServerInfo)(unsafe.Pointer(&msg.Body))
	service.Server.UpdateServerInfo(si)
	return nil
}
