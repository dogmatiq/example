package events

import "github.com/dogmatiq/example/messages"

// TransferStarted is an event indicating that the process of transferring funds
// from one account to another has begun.
type TransferStarted struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
	ScheduledDate string
}

// TransferApproved is an event that indicates a requested transfer has been
// approved.
type TransferApproved struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
}

// TransferDeclined is an event that indicates a requested transfer has been
// declined.
type TransferDeclined struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
	Reason        messages.DebitFailureReason // TODO: does this name/type make sense for Transfer?
}
