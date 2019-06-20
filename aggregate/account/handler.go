package account

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/message/command"
	"github.com/dogmatiq/example/message/event"
)

// AggregateHandler implements the business logic for a bank account.
//
// It centralizes all transactions that are applied to an account in order to
// enforce a strict no-overdraw policy.
type AggregateHandler struct{}

// New returns a new account instance.
func (AggregateHandler) New() dogma.AggregateRoot {
	return &root{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (AggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Name("account")

	c.ConsumesCommandType(command.OpenAccount{})
	c.ConsumesCommandType(command.CreditAccountForDeposit{})
	c.ConsumesCommandType(command.CreditAccountForTransfer{})
	c.ConsumesCommandType(command.DebitAccountForWithdrawal{})
	c.ConsumesCommandType(command.DebitAccountForTransfer{})

	c.ProducesEventType(event.AccountOpened{})
	c.ProducesEventType(event.AccountCreditedForDeposit{})
	c.ProducesEventType(event.AccountCreditedForTransfer{})
	c.ProducesEventType(event.AccountDebitedForWithdrawal{})
	c.ProducesEventType(event.WithdrawalDeclinedDueToInsufficientFunds{})
	c.ProducesEventType(event.AccountDebitedForTransfer{})
	c.ProducesEventType(event.TransferDeclinedDueToInsufficientFunds{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (AggregateHandler) RouteCommandToInstance(m dogma.Message) string {
	switch x := m.(type) {
	case command.OpenAccount:
		return x.AccountID
	case command.CreditAccountForDeposit:
		return x.AccountID
	case command.CreditAccountForTransfer:
		return x.AccountID
	case command.DebitAccountForWithdrawal:
		return x.AccountID
	case command.DebitAccountForTransfer:
		return x.AccountID
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this handler.
func (AggregateHandler) HandleCommand(s dogma.AggregateCommandScope, m dogma.Message) {
	switch x := m.(type) {
	case command.OpenAccount:
		openAccount(s, x)
	case command.CreditAccountForDeposit:
		creditForDeposit(s, x)
	case command.CreditAccountForTransfer:
		creditForTransfer(s, x)
	case command.DebitAccountForWithdrawal:
		debitForWithdrawal(s, x)
	case command.DebitAccountForTransfer:
		debitForTransfer(s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func openAccount(s dogma.AggregateCommandScope, m command.OpenAccount) {
	if !s.Create() {
		s.Log("account has already been opened")
		return
	}

	s.RecordEvent(event.AccountOpened{
		CustomerID:  m.CustomerID,
		AccountID:   m.AccountID,
		AccountName: m.AccountName,
	})
}

func creditForDeposit(s dogma.AggregateCommandScope, m command.CreditAccountForDeposit) {
	s.RecordEvent(event.AccountCreditedForDeposit{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func creditForTransfer(s dogma.AggregateCommandScope, m command.CreditAccountForTransfer) {
	s.RecordEvent(event.AccountCreditedForTransfer{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func debitForWithdrawal(s dogma.AggregateCommandScope, m command.DebitAccountForWithdrawal) {
	r := s.Root().(*root)

	if r.Balance >= m.Amount {
		s.RecordEvent(event.AccountDebitedForWithdrawal{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	} else {
		s.RecordEvent(event.WithdrawalDeclinedDueToInsufficientFunds{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	}
}

func debitForTransfer(s dogma.AggregateCommandScope, m command.DebitAccountForTransfer) {
	r := s.Root().(*root)

	if r.Balance >= m.Amount {
		s.RecordEvent(event.AccountDebitedForTransfer{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	} else {
		s.RecordEvent(event.TransferDeclinedDueToInsufficientFunds{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	}
}
