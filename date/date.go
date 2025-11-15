package date

import (
	"time"
)

type Date time.Time

func NewDate(year int, month time.Month, day int) Date {
	return Date(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

func NewDateWithTime(tm time.Time) Date {
	return NewDate(tm.Year(), tm.Month(), tm.Day())
}

func Today() Date {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return Date(today)
}

func (d Date) Equal(other Date) bool {
	dy, dm, dd := time.Time(d).Date()
	oy, om, od := time.Time(other).Date()

	return dy == oy && dm == om && dd == od
}

func (d Date) String() string {
	return time.Time(d).Format("2006-01-02")
}

func (d Date) BeginningOfDate() time.Time {
	tm := time.Time(d)

	return time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, time.UTC)
}

func (d Date) EndOfDate() time.Time {
	tm := time.Time(d)

	return time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, time.UTC)
}

func (d Date) Time() time.Time {
	return time.Time(d)
}

func (d Date) InWeek(beginOfWeek Date) bool {
	// ISOWeekは月曜日はじまり
	y1, w1 := beginOfWeek.Time().ISOWeek()
	y2, w2 := d.Time().ISOWeek()

	return y1 == y2 && w1 == w2
}
