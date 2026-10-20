package dragon

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/object/monster"
)

type testWorld struct {
	spawned []monster.EventSpawn
	objects []*object.Object
	removed map[*object.Object]bool
	notices []string
	failAt  int
	pending *object.Object
	effects []eventStateCall
}

type eventStateCall struct{ mapNumber, state int }

func (w *testWorld) spawn(s monster.EventSpawn) (*object.Object, error) {
	w.spawned = append(w.spawned, s)
	if len(w.spawned) == w.failAt {
		return nil, errors.New("no capacity")
	}
	o := &object.Object{Index: len(w.objects), Class: s.Class, Live: true, NoRegen: true}
	w.objects = append(w.objects, o)
	return o, nil
}
func (w *testWorld) remove(o *object.Object, _ time.Time) bool {
	if o == w.pending {
		return false
	}
	w.removed[o] = true
	return true
}
func (w *testWorld) announce(s string) { w.notices = append(w.notices, s) }
func (w *testWorld) eventState(mapNumber, state int) {
	w.effects = append(w.effects, eventStateCall{mapNumber, state})
}

func testDragon(w *testWorld) *Dragon {
	w.removed = make(map[*object.Object]bool)
	return &Dragon{config: config{count: 2, maps: []spawnMap{
		{number: 0, spawns: []monster.EventSpawn{{Class: 44, MapNumber: 0}}},
		{number: 2, spawns: []monster.EventSpawn{{Class: 44, MapNumber: 2, StartX: 10}, {Class: 44, MapNumber: 2, StartX: 20}}},
	}}, world: w, intn: func(n int) int { return n - 1 }}
}

func TestDragonLifecycle(t *testing.T) {
	w := &testWorld{}
	d := testDragon(w)
	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)
	if err := d.Start(now); err != nil {
		t.Fatal(err)
	}
	d.Tick(now.Add(2 * time.Second))
	if len(w.spawned) != 0 || len(w.effects) != 0 {
		t.Fatal("spawned during announcement")
	}
	spawnAt := now.Add(4 * time.Second)
	d.Tick(spawnAt)
	if len(w.spawned) != 2 {
		t.Fatal("wrong count")
	}
	for _, s := range w.spawned {
		if s.MapNumber != 2 || s.StartX != 20 {
			t.Fatal("wrong map/area", s)
		}
	}
	if !slices.Equal(w.effects, []eventStateCall{{2, 1}}) {
		t.Fatalf("appearance effects = %v", w.effects)
	}
	// One monster dies just before expiry and still has settlement pending.
	w.objects[0].Live = false
	w.pending = w.objects[0]
	d.Tick(spawnAt.Add(duration - time.Second))
	if len(w.removed) != 0 {
		t.Fatal("removed too early")
	}
	d.Tick(spawnAt.Add(duration))
	if d.Running() || !w.removed[w.objects[1]] || w.removed[w.objects[0]] {
		t.Fatal("expiry lost pending corpse or living monster")
	}
	if !slices.Equal(w.effects, []eventStateCall{{2, 1}, {2, 0}}) {
		t.Fatalf("expiry effects = %v", w.effects)
	}
	w.pending = nil
	d.Tick(spawnAt.Add(duration + 5*time.Second))
	if len(w.removed) != 2 || len(d.monsters) != 0 {
		t.Fatal("idle corpse not drained")
	}
	if len(w.notices) != 2 {
		t.Fatal("wrong announcement count")
	}
}

func TestDragonRollbackAndNoRespawn(t *testing.T) {
	for _, fail := range []bool{false, true} {
		w := &testWorld{}
		if fail {
			w.failAt = 2
		}
		d := testDragon(w)
		now := time.Now()
		if err := d.Start(now); err != nil {
			t.Fatal(err)
		}
		d.Tick(now.Add(delay))
		if fail {
			if d.Running() || len(w.removed) != 1 || len(d.monsters) != 0 {
				t.Fatal("partial spawn not rolled back")
			}
			if len(w.effects) != 0 {
				t.Fatalf("failed spawn emitted an effect: %v", w.effects)
			}
		} else {
			for _, o := range w.objects {
				o.Live = false
			}
			d.Tick(now.Add(10 * time.Second))
			if !d.Running() || len(w.removed) != 2 {
				t.Fatal("all-killed event ended early")
			}
			if !slices.Equal(w.effects, []eventStateCall{{2, 1}}) {
				t.Fatalf("all-killed event ended the effect early: %v", w.effects)
			}
		}
		d.Tick(now.Add(100 * time.Second))
		if len(w.spawned) != 2 {
			t.Fatal("respawned after death/failure")
		}
	}
}

func TestDragonConfigLoadsAndValidates(t *testing.T) {
	const dragon = `<DragonEvent><Monster Index="44" Distance="30" Count="2"><Map Number="0"><Spawn StartX="135" StartY="61" EndX="146" EndY="70"/></Map></Monster></DragonEvent>`
	for _, tt := range []struct {
		name, dragon string
		bad          bool
	}{
		{"valid", dragon, false},
		{"map", strings.Replace(dragon, `Number="0"`, `Number="999"`, 1), true},
		{"area", strings.Replace(dragon, `EndX="146"`, `EndX="135"`, 1), true},
		{"class", strings.Replace(dragon, `Index="44"`, `Index="55"`, 1), true},
		{"count", strings.Replace(dragon, `Count="2"`, `Count="0"`, 1), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			base := t.TempDir()
			if err := os.Mkdir(filepath.Join(base, "Events"), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(base, "Events/IGC_DragonEvent.xml"), []byte(tt.dragon), 0600); err != nil {
				t.Fatal(err)
			}
			var c config
			err := c.load(base)
			if (err != nil) != tt.bad {
				t.Fatalf("load error=%v", err)
			}
			if !tt.bad && c.count != 2 {
				t.Fatal("wrong configured dragon count")
			}
		})
	}
}

func TestDragonStartupConfig(t *testing.T) {
	var c config
	if err := c.load(conf.PathCommon); err != nil {
		t.Fatal(err)
	}
	if len(c.maps) == 0 {
		t.Fatal("startup config empty")
	}
}

func TestDragonRestartPreservesSettlement(t *testing.T) {
	w := &testWorld{}
	d := testDragon(w)
	now := time.Now()
	if err := d.Start(now); err != nil {
		t.Fatal(err)
	}
	// A restart during the announcement starts a fresh countdown.
	restart := now.Add(time.Second)
	if err := d.Start(restart); err != nil {
		t.Fatal(err)
	}
	d.Tick(now.Add(delay))
	if len(w.objects) != 0 || len(w.effects) != 0 {
		t.Fatal("old announcement deadline spawned monsters")
	}
	d.Tick(restart.Add(delay))
	if len(w.objects) != 2 || !slices.Equal(w.effects, []eventStateCall{{2, 1}}) {
		t.Fatal("new announcement did not spawn")
	}
	corpse, living := w.objects[0], w.objects[1]
	corpse.Live = false
	w.pending = corpse
	restart = restart.Add(delay + time.Second)
	d.intn = func(int) int { return 0 }
	if err := d.Start(restart); err != nil {
		t.Fatal(err)
	}
	if !w.removed[living] || w.removed[corpse] || len(d.monsters) != 1 || d.state != announcing {
		t.Fatal("restart failed to preserve pending settlement or remove living monster")
	}
	if !slices.Equal(w.effects, []eventStateCall{{2, 1}, {2, 0}}) {
		t.Fatalf("restart did not clear old map effect: %v", w.effects)
	}
	d.Tick(restart.Add(delay))
	if len(w.objects) != 4 || len(d.monsters) != 3 || !d.deadline.Equal(restart.Add(delay+duration)) {
		t.Fatal("replacement wave has incorrect count or duration")
	}
	if !slices.Equal(w.effects, []eventStateCall{{2, 1}, {2, 0}, {0, 1}}) {
		t.Fatalf("replacement map effect = %v", w.effects)
	}
	w.pending = nil
	d.Tick(restart.Add(delay + time.Second))
	if !w.removed[corpse] || len(d.monsters) != 2 {
		t.Fatal("settled corpse not reclaimed")
	}
	d.Tick(restart.Add(delay + duration))
	if d.Running() || len(d.monsters) != 0 {
		t.Fatal("replacement wave did not expire")
	}
	if !slices.Equal(w.effects, []eventStateCall{{2, 1}, {2, 0}, {0, 1}, {0, 0}}) {
		t.Fatalf("replacement expiry effect = %v", w.effects)
	}
}
