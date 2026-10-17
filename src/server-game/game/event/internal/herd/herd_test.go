package herd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/object/monster"
)

const whiteXML = `<WhiteWizard Enable="1">
<Time ToOpen="3" PlayDuration="30" ToClose="25"/>
<Move RandomSize="15" MinRadius="3" MaxRadius="6"/>
<Schedule><Start Hour="0" Minute="1"/></Schedule>
<SpawnSettings><Monster Index="135" Respawn="0" AttackFirst="1" Name="White Wizard"/>
<Monster Index="136" Respawn="0" AttackFirst="0" Name="Soldier"/></SpawnSettings>
</WhiteWizard>`

const santaXML = `<XMasAttack Enable="1">
<Time ToOpen="3" PlayDuration="10" ToClose="5"/>
<Move RandomSize="5" MinRadius="3" MaxRadius="6"/>
<Reward DropRate="10" SantaBuffDuration="1800" RestBuffDuration="1800"/>
<Schedule><Start Hour="1" Minute="20"/></Schedule>
<SpawnSettings><Monster Index="476" Respawn="0" AttackFirst="1" Name="Santa"/>
<Monster Index="466" Respawn="0" AttackFirst="1" Name="Goblin"/></SpawnSettings>
</XMasAttack>`

func loadTest(t *testing.T, body, filename, root string) (Config, error) {
	t.Helper()
	base := t.TempDir()
	if err := os.WriteFile(filepath.Join(base, filename), []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	return Load(base, filename, root, root)
}

func TestLoadAndSchedules(t *testing.T) {
	white, err := loadTest(t, whiteXML, "white.xml", "WhiteWizard")
	if err != nil {
		t.Fatal(err)
	}
	if !white.Enabled || len(white.Members) != 2 || white.Schedules[0].Hour != 23 || white.Schedules[0].Minute != 58 {
		t.Fatalf("wrong white wizard schedule or members: %+v", white)
	}
	santa, err := loadTest(t, santaXML, "santa.xml", "XMasAttack")
	if err != nil {
		t.Fatal(err)
	}
	if len(santa.Members) != 2 || santa.Schedules[0].Hour != 1 || santa.Schedules[0].Minute != 17 {
		t.Fatal("wrong Santa members or schedule")
	}
	for _, body := range []string{
		strings.Replace(whiteXML, `Index="135"`, `Index="999"`, 1),
		strings.Replace(whiteXML, `ToClose="25"`, `ToClose="31"`, 1),
		strings.Replace(whiteXML, `AttackFirst="0"`, `AttackFirst="2"`, 1),
		strings.Replace(whiteXML, `Name="White Wizard"`, "", 1),
	} {
		if _, err := loadTest(t, body, "invalid.xml", "WhiteWizard"); err == nil {
			t.Fatal("invalid herd config accepted")
		}
	}
	if _, err := loadTest(t, strings.Replace(whiteXML, `Index="135"`, `Index="44"`, 1), "alternate.xml", "WhiteWizard"); err != nil {
		t.Fatal("valid configured monster class was rejected", err)
	}
	if _, err := loadTest(t, strings.Replace(santaXML, `DropRate="10"`, `DropRate="101"`, 1), "legacy.xml", "XMasAttack"); err != nil {
		t.Fatal("unused legacy DropRate rejected", err)
	}
}

func TestLifecycleAndOwnership(t *testing.T) {
	cfg, err := loadTest(t, whiteXML, "white.xml", "WhiteWizard")
	if err != nil {
		t.Fatal(err)
	}
	e := New(cfg)
	now := time.Now()
	t.Cleanup(func() { e.reap(now.Add(time.Hour), true) })
	if err := e.Start(now); err != nil {
		t.Fatal(err)
	}
	if len(e.monsters) != 0 || !e.Running() {
		t.Fatal("event spawned before open time")
	}
	e.Tick(now.Add(cfg.Open))
	if e.state != 2 || len(e.monsters) != 6 {
		t.Fatalf("spawned state=%d monsters=%d", e.state, len(e.monsters))
	}
	for _, entry := range e.monsters {
		if entry.obj.Objecter == nil || !entry.obj.NoRegen ||
			entry.obj.Objecter.(*monster.Monster).PassiveUntilHit != (entry.member.AttackFirst == 0) {
			t.Fatal("monster is not owned by White Wizard event")
		}
	}
	e.Tick(now.Add(cfg.Open + cfg.Play - cfg.Close))
	if e.state != 3 {
		t.Fatal("failed to enter destination phase")
	}
	e.Tick(now.Add(cfg.Open + cfg.Play))
	if e.Running() || len(e.monsters) != 0 {
		t.Fatal("event did not clean up at end")
	}
}

func TestConfiguredHerdCounts(t *testing.T) {
	for _, tt := range []struct {
		file, root string
		count      int
	}{
		{"Events/IGC_WhiteWizardAttack.xml", "WhiteWizard", 93},
		{"Events/IGC_XMasAttack.xml", "XMasAttack", 33},
	} {
		t.Run(tt.root, func(t *testing.T) {
			cfg, err := Load(conf.PathCommon, tt.file, tt.root, tt.root)
			if err != nil {
				t.Fatal(err)
			}
			cfg.Open = 0
			e := New(cfg)
			now := time.Now()
			t.Cleanup(func() { e.reap(now.Add(time.Hour), true) })
			if err := e.Start(now); err != nil {
				t.Fatal(err)
			}
			e.Tick(now)
			if e.state != 2 || len(e.monsters) != tt.count {
				t.Fatalf("state=%d monsters=%d, want %d", e.state, len(e.monsters), tt.count)
			}
		})
	}
}
