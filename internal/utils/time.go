package utils

import (
	"time"
)

func GetUnixMillisecondsByUnixNano(value int64) int64 {
	return value / int64(time.Millisecond)
}

func GetStartOfTheDay(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}

func GetStartOfTheMonth(t time.Time, loc *time.Location) time.Time {
	y, m, _ := t.In(loc).Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, loc)
}

func GetEndOfTheDay(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), loc)
}

func GetPreviousDayAtMidnight(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d-1, 0, 0, 0, 0, loc)
}

func GetNextDayAtMidnight(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d+1, 0, 0, 0, 0, loc)
}

func GetPreviousOrNextDayAtMidnight(t time.Time, day int, loc *time.Location) time.Time {
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d+day, 0, 0, 0, 0, loc)
}

func GetStartOfTheWeek(t time.Time, loc *time.Location) time.Time {
	// Roll back to Monday:
	y, m, d := t.In(loc).Date()
	if wd := t.Weekday(); wd == time.Sunday {
		t = time.Date(y, m, d-6, 0, 0, 0, 0, loc)
	} else {
		t = time.Date(y, m, d-int(wd-1), 0, 0, 0, 0, loc)
	}
	return t
}

func GetEndOfTheWeek(t time.Time, loc *time.Location) time.Time {
	// Move to Monday:
	y, m, d := t.In(loc).Date()
	if wd := t.Weekday(); wd == time.Sunday {
		return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), loc)
	} else {
		return time.Date(y, m, d+(6-int(wd-1)), 23, 59, 59, int(time.Second-time.Nanosecond), loc)
	}
}
