package commands

// ChangeCustomerEmailAddress is a command requesting that a customer email
// address be changed.
type ChangeCustomerEmailAddress struct {
	CustomerID    string
	CustomerEmail string
}
