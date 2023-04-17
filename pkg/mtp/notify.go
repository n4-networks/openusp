package mtp

import (
	"errors"
	"log"

	"github.com/n4-networks/usp/pkg/parser"
	"github.com/n4-networks/usp/pkg/pb/bbf/usp_msg"
)

func prepareUspMsgNotifyRes(agentId string, mData *uspMsgData) ([]byte, error) {
	if mData == nil {
		log.Println("mData is not initialized")
		return nil, errors.New("mData not initialized")
	}
	log.Println("Preparing Notify response for agent:", agentId)

	var notifyResp usp_msg.NotifyResp
	notifyResp.SubscriptionId = mData.notify.subscriptionId
	uspMsg, err := parser.CreateUspNotifyResponse(&notifyResp, mData.id)
	if err != nil {
		log.Println("Error in creating Usp Msg byte stream, err:", err)
		return nil, err
	}
	return uspMsg, nil
}
