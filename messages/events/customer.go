package events

// CustomerAcquired is an event indicating that a new customer has been
// acquired.
type CustomerAcquired struct {
	CustomerID    string
	CustomerName  string
	CustomerEmail string
	AccountID     string
	AccountName   string
}

// CustomerEmailAddressChanged is an event indicating that a customer has
// changed their email address.
type CustomerEmailAddressChanged struct {
	CustomerID    string
	CustomerEmail string
}
