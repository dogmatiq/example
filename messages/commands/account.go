package commands

// OpenAccount is a command requesting that a new bank account be opened.
type OpenAccount struct {
	AccountID string
	Name      string
}
