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
	"regexp"

	"github.com/n4-networks/openusp/pkg/db"
)

type ObjParam struct {
	Path   string   `json:"path"`
	Params []*Param `json:"params"`
}
type Param struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Access string `json:"access"`
}

func (as *ApiServer) getMultipleObjParams(d *uspData) ([]*ObjParam, error) {
	if as.dbH.uspIntf == nil {
		return nil, errors.New("Error: DB interface has not been initilized")
	}

	dmPath := getDmPathFromAbsPath(d.path)
	log.Println("GetParam, path, dmPath:", d.path, dmPath)
	dm, err := as.dbH.uspIntf.GetDm(d.epId, dmPath)
	if err != nil {
		log.Println("GetDm Err:", err)
		return nil, err
	}

	isAbsPath, err := regexp.MatchString(`.[0-9]\.$`, d.path)
	if err != nil {
		log.Println("regex Err:", err)
	}

	var objs []*ObjParam
	if !dm.MultiInstance || isAbsPath {
		obj := &ObjParam{Path: d.path}
		params, err := as.getSingleObjParams(d.epId, d.path, dm)
		if err != nil {
			log.Println("Err:", err)
			return nil, err
		}
		log.Println("received all the params for:", d.path)
		obj.Params = params
		objs = append(objs, obj)
		return objs, nil
	}

	if dm.MultiInstance {
		log.Println("obj is multi-instance, getting all instances")
		instPath := d.path + "\\d+.$"
		log.Println("InstPath search:", instPath)
		insts, err := as.dbH.uspIntf.GetInstancesByRegex(d.epId, instPath)
		if err != nil {
			log.Println("Err:", err)
			return nil, err
		}
		for _, inst := range insts {
			log.Println("Getting params for:", inst.Path)
			obj := &ObjParam{Path: inst.Path}
			params, err := as.getSingleObjParams(d.epId, inst.Path, dm)
			if err != nil {
				log.Println("Err:", err)
				return nil, err
			}
			log.Println("received all the params for:", inst.Path)
			obj.Params = params
			objs = append(objs, obj)
		}
		return objs, nil
	}
	return nil, errors.New("Invalid path")
}

func (as *ApiServer) getSingleObjParams(epId string, path string, dm *db.DmObject) ([]*Param, error) {
	regexPath := path + "\\w+$" // path + word(paramName) $: end of string
	dbParams, err := as.dbH.uspIntf.GetParamsByRegex(epId, regexPath)
	if err != nil {
		log.Println("GetParamByRegex Err:", err)
		return nil, err
	}
	var params []*Param
	var dmParamPath string
	for _, dmParam := range dm.Params {
		param := &Param{Name: dmParam.Name, Access: dmParam.Access}
		dmParamPath = path + dmParam.Name
		for _, dbParam := range dbParams {
			if dbParam.Path == dmParamPath {
				param.Value = dbParam.Value
				params = append(params, param)
				break
			}
		}
	}
	return params, nil
}

func (as *ApiServer) setParamsObj(d *uspData) error {
	return as.CntlrSetParamReq(d.epId, d.path, d.params)
}

func (as *ApiServer) updateParamsObjs(d *uspData) error {
	return as.CntlrGetParamReq(d.epId, d.path)
}
