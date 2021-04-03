package main

import (
	"time"
)

type Day string

const (
	Sun Day = "Sun"
	Mon Day = "Mon"
	Tue Day = "Tue"
	Wed Day = "Wed"
	Thu Day = "Thu"
	Fri Day = "Fri"
	Sat Day = "Sat"
)

func (d Day) Valid() bool {
	switch d {
	case Sun, Mon, Tue, Wed, Thu, Fri, Sat:
		return true
	default:
		return false
	}
}

func (d Day) ToWeekday() time.Weekday {
	switch d {
	case Sun:
		return time.Sunday
	case Mon:
		return time.Monday
	case Tue:
		return time.Tuesday
	case Wed:
		return time.Wednesday
	case Thu:
		return time.Thursday
	case Fri:
		return time.Friday
	case Sat:
		return time.Saturday
	default:
		panic("never happen")
	}
}

var tz *time.Location

func init() {
	tz, _ = time.LoadLocation("Asia/Tokyo") // TODO: Catch error gracefully
}

type Date struct {
	Time time.Time
}

func (d Date) Day() Day {
	switch d.Time.Weekday() {
	case time.Sunday:
		return Sun
	case time.Monday:
		return Mon
	case time.Tuesday:
		return Tue
	case time.Wednesday:
		return Wed
	case time.Thursday:
		return Thu
	case time.Friday:
		return Fri
	case time.Saturday:
		return Sat
	default:
		panic("never happen")
	}
}

func (d Date) In(start, end Date) bool {
	return !d.Time.Before(start.Time) && !d.Time.After(end.Time)
}

func (d Date) Next(day Day) Date {
	wd := day.ToWeekday()
	for i := 0; i < 7; i++ {
		t := d.Time.Add(time.Duration(i*24) * time.Hour)
		if t.Weekday() == wd {
			return Date{Time: t}
		}
	}
	panic("never happen")
}

func (d Date) ToTime(hour, minute int) time.Time {
	return d.Time.Add(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute)
}

func (d Date) MarshalJSON() ([]byte, error) {
	s := d.Time.Format(`"2006-01-02T00:00:00.000Z"`)
	return []byte(s), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	d.Time, err = time.ParseInLocation(`"2006-01-02T00:00:00.000Z"`, string(data), tz)
	return err
}
