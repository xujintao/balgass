package invasion

import (
	"encoding/xml"
	"fmt"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/event"
)

const (
	RedDragon    = 0
	SkeletonKing = 1
)

type StartTime struct{ DayOfWeek, Hour, Minute int }
type invasionManager struct {
	Enabled   bool
	Schedules map[int][]StartTime
}

var InvasionManager invasionManager

func init() { InvasionManager.init() }
func (m *invasionManager) init() {
	if err := m.load(conf.PathCommon); err != nil {
		panic(fmt.Errorf("invasion: %w", err))
	}
	// Load concrete configurations before constructing or registering events.
	dragonSettings.init()
	skeletonSettings.init()

	// register invasion events
	if !m.Enabled {
		return
	}
	for _, entry := range []struct {
		id       int
		name     string
		activity event.Event
	}{
		{RedDragon, "invasion/dragon", NewDragon()},
		{SkeletonKing, "invasion/skeleton", NewSkeleton()},
	} {
		var schedules []event.Schedule
		for _, start := range m.Schedules[entry.id] {
			schedules = append(schedules, event.Schedule{DayOfWeek: start.DayOfWeek, Hour: start.Hour, Minute: start.Minute, AppearanceRate: 100})
		}
		if err := event.EventManager.Register(entry.name, entry.activity, schedules); err != nil {
			panic(err)
		}
	}
}

func (m *invasionManager) load(basePath string) error {
	type timetable struct {
		XMLName   xml.Name `xml:"InvasionManager"`
		Enable    *int     `xml:"Enable,attr"`
		Invasions []struct {
			Type   *int `xml:"Type,attr"`
			Starts []struct {
				Day    *int `xml:"DayOfWeek,attr"`
				Hour   *int `xml:"Hour,attr"`
				Minute *int `xml:"Minute,attr"`
			} `xml:"Start"`
		} `xml:"Invasion"`
	}
	var cfg timetable
	conf.XML(basePath, "Events/IGC_InvasionManager.xml", &cfg)
	if cfg.Enable == nil || (*cfg.Enable != 0 && *cfg.Enable != 1) {
		return fmt.Errorf("InvasionManager requires Enable=0 or 1")
	}
	schedules := make(map[int][]StartTime)
	for _, inv := range cfg.Invasions {
		if inv.Type == nil || *inv.Type < 0 || *inv.Type > 4 || len(inv.Starts) == 0 {
			return fmt.Errorf("invalid invasion type or empty schedule")
		}
		if _, exists := schedules[*inv.Type]; exists {
			return fmt.Errorf("duplicate invasion type %d", *inv.Type)
		}
		var starts []StartTime
		for _, start := range inv.Starts {
			if start.Day == nil || start.Hour == nil || start.Minute == nil || *start.Day < -1 || *start.Day > 6 || *start.Hour < 0 || *start.Hour > 23 || *start.Minute < 0 || *start.Minute > 59 {
				return fmt.Errorf("invalid invasion %d start time", *inv.Type)
			}
			s := StartTime{*start.Day, *start.Hour, *start.Minute}
			for _, previous := range starts {
				if s.Hour == previous.Hour && s.Minute == previous.Minute && (s.DayOfWeek == previous.DayOfWeek || s.DayOfWeek == -1 || previous.DayOfWeek == -1) {
					return fmt.Errorf("overlapping invasion %d start times", *inv.Type)
				}
			}
			starts = append(starts, s)
		}
		schedules[*inv.Type] = starts
	}
	m.Enabled, m.Schedules = *cfg.Enable == 1, schedules
	return nil
}
