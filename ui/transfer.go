package ui

import (
	"fmt"
	"net/http"

	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/ui/templates"
	"github.com/google/uuid"
)

// renderTransferPage renders the transfer form.
func (h *Handler) renderTransferPage(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	accountID := r.PathValue("accountID")

	customerName, err := h.queryCustomerName(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accountName, balance, err := h.queryAccountDetails(r.Context(), accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accountGroups, err := h.queryAllAccountsGrouped(r.Context(), customerID, accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		pageData
		AccountID     string
		AccountName   string
		Balance       money
		AccountGroups []accountGroup
	}{
		pageData: pageData{
			Title:        "Transfer",
			CustomerID:   customerID,
			CustomerName: customerName,
		},
		AccountID:     accountID,
		AccountName:   accountName,
		Balance:       money(balance),
		AccountGroups: accountGroups,
	}

	if err := templates.Get("transfer").ExecuteTemplate(w, "transfer.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// transfer processes a transfer form submission.
func (h *Handler) transfer(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	accountID := r.PathValue("accountID")

	amount, err := parseMoney(r.FormValue("amount"))
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	toAccountID := r.FormValue("to_account_id")
	if toAccountID == "" {
		http.Error(w, "destination account is required", http.StatusBadRequest)
		return
	}

	scheduledTime := parseSchedule(r.FormValue("schedule"))

	err = h.CommandExecutor.ExecuteCommand(
		r.Context(),
		&commands.Transfer{
			TransactionID: uuid.New().String(),
			FromAccountID: accountID,
			ToAccountID:   toAccountID,
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
