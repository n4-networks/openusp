package apiserver

import (
	"errors"
	"log"
)

func (as *ApiServer) getAgentIds() ([]string, error) {
	if as.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	agentIds, err := as.db.uspIntf.GetAllEndpoints()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return agentIds, nil
}
