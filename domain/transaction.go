package domain

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// transaction is the aggregate root for a bank transaction.
type transaction struct {
	// Recorded is the recorded state of the transaction.
	Recorded bool
}

func (t *transaction) ApplyEvent(m dogma.Message) {
	switch m.(type) {
	case events.DepositStarted:
		t.Recorded = true
	case events.WithdrawalStarted:
		t.Recorded = true
	case events.TransferStarted:
		t.Recorded = true
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

	c.ConsumesCommandType(commands.Deposit{})
	c.ConsumesCommandType(commands.ApproveDeposit{})
	c.ConsumesCommandType(commands.Withdraw{})
	c.ConsumesCommandType(commands.ApproveWithdrawal{})
	c.ConsumesCommandType(commands.DeclineWithdrawal{})
	c.ConsumesCommandType(commands.Transfer{})
	c.ConsumesCommandType(commands.ApproveTransfer{})
	c.ConsumesCommandType(commands.DeclineTransfer{})

	c.ProducesEventType(events.DepositStarted{})
	c.ProducesEventType(events.DepositApproved{})
	c.ProducesEventType(events.WithdrawalStarted{})
	c.ProducesEventType(events.WithdrawalApproved{})
	c.ProducesEventType(events.WithdrawalDeclined{})
	c.ProducesEventType(events.TransferStarted{})
	c.ProducesEventType(events.TransferApproved{})
	c.ProducesEventType(events.TransferDeclined{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (TransactionHandler) RouteCommandToInstance(m dogma.Message) string {
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
	m dogma.Message,
) {
	t := r.(*transaction)

	switch x := m.(type) {
	case commands.Deposit:
		startDeposit(t, s, x)
	case commands.ApproveDeposit:
		approveDeposit(t, s, x)
	case commands.Withdraw:
		startWithdraw(t, s, x)
	case commands.ApproveWithdrawal:
		approveWithdrawal(t, s, x)
	case commands.DeclineWithdrawal:
		declineWithdrawal(t, s, x)
	case commands.Transfer:
		startTransfer(t, s, x)
	case commands.ApproveTransfer:
		approveTransfer(t, s, x)
	case commands.DeclineTransfer:
		declineTransfer(t, s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func startDeposit(t *transaction, s dogma.AggregateCommandScope, m commands.Deposit) {
	if t.Recorded {
		s.Log("transaction already exists")
		return
	}

	s.RecordEvent(events.DepositStarted{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func approveDeposit(t *transaction, s dogma.AggregateCommandScope, m commands.ApproveDeposit) {
	s.RecordEvent(events.DepositApproved{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func startWithdraw(t *transaction, s dogma.AggregateCommandScope, m commands.Withdraw) {
	if t.Recorded {
		s.Log("transaction already exists")
		return
	}

	s.RecordEvent(events.WithdrawalStarted{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
		ScheduledDate: m.ScheduledDate,
	})
}

func approveWithdrawal(t *transaction, s dogma.AggregateCommandScope, m commands.ApproveWithdrawal) {
	s.RecordEvent(events.WithdrawalApproved{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func declineWithdrawal(t *transaction, s dogma.AggregateCommandScope, m commands.DeclineWithdrawal) {
	s.RecordEvent(events.WithdrawalDeclined{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
		Reason:        m.Reason,
	})
}

func startTransfer(t *transaction, s dogma.AggregateCommandScope, m commands.Transfer) {
	if m.FromAccountID == m.ToAccountID {
		s.Log("cannot transfer to same account")
		return
	}

	if t.Recorded {
		s.Log("transaction already exists")
		return
	}

	s.RecordEvent(events.TransferStarted{
		TransactionID: m.TransactionID,
		FromAccountID: m.FromAccountID,
		ToAccountID:   m.ToAccountID,
		Amount:        m.Amount,
		ScheduledDate: m.ScheduledDate,
	})
}

func approveTransfer(t *transaction, s dogma.AggregateCommandScope, m commands.ApproveTransfer) {
	s.RecordEvent(events.TransferApproved{
		TransactionID: m.TransactionID,
		FromAccountID: m.FromAccountID,
		ToAccountID:   m.ToAccountID,
		Amount:        m.Amount,
	})
}

func declineTransfer(t *transaction, s dogma.AggregateCommandScope, m commands.DeclineTransfer) {
	s.RecordEvent(events.TransferDeclined{
		TransactionID: m.TransactionID,
		FromAccountID: m.FromAccountID,
		ToAccountID:   m.ToAccountID,
		Amount:        m.Amount,
		Reason:        m.Reason,
	})
}
