package messages

import (
	"fmt"
	"time"
)

// FormatAmount formats a cent amount as dollars.
func FormatAmount(v int64) string {
	f := "$%d.%02d"
	if v < 0 {
		v = -v
		f = "-" + f
	}

	return fmt.Sprintf(f, v/100, v%100)
}

// FormatDate formats a time value as a date.
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}
