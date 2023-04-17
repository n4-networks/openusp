package rest

import (
	"errors"
	"log"

	"github.com/n4-networks/usp/pkg/db"
)

func (re *Rest) getDmObjs(d *uspData) ([]*db.DmObject, error) {
	if re.db.uspIntf == nil {
		return nil, errors.New("Error: DB interface has not been initilized")
	}
	dmObj, err := re.db.uspIntf.GetDmByRegex(d.epId, d.path)
	if err != nil {
		log.Println("Error in getting datamodel from db, err:", err)
		return nil, err
	}
	return dmObj, nil
}

func (re *Rest) updateDmObjs(d *uspData) error {
	if err := re.MtpGetDatamodelReq(d.epId, d.path); err != nil {
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
