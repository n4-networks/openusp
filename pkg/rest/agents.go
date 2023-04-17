package rest

import (
	"errors"
	"log"
)

func (re *Rest) getAgentIds() ([]string, error) {
	if re.db.uspIntf == nil {
		return nil, errors.New("Not connected to DB")
	}
	agentIds, err := re.db.uspIntf.GetAllEndpoints()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return agentIds, nil
}
