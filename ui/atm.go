package ui

import (
	"fmt"
	"net/http"

	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/ui/templates"
	"github.com/google/uuid"
)

// renderATMPage renders the virtual ATM page with deposit and withdrawal forms.
func (h *Handler) renderATMPage(w http.ResponseWriter, r *http.Request) {
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

	data := struct {
		pageData
		AccountID   string
		AccountName string
		Balance     money
	}{
		pageData: pageData{
			Title:        "Virtual ATM: " + accountName,
			CustomerID:   customerID,
			CustomerName: customerName,
		},
		AccountID:   accountID,
		AccountName: accountName,
		Balance:     money(balance),
	}

	if err := templates.Get("atm").ExecuteTemplate(w, "atm.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// deposit processes a deposit form submission.
func (h *Handler) deposit(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	accountID := r.PathValue("accountID")

	amount, err := parseMoney(r.FormValue("amount"))
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	err = h.CommandExecutor.ExecuteCommand(
		r.Context(),
		&commands.Deposit{
			TransactionID: uuid.New().String(),
			AccountID:     accountID,
			Amount:        int64(amount),
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/c/%s/accounts/%s/transactions", customerID, accountID), http.StatusSeeOther)
}

// withdraw processes a withdrawal form submission.
func (h *Handler) withdraw(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	accountID := r.PathValue("accountID")

	amount, err := parseMoney(r.FormValue("amount"))
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	scheduledTime := parseSchedule(r.FormValue("schedule"))

	err = h.CommandExecutor.ExecuteCommand(
		r.Context(),
		&commands.Withdraw{
			TransactionID: uuid.New().String(),
			AccountID:     accountID,
			Amount:        int64(amount),
			ScheduledTime: scheduledTime,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/c/%s/accounts/%s/transactions", customerID, accountID), http.StatusSeeOther)
}
