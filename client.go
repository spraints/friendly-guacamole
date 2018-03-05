package main

import (
	"log"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	log.SetPrefix("[client] ")
	var work int
	pflag.IntVarP(&work, "work", "s", 0, "amount of time to sleep in the job")
	pflag.Parse()
	if work == 0 {
		pflag.Usage()
		os.Exit(1)
	}
	log.Printf("todo")
}
