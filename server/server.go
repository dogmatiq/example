package server

import (
	"github.com/dogmatiq/example/proto"
	"github.com/dogmatiq/testkit/engine"
)

// Server covers all the methods required to expose the API that uses dogmatest
// engine under the hood to
type Server interface {
	proto.AccountServer
	proto.CustomerServer
}

// server is an unexposed type that implements Server interface.
type server struct {
	*accountServer
	*customerServer
}

// NewServer returns a new instance of the object that implements Server.
func NewServer(en *engine.Engine) Server {
	return &server{
		accountServer:  &accountServer{en: en},
		customerServer: &customerServer{en: en},
	}
}
