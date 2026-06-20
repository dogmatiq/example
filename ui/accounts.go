package ui

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"

	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/ui/templates"
)

// account is a summary of a bank account as displayed in the accounts list.
type account struct {
	ID      string
	Name    string
	Balance money
}

// accountsFragment holds the data needed to render the accounts list table.
// It is used both as a standalone HTMX response and composed into the full page.
type accountsFragment struct {
	CustomerID string
	Accounts   []account
}

// renderAccountsPage renders the full page showing a customer's accounts.
func (h *Handler) renderAccountsPage(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")

	customerName, err := h.queryCustomerName(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accounts, err := h.queryAccounts(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		pageData
		AccountsFragment accountsFragment
	}{
		pageData: pageData{
			Title:        "Your Accounts",
			CustomerID:   customerID,
			CustomerName: customerName,
		},
		AccountsFragment: accountsFragment{
			CustomerID: customerID,
			Accounts:   accounts,
		},
	}

	if err := templates.Get("accounts").ExecuteTemplate(w, "accounts.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderAccountsFragment renders only the accounts table, without the
// surrounding page layout. HTMX polls this endpoint to update the table
// without a full page reload.
func (h *Handler) renderAccountsFragment(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")

	accounts, err := h.queryAccounts(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := accountsFragment{
		CustomerID: customerID,
		Accounts:   accounts,
	}

	if err := templates.Get("accounts").ExecuteTemplate(w, "accounts-fragment", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderOpenAccountPage renders the form for opening a new account.
func (h *Handler) renderOpenAccountPage(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")

	customerName, err := h.queryCustomerName(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := pageData{
		Title:        "Open a New Account",
		CustomerID:   customerID,
		CustomerName: customerName,
	}

	if err := templates.Get("openaccount").ExecuteTemplate(w, "openaccount.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// openAccount handles the form submission to open a new account. It dispatches
// an OpenAccount command and redirects back to the accounts list.
func (h *Handler) openAccount(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	accountName := strings.TrimSpace(r.FormValue("account_name"))

	if accountName == "" {
		http.Error(w, "account name is required", http.StatusBadRequest)
		return
	}

	accountID := generateAccountID()

	err := h.CommandExecutor.ExecuteCommand(
		r.Context(),
		&commands.OpenAccount{
			CustomerID:  customerID,
			AccountID:   accountID,
			AccountName: accountName,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/c/%s/accounts", customerID), http.StatusSeeOther)
}

// queryAccounts loads the accounts belonging to a specific customer.
func (h *Handler) queryAccounts(ctx context.Context, customerID string) ([]account, error) {
	rows, err := h.DB.QueryContext(
		ctx,
		`SELECT
			id,
			name,
			balance
		FROM accounts
		WHERE customer_id = ?
		ORDER BY name`,
		customerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []account
	for rows.Next() {
		var a account

		if err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Balance,
		); err != nil {
			return nil, err
		}

		accounts = append(accounts, a)
	}

	return accounts, rows.Err()
}

// accountGroup is a group of accounts belonging to a single customer, used to
// populate optgroups in the transfer destination picker.
type accountGroup struct {
	CustomerID   string
	CustomerName string
	IsOwn        bool
	Accounts     []account
}

// queryAllAccountsGrouped loads every account in the system grouped by
// customer, excluding excludeAccountID. The group for currentCustomerID is
// labelled "My Accounts" and appears first.
func (h *Handler) queryAllAccountsGrouped(ctx context.Context, currentCustomerID, excludeAccountID string) ([]accountGroup, error) {
	rows, err := h.DB.QueryContext(
		ctx,
		`SELECT
			a.id,
			a.name,
			a.balance,
			a.customer_id,
			c.name
		FROM accounts AS a
		INNER JOIN customers AS c
			ON c.id = a.customer_id
		WHERE a.id != ?
		ORDER BY
			CASE WHEN a.customer_id = ? THEN 0 ELSE 1 END,
			c.name,
			a.name`,
		excludeAccountID,
		currentCustomerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		groups []accountGroup
		g      *accountGroup
	)

	for rows.Next() {
		var (
			a            account
			customerID   string
			customerName string
		)

		if err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Balance,
			&customerID,
			&customerName,
		); err != nil {
			return nil, err
		}

		if g == nil || g.CustomerID != customerID {
			label := customerName
			isOwn := customerID == currentCustomerID
			if isOwn {
				label = "My Accounts"
			}

			groups = append(groups, accountGroup{
				CustomerID:   customerID,
				CustomerName: label,
				IsOwn:        isOwn,
			})
			g = &groups[len(groups)-1]
		}

		g.Accounts = append(g.Accounts, a)
	}

	return groups, rows.Err()
}

// generateAccountID produces a random 9-digit account number.
func generateAccountID() string {
	return strconv.Itoa(rand.IntN(900_000_000) + 100_000_000)
}
