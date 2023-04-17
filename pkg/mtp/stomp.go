package mtp

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo/senv"
	"github.com/n4-networks/usp/pkg/pb/bbf/usp_msg"
)

type Stomp struct {
	Conn      *stompngo.Connection
	RxChannel <-chan stompngo.MessageData
}

type agentStomp struct {
	conn      *stompngo.Connection
	destQueue string
	msgCnt    uint64
}

func (s agentStomp) sendMsg(msg []byte) error {
	log.Println("Stomp SendMsg is being called")
	h := stompngo.Headers{}
	id := stompngo.Uuid()
	h = h.Add("id", id)
	h = h.Add("destination", s.destQueue)
	h = h.Add("content-type", "application/vnd.bbf.usp.msg")
	log.Printf("Sending USP record to destination: %v, Success", s.destQueue)
	return s.conn.SendBytes(h, msg)
}

func (s agentStomp) getMsgCnt() uint64 {
	return s.msgCnt
}

func (s agentStomp) incMsgCnt() {
	s.msgCnt++
}

func connectHeaders() stompngo.Headers {
	h := stompngo.Headers{}
	l := senv.Login()
	if l != "" {
		h = h.Add("login", l)
	}
	pc := senv.Passcode()
	if pc != "" {
		h = h.Add("passcode", pc)
	}
	//
	p := senv.Protocol()
	if p != stompngo.SPL_10 { // 1.1 and 1.2
		h = h.Add("accept-version", p).Add("host", senv.Vhost())
		hb := senv.Heartbeats()
		if hb != "" {
			h = h.Add("heart-beat", hb)
		}
	}
	return h
}
func (m *Mtp) stompInit() error {
	cfg := m.Cfg.Stomp
	var d net.Dialer
	var ctx context.Context
	var n net.Conn
	var err error
	for i := 0; i <= cfg.RetryCount; i++ {
		ctx, _ = context.WithTimeout(context.Background(), time.Minute)
		n, err = d.DialContext(ctx, "tcp", cfg.ServerAddr)
		if err != nil {
			if i < cfg.RetryCount {
				log.Printf("Connection STOMP Server failed, retrying (%v of %v)\n", i, cfg.RetryCount)
			} else {
				log.Println("Connection to STOMP Server failed, exiting: ", err.Error())
				return err
			}
		}
	}
	// Create connect headers and connect to STOMP server
	h := stompngo.Headers{}
	conn, err := stompngo.Connect(n, h)
	if err != nil {
		log.Fatalln("Error in connecting to STOMP server: ", err.Error())
		return err
	}

	// Subcribe to receive queue for msgs coming from Agent
	h1 := stompngo.Headers{}
	//id := stompngo.Uuid()
	//h1 = h.Add("id", id)
	h1 = h.Add("destination", cfg.ControllerQueue)
	sub, err := conn.Subscribe(h1)
	if err != nil {
		log.Fatalf("Could not subscribe to: %v, Err: %v: ", cfg.ControllerQueue, err.Error())
		h := stompngo.Headers{"noreceipt", "true"} // no receipt
		conn.Disconnect(h)
		return err
	}
	log.Println("Subscribed to Rx Agent Queue: ", cfg.ControllerQueue)

	s := &Stomp{}
	s.Conn = conn
	s.RxChannel = sub

	m.connH.stomp = s
	return nil
}

func (m *Mtp) StompReceiveThread() {
	for {
		stompMsg := <-m.connH.stomp.RxChannel

		rData, err := parseUspRecord(stompMsg.Message.Body)
		if err != nil {
			log.Println("Error in parsing the USP record")
			continue
		}
		agentId := rData.fromId
		log.Println("Rx Agent EndpointId: ", agentId)

		if err := m.validateUspRecord(rData); err != nil {
			log.Println("Error in validating Rx USP record")
			continue
		}
		if rData.recordType == "STOMP_CONNECT" {
			aStomp := &agentStomp{}
			aStomp.conn = m.connH.stomp.Conn
			aStomp.destQueue = "/queue/agent-1" //agentId

			initData := &agentInitData{}
			initData.epId = agentId
			//params, _ := strToMapWithTwoDelims(mData.notify.evt.params["ParameterMap"], ",", ":")
			//initData.params = params
			initData.mtpIntf = aStomp
			go m.agentInitThread(initData)
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
			aStomp := &agentStomp{}
			aStomp.conn = m.connH.stomp.Conn
			aStomp.destQueue = agentId

			if mData.notify.nType == NotifyEvent && mData.notify.evt.name == "Boot!" {
				log.Println("Received Boot event from agent")
				initData := &agentInitData{}
				initData.epId = agentId
				params, _ := strToMapWithTwoDelims(mData.notify.evt.params["ParameterMap"], ",", ":")
				initData.params = params
				initData.mtpIntf = aStomp
				go m.agentInitThread(initData)

			}
			if mData.notify.sendResp {
				log.Println("Preparing USP Notify Response")
				uspMsg, err := prepareUspMsgNotifyRes(agentId, mData)
				if err != nil {
					log.Println("could not prepare notify response record, err:", err)
					continue
				}
				if err := m.sendUspMsgToAgent(agentId, uspMsg, aStomp); err != nil {
					log.Println("Error in sending USP record, err:", err)
					continue
				}
				log.Println("Sent USP Notify message to agent:", agentId)
			}
		}
		// Non notify messages to be handled here
		if err := m.processRxUspMsg(agentId, mData); err != nil {
			log.Println("Error in processing Rx USP msg, err:", err)
		}
	}
}

func (m *Mtp) StompReceiveUspMsgFromAgentWithTimer(timer int64) error {
	select {
	case <-time.After(1 * time.Second):
		log.Println("Timeout after 1 second in reading msg, exiting...")
		return errors.New("Timeout after 1 second")
	case stompMsg := <-m.connH.stomp.RxChannel:
		rData, err := parseUspRecord(stompMsg.Message.Body)
		if err != nil {
			log.Println("Error in parsing the USP record")
			return err
		}
		agentId := rData.fromId
		log.Println("Rx Agent EndpointId: ", agentId)

		if err := m.validateUspRecord(rData); err != nil {
			log.Println("Error in validating Rx USP record")
			return err
		}

		mData, err1 := parseUspMsg(rData)
		if err1 != nil {
			log.Println("Error in parsing the USP message")
			return err1
		}

		if err := m.processRxUspMsg(rData.fromId, mData); err != nil {
			log.Println("Could not process Rx Msg, err:", err)
		}
		log.Println("Processed Rx USP MSG")
	}
	return nil
}
