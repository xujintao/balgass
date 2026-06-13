package bot

import (
	"errors"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

type testGame struct {
	playerConnErr error
	conns         chan object.Conn
	actions       chan testAction
	closes        chan int
	commands      chan testCommand
}

type testAction struct {
	id     int
	action string
	msg    any
}

type testCommand struct {
	name string
	msg  any
}

func newTestGame() *testGame {
	return &testGame{
		conns:    make(chan object.Conn, 10),
		actions:  make(chan testAction, 10),
		closes:   make(chan int, 10),
		commands: make(chan testCommand, 10),
	}
}

func newTestBotConfig(game game) botConfig {
	return botConfig{
		key:       "account1:char1",
		account:   "account1",
		password:  "password",
		name:      "char1",
		game:      game,
		resources: &resources{},
	}
}

func (g *testGame) PlayerConn(conn object.Conn) (int, error) {
	if g.playerConnErr != nil {
		return -1, g.playerConnErr
	}
	g.conns <- conn
	return 1, nil
}

func (g *testGame) PlayerAction(id int, action string, msg any) {
	g.actions <- testAction{
		id:     id,
		action: action,
		msg:    msg,
	}
}

func (g *testGame) PlayerCloseConn(id int) {
	g.closes <- id
}

func (g *testGame) Command(name string, msg any) (any, error) {
	g.commands <- testCommand{name: name, msg: msg}
	return nil, nil
}

func TestBotManagerDuplicateAndLimit(t *testing.T) {
	g := newTestGame()
	m := botManager{}
	m.init()
	m.maxBotNumber = 1
	m.Register(g)

	if _, err := m.AddBot(&model.MsgAddBot{Account: "account1", Password: "password", Name: "char1"}); err != nil {
		t.Fatalf("AddBot() error = %v", err)
	}
	conn := waitConn(t, g)
	defer func() {
		m.DeleteAllBots()
		_ = conn.Close()
	}()

	if _, err := m.AddBot(&model.MsgAddBot{Account: "account1", Password: "password", Name: "char1"}); err == nil {
		t.Fatal("AddBot() duplicate error = nil")
	}
	if _, err := m.AddBot(&model.MsgAddBot{Account: "account2", Password: "password", Name: "char2"}); err == nil {
		t.Fatal("AddBot() over limit error = nil")
	}
}

func TestBotManagerDeleteMissing(t *testing.T) {
	m := botManager{}
	m.init()
	m.Register(newTestGame())
	reply, err := m.DeleteBot(&model.MsgDeleteBot{Account: "missing", Name: "char"})
	if err != nil {
		t.Fatalf("DeleteBot() error = %v", err)
	}
	if reply.Key != "missing:char" {
		t.Fatalf("DeleteBot() key = %q, want missing:char", reply.Key)
	}
}

func TestBotReachesPlaying(t *testing.T) {
	g := newTestGame()
	b, err := newbot(newTestBotConfig(g))
	if err != nil {
		t.Fatalf("newbot() error = %v", err)
	}

	conn := waitConn(t, g)
	defer func() {
		b.close()
		_ = conn.Close()
	}()

	if err := conn.Write(&model.MsgConnectReply{Result: 1, ID: 1}); err != nil {
		t.Fatalf("Write(connect) error = %v", err)
	}
	login := waitAction(t, g)
	if login.action != "Login" {
		t.Fatalf("login action = %#v", login)
	}
	if msg := login.msg.(*model.MsgLogin); msg.Account != "account1" || msg.Password != "password" {
		t.Fatalf("login msg = %#v", msg)
	}

	if err := conn.Write(&model.MsgLoginReply{Result: 1}); err != nil {
		t.Fatalf("Write(login) error = %v", err)
	}
	load := waitAction(t, g)
	if load.action != "LoadCharacter" {
		t.Fatalf("load action = %#v", load)
	}
	if msg := load.msg.(*model.MsgLoadCharacter); msg.Name != "char1" {
		t.Fatalf("load msg = %#v", msg)
	}

	if err := conn.Write(&model.MsgLoadCharacterReply{}); err != nil {
		t.Fatalf("Write(load character) error = %v", err)
	}
	assertNoCommand(t, g)
	assertNoClose(t, g)
}

func TestBotCloseCallsPlayerCloseConn(t *testing.T) {
	g := newTestGame()
	b, err := newbot(newTestBotConfig(g))
	if err != nil {
		t.Fatalf("newbot() error = %v", err)
	}
	conn := waitConn(t, g)
	waitBotID(t, b, 1)

	b.close()

	closeID := waitClose(t, g)
	if closeID != 1 {
		t.Fatalf("close id = %d, want 1", closeID)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("conn.Close() error = %v", err)
	}
	assertNoCommand(t, g)
}

func TestBotConnCloseCommandsDeleteBot(t *testing.T) {
	g := newTestGame()
	b, err := newbot(newTestBotConfig(g))
	if err != nil {
		t.Fatalf("newbot() error = %v", err)
	}
	conn := waitConn(t, g)
	waitBotID(t, b, 1)

	if err := conn.Close(); err != nil {
		t.Fatalf("conn.Close() error = %v", err)
	}

	command := waitCommand(t, g)
	assertDeleteBotCommand(t, command, "account1", "char1")
	assertNoClose(t, g)
}

func TestBotLoginFailureCommandsDeleteBot(t *testing.T) {
	g := newTestGame()
	_, err := newbot(newTestBotConfig(g))
	if err != nil {
		t.Fatalf("newbot() error = %v", err)
	}

	conn := waitConn(t, g)
	if err := conn.Write(&model.MsgConnectReply{Result: 1, ID: 1}); err != nil {
		t.Fatalf("Write(connect) error = %v", err)
	}
	_ = waitAction(t, g)
	if err := conn.Write(&model.MsgLoginReply{Result: 2}); err != nil {
		t.Fatalf("Write(login failed) error = %v", err)
	}

	command := waitCommand(t, g)
	assertDeleteBotCommand(t, command, "account1", "char1")
}

func TestBotManagerCleansUpAfterConnectFailure(t *testing.T) {
	g := newTestGame()
	g.playerConnErr = errors.New("connect failed")
	m := botManager{}
	m.init()
	m.Register(g)

	if _, err := m.AddBot(&model.MsgAddBot{Account: "account1", Password: "password", Name: "char1"}); err != nil {
		t.Fatalf("AddBot() error = %v", err)
	}
	command := waitCommand(t, g)
	assertDeleteBotCommand(t, command, "account1", "char1")
	if _, err := m.DeleteBot(command.msg.(*model.MsgDeleteBot)); err != nil {
		t.Fatalf("DeleteBot() error = %v", err)
	}
	if len(m.bots) != 0 {
		t.Fatalf("bot count = %d, want 0", len(m.bots))
	}
}

func TestBotLearnsInventorySkillAfterInitialStateArrives(t *testing.T) {
	for _, order := range []string{"items-first", "skills-first"} {
		t.Run(order, func(t *testing.T) {
			g := newTestGame()
			resources := &resources{}
			b := &bot{
				game:   g,
				world:  newWorld(resources),
				policy: newRulePolicy(resources, "account1:char1"),
			}
			b.executor = newExecutor(b)
			b.id.Store(7)
			b.world.phase = PhasePlaying
			b.world.self = Actor{Index: 7, Alive: true}
			b.world.setCharacter(&model.MsgLoadCharacterReply{HP: 100})
			b.world.mergeSelf(Actor{
				Class: int(class.Wizard),
				Level: 1,
				HP:    100,
			})
			base := &item.ItemBase{
				KindA:      item.KindASkill,
				SkillIndex: skill.SkillIndexFireBall,
			}
			base.ReqClass[class.Wizard] = 1
			book := &item.Item{
				ItemBase: base,
				Code:     item.Code(15, 0),
				ID:       123,
			}
			items := make([]*item.Item, 13)
			items[12] = book
			now := time.Unix(0, 0)

			if order == "items-first" {
				b.world.Handle(&model.MsgItemListReply{Items: items})
				b.tick(now)
				assertNoAction(t, g)
				b.world.Handle(&model.MsgSkillListReply{})
			} else {
				b.world.Handle(&model.MsgSkillListReply{})
				b.tick(now)
				assertNoAction(t, g)
				b.world.Handle(&model.MsgItemListReply{Items: items})
			}

			b.tick(now.Add(100 * time.Millisecond))
			action := waitAction(t, g)
			if action.action != "UseItem" {
				t.Fatalf("action = %#v, want UseItem", action)
			}
			msg := action.msg.(*model.MsgUseItem)
			if msg.SrcPosition != 12 || msg.DstPosition != 0 {
				t.Fatalf("use item msg = %#v, want source 12", msg)
			}

			b.tick(now.Add(200 * time.Millisecond))
			assertNoAction(t, g)
			b.world.Handle(&model.MsgSkillOneReply{
				Flag: -2,
				Skill: &skill.Skill{
					SkillBase: &skill.SkillBase{Index: skill.SkillIndexFireBall},
					Index:     skill.SkillIndexFireBall,
				},
			})
			b.tick(now.Add(300 * time.Millisecond))
			assertNoAction(t, g)
		})
	}
}

func waitConn(t *testing.T, g *testGame) object.Conn {
	t.Helper()
	select {
	case conn := <-g.conns:
		return conn
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for PlayerConn")
	}
	return nil
}

func waitAction(t *testing.T, g *testGame) testAction {
	t.Helper()
	select {
	case action := <-g.actions:
		return action
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for PlayerAction")
	}
	return testAction{}
}

func waitClose(t *testing.T, g *testGame) int {
	t.Helper()
	select {
	case id := <-g.closes:
		return id
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for PlayerCloseConn")
	}
	return 0
}

func waitCommand(t *testing.T, g *testGame) testCommand {
	t.Helper()
	select {
	case command := <-g.commands:
		return command
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for Command")
	}
	return testCommand{}
}

func waitBotID(t *testing.T, b *bot, id int) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if int(b.id.Load()) == id {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("bot id = %d, want %d", b.id.Load(), id)
}

func assertDeleteBotCommand(t *testing.T, command testCommand, account, name string) {
	t.Helper()
	if command.name != "DeleteBot" {
		t.Fatalf("command name = %q, want DeleteBot", command.name)
	}
	msg := command.msg.(*model.MsgDeleteBot)
	if msg.Account != account || msg.Name != name {
		t.Fatalf("delete bot msg = %#v", msg)
	}
}

func assertNoCommand(t *testing.T, g *testGame) {
	t.Helper()
	select {
	case command := <-g.commands:
		t.Fatalf("unexpected command = %#v", command)
	case <-time.After(50 * time.Millisecond):
	}
}

func assertNoAction(t *testing.T, g *testGame) {
	t.Helper()
	select {
	case action := <-g.actions:
		t.Fatalf("unexpected PlayerAction = %#v", action)
	case <-time.After(50 * time.Millisecond):
	}
}

func assertNoClose(t *testing.T, g *testGame) {
	t.Helper()
	select {
	case id := <-g.closes:
		t.Fatalf("unexpected PlayerCloseConn id = %d", id)
	case <-time.After(50 * time.Millisecond):
	}
}
