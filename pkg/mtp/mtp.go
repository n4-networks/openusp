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

var rxChannel chan RxChannelData

type MtpHandler struct {
	StompH    MtpStomp
	MqttH     MtpMqtt
	WsH       MtpWs
	CoapH     MtpCoap
	RxChannel chan RxChannelData
}

func (m *MtpHandler) Init() error {

	rxChannel = make(chan RxChannelData, 100)
	m.RxChannel = rxChannel

	// Initialize Stomp client
	if err := m.StompH.Init(); err != nil {
		log.Println("Error in StompInit()")
		return err
	}

	// Initialize Mqtt client
	if err := m.MqttH.Init(); err != nil {
		log.Println("Error in MqttInit()")
		return err
	}

	// Initialize  CoAP Server
	if err := m.CoapH.Init(); err != nil {
		log.Println("Error in CoapServerInit()")
		return err
	}

	// Initialize WebSocket Server
	if err := m.WsH.Init(); err != nil {
		log.Println("Error in WsServerInit()")
		return err
	}

	log.Println("Controller MTP has been initialized successfully!")

	return nil
}

func (m *MtpHandler) MtpRxThreads() error {
	go m.CoapH.ServerThread()

	go m.StompH.RxThread()

	if err := m.MqttH.Start(); err != nil {
		log.Println("Error in Strating MQTT MTP")
	}

	go m.WsH.ServerThread()
	return nil
}
