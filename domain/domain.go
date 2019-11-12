package domain

import "time"

const dateFormat = "2006-01-02"

func businessDayFromTime(t time.Time) string {
	return t.Format(dateFormat)
}

func startOfBusinessDay(date string) time.Time {
	t, err := time.Parse(dateFormat, date)
	if err != nil {
		panic(err)
	}

	return t
}

func mustValidateDate(date string) {
	_, err := time.Parse(dateFormat, date)
	if err != nil {
		panic(err)
	}
}
