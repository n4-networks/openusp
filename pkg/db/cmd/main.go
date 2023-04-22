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

package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/n4-networks/openusp/pkg/db"
	"go.mongodb.org/mongo-driver/mongo"
)

var uspDb db.UspDb
var agentId string = "os::SSGspa-02:42:ac:11:00:06"
var path string = "Device."

func main() {
	log.SetPrefix("N4: ")
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	//redisTest()
	if err := initDb(); err != nil {
		log.Println(err)
		return
	}
	getAllEpIdsTest()
	//cfgInstanceTest()
	//instanceReadByValueTest()
}
func getAllEpIdsTest() {
	epIds, err := uspDb.GetAllEndpoints()
	if err != nil {
		log.Println(err)
		return
	}
	for _, epId := range epIds {
		log.Println("EpId:", epId)
	}
	log.Println("Test of getAllEpIds is completed")
}
func redisTest() {
	addr := "172.17.0.5:6379"
	db.ConnectCache(addr, 5)
}

func cfgInstanceTest() {
	// Write into DB
	inst := &db.CfgInstance{}
	inst.Dev.ProductClass = "MyProductClass"
	inst.Dev.Manufacturer = "MyManufacturer"
	inst.Dev.ModelName = "MyModelName"

	inst.Path = "Device.IP.Interface."
	inst.Params = make(map[string]string)
	inst.Params["mode"] = "dhcp"

	if err := uspDb.WriteCfgInstance(inst); err != nil {
		log.Println(err)
		return
	}
	log.Println("Wrote a cfg instance into db successfully")

	insts, err := uspDb.GetCfgInstances(&inst.Dev)
	if err != nil {
		log.Println(err)
		return
	}
	for _, i := range insts {
		for k, v := range i.Params {
			log.Printf("%v: %v\n", k, v)
		}
	}
}

func instanceReadByValueTest() {
	inst1, err1 := uspDb.GetInstancesByUniqueKeys(agentId, "SSID", "Hyderabad")
	if err1 != nil {
		log.Println("Error in getting instance:", err1)
		return
	}

	log.Printf("instance read, path: %s, Uninique value : %s\n", inst1.Path, inst1.UniqueKeys["SSID"])

}

func initDb() error {
	var dbAddr, dbUser, dbPasswd string
	var ok bool

	if dbAddr, ok = os.LookupEnv("DB_ADDR"); !ok {
		log.Println("DB_ADDR is not set taking default: localhost:27017")
		dbAddr = ":27017"
	}

	if dbUser, ok = os.LookupEnv("DB_USER"); !ok {
		log.Println("DB_USER is not set, returning...")
		return errors.New("DB_USER is not set")
	}

	if dbPasswd, ok = os.LookupEnv("DB_PASSWD"); !ok {
		log.Println("DB_PASSWD is not set, returning...")
		return errors.New("DB_PASSWD is not set")
	}
	dbClient, err := db.Connect(dbAddr, dbUser, dbPasswd, 10*time.Second)
	if err != nil {
		log.Println("Db connect failed, exiting...err: ", err)
		return err
	}
	return uspDb.Init(dbClient, "usp")
}

func instanceTest() {
	log.Printf("\n\n")
	log.Println("Getting all Instance objects")
	log.Println("-----------------------------------------------------")
	instances, _ := uspDb.GetInstances(agentId, path)
	for _, instance := range instances {
		log.Println("{")
		log.Println("  EpID: ", instance.EndpointId)
		log.Println("  Path: ", instance.Path)
		for key, val := range instance.UniqueKeys {
			log.Printf("    Key: %-24s, Value: %12s\n", key, val)
		}
		log.Println("}")
	}
}

func paramTest() {
	log.Printf("\n\n")
	log.Println("Writing a  Parameter objects")
	log.Println("-----------------------------------------------------")
	param := &db.Param{}
	param.EndpointId = agentId
	param.Path = path
	param.Value = "some value"
	uspDb.WriteParamToDb(param)

	log.Printf("\n\n")
	log.Println("Getting all Parameter objects")
	log.Println("-----------------------------------------------------")
	pars, _ := uspDb.GetParams(agentId, path)
	// for _, p := range params {
	for k, v := range pars {
		log.Println("{")
		//log.Println("  EpID: ", p.EndpointId)
		log.Println("  EpID: ", k)
		log.Println("  v: ", v)
		//log.Printf(" %s : %s \n", p.Path, p.Value)
		log.Println("}")
	}

	var params []*db.Param

	log.Printf("\n\n")
	log.Println("Getting Parameters with regex")
	log.Println("-----------------------------------------------------")
	agentId = "os::SSGspa-02:42:ac:11:00:05"
	pattern := "Device\\.WiFi\\.SSID\\.[123]\\.SSID"
	params, _ = uspDb.GetParamsByRegex(agentId, pattern)
	for _, p := range params {
		log.Println("{")
		log.Println("  EpID: ", p.EndpointId)
		log.Printf(" %s : %s \n", p.Path, p.Value)
		log.Println("}")
	}
}

func dmTest() {
	var d db.DmObject

	d.EndpointId = agentId
	d.Path = "Device.WiFi"
	d.MultiInstance = false
	d.Access = "OBJ_READ_ONLY"
	err := uspDb.WriteDmObjectToDb(&d)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Println("Error in writing DM object to DB, err: ", err)
		}
	}

	log.Println("Getting all Datamodel objects for agentId:", agentId)
	log.Println("-----------------------------------------------------")
	dmObjs, _ := uspDb.GetDmByRegex(agentId, path)
	printDmObjects(dmObjs)
	uspDb.DeleteDmObjectManyFromDb(agentId, path)
	log.Println("Data model after deletion")
	dmObjs, _ = uspDb.GetDmByRegex(agentId, path)
	printDmObjects(dmObjs)
}

func printDmObjects(dmObjs []*db.DmObject) {
	for _, obj := range dmObjs {
		log.Println("{")
		log.Println("  EpID: ", obj.EndpointId)
		log.Println("  Path: ", obj.Path)
		log.Println("  AccessType: ", obj.Access)
		log.Println("  MultiInstance: ", obj.MultiInstance)
		for _, param := range obj.Params {
			log.Printf("    Param: %-24s, AccessType : %24s\n", param.Name, param.Access)
		}
		for _, evt := range obj.Events {
			log.Printf("    Event: %-24s, Args: %24s\n", evt.Name, evt.Args)
		}
		for _, cmd := range obj.Cmds {
			log.Printf("    Command: %-24s, Input: %12s, Output: %12s\n", cmd.Name, cmd.Inputs, cmd.Outputs)
		}
		log.Println("}")
	}
}
