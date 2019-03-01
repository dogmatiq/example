package main

import (
	"net/http"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

// TogRPCWeb returns an instance of HTTP server that is capable of
// conveying requests to gRPC servers over gRPC-Web protocol.
func TogRPCWeb(srv *grpc.Server) *http.Server {
	wrapped := grpcweb.WrapServer(
		srv,
		grpcweb.WithOriginFunc(
			func(origin string) bool {
				return true
			},
		),
	)
	return &http.Server{
		// TO-DO:  replace hard-coded values with options
		ReadTimeout: 0 * time.Second,
		// TO-DO:  replace hard-coded values with options
		WriteTimeout: 0 * time.Second,
		Handler: http.HandlerFunc(
			func(
				resp http.ResponseWriter,
				req *http.Request,
			) {
				wrapped.ServeHTTP(resp, req)
			}),
	}
}
