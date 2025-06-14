package handle

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server-connect/service"
	"github.com/xujintao/balgass/src/server-connect/service/object"
	"github.com/xujintao/balgass/src/server-connect/service/server"
)

func TestMain(m *testing.M) {
	service.Service.Start()
	defer service.Service.Close()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		t1s := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-t1s.C:
				in1 := []byte{0x01, 0x00, 0x00, 0x00, 0x30, 0x00}
				in2 := []byte{0x01, 0x00, 0x01, 0x00, 0x40, 0x01}
				C1C2Handle.HandleUDP(&c1c2.Request{Body: in1}, nil)
				C1C2Handle.HandleUDP(&c1c2.Request{Body: in2}, nil)
			case <-ctx.Done():
				return
			}
		}
	}()
	defer cancel()
	time.Sleep(2 * time.Second)
	for _, item := range items {
		switch item.id {
		case "check-version":
			var bw bytes.Buffer
			bw.WriteByte(0x04)
			ver := object.AutoUpdate.Ver
			bw.WriteByte(byte(ver.Major))
			bw.WriteByte(byte(ver.Minor))
			bw.WriteByte(byte(ver.Patch))
			item.in = bw.Bytes()
		case "get-server":
			server := server.ServerManager.GetServer(1)
			var bw bytes.Buffer
			bw.Write([]byte{0xF4, 0x03})
			var ip [16]byte
			copy(ip[:], server.IP)
			bw.Write(ip[:])
			binary.Write(&bw, binary.LittleEndian, uint16(server.Port))
			item.out = bw.Bytes()
		}
	}
	m.Run()
}

type item struct {
	id  string
	in  []byte
	out []byte
}

var items = [...]*item{
	{"connect", nil, []byte{0x00, 0x01}},
	{"check-version", []byte{0x04, 0x00, 0x01, 0x05}, []byte{0x02, 0x01}},
	{"get-server-list", []byte{0xF4, 0x06}, []byte{0xF4, 0x06, 0x00, 0x02, 0x00, 0x00, 0x30, 0x00, 0x01, 0x00, 0x40, 0x01}},
	{"get-server", []byte{0xF4, 0x03, 0x01, 0x00}, nil},
}

func TestC1C2Handle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	item := items[0]
	items := items[1:]
	conn := c1c2.Conn{RemoteAddr: "test"}
	conn.Write = func(resp *c1c2.Response) error {
		buf := resp.Body()
		// compare
		println(item.id)
		if !bytes.Equal(item.out, buf) {
			t.Errorf("failed [id]%s [expect]%v [get]%v\n", item.id, item.out, buf)
			cancel()
			return nil
		}
		// next
		if len(items) == 0 {
			cancel()
			return nil
		}
		item = items[0]
		items = items[1:]
		req := c1c2.Request{Body: item.in}
		C1C2Handle.Handle(ctx, &req)
		return nil
	}
	conn.Close = func() error {
		log.Println("CloseConn is called")
		return nil
	}
	id, err := C1C2Handle.OnConn(&conn)
	if err != nil {
		t.Error(err)
	}
	ctx = context.WithValue(ctx, c1c2.UserContextKey, id)
	defer C1C2Handle.OnClose(ctx)
	// or:
	// defer func() {
	// 	C1C2Handle.OnClose(ctx)
	// }()
	// ctx = context.WithValue(ctx, c1c2.UserContextKey, id)
	<-ctx.Done()
	err = ctx.Err()
	if err == context.DeadlineExceeded {
		t.Error(err)
	}
}
