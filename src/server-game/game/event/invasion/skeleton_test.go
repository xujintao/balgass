package invasion

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/object/monster"
)

func TestSkeletonConfiguration(t *testing.T) {
	valid := `<AttackEvent><Spawn Index="55" MapNumber1="0" MapNumber2="-1" MapNumber3="-1" MapNumber4="-1" Distance="30" StartX="135" StartY="61" EndX="146" EndY="70" Dir="-1" Count="1"/><Spawn Index="56" MapNumber1="0" MapNumber2="-1" MapNumber3="-1" MapNumber4="-1" Distance="30" StartX="135" StartY="61" EndX="146" EndY="70" Dir="-1" Count="2"/></AttackEvent>`
	for _, tt := range []struct {
		name, body string
		bad        bool
	}{
		{"valid", valid, false},
		{"missing guard", strings.Replace(valid, `Index="56"`, `Index="54"`, 1), true},
		{"map", strings.Replace(valid, `MapNumber1="0"`, `MapNumber1="999"`, 1), true},
		{"duplicate", strings.Replace(valid, `Index="56"`, `Index="55"`, 1), true},
		{"count", strings.Replace(valid, `Count="1"`, `Count="0"`, 1), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			base := t.TempDir()
			if err := os.Mkdir(filepath.Join(base, "Events"), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(base, "Events/IGC_AttackEvent.xml"), []byte(tt.body), 0600); err != nil {
				t.Fatal(err)
			}
			var c skeletonConfig
			if err := c.load(base); (err != nil) != tt.bad {
				t.Fatalf("load: %v", err)
			}
		})
	}
}
func TestSkeletonSpawnRefreshAndRollback(t *testing.T) {
	now := time.Now()
	s := NewSkeleton()
	s.config.boss = skeletonSpawn{count: 1, maps: []int{0}, spawn: monster.EventSpawn{Class: 55, MapNumber: 0, StartX: 135, StartY: 61, EndX: 146, EndY: 70, Direction: -1, Distance: 30}}
	s.config.guard = skeletonSpawn{count: 2, spawn: monster.EventSpawn{Class: 56, Direction: -1, Distance: 30}}
	t.Cleanup(func() {
		for _, obj := range s.monsters {
			obj.Live = true
			object.ObjectManager.DeleteEventMonster(obj, time.Now())
		}
	})
	if err := s.Start(now); err != nil {
		t.Fatal(err)
	}
	if len(s.monsters) != 3 || !s.Running() {
		t.Fatal("missing king or guards")
	}
	boss := s.monsters[0]
	for _, guard := range s.monsters[1:] {
		if guard.Class != 56 || guard.MapNumber != boss.MapNumber || !guard.NoRegen || guard.X < boss.X-4 || guard.X >= boss.X+4 || guard.Y < boss.Y-4 || guard.Y >= boss.Y+4 {
			t.Fatal("invalid guard spawn")
		}
	}
	// Pending deaths from the previous wave survive replacement for settlement.
	boss.Live = false
	boss.AddDelayMsg(1, 0, 60000, 0)
	oldGuard := s.monsters[1]
	if err := s.Start(now.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if object.ObjectManager.GetObject(boss.Index) != boss || object.ObjectManager.GetObject(oldGuard.Index) == oldGuard || len(s.monsters) != 4 {
		t.Fatal("refresh lost corpse or kept living guards")
	}
	// An invalid guard fails after the king was created; that king must roll back.
	s.config.guard.spawn.Class = 99999
	before := map[*object.Object]bool{}
	for i := 0; i < conf.Server.MaxMonsterCount; i++ {
		if obj := object.ObjectManager.GetObject(i); obj != nil {
			before[obj] = true
		}
	}
	if err := s.Start(now.Add(2 * time.Hour)); err == nil {
		t.Fatal("expected spawn failure")
	}
	for i := 0; i < conf.Server.MaxMonsterCount; i++ {
		if obj := object.ObjectManager.GetObject(i); obj != nil && !before[obj] {
			t.Fatal("partial spawn leaked")
		}
	}
}
