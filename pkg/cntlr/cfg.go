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

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
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
