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
)

func main() {

	en, err := engine.New(&example.App{})
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	grpclog.SetLogger(log.New(os.Stderr, "grpc: ", log.LstdFlags))
	svr := server.NewServer(grpcServer, en).HTTPServer()

	ln, err := net.Listen("tcp", ":9900")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("about to listen for gRPC call on: %v", ln.Addr())
	log.Fatal(svr.Serve(ln))
}
