package main

import (
	"sort"
	"time"
)

type Schedule struct {
	StartTime time.Time
	EndTime   time.Time

	Day   Day
	Until time.Time

	Exceptions []time.Time
	Additions  []time.Time

	Location string
}

func GetSchedules(modules []Module, cs []CourseSchedule) []Schedule {
	moduleIndex := make(map[string]int, len(modules))
	for i, m := range modules {
		moduleIndex[m.Name] = i
	}

	type item struct {
		ModuleStart int
		ModuleEnd   int

		Day Day

		PeriodStart int
		PeriodEnd   int

		Location string
	}

	items := make([]item, 0, len(cs))
	for _, s := range cs {
		if !s.Day.Valid() || s.Period < 1 || s.Period > 6 { // TODO: Support 7~8 period
			continue
		}
		items = append(items, item{
			ModuleStart: moduleIndex[s.Module],
			ModuleEnd:   moduleIndex[s.Module],
			Day:         s.Day,
			PeriodStart: s.Period,
			PeriodEnd:   s.Period,
			Location:    s.Room,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].ModuleStart != items[j].ModuleStart {
			return items[i].ModuleStart < items[j].ModuleStart
		}
		return items[i].PeriodStart < items[j].PeriodStart
	})

	removed := make(map[int]struct{}, len(items))

	for i, item := range items {
		if _, ok := removed[i]; ok {
			continue
		}
		for j := i + 1; j < len(items); j++ {
			v := items[j]
			if item.Day != v.Day {
				continue
			}
			if item.ModuleStart != v.ModuleStart {
				continue
			}
			if item.PeriodEnd+1 != v.PeriodStart {
				continue
			}
			item.PeriodEnd = v.PeriodEnd
			removed[j] = struct{}{}
		}
		items[i] = item
	}
	for i, item := range items {
		if _, ok := removed[i]; ok {
			continue
		}
		for j := i + 1; j < len(items); j++ {
			v := items[j]
			if item.Day != v.Day {
				continue
			}
			if item.PeriodStart != v.PeriodStart {
				continue
			}
			if item.PeriodEnd != v.PeriodEnd {
				continue
			}
			if item.ModuleEnd+1 != v.ModuleStart {
				continue
			}
			item.ModuleEnd = v.ModuleEnd
			removed[j] = struct{}{}
		}
		items[i] = item
	}

	result := make([]Schedule, 0, len(items)-len(removed))
	for i, item := range items {
		if _, ok := removed[i]; ok {
			continue
		}

		var startTime time.Time
		var endTime time.Time
		var exceptions []time.Time
		var additions []time.Time
		var until time.Time

		for i, m := range modules {
			if i == item.ModuleStart {
				date := m.Start.Next(item.Day)
				startTime = date.ToTime(GetPeriodStart(item.PeriodStart))
				endTime = date.ToTime(GetPeriodEnd(item.PeriodEnd))
			}
			if startTime.IsZero() {
				continue
			}
			if i > item.ModuleStart {
				d := modules[i-1].End.NextDay()
				for !d.Time.Equal(m.Start.Time) {
					if d.Day() == item.Day {
						exceptions = append(exceptions, d.ToTime(GetPeriodStart(item.PeriodStart)))
					}
					d = d.NextDay()
				}
			}
			for _, d := range m.Exceptions[item.Day] {
				exceptions = append(exceptions, d.ToTime(GetPeriodStart(item.PeriodStart)))
			}
			for _, d := range m.Additions[item.Day] {
				additions = append(additions, d.ToTime(GetPeriodStart(item.PeriodStart)))
			}
			if i == item.ModuleEnd {
				until = m.End.ToTime(23, 59)
				break
			}
		}

		result = append(result, Schedule{
			StartTime:  startTime,
			EndTime:    endTime,
			Day:        item.Day,
			Until:      until,
			Exceptions: exceptions,
			Additions:  additions,
			Location:   item.Location,
		})
	}
	return result
}
