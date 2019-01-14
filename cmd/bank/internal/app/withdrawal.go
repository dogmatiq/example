package app

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/cmd/bank/internal/messages"
)

// WithdrawalProcessHandler manages the process of withdrawing funds from an account.
var WithdrawalProcessHandler dogma.ProcessMessageHandler = withdrawalProcessHandler{}

type withdrawalProcessHandler struct {
	dogma.NoTimeoutBehavior
}

func (withdrawalProcessHandler) New() dogma.ProcessRoot {
	return nil
}

func (withdrawalProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Name("withdrawal")
	c.RouteEventType(messages.WithdrawalStarted{})
	c.RouteEventType(messages.AccountDebitedForWithdrawal{})
	c.RouteEventType(messages.WithdrawalDeclined{})
}

func (withdrawalProcessHandler) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case messages.WithdrawalStarted:
		return x.TransactionID, true, nil
	case messages.AccountDebitedForWithdrawal:
		return x.TransactionID, true, nil
	case messages.WithdrawalDeclined:
		return x.TransactionID, true, nil
	default:
		return "", false, nil
	}
}

func (withdrawalProcessHandler) HandleEvent(
	_ context.Context,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case messages.WithdrawalStarted:
		s.Begin()
		s.ExecuteCommand(messages.DebitAccountForWithdrawal{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case messages.AccountDebitedForWithdrawal, messages.WithdrawalDeclined:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}