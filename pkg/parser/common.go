package parser

import (
	"errors"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/n4-networks/openusp/pkg/pb/bbf/usp_msg"
)

func checkIntegrityOfUspErrMsg(uspMsg *usp_msg.Msg) error {

	err := checkIntegrityOfUspMsg(uspMsg)

	if err == nil {

		errMsg, ok := uspMsg.GetBody().MsgBody.(*usp_msg.Body_Error)

		if ok && errMsg != nil && errMsg.Error != nil {
			return nil // success but received error from peer
		}
		return errors.New("Error!!!!! reading usp Err Msg, Msg is not present in usp Body")
	}
	return err
}

func checkIntegrityOfUspReqMsg(uspMsg *usp_msg.Msg) error {

	err := checkIntegrityOfUspMsg(uspMsg)
	if err == nil {
		// This is important
		req, ok := uspMsg.GetBody().MsgBody.(*usp_msg.Body_Request)
		if ok && req != nil && req.Request != nil {
			return nil // success
		}

		return errors.New("Error!!!!! reading usp Req, Msg is not present in usp Body")
	}
	return err
}

func createUspReq(reqType usp_msg.Header_MsgType, id string, request *usp_msg.Request) ([]byte, error) {

	var uspMsg usp_msg.Msg
	var header usp_msg.Header
	var body usp_msg.Body
	var reqBody usp_msg.Body_Request

	uspMsg.Header = &header
	uspMsg.Body = &body

	header.MsgType = reqType
	header.MsgId = id

	body.MsgBody = &reqBody
	reqBody.Request = request

	msg, err := proto.Marshal(&uspMsg)

	if err != nil {

		log.Printf("Marshall USP request Failed")
		return nil, err
	}
	return msg, nil
}

func checkIntegrityOfUspMsg(uspMsg *usp_msg.Msg) error {

	// check if USP Msg Header is present
	if uspMsg.GetHeader() == nil {

		const err = "Error!!! getting USP Header can not process further"
		log.Printf(err)
		return errors.New(err)
	}

	// Check if USP Msg Body is present
	if uspMsg.GetBody() == nil {

		return errors.New("Error!!! reading USP Msg Body, Not present")
	}
	return nil // success
}

func CreateUspNotifyResponse(notifyResp *usp_msg.NotifyResp, id string) ([]byte, error) {

	var resp usp_msg.Response
	var notResp usp_msg.Response_NotifyResp

	notResp.NotifyResp = notifyResp

	resp.RespType = &notResp

	return createUspResp(usp_msg.Header_NOTIFY_RESP, id, &resp)
}
