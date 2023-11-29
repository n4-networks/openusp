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

type Mqtt struct {
	Client mqtt.Client
}

type mqttCfg struct {
	mode       string
	serverAddr string
	topic      string
	userName   string
	passwd     string
}

var mCfg mqttCfg

type agentMqtt struct {
	client mqtt.Client
	topic  string
	msgCnt uint64
}

func loadMqttConfigFromEnv() error {
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
	log.Printf("STOMP Config params: %+v\n", mCfg)
	return nil
}

func (s agentMqtt) sendMsg(msg []byte) error {
	log.Println("Mqtt publishing message to topic:", s.topic)
	token := s.client.Publish(s.topic, 0, false, msg)
	token.Wait()
	return token.Error()
}

func (s agentMqtt) getMsgCnt() uint64 {
	return s.msgCnt
}

func (s agentMqtt) incMsgCnt() {
	s.msgCnt++
}

var publishHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var subcribeHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("inside of sub TOPIC: %s\n", msg.Topic())
	fmt.Printf("inside of sub MSG: %s\n", msg.Payload())
}

func MqttInit() (mqtt.Client, error) {

	if err := loadMqttConfigFromEnv(); err != nil {
		log.Println("Error in loading MQTT config from Env")
		return nil, err
	}
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + mCfg.serverAddr).SetClientID("gotrivial")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(publishHandler)
	opts.SetPingTimeout(1 * time.Second)
	mc := mqtt.NewClient(opts)
	return mc, nil
}

func MqttStart(mc mqtt.Client) error {
	if token := mc.Connect(); token.Wait() && token.Error() != nil {
		log.Println("Mqtt Connect Error:", token.Error())
		return token.Error()
	}

	var rxMsgHandler mqtt.MessageHandler = mqttRxMsgHandler
	log.Println("MQTT subscribing to topic:", mCfg.topic)
	if token := mc.Subscribe(mCfg.topic, 0, rxMsgHandler); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		return token.Error()
	}

	/*
		for i := 0; i < 5; i++ {
			text := fmt.Sprintf("this is msg #%d!", i)
			token := c.Publish("go-mqtt/sample", 0, false, text)
			token.Wait()
		}

				time.Sleep(6 * time.Second)

				if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
					fmt.Println(token.Error())
					os.Exit(1)
				}
			 //c.Disconnect(250)
			//time.Sleep(1 * time.Second)
	*/

	return nil
}

func mqttRxMsgHandler(mc mqtt.Client, msg mqtt.Message) {
	log.Println("MQTT: Received USP msg from agent")

	rxData := &RxChannelData{}
	rxData.Rec = msg.Payload()
	rxData.MtpType = "mqtt"
	rxC <- *rxData
	//	aMqtt.topic = "/usp/endpoint/" + agentId
}
