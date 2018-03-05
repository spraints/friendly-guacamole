package main

import (
	"log"

	"github.com/spf13/pflag"
	"google.golang.org/grpc"

	"github.com/spraints/friendly-guacamole/p"
)

func main() {
	log.SetPrefix("[client] ")
	work := 0
	pflag.IntVarP(&work, "work", "s", work, "amount of time to sleep in the job")
	server := "127.0.0.1:55533"
	pflag.StringVarP(&server, "server", "a", server, "server to connect to (default "+server+")")
	pflag.Parse()

	grpc, err := grpc.Dial(server, grpc.WithInsecure())
	perr(err)
	defer grpc.Close()
	client := p.NewExampleClient(grpc)
	perr(err)
	log.Printf("todo: use %v", client)
}

func perr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
