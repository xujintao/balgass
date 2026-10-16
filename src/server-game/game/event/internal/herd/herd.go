// Package herd implements the shared mechanics of the White Wizard and Santa events.
package herd

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/event"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/object/monster"
)

type Member struct {
	Class, Respawn, AttackFirst, Slot int
}

type Config struct {
	Enabled                          bool
	Name                             string
	Open, Play, Close                time.Duration
	RandomSize, MinRadius, MaxRadius int
	Members                          []Member
	Schedules                        []event.Schedule
}

func Load(basePath, filename, root, name string) (Config, error) {
	// These DTOs mirror the two IGC files; repeated Monster entries are intentional.
	type xmlConfig struct {
		XMLName xml.Name
		Enable  *int `xml:"Enable,attr"`
		Time    struct {
			Open  *int `xml:"ToOpen,attr"`
			Play  *int `xml:"PlayDuration,attr"`
			Close *int `xml:"ToClose,attr"`
		} `xml:"Time"`
		Move struct {
			Random *int `xml:"RandomSize,attr"`
			Min    *int `xml:"MinRadius,attr"`
			Max    *int `xml:"MaxRadius,attr"`
		} `xml:"Move"`
		Starts []struct {
			Hour   *int `xml:"Hour,attr"`
			Minute *int `xml:"Minute,attr"`
		} `xml:"Schedule>Start"`
		Monsters []struct {
			Index       *int   `xml:"Index,attr"`
			Respawn     *int   `xml:"Respawn,attr"`
			AttackFirst *int   `xml:"AttackFirst,attr"`
			Name        string `xml:"Name,attr"`
		} `xml:"SpawnSettings>Monster"`
	}
	var raw xmlConfig
	conf.XML(basePath, filename, &raw)
	if raw.XMLName.Local != root || raw.Enable == nil || (*raw.Enable != 0 && *raw.Enable != 1) ||
		raw.Time.Open == nil || raw.Time.Play == nil || raw.Time.Close == nil ||
		*raw.Time.Open < 0 || *raw.Time.Open > 60 || *raw.Time.Play <= 0 || *raw.Time.Play > 1440 ||
		*raw.Time.Close < 0 || *raw.Time.Close > *raw.Time.Play ||
		raw.Move.Random == nil || raw.Move.Min == nil || raw.Move.Max == nil ||
		*raw.Move.Random <= 0 || *raw.Move.Min < 1 || *raw.Move.Max > 15 || *raw.Move.Min >= *raw.Move.Max ||
		len(raw.Starts) == 0 || len(raw.Monsters) == 0 || len(raw.Monsters)*len(areas) > conf.Server.MaxMonsterCount {
		return Config{}, fmt.Errorf("%s: invalid required configuration", name)
	}
	cfg := Config{Enabled: *raw.Enable == 1, Name: name,
		Open: time.Duration(*raw.Time.Open) * time.Minute, Play: time.Duration(*raw.Time.Play) * time.Minute,
		Close: time.Duration(*raw.Time.Close) * time.Minute, RandomSize: *raw.Move.Random,
		MinRadius: *raw.Move.Min, MaxRadius: *raw.Move.Max}
	for _, start := range raw.Starts {
		if start.Hour == nil || start.Minute == nil || *start.Hour < 0 || *start.Hour > 23 || *start.Minute < 0 || *start.Minute > 59 {
			return Config{}, fmt.Errorf("%s: invalid schedule", name)
		}
		minute := (*start.Hour*60 + *start.Minute - *raw.Time.Open + 1440) % 1440
		cfg.Schedules = append(cfg.Schedules, event.Schedule{DayOfWeek: -1, Hour: minute / 60, Minute: minute % 60, AppearanceRate: 100})
	}
	if err := event.ValidateSchedules(cfg.Schedules); err != nil {
		return Config{}, fmt.Errorf("%s: %w", name, err)
	}
	for _, row := range raw.Monsters {
		if row.Index == nil || row.Respawn == nil || row.AttackFirst == nil ||
			(*row.Respawn != 0 && *row.Respawn != 1) || (*row.AttackFirst != 0 && *row.AttackFirst != 1) || row.Name == "" {
			return Config{}, fmt.Errorf("%s: invalid monster definition", name)
		}
		cfg.Members = append(cfg.Members, Member{Class: *row.Index, Respawn: *row.Respawn, AttackFirst: *row.AttackFirst, Slot: len(cfg.Members)})
	}
	for _, area := range areas {
		if !maps.MapManager.HasMap(area.number) {
			return Config{}, fmt.Errorf("%s: missing map %d", name, area.number)
		}
		for _, member := range cfg.Members {
			spawn := monster.EventSpawn{Class: member.Class, MapNumber: area.number, StartX: area.minX,
				StartY: area.minY, EndX: area.maxX, EndY: area.maxY, Direction: -1, Distance: 15}
			if err := spawn.Validate(); err != nil {
				return Config{}, err
			}
		}
	}
	return cfg, nil
}

type area struct {
	number, minX, minY, maxX, maxY int
	destinations                   [][2]int
}

var areas = []area{
	{maps.Lorencia, 66, 44, 191, 197, [][2]int{{133, 79}, {87, 126}, {133, 178}, {180, 126}}},
	{maps.Devias, 111, 14, 232, 125, [][2]int{{164, 42}, {221, 85}, {164, 42}, {221, 85}}},
	{maps.Noria, 136, 53, 223, 181, [][2]int{{160, 45}, {160, 45}, {152, 117}, {209, 133}}},
}

type tracked struct {
	obj    *object.Object
	member Member
	area   int
}

type Event struct {
	Config     Config
	state      int
	start, end time.Time
	lastMove   time.Time
	centers    [3][2]int
	monsters   []tracked
}

func New(cfg Config) *Event {
	return &Event{Config: cfg}
}

func (e *Event) Running() bool {
	return e.state != 0
}

func (e *Event) Start(now time.Time) error {
	e.state = 0
	e.reap(now, true)
	e.state = 1
	e.start = now.Add(e.Config.Open)
	e.end = e.start.Add(e.Config.Play)
	object.ObjectManager.BroadcastSystemMsg(e.Config.Name + "即将开始！")
	return nil
}

func (e *Event) Tick(now time.Time) {
	e.reap(now, false)
	if e.state == 1 && !now.Before(e.start) {
		if err := e.spawn(now); err != nil {
			e.state = 0
			slog.Error("herd event spawn failed", "event", e.Config.Name, "err", err)
			return
		}
		e.state = 2
		e.lastMove = now
		object.ObjectManager.BroadcastSystemMsg(e.Config.Name + "开始！")
	}
	if e.state != 2 && e.state != 3 {
		return
	}
	if !now.Before(e.end) {
		e.reap(now, true)
		e.state = 0
		object.ObjectManager.BroadcastSystemMsg(e.Config.Name + "结束。")
		return
	}
	if e.state == 2 && e.Config.Close > 0 && !now.Before(e.end.Add(-e.Config.Close)) {
		e.moveToDestinations()
		e.state = 3
	}
	if e.state == 2 && !now.Before(e.lastMove.Add(time.Second)) {
		e.lastMove = now
		if rand.Intn(10) == 0 {
			e.wander()
		}
	}
	if e.state == 2 || e.state == 3 {
		e.respawn(now)
	}
}

func (e *Event) spawn(now time.Time) error {
	var created []tracked
	for ai, area := range areas {
		center, ok := randomCenter(area, max(3, e.Config.MaxRadius), len(e.Config.Members))
		if !ok {
			for _, entry := range created {
				object.ObjectManager.DeleteEventMonster(entry.obj, now)
			}
			return fmt.Errorf("%s: no center on map %d", e.Config.Name, area.number)
		}
		e.centers[ai] = center
		for _, member := range e.Config.Members {
			obj, err := e.spawnMember(ai, member)
			if err != nil {
				for _, entry := range created {
					object.ObjectManager.DeleteEventMonster(entry.obj, now)
				}
				return err
			}
			created = append(created, tracked{obj: obj, member: member, area: ai})
		}
	}
	e.monsters = append(e.monsters, created...)
	return nil
}

func (e *Event) spawnMember(ai int, member Member) (*object.Object, error) {
	center := e.centers[ai]
	radius := max(3, e.Config.MaxRadius)
	spawn := monster.EventSpawn{Class: member.Class, MapNumber: areas[ai].number,
		StartX: max(0, center[0]-radius), StartY: max(0, center[1]-radius),
		EndX: min(256, center[0]+radius+1), EndY: min(256, center[1]+radius+1),
		Direction: -1, Distance: radius,
		EventHerd: true, PassiveUntilHit: member.AttackFirst == 0}
	obj, err := monster.SpawnEventMonster(spawn)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (e *Event) reap(now time.Time, living bool) {
	kept := e.monsters[:0]
	for _, entry := range e.monsters {
		if (living || !entry.obj.Live) && object.ObjectManager.DeleteEventMonster(entry.obj, now) {
			continue
		}
		kept = append(kept, entry)
	}
	clear(e.monsters[len(kept):])
	e.monsters = kept
}

func (e *Event) respawn(now time.Time) {
	for ai := range areas {
		for _, member := range e.Config.Members {
			if member.Respawn == 0 {
				continue
			}
			count := 0
			for _, entry := range e.monsters {
				if entry.area == ai && entry.member == member {
					count++
				}
			}
			if count != 0 {
				continue
			}
			obj, err := e.spawnMember(ai, member)
			if err != nil {
				slog.Error("herd respawn failed", "event", e.Config.Name, "err", err)
				continue
			}
			e.monsters = append(e.monsters, tracked{obj: obj, member: member, area: ai})
		}
	}
}

func (e *Event) wander() {
	for ai, area := range areas {
		center := e.centers[ai]
		minX, minY := max(area.minX, center[0]-e.Config.RandomSize), max(area.minY, center[1]-e.Config.RandomSize)
		maxX, maxY := min(area.maxX, center[0]+e.Config.RandomSize+1), min(area.maxY, center[1]+e.Config.RandomSize+1)
		if next, ok := randomTile(area.number, minX, minY, maxX, maxY); ok {
			e.move(ai, next)
		}
	}
}

func (e *Event) moveToDestinations() {
	for ai, area := range areas {
		for _, destination := range rand.Perm(len(area.destinations)) {
			xy := area.destinations[destination]
			if maps.MapManager.GetMapAttr(area.number, xy[0], xy[1])&15 == 0 {
				e.move(ai, xy)
				break
			}
		}
	}
}

func (e *Event) move(ai int, center [2]int) {
	e.centers[ai] = center
	radius := e.Config.MinRadius + rand.Intn(e.Config.MaxRadius-e.Config.MinRadius)
	for _, entry := range e.monsters {
		if entry.area == ai && entry.obj.Live {
			entry.obj.Objecter.(*monster.Monster).SetEventHerdCenter(center[0], center[1], radius)
		}
	}
}

func randomTile(mapNumber, minX, minY, maxX, maxY int) ([2]int, bool) {
	var tile [2]int
	count := 0
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			if maps.MapManager.GetMapAttr(mapNumber, x, y)&15 != 0 {
				continue
			}
			count++
			if rand.Intn(count) == 0 {
				tile = [2]int{x, y}
			}
		}
	}
	return tile, count > 0
}

func randomCenter(area area, radius, needed int) ([2]int, bool) {
	var center [2]int
	count := 0
	for y := area.minY; y < area.maxY; y++ {
		for x := area.minX; x < area.maxX; x++ {
			if maps.MapManager.GetMapAttr(area.number, x, y)&15 != 0 {
				continue
			}
			free := 0
			for sy := max(0, y-radius); sy <= min(255, y+radius); sy++ {
				for sx := max(0, x-radius); sx <= min(255, x+radius); sx++ {
					if maps.MapManager.GetMapAttr(area.number, sx, sy)&15 == 0 {
						free++
					}
				}
			}
			if free < needed {
				continue
			}
			count++
			if rand.Intn(count) == 0 {
				center = [2]int{x, y}
			}
		}
	}
	return center, count > 0
}
