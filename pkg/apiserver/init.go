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
	"io/ioutil"
	"log"

	"github.com/gorilla/mux"
	"github.com/n4-networks/openusp/pkg/db"
	"github.com/n4-networks/openusp/pkg/pb/cntlrgrpc"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

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
