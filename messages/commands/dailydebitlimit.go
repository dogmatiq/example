package commands

import (
	"fmt"

	"github.com/dogmatiq/example/messages"
)

// ConsumeDailyDebitLimit is a command requesting that an amount of an account
// daily debit limit be consumed.
type ConsumeDailyDebitLimit struct {
	TransactionID string
	AccountID     string
	DebitType     messages.TransactionType
	Amount        int64
	ScheduledDate string
}

// MessageDescription returns a human-readable description of the message.
func (m *ConsumeDailyDebitLimit) MessageDescription() string {
	return fmt.Sprintf(
		"consuming %s from daily debit limit of account %s on date %s",
		messages.FormatAmount(m.Amount),
		messages.FormatID(m.AccountID),
		m.ScheduledDate,
	)
}
