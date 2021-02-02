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
	Date          string
}

// MessageDescription returns a human-readable description of the message.
func (m ConsumeDailyDebitLimit) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: consuming %s from %s daily debit limit of account %s",
		m.DebitType,
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.Date,
		m.AccountID,
	)
}
