package domain

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// TransactionHandler implements the business logic for a transaction of any
// kind against an account.
//
// It's sole purpose is to ensure the global uniqueness of transaction IDs.
type TransactionHandler struct {
	dogma.StatelessAggregateBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (TransactionHandler) Configure(c dogma.AggregateConfigurer) {
	c.Name("transaction")

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
func (TransactionHandler) HandleCommand(s dogma.AggregateCommandScope, m dogma.Message) {
	switch x := m.(type) {
	case commands.Deposit:
		startDeposit(s, x)
	case commands.ApproveDeposit:
		approveDeposit(s, x)
	case commands.Withdraw:
		startWithdraw(s, x)
	case commands.ApproveWithdrawal:
		approveWithdrawal(s, x)
	case commands.DeclineWithdrawal:
		declineWithdrawal(s, x)
	case commands.Transfer:
		startTransfer(s, x)
	case commands.ApproveTransfer:
		approveTransfer(s, x)
	case commands.DeclineTransfer:
		declineTransfer(s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func startDeposit(s dogma.AggregateCommandScope, m commands.Deposit) {
	if !s.Create() {
		s.Log("transaction already exists")
		return
	}

	s.RecordEvent(events.DepositStarted{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func approveDeposit(s dogma.AggregateCommandScope, m commands.ApproveDeposit) {
	s.RecordEvent(events.DepositApproved{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func startWithdraw(s dogma.AggregateCommandScope, m commands.Withdraw) {
	if !s.Create() {
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

func approveWithdrawal(s dogma.AggregateCommandScope, m commands.ApproveWithdrawal) {
	s.RecordEvent(events.WithdrawalApproved{
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

func startTransfer(s dogma.AggregateCommandScope, m commands.Transfer) {
	if m.FromAccountID == m.ToAccountID {
		s.Log("cannot transfer to same account")
		return
	}

	if !s.Create() {
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

func approveTransfer(s dogma.AggregateCommandScope, m commands.ApproveTransfer) {
	s.RecordEvent(events.TransferApproved{
		TransactionID: m.TransactionID,
		FromAccountID: m.FromAccountID,
		ToAccountID:   m.ToAccountID,
		Amount:        m.Amount,
	})
}

func declineTransfer(s dogma.AggregateCommandScope, m commands.DeclineTransfer) {
	s.RecordEvent(events.TransferDeclined{
		TransactionID: m.TransactionID,
		FromAccountID: m.FromAccountID,
		ToAccountID:   m.ToAccountID,
		Amount:        m.Amount,
		Reason:        m.Reason,
	})
}
