package bot

import (
	"errors"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
)

type ExecutorSnapshot struct {
	CurrentAction   Action
	Position        Position
	Dir             int
	ReadyAt         time.Time
	PositionVersion uint64
	Move            MoveSnapshot
}

type MoveSnapshot struct {
	Active     bool
	Path       []Position
	PathNext   int
	SentAt     time.Time
	NextStepAt time.Time
	Target     Position
}

type moveExecution struct {
	path       []Position
	pathNext   int
	sentAt     time.Time
	nextStepAt time.Time
	target     Position
}

type executor struct {
	bot             *bot
	current         Action
	position        Position
	dir             int
	readyAt         time.Time
	positionVersion uint64
	move            moveExecution
}

func newExecutor(bot *bot) *executor {
	return &executor{bot: bot}
}

func (e *executor) Snapshot() ExecutorSnapshot {
	snapshot := ExecutorSnapshot{
		CurrentAction:   cloneAction(e.current),
		Position:        e.position,
		Dir:             e.dir,
		ReadyAt:         e.readyAt,
		PositionVersion: e.positionVersion,
	}
	if len(e.move.path) > 0 {
		snapshot.Move = MoveSnapshot{
			Active:     true,
			Path:       append([]Position(nil), e.move.path...),
			PathNext:   e.move.pathNext,
			SentAt:     e.move.sentAt,
			NextStepAt: e.move.nextStepAt,
			Target:     e.move.target,
		}
	}
	return snapshot
}

func (e *executor) Execute(action Action) {
	switch action.Kind {
	case ActionNone:
	case ActionConnect:
		e.connect(action)
	case ActionLogin:
		e.login(action)
	case ActionLoadCharacter:
		e.loadCharacter(action)
	case ActionStop:
		e.stop(action)
	case ActionLearnSkill:
		e.learnSkill(action)
	case ActionMove:
		e.startMove(action)
	case ActionContinueMove:
		e.continueMove(action)
	case ActionSyncPosition:
		e.syncPosition(action)
	case ActionCancel:
		e.cancel()
	case ActionAttack:
		e.attack(action)
	case ActionUseSkill:
		e.useSkill(action)
	case ActionChat:
		e.chat(action)
	case ActionWhisper:
		e.whisper(action)
	}
}

func (e *executor) connect(action Action) {
	e.current = cloneAction(action)
	id, err := e.bot.game.PlayerConn(&botConn{bot: e.bot})
	if err != nil {
		e.bot.world.Handle(&connectFailed{Err: err})
		return
	}
	e.bot.id.Store(int64(id))
}

func (e *executor) login(action Action) {
	e.current = cloneAction(action)
	e.bot.game.PlayerAction(int(e.bot.id.Load()), "Login", &model.MsgLogin{
		Account:  e.bot.account,
		Password: e.bot.password,
	})
}

func (e *executor) loadCharacter(action Action) {
	e.current = cloneAction(action)
	e.bot.game.PlayerAction(
		int(e.bot.id.Load()),
		"LoadCharacter",
		&model.MsgLoadCharacter{Name: e.bot.name},
	)
}

func (e *executor) stop(action Action) {
	e.current = cloneAction(action)
	e.bot.fail(errors.New(action.Text))
}

func (e *executor) learnSkill(action Action) {
	e.current = cloneAction(action)
	dst := 0
	if action.Position == dst {
		dst = 1
	}
	e.bot.game.PlayerAction(int(e.bot.id.Load()), "UseItem", &model.MsgUseItem{
		SrcPosition: action.Position,
		DstPosition: dst,
	})
}

func (e *executor) startMove(action Action) {
	path := make(maps.Path, len(action.Path))
	dirs := make([]int, len(action.Path))
	from := action.SelfPosition
	for i, pos := range action.Path {
		path[i] = &maps.Pot{X: pos.X, Y: pos.Y}
		dirs[i] = calcDir(from, pos)
		from = pos
	}
	e.current = cloneAction(action)
	e.position = action.SelfPosition
	e.dir = action.Dir
	e.positionVersion = action.PositionVersion
	e.move = moveExecution{
		path:       append([]Position(nil), action.Path...),
		pathNext:   action.PathNext,
		sentAt:     action.SentAt,
		nextStepAt: action.NextStepAt,
		target:     action.Path[len(action.Path)-1],
	}
	e.bot.game.PlayerAction(int(e.bot.id.Load()), "Move", &model.MsgMove{
		X:    action.SelfPosition.X,
		Y:    action.SelfPosition.Y,
		Dir:  action.Dir,
		Dirs: dirs,
		Path: path,
	})
}

func (e *executor) continueMove(action Action) {
	e.position = action.SelfPosition
	e.dir = action.Dir
	e.move.pathNext = action.PathNext
	e.move.nextStepAt = action.NextStepAt
	if e.move.pathNext >= len(e.move.path) {
		e.move = moveExecution{}
		e.current = Action{}
	}
}

func (e *executor) syncPosition(action Action) {
	e.position = action.SelfPosition
	e.dir = action.Dir
	e.positionVersion = action.PositionVersion
	e.move = moveExecution{}
	if action.CancelCurrent {
		e.current = Action{}
	}
}

func (e *executor) cancel() {
	e.current = Action{}
	e.move = moveExecution{}
}

func (e *executor) attack(action Action) {
	e.current = cloneAction(action)
	e.dir = action.Dir
	e.readyAt = action.ReadyAt
	e.bot.game.PlayerAction(int(e.bot.id.Load()), "Attack", &model.MsgAttack{
		Target: action.Target,
		Action: 120,
		Dir:    action.Dir,
	})
}

func (e *executor) useSkill(action Action) {
	e.current = cloneAction(action)
	e.dir = action.Dir
	e.readyAt = action.ReadyAt
	e.bot.game.PlayerAction(int(e.bot.id.Load()), "UseSkill", &model.MsgUseSkill{
		Target: action.Target,
		Skill:  action.Skill,
	})
}

func (e *executor) chat(action Action) {
	e.bot.game.PlayerAction(int(e.bot.id.Load()), "Chat", &model.MsgChat{
		Name: e.bot.name,
		Msg:  action.Text,
	})
	e.current = Action{}
}

func (e *executor) whisper(action Action) {
	e.bot.game.PlayerAction(int(e.bot.id.Load()), "Whisper", &model.MsgWhisper{
		MsgChat: model.MsgChat{Name: action.Name, Msg: action.Text},
	})
	e.current = Action{}
}

func cloneAction(action Action) Action {
	action.Path = append([]Position(nil), action.Path...)
	return action
}
