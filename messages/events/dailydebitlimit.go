package events

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/internal/validation"
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
		return errors.New("DailyDebitLimitConsumed must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("DailyDebitLimitConsumed must not have an empty account ID")
	}
	if err := m.DebitType.Validate(); err != nil && !m.DebitType.IsDebit() {
		return errors.New("DailyDebitLimitConsumed must have a valid debit type")
	}
	if m.Amount < 1 {
		return errors.New("DailyDebitLimitConsumed must have a positive amount")
	}
	if !validation.IsValidDate(m.Date) {
		return errors.New("DailyDebitLimitConsumed must have a valid date")
	}
	if m.TotalDebitsForDay < 1 {
		return errors.New("DailyDebitLimitConsumed must have consumed 1 or more total debits for day")
	}
	if m.DailyLimit < 0 {
		return errors.New("DailyDebitLimitConsumed must not have a negative daily limit")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m DailyDebitLimitExceeded) Validate() error {
	if m.TransactionID == "" {
		return errors.New("DailyDebitLimitExceeded must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("DailyDebitLimitExceeded must not have an empty account ID")
	}
	if err := m.DebitType.Validate(); err != nil && !m.DebitType.IsDebit() {
		return errors.New("DailyDebitLimitExceeded must have a valid debit type")
	}
	if m.Amount < 1 {
		return errors.New("DailyDebitLimitExceeded must have a positive amount")
	}
	if !validation.IsValidDate(m.Date) {
		return errors.New("DailyDebitLimitExceeded must have a valid date")
	}
	if m.TotalDebitsForDay < 0 {
		return errors.New("DailyDebitLimitExceeded must not have a negative total debits for day")
	}
	if m.DailyLimit < 0 {
		return errors.New("DailyDebitLimitExceeded must not have a negative daily limit")
	}

	return nil
}
