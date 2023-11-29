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
	"errors"
	"log"

	"github.com/n4-networks/openusp/pkg/pb/bbf/usp_msg"
	"github.com/n4-networks/openusp/pkg/pb/bbf/usp_record"
	"google.golang.org/protobuf/proto"
)

type uspRecordData struct {
	msg        *usp_msg.Msg
	version    string
	fromId     string
	toId       string
	recordType string
	destQueue  string
}

func (c *Cntlr) parseUspRecord(s []byte) (*uspRecordData, error) {
	var r usp_record.Record

	if err := proto.Unmarshal(s, &r); err != nil {
		log.Println("Could not unpack byte stream", err)
		return nil, err
	}
	rData := &uspRecordData{}
	rData.version = r.GetVersion()
	rData.toId = r.GetToId()
	rData.fromId = r.GetFromId()
	switch r.RecordType.(type) {
	case *usp_record.Record_StompConnect:
		sc := r.GetStompConnect()
		log.Println("Record Type:", r.GetStompConnect())
		rData.destQueue = sc.GetSubscribedDestination()
		log.Println("Subscribed Destination:", sc.GetSubscribedDestination())
		rData.recordType = "STOMP_CONNECT"
	case *usp_record.Record_WebsocketConnect:
		//sc := r.GetStompConnect()
		log.Println("Record Type:", r.GetWebsocketConnect())
		rData.recordType = "WS_CONNECT"
		//rData.recordType = "STOMP_CONNECT"
	default:
		log.Println("Invalid record type")
	}

	log.Println("Record ToId: ", rData.toId)
	log.Println("Record FromId: ", rData.fromId)

	msg := &usp_msg.Msg{}
	if s := r.GetNoSessionContext(); s != nil {
		log.Println("Record has NoSessionContext")
		if err := proto.Unmarshal(s.GetPayload(), msg); err != nil {
			log.Println("Error in unpacking USP Msg from Record")
			return nil, err
		}
	}
	rData.msg = msg
	return rData, nil
}
func (c *Cntlr) validateUspRecord(rData *uspRecordData) error {
	if c.cfg.usp.protoVersionCheck {
		if rData.version != c.cfg.usp.protoVersion {
			log.Printf("Wrong USP Rx Version: %v, supproted Ver: %v", rData.version, c.cfg.usp.protoVersion)
			return errors.New("USP version mismatch")
		}
	}
	if rData.toId != c.cfg.usp.endpointId {
		log.Printf("Wrong USP Rx ToId: %v, controller Id: %v", rData.toId, c.cfg.usp.endpointId)
		return errors.New("USP ToId/Controller id mismatch")
	}
	log.Printf("Rx Record: USP version: %v, toId: %v", rData.version, rData.toId)
	log.Println("Validated controller Id and USP protocol version")
	return nil
}
