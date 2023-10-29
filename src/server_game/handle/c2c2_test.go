package handle

import (
	"bytes"
	"context"
	"encoding/hex"
	"log"
	"testing"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server_game/game"
)

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
	{"login1", "f10188aad888feabfccfabfc88aad888cfabfccfabfccfabfccfabfccfabfccf000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000adc0600302e312e350000004d374234564d3443356938424334396240000000", "F10102"},
	{"login2", "f10188aad888cfabfccfabfc88aad888feabfccfabfccfabfccfabfccfabfccf00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000675b0800302e312e350000004d374234564d3443356938424334396240000000", "F10100"},
	{"login3", "f10188aad888cfabfccfabfc88aad888cfabfccfabfccfabfccfabfccfabfccf000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006dc70200302e312e350000004d374234564d3443356938424334396240000000", "F10101"},
	// {"create-character", "F301726F6C6531000000000010", "f30101726f6c653100000000000001000100000000000000000000000000000000000000000000000000"},
	{"get-character-list", "F300", "f300ff00010100726f6c6531000000000000010000200d0d00000d000000010000000000000000ff0000"},
	// {"delete-character", "F302726F6C6531000000000088aad888cfabfccfabfccfabfccfabfccfabfccf", "F30201"},
}

func TestC1C2Handle(t *testing.T) {
	// panic: Fail in goroutine after TestC1C2Handle has completed
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	item := items[0]
	items := items[1:]
	// log.Println(item.id)
	conn := c1c2.Conn{RemoteAddr: "test"}
	conn.Write = func(resp *c1c2.Response) error {
		buf := resp.Body()

		// compare
		expect, err := hex.DecodeString(item.out)
		if err != nil {
			t.Errorf("hex.DecodeString failed [id]%s [out]%s\n", item.id, item.out)
			cancel()
			return nil
		}
		if len(expect) != 0 && !bytes.Equal(expect, buf) {
			t.Errorf("failed [id]%s [expect]%s [get]%s\n", item.id, item.out, hex.EncodeToString(buf))
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
		// log.Println(item.id)
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
		return
	}
	ctx = context.WithValue(ctx, c1c2.UserContextKey, id)
	defer C1C2Handle.OnClose(ctx)
	<-ctx.Done()
	err = ctx.Err()
	if err == context.DeadlineExceeded {
		t.Error(err)
		return
	}
}
