package main

import "sort"

type CalendarItem struct {
	CourseID   string
	CourseCode string
	CourseName string

	Methods []string

	ModuleStart string
	ModuleEnd   string

	Day Day

	PeriodStart int
	PeriodEnd   int
}

func GetModuleNames(modules []Module) []string {
	s := make([]string, len(modules))
	for i, m := range modules {
		s[i] = m.Name
	}
	return s
}

func ConvertToCalendarItem(c Course, modules []string) []CalendarItem {
	moduleIndex := make(map[string]int, len(modules))
	for i, s := range modules {
		moduleIndex[s] = i
	}

	items := make([]CalendarItem, len(c.Schedules))
	for i, s := range c.Schedules {
		items[i] = CalendarItem{
			CourseID:   c.ID,
			CourseCode: c.Code,
			CourseName: c.Name,
			Methods:    c.Methods,

			ModuleStart: s.Module,
			ModuleEnd:   s.Module,
			Day:         s.Day,
			PeriodStart: s.Period,
			PeriodEnd:   s.Period,
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].ModuleStart != items[j].ModuleStart {
			return moduleIndex[items[i].ModuleStart] < moduleIndex[items[j].ModuleStart]
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
			if item.Day != item.Day {
				continue
			}
			if item.ModuleStart != v.ModuleStart {
				continue
			}
			if item.PeriodEnd+1 != v.PeriodEnd {
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
		end := moduleIndex[item.ModuleStart]
		for j := i + 1; j < len(items); j++ {
			v := items[j]
			if item.Day != item.Day {
				continue
			}
			if item.PeriodStart != v.PeriodStart {
				continue
			}
			if item.PeriodEnd != v.PeriodEnd {
				continue
			}
			if end+1 != moduleIndex[v.ModuleStart] {
				continue
			}
			end++
			removed[j] = struct{}{}
		}
		item.ModuleEnd = modules[end]
		items[i] = item
	}

	result := make([]CalendarItem, 0, len(items)-len(removed))
	for i, item := range items {
		if _, ok := removed[i]; ok {
			continue
		}
		result = append(result, item)
	}
	return result
}
