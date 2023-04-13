package main

import (
	"context"
	"sort"
	"strconv"
)

type Module struct {
	Name  string `json:"module,omitempty"`
	Start Date   `json:"start,omitempty"`
	End   Date   `json:"end,omitempty"`

	Exceptions map[Day][]Date //`json:"-"`
	Additions  map[Day][]Date //`json:"-"`
}

func (m Module) addException(day Day, date Date) {
	for _, d := range m.Exceptions[day] {
		if d.Time.Equal(date.Time) {
			return // Already exists
		}
	}
	m.Exceptions[day] = append(m.Exceptions[day], date)
}

func GetSchoolCalendar(ctx context.Context, year int) ([]Module, error) {
	var ms []Module
	err := GetAPI(ctx, "/school-calendar/modules?year="+strconv.Itoa(year), &ms)
	if err != nil {
		return nil, err
	}

	for i, m := range ms {
		m.Exceptions = make(map[Day][]Date, 7)
		m.Additions = make(map[Day][]Date, 7)
		ms[i] = m
	}

	var es []struct {
		Date        Date   `json:"date,omitempty"`
		EventType   string `json:"eventType,omitempty"`
		Description string `json:"description,omitempty"`
		ChangeTo    *Day   `json:"changeTo,omitempty"`
	}
	err = GetAPI(ctx, "/school-calendar/events?year="+strconv.Itoa(year), &es)
	if err != nil {
		return nil, err
	}

	for _, e := range es {
		if e.EventType == "Other" {
			continue
		}
		for _, m := range ms {
			if !e.Date.In(m.Start, m.End) {
				continue
			}
			if e.EventType == "Exam" && m.Name != "SpringA" && m.Name != "FallA" {
				continue
			}
			if e.ChangeTo != nil {
				day := *e.ChangeTo
				if day != e.Date.Day() { // Do not add exception and addition to same date
					m.addException(e.Date.Day(), e.Date)
					m.Additions[day] = append(m.Additions[day], e.Date)
				}
			} else {
				m.addException(e.Date.Day(), e.Date)
			}
		}
	}

	sort.Slice(ms, func(i, j int) bool {
		return ms[i].Start.Time.Before(ms[j].Start.Time)
	})

	return ms, nil
}
