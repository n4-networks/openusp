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

package mtp

import (
	"log"
)

type MtpIntf interface {
	SetParam(string, string) error
	SendMsg([]byte) error
	GetMsgCnt() uint64
	IncMsgCnt()
}

type RxChannelData struct {
	Rec     []byte
	AgentId string
	MtpType string
	Mtp     MtpIntf
}

var rxC chan RxChannelData

func SetRxChannel(rxChannel chan RxChannelData) {
	rxC = rxChannel
}

type MtpHandler struct {
	StompH    MtpStomp
	MqttH     *Mqtt
	CoapH     *CoAP
	RxChannel chan RxChannelData
}

func Init() (*MtpHandler, error) {

	mtpH := &MtpHandler{}
	mtpH.RxChannel = make(chan RxChannelData, 10)
	rxC = mtpH.RxChannel

	// Initialize Stomp client
	if err := mtpH.StompH.Init(); err != nil {
		log.Println("Error in stompInit()")
		return nil, err
	}

	// Initialize Mqtt client
	mqttClient, err1 := MqttInit()
	if err1 != nil {
		log.Println("Error in mqttInit()")
		return nil, err1
	}
	mtpH.MqttH = &Mqtt{Client: mqttClient}
	//c.mtpH.mqtt.Client = mqttClient

	// Initialize  CoAP Server
	coapHandler, err2 := CoAPServerInit()
	if err2 != nil {
		log.Println("Error in CoapServerInit()")
		return nil, err2
	}
	mtpH.CoapH = coapHandler

	// Initialize WebSocket Server
	if err3 := WsServerInit(); err3 != nil {
		log.Println("Error in WsServerInit()")
		return nil, err3
	}

	log.Println("Controller MTP has been initialized successfully!")

	return mtpH, nil
}

func MtpRxThreads(mtpH *MtpHandler) error {
	go CoAPServerThread(mtpH.CoapH)

	rxChannel := mtpH.RxChannel
	go mtpH.StompH.ReceiveThread(rxChannel)

	if err := MqttStart(mtpH.MqttH.Client); err != nil {
		log.Println("Error in Strating MQTT MTP")
	}

	go WsServerThread()
	return nil
}

/*

func (c *Cntlr) MtpReceiveThread() {
	for {
		chanData := <-c.mtpH.rxChannel
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
			aStomp := &mtp.AgentStomp{}
			aStomp.Conn = c.mtpH.stomp.Conn
			aStomp.DestQueue = rData.destQueue

			initData := &agentInitData{}
			initData.epId = agentId
			//params, _ := strToMapWithTwoDelims(mData.notify.evt.params["ParameterMap"], ",", ":")
			//initData.params = params
			initData.mtpIntf = aStomp
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
			aStomp := &mtp.AgentStomp{}
			aStomp.Conn = c.mtpH.stomp.Conn
			aStomp.DestQueue = rData.destQueue

			if mData.notify.nType == NotifyEvent && mData.notify.evt.name == "Boot!" {
				log.Println("Received Boot event from agent")
				initData := &agentInitData{}
				initData.epId = agentId
				params, _ := strToMapWithTwoDelims(mData.notify.evt.params["ParameterMap"], ",", ":")
				initData.params = params
				initData.mtpIntf = aStomp
				go c.agentInitThread(initData)

			}
			if mData.notify.sendResp {
				log.Println("Preparing USP Notify Response")
				uspMsg, err := prepareUspMsgNotifyRes(agentId, mData)
				if err != nil {
					log.Println("could not prepare notify response record, err:", err)
					continue
				}
				if err := c.sendUspMsgToAgent(agentId, uspMsg, aStomp); err != nil {
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
*/
