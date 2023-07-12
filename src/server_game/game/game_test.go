package game

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/xujintao/balgass/src/server_game/game/model"
)

type conn struct {
	t      *testing.T
	cancel context.CancelFunc
}

func (c *conn) Addr() string {
	return "test"
}

func (c *conn) Write(msg any) error {
	name := reflect.TypeOf(msg).String()
	switch name {
	case "*model.MsgConnectSuccess":
		fmt.Println("connect succeeded")
	case "*model.MsgTest":
		c.cancel()
	default:
		c.t.Error("invalid msg")
	}
	return nil
}

func (c *conn) Close() error {
	return nil
}

func TestGame(t *testing.T) {
	Start()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	conn := conn{
		t:      t,
		cancel: cancel,
	}
	id, err := Conn(&conn)
	if err != nil {
		t.Error(err)
	}

	msg := model.MsgTest{}
	PlayerAction(id, "Test", &msg)
	<-ctx.Done()
	CloseConn(id)

	Close()
}
