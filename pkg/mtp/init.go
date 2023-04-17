package mtp

import (
	"log"

	"github.com/n4-networks/usp/pkg/db"
	"github.com/n4-networks/usp/pkg/pb/mtpgrpc"
)

type Mtp struct {
	Cfg    *Cfg
	dbH    db.UspDb
	connH  connHandler
	cacheH cacheHandler
	agentH agentHandler
	mtpgrpc.UnimplementedMtpGrpcServer
}

func (m *Mtp) Init(confFile string) error {

	// Initialize Logger
	log.SetPrefix("N4: ")
	log.SetFlags(log.Lshortfile)

	var err error

	// Read configuration from yaml file
	err = m.config(confFile)
	if err != nil {
		return err
	}

	// Initialize DB connection
	if err := m.dbInit(); err != nil {
		log.Println("Error in dbInit()")
		return err
	}
	log.Println("Db Init ...successful!")

	// Initialize Stomp client
	if err := m.stompInit(); err != nil {
		log.Println("Error in stompInit()")
		return err
	}

	// Initialize Mqtt client
	if err := m.mqttInit(); err != nil {
		log.Println("Error in mqttInit()")
		return err
	}

	// Initialize Cache handler
	if err := m.cacheInit(); err != nil {
		log.Println("Error in cacheInit()")
		return err
	}

	// Initialize Agent handler
	if err := m.agentHandlerInit(); err != nil {
		return err
	}

	log.Println("MTP has been initialized successfully, Version:", getVer())
	return nil
}
