package rest

import (
	"errors"
	"log"
)

type Instance struct {
	Path       string            `json:"path"`
	UniqueKeys map[string]string `json:"unique_keys"`
}

func (re *Rest) getInstanceObjs(epId string, objPath string) ([]*Instance, error) {
	if re.db.uspIntf == nil {
		return nil, errors.New("Error: DB interface has not been initilized")
	}
	dmPath := getDmPathFromAbsPath(objPath)
	dm, err := re.db.uspIntf.GetDm(epId, dmPath)
	if err != nil {
		log.Println("GetDm Err:", err)
		return nil, err
	}
	if !dm.MultiInstance {
		return nil, errors.New("Not a multi instance object")
	}

	regexPath := objPath + "[0-9]+."
	log.Println("GetInstances regexPath:", regexPath)
	dbInsts, err := re.db.uspIntf.GetInstancesByRegex(epId, regexPath)
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

func (re *Rest) addInstanceObj(d *uspData) (*Instance, error) {
	var objs []*object
	obj := &object{}
	obj.path = d.path
	obj.params = d.params
	objs = append(objs, obj)
	insts, err := re.MtpAddInstanceReq(d.epId, objs)
	if err != nil {
		return nil, err
	}
	// TODO: send update request to update parent node in the db (e.g. Number of
	// entries etc.
	return insts[0], nil
}

func (re *Rest) updateInstancesObjs(d *uspData) error {
	return re.MtpGetInstancesReq(d.epId, d.path, false)
}

func (re *Rest) deleteInstancesObjs(d *uspData) error {
	return re.MtpDeleteInstanceReq(d.epId, d.path)
}
