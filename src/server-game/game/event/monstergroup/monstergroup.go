package monstergroup

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/drop"
	"github.com/xujintao/balgass/src/server-game/game/event"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/object/monster"
)

type StartTime struct{ Hour, Minute, AppearanceRate int }
type member struct {
	class, count                                                        int
	override                                                            bool
	eventID, hp, damageMin, damageMax, defense, attackRate, defenseRate int
}
type Group struct {
	Index                     int
	Starts                    []StartTime
	mapNumber, boss, duration int // Duration is reserved by IGC and currently unused.
	notice                    string
	areas                     []monster.EventSpawn
	members                   []member
	monsters                  []*object.Object
}
type monsterGroupManager struct{ Groups []*Group }

var MonsterGroupManager monsterGroupManager

func init() { MonsterGroupManager.init() }
func (m *monsterGroupManager) init() {
	if err := m.load(conf.PathCommon); err != nil {
		panic(fmt.Errorf("monstergroup: %w", err))
	}
	for _, group := range m.Groups {
		var schedules []event.Schedule
		for _, start := range group.Starts {
			schedules = append(schedules, event.Schedule{DayOfWeek: -1, Hour: start.Hour, Minute: start.Minute, AppearanceRate: start.AppearanceRate})
		}
		if err := event.EventManager.Register(fmt.Sprintf("monstergroup/%d", group.Index), group, schedules); err != nil {
			panic(err)
		}
	}
}
func (m *monsterGroupManager) load(basePath string) error {
	type config struct {
		XMLName xml.Name `xml:"MonsterGroupRegenSystem"`
		Notice  *int     `xml:"SpawnNotice,attr"`
		Groups  []struct {
			Index    *int   `xml:"Index,attr"`
			Map      *int   `xml:"MapNumber,attr"`
			Boss     *int   `xml:"BossMonsterIndex,attr"`
			Duration *int   `xml:"Duration,attr"`
			Notice   string `xml:"SpawnNotice,attr"`
		} `xml:"GroupSettings>Group"`
		Spots []struct {
			Index *int `xml:"Index,attr"`
			Areas []struct {
				X    *int `xml:"StartX,attr"`
				Y    *int `xml:"StartY,attr"`
				EndX *int `xml:"EndX,attr"`
				EndY *int `xml:"EndY,attr"`
			} `xml:"Area"`
		} `xml:"SpotSettings>Group"`
		Monsters []struct {
			Index   *int `xml:"Index,attr"`
			Members []struct {
				Index       *int `xml:"Index,attr"`
				Count       *int `xml:"Count,attr"`
				Override    *int `xml:"OverrideDefaultSettings,attr"`
				EventID     *int `xml:"EventID,attr"`
				HP          *int `xml:"HP,attr"`
				Min         *int `xml:"DamageMin,attr"`
				Max         *int `xml:"DamageMax,attr"`
				Defense     *int `xml:"Defense,attr"`
				AttackRate  *int `xml:"AttackRate,attr"`
				DefenseRate *int `xml:"DefenseRate,attr"`
			} `xml:"Monster"`
		} `xml:"MonsterSettings>Group"`
		Spawns []struct {
			Index  *int `xml:"Index,attr"`
			Starts []struct {
				Hour   *int `xml:"StartHour,attr"`
				Minute *int `xml:"StarMinute,attr"` // Original IGC attribute spelling.
				Rate   *int `xml:"AppearanceRate,attr"`
			} `xml:"Spawn"`
		} `xml:"SpawnSettings>Group"`
	}
	var cfg config
	conf.XML(basePath, "IGC_MonsterGroupRegen.xml", &cfg)
	if cfg.Notice == nil || (*cfg.Notice != 0 && *cfg.Notice != 1) {
		return fmt.Errorf("invalid SpawnNotice")
	}
	groups := make(map[int]*Group)
	for _, row := range cfg.Groups {
		if row.Index == nil || *row.Index < 0 || *row.Index >= 255 || groups[*row.Index] != nil || row.Map == nil || row.Boss == nil || row.Duration == nil || *row.Duration < 0 {
			return fmt.Errorf("invalid or duplicate group definition")
		}
		g := &Group{Index: *row.Index, mapNumber: *row.Map, boss: *row.Boss, duration: *row.Duration}
		if *cfg.Notice == 1 {
			g.notice = row.Notice
		}
		groups[g.Index] = g
	}
	seen := make(map[int]bool)
	for _, row := range cfg.Spots {
		if row.Index == nil || groups[*row.Index] == nil || seen[*row.Index] || len(row.Areas) == 0 || len(row.Areas) > 50 {
			return fmt.Errorf("invalid group areas reference/count")
		}
		seen[*row.Index] = true
		g := groups[*row.Index]
		areas := make(map[monster.EventSpawn]bool)
		for _, area := range row.Areas {
			if area.X == nil || area.Y == nil || area.EndX == nil || area.EndY == nil {
				return fmt.Errorf("missing group coordinates")
			}
			s := monster.EventSpawn{Class: g.boss, MapNumber: g.mapNumber, StartX: *area.X, StartY: *area.Y, EndX: *area.EndX, EndY: *area.EndY, Direction: -1, Distance: 30}
			if err := s.Validate(); err != nil {
				return err
			}
			if areas[s] {
				return fmt.Errorf("duplicate group area")
			}
			areas[s] = true
			g.areas = append(g.areas, s)
		}
	}
	seen = make(map[int]bool)
	for _, row := range cfg.Monsters {
		if row.Index == nil || groups[*row.Index] == nil || seen[*row.Index] || len(row.Members) == 0 || len(row.Members) > 20 {
			return fmt.Errorf("invalid group monsters reference/count")
		}
		seen[*row.Index] = true
		g := groups[*row.Index]
		if len(g.areas) == 0 {
			return fmt.Errorf("group %d missing areas", g.Index)
		}
		classes := make(map[int]bool)
		total := 0
		for _, r := range row.Members {
			if r.Index == nil || r.Count == nil || *r.Count <= 0 || *r.Count > 255 || classes[*r.Index] || r.Override == nil || (*r.Override != 0 && *r.Override != 1) {
				return fmt.Errorf("invalid group monster")
			}
			classes[*r.Index] = true
			total += *r.Count
			spawn := g.areas[0]
			spawn.Class = *r.Index
			if err := spawn.Validate(); err != nil {
				return err
			}
			entry := member{class: *r.Index, count: *r.Count, override: *r.Override == 1}
			if r.EventID != nil {
				if *r.EventID < 0 || *r.EventID == int(^uint(0)>>1) {
					return fmt.Errorf("invalid monster EventID %d", *r.EventID)
				}
				entry.eventID = *r.EventID + 1
				if err := drop.DropManager.LoadEventBag(entry.eventID); err != nil {
					return err
				}
			}
			if entry.override {
				if r.HP == nil || *r.HP <= 0 || r.Min == nil || *r.Min < 0 || r.Max == nil || *r.Max < *r.Min || r.Defense == nil || *r.Defense < 0 || r.AttackRate == nil || *r.AttackRate < 0 || r.DefenseRate == nil || *r.DefenseRate < 0 {
					return fmt.Errorf("invalid monster attribute override")
				}
				entry.hp, entry.damageMin, entry.damageMax, entry.defense, entry.attackRate, entry.defenseRate = *r.HP, *r.Min, *r.Max, *r.Defense, *r.AttackRate, *r.DefenseRate
			}
			g.members = append(g.members, entry)
		}
		if !classes[g.boss] || total > 255 || total > conf.Server.MaxMonsterCount {
			return fmt.Errorf("missing boss or excessive group size")
		}
	}
	seen = make(map[int]bool)
	for _, row := range cfg.Spawns {
		if row.Index == nil || groups[*row.Index] == nil || seen[*row.Index] || len(row.Starts) == 0 || len(row.Starts) > 24 {
			return fmt.Errorf("invalid group schedule reference/count")
		}
		seen[*row.Index] = true
		g := groups[*row.Index]
		slots := make(map[int]bool)
		for _, r := range row.Starts {
			if r.Hour == nil || *r.Hour < 0 || *r.Hour > 23 || r.Minute == nil || *r.Minute < 0 || *r.Minute > 59 || r.Rate == nil || *r.Rate < 0 || *r.Rate > 100 {
				return fmt.Errorf("invalid group schedule")
			}
			slot := *r.Hour*60 + *r.Minute
			if slots[slot] {
				return fmt.Errorf("duplicate group schedule")
			}
			slots[slot] = true
			g.Starts = append(g.Starts, StartTime{*r.Hour, *r.Minute, *r.Rate})
		}
	}
	var next []*Group
	for _, g := range groups {
		if len(g.areas) == 0 || len(g.members) == 0 || len(g.Starts) == 0 {
			return fmt.Errorf("incomplete group %d", g.Index)
		}
		next = append(next, g)
	}
	sort.Slice(next, func(i, j int) bool { return next[i].Index < next[j].Index })
	m.Groups = next
	return nil
}
func (g *Group) Running() bool {
	for _, obj := range g.monsters {
		if object.ObjectManager.GetObject(obj.Index) == obj && obj.Live {
			return true
		}
	}
	return false
}
func (g *Group) Start(now time.Time) error {
	if len(g.areas) == 0 || len(g.members) == 0 {
		return fmt.Errorf("empty group %d", g.Index)
	}
	g.reap(now, true)
	area := g.areas[rand.Intn(len(g.areas))]
	var created []*object.Object
	for _, entry := range g.members {
		for i := 0; i < entry.count; i++ {
			spawn := area
			spawn.Class = entry.class
			spawn.EventID = entry.eventID
			obj, err := monster.SpawnEventMonster(spawn)
			if err != nil {
				for _, previous := range created {
					object.ObjectManager.DeleteEventMonster(previous, now)
				}
				return err
			}
			if entry.override {
				obj.HP, obj.MaxHP = entry.hp, entry.hp
				obj.AttackMin, obj.AttackMax = entry.damageMin, entry.damageMax
				obj.Defense, obj.AttackRate, obj.DefenseRate = entry.defense, entry.attackRate, entry.defenseRate
			}
			created = append(created, obj)
		}
	}
	g.monsters = append(g.monsters, created...)
	if g.notice != "" {
		object.ObjectManager.BroadcastSystemMsg(g.notice)
	}
	return nil
}
func (g *Group) Tick(now time.Time) { g.reap(now, false) }
func (g *Group) reap(now time.Time, living bool) {
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
