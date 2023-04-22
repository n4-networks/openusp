package cntlr

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	_                  = iota
	SERVER_MODE_NORMAL = iota
	SERVER_MODE_TLS
	SERVER_MODE_NORMAL_AND_TLS
)

type grpcCfg struct {
	port string
}

type cacheCfg struct {
	serverAddr string
}

type uspCfg struct {
	endpointId   string
	protoVersion string
}

type cntlrCfg struct {
	cache cacheCfg
	grpc  grpcCfg
	usp   uspCfg
}

func (c *Cntlr) loadConfigFromEnv() error {

	if err := godotenv.Load(); err != nil {
		log.Println("Error in loading .env file")
		return err
	}

	// Cache config
	if env, ok := os.LookupEnv("CACHE_ADDR"); ok {
		c.cfg.cache.serverAddr = env
	} else {
		log.Println("Cache Server Address is not set")
	}

	// Controller gRPC config
	if env, ok := os.LookupEnv("CNTLR_GRPC_PORT"); ok {
		c.cfg.grpc.port = env
	} else {
		log.Println("CNTRL gRPC Server Port is not set")
	}

	// Controller USP config
	if env, ok := os.LookupEnv("CNTLR_EPID"); ok {
		c.cfg.usp.endpointId = env
	} else {
		log.Println("CNTLR Entpoint ID is not set")
	}
	if env, ok := os.LookupEnv("CNTLR_USP_PROTO_VERSION"); ok {
		c.cfg.usp.protoVersion = env
	} else {
		log.Println("CNTLR USP Protocol Version is not set")
	}

	log.Printf("CNTLR Config params: %+v\n", c.cfg)

	return nil
}
