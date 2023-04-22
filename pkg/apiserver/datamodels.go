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

	"github.com/n4-networks/openusp/pkg/db"
)

func (as *ApiServer) getDmObjs(d *uspData) ([]*db.DmObject, error) {
	if as.dbH.uspIntf == nil {
		return nil, errors.New("Error: DB interface has not been initilized")
	}
	dmObj, err := as.dbH.uspIntf.GetDmByRegex(d.epId, d.path)
	if err != nil {
		log.Println("Error in getting datamodel from db, err:", err)
		return nil, err
	}
	return dmObj, nil
}

func (as *ApiServer) updateDmObjs(d *uspData) error {
	if err := as.CntlrGetDatamodelReq(d.epId, d.path); err != nil {
		log.Println("updateDm error:", err)
		return err
	}
	return nil
}

func printDmObjs(dmObjs []*DmObject) {
	for _, d := range dmObjs {
		log.Printf("path: %-24s, MultiInstance: %v Access: %v\n", d.Path, d.MultiInstance, d.Access)
		log.Printf("Commands:\n")
		for _, cmd := range d.Cmds {
			log.Printf("  %-24s, Input: %12s Output: %12s\n", cmd.name, cmd.inputs, cmd.outputs)
		}
		log.Printf("Events:\n")
		for _, evt := range d.Events {
			log.Printf("  %-24s Args: %24s\n", evt.name, evt.args)
		}
		log.Printf("Params:\n")
		for _, param := range d.Params {
			log.Printf("  %-24s AccessType : %24s\n", param.name, param.access)
		}
	}
}
