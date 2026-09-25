package invasion

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
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

type dragonMap struct {
	number int
	spawns []monster.EventSpawn
}

type dragonConfig struct {
	count int
	maps  []dragonMap
}

var dragonSettings dragonConfig

func (c *dragonConfig) init() {
	if err := c.load(conf.PathCommon); err != nil {
		panic(fmt.Errorf("dragon: %w", err))
	}
}

func (c *dragonConfig) load(basePath string) error {
	// Pointer attributes distinguish missing required values from valid zeroes.
	type dragonXML struct {
		XMLName  xml.Name `xml:"DragonEvent"`
		Monsters []struct {
			Index    *int `xml:"Index,attr"`
			Distance *int `xml:"Distance,attr"`
			Count    *int `xml:"Count,attr"`
			Maps     []struct {
				Number *int `xml:"Number,attr"`
				Spawns []struct {
					StartX *int `xml:"StartX,attr"`
					StartY *int `xml:"StartY,attr"`
					EndX   *int `xml:"EndX,attr"`
					EndY   *int `xml:"EndY,attr"`
				} `xml:"Spawn"`
			} `xml:"Map"`
		} `xml:"Monster"`
	}
	next := dragonConfig{}
	var dragons dragonXML
	conf.XML(basePath, "Events/IGC_DragonEvent.xml", &dragons)
	if len(dragons.Monsters) != 1 {
		return fmt.Errorf("DragonEvent requires exactly one Monster definition")
	}
	dragon := dragons.Monsters[0]
	if dragon.Index == nil || *dragon.Index != 44 || dragon.Distance == nil || dragon.Count == nil ||
		*dragon.Count <= 0 || *dragon.Count > conf.Server.GameServerInfo.MaxMonsterCount || len(dragon.Maps) == 0 {
		return fmt.Errorf("invalid red dragon class, distance, count or maps")
	}
	next.count = *dragon.Count
	mapIDs := make(map[int]bool)
	for _, dm := range dragon.Maps {
		if dm.Number == nil || mapIDs[*dm.Number] || len(dm.Spawns) == 0 {
			return fmt.Errorf("missing/duplicate dragon map or empty spawn list")
		}
		mapIDs[*dm.Number] = true
		m := dragonMap{number: *dm.Number}
		areas := make(map[monster.EventSpawn]bool)
		for _, area := range dm.Spawns {
			if area.StartX == nil || area.StartY == nil || area.EndX == nil || area.EndY == nil {
				return fmt.Errorf("dragon map %d: missing spawn coordinates", m.number)
			}
			s := monster.EventSpawn{Class: 44, MapNumber: m.number, StartX: *area.StartX, StartY: *area.StartY,
				EndX: *area.EndX, EndY: *area.EndY, Direction: -1, Distance: *dragon.Distance}
			if err := s.Validate(); err != nil {
				return err
			}
			if areas[s] {
				return fmt.Errorf("duplicate dragon spawn: %+v", s)
			}
			areas[s] = true
			m.spawns = append(m.spawns, s)
		}
		next.maps = append(next.maps, m)
	}
	*c = next
	return nil
}

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
	return &Dragon{config: dragonSettings, world: liveWorld{}, intn: rand.Intn}
}

func (d *Dragon) Running() bool { return d.state != dragonIdle }

func (d *Dragon) Start(now time.Time) error {
	if len(d.config.maps) == 0 || d.config.count <= 0 {
		return fmt.Errorf("red dragon invasion has no spawn configuration")
	}
	d.reap(now, true)
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
