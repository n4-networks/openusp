package cntlr

import (
	"log"

	"github.com/n4-networks/openusp/pkg/mtp"
	"github.com/n4-networks/openusp/pkg/pb/bbf/usp_msg"
)

type MtpCfg struct {
	CoAP  mtp.CoAPCfg  `yaml:"coap"`
	Stomp mtp.StompCfg `yaml:"stomp"`
	Mqtt  mtp.MqttCfg  `yaml:"mqtt"`
}

type mtpHandler struct {
	stomp     *mtp.Stomp
	mqtt      *mtp.Mqtt
	coap      *mtp.CoAP
	rxChannel chan mtp.RxChannelData
}

func (c *Cntlr) mtpInit() error {

	c.mtpH.rxChannel = make(chan mtp.RxChannelData, 10)
	mtp.SetRxChannel(c.mtpH.rxChannel)

	// Initialize Stomp client
	stompHandler, err := mtp.StompInit(&c.Cfg.Mtp.Stomp)
	if err != nil {
		log.Println("Error in stompInit()")
		return err
	}
	c.mtpH.stomp = stompHandler

	// Initialize Mqtt client
	mqttClient, err1 := mtp.MqttInit(&c.Cfg.Mtp.Mqtt)
	if err1 != nil {
		log.Println("Error in mqttInit()")
		return err1
	}
	c.mtpH.mqtt = &mtp.Mqtt{Client: mqttClient}
	//c.mtpH.mqtt.Client = mqttClient

	// Initialize  CoAP Server
	coapHandler, err2 := mtp.CoAPServerInit(&c.Cfg.Mtp.CoAP)
	if err2 != nil {
		log.Println("Error in CoapServerInit()")
		return err2
	}
	c.mtpH.coap = coapHandler
	log.Println("Controller MTP has been initialized successfully!")

	return nil
}

func (c *Cntlr) MtpStart() error {
	addr := ":" + c.Cfg.Mtp.CoAP.Server.Port
	go mtp.CoAPServerThread(c.mtpH.coap, addr)

	rxChannel := c.mtpH.rxChannel
	go mtp.StompReceiveThread(c.mtpH.stomp, rxChannel)

	topic := c.Cfg.Mtp.Mqtt.Topic
	if err := mtp.MqttStart(c.mtpH.mqtt.Client, topic); err != nil {
		log.Println("Error in subscribing to Mqtt Topic: ", topic)
	}
	return nil
}

func (c *Cntlr) MtpReceiveThread() {
	for {
		chanData := <-c.mtpH.rxChannel
		log.Println("Shibu: Rx'd USP record from mtp type: ", chanData.MtpType)

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
			aStomp.DestQueue = "/queue/agent-1" //agentId

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
			aStomp.DestQueue = agentId

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