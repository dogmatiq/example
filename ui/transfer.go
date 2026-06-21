package ui

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	"github.com/dogmatiq/example/ui/templates"
	"github.com/google/uuid"
)

// renderTransferPage renders the transfer form.
func (h *Handler) renderTransferPage(w http.ResponseWriter, r *http.Request) {
	h.renderTransfer(w, r, "")
}

func (h *Handler) renderTransfer(w http.ResponseWriter, r *http.Request, formError string) {
	customerID := r.PathValue("customerID")
	accountID := r.PathValue("accountID")

	customerName, err := h.queryCustomerName(r.Context(), customerID)
	if err != nil {
		renderError(w, http.StatusNotFound)
		return
	}

	accountName, balance, err := h.queryAccountDetails(r.Context(), accountID)
	if err != nil {
		renderError(w, http.StatusNotFound)
		return
	}

	accountGroups, err := h.queryAllAccountsGrouped(r.Context(), customerID, accountID)
	if err != nil {
		renderError(w, http.StatusNotFound)
		return
	}

	data := struct {
		pageData
		AccountID     string
		AccountName   string
		Balance       money
		AccountGroups []accountGroup
		Error         string
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
		Error:         formError,
	}

	if formError != "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	if err := templates.Get("transfer").ExecuteTemplate(w, "transfer.html", data); err != nil {
		renderError(w, http.StatusInternalServerError, err.Error())
	}
}

// transfer processes a transfer form submission.
func (h *Handler) transfer(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	accountID := r.PathValue("accountID")

	amount, err := parseMoney(r.FormValue("amount"))
	if err != nil {
		h.renderTransfer(w, r, "Invalid amount.")
		return
	}

	toAccountID := r.FormValue("to_account_id")
	if toAccountID == "" {
		h.renderTransfer(w, r, "Destination account is required.")
		return
	}

	scheduledTime := parseSchedule(r.FormValue("schedule"))

	var formError string

	err = h.CommandExecutor.ExecuteCommand(
		r.Context(),
		&commands.Transfer{
			TransactionID: uuid.New().String(),
			FromAccountID: accountID,
			ToAccountID:   toAccountID,
			Amount:        int64(amount),
			ScheduledTime: scheduledTime,
		},
		dogma.WithEventObserver(func(context.Context, *events.TransferApproved) (bool, error) {
			return true, nil
		}),
		dogma.WithEventObserver(func(_ context.Context, e *events.TransferDeclined) (bool, error) {
			formError = "Transfer declined — " + string(e.Reason) + "."
			return true, nil
		}),
		dogma.WithEventObserver(func(context.Context, *events.TransferFailed) (bool, error) {
			formError = "Transfer failed."
			return true, nil
		}),
	)
	if err != nil && !errors.Is(err, dogma.ErrEventObserverNotSatisfied) {
		renderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if formError != "" {
		h.renderTransfer(w, r, formError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/c/%s/accounts/%s/transactions", customerID, accountID), http.StatusSeeOther)
}
