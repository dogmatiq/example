package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/events"
	"github.com/dogmatiq/projectionkit/sqlprojection"
)

// accountListHandler returns a list of the authenticated customer's bank
// accounts.
type accountListHandler struct {
	DB *sql.DB
}

type accountListResponse struct {
	Accounts []accountListEntry `json:"accounts"`
}

type accountListEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *accountListHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rows, err := h.DB.QueryContext(
		req.Context(),
		`SELECT
			account_id,
			name
		FROM bank.account
		ORDER BY name`,
		// TODO: add customer_id to WHERE clause
	)
	if err != nil {
		writeError(w, err)
		return
	}
	defer rows.Close()

	var response accountListResponse

	for rows.Next() {
		var entry accountListEntry

		if err := rows.Scan(
			&entry.ID,
			&entry.Name,
		); err != nil {
			writeError(w, err)
			return
		}

		response.Accounts = append(response.Accounts, entry)
	}

	if err := rows.Err(); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, response)
}

type AccountListProjectionHandler struct {
	dogma.NoTimeoutHintBehavior
	sqlprojection.NoCompactBehavior
}

func (h *AccountListProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("account-list", "79bd71fe-ff48-41c7-9c49-10e09e2bca85")

	c.ConsumesEventType(events.AccountOpened{})
}

func (h *AccountListProjectionHandler) HandleEvent(
	ctx context.Context,
	tx *sql.Tx,
	s dogma.ProjectionEventScope,
	m dogma.Message,
) error {
	switch m := m.(type) {
	case events.AccountOpened:
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO bank.account (
				customer_id,
				account_id,
				name
			) VALUES (
				$1, $2, $3
			)`,
			m.CustomerID,
			m.AccountID,
			m.AccountName,
		)
		return err
	default:
		panic(dogma.UnexpectedMessage)
	}
}
