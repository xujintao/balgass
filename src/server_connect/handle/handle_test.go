package handle

import (
	"bytes"
	"context"
	"log"
	"testing"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server_connect/service"
)

type item struct {
	id  string
	in  []byte
	out []byte
}

var items = [...]*item{
	{"connect", nil, []byte{0x00, 0x01}},
	{"check version", []byte{0x04, 0x00, 0x01, 0x05}, []byte{0x02, 0x01}},
}

func TestC1C2Handle(t *testing.T) {
	service.Service.Start()
	defer service.Service.Close()
	ctx, cancel := context.WithCancel(context.Background())
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
}
