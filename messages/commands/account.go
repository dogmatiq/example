package commands

// OpenAccount is a command requesting that a new bank account be opened for an
// existing customer.
type OpenAccount struct {
	CustomerID  string
	AccountID   string
	AccountName string
}
