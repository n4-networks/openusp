package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type uspData struct {
	epId   string
	path   string
	params map[string]string
}

func (re *Rest) addInstance(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		inst, err := re.addInstanceObj(d)
		httpSendRes(w, inst, err)
	}
}

func (re *Rest) setParams(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		err := re.setParamsObj(d)
		httpSendRes(w, nil, err)
	}
}

func (re *Rest) operateCmd(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		err := re.MtpOperateReq(d.epId, d.path, "none", true, d.params)
		httpSendRes(w, nil, err)
	}
}

func (re *Rest) getDm(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		objs, err := re.getDmObjs(d)
		httpSendRes(w, objs, err)
	}
}

func (re *Rest) updateDm(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		err := re.updateDmObjs(d)
		httpSendRes(w, nil, err)
	}
}

func (re *Rest) updateInstances(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		err := re.updateInstancesObjs(d)
		httpSendRes(w, nil, err)
	}
}

func (re *Rest) deleteInstances(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		err := re.deleteInstancesObjs(d)
		httpSendRes(w, nil, err)
	}
}

func (re *Rest) updateParams(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		err := re.updateParamsObjs(d)
		httpSendRes(w, nil, err)
	}
}

func (re *Rest) getAgents(w http.ResponseWriter, r *http.Request) {
	log.Println("inside of getAgents api")
	objs, err := re.getAgentIds()
	httpSendRes(w, objs, err)
}

func (re *Rest) deleteDbColl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var ok bool
	var collName string

	if collName, ok = vars["coll"]; !ok {
		log.Println("Collection name not found in the request")
		httpSendRes(w, nil, errors.New("Collection name not found in the request"))
	}
	log.Println("Collection to be deleted:", collName)
	err := re.dbDeleteColl(collName)
	httpSendRes(w, nil, err)
}

func (re *Rest) getMtpInfo(w http.ResponseWriter, r *http.Request) {
	obj, err := re.getMtpInfoObj()
	httpSendRes(w, obj, err)
}

func (re *Rest) reconnectDb(w http.ResponseWriter, r *http.Request) {
	err := re.connectDb()
	httpSendRes(w, nil, err)
}

func (re *Rest) reconnectMtp(w http.ResponseWriter, r *http.Request) {
	err := re.connectMtp()
	httpSendRes(w, nil, err)
}

func (re *Rest) getParams(w http.ResponseWriter, r *http.Request) {
	d, err := parseUspReq(r)
	if err != nil {
		httpSendRes(w, nil, err)
	}
	objs, err := re.getMultipleObjParams(d)
	httpSendRes(w, objs, err)
}

func (re *Rest) getInstances(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		objs, err := re.getInstanceObjs(d.epId, d.path)
		httpSendRes(w, objs, err)
	}
}

func parseUspReq(r *http.Request) (*uspData, error) {
	vars := mux.Vars(r)
	var ok bool
	usp := &uspData{}

	if usp.epId, ok = vars["epId"]; !ok {
		log.Println("EpId not found in the request")
		return nil, errors.New("EpId not found in the request")
	}
	log.Println("getDm EpId:", usp.epId)

	if usp.path, ok = vars["path"]; !ok {
		log.Println("Path not found in the request")
		return nil, errors.New("Path not found in the request")
	}
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&usp.params)
		log.Printf("params:%+v\n", usp.params)
	}
	return usp, nil
}

func httpSendRes(w http.ResponseWriter, objs interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Access-Control-Allow-Origin", "*")  // require for UI to avoid CORS Policy
	//w.Header().Set("Access-Control-Allow-Headers", "*") // require for UI to avoid CORS Policy

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if objs != nil {
		if err := json.NewEncoder(w).Encode(objs); err != nil {
			log.Println("Json Encoder error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
