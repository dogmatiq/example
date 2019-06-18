package commands

// OpenAccountForNewCustomer is a command requesting that a new bank account be
// opened for a new customer.
type OpenAccountForNewCustomer struct {
	CustomerID    string
	CustomerName  string
	CustomerEmail string
	AccountID     string
	AccountName   string
}

type ChangeCustomerEmailAddress struct {
	CustomerID    string
	CustomerEmail string
}
