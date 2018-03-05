package main

// This could probably move into server.go, server.go could fork itself once to start this process, and then it could send a command to this process to make it fork and start a GRPC server.

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/spraints/friendly-guacamole/p"
	"github.com/spraints/friendly-guacamole/defaults"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: ./realserver SOCKFILE")
	}
	address := os.Args[1]
	log.SetPrefix("["+address+"] ")

	grpcServer := grpc.NewServer()
	p.RegisterExampleServer(grpcServer, &server{})

	go func() {
		time.Sleep(30 * time.Second)
		panic("I should not still be running")
	}()

	os.Remove(address)
	listener, err := net.Listen("unix", address)
	perr(err)
	defer listener.Close()
	log.Printf("running!")
	err = grpcServer.Serve(listener)
	perr(err)
}

type server struct {
}

func (s *server) DoSomeWork(ctx context.Context, req *p.WorkRequest) (*p.WorkResponse, error) {
	reqID := "unknown"
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if id, ok := md[defaults.RequestIDKey]; ok && len(id) > 0 {
			reqID = id[0]
		}
	}
	log.Printf("client (req %q) says: %#v", reqID, req)
	time.Sleep(time.Duration(req.Amount) * time.Second)
	return &p.WorkResponse{Ack: "Hi! Everything worked!"}, nil
}

func perr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
