package debitpolicy

import (
	"fmt"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// debitpolicy is the aggregate root for the debit policy.
type debitpolicy struct {
	// DebitAmount is the current sum of debit amounts for the period, in cents.
	DebitAmount int64
}

func (p *debitpolicy) ApplyEvent(m dogma.Message) {
	switch x := m.(type) {
	case events.AccountDebitedForWithdrawal:
		p.DebitAmount += x.Amount
	case events.AccountDebitedForTransfer:
		p.DebitAmount += x.Amount
	}
}

// Aggregate implements the business logic for a debit policy.
//
// It centralizes all debits that are applied to an account over a period of
// time in order to enforce a limited debit policy.
type Aggregate struct{}

// New returns a new debit policy instance.
func (Aggregate) New() dogma.AggregateRoot {
	return &debitpolicy{}
}

// Configure configures the behaviour of the engine as it relates to this
// handlers.
func (Aggregate) Configure(c dogma.AggregateConfigurer) {
	c.Name("debitpolicy")
	c.RouteCommandType(commands.CheckWithdrawalAllowedByDebitPolicy{})
	c.RouteCommandType(commands.CheckTransferAllowedByDebitPolicy{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (Aggregate) RouteCommandToInstance(m dogma.Message) string {
	switch x := m.(type) {
	case commands.CheckWithdrawalAllowedByDebitPolicy:
		return makeInstanceID(x.Timestamp, x.AccountID)
	case commands.CheckTransferAllowedByDebitPolicy:
		return makeInstanceID(x.Timestamp, x.AccountID)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has bene routed to this handler.
func (Aggregate) HandleCommand(s dogma.AggregateCommandScope, m dogma.Message) {
	switch x := m.(type) {
	case commands.CheckWithdrawalAllowedByDebitPolicy:
		checkWithdrawalAllowed(s, x)
	case commands.CheckTransferAllowedByDebitPolicy:
		checkTransferAllowed(s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func checkWithdrawalAllowed(s dogma.AggregateCommandScope, m commands.CheckWithdrawalAllowedByDebitPolicy) {
	s.Create()

	p := s.Root().(*debitpolicy)

	if p.DebitAmount+m.Amount <= debitLimitPerPeriod {
		s.RecordEvent(events.WithdrawalApprovedByDebitPolicy{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	} else {
		s.RecordEvent(events.WithdrawalDeclinedByDebitPolicy{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	}
}

func checkTransferAllowed(s dogma.AggregateCommandScope, m commands.CheckTransferAllowedByDebitPolicy) {
	s.Create()

	p := s.Root().(*debitpolicy)

	if p.DebitAmount+m.Amount <= debitLimitPerPeriod {
		s.RecordEvent(events.TransferApprovedByDebitPolicy{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	} else {
		s.RecordEvent(events.TransferDeclinedByDebitPolicy{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	}
}

func makeInstanceID(t time.Time, accountID string) string {
	return fmt.Sprintf("%04d%02d%02d-%s", t.Year(), t.Month(), t.Day(), accountID)
}

const debitLimitPerPeriod = 9000
