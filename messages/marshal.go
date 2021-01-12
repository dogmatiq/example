package messages

import "time"

// businessDateFormat is the format used for business dates in messages.
const businessDateFormat = "2006-01-02"

// MarshalDate marshals a time value to a date string.
func MarshalDate(t time.Time) string {
	return t.Format(businessDateFormat)
}

// UnmarshalDate unmarshals a date string to a time value.
func UnmarshalDate(d string) (time.Time, error) {
	return time.Parse(businessDateFormat, d)
}
