package server

import (
	"net/http"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	"github.com/dogmatiq/dogmatest/engine"
	"github.com/dogmatiq/example/web/proto"
)

// Server covers all the methods required to expose the API that uses dogmatest
// engine under the hood to
type Server interface {
	proto.AccountServer
	// HTTPServer returns an instance of HTTP server that is capable of
	// conveying requests to gRPC servers over gRPC-Web spec.
	HTTPServer() *http.Server
}

// server is an unexposed type that implements Server interface.
type server struct {
	*accountServer
	grpcSvr *grpc.Server
}

// NewServer returns a new instance of the object that implements Server.
func NewServer(srv *grpc.Server, en *engine.Engine) Server {
	s := &server{
		accountServer: &accountServer{en: en},
	}

	// register all gRPC servers below
	proto.RegisterAccountServer(srv, s)

	s.grpcSvr = srv
	return s
}

// HTTPServer returns an instance of HTTP server that is capable of
// conveying requests to gRPC servers over gRPC-Web spec.
func (s *server) HTTPServer( /* TO-DO: consider options here */ ) *http.Server {
	wrapped := grpcweb.WrapServer(s.grpcSvr)
	return &http.Server{
		WriteTimeout: 10 * time.Second, // TO-DO:  replace hard-coded values with options
		ReadTimeout:  10 * time.Second, // TO-DO:  replace hard-coded values with options
		Handler: http.HandlerFunc(
			func(
				resp http.ResponseWriter,
				req *http.Request,
			) {
				wrapped.ServeHTTP(resp, req)
			}),
	}
}
