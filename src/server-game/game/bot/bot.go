package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync/atomic"

	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/object"
)

func init() {
	BotManager.init()
}

// bot manager
var BotManager botManager

type botManager struct {
	bots         map[string]*bot
	maxBotNumber int
	game         game
}

func (m *botManager) init() {
	m.bots = make(map[string]*bot)
	m.maxBotNumber = 1000
}

func (m *botManager) Register(game game) {
	m.game = game
}

func (m *botManager) AddBot(msg *model.MsgAddBot) (*model.MsgAddBotReply, error) {
	account := strings.TrimSpace(msg.Account)
	password := strings.TrimSpace(msg.Password)
	name := strings.TrimSpace(msg.Name)
	if account == "" || password == "" || name == "" {
		return nil, errors.New("account, password and name are required")
	}

	if m.game == nil {
		return nil, errors.New("bot manager has no registered game")
	}
	if m.maxBotNumber > 0 && len(m.bots) >= m.maxBotNumber {
		return nil, fmt.Errorf("over max bot number %d", m.maxBotNumber)
	}
	key := botKey(account, name)
	if _, ok := m.bots[key]; ok {
		return nil, fmt.Errorf("bot already exists: %s", key)
	}
	b, err := newbot(botConfig{
		key:      key,
		account:  account,
		password: password,
		name:     name,
	}, m.game)
	if err != nil {
		return nil, err
	}
	m.bots[key] = b
	return &model.MsgAddBotReply{Key: key}, nil
}

func (m *botManager) DeleteBot(msg *model.MsgDeleteBot) (*model.MsgDeleteBotReply, error) {
	account := strings.TrimSpace(msg.Account)
	name := strings.TrimSpace(msg.Name)
	if account == "" || name == "" {
		return nil, errors.New("account and name are required")
	}
	key := botKey(account, name)

	b, ok := m.bots[key]
	if ok {
		delete(m.bots, key)
	}
	if ok {
		b.close()
	}
	return &model.MsgDeleteBotReply{Key: key}, nil
}

func (m *botManager) DeleteAllBots() {
	bots := make([]*bot, 0, len(m.bots))
	for key, b := range m.bots {
		bots = append(bots, b)
		delete(m.bots, key)
	}
	for _, b := range bots {
		b.close()
	}
}

func botKey(account, name string) string {
	return strings.ToLower(strings.TrimSpace(account)) + ":" + strings.ToLower(strings.TrimSpace(name))
}

// bot
func newbot(conf botConfig, game game) (*bot, error) {
	if game == nil {
		return nil, errors.New("game is nil")
	}
	b := &bot{
		key:      conf.key,
		account:  conf.account,
		password: conf.password,
		name:     conf.name,
		game:     game,
		msgChan:  make(chan any, 100),
	}
	b.id.Store(-1)
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	go func() {
		defer func() {
			close(b.msgChan)
		}()
		b.connect()
		for {
			select {
			case msg := <-b.msgChan:
				b.handle(msg)
			case <-ctx.Done():
				b.state = botStateOffline
				return
			}
		}
	}()
	return b, nil
}

type botConfig struct {
	key      string
	account  string
	password string
	name     string
}

type game interface {
	PlayerConn(object.Conn) (int, error)
	PlayerAction(int, string, any)
	PlayerCloseConn(int)
	Command(string, any) (any, error)
}

type botState int

const (
	botStateConnecting botState = iota
	botStateConnected
	botStateLoggingIn
	botStateLoggedIn
	botStateLoadingCharacter
	botStatePlaying
	botStateOffline
)

type bot struct {
	key      string
	account  string
	password string
	name     string
	game     game
	cancel   context.CancelFunc
	id       atomic.Int64
	msgChan  chan any
	state    botState
}

func (b *bot) handle(msg any) {
	switch msg := msg.(type) {
	case *model.MsgConnectReply:
		if b.state != botStateConnecting {
			return
		}
		if msg.Result != 1 {
			b.fail(fmt.Errorf("connect failed: result %d", msg.Result))
			return
		}
		b.state = botStateConnected
		b.login()
	case *model.MsgLoginReply:
		if b.state != botStateLoggingIn {
			return
		}
		if msg.Result != 1 {
			b.fail(fmt.Errorf("login failed: result %d", msg.Result))
			return
		}
		b.state = botStateLoggedIn
		b.loadCharacter()
	case *model.MsgLoadCharacterReply:
		if b.state != botStateLoadingCharacter {
			return
		}
		b.state = botStatePlaying
		slog.Info("bot playing", "key", b.key, "player", b.id.Load())
	}
}

func (b *bot) close() {
	id := b.id.Swap(-1)
	if id == -1 {
		return
	}
	go b.game.PlayerCloseConn(int(id))
}

type botConn struct {
	*bot
}

func (c *botConn) Addr() string {
	return c.name
}

func (c *botConn) Write(msg any) error {
	c.msgChan <- msg
	return nil
}

func (c *botConn) Close() error {
	if c.id.Swap(-1) != -1 {
		// if botConn is closed passively, it means the connection is closed by game,
		// we need to delete the bot again without requesting game.PlayerCloseConn.
		if _, err := c.game.Command("DeleteBot", &model.MsgDeleteBot{Account: c.account, Name: c.name}); err != nil {
			slog.Error("bot delete failed", "err", err, "key", c.key)
		}
	}
	c.cancel()
	return nil
}

func (b *bot) connect() {
	b.state = botStateConnecting
	c := botConn{b}
	id, err := b.game.PlayerConn(&c)
	if err != nil {
		b.fail(fmt.Errorf("connect failed: %w", err))
		return
	}
	b.id.Store(int64(id))
}

func (b *bot) login() {
	b.state = botStateLoggingIn
	b.game.PlayerAction(int(b.id.Load()), "Login", &model.MsgLogin{
		Account:  b.account,
		Password: b.password,
	})
}

func (b *bot) loadCharacter() {
	b.state = botStateLoadingCharacter
	b.game.PlayerAction(int(b.id.Load()), "LoadCharacter", &model.MsgLoadCharacter{Name: b.name})
}

func (b *bot) fail(err error) {
	slog.Error("bot failed", "err", err, "key", b.key, "state", b.state)
	b.game.Command("DeleteBot", &model.MsgDeleteBot{Account: b.account, Name: b.name})
}
