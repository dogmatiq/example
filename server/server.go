package server

import (
	"net/http"
	"time"

	"google.golang.org/grpc"

	"github.com/improbable-eng/grpc-web/go/grpcweb"

	"github.com/dogmatiq/example/proto"
	"github.com/dogmatiq/testkit/engine"
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
					// set the response content type to whatever JS gRPC client
					// original content type was; see issue
					// https://github.com/improbable-eng/grpc-web/issues/162 for
					// details.
					// if ct := req.Header.Get("Content-Type"); ct != "" {
					// 	resp.Header().Set("Content-Type", ct)
					// }
					wrapped.ServeHTTP(resp, req)
				} else {
					// otherwise serve the static content
					http.FileServer(http.Dir("www/dist")).ServeHTTP(resp, req)
				}
			}),
	}
}
