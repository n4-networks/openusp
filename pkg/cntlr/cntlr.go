package cntlr

import (
	"log"

	"github.com/n4-networks/openusp/pkg/db"
	"github.com/n4-networks/openusp/pkg/pb/cntlrgrpc"
)

const (
	_           = iota
	COAP_SERVER = iota
	GRPC_SERVER
	COAP_SERVER_DTLS
	STOMP_SERVER_TLS
)

type Cntlr struct {
	Cfg    *Cfg
	dbH    db.UspDb
	mtpH   mtpHandler
	cacheH cacheHandler
	agentH agentHandler
	cntlrgrpc.UnimplementedGrpcServer
}

func (c *Cntlr) Init(confFile string) error {

	// Initialize Logger
	log.SetPrefix("OpenUsp: ")
	log.SetFlags(log.Lshortfile)

	var err error

	// Read configuration from yaml file
	err = c.config(confFile)
	if err != nil {
		return err
	}

	// Initialize DB connection
	if err := c.dbInit(); err != nil {
		log.Println("Error in dbInit()")
		return err
	}
	log.Println("Db Init ...successful!")

	if err := c.mtpInit(); err != nil {
		log.Println("Error in mtpInit()")
		return err
	}
	log.Println("MTP Init ...successful!")

	// Initialize Cache handler
	if err := c.cacheInit(); err != nil {
		log.Println("Error in cacheInit()")
		return err
	}

	// Initialize Agent handler
	if err := c.agentHandlerInit(); err != nil {
		return err
	}

	log.Println("Cntlr has been initialized successfully, Version:", getVer())
	return nil
}

func (c *Cntlr) Run() (chan int32, error) {

	var exit chan int32
	var err error = nil

	c.MtpStart()

	go c.GrpcServerThread(c.Cfg.Grpc.Port, exit)

	go c.MtpReceiveThread()

	//go m.agentHandlerThread()

	return exit, err
}

func (c *Cntlr) WaitForExit(exit chan int32) {
	switch <-exit {
	case COAP_SERVER:
		log.Println("CoAP Server has exited")
	case COAP_SERVER_DTLS:
		log.Println("CoAP Secure Server has exited")
	case STOMP_SERVER_TLS:
		log.Println("STOMP Secure Server has exited")
	case GRPC_SERVER:
		log.Println("GRPC Server has exited")
	default:
		log.Println(" MTP server is existing, err:")
	}
}
