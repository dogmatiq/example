package domain

import (
	"fmt"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// dailyDebitLimit is the aggregate root for an account daily debit limit
// policy.
type dailyDebitLimit struct {
	// UsedAmount is the sum of debit amounts used during the period, in cents.
	UsedAmount int64
}

func (r *dailyDebitLimit) ApplyEvent(m dogma.Message) {
	switch x := m.(type) {
	case events.DailyDebitLimitConsumed:
		r.UsedAmount = x.Amount
	}
}

// DailyDebitLimitHandler implements the business logic for an account daily
// debit limit policy.
//
// It centralizes all debits that are applied to an account over a calendar day
// in order to enforce a policy of limited daily debits.
type DailyDebitLimitHandler struct{}

// New returns a new daily debit limit instance.
func (DailyDebitLimitHandler) New() dogma.AggregateRoot {
	return &dailyDebitLimit{}
}

// Configure configures the behaviour of the engine as it relates to this
// handler.
func (DailyDebitLimitHandler) Configure(c dogma.AggregateConfigurer) {
	c.Name("daily-debit-limit")

	c.ConsumesCommandType(commands.ConsumeDailyDebitLimit{})

	c.ProducesEventType(events.DailyDebitLimitConsumed{})
	c.ProducesEventType(events.DailyDebitLimitExceeded{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (DailyDebitLimitHandler) RouteCommandToInstance(m dogma.Message) string {
	switch x := m.(type) {
	case commands.ConsumeDailyDebitLimit:
		return makeInstanceID(x.ScheduledDate, x.AccountID)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this handler.
func (DailyDebitLimitHandler) HandleCommand(s dogma.AggregateCommandScope, m dogma.Message) {
	switch x := m.(type) {
	case commands.ConsumeDailyDebitLimit:
		consume(s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func consume(s dogma.AggregateCommandScope, m commands.ConsumeDailyDebitLimit) {
	s.Create()

	r := s.Root().(*dailyDebitLimit)

	if r.isAmountWithinLimit(m.Amount) {
		s.RecordEvent(events.DailyDebitLimitConsumed{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
			LimitUsed:     r.UsedAmount + m.Amount,
			LimitMaximum:  maximumDailyDebitLimit,
		})
	} else {
		s.RecordEvent(events.DailyDebitLimitExceeded{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
			LimitUsed:     r.UsedAmount,
			LimitMaximum:  maximumDailyDebitLimit,
		})
	}
}

func (r *dailyDebitLimit) isAmountWithinLimit(amount int64) bool {
	return r.UsedAmount+amount <= maximumDailyDebitLimit
}

func makeInstanceID(t time.Time, accountID string) string {
	return fmt.Sprintf(
		"%04d-%02d-%02d:%s",
		t.Year(),
		t.Month(),
		t.Day(),
		accountID,
	)
}

// Limit amount, in cents.
const maximumDailyDebitLimit = 900000
