package api

import (
	"encoding/json"
	"net/http"
)

// accountListHandler returns a list of the authenticated customer's bank
// accounts.
type accountListHandler struct {
}

type accountList struct {
	Accounts []account `json:"accounts"`
}

type account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *accountListHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	writeJSON(
		w,
		accountList{
			Accounts: []account{
				{
					ID:   "ddbc4088-f249-40fe-aa92-72dcef7cacd2",
					Name: "Savings",
				},
				{
					ID:   "fce1748b-9d69-4bc6-abe9-5ffe6c378c25",
					Name: "Chequing",
				},
			},
		},
	)
}

// writeJSON writes the JSON representation of v to w.
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
