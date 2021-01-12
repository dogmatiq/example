// Package validation contains some message helpers.
package validation

import "github.com/dogmatiq/example/messages"

// IsValidBusinessDate returns true if the given date is a valid format.
func IsValidBusinessDate(date string) bool {
	if date == "" {
		return false
	}

	_, err := messages.UnmarshalDate(date)

	return err == nil
}
