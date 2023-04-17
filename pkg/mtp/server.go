package mtp

import (
	"log"
)

const (
	_           = iota
	HTTP_SERVER = iota
	COAP_SERVER
	GRPC_SERVER
	HTTP_SERVER_TLS
	COAP_SERVER_DTLS
	STOMP_SERVER_TLS
)

type connHandler struct {
	stomp *Stomp
	mqtt  *Mqtt
}

func (m *Mtp) Server() (chan int32, error) {

	var exit chan int32
	var err error = nil

	// Start HTTP and HTTPS Server instance based on config
	log.Println("Starting HTTP Server...")
	if err := m.HttpServerStart(exit); err != nil {
		log.Println("Error in starting HTTP server")
	}

	// Start CoAP Server instance based on config
	if err := m.CoAPServerStart(exit); err != nil {
		log.Println("Error in starting HTTP server")
	}

	go m.GrpcServer(m.Cfg.Grpc.Port, exit)

	go m.StompReceiveThread()

	//go m.agentHandlerThread()

	return exit, err
}

func (m *Mtp) ServerWait(exit chan int32) {
	switch <-exit {
	case HTTP_SERVER:
		log.Println("HTTP Server has exited")
	case HTTP_SERVER_TLS:
		log.Println("HTTP Secure Server has exited")
	case COAP_SERVER:
		log.Println("CoAP Server has exited")
	case COAP_SERVER_DTLS:
		log.Println("CoAP Secure Server has exited")
	case STOMP_SERVER_TLS:
		log.Println("STOMP Secure Server has exited")
	case GRPC_SERVER:
		log.Println("GRPC Server has exited")
	default:
		log.Println(" MTP server is existing, err:")
	}
}
