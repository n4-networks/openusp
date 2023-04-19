package rest

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

func (re *Rest) initRouter() error {
	re.router = mux.NewRouter()
	re.setMiddlewares()
	//re.setStaticHandler()
	re.setRoutesHandlers()
	return nil
}

func (re *Rest) setRoutesHandlers() error {
	log.Println("Setting routing handlers")
	re.router.HandleFunc(GET_DM+"{epId}/{path}", re.getDm).Methods("GET")
	re.router.HandleFunc(GET_PARAMS+"{epId}/{path}", re.getParams).Methods("GET")
	re.router.HandleFunc(GET_INSTANCES+"{epId}/{path}", re.getInstances).Methods("GET")
	re.router.HandleFunc(UPDATE_DM+"{epId}/{path}", re.updateDm).Methods("GET")
	re.router.HandleFunc(UPDATE_INSTANCES+"{epId}/{path}", re.updateInstances).Methods("GET")
	re.router.HandleFunc(DELETE_INSTANCES+"{epId}/{path}", re.deleteInstances).Methods("GET")
	re.router.HandleFunc(UPDATE_PARAMS+"{epId}/{path}", re.updateParams).Methods("GET")
	re.router.HandleFunc(GET_AGENTS, re.getAgents).Methods("GET")
	re.router.HandleFunc(GET_MTPINFO, re.getMtpInfo).Methods("GET")
	re.router.HandleFunc(DELETE_DBCOLL+"{coll}", re.deleteDbColl).Methods("GET")
	re.router.HandleFunc(RECONNECT_DB, re.reconnectDb).Methods("GET")
	re.router.HandleFunc(RECONNECT_MTP, re.reconnectMtp).Methods("GET")

	re.router.HandleFunc(ADD_INSTANCES+"{epId}/{path}", re.addInstance).Methods("POST")
	re.router.HandleFunc(OPERATE_CMD+"{epId}/{path}", re.operateCmd).Methods("POST")
	re.router.HandleFunc(SET_PARAMS+"{epId}/{path}", re.setParams).Methods("POST")

	//router.HandleFunc("/network/{epId}/{type}", re.getNetwork).Methods("GET")

	return nil
}
