package events

// CustomerAcquired is an event indicating that a new customer has been
// acquired.
type CustomerAcquired struct {
	CustomerID   string
	CustomerName string
	AccountID    string
	AccountName  string
}
