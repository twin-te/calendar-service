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
