package messages

import "time"

// BusinessDateFormat is the format used for business dates in messages.
const BusinessDateFormat = "2006-01-02"

// IsValidBusinessDate returns true if the given date is a valid format.
func IsValidBusinessDate(date string) bool {
	if date == "" {
		return false
	}

	_, err := time.Parse(BusinessDateFormat, date)

	return err == nil
}
