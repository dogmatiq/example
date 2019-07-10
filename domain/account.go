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
	case events.AccountCreditedForDeposit:
		r.Balance += x.Amount
	case events.FundsHeldForWithdrawal:
		r.Balance -= x.Amount
	case events.WithdrawalDeclined:
		if x.Reason == messages.ReasonDailyDebitLimitExceeded {
			r.Balance += x.Amount
		}
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
	c.ConsumesCommandType(commands.CreditAccountForDeposit{})
	c.ConsumesCommandType(commands.HoldFundsForWithdrawal{})
	c.ConsumesCommandType(commands.DeclineWithdrawal{})
	c.ConsumesCommandType(commands.SettleWithdrawal{})
	c.ConsumesCommandType(commands.DebitAccountForTransfer{})
	c.ConsumesCommandType(commands.CreditAccountForTransfer{})

	c.ProducesEventType(events.AccountOpened{})
	c.ProducesEventType(events.AccountCreditedForDeposit{})
	c.ProducesEventType(events.FundsHeldForWithdrawal{})
	c.ProducesEventType(events.WithdrawalDeclined{})
	c.ProducesEventType(events.AccountDebitedForWithdrawal{})
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
	case commands.CreditAccountForDeposit:
		return x.AccountID
	case commands.HoldFundsForWithdrawal:
		return x.AccountID
	case commands.SettleWithdrawal:
		return x.AccountID
	case commands.DebitAccountForTransfer:
		return x.AccountID
	case commands.CreditAccountForTransfer:
		return x.AccountID
	case commands.DeclineWithdrawal:
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
	case commands.CreditAccountForDeposit:
		creditForDeposit(s, x)
	case commands.HoldFundsForWithdrawal:
		holdFundsForWithdrawal(s, x)
	case commands.SettleWithdrawal:
		settleWithdrawal(s, x)
	case commands.DeclineWithdrawal:
		declineWithdrawal(s, x)
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

func creditForDeposit(s dogma.AggregateCommandScope, m commands.CreditAccountForDeposit) {
	s.RecordEvent(events.AccountCreditedForDeposit{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func creditForTransfer(s dogma.AggregateCommandScope, m commands.CreditAccountForTransfer) {
	s.RecordEvent(events.AccountCreditedForTransfer{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func holdFundsForWithdrawal(s dogma.AggregateCommandScope, m commands.HoldFundsForWithdrawal) {
	r := s.Root().(*account)

	if r.Balance >= m.Amount {
		s.RecordEvent(events.FundsHeldForWithdrawal{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
			ScheduledDate: m.ScheduledDate,
		})
	} else {
		s.RecordEvent(events.WithdrawalDeclined{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
			Reason:        messages.ReasonInsufficientFunds,
		})
	}
}

func settleWithdrawal(s dogma.AggregateCommandScope, m commands.SettleWithdrawal) {
	s.RecordEvent(events.AccountDebitedForWithdrawal{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func declineWithdrawal(s dogma.AggregateCommandScope, m commands.DeclineWithdrawal) {
	s.RecordEvent(events.WithdrawalDeclined{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
		Reason:        m.Reason,
	})
}

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
