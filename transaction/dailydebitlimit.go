package transaction

import (
	"fmt"
	"time"

	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"

	"github.com/dogmatiq/dogma"
)

// dailyDebitLimit is the aggregate root for the daily debit limit policy.
type dailyDebitLimit struct {
	// UsedAmount is the sum of debit amounts used during the period, in cents.
	UsedAmount int64
}

func (d *dailyDebitLimit) ApplyEvent(m dogma.Message) {
	switch x := m.(type) {
	case events.DailyDebitAmountConsumed:
		d.UsedAmount += x.Amount
	case events.DailyDebitAmountRestored:
		d.UsedAmount -= x.Amount
	}
}

// DailyDebitLimitAggregate implements the business logic for a daily debit limit
// policy.
//
// It centralizes all debits that are applied to an account over a calendar day
// in order to enforce a policy of limited daily debits.
type DailyDebitLimitAggregate struct{}

// New returns a new daily debit limit instance.
func (DailyDebitLimitAggregate) New() dogma.AggregateRoot {
	return &dailyDebitLimit{}
}

// Configure configures the behaviour of the engine as it relates to this
// handler.
func (DailyDebitLimitAggregate) Configure(c dogma.AggregateConfigurer) {
	c.Name("dailydebitlimit")

	c.ConsumesCommandType(commands.ConsumeDailyDebitAmount{})
	c.ConsumesCommandType(commands.RestoreDailyDebitAmount{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (DailyDebitLimitAggregate) RouteCommandToInstance(m dogma.Message) string {
	switch x := m.(type) {
	case commands.ConsumeDailyDebitAmount:
		return makeInstanceID(x.TransactionTimestamp, x.AccountID)
	case commands.RestoreDailyDebitAmount:
		return makeInstanceID(x.TransactionTimestamp, x.AccountID)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this handler.
func (DailyDebitLimitAggregate) HandleCommand(s dogma.AggregateCommandScope, m dogma.Message) {
	switch x := m.(type) {
	case commands.ConsumeDailyDebitAmount:
		increaseAmount(s, x)
	case commands.RestoreDailyDebitAmount:
		decreaseAmount(s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func increaseAmount(s dogma.AggregateCommandScope, m commands.ConsumeDailyDebitAmount) {
	s.Create()

	d := s.Root().(*dailyDebitLimit)

	if d.isAmountWithinLimit(m.Amount) {
		s.RecordEvent(events.DailyDebitAmountConsumed{
			TransactionID:        m.TransactionID,
			AccountID:            m.AccountID,
			Amount:               m.Amount,
			TransactionTimestamp: m.TransactionTimestamp,
		})
	} else {
		s.RecordEvent(events.DailyDebitAmountConsumtionRejected{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	}
}

func decreaseAmount(s dogma.AggregateCommandScope, m commands.RestoreDailyDebitAmount) {
	s.RecordEvent(events.DailyDebitAmountRestored{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func makeInstanceID(t time.Time, accountID string) string {
	return fmt.Sprintf("%04d-%02d-%02d:%s", t.Year(), t.Month(), t.Day(), accountID)
}

const maximumDailyDebitLimit = 9000

func (d *dailyDebitLimit) isAmountWithinLimit(amount int64) bool {
	return d.UsedAmount+amount <= maximumDailyDebitLimit
}
