package main

import (
	"flag"
	"log"

	"github.com/n4-networks/openusp/pkg/mtp"
)

func main() {

	var confFile string
	var cliMode bool
	var err error

	flag.StringVar(&confFile, "f", "mtp.yaml", "configuration file of mtp")
	flag.BoolVar(&cliMode, "c", false, "run with cli")
	flag.Parse()

	var m mtp.Mtp
	err = m.Init(confFile)
	if err != nil {
		log.Println("Could not initialize Mtp, err:", err)
		return
	}

	var exit chan int32
	exit, err = m.Server()
	if err != nil {
		log.Println("Error in running server, exiting...")
		return
	}
	if cliMode {
		log.Println("Going to wait here")
		go m.ServerWait(exit)
		m.Cli()
	} else {
		log.Printf("Started MTP server successfully")
		m.ServerWait(exit)
	}
}
