package projections

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/cmd/bank/internal/messages"
)

// AccountProjectionHandler is a projection that builds a report of accounts and
// their balances.
type AccountProjectionHandler struct {
	m       sync.RWMutex
	records map[string]*record
}

// record is an entry maintained by the account projection.
type record struct {
	AccountID      string
	Name           string
	DepositsIn     int64
	WithdrawalsOut int64
	TransfersIn    int64
	TransfersOut   int64
	CurrentBalance int64
}

// Configure configs the engine for this projection.
func (h *AccountProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Name("account-projection")
	c.RouteEventType(messages.AccountOpened{})
	c.RouteEventType(messages.AccountCreditedForDeposit{})
	c.RouteEventType(messages.AccountDebitedForWithdrawal{})
	c.RouteEventType(messages.AccountCreditedForTransfer{})
	c.RouteEventType(messages.AccountDebitedForTransfer{})
}

// GenerateCSV writes a CSV report of all accounts to w.
// The CSV is sorted by the current account balance in descending order.
func (h *AccountProjectionHandler) GenerateCSV(w io.Writer) error {
	records := h.slice()

	sort.Slice(records, func(i, j int) bool {
		return records[i].CurrentBalance > records[j].CurrentBalance
	})

	cw := csv.NewWriter(w)

	if err := cw.Write([]string{
		"Account ID",
		"Name",
		"Deposits In",
		"Withdrawals Out",
		"Transfers In",
		"Transfers Out",
		"Current Balance",
	}); err != nil {
		return err
	}

	for _, r := range records {
		if err := cw.Write([]string{
			r.AccountID,
			r.Name,
			formatAmount(r.DepositsIn),
			formatAmount(r.WithdrawalsOut),
			formatAmount(r.TransfersIn),
			formatAmount(r.TransfersOut),
			formatAmount(r.CurrentBalance),
		}); err != nil {
			return err
		}
	}

	cw.Flush()

	return cw.Error()
}

// slice returns a slice of the in-memory records.
func (h *AccountProjectionHandler) slice() []record {
	h.m.RLock()
	defer h.m.RUnlock()

	var records []record
	for _, r := range h.records {
		records = append(records, *r)
	}

	return records
}

// HandleEvent updates the in-memory records to reflect the occurence of m.
func (h *AccountProjectionHandler) HandleEvent(
	_ context.Context,
	s dogma.ProjectionEventScope,
	m dogma.Message,
) error {
	h.m.Lock()
	defer h.m.Unlock()

	switch x := m.(type) {
	case messages.AccountOpened:
		r := h.get(x.AccountID)
		r.Name = x.Name

	case messages.AccountCreditedForDeposit:
		r := h.get(x.AccountID)
		r.DepositsIn += x.Amount
		r.CurrentBalance += x.Amount

	case messages.AccountDebitedForWithdrawal:
		r := h.get(x.AccountID)
		r.WithdrawalsOut += x.Amount
		r.CurrentBalance -= x.Amount

	case messages.AccountCreditedForTransfer:
		r := h.get(x.AccountID)
		r.TransfersIn += x.Amount
		r.CurrentBalance += x.Amount

	case messages.AccountDebitedForTransfer:
		r := h.get(x.AccountID)
		r.TransfersOut += x.Amount
		r.CurrentBalance -= x.Amount

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}

// get returns the record with the given ID, creating it if it does not exist.
func (h *AccountProjectionHandler) get(id string) *record {
	if r, ok := h.records[id]; ok {
		return r
	}

	if h.records == nil {
		h.records = map[string]*record{}
	}

	r := &record{
		AccountID: id,
	}

	h.records[id] = r

	return r
}

// formatAmount formats an amount in cents as dollars.
func formatAmount(amount int64) string {
	s := fmt.Sprintf("%03d", amount)
	l := len(s)
	return s[:l-2] + "." + s[l-2:]
}
