package server

import (
	"context"

	"github.com/dogmatiq/example/proto"
	"github.com/dogmatiq/testkit/engine"
)

// customerServer is an unexposed type that implements proto.CustomerServer
// interface.
type customerServer struct {
	en *engine.Engine
}

// Login is a service handler to process LoginRequest.
func (s *customerServer) Login(
	ctx context.Context,
	req *proto.LoginRequest,
) (*proto.LoginResponse, error) {

	return &proto.LoginResponse{
		CustomerName: req.CustomerName,
		CustomerId:   "test-customer-id",
	}, nil
}
