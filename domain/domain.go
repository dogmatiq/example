package domain

import "time"

const timeZone = "UTC"

func startOfBusinessDay(t time.Time) time.Time {
	x := t.In(time.FixedZone(timeZone, 0))

	return time.Date(
		x.Year(),
		x.Month(),
		x.Day(),
		0,
		0,
		0,
		0,
		x.Location(),
	)
}
