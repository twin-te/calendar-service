package main

import (
	"context"
	"fmt"
)

func main() {
	modules, err := GetSchoolCalendar(context.Background(), 2021)
	if err != nil {
		panic(err)
	}
	cs, err := GetCourses(context.Background(), 2021)
	if err != nil {
		panic(err)
	}
	for _, c := range cs {
		ss := GetSchedules(modules, c.Schedules)

		fmt.Printf("%s %s:\n", c.Code, c.Name)
		for _, s := range ss {
			fmt.Printf("* %s - %s (%s)\n", s.StartTime, s.EndTime, s.Until)
			for _, t := range s.Exceptions {
				fmt.Printf("  - %s\n", t)
			}
			for _, t := range s.Additionals {
				fmt.Printf("  + %s\n", t)
			}
		}
	}
}
