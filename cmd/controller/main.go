package main

import (
	"flag"
	"log"

	"github.com/n4-networks/openusp/pkg/cntlr"
)

func main() {

	var cliMode bool
	var err error

	flag.BoolVar(&cliMode, "c", false, "run with cli")
	flag.Parse()

	var c cntlr.Cntlr
	err = c.Init()
	if err != nil {
		log.Println("Could not initialize Mtp, err:", err)
		return
	}

	var status chan int32
	status, err = c.Run()
	if err != nil {
		log.Println("Error in running server, exiting...")
		return
	}
	if cliMode {
		log.Println("Going to wait here")
		go c.WaitForExit(status)
		c.Cli()
	} else {
		log.Printf("Started MTP server successfully")
		c.WaitForExit(status)
	}
}
