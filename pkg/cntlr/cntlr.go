// Copyright 2023 N4-Networks.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cntlr

import (
	"log"

	"github.com/n4-networks/openusp/pkg/db"
	"github.com/n4-networks/openusp/pkg/mtp"
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
	cfg    cntlrCfg
	dbH    db.UspDb
	mtpH   mtp.MtpHandler
	cacheH cacheHandler
	agentH agentHandler
	cntlrgrpc.UnimplementedGrpcServer
}

func (c *Cntlr) Init() error {

	// Initialize Logger
	log.SetPrefix("OpenUsp: ")
	log.SetFlags(log.Lshortfile)

	// Load config from env
	if err := c.loadConfigFromEnv(); err != nil {
		log.Println("Error in loading controller config")
		return err
	}

	// Initialize DB connection
	if err := c.dbInit(); err != nil {
		log.Println("Error in dbInit()")
		return err
	}
	log.Println("Db Init ...successful!")

	if err := c.mtpH.Init(); err != nil {
		log.Println("Error in MTP Init")
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

	log.Println("Starting MTP threads...")
	c.mtpH.MtpRxThreads()

	log.Println("Starting GRPC Server...")
	go c.GrpcServerThread(c.cfg.grpc.port, exit)

	log.Println("Starting Cntlr MTP Rx Msg Handler thread...")
	go c.MtpRxMessageHandler()

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
