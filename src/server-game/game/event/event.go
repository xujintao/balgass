// Package event schedules activities on the game loop.
package event

import (
	"fmt"
	"log/slog"
	"math/rand"
	"sort"
	"time"
)

type Event interface {
	Start(time.Time) error
	Tick(time.Time)
	Running() bool
}

type Schedule struct {
	DayOfWeek      int // -1 means every day; Sunday=0.
	Hour, Minute   int
	AppearanceRate int // 0..100; evaluated once for each scheduled occurrence.
}

type occurrence struct {
	at   time.Time
	rate int
}

type eventEntry struct {
	name      string
	activity  Event
	schedules []Schedule
}

// All registration and activity operations run synchronously on the game loop.
type eventManager struct {
	activities []eventEntry
	last       time.Time
}

var EventManager eventManager

// Register adds an event during startup, before the game loop begins.
func (m *eventManager) Register(name string, activity Event, schedules []Schedule) error {
	if name == "" || activity == nil {
		return fmt.Errorf("invalid event registration: %q", name)
	}
	for _, registered := range m.activities {
		if registered.name == name {
			return fmt.Errorf("duplicate event registration: %s", name)
		}
	}
	if err := ValidateSchedules(schedules); err != nil {
		return err
	}
	m.activities = append(m.activities, eventEntry{name: name, activity: activity, schedules: append([]Schedule(nil), schedules...)})
	return nil
}

func (s Schedule) Validate() error {
	if s.DayOfWeek < -1 || s.DayOfWeek > 6 || s.Hour < 0 || s.Hour > 23 || s.Minute < 0 || s.Minute > 59 || s.AppearanceRate < 0 || s.AppearanceRate > 100 {
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
			if s.Hour == previous.Hour && s.Minute == previous.Minute && (s.DayOfWeek == previous.DayOfWeek || s.DayOfWeek == -1 || previous.DayOfWeek == -1) {
				return fmt.Errorf("overlapping event schedules: %+v and %+v", previous, s)
			}
		}
	}
	return nil
}

// The first Tick establishes the baseline without replaying past schedules.
func (m *eventManager) Tick(now time.Time) {
	if m.last.IsZero() {
		m.last = now.Round(0)
		return
	}
	if now.After(m.last) {
		for _, entry := range m.activities {
			for _, occurrence := range due(entry.schedules, m.last, now) {
				if occurrence.rate == 0 || (occurrence.rate < 100 && rand.Intn(100) >= occurrence.rate) {
					continue
				}
				if err := entry.activity.Start(now); err != nil {
					slog.Error("event start failed", "event", entry.name, "scheduled", occurrence.at, "err", err)
				}
			}
		}
		// Never rewind the baseline; failures and skipped occurrences are consumed.
		m.last = now.Round(0)
	}
	for _, entry := range m.activities {
		entry.activity.Tick(now)
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
				result = append(result, occurrence{at: at, rate: s.AppearanceRate})
			}
		}
		day = day.AddDate(0, 0, 1)
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].at.Before(result[j].at) })
	return result
}
