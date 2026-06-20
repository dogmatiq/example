package ui

import (
	"context"
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/dogmatiq/dogma"
)

// Handler is an [http.Handler] that serves the Dogmatiq Bank UI.
type Handler struct {
	DB              *sql.DB
	CommandExecutor dogma.CommandExecutor

	once sync.Once
	mux  http.ServeMux
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(func() {
		h.mux.HandleFunc("GET  /", h.renderLoginPage)
		h.mux.HandleFunc("GET  /signup", h.renderSignupPage)
		h.mux.HandleFunc("POST /signup", h.openAccountForNewCustomer)
		h.mux.HandleFunc("GET  /c/{customerID}/accounts", h.renderAccountsPage)
		h.mux.HandleFunc("GET  /c/{customerID}/accounts/fragment", h.renderAccountsFragment)
		h.mux.HandleFunc("GET  /c/{customerID}/accounts/new", h.renderOpenAccountPage)
		h.mux.HandleFunc("POST /c/{customerID}/accounts", h.openAccount)
		h.mux.HandleFunc("GET  /c/{customerID}/accounts/{accountID}/transactions", h.renderTransactionsPage)
		h.mux.HandleFunc("GET  /c/{customerID}/accounts/{accountID}/transactions/fragment", h.renderTransactionsFragment)
		h.mux.HandleFunc("GET  /c/{customerID}/accounts/{accountID}/atm", h.renderATMPage)
		h.mux.HandleFunc("POST /c/{customerID}/accounts/{accountID}/deposit", h.deposit)
		h.mux.HandleFunc("POST /c/{customerID}/accounts/{accountID}/withdraw", h.withdraw)
		h.mux.HandleFunc("GET  /c/{customerID}/accounts/{accountID}/transfer", h.renderTransferPage)
		h.mux.HandleFunc("POST /c/{customerID}/accounts/{accountID}/transfer", h.transfer)
	})

	h.mux.ServeHTTP(w, r)
}

// parseSchedule converts a schedule radio value into a time.
func parseSchedule(s string) time.Time {
	switch s {
	case "10s":
		return time.Now().Add(10 * time.Second)
	case "1m":
		return time.Now().Add(1 * time.Minute)
	default:
		return time.Now()
	}
}

// queryCustomerName returns the name of a customer by ID.
func (h *Handler) queryCustomerName(
	ctx context.Context,
	customerID string,
) (name string, err error) {
	err = h.DB.QueryRowContext(
		ctx,
		`SELECT name
		FROM customers
		WHERE id = ?`,
		customerID,
	).Scan(&name)
	return name, err
}

// queryAccountDetails returns the name and balance for a single account.
func (h *Handler) queryAccountDetails(
	ctx context.Context,
	accountID string,
) (name string, balance int64, err error) {
	err = h.DB.QueryRowContext(
		ctx,
		`SELECT
			name,
			balance
		FROM accounts
		WHERE id = ?`,
		accountID,
	).Scan(
		&name,
		&balance,
	)

	return name, balance, err
}
