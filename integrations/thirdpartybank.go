package integrations

import (
	"context"
	"strconv"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// ThirdPartyBankIntegrationHandler handles commands that interact with
// a hypothetical third-party bank's API on behalf of the application.
type ThirdPartyBankIntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this handler.
func (ThirdPartyBankIntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("third-party-bank", "f2a7e4b1-9c3d-4f8a-b6e5-1d0c2a9f7e3b")

	c.Routes(
		dogma.HandlesCommand[*commands.CreditThirdPartyAccount](),
		dogma.RecordsEvent[*events.ThirdPartyAccountCredited](),
		dogma.RecordsEvent[*events.ThirdPartyAccountCreditFailed](),
	)
}

// HandleCommand handles a command message that has been routed to this handler.
func (ThirdPartyBankIntegrationHandler) HandleCommand(
	_ context.Context,
	s dogma.IntegrationCommandScope,
	c dogma.Command,
) error {
	switch x := c.(type) {
	case *commands.CreditThirdPartyAccount:
		s.Log(
			"crediting third-party account %s with %s (transaction %s)",
			x.AccountID,
			messages.FormatAmount(x.Amount),
			x.TransactionID,
		)

		if _, err := strconv.ParseUint(x.AccountID, 10, 64); err != nil {
			s.Log("third-party bank rejected the credit: account %s not found", x.AccountID)
			s.RecordEvent(&events.ThirdPartyAccountCreditFailed{
				TransactionID: x.TransactionID,
				AccountID:     x.AccountID,
				Amount:        x.Amount,
			})
		} else {
			s.RecordEvent(&events.ThirdPartyAccountCredited{
				TransactionID: x.TransactionID,
				AccountID:     x.AccountID,
				Amount:        x.Amount,
			})
		}

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
