package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"github.com/mwitkow/grpc-proxy/proxy"

	"github.com/spraints/friendly-guacamole/defaults"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[server] ")
	var server server
	pflag.IntVarP(&server.timeout, "timeout", "t", 1, "amount of time to allow for each job")
	address := defaults.ServerAddr
	pflag.StringVarP(&address, "address", "a", address, "server to listen on (default "+address+")")
	pflag.Parse()

	grpcServer := grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(server.proxy)),
	)

	listener, err := net.Listen("tcp", address)
	perr(err)
	defer listener.Close()
	err = grpcServer.Serve(listener)
	perr(err)
}

type server struct {
	timeout int
}

func (s *server) proxy(ctx context.Context, method string) (context.Context, *grpc.ClientConn, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, nil, fmt.Errorf("no metadata in request")
	}
	ids, ok := md[defaults.RequestIDKey]
	if !ok || len(ids) < 1 {
		return nil, nil, fmt.Errorf("no request-id in request")
	}
	reqID := ids[0]
	log.Printf("proxy %s to %s", reqID, method)
	ctx, cancelDl := context.WithTimeout(ctx, time.Duration(2 + s.timeout) * time.Second)
	sock := reqID + ".sock" // yolo
	err := runRealServerWithTimeout(ctx, cancelDl, sock)
	if err != nil {
		log.Printf("Error starting realserver: %s", err.Error())
		cancelDl()
		return ctx, nil, err
	}
	conn, err := grpc.DialContext(ctx, sock,
		grpc.WithCodec(proxy.Codec()),
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		log.Printf("Error dialing realserver: %s", err.Error())
		cancelDl()
		return ctx, nil, err
	}
	log.Printf("returning connection to proxy!")
	return ctx, conn, nil
}

func runRealServerWithTimeout(ctx context.Context, cancel context.CancelFunc, sock string) error {
	cmd := exec.Command("./realserver", sock)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	log.Printf("start realserver %s", sock)
	err := cmd.Start()
	if err != nil {
		return err
	}
	done := make(chan struct{}, 1)
	go func(done chan<- struct{}) {
		cmd.Wait()
		close(done)
	}(done)
	go func(done <-chan struct{}) {
		select {
		case v, ok := <-done:
			log.Printf("realserver finished (%#v, %s)", v, ok)
			cancel()
		case v, ok := <-ctx.Done():
			log.Printf("timed out or request finished (%#v, %s), killing realserver %d", v, ok, cmd.Process.Pid)
			cmd.Process.Kill()
		}
	}(done)
	time.Sleep(10 * time.Millisecond)
	return nil
}

func perr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
