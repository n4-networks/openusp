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

func (as *ApiServer) addInstance(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		inst, err := as.addInstanceObj(d)
		httpSendRes(w, inst, err)
	}
}

func (as *ApiServer) setParams(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		err := as.setParamsObj(d)
		httpSendRes(w, nil, err)
	}
}

func (as *ApiServer) operateCmd(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		err := as.CntlrOperateReq(d.epId, d.path, "none", true, d.params)
		httpSendRes(w, nil, err)
	}
}

func (as *ApiServer) getDm(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		objs, err := as.getDmObjs(d)
		httpSendRes(w, objs, err)
	}
}

func (as *ApiServer) updateDm(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		err := as.updateDmObjs(d)
		httpSendRes(w, nil, err)
	}
}

func (as *ApiServer) updateInstances(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		err := as.updateInstancesObjs(d)
		httpSendRes(w, nil, err)
	}
}

func (as *ApiServer) deleteInstances(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		err := as.deleteInstancesObjs(d)
		httpSendRes(w, nil, err)
	}
}

func (as *ApiServer) updateParams(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		err := as.updateParamsObjs(d)
		httpSendRes(w, nil, err)
	}
}

func (as *ApiServer) getAgents(w http.ResponseWriter, r *http.Request) {
	log.Println("inside of getAgents api")
	objs, err := as.getAgentIds()
	httpSendRes(w, objs, err)
}

func (as *ApiServer) deleteDbColl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var ok bool
	var collName string

	if collName, ok = vars["coll"]; !ok {
		log.Println("Collection name not found in the request")
		httpSendRes(w, nil, errors.New("Collection name not found in the request"))
	}
	log.Println("Collection to be deleted:", collName)
	err := as.dbDeleteColl(collName)
	httpSendRes(w, nil, err)
}

func (as *ApiServer) getCntlrInfo(w http.ResponseWriter, r *http.Request) {
	obj, err := as.getCntlrInfoObj()
	httpSendRes(w, obj, err)
}

func (as *ApiServer) reconnectDb(w http.ResponseWriter, r *http.Request) {
	err := as.connectDb()
	httpSendRes(w, nil, err)
}

func (as *ApiServer) reconnectCntlr(w http.ResponseWriter, r *http.Request) {
	err := as.connectToController()
	httpSendRes(w, nil, err)
}

func (as *ApiServer) getParams(w http.ResponseWriter, r *http.Request) {
	d, err := parseUspReq(r)
	if err != nil {
		httpSendRes(w, nil, err)
	}
	objs, err := as.getMultipleObjParams(d)
	httpSendRes(w, objs, err)
}

func (as *ApiServer) getInstances(w http.ResponseWriter, r *http.Request) {
	if d, err := parseUspReq(r); err != nil {
		httpSendRes(w, nil, err)
	} else {
		objs, err := as.getInstanceObjs(d.epId, d.path)
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
