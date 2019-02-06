package main

import (
	"log"
	"net"

	"github.com/dogmatiq/testkit/engine"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/server"
	"google.golang.org/grpc"
	"github.com/mwitkow/grpc-proxy/proxy"
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

	grpcServer := buildGrpcProxyServer()
	svr := server.NewServer(grpcServer, en).HTTPServer()

	ln, err := net.Listen("tcp", ":9900")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("about to listen for gRPC call on: %v", ln.Addr())
	log.Fatal(svr.Serve(ln))
}


func buildGrpcProxyServer() *grpc.Server {
	// gRPC proxy logic.
	backendConn := dialBackendOrFail()
	director := func(ctx context.Context, fullMethodName string)
						(context.Context, *grpc.ClientConn, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		outCtx, _ := context.WithCancel(ctx)
		mdCopy := md.Copy()
		delete(mdCopy, "user-agent")
		outCtx = metadata.NewOutgoingContext(outCtx, mdCopy)
		return outCtx, backendConn, nil
	}

	return grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()), // needed for proxy to function.
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)),
	)
}
