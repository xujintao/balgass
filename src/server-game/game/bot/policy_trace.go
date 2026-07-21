package bot

import (
	"log/slog"
	"os"
	"strings"

	"github.com/xujintao/balgass/src/server-game/conf"
)

var (
	policyTraceLogger *slog.Logger
)

const defaultPolicyTraceFile = "/tmp/server-game-bot-policy.jsonl"

func init() {
	initPolicyTrace()
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

func tracePolicyDecision(bot string, world WorldSnapshot, execution ExecutorSnapshot, reason string, action Action, extra map[string]interface{}) {
	if policyTraceLogger == nil {
		return
	}
	args := []any{
		"bot", bot,
		"src", "policy",
		"phase", phaseName(world.Phase),
		"reason", reason,
		"map", world.Self.MapNumber,
		"x", world.Self.X,
		"y", world.Self.Y,
		"dir", world.Self.Dir,
		"alive", world.Self.Alive,
		"exec_x", execution.Position.X,
		"exec_y", execution.Position.Y,
		"exec_dir", execution.Dir,
		"exec_v", execution.PositionVersion,
		"exec_action", actionKindName(execution.CurrentAction.Kind),
		"act", actionKindName(action.Kind),
	}
	if len(action.Path) > 0 {
		end := action.Path[len(action.Path)-1]
		args = append(args,
			"act_len", len(action.Path),
			"act_end_x", end.X,
			"act_end_y", end.Y,
		)
	}
	if action.Target != 0 {
		args = append(args, "act_target", action.Target)
	}
	if action.Skill != 0 {
		args = append(args, "act_skill", action.Skill)
	}
	if action.Dir != 0 {
		args = append(args, "act_dir", action.Dir)
	}
	if action.PositionVersion != 0 {
		args = append(args, "act_v", action.PositionVersion)
	}

	if len(extra) > 0 {
		if pathLen, ok := extra["path_len"]; ok {
			args = append(args, "full_len", pathLen)
		}
		if pathEnd, ok := extra["path_end"].(Position); ok {
			args = append(args,
				"full_end_x", pathEnd.X,
				"full_end_y", pathEnd.Y,
			)
		}
		if target, ok := extra["target"].(Actor); ok {
			args = append(args,
				"target_index", target.Index,
				"target_class", target.Class,
				"target_x", target.X,
				"target_y", target.Y,
			)
		}
		if targets, ok := extra["targets"]; ok {
			args = append(args, "targets", targets)
		}
		if attackRange, ok := extra["attack_range"]; ok {
			args = append(args, "attack_range", attackRange)
		}
		if skill, ok := extra["skill"]; ok {
			args = append(args, "skill", skill)
		}
		if worldPositionVersion, ok := extra["world_position_version"]; ok {
			args = append(args, "world_v", worldPositionVersion)
		}
		if learnSkills, ok := extra["learn_skills"]; ok {
			args = append(args, "learn_skills", learnSkills)
		}
		if position, ok := extra["position"]; ok {
			args = append(args, "position", position)
		}
		if failure, ok := extra["failure"]; ok {
			args = append(args, "failure", failure)
		}
	}
	policyTraceLogger.Debug("bot trace", args...)
}

func traceContinueAction(bot string, world WorldSnapshot, execution ExecutorSnapshot, action Action) {
	if policyTraceLogger == nil || action.Kind == ActionNone {
		return
	}
	reason := actionKindName(action.Kind)
	if action.Kind == ActionCancel {
		reason = "continue_cancel"
	}
	args := []any{
		"bot", bot,
		"src", "continue",
		"phase", phaseName(world.Phase),
		"reason", reason,
		"map", world.Self.MapNumber,
		"x", world.Self.X,
		"y", world.Self.Y,
		"dir", world.Self.Dir,
		"alive", world.Self.Alive,
		"exec_x", execution.Position.X,
		"exec_y", execution.Position.Y,
		"exec_dir", execution.Dir,
		"exec_v", execution.PositionVersion,
		"exec_action", actionKindName(execution.CurrentAction.Kind),
		"act", actionKindName(action.Kind),
	}
	if execution.Move.Active {
		args = append(args,
			"move_next", execution.Move.PathNext,
			"move_len", len(execution.Move.Path),
			"move_end_x", execution.Move.Target.X,
			"move_end_y", execution.Move.Target.Y,
		)
	}
	if action.Target != 0 {
		args = append(args, "act_target", action.Target)
	}
	if action.Skill != 0 {
		args = append(args, "act_skill", action.Skill)
	}
	if action.Dir != 0 {
		args = append(args, "act_dir", action.Dir)
	}
	policyTraceLogger.Debug("bot trace", args...)
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
	case ActionContinueAttack:
		return "continue_attack"
	case ActionUseSkill:
		return "use_skill"
	case ActionContinueUseSkill:
		return "continue_use_skill"
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
