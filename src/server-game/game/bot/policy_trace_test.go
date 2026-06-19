package bot

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPolicyTraceWritesDecisionJSONLines(t *testing.T) {
	file := filepath.Join(t.TempDir(), "policy.jsonl")
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		t.Fatalf("open trace file: %v", err)
	}
	logger := slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug}))

	world := WorldSnapshot{
		Phase: PhasePlaying,
		Self:  Actor{Index: 10, Class: 4, MapNumber: 0, X: 1, Y: 1, Dir: 3, Alive: true},
	}
	execution := ExecutorSnapshot{
		CurrentAction:   Action{Kind: ActionMove},
		Position:        Position{X: 1, Y: 1},
		Dir:             3,
		ReadyAt:         time.Unix(2, 3),
		PositionVersion: 7,
		Move: MoveSnapshot{
			Active:     true,
			Target:     Position{X: 4, Y: 5},
			Path:       []Position{{X: 2, Y: 1}, {X: 3, Y: 1}},
			PathNext:   1,
			NextStepAt: time.Unix(4, 5),
		},
	}
	action := Action{
		Kind:            ActionMove,
		Path:            []Position{{X: 2, Y: 1}, {X: 3, Y: 1}},
		Dir:             3,
		PositionVersion: 7,
	}

	logger.Debug("bot policy decision",
		"bot", "account1:char1",
		"phase", phaseName(world.Phase),
		"reason", "move_to_spawn_area",
		"self", traceActor(world.Self),
		"execution", traceExecution(execution),
		"action", traceAction(action),
		"extra", map[string]interface{}{"target_class": 1},
	)
	if err := f.Close(); err != nil {
		t.Fatalf("close trace file: %v", err)
	}

	entries := readPolicyTraceEntries(t, file)
	if len(entries) != 1 {
		t.Fatalf("trace entries = %d, want 1", len(entries))
	}
	entry := entries[0]
	if entry["msg"] != "bot policy decision" ||
		entry["bot"] != "account1:char1" ||
		entry["phase"] != "playing" ||
		entry["reason"] != "move_to_spawn_area" {
		t.Fatalf("trace entry = %#v", entry)
	}
	if nestedNumber(entry, "self", "x") != 1 ||
		nestedNumber(entry, "self", "y") != 1 ||
		nestedNumber(entry, "self", "class") != 4 {
		t.Fatalf("trace self = %#v", entry["self"])
	}
	if nestedString(entry, "execution", "current_action") != "move" ||
		nestedNumber(entry, "execution", "position_version") != 7 {
		t.Fatalf("trace execution = %#v", entry["execution"])
	}
	if nestedString(entry, "action", "kind") != "move" ||
		nestedNumber(entry, "action", "path_len") != 2 ||
		nestedNumber(entry, "action", "position_version") != 7 {
		t.Fatalf("trace action = %#v", entry["action"])
	}
}

func readPolicyTraceEntries(t *testing.T, file string) []map[string]interface{} {
	t.Helper()
	f, err := os.Open(file)
	if err != nil {
		t.Fatalf("open trace file: %v", err)
	}
	defer f.Close()

	var entries []map[string]interface{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var entry map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			t.Fatalf("unmarshal trace line %q: %v", scanner.Text(), err)
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan trace file: %v", err)
	}
	return entries
}

func nestedString(entry map[string]interface{}, object, field string) string {
	nested, _ := entry[object].(map[string]interface{})
	value, _ := nested[field].(string)
	return value
}

func nestedNumber(entry map[string]interface{}, object, field string) float64 {
	nested, _ := entry[object].(map[string]interface{})
	value, _ := nested[field].(float64)
	return value
}
