package events

// AccountOpened is an event indicating that a new bank account has been opened.
type AccountOpened struct {
	AccountID string
	Name      string
}
