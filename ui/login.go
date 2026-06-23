package ui

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/ui/templates"
	"github.com/google/uuid"
)

// customer is an existing bank customer shown on the login page.
type customer struct {
	ID   string
	Name string
}

// renderLoginPage renders the login page, which lists existing customers and
// provides a form to open a new account as a new customer.
func (h *Handler) renderLoginPage(w http.ResponseWriter, r *http.Request) {
	customers, err := h.queryCustomers(r.Context())
	if err != nil {
		renderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		pageData
		Customers []customer
	}{
		pageData:  pageData{Title: "Welcome"},
		Customers: customers,
	}

	if err := templates.Get("login").ExecuteTemplate(w, "login.html", data); err != nil {
		renderError(w, http.StatusInternalServerError, err.Error())
	}
}

// renderSignupPage renders the form for new customers to open their first account.
func (h *Handler) renderSignupPage(w http.ResponseWriter, _ *http.Request) {
	h.renderSignup(w, "")
}

func (h *Handler) renderSignup(w http.ResponseWriter, formError string) {
	data := struct {
		pageData
		Error string
	}{
		pageData: pageData{Title: "Sign Up"},
		Error:    formError,
	}

	if formError != "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	if err := templates.Get("signup").ExecuteTemplate(w, "signup.html", data); err != nil {
		renderError(w, http.StatusInternalServerError, err.Error())
	}
}

// openAccountForNewCustomer handles the form submission to create a new
// customer and open their first account. It dispatches an
// OpenAccountForNewCustomer command and redirects to the accounts list.
func (h *Handler) openAccountForNewCustomer(w http.ResponseWriter, r *http.Request) {
	customerName := strings.TrimSpace(r.FormValue("customer_name"))
	accountName := strings.TrimSpace(r.FormValue("account_name"))

	if customerName == "" || accountName == "" {
		h.renderSignup(w, "Name is required.")
		return
	}

	customerID := uuid.NewString()
	accountID := generateAccountID()

	if err := h.CommandExecutor.ExecuteCommand(
		r.Context(),
		&commands.OpenAccountForNewCustomer{
			CustomerID:   customerID,
			CustomerName: customerName,
			AccountID:    accountID,
			AccountName:  accountName,
		},
	); err != nil {
		renderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/c/%s/accounts", customerID), http.StatusSeeOther)
}

// queryCustomers loads all customers for display on the login page.
func (h *Handler) queryCustomers(ctx context.Context) ([]customer, error) {
	rows, err := h.DB.QueryContext(
		ctx,
		`SELECT
			id,
			name
		FROM customers
		ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []customer
	for rows.Next() {
		var c customer

		if err := rows.Scan(
			&c.ID,
			&c.Name,
		); err != nil {
			return nil, err
		}

		customers = append(customers, c)
	}

	return customers, rows.Err()
}
