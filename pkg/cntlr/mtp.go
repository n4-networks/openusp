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

package cntlr

import (
	"log"

	"github.com/n4-networks/openusp/pkg/pb/bbf/usp_msg"
)

func (c *Cntlr) MtpRxMessageHandler() {
	for {
		chanData := <-c.mtpH.RxChannel
		log.Println("Rx'd USP record from mtp type: ", chanData.MtpType)

		rData, err := c.parseUspRecord(chanData.Rec)
		if err != nil {
			log.Println("Error in parsing the USP record")
			continue
		}
		agentId := rData.fromId
		log.Println("Rx Agent EndpointId: ", agentId)

		if err := c.validateUspRecord(rData); err != nil {
			log.Println("Error in validating Rx USP record")
			continue
		}
		if rData.recordType == "STOMP_CONNECT" {
			initData := &agentInitData{}
			initData.epId = agentId
			chanData.Mtp.SetParam("DestQueue", rData.destQueue)
			initData.mtpIntf = chanData.Mtp
			go c.agentInitThread(initData)
			continue

		} else if rData.recordType == "WS_CONNECT" {
			initData := &agentInitData{}
			initData.epId = agentId
			initData.mtpIntf = chanData.Mtp
			go c.agentInitThread(initData)
			continue

		}
		mData, err := parseUspMsg(rData)
		if err != nil {
			log.Println("Error in parsing the USP message")
			continue
		}
		log.Println("Parsed Rx USP MSG")

		if mData.mType == usp_msg.Header_NOTIFY {
			if mData.notify == nil {
				log.Println("mData.notify is nil")
				continue
			}
			chanData.Mtp.SetParam("DestQueue", rData.destQueue)
			if mData.notify.nType == NotifyEvent && mData.notify.evt.name == "Boot!" {
				log.Println("Received Boot event from agent")
				initData := &agentInitData{}
				initData.epId = agentId
				params, _ := strToMapWithTwoDelims(mData.notify.evt.params["ParameterMap"], ",", ":")
				initData.params = params
				initData.mtpIntf = chanData.Mtp
				go c.agentInitThread(initData)
				continue

			}
			if mData.notify.sendResp {
				log.Println("Preparing USP Notify Response")
				uspMsg, err := prepareUspMsgNotifyRes(agentId, mData)
				if err != nil {
					log.Println("could not prepare notify response record, err:", err)
					continue
				}
				if err := c.sendUspMsgToAgent(agentId, uspMsg, chanData.Mtp); err != nil {
					log.Println("Error in sending USP record, err:", err)
					continue
				}
				log.Println("Sent USP Notify message to agent:", agentId)
			}
		}
		// Non notify messages to be handled here
		if err := c.processRxUspMsg(agentId, mData); err != nil {
			log.Println("Error in processing Rx USP msg, err:", err)
		}
	}
}
