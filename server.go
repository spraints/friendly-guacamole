package main

import (
	"context"
	"log"
	"net"

	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/spraints/friendly-guacamole/p"
	"github.com/spraints/friendly-guacamole/defaults"
)

func main() {
	log.SetPrefix("[server] ")
	var server server
	pflag.IntVarP(&server.timeout, "timeout", "t", 1, "amount of time to allow for each job")
	address := defaults.ServerAddr
	pflag.StringVarP(&address, "address", "a", address, "server to listen on (default "+address+")")
	pflag.Parse()

	grpcServer := grpc.NewServer()
	p.RegisterExampleServer(grpcServer, &server)

	listener, err := net.Listen("tcp", address)
	perr(err)
	defer listener.Close()
	err = grpcServer.Serve(listener)
	perr(err)
}

type server struct {
	timeout int
}

func (s *server) DoSomeWork(ctx context.Context, req *p.WorkRequest) (*p.WorkResponse, error) {
	reqID := "unknown"
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if id, ok := md["request-id"]; ok && len(id) > 0 {
			reqID = id[0]
		}
	}
	log.Printf("client (req %q) says: %#v", reqID, req)
	return &p.WorkResponse{Ack: "Hi! Everything worked!"}, nil
}

func perr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
