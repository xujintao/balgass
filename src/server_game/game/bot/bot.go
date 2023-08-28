package bot

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/xujintao/balgass/src/server_game/game/model"
	"github.com/xujintao/balgass/src/server_game/game/object"
)

func init() {
	BotManager.init()
}

var BotManager botManager

type botManager struct {
	bots         map[string]*bot
	maxBotNumber int
	boter        boter
}

func (m *botManager) init() {
	m.bots = make(map[string]*bot)
}

func (m *botManager) Register(boter boter) {
	m.boter = boter
}

func (m *botManager) AddBot(name string) error {
	b, err := newbot(name, m.boter)
	if err != nil {
		return err
	}
	m.bots[b.name] = b
	return nil
}

func (m *botManager) DeleteBot(name string) {
	if b, ok := m.bots[name]; ok {
		b.close()
		delete(m.bots, name)
	}
}

func (m *botManager) DeleteAllBots() {
	for name, b := range m.bots {
		b.close()
		delete(m.bots, name)
	}
}

func newbot(name string, boter boter) (*bot, error) {
	b := &bot{name: name, boter: boter}
	b.msgChan = make(chan any, 100)
	b.closeConnChan = make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	go func() {
		t100ms := time.NewTicker(time.Millisecond * 100)
		defer func() {
			close(b.msgChan)
			close(b.closeConnChan)
		}()
		for {
			select {
			case msg := <-b.msgChan:
				b.handle(msg)
			case <-b.closeConnChan:
				if b.exit {
					return
				}
				b.state = botStateOffline
			case <-t100ms.C:
				b.stateMachine()
			case <-ctx.Done():
				if b.state < botStateOnline {
					return
				}
				// excute offline and wait close
				b.state = botStateOffline
				b.nextTime = 0
				b.stateMachine()
				b.exit = true
			}
		}
	}()
	return b, nil
}

type boter interface {
	PlayerConn(object.Conn) (int, error)
	PlayerAction(int, string, any)
	PlayerCloseConn(int)
}

type botState int

const (
	botStateInit botState = iota
	botStateOnline
	botStatePickCharacter
	botStatePlay
	botStateOffline
)

type bot struct {
	name string
	boter
	cancel        context.CancelFunc
	id            int
	msgChan       chan any
	closeConnChan chan struct{}
	exit          bool
	state         botState
	curTime       int64
	nextTime      int64
}

type botConn struct {
	*bot
}

func (c *botConn) Addr() string {
	return c.name
}

func (c *botConn) Write(msg any) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("*bot.Write [panic]%v [bot]%s [msg]%v\n", r, c.name, msg)
		}
	}()
	c.msgChan <- msg
	return nil
}

func (c *botConn) Close() error {
	c.closeConnChan <- struct{}{}
	return nil
}

func (b *bot) handle(msg any) {
	name := reflect.TypeOf(msg).String()
	switch name {
	case "*model.MsgConnectSuccess":
		b.state = botStatePickCharacter
	default:
	}
}

func (b *bot) stateMachine() {
	curTime := time.Now().UnixMilli()
	if curTime-b.curTime+1 < b.nextTime {
		return
	}
	b.curTime = curTime
	switch b.state {
	case botStateInit:
		b.nextTime = 1000
		b.state = botStateOnline
	case botStateOnline:
		b.online()
	case botStatePickCharacter:
		b.pickCharacter()
	case botStatePlay:
		b.play()
	case botStateOffline:
		b.offline()
	}
}

func (b *bot) close() {
	b.cancel()
}

func (b *bot) online() {
	c := botConn{b}
	id, err := b.PlayerConn(&c)
	if err != nil {
		return
	}
	b.id = id
	b.nextTime = 2000
}

func (b *bot) pickCharacter() {
	b.PlayerAction(b.id, "PickCharacter", &model.MsgPickCharacter{Name: b.name})
	b.nextTime = 2000
	b.state = botStatePlay // should be set in handle
}

func (b *bot) play() {

}

func (b *bot) offline() {
	b.PlayerCloseConn(b.id)
	b.nextTime = 10 * 1000
	b.state = botStateOnline
}
