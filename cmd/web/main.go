package main

import (
	"log"
	"net"
	"os"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/server"
	"github.com/dogmatiq/testkit/engine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"github.com/dogmatiq/example/proto"
)

func main() {
	en, err := engine.New(&example.App{})
	if err != nil {
		panic(err)
	}
	// set a global groc logger
	grpclog.SetLogger(log.New(os.Stderr, "grpc: ", log.LstdFlags))

	svr:= server.NewServer(en)
	grpcSvr :=  grpc.NewServer()

	proto.RegisterAccountServer(grpcSvr, svr)
	httpSvr := gRPC2HTTP(grpcSvr)

	ln, err := net.Listen("tcp", ":9900")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %v", ln.Addr())
	log.Fatal(httpSvr.Serve(ln))
}
