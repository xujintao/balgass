package handle

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/xujintao/balgass/cmd/server_connect/model"
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
	0x05:   checkVersion,
	0xf406: getServerList,
	0xf403: getServerInfo,
}

func checkVersion(req *protocol.Message) interface{} {
	// unmarshal
	ver := model.Version{}
	binary.Read(bytes.NewReader(req.Body), binary.LittleEndian, &ver)

	// service
	result := service.Update.CheckVersion(&ver)

	msg := &protocol.Message{Flag: 0xC1, Code: 0x05}
	// marshal
	// var body []byte
	// switch result := result.(type) {
	// case bool:
	// 	body = append(body, 1)
	// case *model.AutoUpdateConfig:
	// 	buf := new(bytes.Buffer)
	// 	binary.Write(buf, binary.LittleEndian, result)
	// 	body = buf.Bytes()
	// }
	// msg.Body = body

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, result)
	msg.Body = buf.Bytes()
	return msg
}

func getServerList(req *protocol.Message) interface{} {
	// service
	slis := service.Server.GetServerList()

	// write head
	msg := &protocol.Message2{Flag: 0xc1, Code: 0xf4, SubCode: 0x06}

	// marshal body and write
	res := model.ServerListRes{
		Size: int16(len(slis)),
		Slis: slis,
	}
	buf, err := res.Marshal()
	if err != nil {
		log.Println("Marshal failed", err)
	}
	msg.Body = buf
	return msg
}

func getServerInfo(req *protocol.Message) interface{} {
	// get param
	code := binary.LittleEndian.Uint16(req.Body)

	// service
	sci := service.Server.GetServerInfo(int(code))
	if sci == nil {
		log.Println("GetServerInfo failed, code: %u", code)
		return nil
	}

	// write head
	msg := &protocol.Message2{Flag: 0xc1, Code: 0xf4, SubCode: 0x03}
	// marshal body and write
	buf, err := sci.Marshal()
	if err != nil {
		log.Println("Marshal failed", err)
	}
	msg.Body = buf
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
	0x01: registerServer,
}

func registerServer(msg *wzudp.Message) *wzudp.Message {
	// unmarshal
	sri := model.ServerRegisterInfo{}
	if err := sri.Unmarshal(msg.Body); err != nil {
		log.Println(err)
		return nil
	}

	// service
	service.Server.RegisterServer(&sri)
	return nil
}
