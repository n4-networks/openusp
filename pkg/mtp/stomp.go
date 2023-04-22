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
	"github.com/joho/godotenv"
)

type Stomp struct {
	Conn      *stompngo.Connection
	RxChannel <-chan stompngo.MessageData
}

type stompCfg struct {
	mode            string
	serverAddr      string
	serverAddrTLS   string
	agentQueue      string
	controllerQueue string
	userName        string
	passwd          string
	retryCount      int
}

var sCfg stompCfg

type AgentStomp struct {
	Conn      *stompngo.Connection
	DestQueue string
	MsgCnt    uint64
}

func loadStompConfigFromEnv() error {

	if err := godotenv.Load(); err != nil {
		log.Println("Error in loading .env file")
		return err
	}

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

	if env, ok := os.LookupEnv("STOMP_AGENT_QUEUE"); ok {
		sCfg.agentQueue = env
	} else {
		log.Println("STOMP Agent Queue is not set")
		return errors.New("STOMP_AGENT_QUEUE is not set")
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

func (s AgentStomp) SendMsg(msg []byte) error {
	log.Println("Stomp SendMsg is being called")
	h := stompngo.Headers{}
	id := stompngo.Uuid()
	h = h.Add("id", id)
	h = h.Add("destination", s.DestQueue)
	h = h.Add("content-type", "application/vnd.bbf.usp.msg")
	log.Printf("Sending USP record to destination: %v, Success", s.DestQueue)
	return s.Conn.SendBytes(h, msg)
}

func (s AgentStomp) GetMsgCnt() uint64 {
	return s.MsgCnt
}

func (s AgentStomp) IncMsgCnt() {
	s.MsgCnt++
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
func StompInit() (*Stomp, error) {

	if err := loadStompConfigFromEnv(); err != nil {
		log.Println("Error in loading STOMP config from Env")
		return nil, err
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
				return nil, err
			}
		}
	}
	// Create connect headers and connect to STOMP server
	h := stompngo.Headers{}
	conn, err := stompngo.Connect(n, h)
	if err != nil {
		log.Fatalln("Error in connecting to STOMP server: ", err.Error())
		return nil, err
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
		return nil, err
	}
	log.Println("Subscribed to Rx Agent Queue: ", sCfg.controllerQueue)

	s := &Stomp{}
	s.Conn = conn
	s.RxChannel = sub

	return s, nil
}

func StompReceiveThread(s *Stomp, rxChannel chan RxChannelData) {
	for {
		stompMsg := <-s.RxChannel
		rxData := &RxChannelData{}
		rxData.Rec = stompMsg.Message.Body
		rxData.MtpType = "stomp"
		rxChannel <- *rxData
	}
}

/*
func (c *Cntlr) StompReceiveUspMsgFromAgentWithTimer(timer int64) error {
	select {
	case <-time.After(1 * time.Second):
		log.Println("Timeout after 1 second in reading msg, exiting...")
		return errors.New("Timeout after 1 second")
	case stompMsg := <-c.mtpH.stomp.RxChannel:
		rData, err := parseUspRecord(stompMsg.Message.Body)
		if err != nil {
			log.Println("Error in parsing the USP record")
			return err
		}
		agentId := rData.fromId
		log.Println("Rx Agent EndpointId: ", agentId)

		if err := validateUspRecord(rData); err != nil {
			log.Println("Error in validating Rx USP record")
			return err
		}

		mData, err1 := parseUspMsg(rData)
		if err1 != nil {
			log.Println("Error in parsing the USP message")
			return err1
		}

		if err := c.processRxUspMsg(rData.fromId, mData); err != nil {
			log.Println("Could not process Rx Msg, err:", err)
		}
		log.Println("Processed Rx USP MSG")
	}
	return nil
}
*/
