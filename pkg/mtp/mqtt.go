package mtp

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Mqtt struct {
	Client mqtt.Client
}

type MqttCfg struct {
	Mode       string `yaml:"mode"`
	BrokerAddr string `yaml:"brokerAddr"`
	Topic      string `yaml:"topic"`
	UserName   string `yaml:"userName"`
	Passwd     string `yaml:"passwd"`
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

var publishHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var subcribeHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("inside of sub TOPIC: %s\n", msg.Topic())
	fmt.Printf("inside of sub MSG: %s\n", msg.Payload())
}

func MqttInit(cfg *MqttCfg) (mqtt.Client, error) {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + cfg.BrokerAddr).SetClientID("gotrivial")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(publishHandler)
	opts.SetPingTimeout(1 * time.Second)
	mc := mqtt.NewClient(opts)
	return mc, nil
}

func MqttStart(mc mqtt.Client, topic string) error {
	if token := mc.Connect(); token.Wait() && token.Error() != nil {
		log.Println("Mqtt Connect Error:", token.Error())
		return token.Error()
	}

	var rxMsgHandler mqtt.MessageHandler = mqttRxMsgHandler
	log.Println("MQTT subscribing to topic:", topic)
	if token := mc.Subscribe(topic, 0, rxMsgHandler); token.Wait() && token.Error() != nil {
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
