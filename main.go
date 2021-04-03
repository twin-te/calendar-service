package main

import (
	"context"
	"os"
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
	err = WriteICalendar(os.Stdout, modules, cs)
	if err != nil {
		panic(err)
	}
}
