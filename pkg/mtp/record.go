package mtp

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
}

func parseUspRecord(s []byte) (*uspRecordData, error) {
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
		log.Println("Record Type:", r.GetStompConnect())
		rData.recordType = "STOMP_CONNECT"
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
func (m *Mtp) validateUspRecord(rData *uspRecordData) error {
	if rData.version != m.Cfg.Usp.ProtoVersion {
		log.Printf("Wrong USP Rx Version: %v, supproted Ver: ", rData.version, m.Cfg.Usp.ProtoVersion)
		return errors.New("USP version mismatch")
	}
	if rData.toId != m.Cfg.Usp.EndpointId {
		log.Printf("Wrong USP Rx ToId: %v, controller Id: %v", rData.toId, m.Cfg.Usp.EndpointId)
		return errors.New("USP ToId/Controller id mismatch")
	}
	log.Printf("Record: USP version: %v, toId: %v", rData.version, rData.toId)
	log.Println("Validated controller Id and USP protocol version")
	return nil
}