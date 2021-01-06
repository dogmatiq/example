package events

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/example/messages"
)

// DailyDebitLimitConsumed is an event that indicates an amount of an account
// daily debit limit has been consumed.
type DailyDebitLimitConsumed struct {
	TransactionID     string
	AccountID         string
	DebitType         messages.TransactionType
	Amount            int64
	Date              string
	TotalDebitsForDay int64
	DailyLimit        int64
}

// DailyDebitLimitExceeded is an event that indicates an attempt to consume from
// an account daily debit limit has been rejected due to reaching the limit.
type DailyDebitLimitExceeded struct {
	TransactionID     string
	AccountID         string
	DebitType         messages.TransactionType
	Amount            int64
	Date              string
	TotalDebitsForDay int64
	DailyLimit        int64
}

// MessageDescription returns a human-readable description of the message.
func (m DailyDebitLimitConsumed) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: consumed %s from %s daily debit limit of account %s",
		m.DebitType,
		m.TransactionID,
		m.Date,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m DailyDebitLimitExceeded) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: exceeded %s daily debit limit of account %s by %s",
		m.DebitType,
		m.TransactionID,
		m.Date,
		m.AccountID,
		messages.FormatAmount((m.TotalDebitsForDay+m.Amount)-m.DailyLimit),
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m DailyDebitLimitConsumed) Validate() error {
	if m.TransactionID == "" {
		return errors.New("DailyDebitLimitConsumed needs a valid transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("DailyDebitLimitConsumed needs a valid account ID")
	}
	if m.DebitType == "" {
		return errors.New("DailyDebitLimitConsumed needs a valid debit type")
	}
	if m.Amount < 1 {
		return errors.New("DailyDebitLimitConsumed needs a valid amount")
	}
	if !messages.IsValidBusinessDate(m.Date) {
		return errors.New("DailyDebitLimitConsumed needs a valid date")
	}
	if m.TotalDebitsForDay < 1 {
		return errors.New("DailyDebitLimitConsumed needs a valid total debits for day")
	}
	if m.DailyLimit < 0 {
		return errors.New("DailyDebitLimitConsumed needs a valid daily limit")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m DailyDebitLimitExceeded) Validate() error {
	if m.TransactionID == "" {
		return errors.New("DailyDebitLimitExceeded needs a valid transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("DailyDebitLimitExceeded needs a valid account ID")
	}
	if m.DebitType == "" {
		return errors.New("DailyDebitLimitExceeded needs a valid debit type")
	}
	if m.Amount < 1 {
		return errors.New("DailyDebitLimitExceeded needs a valid amount")
	}
	if !messages.IsValidBusinessDate(m.Date) {
		return errors.New("DailyDebitLimitExceeded needs a valid date")
	}
	if m.TotalDebitsForDay < 0 {
		return errors.New("DailyDebitLimitExceeded needs a valid total debits for day")
	}
	if m.DailyLimit < 0 {
		return errors.New("DailyDebitLimitExceeded needs a valid daily limit")
	}

	return nil
}
