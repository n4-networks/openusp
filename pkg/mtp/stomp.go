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
	"context"
	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo/senv"
)

type stompCfg struct {
	mode            string
	serverAddr      string
	serverAddrTLS   string
	controllerQueue string
	userName        string
	passwd          string
	retryCount      int
}

var sCfg stompCfg

type MtpStomp struct {
	Conn      *stompngo.Connection
	DestQueue string
	MsgCnt    uint64
	RxChannel <-chan stompngo.MessageData
}

func (m *MtpStomp) configFromEnv() error {
	if env, ok := os.LookupEnv("STOMP_MODE"); ok {
		sCfg.mode = env
	} else {
		log.Println("STOMP mode is not set")
	}

	if env, ok := os.LookupEnv("STOMP_ADDR"); ok {
		sCfg.serverAddr = env
	} else {
		log.Println("STOMP Server Addr is not set")
		return errors.New("STOMP_ADDR is not set")
	}

	if env, ok := os.LookupEnv("STOMP_CNTLR_QUEUE"); ok {
		sCfg.controllerQueue = env
	} else {
		log.Println("STOMP Controller Queue is not set")
		return errors.New("STOMP_CNTLR_QUEUE is not set")
	}

	if env, ok := os.LookupEnv("STOMP_USER"); ok {
		sCfg.userName = env
	} else {
		log.Println("STOMP User Name is not set")
		return errors.New("STOMP_USER is not set")
	}

	if env, ok := os.LookupEnv("STOMP_PASSWD"); ok {
		sCfg.passwd = env
	} else {
		log.Println("STOMP Password is not set")
	}

	if env, ok := os.LookupEnv("STOMP_CONN_RETRY"); ok {
		x, _ := strconv.ParseInt(env, 10, 0)
		sCfg.retryCount = int(x)
	} else {
		log.Println("STOMP Connection retry count is not set, default is 5")
		sCfg.retryCount = 5
	}

	log.Printf("STOMP Config params: %+v\n", sCfg)
	return nil

}

func (m *MtpStomp) SendMsg(msg []byte) error {
	log.Println("Stomp SendMsg is being called")
	h := stompngo.Headers{}
	id := stompngo.Uuid()
	h = h.Add("id", id)
	h = h.Add("destination", m.DestQueue)
	log.Printf("Stomp destination: %v", m.DestQueue)
	h = h.Add("content-type", "application/vnd.bbf.usp.msg")
	return m.Conn.SendBytes(h, msg)
}

func (m *MtpStomp) GetMsgCnt() uint64 {
	return m.MsgCnt
}

func (m *MtpStomp) IncMsgCnt() {
	m.MsgCnt++
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
func (m *MtpStomp) Init() error {

	if err := m.configFromEnv(); err != nil {
		log.Println("Error in loading STOMP config from Env")
		return err
	}

	var d net.Dialer
	var ctx context.Context
	var n net.Conn
	var err error
	for i := 0; i <= sCfg.retryCount; i++ {
		ctx, _ = context.WithTimeout(context.Background(), time.Minute)
		n, err = d.DialContext(ctx, "tcp", sCfg.serverAddr)
		if err != nil {
			if i < sCfg.retryCount {
				log.Printf("Connection STOMP Server failed, retrying (%v of %v)\n", i, sCfg.retryCount)
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
	h1 = h.Add("destination", sCfg.controllerQueue)
	sub, err := conn.Subscribe(h1)
	if err != nil {
		log.Fatalf("Could not subscribe to: %v, Err: %v: ", sCfg.controllerQueue, err.Error())
		h := stompngo.Headers{"noreceipt", "true"} // no receipt
		conn.Disconnect(h)
		return err
	}
	log.Println("Subscribed to Rx Agent Queue: ", sCfg.controllerQueue)

	m.Conn = conn
	m.RxChannel = sub

	return nil
}

func (m *MtpStomp) SetParam(name string, value string) error {
	if name == "DestQueue" {
		m.DestQueue = value
	}
	return nil
}
func (m *MtpStomp) RxThread() {
	log.Println("Starting Stomp MTP thread")
	for {
		stompMsg := <-m.RxChannel
		log.Println("Stomp rx msg")
		rxData := &RxChannelData{}
		rxData.Rec = stompMsg.Message.Body
		rxData.MtpType = "stomp"
		rxData.Mtp = m
		rxChannel <- *rxData
	}
}
