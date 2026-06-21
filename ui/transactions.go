package ui

import (
	"context"
	"net/http"
	"time"

	"github.com/dogmatiq/example/ui/templates"
)

// transaction represents a single entry in an account's transaction history.
type transaction struct {
	OccurredAt  time.Time
	Description string
	Debit       money
	Credit      money
	Balance     money
}

// transactionsFragment holds the data needed to render the transactions table.
// It is used both as a standalone HTMX response and composed into the full page.
type transactionsFragment struct {
	CustomerID   string
	AccountID    string
	Transactions []transaction
}

// renderTransactionsPage renders the full page showing an account's transaction
// history along with deposit, withdraw and transfer actions.
func (h *Handler) renderTransactionsPage(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	accountID := r.PathValue("accountID")

	customerName, err := h.queryCustomerName(r.Context(), customerID)
	if httpError(w, err) {
		return
	}

	accountName, balance, err := h.queryAccountDetails(r.Context(), accountID)
	if httpError(w, err) {
		return
	}

	transactions, err := h.queryTransactions(r.Context(), accountID)
	if httpError(w, err) {
		return
	}

	data := struct {
		pageData
		AccountID            string
		AccountName          string
		Balance              money
		TransactionsFragment transactionsFragment
	}{
		pageData: pageData{
			Title:        "Transactions: " + accountName,
			CustomerID:   customerID,
			CustomerName: customerName,
		},
		AccountID:   accountID,
		AccountName: accountName,
		Balance:     money(balance),
		TransactionsFragment: transactionsFragment{
			CustomerID:   customerID,
			AccountID:    accountID,
			Transactions: transactions,
		},
	}

	if err := templates.Get("transactions").ExecuteTemplate(w, "transactions.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderTransactionsFragment renders only the transactions table, without the
// surrounding page layout. HTMX polls this endpoint to update the table
// without a full page reload.
func (h *Handler) renderTransactionsFragment(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("accountID")

	transactions, err := h.queryTransactions(r.Context(), accountID)
	if httpError(w, err) {
		return
	}

	data := transactionsFragment{
		CustomerID:   r.PathValue("customerID"),
		AccountID:    accountID,
		Transactions: transactions,
	}

	if err := templates.Get("transactions").ExecuteTemplate(w, "transactions-fragment", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// queryTransactions loads the transaction history for an account.
func (h *Handler) queryTransactions(ctx context.Context, accountID string) ([]transaction, error) {
	rows, err := h.DB.QueryContext(
		ctx,
		`SELECT
			created_at,
			description,
			debit,
			credit,
			balance
		FROM ledger
		WHERE account_id = ?
		ORDER BY created_at DESC`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []transaction
	for rows.Next() {
		var t transaction

		if err := rows.Scan(
			&t.OccurredAt,
			&t.Description,
			&t.Debit,
			&t.Credit,
			&t.Balance,
		); err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, rows.Err()
}
