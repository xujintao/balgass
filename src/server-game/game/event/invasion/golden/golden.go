package golden

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/object/monster"
)

func init() {
	settings.init()
}

var settings config

type spawnConfig struct {
	count int
	maps  []int
	spawn monster.EventSpawn
}

type config struct {
	spawns map[int]spawnConfig
	order  []int
}

var classes = map[int]bool{
	53: true, 54: true, 78: true, 79: true, 80: true, 81: true, 82: true, 83: true,
	493: true, 494: true, 495: true, 496: true, 497: true, 498: true, 499: true,
	500: true, 501: true, 502: true,
}

var guards = map[int]int{54: 53, 81: 80, 83: 82}

func (c *config) init() {
	if err := c.load(conf.PathCommon); err != nil {
		panic(fmt.Errorf("golden: %w", err))
	}
}

func (c *config) load(basePath string) error {
	type xmlConfig struct {
		XMLName xml.Name `xml:"GoldenInvasion"`
		Spawns  []struct {
			Index    *int `xml:"Index,attr"`
			Map1     *int `xml:"MapNumber1,attr"`
			Map2     *int `xml:"MapNumber2,attr"`
			Map3     *int `xml:"MapNumber3,attr"`
			Map4     *int `xml:"MapNumber4,attr"`
			Distance *int `xml:"Distance,attr"`
			StartX   *int `xml:"StartX,attr"`
			StartY   *int `xml:"StartY,attr"`
			EndX     *int `xml:"EndX,attr"`
			EndY     *int `xml:"EndY,attr"`
			Dir      *int `xml:"Dir,attr"`
			Count    *int `xml:"Count,attr"`
		} `xml:"Spot>Spawn"`
	}
	var cfg xmlConfig
	conf.XML(basePath, "Events/IGC_GoldenInvasion.xml", &cfg)
	if len(cfg.Spawns) == 0 {
		return fmt.Errorf("GoldenInvasion requires spawn definitions")
	}
	next := config{spawns: make(map[int]spawnConfig, len(cfg.Spawns))}
	total := 0
	for _, row := range cfg.Spawns {
		if row.Index == nil || !classes[*row.Index] {
			return fmt.Errorf("invalid golden monster class")
		}
		class := *row.Index
		if _, exists := next.spawns[class]; exists {
			return fmt.Errorf("duplicate golden monster class %d", class)
		}
		if row.Distance == nil || row.StartX == nil || row.StartY == nil || row.EndX == nil || row.EndY == nil || row.Dir == nil || row.Count == nil ||
			*row.Count <= 0 || *row.Count > conf.Server.MaxMonsterCount {
			return fmt.Errorf("golden monster %d: missing or invalid spawn fields", class)
		}
		spawn := spawnConfig{count: *row.Count, spawn: monster.EventSpawn{
			Class: class, StartX: *row.StartX, StartY: *row.StartY, EndX: *row.EndX, EndY: *row.EndY,
			Direction: *row.Dir, Distance: *row.Distance,
		}}
		seenMaps := make(map[int]bool)
		for _, number := range []*int{row.Map1, row.Map2, row.Map3, row.Map4} {
			if number == nil {
				return fmt.Errorf("golden monster %d: missing map", class)
			}
			if *number == -1 {
				continue
			}
			if seenMaps[*number] {
				return fmt.Errorf("golden monster %d: duplicate map %d", class, *number)
			}
			seenMaps[*number] = true
			spawn.spawn.MapNumber = *number
			if err := spawn.spawn.Validate(); err != nil {
				return fmt.Errorf("golden monster %d: %w", class, err)
			}
			spawn.maps = append(spawn.maps, *number)
		}
		if len(spawn.maps) == 0 {
			return fmt.Errorf("golden monster %d: no enabled map", class)
		}
		total += spawn.count
		if total > conf.Server.MaxMonsterCount {
			return fmt.Errorf("golden invasion exceeds monster capacity")
		}
		next.spawns[class] = spawn
		next.order = append(next.order, class)
	}
	for guard, boss := range guards {
		if _, exists := next.spawns[guard]; exists {
			if _, exists := next.spawns[boss]; !exists {
				return fmt.Errorf("golden guard %d requires boss %d", guard, boss)
			}
		}
	}
	*c = next
	return nil
}

type Golden struct {
	config   config
	monsters []*object.Object
}

func New() *Golden {
	return &Golden{config: settings}
}

func (g *Golden) Running() bool {
	for _, obj := range g.monsters {
		if object.ObjectManager.GetObject(obj.Index) == obj && obj.Live {
			return true
		}
	}
	return false
}

func (g *Golden) Start(now time.Time) error {
	g.reap(now, true)
	if len(g.config.order) == 0 {
		return fmt.Errorf("empty golden invasion configuration")
	}
	var created []*object.Object
	rollback := func(err error) error {
		for _, obj := range created {
			object.ObjectManager.DeleteEventMonster(obj, now)
		}
		return err
	}
	bosses := make(map[int]*object.Object)
	for _, class := range g.config.order {
		if guards[class] != 0 {
			continue
		}
		entry := g.config.spawns[class]
		for i := 0; i < entry.count; i++ {
			spawn := entry.spawn
			spawn.MapNumber = entry.maps[rand.Intn(len(entry.maps))]
			obj, err := monster.SpawnEventMonster(spawn)
			if err != nil {
				return rollback(err)
			}
			created = append(created, obj)
			bosses[class] = obj
		}
	}
	for _, class := range g.config.order {
		bossClass := guards[class]
		if bossClass == 0 {
			continue
		}
		entry := g.config.spawns[class]
		boss := bosses[bossClass]
		radius := 4
		if class == 83 {
			radius = 10
		}
		for i := 0; i < entry.count; i++ {
			spawn := entry.spawn
			spawn.MapNumber = boss.MapNumber
			spawn.StartX, spawn.StartY = max(0, boss.X-radius), max(0, boss.Y-radius)
			spawn.EndX, spawn.EndY = min(256, boss.X+radius), min(256, boss.Y+radius)
			obj, err := monster.SpawnEventMonster(spawn)
			if err != nil {
				return rollback(err)
			}
			created = append(created, obj)
		}
	}
	g.monsters = append(g.monsters, created...)
	object.ObjectManager.BroadcastSystemMsg("黄金怪物入侵开始！")
	return nil
}

func (g *Golden) Tick(now time.Time) { g.reap(now, false) }

func (g *Golden) reap(now time.Time, living bool) {
	kept := g.monsters[:0]
	for _, obj := range g.monsters {
		if (living || !obj.Live) && object.ObjectManager.DeleteEventMonster(obj, now) {
			continue
		}
		kept = append(kept, obj)
	}
	clear(g.monsters[len(kept):])
	g.monsters = kept
}
