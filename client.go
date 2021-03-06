package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/spraints/friendly-guacamole/p"
	"github.com/spraints/friendly-guacamole/defaults"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[client] ")
	var work int32
	pflag.Int32VarP(&work, "work", "s", work, "amount of time to sleep in the job")
	server := defaults.ServerAddr
	pflag.StringVarP(&server, "server", "a", server, "server to connect to (default "+server+")")
	pflag.Parse()

	conn, err := grpc.Dial(server, grpc.WithInsecure())
	perr(err)
	defer conn.Close()
	client := p.NewExampleClient(conn)
	perr(err)

	reqID := fmt.Sprintf("req-%d", time.Now().Unix())
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx,
		defaults.RequestIDKey, reqID,
	)
	log.Printf("sending request %q...", reqID)
	res, err := client.DoSomeWork(ctx, &p.WorkRequest{Amount:work})
	log.Printf("res: %#v", res)
	if err != nil {
		log.Printf("err: %s", err.Error())
	}
}

func perr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
