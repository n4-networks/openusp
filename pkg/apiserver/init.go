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

package apiserver

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/n4-networks/openusp/pkg/db"
	"github.com/n4-networks/openusp/pkg/pb/cntlrgrpc"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type apiServerCfg struct {
	httpPort    string
	isTlsOn     bool
	cntlrAddr   string
	dbAddr      string
	dbUserName  string
	dbPasswd    string
	connTimeout time.Duration
	logSetting  string
}

type grpcHandle struct {
	intf     cntlrgrpc.GrpcClient
	conn     *grpc.ClientConn
	txMsgCnt uint64
}

func (g *grpcHandle) incTxMsgCnt() uint64 {
	g.txMsgCnt++
	return g.txMsgCnt
}

type dbHandle struct {
	client  *mongo.Client
	uspIntf *db.UspDb
}

type ApiServer struct {
	grpcH  grpcHandle
	dbH    dbHandle
	cfg    apiServerCfg
	router *mux.Router
}

func (as *ApiServer) Init() error {

	log.Println("Running Api Server version:", getVer())

	log.Println("Reading config parameters...")
	if err := as.config(); err != nil {
		log.Println("Could not configure Api Server, err:", err)
		return err
	}

	// Initialize logging
	log.Println("Initializing logging module...")
	if err := as.loggingInit(); err != nil {
		log.Println("Logging settings could not be applied")
	}
	// Connect o Db
	log.Println("Connecting to DB server @", as.cfg.dbAddr)
	if err := as.connectDb(); err != nil {
		log.Println("Error in connecting to DB:", err)
	}

	// Connect to Controller
	log.Println("Connecting to Controller @", as.cfg.cntlrAddr)
	if err := as.connectToController(); err != nil {
		log.Println("Error in connecting to Controller:", err)
	} else {
		log.Println("Connection to Controller...Success")
	}

	// Initialize Router
	if err := as.initRouter(); err != nil {
		log.Println("Error in initializing Router:", err)
	} else {
		log.Println("Initializing Router...Success")
	}
	log.Println("API Server has been initialized")
	return nil
}

func (as *ApiServer) config() error {

	if httpPort, ok := os.LookupEnv("HTTP_PORT"); ok {
		as.cfg.httpPort = httpPort
	} else {
		as.cfg.httpPort = "8080"
	}

	isTlsOn, ok := os.LookupEnv("HTTP_TLS")
	if ok && isTlsOn == "1" {
		as.cfg.isTlsOn = true
	} else {
		as.cfg.isTlsOn = false
	}

	if dbAddr, ok := os.LookupEnv("DB_ADDR"); ok {
		as.cfg.dbAddr = dbAddr
	} else {
		as.cfg.dbAddr = ":27017"
	}

	if dbUserName, ok := os.LookupEnv("DB_USER"); ok {
		as.cfg.dbUserName = dbUserName
	} else {
		log.Println("DB_USER is not set")
		return errors.New("DB_USER not set")
	}

	if dbPasswd, ok := os.LookupEnv("DB_PASSWD"); ok {
		as.cfg.dbPasswd = dbPasswd
	} else {
		log.Println("DB_PASSWD is not set")
		return errors.New("DB_PASSWD not set")
	}

	if cntlrGrpcAddr, ok := os.LookupEnv("CNTLR_GRPC_ADDR"); ok {
		as.cfg.cntlrAddr = cntlrGrpcAddr
	} else {
		as.cfg.cntlrAddr = ":9001"
	}

	as.cfg.connTimeout = 10 * time.Second

	if logging, ok := os.LookupEnv("LOGGING"); ok {
		as.cfg.logSetting = logging
	} else {
		as.cfg.logSetting = "none"
	}
	return nil
}

func (as *ApiServer) loggingInit() error {
	log.SetPrefix("N4: ")
	switch as.cfg.logSetting {
	case "short":
		log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	case "long":
		log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
	case "all":
		log.Println("Setting log for all")
		log.SetFlags(log.Lshortfile | log.Llongfile | log.Ldate | log.Ltime)
	default:
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	return nil
}
