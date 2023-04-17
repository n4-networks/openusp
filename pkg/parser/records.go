package parser

import (
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/n4-networks/openusp/pkg/pb/bbf/usp_record"
)

// This package adds usp secured/plaintext recods communicaion

// CreateNewPlainTextRecord ...
func CreateNewPlainTextRecord(to *string, from *string,
	signature []byte, cert []byte, msgStream []byte) ([]byte, error) {
	return createNewPlainTextRecord(to, from, signature, cert, createUspRecordWithoutContext(msgStream))
}
func createNewPlainTextRecord(to *string, from *string,
	signature []byte, cert []byte, recordContext interface{}) ([]byte, error) {
	var record usp_record.Record

	record.Version = "1.0"
	record.ToId = *to
	record.FromId = *from
	record.MacSignature = signature
	record.SenderCert = cert

	record.PayloadSecurity = usp_record.Record_PLAINTEXT

	if x, ok := recordContext.(*usp_record.Record_NoSessionContext); ok {

		record.RecordType = x
	} else if x, ok := recordContext.(*usp_record.Record_NoSessionContext); ok {

		record.RecordType = x

	} else {

		return nil, errors.New("Invalid session context passed to the createNewPlainTestRecord")
	}

	return proto.Marshal(&record)
}

// createUspRecordWithoutContext ...
func createUspRecordWithoutContext(msgStream []byte) interface{} {
	var noSesCtxRec usp_record.NoSessionContextRecord
	noSesCtxRec.Payload = msgStream
	var noSessionContextRecord usp_record.Record_NoSessionContext
	noSessionContextRecord.NoSessionContext = &noSesCtxRec

	return &noSessionContextRecord
}

// CreateUspRecordWithContext ...
func CreateUspRecordWithContext(msgStream []byte) interface{} {
	//var sessionCtxRec usp_record.SessionContextRecord
	return nil
}

// GetUspMsgStreamFromRecord ...
func GetUspMsgStreamFromRecord(stream []byte) ([]byte, error) {

	var record usp_record.Record

	err := proto.Unmarshal(stream, &record)

	if err != nil {

		return nil, errors.New("Error during Unmarshal of USP Record ")
	}

	if record.PayloadSecurity == usp_record.Record_PLAINTEXT {

		if x, ok := record.RecordType.(*usp_record.Record_NoSessionContext); ok {

			if x.NoSessionContext != nil {
				return x.NoSessionContext.Payload, nil
			}
		}
	}

	return nil, err
}
