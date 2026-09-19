package invasion

import (
	"errors"
	"testing"
	"time"

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
}

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

func testDragon(w *testWorld) *Dragon {
	w.removed = make(map[*object.Object]bool)
	return &Dragon{config: dragonConfig{count: 2, maps: []dragonMap{
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
	if err := d.Start(now); err == nil {
		t.Fatal("allowed overlapping start")
	}
	d.Tick(now.Add(2 * time.Second))
	if len(w.spawned) != 0 {
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
	// One monster dies just before expiry and still has settlement pending.
	w.objects[0].Live = false
	w.pending = w.objects[0]
	d.Tick(spawnAt.Add(dragonDuration - time.Second))
	if len(w.removed) != 0 {
		t.Fatal("removed too early")
	}
	d.Tick(spawnAt.Add(dragonDuration))
	if d.Running() || !w.removed[w.objects[1]] || w.removed[w.objects[0]] {
		t.Fatal("expiry lost pending corpse or living monster")
	}
	w.pending = nil
	d.Tick(spawnAt.Add(dragonDuration + 5*time.Second))
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
		d.Tick(now.Add(dragonDelay))
		if fail {
			if d.Running() || len(w.removed) != 1 || len(d.monsters) != 0 {
				t.Fatal("partial spawn not rolled back")
			}
		} else {
			for _, o := range w.objects {
				o.Live = false
			}
			d.Tick(now.Add(10 * time.Second))
			if !d.Running() || len(w.removed) != 2 {
				t.Fatal("all-killed event ended early")
			}
		}
		d.Tick(now.Add(100 * time.Second))
		if len(w.spawned) != 2 {
			t.Fatal("respawned after death/failure")
		}
	}
}
