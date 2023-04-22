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
	"log"

	"github.com/gorilla/mux"
)

const (
	GET_DM    = "/get/dm/"
	UPDATE_DM = "/update/dm/"

	GET_PARAMS    = "/get/params/"
	SET_PARAMS    = "/set/params/"
	UPDATE_PARAMS = "/update/params/"

	GET_INSTANCES    = "/get/instances/"
	ADD_INSTANCES    = "/add/instances/"
	DELETE_INSTANCES = "/delete/instances/"
	UPDATE_INSTANCES = "/update/instances/"

	OPERATE_CMD = "/operate/cmd/"

	RECONNECT_MTP = "/reconnect/mtp/"
	RECONNECT_DB  = "/reconnect/db/"

	DELETE_DBCOLL = "/delete/dbcoll/"

	GET_AGENTS  = "/get/agents/"
	GET_MTPINFO = "/get/mtpinfo/"
)

func (as *ApiServer) initRouter() error {
	as.router = mux.NewRouter()
	as.setMiddlewares()
	//as.setStaticHandler()
	as.setRoutesHandlers()
	return nil
}

func (as *ApiServer) setRoutesHandlers() error {
	log.Println("Setting routing handlers")
	as.router.HandleFunc(GET_DM+"{epId}/{path}", as.getDm).Methods("GET")
	as.router.HandleFunc(GET_PARAMS+"{epId}/{path}", as.getParams).Methods("GET")
	as.router.HandleFunc(GET_INSTANCES+"{epId}/{path}", as.getInstances).Methods("GET")
	as.router.HandleFunc(UPDATE_DM+"{epId}/{path}", as.updateDm).Methods("GET")
	as.router.HandleFunc(UPDATE_INSTANCES+"{epId}/{path}", as.updateInstances).Methods("GET")
	as.router.HandleFunc(DELETE_INSTANCES+"{epId}/{path}", as.deleteInstances).Methods("GET")
	as.router.HandleFunc(UPDATE_PARAMS+"{epId}/{path}", as.updateParams).Methods("GET")
	as.router.HandleFunc(GET_AGENTS, as.getAgents).Methods("GET")
	as.router.HandleFunc(GET_MTPINFO, as.getCntlrInfo).Methods("GET")
	as.router.HandleFunc(DELETE_DBCOLL+"{coll}", as.deleteDbColl).Methods("GET")
	as.router.HandleFunc(RECONNECT_DB, as.reconnectDb).Methods("GET")
	as.router.HandleFunc(RECONNECT_MTP, as.reconnectCntlr).Methods("GET")

	as.router.HandleFunc(ADD_INSTANCES+"{epId}/{path}", as.addInstance).Methods("POST")
	as.router.HandleFunc(OPERATE_CMD+"{epId}/{path}", as.operateCmd).Methods("POST")
	as.router.HandleFunc(SET_PARAMS+"{epId}/{path}", as.setParams).Methods("POST")

	//router.HandleFunc("/network/{epId}/{type}", as.getNetwork).Methods("GET")

	return nil
}
