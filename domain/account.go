package domain

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// account is the aggregate root for a bank account.
type account struct {
	// Balance is the current account balance, in cents.
	Balance int64
}

func (r *account) ApplyEvent(m dogma.Message) {
	switch x := m.(type) {
	case events.AccountCredited:
		r.Balance += x.Amount
	case events.AccountDebited:
		r.Balance -= x.Amount

	// TODO: later these below will be merged with generic above

	case events.AccountCreditedForDeposit:
		r.Balance += x.Amount
	case events.AccountDebitedForTransfer:
		r.Balance -= x.Amount
	case events.AccountCreditedForTransfer:
		r.Balance += x.Amount
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
	c.Name("account")

	c.ConsumesCommandType(commands.OpenAccount{})
	c.ConsumesCommandType(commands.CreditAccount{})
	c.ConsumesCommandType(commands.DebitAccount{})

	c.ProducesEventType(events.AccountOpened{})
	c.ProducesEventType(events.AccountCredited{})
	c.ProducesEventType(events.AccountDebited{})
	c.ProducesEventType(events.AccountDebitDeclined{})

	// TODO: later these below will be merged with generic above
	c.ConsumesCommandType(commands.CreditAccountForDeposit{})
	c.ConsumesCommandType(commands.CreditAccountForTransfer{})
	c.ConsumesCommandType(commands.DebitAccountForTransfer{})

	// TODO: later these below will be merged with generic above
	c.ProducesEventType(events.AccountCreditedForDeposit{})
	c.ProducesEventType(events.AccountDebitedForTransfer{})
	c.ProducesEventType(events.TransferDeclinedDueToInsufficientFunds{})
	c.ProducesEventType(events.AccountCreditedForTransfer{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (AccountHandler) RouteCommandToInstance(m dogma.Message) string {
	switch x := m.(type) {
	case commands.OpenAccount:
		return x.AccountID
	case commands.CreditAccount:
		return x.AccountID
	case commands.DebitAccount:
		return x.AccountID

	// TODO: later these will be merged with generic above

	case commands.CreditAccountForDeposit:
		return x.AccountID
	case commands.DebitAccountForTransfer:
		return x.AccountID
	case commands.CreditAccountForTransfer:
		return x.AccountID

	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this handler.
func (AccountHandler) HandleCommand(s dogma.AggregateCommandScope, m dogma.Message) {
	switch x := m.(type) {
	case commands.OpenAccount:
		openAccount(s, x)
	case commands.CreditAccount:
		creditAccount(s, x)
	case commands.DebitAccount:
		debitAccount(s, x)

	// TODO: later these will be merged with generic above

	case commands.CreditAccountForDeposit:
		creditForDeposit(s, x)
	case commands.DebitAccountForTransfer:
		debitForTransfer(s, x)
	case commands.CreditAccountForTransfer:
		creditForTransfer(s, x)

	default:
		panic(dogma.UnexpectedMessage)
	}
}

func openAccount(s dogma.AggregateCommandScope, m commands.OpenAccount) {
	if !s.Create() {
		s.Log("account has already been opened")
		return
	}

	s.RecordEvent(events.AccountOpened{
		CustomerID:  m.CustomerID,
		AccountID:   m.AccountID,
		AccountName: m.AccountName,
	})
}

func creditAccount(s dogma.AggregateCommandScope, m commands.CreditAccount) {
	s.RecordEvent(events.AccountCredited{
		TransactionID:   m.TransactionID,
		AccountID:       m.AccountID,
		TransactionType: m.TransactionType,
		Amount:          m.Amount,
	})
}

func debitAccount(s dogma.AggregateCommandScope, m commands.DebitAccount) {
	r := s.Root().(*account)

	if r.hasAvailableAmount(m.Amount) {
		s.RecordEvent(events.AccountDebited{
			TransactionID:   m.TransactionID,
			AccountID:       m.AccountID,
			TransactionType: m.TransactionType,
			Amount:          m.Amount,
			ScheduledDate:   m.ScheduledDate,
		})
	} else {
		s.RecordEvent(events.AccountDebitDeclined{
			TransactionID:   m.TransactionID,
			AccountID:       m.AccountID,
			TransactionType: m.TransactionType,
			Amount:          m.Amount,
			Reason:          messages.InsufficientFunds,
		})
	}
}

func (r *account) hasAvailableAmount(amount int64) bool {
	return r.Balance-amount >= 0
}

// TODO: later this will be merged with generic above
func creditForDeposit(s dogma.AggregateCommandScope, m commands.CreditAccountForDeposit) {
	s.RecordEvent(events.AccountCreditedForDeposit{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

// TODO: later this will be merged with generic above
func creditForTransfer(s dogma.AggregateCommandScope, m commands.CreditAccountForTransfer) {
	s.RecordEvent(events.AccountCreditedForTransfer{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

// TODO: later this will be merged with generic above
func debitForTransfer(s dogma.AggregateCommandScope, m commands.DebitAccountForTransfer) {
	r := s.Root().(*account)

	if r.Balance >= m.Amount {
		s.RecordEvent(events.AccountDebitedForTransfer{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	} else {
		s.RecordEvent(events.TransferDeclinedDueToInsufficientFunds{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	}
}
