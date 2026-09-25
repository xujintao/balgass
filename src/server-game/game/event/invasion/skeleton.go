package invasion

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/object/monster"
)

type skeletonSpawn struct {
	count int
	maps  []int
	spawn monster.EventSpawn
}
type skeletonConfig struct{ boss, guard skeletonSpawn }

var skeletonSettings skeletonConfig

func (c *skeletonConfig) init() {
	if err := c.load(conf.PathCommon); err != nil {
		panic(fmt.Errorf("skeleton: %w", err))
	}
}
func (c *skeletonConfig) load(basePath string) error {
	type config struct {
		XMLName xml.Name `xml:"AttackEvent"`
		Spawns  []struct {
			Index    *int `xml:"Index,attr"`
			Map1     *int `xml:"MapNumber1,attr"`
			Map2     *int `xml:"MapNumber2,attr"`
			Map3     *int `xml:"MapNumber3,attr"`
			Map4     *int `xml:"MapNumber4,attr"`
			Distance *int `xml:"Distance,attr"`
			X        *int `xml:"StartX,attr"`
			Y        *int `xml:"StartY,attr"`
			EndX     *int `xml:"EndX,attr"`
			EndY     *int `xml:"EndY,attr"`
			Dir      *int `xml:"Dir,attr"`
			Count    *int `xml:"Count,attr"`
		} `xml:"Spawn"`
	}
	var cfg config
	conf.XML(basePath, "Events/IGC_AttackEvent.xml", &cfg)
	var next skeletonConfig
	seen := make(map[int]bool)
	total := 0
	for _, row := range cfg.Spawns {
		if row.Index == nil || *row.Index < 53 || *row.Index > 56 || seen[*row.Index] {
			return fmt.Errorf("invalid or duplicate AttackEvent monster")
		}
		seen[*row.Index] = true
		// Golden Titan/Soldier belong to a separate, not yet registered activity.
		if *row.Index < 55 {
			continue
		}
		if row.Distance == nil || row.X == nil || row.Y == nil || row.EndX == nil || row.EndY == nil || row.Dir == nil || row.Count == nil || *row.Count <= 0 || *row.Count > conf.Server.GameServerInfo.MaxMonsterCount {
			return fmt.Errorf("invalid skeleton spawn fields")
		}
		spawn := skeletonSpawn{count: *row.Count, spawn: monster.EventSpawn{Class: *row.Index, StartX: *row.X, StartY: *row.Y, EndX: *row.EndX, EndY: *row.EndY, Direction: *row.Dir, Distance: *row.Distance}}
		maps := make(map[int]bool)
		for _, number := range []*int{row.Map1, row.Map2, row.Map3, row.Map4} {
			if number == nil {
				return fmt.Errorf("missing skeleton map")
			}
			if *number == -1 {
				continue
			}
			if maps[*number] {
				return fmt.Errorf("duplicate skeleton map %d", *number)
			}
			maps[*number] = true
			spawn.spawn.MapNumber = *number
			if err := spawn.spawn.Validate(); err != nil {
				return err
			}
			spawn.maps = append(spawn.maps, *number)
		}
		if len(spawn.maps) == 0 {
			return fmt.Errorf("skeleton spawn has no map")
		}
		total += spawn.count
		if *row.Index == 55 {
			next.boss = spawn
		} else {
			next.guard = spawn
		}
	}
	if !seen[55] || !seen[56] || total > conf.Server.GameServerInfo.MaxMonsterCount {
		return fmt.Errorf("skeleton requires boss and guard within monster capacity")
	}
	*c = next
	return nil
}

type Skeleton struct {
	config   skeletonConfig
	monsters []*object.Object
}

func NewSkeleton() *Skeleton { return &Skeleton{config: skeletonSettings} }
func (s *Skeleton) Running() bool {
	for _, obj := range s.monsters {
		if object.ObjectManager.GetObject(obj.Index) == obj && obj.Live {
			return true
		}
	}
	return false
}
func (s *Skeleton) Start(now time.Time) error {
	s.reap(now, true)
	var created []*object.Object
	rollback := func(err error) error {
		for _, obj := range created {
			object.ObjectManager.DeleteEventMonster(obj, now)
		}
		return err
	}
	if len(s.config.boss.maps) == 0 || s.config.boss.count <= 0 {
		return fmt.Errorf("empty skeleton configuration")
	}
	var boss *object.Object
	for i := 0; i < s.config.boss.count; i++ {
		spawn := s.config.boss.spawn
		spawn.MapNumber = s.config.boss.maps[rand.Intn(len(s.config.boss.maps))]
		obj, err := monster.SpawnEventMonster(spawn)
		if err != nil {
			return rollback(err)
		}
		created = append(created, obj)
		boss = obj
	}
	// Match AttackEvent: the configured guards surround the last spawned king.
	for i := 0; i < s.config.guard.count; i++ {
		spawn := s.config.guard.spawn
		spawn.MapNumber = boss.MapNumber
		spawn.StartX, spawn.StartY = max(0, boss.X-4), max(0, boss.Y-4)
		spawn.EndX, spawn.EndY = min(256, boss.X+4), min(256, boss.Y+4)
		obj, err := monster.SpawnEventMonster(spawn)
		if err != nil {
			return rollback(err)
		}
		created = append(created, obj)
	}
	s.monsters = append(s.monsters, created...)
	object.ObjectManager.BroadcastSystemMsg("骷髅王入侵开始！")
	return nil
}
func (s *Skeleton) Tick(now time.Time) { s.reap(now, false) }
func (s *Skeleton) reap(now time.Time, living bool) {
	kept := s.monsters[:0]
	for _, obj := range s.monsters {
		if (living || !obj.Live) && object.ObjectManager.DeleteEventMonster(obj, now) {
			continue
		}
		kept = append(kept, obj)
	}
	clear(s.monsters[len(kept):])
	s.monsters = kept
}
