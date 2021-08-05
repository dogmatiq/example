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
	Number string `json:"number"`
	Name   string `json:"name"`
}

func (h *accountListHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	result := accountList{
		Accounts: []account{
			{
				Number: "100",
				Name:   "Savings",
			},
		},
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
