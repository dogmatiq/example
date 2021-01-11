package validation

import (
	"time"

	"github.com/dogmatiq/example/messages"
)

// IsValidBusinessDate returns true if the given date is a valid format.
func IsValidBusinessDate(date string) bool {
	if date == "" {
		return false
	}

	_, err := time.Parse(messages.BusinessDateFormat, date)

	return err == nil
}
