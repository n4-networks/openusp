package cli

import (
	"log"

	"github.com/abiosoft/ishell"
)

func (cli *Cli) registerNounsDb() {
	cmds := []noun{
		{"remove", "db", removeDbHelp, cli.removeDbColl},
	}
	cli.registerNouns(cmds)
}

const removeDbHelp = "remove db datamodel|instances|params|cfginstances|cfgparams"

func (cli *Cli) removeDbColl(c *ishell.Context) {
	if len(c.Args) < 1 {
		c.Println("Wrong input.", removeDbHelp)
		return
	}
	collName := c.Args[0]
	if collName != "datamodel" && collName != "instances" && collName != "params" &&
		collName != "cfginstances" && collName != "cfgparams" {
		c.Println("Invalid db/collection name.", removeDbHelp)
		return
	}
	if err := cli.restDeleteCollection(collName); err != nil {
		log.Printf("Error in deleteing db/collection: %v, err: %v\n", collName, err)
		return
	}
	c.Printf("Db/Collection %v has been removed successfully", collName)
	c.Println("-------------------------------------------------\n")
}

/*
func (cli *Cli) dbGetInstanceByAlias(aliasName string) (*Instance, error) {
	if cli.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbInsts, err := cli.db.uspIntf.GetInstancesByUniqueKeys(cli.agent.epId, "Alias", aliasName)
	if err != nil {
		return nil, err
	}
	inst := &Instance{}
	inst.Path = dbInsts.Path
	inst.UniqueKeys = dbInsts.UniqueKeys
	return inst, nil
}
func (cli *Cli) dbDeleteInstances(paths []*string) error {
	if cli.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	for _, path := range paths {
		log.Println("Affected path:", path)
		if err := cli.db.uspIntf.DeleteInstanceFromDb(cli.agent.epId, *path); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}
func (cli *Cli) dbDeleteInstanceByAlias(value string) error {
	if cli.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	return cli.db.uspIntf.DeleteInstanceByUniqueKey(cli.agent.epId, "Alias", value)
}
*/

/*
func (cli *Cli) dbWriteCfgInstance(path string, level int, key string, params map[string]string) error {
	if cli.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	inst := &db.CfgInstance{}
	inst.Dev.ProductClass = cli.agent.dev.productClass
	inst.Dev.Manufacturer = cli.agent.dev.manufacturer
	inst.Dev.ModelName = cli.agent.dev.modelName
	inst.Path = path
	inst.Params = params
	inst.Key = key
	inst.Level = level
	return cli.db.uspIntf.WriteCfgInstance(inst)
}

func (cli *Cli) dbGetCfgInstancesByPath(path string) ([]*cfgInstance, error) {
	if cli.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbDevInfo := &db.DevType{
		ProductClass: cli.agent.dev.productClass,
		Manufacturer: cli.agent.dev.manufacturer,
		ModelName:    cli.agent.dev.modelName,
	}
	dbInsts, err := cli.db.uspIntf.GetCfgInstancesByPath(dbDevInfo, path)
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

func (cli *Cli) dbGetCfgInstancesByRegex(path string) ([]*cfgInstance, error) {
	if cli.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbDevInfo := &db.DevType{
		ProductClass: cli.agent.dev.productClass,
		Manufacturer: cli.agent.dev.manufacturer,
		ModelName:    cli.agent.dev.modelName,
	}
	dbInsts, err := cli.db.uspIntf.GetCfgInstancesByRegex(dbDevInfo, path)
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

func (cli *Cli) dbGetCfgParams(path string) (map[string]string, error) {
	if cli.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbDevInfo := &db.DevType{
		ProductClass: cli.agent.dev.productClass,
		Manufacturer: cli.agent.dev.manufacturer,
		ModelName:    cli.agent.dev.modelName,
	}
	return cli.db.uspIntf.GetCfgParams(dbDevInfo, path)
}

func (cli *Cli) dbGetCfgParamNodesByRegex(path string) ([]*cfgParamNode, error) {
	if cli.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	dbDevInfo := &db.DevType{
		ProductClass: cli.agent.dev.productClass,
		Manufacturer: cli.agent.dev.manufacturer,
		ModelName:    cli.agent.dev.modelName,
	}
	dbCfgParamNodes, err := cli.db.uspIntf.GetCfgParamsByRegex(dbDevInfo, path)
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
func (cli *Cli) dbWriteCfgParamNode(path string, params map[string]string) error {
	if cli.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	dbNode := &db.CfgParamNode{}
	dbNode.Dev.ProductClass = cli.agent.dev.productClass
	dbNode.Dev.Manufacturer = cli.agent.dev.manufacturer
	dbNode.Dev.ModelName = cli.agent.dev.modelName
	dbNode.Path = path
	dbNode.Params = params
	return cli.db.uspIntf.WriteCfgParamNode(dbNode)
}

func (cli *Cli) dbDeleteCfgInstancesByRegex(path string) error {
	if cli.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	dbDev := &db.DevType{
		ProductClass: cli.agent.dev.productClass,
		Manufacturer: cli.agent.dev.manufacturer,
		ModelName:    cli.agent.dev.modelName,
	}
	if err := cli.db.uspIntf.DeleteCfgInstancesByRegex(dbDev, path); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (cli *Cli) dbDeleteCfgParamNodesByRegex(path string) error {
	if cli.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	dbDev := &db.DevType{
		ProductClass: cli.agent.dev.productClass,
		Manufacturer: cli.agent.dev.manufacturer,
		ModelName:    cli.agent.dev.modelName,
	}
	if err := cli.db.uspIntf.DeleteCfgParamNodesByRegex(dbDev, path); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (cli *Cli) dbDeleteCfgInstanceByKey(path string, key string) error {
	if cli.db.uspIntf == nil {
		return errors.New("Not connected to DB")
	}
	dbDev := &db.DevType{
		ProductClass: cli.agent.dev.productClass,
		Manufacturer: cli.agent.dev.manufacturer,
		ModelName:    cli.agent.dev.modelName,
	}
	if err := cli.db.uspIntf.DeleteCfgInstance(dbDev, path, key); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
*/
