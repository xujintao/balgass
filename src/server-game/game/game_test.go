package game

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/object"
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
	defer startGameWithoutFixture(t)()

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

func TestPlayerActionPolicyAllow(t *testing.T) {
	tests := []struct {
		name   string
		policy PlayerActionPolicy
		state  object.ConnectState
		want   bool
	}{
		{
			name:   "connection command allows connected state",
			policy: ConnPhase,
			state:  object.ConnectStateConnected,
			want:   true,
		},
		{
			name:   "login command rejects logged state",
			policy: SignPhase,
			state:  object.ConnectStateLogged,
			want:   false,
		},
		{
			name:   "account command rejects connected state",
			policy: AcctPhase,
			state:  object.ConnectStateConnected,
			want:   false,
		},
		{
			name:   "account command allows logged state",
			policy: AcctPhase,
			state:  object.ConnectStateLogged,
			want:   true,
		},
		{
			name:   "account command rejects playing state",
			policy: AcctPhase,
			state:  object.ConnectStatePlaying,
			want:   false,
		},
		{
			name:   "world command rejects logged state",
			policy: PlayPhase,
			state:  object.ConnectStateLogged,
			want:   false,
		},
		{
			name:   "world command allows playing state",
			policy: PlayPhase,
			state:  object.ConnectStatePlaying,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.policy.Allow(tt.state); got != tt.want {
				t.Fatalf("Allow(%v) = %v, want %v", tt.state, got, tt.want)
			}
		})
	}
}

func TestPlayerActionAuth(t *testing.T) {
	defer startGameWithoutFixture(t)()

	ctx, cancel := context.WithCancel(context.Background())
	conn := conn{
		t:      t,
		cancel: cancel,
	}
	id, err := Game.PlayerConn(&conn)
	if err != nil {
		t.Error(err)
		return
	}
	defer Game.PlayerCloseConn(id)

	msg := model.MsgTest{}
	Game.PlayerActionAuth(id, "Test", &msg, PlayPhase)
	select {
	case <-ctx.Done():
		t.Fatal("connected player should not pass PlayPhase")
	case <-time.After(50 * time.Millisecond):
	}

	Game.PlayerActionAuth(id, "Test", &msg, ConnPhase)
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("connected player should pass ConnPhase")
	}
}

func TestPlayerActionWithoutPolicy(t *testing.T) {
	defer startGameWithoutFixture(t)()

	ctx, cancel := context.WithCancel(context.Background())
	conn := conn{
		t:      t,
		cancel: cancel,
	}
	id, err := Game.PlayerConn(&conn)
	if err != nil {
		t.Error(err)
		return
	}
	defer Game.PlayerCloseConn(id)

	msg := model.MsgTest{}
	Game.PlayerAction(id, "Test", &msg)
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("internal player action should be trusted without policy")
	}
}

func startGameWithoutFixture(t *testing.T) func() {
	t.Helper()
	pathCommon := conf.ServerEnv.PathCommon
	conf.ServerEnv.PathCommon = t.TempDir()
	Game.Start()
	return func() {
		Game.Close()
		conf.ServerEnv.PathCommon = pathCommon
	}
}
