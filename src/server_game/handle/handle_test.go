package handle

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server_game/game"
	"github.com/xujintao/balgass/src/server_game/game/model"
)

func TestMarshal(t *testing.T) {
	tm := model.MsgConnectFailed{Result: 1}
	_, err := APIHandleDefault.Marshal(&tm)
	if err != nil {
		t.Error(err)
	}
}

func TestHandle(t *testing.T) {
	game.Start()
	testData := []byte{0xFF, 0xFF, 0x01}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	conn := c1c2.Conn{RemoteAddr: "test"}
	conn.Write = func(resp *c1c2.Response) error {
		buf := resp.Body()
		if buf[0] != 0xF1 {
			if !reflect.DeepEqual(buf, testData) {
				t.Error()
			}
			cancel()
		}
		return nil
	}
	conn.Close = func() error {
		log.Println("CloseConn is called")
		return nil
	}
	id, err := APIHandleDefault.OnConn(&conn)
	if err != nil {
		t.Error(err)
	}
	ctx = context.WithValue(ctx, c1c2.UserContextKey, id)
	req := c1c2.Request{
		Body: testData,
	}
	APIHandleDefault.Handle(ctx, &req)
	<-ctx.Done()
	APIHandleDefault.OnClose(ctx)
	game.Close()
}
