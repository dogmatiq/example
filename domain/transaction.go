package domain

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// transaction is the aggregate root for a bank transaction.
type transaction struct {
	// Started is true if the transaction has started.
	Started bool
}

func (t *transaction) StartDeposit(s dogma.AggregateCommandScope, m commands.Deposit) {
	if t.Started {
		s.Log("transaction already started")
		return
	}

	s.RecordEvent(events.DepositStarted{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func (t *transaction) ApproveDeposit(s dogma.AggregateCommandScope, m commands.ApproveDeposit) {
	s.RecordEvent(events.DepositApproved{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func (t *transaction) StartWithdraw(s dogma.AggregateCommandScope, m commands.Withdraw) {
	if t.Started {
		s.Log("transaction already started")
		return
	}

	s.RecordEvent(events.WithdrawalStarted{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
		ScheduledTime: m.ScheduledTime,
	})
}

func (t *transaction) ApproveWithdrawal(s dogma.AggregateCommandScope, m commands.ApproveWithdrawal) {
	s.RecordEvent(events.WithdrawalApproved{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func (t *transaction) DeclineWithdrawal(s dogma.AggregateCommandScope, m commands.DeclineWithdrawal) {
	s.RecordEvent(events.WithdrawalDeclined{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
		Reason:        m.Reason,
	})
}

func (t *transaction) StartTransfer(s dogma.AggregateCommandScope, m commands.Transfer) {
	if m.FromAccountID == m.ToAccountID {
		s.Log("cannot transfer to same account")
		return
	}

	if t.Started {
		s.Log("transaction already started")
		return
	}

	s.RecordEvent(events.TransferStarted{
		TransactionID: m.TransactionID,
		FromAccountID: m.FromAccountID,
		ToAccountID:   m.ToAccountID,
		Amount:        m.Amount,
		ScheduledTime: m.ScheduledTime,
	})
}

func (t *transaction) ApproveTransfer(s dogma.AggregateCommandScope, m commands.ApproveTransfer) {
	s.RecordEvent(events.TransferApproved{
		TransactionID: m.TransactionID,
		FromAccountID: m.FromAccountID,
		ToAccountID:   m.ToAccountID,
		Amount:        m.Amount,
	})
}

func (t *transaction) DeclineTransfer(s dogma.AggregateCommandScope, m commands.DeclineTransfer) {
	s.RecordEvent(events.TransferDeclined{
		TransactionID: m.TransactionID,
		FromAccountID: m.FromAccountID,
		ToAccountID:   m.ToAccountID,
		Amount:        m.Amount,
		Reason:        m.Reason,
	})
}

func (t *transaction) ApplyEvent(m dogma.Event) {
	switch m.(type) {
	case events.DepositStarted:
		t.Started = true
	case events.WithdrawalStarted:
		t.Started = true
	case events.TransferStarted:
		t.Started = true
	}
}

// TransactionHandler implements the business logic for a transaction of any
// kind against an account.
//
// It's sole purpose is to ensure the global uniqueness of transaction IDs.
type TransactionHandler struct{}

// New returns a new transaction instance.
func (TransactionHandler) New() dogma.AggregateRoot {
	return &transaction{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (TransactionHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("transaction", "2afe7484-8eb4-4c02-9c39-c2493e0defb0")

	c.Routes(
		dogma.HandlesCommand[commands.Deposit](),
		dogma.HandlesCommand[commands.ApproveDeposit](),
		dogma.HandlesCommand[commands.Withdraw](),
		dogma.HandlesCommand[commands.ApproveWithdrawal](),
		dogma.HandlesCommand[commands.DeclineWithdrawal](),
		dogma.HandlesCommand[commands.Transfer](),
		dogma.HandlesCommand[commands.ApproveTransfer](),
		dogma.HandlesCommand[commands.DeclineTransfer](),
		dogma.RecordsEvent[events.DepositStarted](),
		dogma.RecordsEvent[events.DepositApproved](),
		dogma.RecordsEvent[events.WithdrawalStarted](),
		dogma.RecordsEvent[events.WithdrawalApproved](),
		dogma.RecordsEvent[events.WithdrawalDeclined](),
		dogma.RecordsEvent[events.TransferStarted](),
		dogma.RecordsEvent[events.TransferApproved](),
		dogma.RecordsEvent[events.TransferDeclined](),
	)
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (TransactionHandler) RouteCommandToInstance(m dogma.Command) string {
	switch x := m.(type) {
	case commands.Deposit:
		return x.TransactionID
	case commands.ApproveDeposit:
		return x.TransactionID
	case commands.Withdraw:
		return x.TransactionID
	case commands.ApproveWithdrawal:
		return x.TransactionID
	case commands.DeclineWithdrawal:
		return x.TransactionID
	case commands.Transfer:
		return x.TransactionID
	case commands.ApproveTransfer:
		return x.TransactionID
	case commands.DeclineTransfer:
		return x.TransactionID
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this
// handler.
func (TransactionHandler) HandleCommand(
	r dogma.AggregateRoot,
	s dogma.AggregateCommandScope,
	m dogma.Command,
) {
	t := r.(*transaction)

	switch x := m.(type) {
	case commands.Deposit:
		t.StartDeposit(s, x)
	case commands.ApproveDeposit:
		t.ApproveDeposit(s, x)
	case commands.Withdraw:
		t.StartWithdraw(s, x)
	case commands.ApproveWithdrawal:
		t.ApproveWithdrawal(s, x)
	case commands.DeclineWithdrawal:
		t.DeclineWithdrawal(s, x)
	case commands.Transfer:
		t.StartTransfer(s, x)
	case commands.ApproveTransfer:
		t.ApproveTransfer(s, x)
	case commands.DeclineTransfer:
		t.DeclineTransfer(s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}
