// Package event provides scheduling and synchronous advancement of activities.
// EventManager loads schedules and registers activities during package initialization.
package event

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/event/invasion"
)

type Activity interface {
	Start(time.Time) error
	Tick(time.Time)
	Running() bool
}

type Schedule struct {
	ActivityID int
	DayOfWeek  int // -1 means every day; otherwise Sunday=0.
	Hour       int
	Minute     int
}

type occurrence struct {
	id int
	at time.Time
}

// eventManager is advanced by the game loop, never concurrently.
type eventManager struct {
	enabled    bool
	schedules  []Schedule
	activities map[int]Activity
	ids        []int
	last       time.Time
}

var EventManager eventManager

func init() {
	EventManager.init()
}

func (m *eventManager) init() {
	if err := m.loadSchedules(conf.PathCommon); err != nil {
		panic(fmt.Errorf("event: %w", err))
	}
	m.activities = make(map[int]Activity)
	if err := m.register(invasion.RedDragon, invasion.NewDragon()); err != nil {
		panic(err)
	}
}

func newManager(enabled bool, schedules []Schedule) (*eventManager, error) {
	if err := ValidateSchedules(schedules); err != nil {
		return nil, err
	}
	m := &eventManager{enabled: enabled, schedules: append([]Schedule(nil), schedules...), activities: make(map[int]Activity)}
	return m, nil
}

func (m *eventManager) register(id int, activity Activity) error {
	if id < 0 || activity == nil || m.activities[id] != nil {
		return fmt.Errorf("invalid or duplicate event registration: %d", id)
	}
	m.activities[id] = activity
	m.ids = append(m.ids, id)
	sort.Ints(m.ids)
	return nil
}

func (m *eventManager) loadSchedules(basePath string) error {
	// Pointer attributes distinguish missing required values from valid zeroes.
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
	var times timetable
	conf.XML(basePath, "Events/IGC_InvasionManager.xml", &times)
	if times.Enable == nil || (*times.Enable != 0 && *times.Enable != 1) {
		return fmt.Errorf("InvasionManager requires Enable=0 or 1")
	}
	enabled := *times.Enable == 1
	var schedules []Schedule
	ids := make(map[int]bool)
	for _, inv := range times.Invasions {
		if inv.Type == nil || *inv.Type < 0 || *inv.Type > 4 || ids[*inv.Type] || len(inv.Starts) == 0 {
			return fmt.Errorf("missing, invalid or duplicate invasion type/time definitions")
		}
		ids[*inv.Type] = true
		for _, start := range inv.Starts {
			if start.Day == nil || start.Hour == nil || start.Minute == nil {
				return fmt.Errorf("invasion %d: missing start time attribute", *inv.Type)
			}
			s := Schedule{ActivityID: *inv.Type, DayOfWeek: *start.Day, Hour: *start.Hour, Minute: *start.Minute}
			schedules = append(schedules, s)
		}
	}
	if err := ValidateSchedules(schedules); err != nil {
		return err
	}
	m.enabled = enabled
	m.schedules = schedules
	return nil
}

func (s Schedule) Validate() error {
	if s.ActivityID < 0 || s.DayOfWeek < -1 || s.DayOfWeek > 6 || s.Hour < 0 || s.Hour > 23 || s.Minute < 0 || s.Minute > 59 {
		return fmt.Errorf("invalid event schedule: %+v", s)
	}
	return nil
}

func ValidateSchedules(schedules []Schedule) error {
	for i, s := range schedules {
		if err := s.Validate(); err != nil {
			return err
		}
		for _, previous := range schedules[:i] {
			if s.ActivityID == previous.ActivityID && s.Hour == previous.Hour && s.Minute == previous.Minute &&
				(s.DayOfWeek == previous.DayOfWeek || s.DayOfWeek == -1 || previous.DayOfWeek == -1) {
				return fmt.Errorf("overlapping event schedules: %+v and %+v", previous, s)
			}
		}
	}
	return nil
}

// Tick establishes the baseline on its first call, without replaying earlier schedules.
func (m *eventManager) Tick(now time.Time) {
	if m.last.IsZero() {
		m.last = now.Round(0)
		return
	}
	// Keep a wall-clock high-water mark on rollback. Strip monotonic readings
	// from last, so time.Now's monotonic clock cannot hide a wall-clock rollback.
	// Failed/skipped occurrences are consumed too.
	if now.After(m.last) {
		if m.enabled {
			for _, occurrence := range due(m.schedules, m.last, now) {
				activity := m.activities[occurrence.id]
				if activity == nil {
					continue
				}
				if activity.Running() {
					slog.Info("event start skipped: already running", "event", occurrence.id, "scheduled", occurrence.at)
					continue
				}
				if err := activity.Start(now); err != nil {
					slog.Error("event start failed", "event", occurrence.id, "scheduled", occurrence.at, "err", err)
				}
			}
		}
		m.last = now.Round(0)
	}
	// Idle activities can still have dead monsters awaiting death settlement.
	for _, id := range m.ids {
		m.activities[id].Tick(now)
	}
}

// due returns occurrences in (last, now], in chronological order. Dates are
// constructed in the server's location rather than adding 24-hour durations.
func due(schedules []Schedule, last, now time.Time) []occurrence {
	var result []occurrence
	day := time.Date(last.Year(), last.Month(), last.Day(), 0, 0, 0, 0, time.Local)
	for !day.After(now) {
		for _, s := range schedules {
			if s.DayOfWeek != -1 && s.DayOfWeek != int(day.Weekday()) {
				continue
			}
			at := time.Date(day.Year(), day.Month(), day.Day(), s.Hour, s.Minute, 0, 0, time.Local)
			if at.After(last) && !at.After(now) {
				result = append(result, occurrence{id: s.ActivityID, at: at})
			}
		}
		day = day.AddDate(0, 0, 1)
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].at.Before(result[j].at) })
	return result
}
