package account

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/message/event"
)

// root is the aggregate root for a bank account.
type root struct {
	// Balance is the current account balance, in cents.
	Balance int64
}

func (r *root) ApplyEvent(m dogma.Message) {
	switch x := m.(type) {
	case event.AccountCreditedForDeposit:
		r.Balance += x.Amount
	case event.AccountCreditedForTransfer:
		r.Balance += x.Amount
	case event.AccountDebitedForWithdrawal:
		r.Balance -= x.Amount
	case event.AccountDebitedForTransfer:
		r.Balance -= x.Amount
	}
}
