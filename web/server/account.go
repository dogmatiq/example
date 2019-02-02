package server

import (
	"context"

	"github.com/dogmatiq/dogmatest/engine"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/web/proto"
)

// accountServer is an unexposed type that implements proto.AccountServer
// interface.
type accountServer struct {
	en *engine.Engine
}

// OpenAccount is a service handler to process OpenAccountRequest.
func (s *accountServer) OpenAccount(
	ctx context.Context,
	req *proto.OpenAccountRequest,
) (*proto.OpenAccountResponse, error) {
	if err := s.en.Dispatch(
		context.Background(),
		messages.OpenAccount{
			AccountID: req.AccountId,
			Name:      req.Name,
		},
		// engine.WithObserver(
		// 	fact.ObserverFunc(func(f fact.Fact) {
		// 		dapper.Print(f)
		// 		fmt.Print("\n\n")
		// 	}),
		// ),
		engine.EnableProjections(true),
	); err != nil {
		return nil, err
	}

	return &proto.OpenAccountResponse{
		AccountId: req.AccountId,
		Name:      req.Name,
	}, nil
}
