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
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type wsCfg struct {
	mode    string
	addr    string
	path    string
	addrTLS string
}

type MtpWs struct {
	Conn   *websocket.Conn
	MsgCnt uint64
}

func (m *MtpWs) SendMsg(data []byte) error {
	return m.Conn.WriteMessage(websocket.BinaryMessage, data)
}

func (m *MtpWs) GetMsgCnt() uint64 {
	return m.MsgCnt
}

func (m *MtpWs) IncMsgCnt() {
	m.MsgCnt++
}

func (m *MtpWs) SetParam(name string, value string) error {
	return nil
}

var wCfg wsCfg

func (m *MtpWs) Init() error {
	if err := m.configFromEnv(); err != nil {
		log.Println("Error in loading WebSocket config from Env")
		return err
	}
	log.Println("Configuring Ws Receive handler with path: ", wCfg.path)
	http.HandleFunc(wCfg.path, m.ReceiveHandler)
	return nil
}

func (m *MtpWs) configFromEnv() error {
	if env, ok := os.LookupEnv("WS_MODE"); ok {
		wCfg.mode = env
	} else {
		log.Println("WebSocket mode is not set")
	}

	if env, ok := os.LookupEnv("WS_PATH"); ok {
		wCfg.path = env
	} else {
		log.Println("WebSocket path is not set")
	}

	if env, ok := os.LookupEnv("WS_SERVER_PORT"); ok {
		wCfg.addr = ":" + env
	} else {
		log.Println("WS Server PORT is not set")
		return errors.New("WS_SERVER_PORT is not set")
	}

	if env, ok := os.LookupEnv("WS_SERVER_TLS_PORT"); ok {
		wCfg.addrTLS = ":" + env
	} else {
		log.Println("WS Server TLS PORT is not set")
		return errors.New("WS_SERVER_TLS_PORT is not set")
	}

	log.Printf("WebSocket Config params: %+v\n", wCfg)
	return nil

}

func (m *MtpWs) ReceiveHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{} // use default options
	header := []string{"v1.usp"}
	conn, err := upgrader.Upgrade(w, r, http.Header{
		"Sec-websocket-Protocol": header})
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	m.Conn = conn
	for {
		if _, message, err := conn.ReadMessage(); err != nil {
			log.Println("WS Read Error:", err)
			return
		} else {
			//log.Println("recv: %s", message)
			rxData := &RxChannelData{}
			rxData.Rec = message
			rxData.MtpType = "ws"
			rxData.Mtp = m
			rxChannel <- *rxData
		}
	}
}

func (m *MtpWs) ServerThread() {
	log.Println("Starting WebSocket server at: ", wCfg.addr)
	http.ListenAndServe(wCfg.addr, nil)

}
