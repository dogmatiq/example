// Package validation contains some message validation helpers.
package validation

import "time"

// IsValidDate returns true if the given date is a valid format.
func IsValidDate(date string) bool {
	if date == "" {
		return false
	}

	_, err := time.Parse("2006-01-02", date)

	return err == nil
}
