package mtp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/n4-networks/usp/pkg/pb/bbf/usp_msg"
	piondtls "github.com/pion/dtls/v2"
	coap "github.com/plgd-dev/go-coap/v2"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/client"
)

type coapMsgData struct {
	confirm  bool
	uriQuery string
	uriHost  string
	uriPort  string
	uriPath  string
	pdu      []byte
}

type agentCoap struct {
	addr        string
	port        string
	path        string
	isEncrypted string

	selfUriQuery *message.Option

	conn   *client.ClientConn
	msgCnt uint64
}

func (c agentCoap) sendMsg(msg []byte) error {
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
	resp, err := c.conn.Post(ctx, c.path, message.AppOctets, r, *c.selfUriQuery)
	if err != nil {
		log.Printf("Cannot get response: %v", err)
		return err
	}
	log.Printf("POST Response: %v", resp.String())
	return nil
}
func (c agentCoap) getMsgCnt() uint64 {
	return c.msgCnt
}
func (c agentCoap) incMsgCnt() {
	c.msgCnt++ //TODO: use lock here
}

func (m *Mtp) coapReceiveHandler(w mux.ResponseWriter, req *mux.Message) {
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

	// Parse USP Record
	log.Println("Parsing USP record")
	rData, err := parseUspRecord(cData.pdu)
	if err != nil {
		log.Println("Error in parsing the USP record")
		return
	}
	agentId := rData.fromId
	log.Println("Agent EndpointId: ", agentId)

	log.Println("Validating USP record parameters")
	if err := m.validateUspRecord(rData); err != nil {
		log.Println("Error in validating Rx USP record")
		return
	}
	log.Println("Parsing Rx USP Msg")
	mData, err := parseUspMsg(rData)
	if err != nil {
		log.Println("Error in parsing the USP message")
		return
	}
	aCoap, _ := getAgentInfoCoap(cData)
	if mData.mType == usp_msg.Header_NOTIFY {
		log.Println("Rxed Notify message")
		if mData.notify == nil {
			log.Println("mData.notify is nil")
			return
		}
		if mData.notify.nType == NotifyEvent && mData.notify.evt.name == "Boot!" {
			log.Println("Received Boot event from agent")

			initData := &agentInitData{}
			initData.epId = agentId
			initData.mtpIntf = aCoap
			params, _ := strToMapWithTwoDelims(mData.notify.evt.params["ParameterMap"], ",", ":")
			initData.params = params
			go m.agentInitThread(initData)
			//m.agentH.rxChannel <- initData
		}
		if mData.notify.sendResp {
			uspMsg, err := prepareUspMsgNotifyRes(agentId, mData)
			if err != nil {
				log.Println("could not prepare notify response, err:", err)
				return
			}
			log.Println("Sending Notify response to agent...")
			if err := m.sendUspMsgToAgent(agentId, uspMsg, aCoap); err != nil {
				log.Println("Error in sending USP record, err:", err)
				return
			}
			log.Println("Sent Notify message to agent:", agentId)
		}
	}
	// Non notify messages to be handled here
	if err := m.processRxUspMsg(agentId, mData); err != nil {
		log.Println("Error in processing Rx USP msg, err:", err)
	}
}

func getAgentInfoCoap(cData *coapMsgData) (*agentCoap, error) {
	aCoap := &agentCoap{}
	aCoap.conn = nil
	u, err := url.Parse(cData.uriQuery)
	if err != nil {
		log.Println("Could not parse URIQuery:", cData.uriQuery)
		return nil, err
	}
	aCoap.addr = u.Host
	aCoap.path = u.Path

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

func (m *Mtp) CoAPServerStart(exit chan int32) error {
	switch m.Cfg.CoAP.Server.Mode {
	case "nondtls":
		go m.CoAPServer(exit)
	case "dtls":
		go m.CoAPServerDTLS(exit)
	case "both":
		go m.CoAPServer(exit)
		go m.CoAPServerDTLS(exit)
	}
	return nil
}

// Middleware function, which will be called for each request.
func loggingMiddleware(next mux.Handler) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		log.Printf("ClientAddress %v, %v\n", w.Client().RemoteAddr(), r.String())
		next.ServeCOAP(w, r)
	})
}

func (m *Mtp) CoAPServer(exit chan int32) error {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.Handle("/a", mux.HandlerFunc(m.coapReceiveHandler))
	r.Handle("/b", mux.HandlerFunc(handleB))

	addr := ":" + m.Cfg.CoAP.Server.Port
	log.Println("Starting CoAP server at:", addr)
	log.Fatal(coap.ListenAndServe("udp", addr, r))

	log.Fatalf("CoAP Server is exiting...")
	exit <- COAP_SERVER
	return errors.New("Error in strating CoAP server")
}

func (m *Mtp) CoAPServerDTLS(exit chan int32) {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)
	r.Handle("/a", mux.HandlerFunc(m.coapReceiveHandler))
	//r.Handle("/b", r.HandlerFunc(handleB))

	dtlsConfig := piondtls.Config{
		PSK: func(hint []byte) ([]byte, error) {
			fmt.Printf("Client's hint: %s\n", hint)
			return []byte{0xAB, 0xC1, 0x23}, nil
		},
		PSKIdentityHint: []byte("N4 DTLS client"),
		CipherSuites:    []piondtls.CipherSuiteID{piondtls.TLS_PSK_WITH_AES_128_CCM_8},
	}

	addr := ":" + m.Cfg.CoAP.Server.DTLSPort
	log.Println("Starting CoAP server at:", addr)
	log.Fatal(coap.ListenAndServeDTLS("udp", addr, &dtlsConfig, r))

	log.Fatalf("CoAP Server is exiting...")
	exit <- COAP_SERVER_DTLS
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
