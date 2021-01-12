package domain

import (
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// Limit amount, in cents.
const maximumDailyDebitLimit = 900000

// dailyDebitLimit is the aggregate root for an account daily debit limit
// policy.
type dailyDebitLimit struct {
	// TotalDebitsForDay is the total of all debits for the day, in cents.
	TotalDebitsForDay int64
}

func (d *dailyDebitLimit) Consume(s dogma.AggregateCommandScope, m commands.ConsumeDailyDebitLimit) {
	if d.wouldExceedLimit(m.Amount) {
		s.RecordEvent(events.DailyDebitLimitExceeded{
			TransactionID:     m.TransactionID,
			AccountID:         m.AccountID,
			DebitType:         m.DebitType,
			Amount:            m.Amount,
			Date:              m.ScheduledDate,
			TotalDebitsForDay: d.TotalDebitsForDay,
			DailyLimit:        maximumDailyDebitLimit,
		})
	} else {
		s.RecordEvent(events.DailyDebitLimitConsumed{
			TransactionID:     m.TransactionID,
			AccountID:         m.AccountID,
			DebitType:         m.DebitType,
			Amount:            m.Amount,
			Date:              m.ScheduledDate,
			TotalDebitsForDay: d.TotalDebitsForDay + m.Amount,
			DailyLimit:        maximumDailyDebitLimit,
		})
	}
}

func (d *dailyDebitLimit) wouldExceedLimit(amount int64) bool {
	return d.TotalDebitsForDay+amount > maximumDailyDebitLimit
}

func (d *dailyDebitLimit) ApplyEvent(m dogma.Message) {
	switch x := m.(type) {
	case events.DailyDebitLimitConsumed:
		d.TotalDebitsForDay = x.Amount
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
	c.Identity("daily-debit-limit", "238c5a7b-b51d-42d8-ac8d-a8c81b780230")

	c.ConsumesCommandType(commands.ConsumeDailyDebitLimit{})

	c.ProducesEventType(events.DailyDebitLimitConsumed{})
	c.ProducesEventType(events.DailyDebitLimitExceeded{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (DailyDebitLimitHandler) RouteCommandToInstance(m dogma.Message) string {
	switch x := m.(type) {
	case commands.ConsumeDailyDebitLimit:
		return fmt.Sprintf("%s:%s", messages.FormatDate(x.ScheduledDate), x.AccountID)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this handler.
func (DailyDebitLimitHandler) HandleCommand(
	r dogma.AggregateRoot,
	s dogma.AggregateCommandScope,
	m dogma.Message,
) {
	d := r.(*dailyDebitLimit)

	switch x := m.(type) {
	case commands.ConsumeDailyDebitLimit:
		d.Consume(s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}
