package mtp

import (
	"log"
	"time"

	"github.com/n4-networks/openusp/pkg/db"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *Mtp) dbInit() error {
	var dbClient *mongo.Client
	dbClient, err := db.Connect(m.Cfg.Db.ServerAddr, m.Cfg.Db.UserName, m.Cfg.Db.Passwd, 3*time.Minute)
	if err != nil {
		return err
	}
	err = m.dbH.Init(dbClient, "usp")
	if err != nil {
		return err
	}
	return nil
}

func (m *Mtp) dbGetInstancesByRegex(epId string, path string) ([]*instance, error) {
	dbInstances, err := m.dbH.GetInstancesByRegex(epId, path)
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

func (m *Mtp) dbGetParamsByRegex(epId string, path string) (map[string]string, error) {
	dbParams, err := m.dbH.GetParamsByRegex(epId, path)
	if err != nil {
		return nil, err
	}
	params := make(map[string]string)
	for _, dbParam := range dbParams {
		params[dbParam.Path] = dbParam.Value
	}
	return params, nil
}

func (m *Mtp) dbWriteInstances(epId string, instances []*instance) error {
	var dbInst db.Instance
	for _, inst := range instances {
		dbInst.EndpointId = epId
		dbInst.UniqueKeys = inst.uniqueKeys
		dbInst.Path = inst.path
		log.Println("Writing instance with path:", inst.path)
		if err := m.dbH.WriteInstanceToDb(dbInst); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (m *Mtp) dbDeleteInstancesByRegex(epId string, path string) error {
	log.Println("Delete Instance affected path:", path)
	if err := m.dbH.DeleteInstancesByRegex(epId, path); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (m *Mtp) dbDeleteInstances(epId string, paths []string) error {
	for _, path := range paths {
		log.Println("Delete Instance affected path:", path)
		if err := m.dbH.DeleteInstanceFromDb(epId, path); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

func (m *Mtp) dbWriteParams(epId string, params []*param) error {
	var dbParam db.Param
	for _, param := range params {
		dbParam.EndpointId = epId
		dbParam.Path = param.path
		//log.Printf("{\n%s : %s\n}\n", param.path, param.value)
		dbParam.Value = param.value
		if err := m.dbH.WriteParamToDb(&dbParam); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (m *Mtp) dbDeleteParams(epId string, paths []string) error {
	for _, path := range paths {
		log.Println("Affected path:", path)
		if err := m.dbH.DeleteParamManyFromDb(epId, path); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

func (m *Mtp) dbWriteDatamodels(dmObjs []*db.DmObject) error {
	for _, obj := range dmObjs {
		if err := m.dbH.WriteDmObjectToDb(obj); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (m *Mtp) dbDeleteDatamodels(epId string, paths []string) error {
	for _, path := range paths {
		log.Println("Deleting datamodel, given path:", path)
		if err := m.dbH.DeleteParamManyFromDb(epId, path); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

func (m *Mtp) dbGetCfgInstances(dev *agentDeviceInfo) ([]*cfgInstance, error) {
	dbDev := &db.DevType{
		ProductClass: dev.productClass,
		Manufacturer: dev.manufacturer,
		ModelName:    dev.modelName,
	}

	dbInsts, err := m.dbH.GetCfgInstances(dbDev)
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

func (m *Mtp) dbGetCfgParamNodes(dev *agentDeviceInfo) ([]*cfgParamNode, error) {
	dbDev := &db.DevType{
		ProductClass: dev.productClass,
		Manufacturer: dev.manufacturer,
		ModelName:    dev.modelName,
	}

	dbParamNodes, err := m.dbH.GetCfgParamNodes(dbDev)
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
