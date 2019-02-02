package main

import (
	"log"
	"net"

	"github.com/dogmatiq/dogmatest/engine"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/web/server"
	"google.golang.org/grpc"
)

func main() {
	app := &example.App{}

	cfg, err := config.NewApplicationConfig(app)
	if err != nil {
		log.Fatal(err)
	}

	en, err := engine.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	svr := server.NewServer(grpcServer, en).HTTPServer()

	ln, err := net.Listen("tcp", ":9900")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("about to listen for gRPC call on: %v", ln.Addr())
	log.Fatal(svr.Serve(ln))
}
