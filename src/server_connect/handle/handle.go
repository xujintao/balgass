package handle

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"reflect"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server_connect/service"
	"github.com/xujintao/balgass/src/server_connect/service/model"
)

func init() {
	C1C2Handle.init()
}

var C1C2Handle c1c2Handle

type c1c2Handle struct {
	apiIns  map[int]*apiIn
	apiOuts map[any]*apiOut
}

func (h *c1c2Handle) init() {
	// ingress
	h.apiIns = make(map[int]*apiIn)
	for _, v := range apiIns {
		if vv, ok := h.apiIns[v.code]; ok {
			log.Fatalf("duplicated api code[%d] action[%s] with code[%d] action[%s]\n",
				v.code, v.action, vv.code, vv.action)
		}
		h.apiIns[v.code] = v
	}

	// egress
	h.apiOuts = make(map[any]*apiOut)
	for _, v := range apiOuts {
		t := reflect.TypeOf(v.msg)
		if t.Kind() != reflect.Ptr {
			log.Fatalf("api code[%d] name[%s] msg field must be a pointer\n",
				v.code, v.name)
		}
		h.apiOuts[t] = v
	}
}

func (h *c1c2Handle) Handle(ctx context.Context, req *c1c2.Request) {
	v := ctx.Value(c1c2.UserContextKey)
	if v == nil {
		return
	}
	id := v.(int)

	var api *apiIn
	var ok bool
	code := int(req.Body[0])
	if api, ok = h.apiIns[code]; !ok {
		if len(req.Body) < 2 {
			log.Printf("invalid api [body]%s\n", hex.EncodeToString(req.Body))
			return
		}
		codes := []byte{req.Body[0], req.Body[1]}
		code = int(binary.BigEndian.Uint16(codes))
		if api, ok = h.apiIns[code]; !ok {
			log.Printf("invalid api [body]%s\n", hex.EncodeToString(req.Body))
			return
		}
		req.Body = req.Body[1:]
	}
	req.Body = req.Body[1:]

	t := reflect.TypeOf(api.msg)
	if _, ok := t.MethodByName("Unmarshal"); !ok {
		log.Printf("can't find Unmarshal method [msg]%s\n", t.String())
		return
	}
	msg := reflect.New(t.Elem())
	in := []reflect.Value{reflect.ValueOf(req.Body)}
	out := msg.MethodByName("Unmarshal").Call(in)
	err := out[0].Interface()
	if err != nil {
		log.Printf("Unmarshal failed [msg]%s [err]%v", msg.String(), err)
		return
	}
	service.Service.PlayerAction(id, api.action, msg.Interface())
}

func (h *c1c2Handle) marshal(msg any) (*c1c2.Response, error) {
	v := reflect.ValueOf(msg)
	t := v.Type()
	api, ok := h.apiOuts[t]
	if !ok {
		err := fmt.Errorf("%s has not yet be registered to api table", t.String())
		return nil, err
	}
	if _, ok := t.MethodByName("Marshal"); !ok {
		err := fmt.Errorf("%s has no Marshal Method", t.String())
		return nil, err
	}
	rets := v.MethodByName("Marshal").Call(nil)
	if len(rets) != 2 {
		err := fmt.Errorf("%s Marshal Method signature is invalid", t.String())
		return nil, err
	}
	data := rets[0].Bytes()
	err := rets[1].Interface()
	if err != nil {
		return nil, err.(error)
	}
	var buf bytes.Buffer
	if api.code>>8 != 0 {
		var codes [2]uint8
		binary.BigEndian.PutUint16(codes[:], uint16(api.code))
		buf.Write(codes[:])
	} else {
		buf.WriteByte(uint8(api.code))
	}
	buf.Write(data)

	var resp c1c2.Response
	resp.WriteHead(uint8(api.flag))
	resp.Write(buf.Bytes())
	return &resp, nil
}

type conn struct {
	*c1c2.Conn
}

func (c *conn) Addr() string {
	return c.RemoteAddr
}

func (c *conn) Write(msg any) error {
	resp, err := C1C2Handle.marshal(msg)
	if err != nil {
		return err
	}
	return c.Conn.Write(resp)
}

func (c *conn) Close() error {
	return c.Conn.Close()
}

// OnConn implements c1c2.Handler.OnConn
func (h *c1c2Handle) OnConn(c *c1c2.Conn) (any, error) {
	conn := conn{c}
	return service.Service.PlayerConn(&conn)
}

// OnClose implements c1c2.Handler.OnConn
func (h *c1c2Handle) OnClose(ctx context.Context) {
	v := ctx.Value(c1c2.UserContextKey)
	if v == nil {
		return
	}
	id := v.(int)
	service.Service.PlayerCloseConn(id)
}

func (h *c1c2Handle) HandleUDP(req *c1c2.Request, res *c1c2.Response) bool {
	apiIns := map[int]*apiIn{
		0x01: {0, 0x01, "Register", (*model.MsgRegister)(nil)},
	}
	var api *apiIn
	var ok bool
	code := int(req.Body[0])
	if api, ok = apiIns[code]; !ok {
		codes := []byte{req.Body[0], req.Body[1]}
		code = int(binary.BigEndian.Uint16(codes))
		if api, ok = h.apiIns[code]; !ok {
			log.Printf("invalid api [body]%v\n", req.Body)
			return false
		}
		req.Body = req.Body[1:]
	}
	req.Body = req.Body[1:]

	in := []reflect.Value{reflect.ValueOf(req.Body)}
	msg := reflect.New(reflect.TypeOf(api.msg).Elem())
	out := msg.MethodByName("Unmarshal").Call(in)
	err := out[0].Interface()
	if err != nil {
		log.Printf("%s UnMarshal failed", msg.String())
		return false
	}
	service.Service.ServerAction(api.action, msg.Interface())

	return false
}

type apiIn struct {
	id     int
	code   int
	action string
	msg    any
}

type apiOut struct {
	id   int
	flag int
	code int
	name string
	msg  any
}

var apiIns = [...]*apiIn{
	// {0, 0x01, "Register", (*model.MsgRegister)(nil)}, // para segridad, move to HandleUDP func
	{0, 0x04, "CheckVersion", (*model.MsgCheckVersion)(nil)}, // 1.04.44
	// {0, 0x05, "CheckVersion", (*model.MsgCheckVersion)(nil)}, // 1.05.25

	// TODO:
	// If we close main.exe suddently after GetServerList but before GetServer, it will send:
	// c122f33dda063c605599a7940d954b478a678d47919190722b3b09b787440a47ecfd from f006B4206
	// xor.Dec
	// c122f330ffffffffffffffffffffffffffffffffffffffff1dffffff16ff00000000
	// omit

	{0, 0xF406, "GetServerList", (*model.MsgGetServerList)(nil)},
	{0, 0xF403, "GetServer", (*model.MsgGetServer)(nil)},
}

var apiOuts = [...]*apiOut{
	{0, 0xC1, 0x00, "ConnectReply", (*model.MsgConnectReply)(nil)},
	{0, 0xC1, 0x02, "CheckVersionSuccess", (*model.MsgCheckVersionSuccess)(nil)},
	{0, 0xC1, 0x04, "CheckVersionFailed", (*model.MsgCheckVersionFailed)(nil)},
	{0, 0xC2, 0xF406, "GetServerListReply", (*model.MsgGetServerListReply)(nil)},
	{0, 0xC1, 0xF403, "GetServerReply", (*model.MsgGetServerReply)(nil)},
}
