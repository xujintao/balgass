package handle

import (
	"bytes"
	"context"
	"encoding/hex"
	"log"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server_game/game"
	"github.com/xujintao/balgass/src/server_game/game/model"
)

func TestMarshal(t *testing.T) {
	tm := model.MsgConnectFailed{Result: 1}
	_, err := C1C2Handle.marshal(&tm)
	if err != nil {
		t.Error(err)
	}
}

func TestMain(m *testing.M) {
	game.Game.Start()
	defer game.Game.Close()
	m.Run()
}

type item struct {
	id  string
	in  string
	out string
}

var items = [...]*item{
	{"connect", "", ""},
	{"login", "f101cdfd98c8faabfccfabfccdfd98c8faabfccfabfccfabfccfabfccfabfccf000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007dfaa614302e312e350000004d374234564d3443356938424334396240000000", ""},
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
		expect, err := hex.DecodeString(item.out)
		if err != nil {
			t.Errorf("hex.DecodeString failed [id]%s [out]%s\n", item.id, item.out)
			cancel()
			return nil
		}
		if len(expect) != 0 && !bytes.Equal(expect, buf) {
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
		in, err := hex.DecodeString(item.in)
		if err != nil {
			t.Errorf("hex.DecodeString failed [id]%s [in]%s\n", item.id, item.in)
			cancel()
			return nil
		}
		req := c1c2.Request{Body: in}
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
	<-ctx.Done()
	err = ctx.Err()
	if err == context.DeadlineExceeded {
		t.Error(err)
	}
}
