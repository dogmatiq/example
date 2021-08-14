package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// NewHandler returns an HTTP handler that serves API requests.
func NewHandler(db *sql.DB, sub Subscribable) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/x/accounts", &accountListHandler{DB: db})
	mux.Handle("/x/accounts.sse", &accountListSSEHandler{Subscribable: sub})

	return mux
}

// writeJSON writes the JSON representation of v to w.
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// writeError writes an internal server error response and print information
// about the causal error to stdout.
func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error!"))

	fmt.Println(err)
}
