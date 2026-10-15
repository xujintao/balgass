package golden

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/object"
)

const goldenXML = `<GoldenInvasion><Spot>
<Spawn Index="53" MapNumber1="0" MapNumber2="-1" MapNumber3="-1" MapNumber4="-1" Distance="30" StartX="135" StartY="61" EndX="146" EndY="70" Dir="-1" Count="1"/>
<Spawn Index="54" MapNumber1="0" MapNumber2="-1" MapNumber3="-1" MapNumber4="-1" Distance="30" StartX="135" StartY="61" EndX="146" EndY="70" Dir="-1" Count="2"/>
<Spawn Index="78" MapNumber1="0" MapNumber2="-1" MapNumber3="-1" MapNumber4="-1" Distance="30" StartX="135" StartY="61" EndX="146" EndY="70" Dir="-1" Count="1"/>
<Spawn Index="82" MapNumber1="0" MapNumber2="-1" MapNumber3="-1" MapNumber4="-1" Distance="30" StartX="135" StartY="61" EndX="146" EndY="70" Dir="-1" Count="1"/>
<Spawn Index="83" MapNumber1="0" MapNumber2="-1" MapNumber3="-1" MapNumber4="-1" Distance="30" StartX="135" StartY="61" EndX="146" EndY="70" Dir="-1" Count="1"/>
</Spot></GoldenInvasion>`

func loadGolden(t *testing.T, body string) (config, error) {
	t.Helper()
	base := t.TempDir()
	if err := os.Mkdir(filepath.Join(base, "Events"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "Events/IGC_GoldenInvasion.xml"), []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	var cfg config
	err := cfg.load(base)
	return cfg, err
}

func TestConfigValidation(t *testing.T) {
	for _, tt := range []struct {
		name, body string
		bad        bool
	}{
		{"valid", goldenXML, false},
		{"unknown class", strings.Replace(goldenXML, `Index="78"`, `Index="999"`, 1), true},
		{"missing map", strings.Replace(goldenXML, `MapNumber1="0"`, "", 1), true},
		{"duplicate map", strings.Replace(goldenXML, `MapNumber2="-1"`, `MapNumber2="0"`, 1), true},
		{"invalid area", strings.Replace(goldenXML, `EndX="146"`, `EndX="135"`, 1), true},
		{"duplicate class", strings.Replace(goldenXML, `Index="78"`, `Index="53"`, 1), true},
		{"missing boss", strings.Replace(goldenXML, `Index="53"`, `Index="79"`, 1), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := loadGolden(t, tt.body)
			if (err != nil) != tt.bad {
				t.Fatalf("load error = %v", err)
			}
			if !tt.bad && (len(cfg.order) != 5 || cfg.spawns[54].count != 2) {
				t.Fatal("lost golden spawn configuration")
			}
		})
	}
}

func TestSpawnAndRefresh(t *testing.T) {
	cfg, err := loadGolden(t, goldenXML)
	if err != nil {
		t.Fatal(err)
	}
	g := &Golden{config: cfg}
	now := time.Now()
	t.Cleanup(func() { g.reap(now.Add(time.Hour), true) })
	if err := g.Start(now); err != nil {
		t.Fatal(err)
	}
	if !g.Running() || len(g.monsters) != 6 {
		t.Fatalf("running=%v, count=%d", g.Running(), len(g.monsters))
	}
	bosses := map[int]*object.Object{}
	count := map[int]int{}
	for _, obj := range g.monsters {
		count[obj.Class]++
		if obj.Class == 53 || obj.Class == 82 {
			bosses[obj.Class] = obj
		}
		if !obj.NoRegen || !obj.Live {
			t.Fatal("golden monster must be live and event-owned")
		}
	}
	if count[53] != 1 || count[54] != 2 || count[78] != 1 || count[82] != 1 || count[83] != 1 {
		t.Fatalf("unexpected classes: %v", count)
	}
	for _, obj := range g.monsters {
		bossClass, radius := 0, 0
		switch obj.Class {
		case 54:
			bossClass, radius = 53, 4
		case 83:
			bossClass, radius = 82, 10
		}
		if bossClass != 0 {
			boss := bosses[bossClass]
			if obj.MapNumber != boss.MapNumber || abs(obj.X-boss.X) > radius || abs(obj.Y-boss.Y) > radius {
				t.Fatalf("guard %d at %d,%d is not near boss %d at %d,%d", obj.Class, obj.X, obj.Y, bossClass, boss.X, boss.Y)
			}
		}
	}
	old := append([]*object.Object(nil), g.monsters...)
	if err := g.Start(now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if len(g.monsters) != 6 {
		t.Fatalf("refresh count = %d", len(g.monsters))
	}
	for _, obj := range old {
		if object.ObjectManager.GetObject(obj.Index) == obj {
			t.Fatal("refresh left an old golden monster")
		}
	}
	dead := g.monsters[0]
	dead.Live = false
	g.Tick(now.Add(2 * time.Minute))
	if object.ObjectManager.GetObject(dead.Index) == dead || len(g.monsters) != 5 {
		t.Fatal("dead golden monster was not reaped")
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
