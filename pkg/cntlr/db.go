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
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Cntlr) dbInit() error {
	var dbClient *mongo.Client
	dbClient, err := db.Connect()
	if err != nil {
		return err
	}
	err = c.dbH.Init(dbClient)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cntlr) dbGetInstancesByRegex(epId string, path string) ([]*instance, error) {
	dbInstances, err := c.dbH.GetInstancesByRegex(epId, path)
	if err != nil {
		return nil, err
	}
	var instances []*instance
	for _, dbInst := range dbInstances {
		inst := &instance{}
		inst.path = dbInst.Path
		inst.uniqueKeys = dbInst.UniqueKeys
		instances = append(instances, inst)
	}
	return instances, nil
}

func (c *Cntlr) dbGetParamsByRegex(epId string, path string) (map[string]string, error) {
	dbParams, err := c.dbH.GetParamsByRegex(epId, path)
	if err != nil {
		return nil, err
	}
	params := make(map[string]string)
	for _, dbParam := range dbParams {
		params[dbParam.Path] = dbParam.Value
	}
	return params, nil
}

func (c *Cntlr) dbWriteInstances(epId string, instances []*instance) error {
	var dbInst db.Instance
	for _, inst := range instances {
		dbInst.EndpointId = epId
		dbInst.UniqueKeys = inst.uniqueKeys
		dbInst.Path = inst.path
		log.Println("Writing instance with path:", inst.path)
		if err := c.dbH.WriteInstanceToDb(dbInst); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (c *Cntlr) dbDeleteInstancesByRegex(epId string, path string) error {
	log.Println("Delete Instance affected path:", path)
	if err := c.dbH.DeleteInstancesByRegex(epId, path); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (c *Cntlr) dbDeleteInstances(epId string, paths []string) error {
	for _, path := range paths {
		log.Println("Delete Instance affected path:", path)
		if err := c.dbH.DeleteInstanceFromDb(epId, path); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

func (c *Cntlr) dbWriteParams(epId string, params []*param) error {
	var dbParam db.Param
	for _, param := range params {
		dbParam.EndpointId = epId
		dbParam.Path = param.path
		//log.Printf("{\n%s : %s\n}\n", param.path, param.value)
		dbParam.Value = param.value
		if err := c.dbH.WriteParamToDb(&dbParam); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (c *Cntlr) dbDeleteParams(epId string, paths []string) error {
	for _, path := range paths {
		log.Println("Affected path:", path)
		if err := c.dbH.DeleteParamManyFromDb(epId, path); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

func (c *Cntlr) dbWriteDatamodels(dmObjs []*db.DmObject) error {
	for _, obj := range dmObjs {
		if err := c.dbH.WriteDmObjectToDb(obj); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (c *Cntlr) dbDeleteDatamodels(epId string, paths []string) error {
	for _, path := range paths {
		log.Println("Deleting datamodel, given path:", path)
		if err := c.dbH.DeleteParamManyFromDb(epId, path); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

func (c *Cntlr) dbGetCfgInstances(dev *agentDeviceInfo) ([]*cfgInstance, error) {
	dbDev := &db.DevType{
		ProductClass: dev.productClass,
		Manufacturer: dev.manufacturer,
		ModelName:    dev.modelName,
	}

	dbInsts, err := c.dbH.GetCfgInstances(dbDev)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var insts []*cfgInstance
	for _, i := range dbInsts {
		inst := &cfgInstance{
			path:   i.Path,
			params: i.Params,
			level:  i.Level,
		}
		insts = append(insts, inst)
	}
	return insts, nil
}

func (c *Cntlr) dbGetCfgParamNodes(dev *agentDeviceInfo) ([]*cfgParamNode, error) {
	dbDev := &db.DevType{
		ProductClass: dev.productClass,
		Manufacturer: dev.manufacturer,
		ModelName:    dev.modelName,
	}

	dbParamNodes, err := c.dbH.GetCfgParamNodes(dbDev)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var paramNodes []*cfgParamNode
	for _, i := range dbParamNodes {
		paramNode := &cfgParamNode{
			path:   i.Path,
			params: i.Params,
		}
		paramNodes = append(paramNodes, paramNode)
	}
	return paramNodes, nil
}
