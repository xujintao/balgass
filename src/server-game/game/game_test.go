package game

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/xujintao/balgass/src/server-game/game/model"
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
	case "*model.MsgConnectReply":
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
	Game.Start()
	defer Game.Close()

	ctx, cancel := context.WithCancel(context.Background())
	conn := conn{
		t:      t,
		cancel: cancel,
	}
	msg := model.MsgTest{}

	// player
	id, err := Game.PlayerConn(&conn)
	if err != nil {
		t.Error(err)
	}
	Game.PlayerAction(id, "Test", &msg)
	<-ctx.Done()
	Game.PlayerCloseConn(id)

	// user
	id, err = Game.UserConn(&conn)
	if err != nil {
		t.Error(err)
	}
	Game.UserAction(id, "Test", &msg)
	<-ctx.Done()
	Game.UserCloseConn(id)
}
