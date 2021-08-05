package api

import "net/http"

// NewHandler returns an HTTP handler that serves API requests.
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/x/accounts", &accountListHandler{})

	return mux
}
