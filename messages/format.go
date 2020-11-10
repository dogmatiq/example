package messages

import "fmt"

// FormatAmount formats a cent amount as dollars.
func FormatAmount(v int64) string {
	f := "-$%d.%02d"
	if v < 0 {
		v = -v
		f = "-" + f
	}

	return fmt.Sprintf(f, v/100, v%100)
}

// FormatID returns a compact rendering of ID for use in log messages and other
// human-readable strings.
// It returns '<unidentified>' if id is an empty string.
func FormatID(id string) string {
	if looksLikeUUID(id) {
		return id[:uuidSep1] + id[uuidLen:]
	}

	if id == "" {
		return "<unidentified>"
	}

	return id
}

func looksLikeUUID(id string) bool {
	if len(id) < uuidLen {
		return false
	}

	return id[uuidSep1] == '-' &&
		id[uuidSep2] == '-' &&
		id[uuidSep3] == '-' &&
		id[uuidSep4] == '-'
}

const (
	uuidLen  = 36
	uuidSep1 = 8
	uuidSep2 = 13
	uuidSep3 = 18
	uuidSep4 = 23
)
