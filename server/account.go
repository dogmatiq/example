package server

import (
	"context"
	"fmt"
	"time"

	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/proto"
	"github.com/dogmatiq/testkit/engine"
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
	if err := s.en.ExecuteCommand(
		context.Background(),
		commands.OpenAccount{
			AccountID: req.AccountId,
			Name:      req.Name,
		},
	); err != nil {
		return nil, err
	}

	return &proto.OpenAccountResponse{
		AccountId: req.AccountId,
		Name:      req.Name,
	}, nil
}

func (s *accountServer) TestStreaming(
	req *proto.TestStreamingRequest,
	out proto.Account_TestStreamingServer,
) error {
	now := time.Now()
	for {
		since := time.Since(now).Round(time.Second)
		if err := out.Send(&proto.TestStreamingResponse{
			Message: fmt.Sprintf("The stream has been opened for %v", since),
		}); err != nil {
			return err
		}
		<-time.After(time.Second * 1)
	}
}
