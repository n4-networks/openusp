package main

import (
	"log"

	"github.com/n4-networks/openusp/pkg/cli"
)

func main() {

	cli := &cli.Cli{}
	if err := cli.Init(); err != nil {
		log.Println("Error in initializing shell, exiting...:", err)
		return
	}

	cli.Run()
}
