package invasion

import (
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/object/monster"
)

const (
	dragonIdle = iota
	dragonAnnouncing
	dragonActive
	dragonDelay    = 3 * time.Second
	dragonDuration = 5 * time.Minute
)

type dragonWorld interface {
	spawn(monster.EventSpawn) (*object.Object, error)
	remove(*object.Object, time.Time) bool
	announce(string)
}

type liveWorld struct{}

func (liveWorld) spawn(s monster.EventSpawn) (*object.Object, error) {
	return monster.SpawnEventMonster(s)
}
func (liveWorld) remove(obj *object.Object, now time.Time) bool {
	return object.ObjectManager.DeleteEventMonster(obj, now)
}
func (liveWorld) announce(msg string) { object.ObjectManager.BroadcastSystemMsg(msg) }

type Dragon struct {
	config   dragonConfig
	world    dragonWorld
	intn     func(int) int
	state    int
	deadline time.Time
	mapIndex int
	monsters []*object.Object
}

func NewDragon() *Dragon {
	return &Dragon{config: config.dragon, world: liveWorld{}, intn: rand.Intn}
}

func (d *Dragon) Running() bool { return d.state != dragonIdle }

func (d *Dragon) Start(now time.Time) error {
	if d.Running() {
		return fmt.Errorf("red dragon invasion is already running")
	}
	if len(d.config.maps) == 0 || d.config.count <= 0 {
		return fmt.Errorf("red dragon invasion has no spawn configuration")
	}
	d.mapIndex = d.intn(len(d.config.maps))
	d.state = dragonAnnouncing
	d.deadline = now.Add(dragonDelay)
	d.world.announce(fmt.Sprintf("红龙即将入侵%s！", maps.MapManager.GetMapName(d.config.maps[d.mapIndex].number)))
	slog.Info("red dragon invasion announced", "map", d.config.maps[d.mapIndex].number)
	return nil
}

func (d *Dragon) Tick(now time.Time) {
	// Always drain corpses, including those killed just before the event ended.
	d.reap(now, false)
	if !d.Running() || now.Before(d.deadline) {
		return
	}
	if d.state == dragonAnnouncing {
		m := d.config.maps[d.mapIndex]
		var created []*object.Object
		for i := 0; i < d.config.count; i++ {
			obj, err := d.world.spawn(m.spawns[d.intn(len(m.spawns))])
			if err != nil {
				for _, previous := range created {
					d.world.remove(previous, now)
				}
				d.state = dragonIdle
				d.world.announce("红龙入侵取消。")
				slog.Error("red dragon invasion spawn failed", "map", m.number, "err", err)
				return
			}
			created = append(created, obj)
		}
		d.monsters = append(d.monsters, created...)
		d.state = dragonActive
		d.deadline = now.Add(dragonDuration)
		slog.Info("red dragon invasion started", "map", m.number, "count", len(created))
		return
	}
	d.reap(now, true)
	d.state = dragonIdle
	d.world.announce("红龙入侵结束。")
	slog.Info("red dragon invasion ended", "map", d.config.maps[d.mapIndex].number)
}

func (d *Dragon) reap(now time.Time, includeLiving bool) {
	kept := d.monsters[:0]
	for _, obj := range d.monsters {
		if (includeLiving || !obj.Live) && d.world.remove(obj, now) {
			continue
		}
		kept = append(kept, obj)
	}
	clear(d.monsters[len(kept):])
	d.monsters = kept
}
