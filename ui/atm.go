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

// renderATMPage renders the virtual ATM page with deposit and withdrawal forms.
func (h *Handler) renderATMPage(w http.ResponseWriter, r *http.Request) {
	h.renderATM(w, r, "")
}

func (h *Handler) renderATM(w http.ResponseWriter, r *http.Request, formError string) {
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

	data := struct {
		pageData
		AccountID   string
		AccountName string
		Balance     money
		Error       string
	}{
		pageData: pageData{
			Title:        "Virtual ATM: " + accountName,
			CustomerID:   customerID,
			CustomerName: customerName,
		},
		AccountID:   accountID,
		AccountName: accountName,
		Balance:     money(balance),
		Error:       formError,
	}

	if formError != "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	if err := templates.Get("atm").ExecuteTemplate(w, "atm.html", data); err != nil {
		renderError(w, http.StatusInternalServerError, err.Error())
	}
}

// deposit processes a deposit form submission.
func (h *Handler) deposit(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	accountID := r.PathValue("accountID")

	amount, err := parseMoney(r.FormValue("amount"))
	if err != nil {
		h.renderATM(w, r, "Invalid amount.")
		return
	}

	err = h.CommandExecutor.ExecuteCommand(
		r.Context(),
		&commands.Deposit{
			TransactionID: uuid.New().String(),
			AccountID:     accountID,
			Amount:        int64(amount),
		},
		dogma.WithEventObserver(func(context.Context, *events.DepositApproved) (bool, error) {
			return true, nil
		}),
	)
	if err != nil && !errors.Is(err, dogma.ErrEventObserverNotSatisfied) {
		renderError(w, http.StatusInternalServerError, err.Error())
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
		h.renderATM(w, r, "Invalid amount.")
		return
	}

	scheduledTime := parseSchedule(r.FormValue("schedule"))

	var formError string

	err = h.CommandExecutor.ExecuteCommand(
		r.Context(),
		&commands.Withdraw{
			TransactionID: uuid.New().String(),
			AccountID:     accountID,
			Amount:        int64(amount),
			ScheduledTime: scheduledTime,
		},
		dogma.WithEventObserver(func(context.Context, *events.WithdrawalApproved) (bool, error) {
			return true, nil
		}),
		dogma.WithEventObserver(func(_ context.Context, e *events.WithdrawalDeclined) (bool, error) {
			formError = "Withdrawal declined — " + string(e.Reason) + "."
			return true, nil
		}),
	)
	if err != nil && !errors.Is(err, dogma.ErrEventObserverNotSatisfied) {
		renderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if formError != "" {
		h.renderATM(w, r, formError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/c/%s/accounts/%s/transactions", customerID, accountID), http.StatusSeeOther)
}
