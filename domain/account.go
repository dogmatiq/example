package domain

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// account is the aggregate root for a bank account.
type account struct {
	dogma.NoSnapshotBehavior

	// Opened is true if the account has been opened.
	Opened bool

	// Balance is the current account balance, in cents.
	Balance int64
}

func (a *account) OpenAccount(s dogma.AggregateCommandScope, m *commands.OpenAccount) {
	if a.Opened {
		s.Log("account has already been opened")
		return
	}

	s.RecordEvent(&events.AccountOpened{
		CustomerID:  m.CustomerID,
		AccountID:   m.AccountID,
		AccountName: m.AccountName,
	})
}

func (a *account) CreditAccount(s dogma.AggregateCommandScope, m *commands.CreditAccount) {
	s.RecordEvent(&events.AccountCredited{
		TransactionID:   m.TransactionID,
		AccountID:       m.AccountID,
		TransactionType: m.TransactionType,
		Amount:          m.Amount,
	})
}

func (a *account) DebitAccount(s dogma.AggregateCommandScope, m *commands.DebitAccount) {
	if a.hasSufficientFunds(m.Amount) {
		s.RecordEvent(&events.AccountDebited{
			TransactionID:   m.TransactionID,
			AccountID:       m.AccountID,
			TransactionType: m.TransactionType,
			Amount:          m.Amount,
			ScheduledTime:   m.ScheduledTime,
		})
	} else {
		s.RecordEvent(&events.AccountDebitDeclined{
			TransactionID:   m.TransactionID,
			AccountID:       m.AccountID,
			TransactionType: m.TransactionType,
			Amount:          m.Amount,
			Reason:          messages.InsufficientFunds,
		})
	}
}

func (a *account) hasSufficientFunds(amount int64) bool {
	return a.Balance >= amount
}

func (a *account) ApplyEvent(m dogma.Event) {
	switch x := m.(type) {
	case *events.AccountOpened:
		a.Opened = true
	case *events.AccountCredited:
		a.Balance += x.Amount
	case *events.AccountDebited:
		a.Balance -= x.Amount
	}
}

// AccountHandler implements the business logic for a bank account.
//
// It centralizes all transactions that are applied to an account in order to
// enforce a strict no-overdraw policy.
type AccountHandler struct{}

// New returns a new account instance.
func (AccountHandler) New() dogma.AggregateRoot {
	return &account{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (AccountHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("account", "fcce9a78-23a3-4211-b608-ecbe21ea446f")

	c.Routes(
		dogma.HandlesCommand[*commands.OpenAccount](),
		dogma.HandlesCommand[*commands.CreditAccount](),
		dogma.HandlesCommand[*commands.DebitAccount](),
		dogma.RecordsEvent[*events.AccountOpened](),
		dogma.RecordsEvent[*events.AccountCredited](),
		dogma.RecordsEvent[*events.AccountDebited](),
		dogma.RecordsEvent[*events.AccountDebitDeclined](),
	)
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (AccountHandler) RouteCommandToInstance(m dogma.Command) string {
	switch x := m.(type) {
	case *commands.OpenAccount:
		return x.AccountID
	case *commands.CreditAccount:
		return x.AccountID
	case *commands.DebitAccount:
		return x.AccountID
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this handler.
func (AccountHandler) HandleCommand(
	r dogma.AggregateRoot,
	s dogma.AggregateCommandScope,
	m dogma.Command,
) {
	a := r.(*account)

	switch x := m.(type) {
	case *commands.OpenAccount:
		a.OpenAccount(s, x)
	case *commands.CreditAccount:
		a.CreditAccount(s, x)
	case *commands.DebitAccount:
		a.DebitAccount(s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}
