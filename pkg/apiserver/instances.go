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
	"log"
)

type Instance struct {
	Path       string            `json:"path"`
	UniqueKeys map[string]string `json:"unique_keys"`
}

func (as *ApiServer) getInstanceObjs(epId string, objPath string) ([]*Instance, error) {
	if as.dbH.uspIntf == nil {
		return nil, errors.New("Error: DB interface has not been initilized")
	}
	dmPath := getDmPathFromAbsPath(objPath)
	dm, err := as.dbH.uspIntf.GetDm(epId, dmPath)
	if err != nil {
		log.Println("GetDm Err:", err)
		return nil, err
	}
	if !dm.MultiInstance {
		return nil, errors.New("Not a multi instance object")
	}

	regexPath := objPath + "[0-9]+."
	log.Println("GetInstances regexPath:", regexPath)
	dbInsts, err := as.dbH.uspIntf.GetInstancesByRegex(epId, regexPath)
	if err != nil {
		log.Println("GetInstances from DB failed", err)
		return nil, err
	}
	var insts []*Instance
	for _, dbInst := range dbInsts {
		inst := &Instance{
			Path:       dbInst.Path,
			UniqueKeys: dbInst.UniqueKeys,
		}
		insts = append(insts, inst)
	}
	return insts, nil
}

func (as *ApiServer) addInstanceObj(d *uspData) (*Instance, error) {
	var objs []*object
	obj := &object{}
	obj.path = d.path
	obj.params = d.params
	objs = append(objs, obj)
	insts, err := as.CntlrAddInstanceReq(d.epId, objs)
	if err != nil {
		return nil, err
	}
	// TODO: send update request to update parent node in the db (e.g. Number of
	// entries etc.
	return insts[0], nil
}

func (as *ApiServer) updateInstancesObjs(d *uspData) error {
	return as.CntlrGetInstancesReq(d.epId, d.path, false)
}

func (as *ApiServer) deleteInstancesObjs(d *uspData) error {
	return as.CntlrDeleteInstanceReq(d.epId, d.path)
}
