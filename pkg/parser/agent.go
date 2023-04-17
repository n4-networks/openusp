package parser

import (
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/n4-networks/openusp/pkg/pb/bbf/usp_msg"
)

// ProcessAgentMsgsHandler ...
func ProcessAgentMsgsHandler(byteStream []byte, handler func(usp_msg.Msg) ([]byte, error)) ([]byte, error) {
	var agentErr error = nil
	var uspMsg usp_msg.Msg

	agentErr = proto.Unmarshal(byteStream, &uspMsg)

	if agentErr != nil {
		return nil, agentErr
	}

	// Write special processing for Error
	if uspMsg.GetHeader().MsgType == usp_msg.Header_ERROR {
		err := checkIntegrityOfUspErrMsg(&uspMsg)

		if err != nil {
			log.Printf("USP Error Process, Error check failed")
			return nil, err
		}
		log.Printf("ERROR Msg Received")
		// Process Error Here
		return nil, nil
		// no further processing required
	}

	// Check for errors
	err := checkIntegrityOfUspReqMsg(&uspMsg)

	if err != nil {
		log.Printf("USP Request Process, Error check failed")
		return nil, err
	}

	return handler(uspMsg)
}

func defaultHandler(uspMsg usp_msg.Msg) ([]byte, error) {
	switch uspMsg.GetHeader().MsgType {

	case usp_msg.Header_GET:
		{
			log.Printf("GET Msg Received")
		}
	case usp_msg.Header_NOTIFY:
		{
			log.Printf("NOTIFY Msg Received ")
		}
	case usp_msg.Header_SET:
		{
			log.Printf("SET Msg Received")
		}
	case usp_msg.Header_OPERATE:
		{
			log.Printf("OPERATE Msg Received")
			// respond if response needed
		}
	case usp_msg.Header_ADD:
		{
			log.Printf("ADD Msg Received ")
		}
	case usp_msg.Header_DELETE:
		{
			log.Printf("DELETE Msg Received")
		}
	case usp_msg.Header_GET_SUPPORTED_DM:
		{
			log.Printf("GETSUPPORTED_DM Msg Received")
		}
	case usp_msg.Header_GET_INSTANCES:
		{
			log.Printf("GET INSTANCES Msg Received")
		}
	case usp_msg.Header_NOTIFY_RESP:
		{
			log.Printf("NOTIFY RESP Msg Received")
		}
	case usp_msg.Header_GET_SUPPORTED_PROTO:
		{
			log.Printf("GET SUPPORTED PROTO Msg Received")
		}
	default:
		log.Printf("Agent cannot process incoming USP msg, unknown/unsupported msg type = %v ", uspMsg.GetHeader().MsgType)
	}

	return nil, nil
}

// ProcessAgentMsgs ... This function is depricated
func ProcessAgentMsgs(byteStream []byte) ([]byte, error) {

	return ProcessAgentMsgsHandler(byteStream, defaultHandler)
}

// CreateUspNotifyEventMsg ...
func CreateUspNotifyEventMsg(eventNotification *usp_msg.Notify_Event_, sendResp bool, subscriptionID string, id string) ([]byte, error) {

	var notify usp_msg.Notify

	notify.Notification = eventNotification
	notify.SendResp = sendResp
	notify.SubscriptionId = subscriptionID

	return createNotifyMsg(&notify, id)
}

// CreateUspNotifyValueChangeMsg ...
func CreateUspNotifyValueChangeMsg(valChangeNotify *usp_msg.Notify_ValueChange_, sendResp bool, subscriptionID string, id string) ([]byte, error) {

	var notify usp_msg.Notify

	notify.Notification = valChangeNotify
	notify.SendResp = sendResp
	notify.SubscriptionId = subscriptionID

	return createNotifyMsg(&notify, id)
}

// CreateUspNotifyObjectCreationMsg ...
func CreateUspNotifyObjectCreationMsg(objCreateNotify *usp_msg.Notify_ObjCreation, sendResp bool, subscriptionID string, id string) ([]byte, error) {

	var notify usp_msg.Notify

	notify.Notification = objCreateNotify
	notify.SendResp = sendResp
	notify.SubscriptionId = subscriptionID

	return createNotifyMsg(&notify, id)
}

// CreateUspNotifyObjectDeletionMsg ...
func CreateUspNotifyObjectDeletionMsg(objDeleteNotify *usp_msg.Notify_ObjDeletion, sendResp bool, subscriptionID string, id string) ([]byte, error) {

	var notify usp_msg.Notify

	notify.Notification = objDeleteNotify
	notify.SendResp = sendResp
	notify.SubscriptionId = subscriptionID

	return createNotifyMsg(&notify, id)
}

// CreateUspNotifyOperationCompleteMsg ...
func CreateUspNotifyOperationCompleteMsg(oprCompleteNotify *usp_msg.Notify_OperComplete, sendResp bool, subscriptionID string, id string) ([]byte, error) {

	var notify usp_msg.Notify

	notify.Notification = oprCompleteNotify
	notify.SendResp = sendResp
	notify.SubscriptionId = subscriptionID

	return createNotifyMsg(&notify, id)
}

// CreateUspNotifyOnboardReqMsg ...
func CreateUspNotifyOnboardReqMsg(onBoardRqNotify *usp_msg.Notify_OnBoardReq, sendResp bool, subscriptionID string, id string) ([]byte, error) {

	var notify usp_msg.Notify

	notify.Notification = onBoardRqNotify
	notify.SendResp = sendResp
	notify.SubscriptionId = subscriptionID

	return createNotifyMsg(&notify, id)
}

// CreateErrorMsg ...
func CreateErrorMsg(errorMsg *usp_msg.Error, id string) ([]byte, error) {
	var uspMsg usp_msg.Msg
	var header usp_msg.Header
	var body usp_msg.Body
	var errBody usp_msg.Body_Error

	uspMsg.Header = &header
	uspMsg.Body = &body

	header.MsgType = usp_msg.Header_ERROR
	header.MsgId = id

	body.MsgBody = &errBody

	errBody.Error = errorMsg

	msg, err := proto.Marshal(&uspMsg)

	if err != nil {

		log.Printf("Marshall USP request Failed")
		return nil, err
	}
	return msg, nil
}

func createNotifyMsg(notify *usp_msg.Notify, id string) ([]byte, error) {
	var req usp_msg.Request
	var notifyReq usp_msg.Request_Notify
	notifyReq.Notify = notify
	req.ReqType = &notifyReq

	return createUspReq(usp_msg.Header_NOTIFY, id, &req)
}

// CreateUspOperateRespMsg ...
func CreateUspOperateRespMsg(id string, result []*usp_msg.OperateResp_OperationResult) ([]byte, error) {
	var resp usp_msg.Response
	var oprResp usp_msg.Response_OperateResp
	var operateRsp usp_msg.OperateResp

	operateRsp.OperationResults = result

	resp.RespType = &oprResp

	return createUspResp(usp_msg.Header_OPERATE_RESP, id, &resp)
}

// CreateUspGetRespMsg ...
func CreateUspGetRespMsg(result []*usp_msg.GetResp_RequestedPathResult, id string) ([]byte, error) {

	var resp usp_msg.Response
	var getResp usp_msg.Response_GetResp
	var getRsp usp_msg.GetResp

	//getRsp.ParamPaths = params
	getRsp.ReqPathResults = result

	getResp.GetResp = &getRsp

	resp.RespType = &getResp

	return createUspResp(usp_msg.Header_GET_RESP, id, &resp)
}

// CreateUspSetRespMsg ...
func CreateUspSetRespMsg(id string, result []*usp_msg.SetResp_UpdatedObjectResult) ([]byte, error) {
	var resp usp_msg.Response
	var setResp usp_msg.Response_SetResp
	var setRsp usp_msg.SetResp

	setRsp.UpdatedObjResults = result

	setResp.SetResp = &setRsp

	resp.RespType = &setResp

	return createUspResp(usp_msg.Header_SET_RESP, id, &resp)
}

// CreateUspAddRespMsg ...
func CreateUspAddRespMsg(id string, result []*usp_msg.AddResp_CreatedObjectResult) ([]byte, error) {
	var resp usp_msg.Response
	var addResp usp_msg.Response_AddResp
	var addRsp usp_msg.AddResp

	addRsp.CreatedObjResults = result

	addResp.AddResp = &addRsp

	resp.RespType = &addResp

	return createUspResp(usp_msg.Header_ADD_RESP, id, &resp)
}

// CreateUspDeleteRespMsg ...
func CreateUspDeleteRespMsg(result []*usp_msg.DeleteResp_DeletedObjectResult, id string) ([]byte, error) {

	var resp usp_msg.Response
	var delResp usp_msg.Response_DeleteResp
	var deleteRsp usp_msg.DeleteResp

	deleteRsp.DeletedObjResults = result

	delResp.DeleteResp = &deleteRsp

	resp.RespType = &delResp

	return createUspResp(usp_msg.Header_DELETE_RESP, id, &resp)
}

// CreateUspGetSupportedDmRespMsg ...
func CreateUspGetSupportedDmRespMsg(result []*usp_msg.GetSupportedDMResp_RequestedObjectResult, id string) ([]byte, error) {

	var resp usp_msg.Response
	var getSupportedDmResp usp_msg.Response_GetSupportedDmResp
	var getSupportedDmRsp usp_msg.GetSupportedDMResp

	getSupportedDmRsp.ReqObjResults = result

	getSupportedDmResp.GetSupportedDmResp = &getSupportedDmRsp

	resp.RespType = &getSupportedDmResp

	return createUspResp(usp_msg.Header_GET_SUPPORTED_DM_RESP, id, &resp)
}

// CreateUspGetInstancesRespMsg ...
func CreateUspGetInstancesRespMsg(result []*usp_msg.GetInstancesResp_RequestedPathResult, id string) ([]byte, error) {

	var resp usp_msg.Response
	var getInstancesResp usp_msg.Response_GetInstancesResp
	var getInstancesRsp usp_msg.GetInstancesResp

	getInstancesRsp.ReqPathResults = result

	getInstancesResp.GetInstancesResp = &getInstancesRsp

	resp.RespType = &getInstancesResp

	return createUspResp(usp_msg.Header_GET_INSTANCES_RESP, id, &resp)
}

// CreateUspGetSupportedProtoRespMsg ...
func CreateUspGetSupportedProtoRespMsg(supportedProtoVersion string, id string) ([]byte, error) {

	var resp usp_msg.Response
	var getSupportedProtocolResp usp_msg.Response_GetSupportedProtocolResp
	var getSupportedProtocolRsp usp_msg.GetSupportedProtocolResp

	getSupportedProtocolRsp.AgentSupportedProtocolVersions = supportedProtoVersion

	getSupportedProtocolResp.GetSupportedProtocolResp = &getSupportedProtocolRsp

	resp.RespType = &getSupportedProtocolResp

	return createUspResp(usp_msg.Header_GET_SUPPORTED_PROTO_RESP, id, &resp)
}

func createUspResp(respType usp_msg.Header_MsgType, id string, response *usp_msg.Response) ([]byte, error) {

	var uspMsg usp_msg.Msg
	var header usp_msg.Header
	var body usp_msg.Body
	var respBody usp_msg.Body_Response

	uspMsg.Header = &header
	uspMsg.Body = &body

	header.MsgType = respType
	header.MsgId = id

	body.MsgBody = &respBody
	respBody.Response = response

	msg, err := proto.Marshal(&uspMsg)

	if err != nil {

		log.Printf("Marshall USP request Failed")
		return nil, err
	}
	return msg, nil
}
