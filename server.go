package main

import (
	"log"

	"github.com/spf13/pflag"
)

func main() {
	log.SetPrefix("[server] ")
	var timeout int
	pflag.IntVarP(&timeout, "timeout", "t", 10, "amount of time to allow for each job")
	pflag.Parse()
	log.Printf("todo")
}
