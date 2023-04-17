package mtp

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/n4-networks/openusp/pkg/pb/bbf/usp_msg"
)

type Mqtt struct {
	client mqtt.Client
}

type agentMqtt struct {
	client mqtt.Client
	topic  string
	msgCnt uint64
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

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var subHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("inside of sub TOPIC: %s\n", msg.Topic())
	fmt.Printf("inside of sub MSG: %s\n", msg.Payload())
}

func (m *Mtp) mqttInit() error {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + m.Cfg.Mqtt.BrokerAddr).SetClientID("gotrivial")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Println("Mqtt Connect Error:", token.Error())
		return token.Error()
	}

	var rxMsgHandler mqtt.MessageHandler = m.mqttRxMsgHandler
	log.Println("MQTT subscribing to topic:", m.Cfg.Mqtt.Topic)
	if token := c.Subscribe(m.Cfg.Mqtt.Topic, 0, rxMsgHandler); token.Wait() && token.Error() != nil {
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

	m.connH.mqtt = &Mqtt{}
	m.connH.mqtt.client = c

	return nil
}

func (m *Mtp) mqttRxMsgHandler(c mqtt.Client, msg mqtt.Message) {
	log.Println("MQTT: Received USP msg from agent")
	rData, err := parseUspRecord(msg.Payload())
	if err != nil {
		log.Println("Error in parsing the USP record")
		return
	}
	agentId := rData.fromId
	log.Println("Rx Agent EndpointId: ", agentId)

	if err := m.validateUspRecord(rData); err != nil {
		log.Println("Error in validating Rx USP record")
		return
	}

	mData, err := parseUspMsg(rData)
	if err != nil {
		log.Println("Error in parsing the USP message")
		return
	}
	log.Println("Parsed Rx USP MSG")

	if mData.mType == usp_msg.Header_NOTIFY {
		if mData.notify == nil {
			log.Println("mData.notify is nil")
			return
		}
		aMqtt := &agentMqtt{}
		aMqtt.client = c
		aMqtt.topic = "/usp/endpoint/" + agentId

		if mData.notify.nType == NotifyEvent && mData.notify.evt.name == "Boot!" {
			log.Println("Received Boot event from agent")
			initData := &agentInitData{}
			initData.epId = agentId
			params, _ := strToMapWithTwoDelims(mData.notify.evt.params["ParameterMap"], ",", ":")
			initData.params = params
			initData.mtpIntf = aMqtt
			go m.agentInitThread(initData)

		}
		if mData.notify.sendResp {
			log.Println("Preparing USP Notify Response")
			uspMsg, err := prepareUspMsgNotifyRes(agentId, mData)
			if err != nil {
				log.Println("could not prepare notify response record, err:", err)
				return
			}
			if err := m.sendUspMsgToAgent(agentId, uspMsg, aMqtt); err != nil {
				log.Println("Error in sending USP record, err:", err)
				return
			}
			log.Println("Sent USP Notify message to agent:", agentId)
		}
	}
	// Non notify messages to be handled here
	if err := m.processRxUspMsg(agentId, mData); err != nil {
		log.Println("Error in processing Rx USP msg, err:", err)
	}
}
