package domain

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma-examples/cmd/bank/internal/messages"
)

// DepositProcessHandler manages the process of depositing funds into an account.
var DepositProcessHandler dogma.ProcessMessageHandler = depositProcessHandler{}

type depositProcessHandler struct {
	dogma.StatelessProcess
	dogma.NoTimeouts
}

func (depositProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.RouteEventType(messages.DepositStarted{})
	c.RouteEventType(messages.AccountCreditedForDeposit{})
}

func (depositProcessHandler) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case messages.DepositStarted:
		return x.TransactionID, true, nil
	case messages.AccountCreditedForDeposit:
		return x.TransactionID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func (depositProcessHandler) HandleEvent(_ context.Context, s dogma.ProcessEventScope, m dogma.Message) error {
	switch x := m.(type) {
	case messages.DepositStarted:
		s.Begin()
		s.ExecuteCommand(messages.CreditAccountForDeposit{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case messages.AccountCreditedForDeposit:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
