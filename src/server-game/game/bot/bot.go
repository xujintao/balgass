package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/maps"
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
	apis         map[any]*api
	resources    *resources
}

func (m *botManager) init() {
	m.bots = make(map[string]*bot)
	m.maxBotNumber = 1000
	m.apis = make(map[any]*api)
	for _, v := range apis {
		t := reflect.TypeOf(v.msg)
		if t.Kind() != reflect.Ptr {
			slog.Error("api msg field must be a pointer",
				"handle", v.handle)
			os.Exit(1)
		}
		m.apis[t] = v
	}
	resources, err := getDefaultResources()
	if err != nil {
		panic(fmt.Errorf("load bot resources: %w", err))
	}
	m.resources = resources
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
		key:       key,
		account:   account,
		password:  password,
		name:      name,
		game:      m.game,
		apis:      m.apis,
		resources: m.resources,
	})
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
func newbot(conf botConfig) (*bot, error) {
	if conf.game == nil {
		return nil, errors.New("game is nil")
	}
	b := &bot{
		key:      conf.key,
		account:  conf.account,
		password: conf.password,
		name:     conf.name,
		game:     conf.game,
		apis:     conf.apis,
		msgChan:  make(chan any, 100),
		world:    newWorld(conf.resources),
		policy:   newRulePolicy(conf.resources, conf.key),
	}
	b.id.Store(-1)
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	go func() {
		defer func() {
			close(b.msgChan)
		}()
		b.connect()
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case msg := <-b.msgChan:
				b.handle(msg)
			case now := <-ticker.C:
				b.tick(now)
			case <-ctx.Done():
				return
			}
		}
	}()
	return b, nil
}

type botConfig struct {
	key       string
	account   string
	password  string
	name      string
	game      game
	apis      map[any]*api
	resources *resources
}

type game interface {
	PlayerConn(object.Conn) (int, error)
	PlayerAction(int, string, any)
	PlayerCloseConn(int)
	Command(string, any) (any, error)
}

type bot struct {
	key          string
	account      string
	password     string
	name         string
	game         game
	apis         map[any]*api
	cancel       context.CancelFunc
	id           atomic.Int64
	msgChan      chan any
	world        *world
	policy       Policy
	action       actionExecution
	nextDecision time.Time
}

type actionExecutionKind int

const (
	actionExecutionIdle actionExecutionKind = iota
	actionExecutionMoving
	actionExecutionAttacking
)

const defaultAttackInterval = 800 * time.Millisecond

type actionExecution struct {
	kind       actionExecutionKind
	target     int
	nextAttack time.Time
}

func (b *bot) handle(msg any) {
	if b.id.Load() < 0 || msg == nil {
		return
	}
	t := reflect.TypeOf(msg)
	api, ok := b.apis[t]
	if !ok {
		// slog.Error("bot received unregistered msg", "key", b.key, "msg", t.String())
		return
	}
	handler := reflect.ValueOf(b).MethodByName(api.handle)
	if !handler.IsValid() {
		handler = reflect.ValueOf(b.world).MethodByName(api.handle)
	}
	if !handler.IsValid() {
		slog.Error("bot has no method to handle msg", "key", b.key, "handle", api.handle)
		return
	}
	handler.Call([]reflect.Value{reflect.ValueOf(msg)})
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
	c := botConn{b}
	id, err := b.game.PlayerConn(&c)
	if err != nil {
		b.fail(fmt.Errorf("connect failed: %w", err))
		return
	}
	b.id.Store(int64(id))
}

func (b *bot) fail(err error) {
	slog.Error("bot failed", "err", err, "key", b.key)
	b.game.Command("DeleteBot", &model.MsgDeleteBot{Account: b.account, Name: b.name})
}

func (b *bot) HandleConnectReply(msg *model.MsgConnectReply) {
	if msg.Result != 1 {
		b.fail(fmt.Errorf("connect failed: result %d", msg.Result))
		return
	}
	b.world.setSelfIndex(msg.ID)
	b.game.PlayerAction(int(b.id.Load()), "Login", &model.MsgLogin{
		Account:  b.account,
		Password: b.password,
	})
}

func (b *bot) HandleLoginReply(msg *model.MsgLoginReply) {
	if msg.Result != 1 {
		b.fail(fmt.Errorf("login failed: result %d", msg.Result))
		return
	}
	b.game.PlayerAction(int(b.id.Load()), "LoadCharacter", &model.MsgLoadCharacter{Name: b.name})
}

func (b *bot) HandleLoadCharacterReply(msg *model.MsgLoadCharacterReply) {
	b.world.setCharacter(msg)
	slog.Info("bot playing", "key", b.key, "player", b.id.Load())
	b.nextDecision = time.Now().Add(b.decisionDelay())
}

type api struct {
	msg    any
	handle string
}

var apis = [...]*api{
	{(*model.MsgConnectReply)(nil), "HandleConnectReply"},
	{(*model.MsgLoginReply)(nil), "HandleLoginReply"},
	{(*model.MsgLoadCharacterReply)(nil), "HandleLoadCharacterReply"},
	{(*model.MsgCreateViewportPlayerReply)(nil), "HandleCreateViewportPlayerReply"},
	{(*model.MsgCreateViewportMonsterReply)(nil), "HandleCreateViewportMonsterReply"},
	{(*model.MsgDestroyViewportObjectReply)(nil), "HandleDestroyViewportObjectReply"},
	{(*model.MsgMoveReply)(nil), "HandleMoveReply"},
	{(*model.MsgSetPositionReply)(nil), "HandleSetPositionReply"},
	{(*model.MsgAttackHPReply)(nil), "HandleAttackHPReply"},
	{(*model.MsgAttackDieReply)(nil), "HandleAttackDieReply"},
	{(*model.MsgTeleportReply)(nil), "HandleTeleportReply"},
	{(*model.MsgReloadCharacterReply)(nil), "HandleReloadCharacterReply"},
	{(*model.MsgHPReply)(nil), "HandleHPReply"},
}

func (b *bot) tick(now time.Time) {
	if !b.world.playing {
		b.action = actionExecution{}
		return
	}
	if b.advanceAction(now) {
		return
	}
	if now.Before(b.nextDecision) {
		return
	}
	b.nextDecision = now.Add(b.decisionDelay())
	b.execute(b.policy.Decide(b.world.snapshot(now)), now)
}

func (b *bot) decisionDelay() time.Duration {
	return 500 * time.Millisecond
}

func (b *bot) advanceAction(now time.Time) bool {
	switch b.action.kind {
	case actionExecutionMoving:
		b.world.advance(now)
		if !b.world.self.Alive || !b.world.moving() {
			b.action = actionExecution{}
		}
		return true
	case actionExecutionAttacking:
		target, ok := b.world.objects[b.action.target]
		if !b.world.self.Alive ||
			!ok ||
			!target.Alive ||
			!target.Attackable ||
			pathDistance(b.world.self.position(), target.position()) > 1 {
			b.action = actionExecution{}
			return true
		}
		if now.Before(b.action.nextAttack) {
			return true
		}
		if b.attack(target) {
			b.action.nextAttack = now.Add(defaultAttackInterval)
		} else {
			b.action = actionExecution{}
		}
		return true
	default:
		return false
	}
}

func (b *bot) execute(action Action, now time.Time) {
	id := int(b.id.Load())
	if id < 0 {
		return
	}
	switch action.Kind {
	case ActionMove:
		if len(action.Path) == 0 || b.world.moving() {
			return
		}
		path := make(maps.Path, len(action.Path))
		dirs := make([]int, len(action.Path))
		from := b.world.self.position()
		for i, pos := range action.Path {
			path[i] = &maps.Pot{X: pos.X, Y: pos.Y}
			dirs[i] = calcDir(from, pos)
			from = pos
		}
		b.world.startMove(action.Path, now)
		b.game.PlayerAction(id, "Move", &model.MsgMove{
			X:    b.world.self.X,
			Y:    b.world.self.Y,
			Dir:  b.world.self.Dir,
			Dirs: dirs,
			Path: path,
		})
		b.action = actionExecution{kind: actionExecutionMoving}
	case ActionAttack:
		target, ok := b.world.objects[action.Target]
		if !ok ||
			!target.Alive ||
			!target.Attackable ||
			pathDistance(b.world.self.position(), target.position()) > 1 {
			return
		}
		if b.attack(target) {
			b.action = actionExecution{
				kind:       actionExecutionAttacking,
				target:     action.Target,
				nextAttack: now.Add(defaultAttackInterval),
			}
		}
	case ActionChat:
		b.game.PlayerAction(id, "Chat", &model.MsgChat{Name: b.name, Msg: action.Text})
	case ActionWhisper:
		b.game.PlayerAction(id, "Whisper", &model.MsgWhisper{MsgChat: model.MsgChat{Name: action.Name, Msg: action.Text}})
	}
}

func (b *bot) attack(target Actor) bool {
	id := int(b.id.Load())
	if id < 0 {
		return false
	}
	dir := calcDir(b.world.self.position(), target.position())
	b.world.self.Dir = dir
	b.game.PlayerAction(id, "Attack", &model.MsgAttack{
		Target: target.Index,
		Action: 120,
		Dir:    dir,
	})
	return true
}
