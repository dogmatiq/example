package command

// OpenAccountForNewCustomer is a command requesting that a new bank account be
// opened for a new customer.
type OpenAccountForNewCustomer struct {
	CustomerID    string
	CustomerName  string
	CustomerEmail string
	AccountID     string
	AccountName   string
}

// OpenAccount is a command requesting that a new bank account be opened for an
// existing customer.
type OpenAccount struct {
	CustomerID  string
	AccountID   string
	AccountName string
}
