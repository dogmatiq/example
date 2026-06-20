package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// money represents a monetary amount in cents.
type money int64

func (m money) String() string {
	v := int64(m)
	prefix := "$"

	if v < 0 {
		v = -v
		prefix = "-$"
	}

	p := message.NewPrinter(language.English)
	return p.Sprintf("%s%d.%02d", prefix, v/100, v%100)
}

// parseMoney parses a string like "100.00" or "100" into cents.
func parseMoney(s string) (money, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "$")

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	if f <= 0 {
		return 0, fmt.Errorf("amount must be positive")
	}

	return money(int64(math.Round(f * 100))), nil
}
