package rest

import (
	"context"
	"errors"
	"log"

	"github.com/n4-networks/usp/pkg/db"
)

func (re *Rest) IsConnectedToDb() bool {
	if re.db.client == nil {
		return false
	}
	return true
}

func (re *Rest) connectDb() error {
	if re.db.client != nil {
		ctx, _ := context.WithTimeout(context.Background(), re.cfg.connTimeout)
		re.db.client.Disconnect(ctx)
	}
	/* Connect to DB */
	//log.Println("Connecting to Database @", re.cfg.dbAddr)
	dbClient, err := db.Connect(re.cfg.dbAddr, re.cfg.dbUserName, re.cfg.dbPasswd, re.cfg.connTimeout)
	if err != nil {
		return err
	}
	/* Initialize USP collection connection */
	usp := &db.UspDb{}
	if err := usp.Init(dbClient, "usp"); err != nil {
		return err
	}
	re.db.client = dbClient
	re.db.uspIntf = usp
	log.Println("Connection to DB..SUCCESS")
	return nil
}

func (re *Rest) dbDeleteColl(collName string) error {

	if collName != "datamodel" && collName != "instances" && collName != "params" &&
		collName != "cfginstances" && collName != "cfgparams" {
		log.Println("Invalid db/collection name.", collName)
		return errors.New("Invalid collection name")
	}
	if err := re.db.uspIntf.DeleteCollection(collName); err != nil {
		log.Printf("Error in deleteing db/collection: %v, err: %v\n", collName, err)
		return err
	}
	log.Printf("Db/Collection %v has been removed successfully", collName)
	return nil
}

func (re *Rest) dbGetParamsByRegex(agentId string, path string) (map[string]string, error) {
	//log.Println("Path:", path)
	if re.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbParams, err := re.db.uspIntf.GetParamsByRegex(agentId, path)
	if err != nil {
		return nil, err
	}
	params := make(map[string]string)
	for _, dbParam := range dbParams {
		params[dbParam.Path] = dbParam.Value
	}
	return params, nil
}

func (re *Rest) dbGetParams(agentId string, path string) (map[string]string, error) {
	//log.Println("Path:", path)
	if re.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbParams, err := re.db.uspIntf.GetParams(agentId, path)
	if err != nil {
		return nil, err
	}
	params := make(map[string]string)
	for _, dbParam := range dbParams {
		params[dbParam.Path] = dbParam.Value
	}
	return params, nil
}

func (re *Rest) dbGetDmByRegex(agentId string, path string) ([]*DmObject, error) {
	if re.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbDmObjects, err := re.db.uspIntf.GetDmByRegex(agentId, path)
	if err != nil {
		return nil, err
	}
	var objs []*DmObject
	for _, dbDmObj := range dbDmObjects {
		dmObj := &DmObject{
			Path:          dbDmObj.Path,
			MultiInstance: dbDmObj.MultiInstance,
			Access:        dbDmObj.Access,
		}
		for _, dbParam := range dbDmObj.Params {
			p := dmParam{
				name:   dbParam.Name,
				access: dbParam.Access,
			}
			dmObj.Params = append(dmObj.Params, p)
		}
		for _, dbEvent := range dbDmObj.Events {
			e := dmEvent{}
			e.name = dbEvent.Name
			for _, a := range dbEvent.Args {
				e.args = append(e.args, a)
			}
			dmObj.Events = append(dmObj.Events, e)
		}
		for _, dbCmd := range dbDmObj.Cmds {
			c := dmCmd{}
			c.name = dbCmd.Name
			for _, in := range dbCmd.Inputs {
				c.inputs = append(c.inputs, in)
			}
			for _, out := range dbCmd.Outputs {
				c.outputs = append(c.outputs, out)
			}
			dmObj.Cmds = append(dmObj.Cmds, c)
		}
		objs = append(objs, dmObj)
	}
	return objs, nil
}

func (re *Rest) dbGetDm(agentId string, path string) (*DmObject, error) {
	if re.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbDmObj, err := re.db.uspIntf.GetDm(agentId, path)
	if err != nil {
		return nil, err
	}

	dmObj := &DmObject{
		Path:          dbDmObj.Path,
		MultiInstance: dbDmObj.MultiInstance,
		Access:        dbDmObj.Access,
	}
	for _, dbParam := range dbDmObj.Params {
		p := dmParam{
			name:   dbParam.Name,
			access: dbParam.Access,
		}
		dmObj.Params = append(dmObj.Params, p)
	}
	for _, dbEvent := range dbDmObj.Events {
		e := dmEvent{}
		e.name = dbEvent.Name
		for _, a := range dbEvent.Args {
			e.args = append(e.args, a)
		}
		dmObj.Events = append(dmObj.Events, e)
	}
	for _, dbCmd := range dbDmObj.Cmds {
		c := dmCmd{}
		c.name = dbCmd.Name
		for _, in := range dbCmd.Inputs {
			c.inputs = append(c.inputs, in)
		}
		for _, out := range dbCmd.Outputs {
			c.outputs = append(c.outputs, out)
		}
		dmObj.Cmds = append(dmObj.Cmds, c)
	}
	return dmObj, nil
}

func (re *Rest) dbGetInstancesByRegex(agentId string, path string) ([]*Instance, error) {
	if re.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbInstances, err := re.db.uspIntf.GetInstancesByRegex(agentId, path)
	if err != nil {
		return nil, err
	}
	var instances []*Instance
	for _, dbInst := range dbInstances {
		inst := &Instance{}
		inst.Path = dbInst.Path
		inst.UniqueKeys = dbInst.UniqueKeys
		instances = append(instances, inst)
	}
	return instances, nil
}

func (re *Rest) dbGetInstances(agentId string, path string) ([]*Instance, error) {
	if re.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbInstances, err := re.db.uspIntf.GetInstances(agentId, path)
	if err != nil {
		return nil, err
	}
	var instances []*Instance
	for _, dbInst := range dbInstances {
		inst := &Instance{}
		inst.Path = dbInst.Path
		inst.UniqueKeys = dbInst.UniqueKeys
		instances = append(instances, inst)
	}
	return instances, nil
}

func (re *Rest) dbGetInstanceByAlias(agentId string, aliasName string) (*Instance, error) {
	if re.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbInsts, err := re.db.uspIntf.GetInstancesByUniqueKeys(agentId, "Alias", aliasName)
	if err != nil {
		return nil, err
	}
	inst := &Instance{}
	inst.Path = dbInsts.Path
	inst.UniqueKeys = dbInsts.UniqueKeys
	return inst, nil
}
func (re *Rest) dbDeleteInstances(agentId string, paths []*string) error {
	if re.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	for _, path := range paths {
		log.Println("Affected path:", path)
		if err := re.db.uspIntf.DeleteInstanceFromDb(agentId, *path); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}
func (re *Rest) dbDeleteInstanceByAlias(agentId string, value string) error {
	if re.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	return re.db.uspIntf.DeleteInstanceByUniqueKey(agentId, "Alias", value)
}

type agentInfo struct {
	dev  devInfo
	epId string
}

func (re *Rest) dbWriteCfgInstance(agent agentInfo, path string, level int, key string, params map[string]string) error {
	if re.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	inst := &db.CfgInstance{}
	inst.Dev.ProductClass = agent.dev.productClass
	inst.Dev.Manufacturer = agent.dev.manufacturer
	inst.Dev.ModelName = agent.dev.modelName
	inst.Path = path
	inst.Params = params
	inst.Key = key
	inst.Level = level
	return re.db.uspIntf.WriteCfgInstance(inst)
}

func (re *Rest) dbGetCfgInstancesByPath(agent agentInfo, path string) ([]*cfgInstance, error) {
	if re.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbDevInfo := &db.DevType{
		ProductClass: agent.dev.productClass,
		Manufacturer: agent.dev.manufacturer,
		ModelName:    agent.dev.modelName,
	}
	dbInsts, err := re.db.uspIntf.GetCfgInstancesByPath(dbDevInfo, path)
	if err != nil {
		return nil, err
	}
	var instances []*cfgInstance
	for _, dbInst := range dbInsts {
		inst := &cfgInstance{}
		inst.path = dbInst.Path
		inst.params = dbInst.Params
		inst.level = dbInst.Level
		instances = append(instances, inst)
	}
	return instances, nil
}

func (re *Rest) dbGetCfgInstancesByRegex(agent agentInfo, path string) ([]*cfgInstance, error) {
	if re.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbDevInfo := &db.DevType{
		ProductClass: agent.dev.productClass,
		Manufacturer: agent.dev.manufacturer,
		ModelName:    agent.dev.modelName,
	}
	dbInsts, err := re.db.uspIntf.GetCfgInstancesByRegex(dbDevInfo, path)
	if err != nil {
		return nil, err
	}
	var instances []*cfgInstance
	for _, dbInst := range dbInsts {
		inst := &cfgInstance{}
		inst.path = dbInst.Path
		inst.params = dbInst.Params
		inst.level = dbInst.Level
		inst.key = dbInst.Key
		instances = append(instances, inst)
	}
	return instances, nil
}

func (re *Rest) dbGetCfgParams(agent agentInfo, path string) (map[string]string, error) {
	if re.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbDevInfo := &db.DevType{
		ProductClass: agent.dev.productClass,
		Manufacturer: agent.dev.manufacturer,
		ModelName:    agent.dev.modelName,
	}
	return re.db.uspIntf.GetCfgParams(dbDevInfo, path)
}

func (re *Rest) dbGetCfgParamNodesByRegex(agent agentInfo, path string) ([]*cfgParamNode, error) {
	if re.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbDevInfo := &db.DevType{
		ProductClass: agent.dev.productClass,
		Manufacturer: agent.dev.manufacturer,
		ModelName:    agent.dev.modelName,
	}
	dbCfgParamNodes, err := re.db.uspIntf.GetCfgParamsByRegex(dbDevInfo, path)
	if err != nil {
		return nil, err
	}
	var paramNodes []*cfgParamNode
	for _, dbCfgParamNode := range dbCfgParamNodes {
		paramNode := &cfgParamNode{}
		paramNode.path = dbCfgParamNode.Path
		paramNode.params = dbCfgParamNode.Params
		paramNodes = append(paramNodes, paramNode)
	}
	return paramNodes, nil
}
func (re *Rest) dbWriteCfgParamNode(agent agentInfo, path string, params map[string]string) error {
	if re.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	dbNode := &db.CfgParamNode{}
	dbNode.Dev.ProductClass = agent.dev.productClass
	dbNode.Dev.Manufacturer = agent.dev.manufacturer
	dbNode.Dev.ModelName = agent.dev.modelName
	dbNode.Path = path
	dbNode.Params = params
	return re.db.uspIntf.WriteCfgParamNode(dbNode)
}

func (re *Rest) dbDeleteCfgInstancesByRegex(agent agentInfo, path string) error {
	if re.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	dbDev := &db.DevType{
		ProductClass: agent.dev.productClass,
		Manufacturer: agent.dev.manufacturer,
		ModelName:    agent.dev.modelName,
	}
	if err := re.db.uspIntf.DeleteCfgInstancesByRegex(dbDev, path); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (re *Rest) dbDeleteCfgParamNodesByRegex(agent agentInfo, path string) error {
	if re.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	dbDev := &db.DevType{
		ProductClass: agent.dev.productClass,
		Manufacturer: agent.dev.manufacturer,
		ModelName:    agent.dev.modelName,
	}
	if err := re.db.uspIntf.DeleteCfgParamNodesByRegex(dbDev, path); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (re *Rest) dbDeleteCfgInstanceByKey(agent agentInfo, path string, key string) error {
	if re.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	dbDev := &db.DevType{
		ProductClass: agent.dev.productClass,
		Manufacturer: agent.dev.manufacturer,
		ModelName:    agent.dev.modelName,
	}
	if err := re.db.uspIntf.DeleteCfgInstance(dbDev, path, key); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (re *Rest) dbGetAllEndpoints() ([]string, error) {
	if re.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	epIds, err := re.db.uspIntf.GetAllEndpoints()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return epIds, nil
}
