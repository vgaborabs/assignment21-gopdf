package types

import (
	"strings"
	"time"
)

type DateTime struct {
	time.Time
}

func (d *DateTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		d.Time = time.Time{}
		return
	}
	if len(s) == 10 {
		d.Time, err = time.Parse("2006-01-02", s)
	}
	return d.Time.UnmarshalJSON(b)
}

func NewDate(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) DateTime {
	return DateTime{
		Time: time.Date(year, month, day, hour, min, sec, nsec, loc),
	}
}
