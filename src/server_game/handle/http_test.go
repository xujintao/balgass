package handle

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/xujintao/balgass/src/server_game/game"
	"github.com/xujintao/balgass/src/server_game/game/model"
)

type testWSConn struct {
	t      *testing.T
	cancel context.CancelFunc
}

func (c *testWSConn) Addr() string {
	return "testWSConn"
}

func (c *testWSConn) Write(msg any) error {
	reply, err := wsHandleDefault.addAction(msg)
	if err != nil {
		c.t.Errorf("addAction failed [err]%v\n", err)
		return err
	}
	if reply["action"] == testData[len(testData)-1].out.action {
		c.cancel()
	}
	return nil
}

func (c *testWSConn) Close() error {
	return nil
}

type testItem struct {
	in  wsIn
	out wsOut
}

var testData = [...]testItem{
	{wsIn{"SubscribeMap", &model.MsgSubscribeMap{Number: 0}},
		wsOut{"SubscribeMapReply", &model.MsgSubscribeMapReply{}}},
}

func TestWSHandle(t *testing.T) {
	game.Game.Start()
	defer game.Game.Close()
	ctx, cancel := context.WithCancel(context.Background())
	conn := testWSConn{t: t, cancel: cancel}
	id, err := game.Game.UserConn(&conn)
	if err != nil {
		conn.Close()
		t.Error(err)
		return
	}
	defer game.Game.UserCloseConn(id)

	for _, item := range testData {
		in := item.in
		action := in.action
		data, err := json.Marshal(in.msg)
		if err != nil {
			t.Error(err)
			break
		}
		wsHandleDefault.Handle(id, action, data)
	}
	<-ctx.Done()
}
