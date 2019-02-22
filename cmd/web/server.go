package main

import (
	"net/http"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

// gRPC2HTTP returns an instance of HTTP server that is capable of
// conveying requests to gRPC servers over gRPC-Web spec.
func gRPC2HTTP(srv *grpc.Server) *http.Server {
	wrapped := grpcweb.WrapServer(srv)
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
				if wrapped.IsGrpcWebRequest(req) {
					wrapped.ServeHTTP(resp, req)
				} else {
					// otherwise serve the static content
					http.FileServer(http.Dir("www/dist")).ServeHTTP(resp, req)
				}
			}),
	}
}
