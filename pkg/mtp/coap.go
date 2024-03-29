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
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	coap "github.com/plgd-dev/go-coap/v2"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/client"
)

type coapServerCfg struct {
	mode     string
	port     string
	dtlsPort string
}

type coapClientCfg struct {
	mode           string
	serverPort     string
	serverAddrDTLS string
}

type coapCfg struct {
	server coapServerCfg
	client coapClientCfg
}

var cCfg coapCfg

type coapMsgData struct {
	confirm  bool
	uriQuery string
	uriHost  string
	uriPort  string
	uriPath  string
	pdu      []byte
}

type MtpCoap struct {
	addr        string
	Port        string
	Path        string
	IsEncrypted string

	selfUriQuery *message.Option

	Router *mux.Router
	conn   *client.ClientConn
	MsgCnt uint64
}

func (m *MtpCoap) configFromEnv() error {
	if env, ok := os.LookupEnv("COAP_SERVER_MODE"); ok {
		cCfg.server.mode = env
	} else {
		log.Println("CoAP mode is not set, default is nondtls")
		cCfg.server.mode = "nondtls"
	}

	if env, ok := os.LookupEnv("COAP_SERVER_PORT"); ok {
		cCfg.server.port = env
	} else {
		log.Println("COAP Server Port is not set, default is 5683")
		cCfg.server.port = "5683"
	}

	log.Printf("CoAP Config params: %+v\n", cCfg)
	return nil
}

func (m *MtpCoap) Init() error {

	if err := m.configFromEnv(); err != nil {
		log.Println("Error in loading CoAP config from Env")
		return err
	}

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.Handle("/a", mux.HandlerFunc(m.coapReceiveHandler))
	r.Handle("/b", mux.HandlerFunc(handleB))
	m.Router = r

	return nil
}

func (m *MtpCoap) ServerThread() {

	addr := ":" + cCfg.server.port
	log.Println("Starting CoAP server at:", addr)
	log.Fatal(coap.ListenAndServe("udp", addr, m.Router))

	log.Fatalf("CoAP Server is exiting...")
}

/*
func CoAPServerDTLS(cfg *CoAPCfg, exit chan int32) {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)
	r.Handle("/a", mux.HandlerFunc(coapReceiveHandler))
	//r.Handle("/b", r.HandlerFunc(handleB))

	dtlsConfig := piondtls.Config{
		PSK: func(hint []byte) ([]byte, error) {
			fmt.Printf("Client's hint: %s\n", hint)
			return []byte{0xAB, 0xC1, 0x23}, nil
		},
		PSKIdentityHint: []byte("N4 DTLS client"),
		CipherSuites:    []piondtls.CipherSuiteID{piondtls.TLS_PSK_WITH_AES_128_CCM_8},
	}

	addr := ":" + cfg.Server.DTLSPort
	log.Println("Starting CoAP server at:", addr)
	log.Fatal(coap.ListenAndServeDTLS("udp", addr, &dtlsConfig, r))

	log.Fatalf("CoAP Server is exiting...")
	exit <- COAP_SERVER_DTLS
}
*/

func (c *MtpCoap) SetParam(name string, value string) error {
	return nil
}

func (c *MtpCoap) SendMsg(msg []byte) error {
	var err error
	if c.conn == nil {
		c.conn, err = udp.Dial(c.addr)
		if err != nil {
			log.Println("Error dialing: %v", err)
			return err
		}
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	r := bytes.NewReader(msg)
	log.Println("Sending post msg")
	resp, err := c.conn.Post(ctx, c.Path, message.AppOctets, r, *c.selfUriQuery)
	if err != nil {
		log.Printf("Cannot get response: %v", err)
		return err
	}
	log.Printf("POST Response: %v", resp.String())
	return nil
}
func (m *MtpCoap) GetMsgCnt() uint64 {
	return m.MsgCnt
}
func (m *MtpCoap) IncMsgCnt() {
	m.MsgCnt++ //TODO: use lock here
}

func (m *MtpCoap) coapReceiveHandler(w mux.ResponseWriter, req *mux.Message) {
	log.Println("remote addr:", w.Client().RemoteAddr())
	if req.IsConfirmable {
		log.Println("Sending ACK null msg through setResponse")
		if err := w.SetResponse(codes.Changed, message.TextPlain, nil); err != nil {
			log.Println("Could not send CoAP response, err:", err)
			return
		}
	}
	// Parse CoAP message
	log.Println("Parsing CoAPRxMsg")
	cData, err := parseCoapRxMsg(req)
	if err != nil {
		log.Printf("Could not parse CoAP Msg, err: %v", err)
		return
	}

	rxData := &RxChannelData{}
	rxData.Rec = cData.pdu
	rxData.MtpType = "coap"
	rxData.Mtp = m
	rxChannel <- *rxData

}

func getAgentInfoCoap(cData *coapMsgData) (*MtpCoap, error) {
	aCoap := &MtpCoap{}
	aCoap.conn = nil
	u, err := url.Parse(cData.uriQuery)
	if err != nil {
		log.Println("Could not parse URIQuery:", cData.uriQuery)
		return nil, err
	}
	aCoap.addr = u.Host
	aCoap.Path = u.Path

	opt := &message.Option{}
	opt.ID = message.URIQuery
	uriQuery := "reply-to=coap://" + cData.uriHost + ":" + cData.uriPort + "/" + cData.uriPath
	opt.Value = []byte(uriQuery)
	aCoap.selfUriQuery = opt

	return aCoap, nil
}

func parseCoapRxMsg(req *mux.Message) (*coapMsgData, error) {
	log.Println("Rx Message:", req)
	cData := &coapMsgData{}
	cData.confirm = req.IsConfirmable

	log.Println("Seq Num:", req.SequenceNumber)
	log.Println("Is confirmable:", cData.confirm)

	if uriQuery, err := req.Message.Options.GetString(message.URIQuery); err != nil {
		log.Println("Could not get URI Query")
		return nil, err
	} else {
		str := strings.Split(uriQuery, "reply-to=")
		cData.uriQuery = str[1]
		log.Println("Client Uri-Query:", cData.uriQuery)
	}

	if uriHost, err := req.Message.Options.GetString(message.URIHost); err != nil {
		log.Println("Could not get URI Host")
		return nil, err
	} else {
		cData.uriHost = uriHost
	}

	if uriPort, err := req.Message.Options.GetUint32(message.URIPort); err != nil {
		log.Println("Could not get URI Port")
		return nil, err
	} else {
		cData.uriPort = strconv.Itoa(int(uriPort))
	}

	if uriPath, err := req.Message.Options.GetString(message.URIPath); err != nil {
		log.Println("Could not get URI Path, err:", err)
		return nil, err
	} else {
		cData.uriPath = uriPath
	}

	/*
		size1, err := req.Message.Options.GetUint32(message.Size1)
		if err != nil {
			log.Println("Could not get size1, err:", err)
			return nil, err
		}
		log.Println("Size1:", size1)
	*/

	data, err := ioutil.ReadAll(req.Message.Body)
	if err != nil {
		log.Println("Error in reading request body, err:", err)
		return nil, err
	}
	cData.pdu = data
	log.Println("CoAP Msg Body len:", len(data))
	//log.Println("Data:", string(data))
	return cData, nil
}

// Middleware function, which will be called for each request.
func loggingMiddleware(next mux.Handler) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		log.Printf("ClientAddress %v, %v\n", w.Client().RemoteAddr(), r.String())
		next.ServeCOAP(w, r)
	})
}

func handleB(w mux.ResponseWriter, r *mux.Message) {
	log.Printf("got message in handleB:  %+v from %v\n", r, w.Client().RemoteAddr())
	customResp := message.Message{
		Code:    codes.Content,
		Token:   r.Token,
		Context: r.Context,
		Options: make(message.Options, 0, 16),
		Body:    bytes.NewReader([]byte("B hello world")),
	}
	optsBuf := make([]byte, 32)
	opts, used, err := customResp.Options.SetContentFormat(optsBuf, message.TextPlain)
	if err == message.ErrTooSmall {
		optsBuf = append(optsBuf, make([]byte, used)...)
		opts, used, err = customResp.Options.SetContentFormat(optsBuf, message.TextPlain)
	}
	if err != nil {
		log.Printf("cannot set options to response: %v", err)
		return
	}
	optsBuf = optsBuf[:used]
	customResp.Options = opts

	err = w.Client().WriteMessage(&customResp)
	if err != nil {
		log.Printf("cannot set response: %v", err)
	}
}

func coAPDTLSClient() {
	/*
		co, err := coap.DialDTLS("udp", "localhost:5688", &dtls.Config{
			PSK: func(hint []byte) ([]byte, error) {
				fmt.Printf("Server's hint: %s \n", hint)
				return []byte{0xAB, 0xC1, 0x23}, nil
			},
			PSKIdentityHint: []byte("N4 CoAP DTLS Client"),
			CipherSuites:    []dtls.CipherSuiteID{dtls.TLS_PSK_WITH_AES_128_CCM_8},
		})
		if err != nil {
			log.Fatalf("Error dialing: %v", err)
		}
		path := "/b"
		if len(os.Args) > 1 {
			path = os.Args[1]
		}
		resp, err := co.Get(path)

		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}

		log.Printf("Response payload: %v", string(resp.Payload()))
	*/
}
