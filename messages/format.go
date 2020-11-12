package messages

import (
	"fmt"
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
