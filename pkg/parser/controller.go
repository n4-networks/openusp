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

package parser

import (
	"errors"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/n4-networks/openusp/pkg/pb/bbf/usp_msg"
)

// ProcessControllerMsgs ...
func ProcessControllerMsgs(byteStream []byte) ([]byte, error) {

	var cntrlrErr error = nil
	var uspMsg usp_msg.Msg

	cntrlrErr = proto.Unmarshal(byteStream, &uspMsg)

	if cntrlrErr != nil {
		return nil, cntrlrErr
	}

	// Write special processing for Error
	if uspMsg.GetHeader().MsgType == usp_msg.Header_ERROR {
		err := checkIntegrityOfUspErrMsg(&uspMsg)

		if err != nil {
			log.Printf("USP Error Process, Error check failed")
			return nil, err
		}
		// Process Error Here
		return nil, nil
		// no further processing required
	}

	// Check for errors
	err := checkIntegrityOfUspRespMsg(&uspMsg)

	if err != nil {
		log.Printf("USP Response Process, Error check failed")
		return nil, err
	}

	switch uspMsg.GetHeader().MsgType {
	case usp_msg.Header_NOTIFY:
		{
			log.Printf("NOTIFY Msg Received ")
		}
	case usp_msg.Header_NOTIFY_RESP:
		{
			log.Printf("NOTIFY_RESP Msg Received ")
		}
	case usp_msg.Header_GET_RESP:
		{
			log.Printf("GET_RESP Msg Received ")
		}
	case usp_msg.Header_SET_RESP:
		{
			log.Printf("SET_RESP Msg Received ")
		}
	case usp_msg.Header_OPERATE_RESP:
		{
			log.Printf("OPERATE_RESP Msg Received ")
		}
	case usp_msg.Header_ADD_RESP:
		{
			log.Printf("ADD_RESP Msg Received ")
		}
	case usp_msg.Header_DELETE_RESP:
		{
			log.Printf("DELETE_RESP Msg Received ")
		}
	case usp_msg.Header_GET_SUPPORTED_DM_RESP:
		{
			log.Printf("GET_SUPPORTED_DM_RESP Msg Received ")
		}
	case usp_msg.Header_GET_INSTANCES_RESP:
		{
			log.Printf("GET_INSTANCES_RESP Msg Received ")
		}
	case usp_msg.Header_GET_SUPPORTED_PROTO_RESP:
		{
			log.Printf("GET_SUPPORTED_PROTO_RESP Msg Received ")
		}
	default:
		log.Printf("Controller cannot process incoming USP msg, unknown/unsupported msg type = %v", uspMsg.GetHeader().MsgType)
	}

	return nil, cntrlrErr
}

func CreateUspOperateReqMsg(cmd string, cmdKey string, resp bool, id string, args map[string]string) ([]byte, error) {

	var req usp_msg.Request
	var operateReq usp_msg.Request_Operate
	var operate usp_msg.Operate

	operate.Command = cmd
	operate.CommandKey = cmdKey
	operate.SendResp = resp
	operate.InputArgs = args

	operateReq.Operate = &operate

	req.ReqType = &operateReq

	return createUspReq(usp_msg.Header_OPERATE, id, &req)
}

func CreateUspGetReqMsg(params []string, id string) ([]byte, error) {

	var req usp_msg.Request
	var getReq usp_msg.Request_Get
	var get usp_msg.Get

	get.ParamPaths = params

	getReq.Get = &get

	req.ReqType = &getReq

	return createUspReq(usp_msg.Header_GET, id, &req)
}

func CreateUspSetReqMsg(objects []*usp_msg.Set_UpdateObject, id string) ([]byte, error) {

	var req usp_msg.Request
	var setReq usp_msg.Request_Set
	var set usp_msg.Set

	set.AllowPartial = true
	set.UpdateObjs = objects

	setReq.Set = &set

	req.ReqType = &setReq

	return createUspReq(usp_msg.Header_SET, id, &req)
}

func CreateUspAddReqMsg(objects []*usp_msg.Add_CreateObject, id string) ([]byte, error) {

	var req usp_msg.Request
	var addReq usp_msg.Request_Add
	var add usp_msg.Add

	add.AllowPartial = false
	add.CreateObjs = objects

	addReq.Add = &add

	req.ReqType = &addReq

	return createUspReq(usp_msg.Header_ADD, id, &req)
}

func CreateUspDeleteReqMsg(objects []string, id string) ([]byte, error) {

	var req usp_msg.Request
	var delReq usp_msg.Request_Delete
	var delete usp_msg.Delete

	delete.AllowPartial = false
	delete.ObjPaths = objects

	delReq.Delete = &delete

	req.ReqType = &delReq

	return createUspReq(usp_msg.Header_DELETE, id, &req)
}

func CreateUspGetSupportedDmMsg(objects []string, retCmd bool, retEvents bool, retParams bool, id string) ([]byte, error) {

	var req usp_msg.Request
	var getSupportedDmReq usp_msg.Request_GetSupportedDm
	var getSupportedDm usp_msg.GetSupportedDM

	getSupportedDm.ObjPaths = objects
	getSupportedDm.ReturnCommands = retCmd
	getSupportedDm.ReturnEvents = retEvents
	getSupportedDm.ReturnParams = retParams

	getSupportedDmReq.GetSupportedDm = &getSupportedDm

	req.ReqType = &getSupportedDmReq

	return createUspReq(usp_msg.Header_GET_SUPPORTED_DM, id, &req)
}

func CreateUspGetInstancesMsg(objects []string, firstLevel bool, id string) ([]byte, error) {

	var req usp_msg.Request
	var getInstancesReq usp_msg.Request_GetInstances
	var getInstances usp_msg.GetInstances

	getInstances.ObjPaths = objects
	getInstances.FirstLevelOnly = firstLevel

	getInstancesReq.GetInstances = &getInstances

	req.ReqType = &getInstancesReq

	return createUspReq(usp_msg.Header_GET_INSTANCES, id, &req)
}

func CreateUspGetSupportedProtoMsg(supportedProtoVersion string, id string) ([]byte, error) {

	var req usp_msg.Request
	var getSupportedProtocolReq usp_msg.Request_GetSupportedProtocol
	var getSupportedProtocol usp_msg.GetSupportedProtocol

	getSupportedProtocol.ControllerSupportedProtocolVersions = supportedProtoVersion

	getSupportedProtocolReq.GetSupportedProtocol = &getSupportedProtocol

	req.ReqType = &getSupportedProtocolReq

	return createUspReq(usp_msg.Header_GET_SUPPORTED_PROTO, id, &req)
}

func checkIntegrityOfUspRespMsg(uspMsg *usp_msg.Msg) error {

	err := checkIntegrityOfUspMsg(uspMsg)

	if err == nil {

		req, ok := uspMsg.GetBody().MsgBody.(*usp_msg.Body_Response)

		if ok && req != nil && req.Response != nil {
			return nil // success
		}
		return errors.New("Error!!!!! reading usp Response, Msg is not present in usp Body")
	}
	return err
}
