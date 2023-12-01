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
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttCfg struct {
	mode       string
	serverAddr string
	topic      string
	userName   string
	passwd     string
}

var mCfg mqttCfg

type MtpMqtt struct {
	Client mqtt.Client
	Topic  string
	MsgCnt uint64
}

func (m *MtpMqtt) configFromEnv() error {
	if env, ok := os.LookupEnv("MQTT_MODE"); ok {
		mCfg.mode = env
	} else {
		log.Println("MQTT mode is not set")
	}

	if env, ok := os.LookupEnv("MQTT_ADDR"); ok {
		mCfg.serverAddr = env
	} else {
		log.Println("MQTT Server Addr is not set")
		return errors.New("MQTT is not set")
	}

	if env, ok := os.LookupEnv("MQTT_TOPIC"); ok {
		mCfg.topic = env
	} else {
		log.Println("MQTT Queue is not set")
		return errors.New("MQTT_QUEUE is not set")
	}

	if env, ok := os.LookupEnv("MQTT_USER"); ok {
		mCfg.userName = env
	} else {
		log.Println("MQTT User Name is not set")
		return errors.New("MQTT_USER is not set")
	}

	if env, ok := os.LookupEnv("MQTT_PASSWD"); ok {
		mCfg.passwd = env
	} else {
		log.Println("MQTT Password is not set")
	}
	log.Printf("MQTT Config params: %+v\n", mCfg)
	return nil
}

func (s *MtpMqtt) SendMsg(msg []byte) error {
	log.Println("Mqtt publishing message to topic:", s.Topic)
	token := s.Client.Publish(s.Topic, 0, false, msg)
	token.Wait()
	return token.Error()
}

func (s *MtpMqtt) GetMsgCnt() uint64 {
	return s.MsgCnt
}

func (s *MtpMqtt) IncMsgCnt() {
	s.MsgCnt++
}

func (s *MtpMqtt) SetParam(name string, value string) error {
	if name == "SubscribedTopic" {
		s.Topic = value
	}
	return nil
}

var publishHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Pub TOPIC: %s\n", msg.Topic())
	fmt.Printf("Pub MSG: %s\n", msg.Payload())
}

var subcribeHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Sub TOPIC: %s\n", msg.Topic())
	fmt.Printf("Sub MSG: %s\n", msg.Payload())
}

func (s *MtpMqtt) Init() error {

	if err := s.configFromEnv(); err != nil {
		log.Println("Error in loading MQTT config from Env")
		return err
	}
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + mCfg.serverAddr).SetClientID("openusp")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(publishHandler)
	opts.SetPingTimeout(1 * time.Second)
	s.Client = mqtt.NewClient(opts)
	return nil
}

func (s *MtpMqtt) Start() error {
	if token := s.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("Mqtt Connect Error:", token.Error())
		return token.Error()
	}

	var rxMsgHandler mqtt.MessageHandler = s.mqttRxMsgHandler
	log.Println("MQTT subscribing to topic:", mCfg.topic)
	if token := s.Client.Subscribe(mCfg.topic, 0, rxMsgHandler); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		return token.Error()
	}
	return nil
}

func (s *MtpMqtt) mqttRxMsgHandler(mc mqtt.Client, msg mqtt.Message) {
	log.Println("MQTT: Received USP msg from agent")

	rxData := &RxChannelData{}
	rxData.Rec = msg.Payload()
	rxData.MtpType = "mqtt"
	rxData.Mtp = s
	rxChannel <- *rxData
	//	aMqtt.topic = "/usp/endpoint/" + agentId
}
