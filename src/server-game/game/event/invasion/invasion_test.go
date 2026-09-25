package invasion

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/event"
)

func TestScheduleConfigLoadsAndValidates(t *testing.T) {
	const times = `<InvasionManager Enable="0"><Invasion Type="0"><Start DayOfWeek="-1" Hour="12" Minute="0"/></Invasion><Invasion Type="1"><Start DayOfWeek="0" Hour="13" Minute="0"/></Invasion></InvasionManager>`
	for _, tt := range []struct {
		name, body string
		bad        bool
	}{
		{"valid", times, false},
		{"enum", strings.Replace(times, `Type="1"`, `Type="9"`, 1), true},
		{"required", strings.Replace(times, ` Minute="0"`, "", 1), true},
		{"time", strings.Replace(times, `Hour="12"`, `Hour="24"`, 1), true},
		{"duplicate", strings.Replace(times, `</Invasion>`, `<Start DayOfWeek="-1" Hour="12" Minute="0"/></Invasion>`, 1), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			base := t.TempDir()
			if err := os.Mkdir(filepath.Join(base, "Events"), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(base, "Events/IGC_InvasionManager.xml"), []byte(tt.body), 0600); err != nil {
				t.Fatal(err)
			}
			var m invasionManager
			err := m.load(base)
			if (err != nil) != tt.bad {
				t.Fatalf("load error=%v", err)
			}
			if !tt.bad && (m.Enabled || len(m.Schedules) != 2) {
				t.Fatal("lost disabled state or unsupported schedules")
			}
		})
	}
}

func TestInvasionRegistrationEnable(t *testing.T) {
	savedManager := event.EventManager
	savedPath := conf.PathCommon
	savedDragon, savedSkeleton := dragonSettings, skeletonSettings
	t.Cleanup(func() {
		event.EventManager = savedManager
		conf.PathCommon = savedPath
		dragonSettings, skeletonSettings = savedDragon, savedSkeleton
	})
	base := t.TempDir()
	if err := os.Mkdir(filepath.Join(base, "Events"), 0755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"IGC_DragonEvent.xml", "IGC_AttackEvent.xml"} {
		data, err := os.ReadFile(filepath.Join(savedPath, "Events", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(base, "Events", name), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	conf.PathCommon = base
	writeSchedule := func(enabled string) {
		t.Helper()
		data := `<InvasionManager Enable="` + enabled + `"><Invasion Type="0"><Start DayOfWeek="-1" Hour="12" Minute="0"/></Invasion><Invasion Type="1"><Start DayOfWeek="-1" Hour="13" Minute="0"/></Invasion></InvasionManager>`
		if err := os.WriteFile(filepath.Join(base, "Events/IGC_InvasionManager.xml"), []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	writeSchedule("0")
	var m invasionManager
	m.init()
	if m.Enabled || len(dragonSettings.maps) == 0 || skeletonSettings.boss.count == 0 {
		t.Fatal("disabled initialization did not load configurations")
	}
	if !InvasionManager.Enabled {
		// Successful registration proves the disabled initialization registered nothing.
		writeSchedule("1")
		m.init()
		if !m.Enabled {
			t.Fatal("enabled configuration was not loaded")
		}
	}
	for _, name := range []string{"invasion/dragon", "invasion/skeleton"} {
		if err := event.EventManager.Register(name, NewDragon(), nil); err == nil {
			t.Fatalf("missing registration: %s", name)
		}
	}
	// Disabled initialization must not attempt to register duplicates.
	writeSchedule("0")
	m.init()
	// Enabled initialization reports duplicate registration as a startup failure.
	writeSchedule("1")
	func() {
		defer func() {
			if recover() == nil {
				t.Error("duplicate registration did not panic")
			}
		}()
		m.init()
	}()
}
