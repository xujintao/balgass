package bot

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
)

var (
	policyTraceLogger *slog.Logger
)

const defaultPolicyTraceFile = "/tmp/server-game-bot-policy.jsonl"

func init() {
	initPolicyTrace()
}

type policyTraceActor struct {
	Index     int  `json:"index"`
	Class     int  `json:"class"`
	MapNumber int  `json:"map"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Dir       int  `json:"dir"`
	Alive     bool `json:"alive"`
}

type policyTraceExecution struct {
	CurrentAction   string              `json:"current_action"`
	Position        Position            `json:"position"`
	Dir             int                 `json:"dir"`
	MoveActive      bool                `json:"move_active"`
	MoveTarget      Position            `json:"move_target,omitempty"`
	ReadyAt         string              `json:"ready_at,omitempty"`
	PositionVersion uint64              `json:"position_version"`
	Move            policyTraceMoveInfo `json:"move,omitempty"`
}

type policyTraceMoveInfo struct {
	PathLen    int    `json:"path_len"`
	PathNext   int    `json:"path_next"`
	NextStepAt string `json:"next_step_at,omitempty"`
}

type policyTraceAction struct {
	Kind            string   `json:"kind"`
	Target          int      `json:"target,omitempty"`
	Skill           int      `json:"skill,omitempty"`
	Position        int      `json:"position,omitempty"`
	PathLen         int      `json:"path_len,omitempty"`
	PathEnd         Position `json:"path_end,omitempty"`
	Dir             int      `json:"dir,omitempty"`
	CancelCurrent   bool     `json:"cancel_current,omitempty"`
	PositionVersion uint64   `json:"position_version,omitempty"`
	Text            string   `json:"text,omitempty"`
}

func tracePolicyDecision(bot string, world WorldSnapshot, execution ExecutorSnapshot, reason string, action Action, extra map[string]interface{}) {
	if policyTraceLogger == nil {
		return
	}
	policyTraceLogger.Debug("bot policy decision",
		"bot", bot,
		"phase", phaseName(world.Phase),
		"reason", reason,
		"self", traceActor(world.Self),
		"execution", traceExecution(execution),
		"action", traceAction(action),
		"extra", extra,
	)
}

func initPolicyTrace() {
	policyTraceLogger = nil
	if !conf.ServerEnv.TraceBotPolicyEnable {
		return
	}
	file := strings.TrimSpace(conf.ServerEnv.TraceBotPolicyFile)
	if file == "" {
		file = defaultPolicyTraceFile
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		slog.Error("open bot policy trace file", "file", file, "err", err)
		return
	}
	policyTraceLogger = slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func traceActor(actor Actor) policyTraceActor {
	return policyTraceActor{
		Index:     actor.Index,
		Class:     actor.Class,
		MapNumber: actor.MapNumber,
		X:         actor.X,
		Y:         actor.Y,
		Dir:       actor.Dir,
		Alive:     actor.Alive,
	}
}

func traceExecution(execution ExecutorSnapshot) policyTraceExecution {
	info := policyTraceExecution{
		CurrentAction:   actionKindName(execution.CurrentAction.Kind),
		Position:        execution.Position,
		Dir:             execution.Dir,
		MoveActive:      execution.Move.Active,
		MoveTarget:      execution.Move.Target,
		PositionVersion: execution.PositionVersion,
	}
	if !execution.ReadyAt.IsZero() {
		info.ReadyAt = execution.ReadyAt.Format(time.RFC3339Nano)
	}
	if execution.Move.Active {
		info.Move = policyTraceMoveInfo{
			PathLen:  len(execution.Move.Path),
			PathNext: execution.Move.PathNext,
		}
		if !execution.Move.NextStepAt.IsZero() {
			info.Move.NextStepAt = execution.Move.NextStepAt.Format(time.RFC3339Nano)
		}
	}
	return info
}

func traceAction(action Action) policyTraceAction {
	info := policyTraceAction{
		Kind:            actionKindName(action.Kind),
		Target:          action.Target,
		Skill:           action.Skill,
		Position:        action.Position,
		Dir:             action.Dir,
		CancelCurrent:   action.CancelCurrent,
		PositionVersion: action.PositionVersion,
		Text:            action.Text,
	}
	if len(action.Path) > 0 {
		info.PathLen = len(action.Path)
		info.PathEnd = action.Path[len(action.Path)-1]
	}
	return info
}

func actionKindName(kind ActionKind) string {
	switch kind {
	case ActionNone:
		return "none"
	case ActionConnect:
		return "connect"
	case ActionLogin:
		return "login"
	case ActionLoadCharacter:
		return "load_character"
	case ActionStop:
		return "stop"
	case ActionLearnSkill:
		return "learn_skill"
	case ActionMove:
		return "move"
	case ActionContinueMove:
		return "continue_move"
	case ActionSyncPosition:
		return "sync_position"
	case ActionCancel:
		return "cancel"
	case ActionAttack:
		return "attack"
	case ActionUseSkill:
		return "use_skill"
	case ActionChat:
		return "chat"
	case ActionWhisper:
		return "whisper"
	default:
		return "unknown"
	}
}

func phaseName(phase Phase) string {
	switch phase {
	case PhaseDisconnected:
		return "disconnected"
	case PhaseConnected:
		return "connected"
	case PhaseLoggedIn:
		return "logged_in"
	case PhasePlaying:
		return "playing"
	case PhaseFailed:
		return "failed"
	default:
		return "unknown"
	}
}
