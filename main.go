package main

import (
	"context"
	"fmt"
)

func main() {
	sc, err := GetSchoolCalendar(context.Background(), 2021)
	if err != nil {
		panic(err)
	}
	cs, err := GetCourses(context.Background(), 2021)
	if err != nil {
		panic(err)
	}
	modules := GetModuleNames(sc)
	for _, c := range cs {
		items := ConvertToCalendarItem(c, modules)

		fmt.Printf("%s %s:\n", c.Code, c.Name)
		for _, i := range items {
			fmt.Printf("- %s %d-%d (%s - %s)\n", i.Day, i.PeriodStart, i.PeriodEnd, i.ModuleStart, i.ModuleEnd)
		}
	}
}
