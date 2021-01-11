package domain

import (
	"time"

	"github.com/dogmatiq/example/messages"
)

func startOfBusinessDay(date string) time.Time {
	t, err := time.Parse(messages.BusinessDateFormat, date)
	if err != nil {
		panic(err)
	}

	return t
}
